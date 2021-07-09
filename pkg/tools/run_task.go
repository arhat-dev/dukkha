package tools

import (
	"fmt"
	"strings"
	"sync"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/matrix"
	"arhat.dev/dukkha/pkg/output"
)

// nolint:gocyclo
func runTask(ctx dukkha.TaskExecContext, tool dukkha.Tool, task dukkha.Task) error {
	defer ctx.Cancel()

	matrixSpecs, err := task.GetMatrixSpecs(ctx)
	if err != nil {
		return fmt.Errorf("failed to create build matrix: %w", err)
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

	wg := &sync.WaitGroup{}

	// resolve hooks for whole task
	options := dukkha.TaskExecOptions{
		UseShell:        tool.UseShell(),
		ShellName:       tool.ShellName(),
		ToolCmd:         tool.GetCmd(),
		ContinueOnError: !ctx.FailFast(),
	}

	hookBefore, err := task.GetHookExecSpecs(
		ctx, dukkha.StageBefore, options,
	)
	if err != nil {
		return err
	}

	hookAfterSuccess, err := task.GetHookExecSpecs(
		ctx, dukkha.StageAfterSuccess, options,
	)
	if err != nil {
		return err
	}

	hookAfterFailure, err := task.GetHookExecSpecs(
		ctx, dukkha.StageAfterFailure, options,
	)
	if err != nil {
		return err
	}

	hookAfter, err := task.GetHookExecSpecs(
		ctx, dukkha.StageAfter, options,
	)
	if err != nil {
		return err
	}

	// ensure hook `after` always run
	defer func() {
		// TODO: handle hook error
		_ = runHook(ctx, dukkha.StageAfter, hookAfter)
	}()

	// run hook `before`
	err = runHook(ctx, dukkha.StageBefore, hookBefore)
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

		// mCtx is the matrix execution context

		mCtx := ctx.DeriveNew()
		mCtx.SetMatrixFilter(mFilter)

		select {
		case <-mCtx.Done():
			break matrixRun
		case <-waitCh:
		}

		err2 := func() error {
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

			output.WriteTaskStart(mCtx.PrefixColor(), tool.Key(), task.Key(), ms)

			// tool may have reference to MATRIX_ values
			err3 := tool.ResolveFields(mCtx, -1, "")
			if err3 != nil {
				return fmt.Errorf("failed to resolve tool fields: %w", err3)
			}

			mCtx.AddEnv(tool.GetEnv()...)

			// resolve task fields
			err3 = task.ResolveFields(mCtx, -1, "")
			if err3 != nil {
				return fmt.Errorf("failed to resolve task fields: %w", err3)
			}

			// produce a snapshot of what to do
			execSpecs, err3 := task.GetExecSpecs(mCtx, options)

			if err3 != nil {
				return fmt.Errorf("failed to generate task exec specs: %w", err3)
			}

			hookBeforeMatrix, err3 := task.GetHookExecSpecs(
				mCtx, dukkha.StageBeforeMatrix, options,
			)
			if err3 != nil {
				return err3
			}

			hookAfterMatrixSuccess, err3 := task.GetHookExecSpecs(
				mCtx, dukkha.StageAfterMatrixSuccess, options,
			)
			if err3 != nil {
				return err3
			}

			hookAfterMatrixFailure, err3 := task.GetHookExecSpecs(
				mCtx, dukkha.StageAfterMatrixFailure, options,
			)
			if err3 != nil {
				return err3
			}

			hookAfterMatrix, err3 := task.GetHookExecSpecs(
				mCtx, dukkha.StageAfterMatrix, options,
			)
			if err3 != nil {
				return err3
			}

			wg.Add(1)
			go func(ms matrix.Entry) {
				defer func() {
					// TODO: handle hook error
					_ = runHook(mCtx, dukkha.StageAfterMatrix, hookAfterMatrix)

					wg.Done()

					select {
					case waitCh <- struct{}{}:
					case <-mCtx.Done():
						return
					}
				}()

				err4 := runHook(mCtx, dukkha.StageBeforeMatrix, hookBeforeMatrix)
				if err4 != nil {
					appendErrorResult(ms, err4)
					return
				}

				err4 = doRun(mCtx, execSpecs, nil)

				output.WriteExecResult(mCtx.PrefixColor(), tool.Key(), task.Key(), ms.String(), err4)

				if err4 != nil {
					// cancel other tasks if in fail-fast mode
					if ctx.FailFast() {
						ctx.Cancel()
					}

					appendErrorResult(ms, err4)

					err4 = runHook(mCtx, dukkha.StageAfterMatrixFailure, hookAfterMatrixFailure)
					if err4 != nil {
						// TODO: handle hook error
						_ = err4
					}

					return
				}

				err4 = runHook(mCtx, dukkha.StageAfterMatrixSuccess, hookAfterMatrixSuccess)
				if err4 != nil {
					// TODO: handle hook error
					_ = err4
				}
			}(ms)

			return nil
		}()

		if err2 != nil {
			// failed before execution
			if ctx.FailFast() {
				ctx.Cancel()
			}

			appendErrorResult(ms, err2)
		}
	}

	wg.Wait()

	if len(errCollection) != 0 {
		err = runHook(ctx, dukkha.StageAfterFailure, hookAfterFailure)
		if err != nil {
			// TODO: handle hook error
			_ = err
		}

		// TODO: handle execution error
		return fmt.Errorf("task execution failed: %v", errCollection)
	}

	err = runHook(ctx, dukkha.StageAfterSuccess, hookAfterSuccess)
	if err != nil {
		// TODO: handle hook error
		return err
	}

	return nil
}
