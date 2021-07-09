package tools

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/matrix"
	"arhat.dev/dukkha/pkg/output"
	"go.uber.org/multierr"
)

type TaskExecSpecWithContext struct {
	Matrix  matrix.Entry
	Context dukkha.TaskExecContext

	HookBeofreMatrix []dukkha.RunTaskOrRunShell

	Specs []dukkha.TaskExecSpec

	HookAfterMatrixSuccess []dukkha.RunTaskOrRunShell
	HookAfterMatrixFailure []dukkha.RunTaskOrRunShell
	HookAfterMatrix        []dukkha.RunTaskOrRunShell
}

type CompleteTaskExecSpecs struct {
	Context    dukkha.TaskExecContext
	CancelTask context.CancelFunc

	HookBefore []dukkha.RunTaskOrRunShell

	TaskExec []TaskExecSpecWithContext

	HookAfterSuccess []dukkha.RunTaskOrRunShell
	hookAfterFailure []dukkha.RunTaskOrRunShell
	HookAfter        []dukkha.RunTaskOrRunShell
}

func GenCompleteTaskExecSpecs(
	_ctx dukkha.TaskExecContext,
	tool dukkha.Tool,
	task dukkha.Task,
) (*CompleteTaskExecSpecs, error) {
	ctx := _ctx.DeriveNew()

	matrixSpecs, err := task.GetMatrixSpecs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get execution matrix: %w", err)
	}

	if len(matrixSpecs) == 0 {
		return nil, fmt.Errorf("no matrix spec match")
	}

	// resolve hooks for whole task
	options := dukkha.TaskExecOptions{
		UseShell:        tool.UseShell(),
		ShellName:       tool.ShellName(),
		ToolCmd:         tool.GetCmd(),
		ContinueOnError: !ctx.FailFast(),
	}

	var err2 error

	ret := &CompleteTaskExecSpecs{
		Context:    ctx,
		CancelTask: ctx.Cancel,
	}

	ret.Context.SetTask(tool.Key(), task.Key())

	ret.HookBefore, err2 = task.GetHookExecSpecs(
		ctx, dukkha.StageBefore, options,
	)
	err = multierr.Append(err, err2)

	ret.HookAfterSuccess, err2 = task.GetHookExecSpecs(
		ctx, dukkha.StageAfterSuccess, options,
	)
	err = multierr.Append(err, err2)

	ret.hookAfterFailure, err2 = task.GetHookExecSpecs(
		ctx, dukkha.StageAfterFailure, options,
	)
	err = multierr.Append(err, err2)

	ret.HookAfter, err2 = task.GetHookExecSpecs(
		ctx, dukkha.StageAfter, options,
	)
	err = multierr.Append(err, err2)

	if err != nil {
		return nil, fmt.Errorf("failed to get task hooks exec specs: %w", err)
	}

	for i, ms := range matrixSpecs {
		// set default matrix filter for referenced hook tasks
		mFilter := make(map[string][]string)
		for k, v := range ms {
			mFilter[k] = []string{v}
		}

		// mCtx is the matrix execution context

		mCtx := ctx.DeriveNew()
		mCtx.SetMatrixFilter(mFilter)

		mSpec := &TaskExecSpecWithContext{
			Matrix:  ms,
			Context: mCtx,
		}

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
		err = tool.ResolveFields(mCtx, -1, "")
		if err != nil {
			return nil, fmt.Errorf("failed to resolve tool fields: %w", err)
		}

		mCtx.AddEnv(tool.GetEnv()...)

		// resolve task fields
		err = task.ResolveFields(mCtx, -1, "")
		if err != nil {
			return nil, fmt.Errorf("failed to resolve task fields: %w", err)
		}

		// produce a snapshot of what to do
		mSpec.Specs, err = task.GetExecSpecs(mCtx, options)
		if err != nil {
			return nil, fmt.Errorf("failed to generate task exec specs: %w", err)
		}

		mSpec.HookBeofreMatrix, err2 = task.GetHookExecSpecs(
			mCtx, dukkha.StageBeforeMatrix, options,
		)
		err = multierr.Append(err, err2)

		mSpec.HookAfterMatrixSuccess, err2 = task.GetHookExecSpecs(
			mCtx, dukkha.StageAfterMatrixSuccess, options,
		)
		err = multierr.Append(err, err2)

		mSpec.HookAfterMatrixFailure, err2 = task.GetHookExecSpecs(
			mCtx, dukkha.StageAfterMatrixFailure, options,
		)
		err = multierr.Append(err, err2)

		mSpec.HookAfterMatrix, err2 = task.GetHookExecSpecs(
			mCtx, dukkha.StageAfterMatrix, options,
		)
		err = multierr.Append(err, err2)

		if err != nil {
			return nil, fmt.Errorf("failed to get task matrix hooks exec specs: %w", err)
		}

		ret.TaskExec = append(ret.TaskExec, *mSpec)
	}

	return ret, nil
}

