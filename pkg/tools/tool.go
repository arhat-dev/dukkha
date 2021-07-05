package tools

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"

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

	baseCtx := field.WithRenderingValues(ctx, os.Environ(), t.Env)

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

	appendErrorResult := func(spec MatrixSpec, err error) {
		resultMU.Lock()
		defer resultMU.Unlock()

		errCollection = append(errCollection, taskResult{
			matrixSpec: spec,
			err:        err,
		})
	}

	failFast := constant.IsFailFast(baseCtx.Context())

	wg := &sync.WaitGroup{}

	// ready

	taskPrefix := fmt.Sprintf(
		"%s [ %s ]",
		output.AssembleTaskKindID(task.ToolKind(), task.ToolName(), task.TaskKind()),
		task.TaskName(),
	)

	// ensure hook `after` always run
	defer func() {
		// TODO: handle hook error
		_ = HandleHookRunError(
			StageAfter,
			task.RunHooks(
				baseCtx, t.RenderingFunc,
				StageAfter,
				StageAfter.String()+" "+taskPrefix,
				nil, nil,
				thisTool, allTools, allShells,
			),
		)
	}()

	// run hook `before`
	err = HandleHookRunError(
		StageBefore,
		task.RunHooks(
			baseCtx, t.RenderingFunc,
			StageBefore,
			StageBefore.String()+" "+taskPrefix,
			nil, nil,
			thisTool, allTools, allShells,
		),
	)
	if err != nil {
		// cancel task execution
		return err
	}

	for i, ms := range matrixSpecs {
		// set default matrix filter for referenced hook tasks
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

			prefix := ms.BriefString() + ": "

			prefixColor, outputColor := output.PickColor(i)
			err3 := HandleHookRunError(
				StageBeforeMatrix,
				task.RunHooks(
					taskCtx, t.RenderingFunc,
					StageBeforeMatrix,
					StageBeforeMatrix.String()+" "+prefix,
					prefixColor, outputColor,
					thisTool, allTools, allShells,
				),
			)
			if err3 != nil {
				return err3
			}

			// tool may have reference to MATRIX_ values
			err3 = t.ResolveFields(taskCtx, t.RenderingFunc, -1, false)
			if err3 != nil {
				return fmt.Errorf("failed to resolve tool fields: %w", err3)
			}

			taskCtx.AddEnv(t.Env...)

			// resolve task fields
			err3 = task.ResolveFields(taskCtx, t.RenderingFunc, -1, false)
			if err3 != nil {
				return fmt.Errorf("failed to resolve task fields: %w", err3)
			}

			toolCmd := sliceutils.NewStrings(t.Cmd)
			if len(toolCmd) == 0 {
				toolCmd = append(toolCmd, t.defaultExecutable)
			}

			// produce a snapshot of what to do
			execSpecs, err3 := task.GetExecSpecs(taskCtx, toolCmd)
			if err3 != nil {
				return fmt.Errorf("failed to generate task exec specs: %w", err3)
			}

			wg.Add(1)
			go func(ms MatrixSpec) {
				defer func() {
					// TODO: handle hook error
					_ = HandleHookRunError(
						StageAfterMatrix,
						task.RunHooks(
							taskCtx, t.RenderingFunc,
							StageAfterMatrix,
							StageAfterMatrix.String()+" "+prefix,
							prefixColor, outputColor,
							thisTool, allTools, allShells,
						),
					)

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

					appendErrorResult(ms, err4)

					err4 = HandleHookRunError(
						StageAfterMatrixFailure,
						task.RunHooks(
							taskCtx, t.RenderingFunc,
							StageAfterMatrixFailure,
							StageAfterMatrixFailure.String()+" "+prefix,
							prefixColor, outputColor,
							thisTool, allTools, allShells,
						),
					)
					if err4 != nil {
						// TODO: handle hook error
						_ = err4
					}

					return
				}

				err4 = HandleHookRunError(
					StageAfterMatrixSuccess,
					task.RunHooks(
						taskCtx, t.RenderingFunc,
						StageAfterMatrixSuccess,
						StageAfterMatrixSuccess.String()+" "+prefix,
						prefixColor, outputColor,
						thisTool, allTools, allShells,
					),
				)
				if err4 != nil {
					// TODO: handle hook error
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
		err = HandleHookRunError(
			StageAfterFailure,
			task.RunHooks(
				baseCtx, t.RenderingFunc,
				StageAfterFailure,
				StageAfterFailure.String()+" "+taskPrefix,
				nil, nil,
				thisTool, allTools, allShells,
			),
		)
		if err != nil {
			// TODO: handle hook error
			_ = err
		}

		// TODO: handle execution error
		return fmt.Errorf("task execution failed: %v", errCollection)
	}

	err = HandleHookRunError(
		StageAfterSuccess,
		task.RunHooks(
			baseCtx, t.RenderingFunc,
			StageAfterSuccess,
			StageAfterSuccess.String()+" "+taskPrefix,
			nil, nil,
			thisTool, allTools, allShells,
		),
	)
	if err != nil {
		// TODO: handle hook error
		_ = err
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
