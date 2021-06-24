package tools

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"arhat.dev/pkg/exechelper"
	"arhat.dev/pkg/log"
	"github.com/fatih/color"
	"golang.org/x/term"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/output"
)

// ToolType for interface type registration
var ToolType = reflect.TypeOf((*Tool)(nil)).Elem()

// nolint:revive
type Tool interface {
	field.Interface

	// Kind of the tool, e.g. golang, docker
	ToolKind() string

	ToolName() string

	Init(
		cacheDir string,
		rf field.RenderingFunc,
		getBaseExecSpec field.ExecSpecGetFunc,
	) error

	ResolveTasks(tasks []Task) error

	Run(ctx context.Context, taskKind, taskName string) error
}

type BaseTool struct {
	field.BaseField

	Name string   `yaml:"name"`
	Path string   `yaml:"path"`
	Env  []string `yaml:"env"`

	Args []string `yaml:"args"`

	cacheDir             string                `json:"-" yaml:"-"`
	defaultExecutable    string                `json:"-" yaml:"-"`
	RenderingFunc        field.RenderingFunc   `json:"-" yaml:"-"`
	getBootstrapExecSpec field.ExecSpecGetFunc `json:"-" yaml:"-"`

	stdoutIsTty bool `json:"-" yaml:"-"`
	stderrIsTty bool `json:"-" yaml:"-"`
}

func (t *BaseTool) InitBaseTool(
	cacheDir string,
	defaultExecutable string,
	rf field.RenderingFunc,
	getBaseExecSpec field.ExecSpecGetFunc,
) error {
	t.cacheDir = cacheDir
	t.defaultExecutable = defaultExecutable
	t.RenderingFunc = rf
	t.getBootstrapExecSpec = getBaseExecSpec

	t.stdoutIsTty = term.IsTerminal(int(os.Stdout.Fd()))
	t.stderrIsTty = term.IsTerminal(int(os.Stderr.Fd()))

	return nil
}

func (t *BaseTool) ToolName() string { return t.Name }

func (t *BaseTool) RunTask(ctx context.Context, toolKind string, task Task) error {
	execCtx, cancelExec := context.WithCancel(ctx)
	defer cancelExec()

	baseCtx := field.WithRenderingValues(execCtx, t.Env)

	matrixSpecs, err := task.GetMatrixSpecs(baseCtx, t.RenderingFunc)
	if err != nil {
		return fmt.Errorf("failed to create build matrix: %w", err)
	}

	if len(matrixSpecs) == 0 {
		return fmt.Errorf("no matrix spec match")
	}

	workerCount := constant.GetWorkerCount(ctx)
	if workerCount > len(matrixSpecs) {
		workerCount = len(matrixSpecs)
	}

	waitCh := make(chan struct{}, workerCount)
	for i := 0; i < workerCount; i++ {
		waitCh <- struct{}{}
	}

	type taskResult struct {
		matrixSpec MatrixSpec
		err        error
	}

	var (
		errCollection []taskResult

		resultMU = &sync.Mutex{}
	)

	appendResult := func(spec MatrixSpec, err error) {
		resultMU.Lock()
		defer resultMU.Unlock()

		errCollection = append(errCollection, taskResult{
			matrixSpec: spec,
			err:        err,
		})
	}

	var colorList = [][2]*color.Color{
		{color.New(color.FgHiCyan), color.New(color.FgCyan)},
		{color.New(color.FgHiGreen), color.New(color.FgGreen)},
		{color.New(color.FgHiMagenta), color.New(color.FgMagenta)},
		{color.New(color.FgHiYellow), color.New(color.FgYellow)},
		{color.New(color.FgHiBlue), color.New(color.FgBlue)},
		{color.New(color.FgHiRed), color.New(color.FgRed)},
	}

	failFast := constant.IsFailFast(baseCtx.Context())

	wg := &sync.WaitGroup{}
	for i, ms := range matrixSpecs {
		err2 := func() error {
			color := colorList[i%len(colorList)]
			prefixColor, outputColor := color[0], color[1]

			taskCtx := baseCtx.Clone()

			select {
			case <-taskCtx.Context().Done():
				return taskCtx.Context().Err()
			case <-waitCh:
			}

			output.WriteTaskStart(
				taskCtx.Context(),
				task.ToolKind(), task.ToolName(), task.TaskKind(), task.TaskName(),
				ms.String(),
			)

			for k, v := range ms {
				taskCtx.AddEnv("MATRIX_" + strings.ToUpper(k) + "=" + v)
			}

			err = task.ResolveFields(taskCtx, t.RenderingFunc, -1)
			if err != nil {
				return fmt.Errorf("failed to resolve task fields: %w", err)
			}

			var toolCmd []string
			if len(t.Path) != 0 {
				toolCmd = append(toolCmd, t.Path)
			} else {
				toolCmd = append(toolCmd, t.defaultExecutable)
			}

			toolCmd = append(toolCmd, t.Args...)

			execSpecs, err := task.GetExecSpecs(taskCtx, toolCmd)
			if err != nil {
				return fmt.Errorf("failed to generate task args: %w", err)
			}

			wg.Add(1)
			go func(ms MatrixSpec) {
				defer func() {
					wg.Done()

					select {
					case waitCh <- struct{}{}:
					case <-taskCtx.Context().Done():
						return
					}
				}()

				err := t.doRunTask(
					taskCtx,
					fmt.Sprint("{", ms.String(), "}: "),
					prefixColor, outputColor,
					execSpecs,
				)
				output.WriteExecResult(
					taskCtx.Context(),
					task.ToolKind(), task.ToolName(), task.TaskKind(), task.TaskName(),
					ms.String(),
					err,
				)

				if err != nil {
					if failFast {
						cancelExec()
					}

					appendResult(ms, err)
				}
			}(ms)

			return nil
		}()

		if err2 != nil {
			// failed before execution
			if failFast {
				return fmt.Errorf("failed to prepare task execution: %w", err2)
			}
		}
	}

	wg.Wait()

	if len(errCollection) != 0 {
		return fmt.Errorf("task execution failed")
	}

	return nil
}

