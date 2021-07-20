package tools

import (
	"fmt"
	"strings"
	"sync"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/matrix"
	"arhat.dev/dukkha/pkg/output"
	"arhat.dev/dukkha/pkg/sliceutils"
)

type TaskExecRequest struct {
	Context dukkha.TaskExecContext

	Tool dukkha.Tool
	Task dukkha.Task

	IgnoreError bool
}

// nolint:gocyclo
func RunTask(req *TaskExecRequest) error {
	type taskResult struct {
		matrixSpec string
		errMsg     string
	}

	var (
		errCollection []taskResult

		resultMU = &sync.Mutex{}
	)

	appendErrorResult := func(spec matrix.Entry, err error) {
		resultMU.Lock()
		defer resultMU.Unlock()

		errCollection = append(errCollection, taskResult{
			matrixSpec: spec.BriefString(),
			errMsg:     err.Error(),
		})
	}

	// task may need tool specific env, resolve tool env first

	err := req.Tool.DoAfterFieldsResolved(req.Context, -1, func() error {
		req.Context.AddEnv(req.Tool.GetEnv()...)
		return nil
	}, "BaseTool.Env")
	if err != nil {
		return fmt.Errorf("failed to resolve tool specific env: %w", err)
	}

	// resolve hooks for whole task

	req.Context.SetTask(req.Tool.Key(), req.Task.Key())

	wg := &sync.WaitGroup{}

	// ensure hook `after` always run
	defer func() {
		// TODO: handle hook error
		hookAfter, err2 := req.Task.GetHookExecSpecs(
			req.Context, dukkha.StageAfter,
		)
		if err2 != nil {
			appendErrorResult(make(matrix.Entry), err2)
		}

		err2 = runHook(req.Context, dukkha.StageAfter, hookAfter)
		if err2 != nil {
			appendErrorResult(make(matrix.Entry), err2)
		}
	}()

	// run hook `before`
	hookBefore, err := req.Task.GetHookExecSpecs(
		req.Context, dukkha.StageBefore,
	)
	if err != nil {
		// cancel task execution
		return err
	}

	err = runHook(req.Context, dukkha.StageBefore, hookBefore)
	if err != nil {
		// cancel task execution
		return err
	}

	matrixSpecs, err := req.Task.GetMatrixSpecs(req.Context)
	if err != nil {
		return fmt.Errorf("failed to get execution matrix: %w", err)
	}

	if len(matrixSpecs) == 0 {
		return fmt.Errorf("no matrix spec match")
	}

	// TODO: do real global limit
	workerCount := req.Context.ClaimWorkers(len(matrixSpecs))
	waitCh := make(chan struct{}, workerCount)
	for i := 0; i < workerCount; i++ {
		waitCh <- struct{}{}
	}

matrixRun:
	for i, ms := range matrixSpecs {
		mCtx, options, err2 := createTaskMatrixContext(req, i, ms)

		if err2 != nil {
			appendErrorResult(ms, err2)
			if req.Context.FailFast() {
				req.Context.Cancel()
				break matrixRun
			}

			continue
		}

		select {
		case <-mCtx.Done():
			break matrixRun
		case <-waitCh:
		}

		output.WriteTaskStart(
			mCtx.PrefixColor(),
			mCtx.CurrentTool(),
			mCtx.CurrentTask(),
			ms,
		)

		wg.Add(1)
		go func(ms matrix.Entry) {
			var err3 error
			defer func() {
				defer func() {
					wg.Done()

					select {
					case waitCh <- struct{}{}:
					case <-mCtx.Done():
						return
					}
				}()

				if err3 != nil && req.Context.FailFast() {
					req.Context.Cancel()
				}

				hookAfterMatrix, err4 := req.Task.GetHookExecSpecs(
					mCtx, dukkha.StageAfterMatrix,
				)
				if err4 != nil {
					appendErrorResult(ms, err4)
					return
				}

				// TODO: handle hook error
				err4 = runHook(mCtx, dukkha.StageAfterMatrix, hookAfterMatrix)
				if err4 != nil {
					appendErrorResult(ms, err4)
					return
				}
			}()

			hookBeofreMatrix, err3 := req.Task.GetHookExecSpecs(
				mCtx, dukkha.StageBeforeMatrix,
			)
			if err3 != nil {
				appendErrorResult(ms, err3)
				return
			}

			err3 = runHook(mCtx, dukkha.StageBeforeMatrix, hookBeofreMatrix)
			if err3 != nil {
				appendErrorResult(ms, err3)
				return
			}

			// produce a snapshot of what to do
			execSpecs, err3 := req.Task.GetExecSpecs(mCtx, options)
			if err3 != nil {
				appendErrorResult(
					ms,
					fmt.Errorf("failed to generate task exec specs: %w", err3),
				)
				return
			}

			err3 = doRun(mCtx, execSpecs, nil)

			output.WriteExecResult(
				mCtx.PrefixColor(),
				mCtx.CurrentTool(),
				mCtx.CurrentTask(),
				ms.String(),
				err3,
			)

			if err3 != nil {
				// cancel other tasks if in fail-fast mode
				if req.Context.FailFast() {
					req.Context.Cancel()
				}

				appendErrorResult(ms, err3)

				hookAfterMatrixFailure, err4 := req.Task.GetHookExecSpecs(
					mCtx, dukkha.StageAfterMatrixFailure,
				)
				if err4 != nil {
					appendErrorResult(ms, err4)
					return
				}

				err4 = runHook(
					mCtx,
					dukkha.StageAfterMatrixFailure,
					hookAfterMatrixFailure,
				)
				if err4 != nil {
					appendErrorResult(ms, err4)
					return
				}

				return
			}

			hookAfterMatrixSuccess, err3 := req.Task.GetHookExecSpecs(
				mCtx, dukkha.StageAfterMatrixSuccess,
			)
			if err3 != nil {
				appendErrorResult(ms, err3)
				return
			}

			err3 = runHook(mCtx, dukkha.StageAfterMatrixSuccess, hookAfterMatrixSuccess)
			if err3 != nil {
				appendErrorResult(ms, err3)
				return
			}
		}(ms)
	}

	wg.Wait()

	if len(errCollection) != 0 {
		hookAfterFailure, err2 := req.Task.GetHookExecSpecs(
			req.Context, dukkha.StageAfterFailure,
		)
		if err2 != nil {
			appendErrorResult(make(matrix.Entry), err2)
		}

		err2 = runHook(req.Context, dukkha.StageAfterFailure, hookAfterFailure)
		if err2 != nil {
			appendErrorResult(make(matrix.Entry), err2)
		}

		return fmt.Errorf("task execution failed: %v", errCollection)
	}

	hookAfterSuccess, err := req.Task.GetHookExecSpecs(
		req.Context, dukkha.StageAfterSuccess,
	)
	if err != nil {
		appendErrorResult(make(matrix.Entry), err)
	}

	err = runHook(req.Context, dukkha.StageAfterSuccess, hookAfterSuccess)
	if err != nil {
		appendErrorResult(make(matrix.Entry), err)
	}

	return nil
}