// nolint:gocyclo
func RunTask(specs *CompleteTaskExecSpecs) error {
	defer specs.CancelTask()

	// TODO: do real global limit
	workerCount := specs.Context.ClaimWorkers(len(specs.TaskExec))
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

	// ensure hook `after` always run
	defer func() {
		// TODO: handle hook error
		_ = runHook(specs.Context, dukkha.StageAfter, specs.HookAfter)
	}()

	// run hook `before`
	err := runHook(specs.Context, dukkha.StageBefore, specs.HookBefore)
	if err != nil {
		// cancel task execution
		return err
	}

matrixRun:
	for _, mSpec := range specs.TaskExec {
		select {
		case <-mSpec.Context.Done():
			break matrixRun
		case <-waitCh:
		}

		output.WriteTaskStart(
			mSpec.Context.PrefixColor(),
			mSpec.Context.CurrentTool(),
			mSpec.Context.CurrentTask(),
			mSpec.Matrix,
		)

		wg.Add(1)
		go func(mSpec TaskExecSpecWithContext) {
			defer func() {
				// TODO: handle hook error
				_ = runHook(
					mSpec.Context,
					dukkha.StageAfterMatrix,
					mSpec.HookAfterMatrix,
				)

				wg.Done()

				select {
				case waitCh <- struct{}{}:
				case <-mSpec.Context.Done():
					return
				}
			}()

			err2 := runHook(
				mSpec.Context,
				dukkha.StageBeforeMatrix,
				mSpec.HookBeofreMatrix,
			)
			if err2 != nil {
				appendErrorResult(mSpec.Matrix, err2)
				return
			}

			err2 = doRun(mSpec.Context, mSpec.Specs, nil)

			output.WriteExecResult(
				mSpec.Context.PrefixColor(),
				mSpec.Context.CurrentTool(),
				mSpec.Context.CurrentTask(),
				mSpec.Matrix.String(),
				err2,
			)

			if err2 != nil {
				// cancel other tasks if in fail-fast mode
				if specs.Context.FailFast() {
					specs.CancelTask()
				}

				appendErrorResult(mSpec.Matrix, err2)

				err2 = runHook(
					mSpec.Context,
					dukkha.StageAfterMatrixFailure,
					mSpec.HookAfterMatrixFailure,
				)
				if err2 != nil {
					// TODO: handle hook error
					_ = err2
				}

				return
			}

			err2 = runHook(
				mSpec.Context,
				dukkha.StageAfterMatrixSuccess,
				mSpec.HookAfterMatrixSuccess,
			)
			if err2 != nil {
				// TODO: handle hook error
				_ = err2
			}
		}(mSpec)
	}

	wg.Wait()

	if len(errCollection) != 0 {
		err = runHook(specs.Context, dukkha.StageAfterFailure, specs.hookAfterFailure)
		if err != nil {
			// TODO: handle hook error
			_ = err
		}

		// TODO: handle execution error
		return fmt.Errorf("task execution failed: %v", errCollection)
	}

	err = runHook(specs.Context, dukkha.StageAfterSuccess, specs.HookAfterSuccess)
	if err != nil {
		// TODO: handle hook error
		return err
	}

	return nil
}
