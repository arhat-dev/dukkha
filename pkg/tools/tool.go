package tools

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"arhat.dev/pkg/exechelper"
	"arhat.dev/pkg/log"
	"github.com/fatih/color"
	"golang.org/x/term"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/output"
	"arhat.dev/dukkha/pkg/sliceutils"
)

// ToolType for interface type registration
var ToolType = reflect.TypeOf((*Tool)(nil)).Elem()

type ToolKey struct {
	ToolKind string
	ToolName string
}

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

	Run(
		ctx context.Context,
		allTools map[ToolKey]Tool,
		allShells map[ToolKey]*BaseTool,
		taskKind, taskName string,
	) error

	GetEnv() []string
}

type BaseTool struct {
	field.BaseField

	Name string   `yaml:"name"`
	Env  []string `yaml:"env"`
	Cmd  []string `yaml:"cmd"`

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
func (t *BaseTool) GetEnv() []string { return sliceutils.NewStrings(t.Env) }

func (t *BaseTool) RunTask(
	ctx context.Context,
	thisTool Tool,
	allTools map[ToolKey]Tool,
	allShells map[ToolKey]*BaseTool,
	task Task,
) error {
	workerCount := constant.GetWorkerCount(ctx)

	ctx, cancelExec := context.WithCancel(
		// all sub tasks for this task should only have one worker
		constant.WithWorkerCount(ctx, 1),
	)
	defer cancelExec()

	baseCtx := field.WithRenderingValues(ctx, t.Env)

	matrixSpecs, err := task.GetMatrixSpecs(
		baseCtx, t.RenderingFunc,
		constant.GetMatrixFilter(baseCtx.Context()),
	)
	if err != nil {
		return fmt.Errorf("failed to create build matrix: %w", err)
	}

	if len(matrixSpecs) == 0 {
		return fmt.Errorf("no matrix spec match")
	}

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
		{color.New(color.FgHiWhite), color.New(color.FgWhite)},
		{color.New(color.FgHiCyan), color.New(color.FgCyan)},
		{color.New(color.FgHiGreen), color.New(color.FgGreen)},
		{color.New(color.FgHiMagenta), color.New(color.FgMagenta)},
		{color.New(color.FgHiYellow), color.New(color.FgYellow)},
		{color.New(color.FgHiBlue), color.New(color.FgBlue)},
		{color.New(color.FgHiRed), color.New(color.FgRed)},
	}

	failFast := constant.IsFailFast(baseCtx.Context())

	wg := &sync.WaitGroup{}

	// ready

	taskPrefix := fmt.Sprintf(
		"%s [ %s ]",
		output.AssembleTaskKindID(task.ToolKind(), task.ToolName(), task.TaskKind()),
		task.TaskName(),
	)

	err = task.RunHooks(
		baseCtx, t.RenderingFunc,
		taskExecBeforeStart,
		taskExecBeforeStart.String()+" "+taskPrefix,
		nil, nil,
		thisTool, allTools, allShells,
	)
	if err != nil {
		return fmt.Errorf("failed to run hooks %q: %w",
			taskExecBeforeStart.String(), err,
		)
	}

	for i, ms := range matrixSpecs {
		color := colorList[i%len(colorList)]
		prefixColor, outputColor := color[0], color[1]

		mFilter := make(map[string][]string)
		for k, v := range ms {
			mFilter[k] = []string{v}
		}

		ctx := constant.WithMatrixFilter(baseCtx.Context(), mFilter)
		taskCtx := field.WithRenderingValues(ctx, nil)
		taskCtx.SetEnv(baseCtx.Values().Env)

		select {
		case <-taskCtx.Context().Done():
			return taskCtx.Context().Err()
		case <-waitCh:
		}

		err2 := func() error {
			output.WriteTaskStart(
				taskCtx.Context(),
				task.ToolKind(), task.ToolName(), task.TaskKind(), task.TaskName(),
				ms.String(),
			)

			for k, v := range ms {
				taskCtx.AddEnv("MATRIX_" + strings.ToUpper(k) + "=" + v)
			}

			// tool may have reference to MATRIX_ values
			err3 := t.ResolveFields(taskCtx, t.RenderingFunc, -1)
			if err3 != nil {
				return fmt.Errorf("failed to resolve tool fields: %w", err3)
			}

			taskCtx.AddEnv(t.Env...)

			// resolve tasks
			err3 = task.ResolveFields(taskCtx, t.RenderingFunc, -1)
			if err3 != nil {
				return fmt.Errorf("failed to resolve task fields: %w", err3)
			}

			toolCmd := sliceutils.NewStrings(t.Cmd)
			if len(toolCmd) == 0 {
				toolCmd = append(toolCmd, t.defaultExecutable)
			}

			execSpecs, err3 := task.GetExecSpecs(taskCtx, toolCmd)
			if err3 != nil {
				return fmt.Errorf("failed to generate task args: %w", err3)
			}

			prefix := ms.BriefString() + ": "

			err3 = task.RunHooks(
				taskCtx, t.RenderingFunc,
				taskExecBeforeMatrixStart,
				taskExecBeforeMatrixStart.String()+" "+prefix,
				prefixColor, outputColor,
				thisTool, allTools, allShells,
			)
			if err3 != nil {
				return fmt.Errorf("failed to run hooks %q: %w",
					taskExecBeforeMatrixStart.String(), err3,
				)
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

				err4 := t.doRunTask(
					taskCtx,
					prefix,
					prefixColor, outputColor,
					execSpecs, nil,
				)
				output.WriteExecResult(
					taskCtx.Context(),
					task.ToolKind(), task.ToolName(), task.TaskKind(), task.TaskName(),
					ms.String(),
					err4,
				)

				if err4 != nil {
					if failFast {
						cancelExec()
					}

					appendResult(ms, err4)

					err4 = task.RunHooks(
						taskCtx, t.RenderingFunc,
						taskExecAfterMatrixFailure,
						taskExecAfterMatrixFailure.String()+" "+prefix,
						prefixColor, outputColor,
						thisTool, allTools, allShells,
					)
					if err4 != nil {
						// TODO: handle hook error
						err4 = fmt.Errorf(
							"failed to run hooks %q: %w",
							taskExecAfterMatrixFailure.String(), err4,
						)
						_ = err4
					}

					return
				}

				err4 = task.RunHooks(
					taskCtx, t.RenderingFunc,
					taskExecAfterMatrixSuccess,
					taskExecAfterMatrixSuccess.String()+" "+prefix,
					prefixColor, outputColor,
					thisTool, allTools, allShells,
				)
				if err4 != nil {
					// TODO: handle hook error
					err4 = fmt.Errorf("failed to run hooks %q: %w",
						taskExecAfterMatrixSuccess.String(), err4,
					)
					_ = err4
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
		err = task.RunHooks(
			baseCtx, t.RenderingFunc,
			taskExecAfterFailure,
			taskExecAfterFailure.String()+" "+taskPrefix,
			nil, nil,
			thisTool, allTools, allShells,
		)
		if err != nil {
			// TODO: handle hook error
			err = fmt.Errorf("failed to run hooks %q: %w",
				taskExecAfterFailure.String(), err,
			)
			_ = err
		}

		return fmt.Errorf("task execution failed")
	}

	err = task.RunHooks(
		baseCtx, t.RenderingFunc,
		taskExecAfterSuccess,
		taskExecAfterSuccess.String()+" "+taskPrefix,
		nil, nil,
		thisTool, allTools, allShells,
	)
	if err != nil {
		// TODO: handle hook error
		err = fmt.Errorf("failed to run hooks %q: %w",
			taskExecAfterSuccess.String(), err,
		)
		_ = err
	}

	return nil
}

func (t *BaseTool) doRunTask(
	taskCtx *field.RenderingContext,
	outputPrefix string,
	prefixColor, outputColor *color.Color,
	execSpecs []TaskExecSpec,
	_replaceEntries *map[string][]byte,
) error {
	timer := time.NewTimer(0)
	if !timer.Stop() {
		<-timer.C
	}

	var replace map[string][]byte
	if _replaceEntries != nil {
		replace = *_replaceEntries
	} else {
		replace = make(map[string][]byte)
	}

	for _, es := range execSpecs {
		ctx := taskCtx.Clone()

		if es.Delay > 0 {
			_ = timer.Reset(es.Delay)

			select {
			case <-timer.C:
			case <-ctx.Context().Done():
				if !timer.Stop() {
					<-timer.C
				}

				return ctx.Context().Err()
			}
		}

		var (
			stdin          io.Reader
			stdout, stderr io.Writer
		)

		if es.Stdin != nil {
			stdin = io.MultiReader(es.Stdin, os.Stdin)
		} else {
			stdin = os.Stdin
		}

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

		var buf *bytes.Buffer
		if len(es.OutputAsReplace) != 0 {
			buf = &bytes.Buffer{}

			stdout = io.MultiWriter(stdout, buf)
		}

		// alter exec func can generate sub exec specs
		if es.AlterExecFunc != nil {
			subSpecs, err := es.AlterExecFunc(replace, stdin, stdout, stderr)
			if err != nil {
				return fmt.Errorf("failed to execute alter exec func: %w", err)
			}

			if buf != nil {
				newValue := buf.Bytes()
				if es.FixOutputForReplace != nil {
					newValue = es.FixOutputForReplace(newValue)
				}

				replace[es.OutputAsReplace] = newValue
			}

			if len(subSpecs) != 0 {
				err = t.doRunTask(taskCtx, outputPrefix, prefixColor, outputColor, subSpecs, &replace)
				if err != nil {
					return fmt.Errorf("failed to run sub tasks: %w", err)
				}
			}

			continue
		}

		var cmd []string
		if len(replace) != 0 {
			pairs := make([]string, 2*len(replace))
			i := 0
			for toReplace, newValue := range replace {
				pairs[i], pairs[i+1] = toReplace, string(newValue)
				i += 2
			}

			replacer := strings.NewReplacer(pairs...)
			for _, rawEnvPart := range es.Env {
				ctx.AddEnv(replacer.Replace(rawEnvPart))
			}

			for _, rawCmdPart := range es.Command {
				cmd = append(cmd, replacer.Replace(rawCmdPart))
			}
		} else {
			cmd = sliceutils.NewStrings(es.Command)
		}

		_, runScriptCmd, err := t.getBootstrapExecSpec(cmd, false)
		if err != nil {
			return fmt.Errorf("failed to get exec spec from bootstrap config: %w", err)
		}

		output.WriteExecStart(
			ctx.Context(),
			t.ToolName(),
			cmd,
			filepath.Base(runScriptCmd[len(runScriptCmd)-1]),
		)

		p, err := exechelper.Do(exechelper.Spec{
			Context: ctx.Context(),
			Command: runScriptCmd,
			Env:     ctx.Values().Env,
			Dir:     es.Chdir,

			Stdin: stdin,

			Stdout: stdout,
			Stderr: stderr,
		})
		if err != nil {
			if !es.IgnoreError {
				return fmt.Errorf("failed to prepare command [ %s ]: %w", strings.Join(cmd, " "), err)
			}

			// TODO: log error in detail
			log.Log.I("error ignored", log.Error(err))

			delete(replace, es.OutputAsReplace)

			continue
		}

		_, err = p.Wait()
		if err != nil {
			if !es.IgnoreError {
				return fmt.Errorf("command exited with error: %w", err)
			}

			// TODO: log error in detail
			log.Log.I("error ignored", log.Error(err))

			delete(replace, es.OutputAsReplace)

			continue
		}

		if buf != nil {
			newValue := buf.Bytes()
			if es.FixOutputForReplace != nil {
				newValue = es.FixOutputForReplace(newValue)
			}

			replace[es.OutputAsReplace] = newValue
		}
	}

	return nil
}

// GetExecSpec is a helper func for shell renderer
func (t *BaseTool) GetExecSpec(toExec []string, isFilePath bool) (env, cmd []string, err error) {
	if len(toExec) == 0 {
		return nil, nil, fmt.Errorf("invalid empty exec spec")
	}

	scriptPath := ""
	if !isFilePath {
		scriptPath, err = GetScriptCache(t.cacheDir, strings.Join(toExec, " "))
		if err != nil {
			return nil, nil, fmt.Errorf("tools: failed to ensure script cache: %w", err)
		}
	} else {
		scriptPath = toExec[0]
	}

	cmd = sliceutils.NewStrings(t.Cmd)
	if len(cmd) == 0 {
		cmd = append(cmd, t.defaultExecutable)
	}

	return t.Env, append(cmd, scriptPath), nil
}