// createTaskMatrixContext creates a per matrix entry task exec options
// with context resolved
func createTaskMatrixContext(
	req *TaskExecRequest,
	i int,
	ms matrix.Entry,
) (dukkha.TaskExecContext, *dukkha.TaskExecOptions, error) {
	mCtx := req.Context.DeriveNew()

	// set default matrix filter for referenced hook tasks
	mFilter := make(map[string][]string)
	for k, v := range ms {
		mFilter[k] = []string{v}
	}

	mCtx.SetMatrixFilter(mFilter)

	for k, v := range ms {
		mCtx.AddEnv("MATRIX_" + strings.ToUpper(k) + "=" + v)
	}

	existingPrefix := mCtx.OutputPrefix()
	if len(existingPrefix) != 0 {
		if !strings.HasPrefix(existingPrefix, ms.BriefString()) {
			// not same matrix, add this matrix prefix
			mCtx.SetOutputPrefix(existingPrefix + ms.BriefString() + ": ")
		}
	} else {
		mCtx.SetOutputPrefix(ms.BriefString() + ": ")
	}

	mCtx.SetTaskColors(output.PickColor(i))

	// tool may have reference to MATRIX_ values
	// but MUST not have reference to task specific env

	// now everything prepared for the tool, resolve all of it

	options := &dukkha.TaskExecOptions{}
	err := req.Tool.DoAfterFieldsResolved(mCtx, -1, func() error {
		mCtx.AddEnv(req.Tool.GetEnv()...)

		options.ToolCmd = sliceutils.NewStrings(req.Tool.GetCmd())
		options.UseShell = req.Tool.UseShell()
		options.ShellName = req.Tool.ShellName()

		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to add tool env: %w", err)
	}

	return mCtx, options, nil
}
