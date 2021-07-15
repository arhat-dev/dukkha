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

type CompleteTaskExecSpecs struct {
	Context dukkha.TaskExecContext

	Tool dukkha.Tool
	Task dukkha.Task
}

// nolint:gocyclo
func RunTask(ctx dukkha.TaskExecContext, tool dukkha.Tool, task dukkha.Task) error {
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

	err := tool.DoAfterFieldsResolved(ctx, -1, func() error {
		ctx.AddEnv(tool.GetEnv()...)
		return nil
	}, "BaseTool.Env")
	if err != nil {
		return fmt.Errorf("failed to resolve tool specific env: %w", err)
	}

	// resolve hooks for whole task

	ctx.SetTask(tool.Key(), task.Key())

	wg := &sync.WaitGroup{}

	// ensure hook `after` always run
	defer func() {
		// TODO: handle hook error
		hookAfter, err2 := task.GetHookExecSpecs(
			ctx, dukkha.StageAfter,
		)
		if err2 != nil {
			appendErrorResult(make(matrix.Entry), err2)
		}

		err2 = runHook(ctx, dukkha.StageAfter, hookAfter)
		if err2 != nil {
			appendErrorResult(make(matrix.Entry), err2)
		}
	}()

	// run hook `before`
	hookBefore, err := task.GetHookExecSpecs(
		ctx, dukkha.StageBefore,
	)
	if err != nil {
		// cancel task execution
		return err
	}

	err = runHook(ctx, dukkha.StageBefore, hookBefore)
	if err != nil {
		// cancel task execution
		return err
	}

	matrixSpecs, err := task.GetMatrixSpecs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get execution matrix: %w", err)
	}

	if len(matrixSpecs) == 0 {
		return fmt.Errorf("no matrix spec match")
	}

	// TODO: do real global limit
	workerCount := ctx.ClaimWorkers(len(matrixSpecs))
	waitCh := make(chan struct{}, workerCount)
	for i := 0; i < workerCount; i++ {
		waitCh <- struct{}{}
	}

matrixRun:
	for i, ms := range matrixSpecs {
		mCtx, options, err2 := createTaskMatrixContext(ctx, i, ms, tool)

		if err2 != nil {
			appendErrorResult(ms, err2)
			if ctx.FailFast() {
				ctx.Cancel()
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

				if err3 != nil && ctx.FailFast() {
					ctx.Cancel()
				}

				hookAfterMatrix, err4 := task.GetHookExecSpecs(
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

			hookBeofreMatrix, err3 := task.GetHookExecSpecs(
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
			execSpecs, err3 := task.GetExecSpecs(mCtx, options)
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
				if ctx.FailFast() {
					ctx.Cancel()
				}

				appendErrorResult(ms, err3)

				hookAfterMatrixFailure, err4 := task.GetHookExecSpecs(
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

			hookAfterMatrixSuccess, err3 := task.GetHookExecSpecs(
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
		hookAfterFailure, err2 := task.GetHookExecSpecs(
			ctx, dukkha.StageAfterFailure,
		)
		if err2 != nil {
			appendErrorResult(make(matrix.Entry), err2)
		}

		err2 = runHook(ctx, dukkha.StageAfterFailure, hookAfterFailure)
		if err2 != nil {
			appendErrorResult(make(matrix.Entry), err2)
		}

		return fmt.Errorf("task execution failed: %v", errCollection)
	}

	hookAfterSuccess, err := task.GetHookExecSpecs(
		ctx, dukkha.StageAfterSuccess,
	)
	if err != nil {
		appendErrorResult(make(matrix.Entry), err)
	}

	err = runHook(ctx, dukkha.StageAfterSuccess, hookAfterSuccess)
	if err != nil {
		appendErrorResult(make(matrix.Entry), err)
	}

	return nil
}

func createTaskMatrixContext(
	ctx dukkha.TaskExecContext,
	i int,
	ms matrix.Entry,
	tool dukkha.Tool,
) (dukkha.TaskExecContext, *dukkha.TaskExecOptions, error) {
	mCtx := ctx.DeriveNew()

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
	err := tool.DoAfterFieldsResolved(mCtx, -1, func() error {
		mCtx.AddEnv(tool.GetEnv()...)

		options.ToolCmd = sliceutils.NewStrings(tool.GetCmd())
		options.UseShell = tool.UseShell()
		options.ShellName = tool.ShellName()

		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to add tool env: %w", err)
	}

	return mCtx, options, nil
}
