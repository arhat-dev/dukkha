package tools

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"golang.org/x/term"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/matrix"
	"arhat.dev/dukkha/pkg/output"
	"arhat.dev/dukkha/pkg/sliceutils"
)

var _ dukkha.Tool = (*baseToolWithKind)(nil)

type baseToolWithKind struct{ BaseTool }

func (*baseToolWithKind) Kind() dukkha.ToolKind { return "" }

type BaseTool struct {
	field.BaseField

	ToolName string   `yaml:"name"`
	Env      []string `yaml:"env"`
	Cmd      []string `yaml:"cmd"`

	cacheDir          string `json:"-" yaml:"-"`
	defaultExecutable string `json:"-" yaml:"-"`
	stdoutIsTty       bool   `json:"-" yaml:"-"`

	Tasks map[dukkha.TaskKey]dukkha.Task `json:"-" yaml:"-"`
}

// Init the tool, override it if your tool name is different from
// its default executable name
func (t *BaseTool) Init(cachdDir string) error {
	return t.InitBaseTool(t.ToolName, cachdDir)
}

// InitBaseTool must be called in your own version of Init()
// with correct defaultExecutable name
func (t *BaseTool) InitBaseTool(defaultExecutable, cacheDir string) error {
	t.defaultExecutable = defaultExecutable
	t.stdoutIsTty = term.IsTerminal(int(os.Stdout.Fd()))
	return nil
}

// ResolveTasks accpets all tasks, override this function if your tool need
// different handling of tasks
func (t *BaseTool) ResolveTasks(tasks []dukkha.Task) error {
	for i, tsk := range tasks {
		t.Tasks[dukkha.TaskKey{Kind: tsk.Kind(), Name: tsk.Name()}] = tasks[i]
	}

	return nil
}

// Run task
func (t *BaseTool) Run(taskCtx dukkha.TaskExecContext) error {
	tsk, ok := t.Tasks[dukkha.TaskKey{Kind: "", Name: ""}]
	if !ok {
		return fmt.Errorf("task %q not found")
	}

	return t.RunTask(taskCtx, tsk)
}

func (t *BaseTool) Name() dukkha.ToolName { return dukkha.ToolName(t.ToolName) }
func (t *BaseTool) GetEnv() []string      { return sliceutils.NewStrings(t.Env) }

func (t *BaseTool) RunTask(taskCtx dukkha.TaskExecContext, task dukkha.Task) error {
	defer taskCtx.Cancel()

	matrixSpecs, err := task.GetMatrixSpecs(taskCtx)
	if err != nil {
		return fmt.Errorf("failed to create build matrix: %w", err)
	}

	if len(matrixSpecs) == 0 {
		return fmt.Errorf("no matrix spec match")
	}

	workerCount := taskCtx.ClaimWorkers(len(matrixSpecs))
	waitCh := make(chan struct{}, workerCount)
	for i := 0; i < workerCount; i++ {
		waitCh <- struct{}{}
	}

	type taskResult struct {
		matrixSpec matrix.Spec
		err        error
	}

	var (
		errCollection []taskResult

		resultMU = &sync.Mutex{}
	)

	appendErrorResult := func(spec matrix.Spec, err error) {
		resultMU.Lock()
		defer resultMU.Unlock()

		errCollection = append(errCollection, taskResult{
			matrixSpec: spec,
			err:        err,
		})
	}

	wg := &sync.WaitGroup{}

	// ensure hook `after` always run
	defer func() {
		// TODO: handle hook error
		_ = task.RunHooks(taskCtx, dukkha.StageAfter)
	}()

	// run hook `before`
	err = task.RunHooks(taskCtx, dukkha.StageBefore)
	if err != nil {
		// cancel task execution
		return err
	}

matrixRun:
	for i, ms := range matrixSpecs {
		// set default matrix filter for referenced hook tasks
		mFilter := make(map[string][]string)
		for k, v := range ms {
			mFilter[k] = []string{v}
		}

		mCtx := taskCtx.DeriveNew()
		mCtx.SetMatrixFilter(mFilter)

		select {
		case <-mCtx.Done():
			break matrixRun
		case <-waitCh:
		}

		err2 := func() error {
			output.WriteTaskStart(
				task.ToolKind(), task.ToolName(),
				task.Kind(), task.Name(),
				ms,
			)

			for k, v := range ms {
				mCtx.AddEnv("MATRIX_" + strings.ToUpper(k) + "=" + v)
			}

			mCtx.SetOutputPrefix(ms.BriefString() + ":")
			mCtx.SetTaskColors(output.PickColor(i))

			err3 := task.RunHooks(mCtx, dukkha.StageBeforeMatrix)
			if err3 != nil {
				return err3
			}

			// tool may have reference to MATRIX_ values
			err3 = t.ResolveFields(mCtx, -1, "")
			if err3 != nil {
				return fmt.Errorf("failed to resolve tool fields: %w", err3)
			}

			mCtx.AddEnv(t.Env...)

			// resolve task fields
			err3 = task.ResolveFields(mCtx, -1, "")
			if err3 != nil {
				return fmt.Errorf("failed to resolve task fields: %w", err3)
			}

			toolCmd := sliceutils.NewStrings(t.Cmd)
			if len(toolCmd) == 0 {
				toolCmd = append(toolCmd, t.defaultExecutable)
			}

			// produce a snapshot of what to do
			execSpecs, err3 := task.GetExecSpecs(mCtx, toolCmd)
			if err3 != nil {
				return fmt.Errorf("failed to generate task exec specs: %w", err3)
			}

			wg.Add(1)
			go func(ms matrix.Spec) {
				defer func() {
					// TODO: handle hook error
					_ = task.RunHooks(mCtx, dukkha.StageAfterMatrix)

					wg.Done()

					select {
					case waitCh <- struct{}{}:
					case <-mCtx.Done():
						return
					}
				}()

				err4 := t.doRunTask(mCtx, execSpecs, nil)

				output.WriteExecResult(
					taskCtx,
					task.ToolKind(), task.ToolName(), task.Kind(), task.Name(),
					ms.String(),
					err4,
				)

				if err4 != nil {
					// cancel other tasks if in fail-fast mode
					if taskCtx.FailFast() {
						taskCtx.Cancel()
					}

					appendErrorResult(ms, err4)

					err4 = task.RunHooks(mCtx, dukkha.StageAfterMatrixFailure)
					if err4 != nil {
						// TODO: handle hook error
						_ = err4
					}

					return
				}

				err4 = task.RunHooks(mCtx, dukkha.StageAfterMatrixSuccess)
				if err4 != nil {
					// TODO: handle hook error
					_ = err4
				}
			}(ms)

			return nil
		}()

		if err2 != nil {
			// failed before execution
			if taskCtx.FailFast() {
				taskCtx.Cancel()
			}

			appendErrorResult(ms, err2)
		}
	}

	wg.Wait()

	if len(errCollection) != 0 {
		err = task.RunHooks(taskCtx, dukkha.StageAfterFailure)
		if err != nil {
			// TODO: handle hook error
			_ = err
		}

		// TODO: handle execution error
		return fmt.Errorf("task execution failed: %v", errCollection)
	}

	err = task.RunHooks(taskCtx, dukkha.StageAfterSuccess)
	if err != nil {
		// TODO: handle hook error
		_ = err
	}

	return nil
}

// GetExecSpec is a helper func for shells
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
