package tools

import (
	"fmt"
	"strings"
	"sync"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/matrix"
	"arhat.dev/dukkha/pkg/output"
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

	// resolve hooks for whole task
	options := dukkha.TaskExecOptions{
		UseShell:        tool.UseShell(),
		ShellName:       tool.ShellName(),
		ToolCmd:         tool.GetCmd(),
		ContinueOnError: !ctx.FailFast(),
	}

	ctx.SetTask(tool.Key(), task.Key())

	wg := &sync.WaitGroup{}

	// ensure hook `after` always run
	defer func() {
		// TODO: handle hook error
		hookAfter, err2 := task.GetHookExecSpecs(
			ctx, dukkha.StageAfter, options,
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
		ctx, dukkha.StageBefore, options,
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

matrixRun:
	for i, ms := range matrixSpecs {
		mCtx, err3 := createTaskMatrixContext(ctx, i, ms, tool)

		if err3 != nil {
			appendErrorResult(ms, err3)
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
			defer func() {
				if err != nil && ctx.FailFast() {
					ctx.Cancel()
				}

				hookAfterMatrix, err3 := task.GetHookExecSpecs(
					mCtx, dukkha.StageAfterMatrix, options,
				)
				if err3 != nil {
					appendErrorResult(ms, err3)
					return
				}

				// TODO: handle hook error
				err3 = runHook(mCtx, dukkha.StageAfterMatrix, hookAfterMatrix)
				if err3 != nil {
					appendErrorResult(ms, err3)
					return
				}

				wg.Done()

				select {
				case waitCh <- struct{}{}:
				case <-mCtx.Done():
					return
				}
			}()

			hookBeofreMatrix, err3 := task.GetHookExecSpecs(
				mCtx, dukkha.StageBeforeMatrix, options,
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
					mCtx, dukkha.StageAfterMatrixFailure, options,
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
				mCtx, dukkha.StageAfterMatrixSuccess, options,
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
			ctx, dukkha.StageAfterFailure, options,
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
		ctx, dukkha.StageAfterSuccess, options,
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
) (dukkha.TaskExecContext, error) {
	// set default matrix filter for referenced hook tasks
	mFilter := make(map[string][]string)
	for k, v := range ms {
		mFilter[k] = []string{v}
	}

	// mCtx is the matrix execution context

	mCtx := ctx.DeriveNew()
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
	err := tool.DoAfterFieldsResolved(mCtx, func() error {
		mCtx.AddEnv(tool.GetEnv()...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add tool env: %w", err)
	}

	return mCtx, nil
}