func (t *BaseTool) doRunTask(
	taskCtx *field.RenderingContext,
	outputPrefix string,
	prefixColor, outputColor *color.Color,
	execSpecs []TaskExecSpec,
) error {
	for _, es := range execSpecs {
		ctx := taskCtx.Clone()

		ctx.AddEnv(es.Env...)

		_, runScriptCmd, err := t.getBootstrapExecSpec(strings.Join(es.Command, " "), false)
		if err != nil {
			return fmt.Errorf("failed to get exec spec from bootstrap config: %w", err)
		}

		output.WriteExecStart(
			ctx.Context(),
			t.ToolName(),
			es.Command,
			filepath.Base(runScriptCmd[len(runScriptCmd)-1]),
		)

		var (
			stdout io.Writer
			stderr io.Writer
		)

		if t.stderrIsTty {
			stderr = output.PrefixWriter(outputPrefix, prefixColor, outputColor, os.Stderr)
		} else {
			stderr = output.PrefixWriter(outputPrefix, nil, nil, os.Stderr)
		}

		if t.stdoutIsTty {
			stdout = output.PrefixWriter(outputPrefix, prefixColor, outputColor, os.Stdout)
		} else {
			stdout = output.PrefixWriter(outputPrefix, nil, nil, os.Stdout)
		}

		p, err := exechelper.Do(exechelper.Spec{
			Context: ctx.Context(),
			Command: runScriptCmd,
			Env:     ctx.Values().Env,

			Stdin: os.Stdin,

			Stdout: stdout,
			Stderr: stderr,
		})

		if err != nil {
			if !es.IgnoreError {
				return fmt.Errorf("failed to execute command [ %s ]: %w", strings.Join(es.Command, " "), err)
			}

			// TODO: log error in detail
			log.Log.I("error ignored", log.Error(err))
		}

		code, err := p.Wait()
		if err != nil {
			if !es.IgnoreError {
				return fmt.Errorf("command exited with code %d: %w", code, err)
			}

			// TODO: log error in detail
			log.Log.I("error ignored", log.Error(err))
		}
	}

	return nil
}

// GetExecSpec is a helper func for shell renderer
func (t *BaseTool) GetExecSpec(script string, isFilePath bool) (env, cmd []string, err error) {
	scriptPath := script
	if !isFilePath {
		scriptPath, err = GetScriptCache(t.cacheDir, script)
		if err != nil {
			return nil, nil, fmt.Errorf("tools: failed to ensure script cache: %w", err)
		}
	}

	if len(t.Path) != 0 {
		cmd = append(cmd, t.Path)
	} else {
		cmd = append(cmd, t.ToolName())
	}

	cmd = append(cmd, t.Args...)

	return t.Env, append(cmd, scriptPath), nil
}
