package tools

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"go.uber.org/multierr"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/matrix"
	"arhat.dev/dukkha/pkg/output"
)

type TaskExecRequest struct {
	Context dukkha.TaskExecContext

	Tool dukkha.Tool
	Task dukkha.Task

	IgnoreError bool

	// DryRun do not actually run any thing, just evaluate values
	DryRun bool
}

// nolint:gocyclo
func RunTask(req *TaskExecRequest) (err error) {
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

		res := &taskResult{
			errMsg: err.Error(),
		}
		if spec != nil {
			res.matrixSpec = spec.BriefString()
		}

		errCollection = append(errCollection, *res)
	}

	// task may need tool specific env, resolve tool env first

	toolCmd := func(ctx dukkha.RenderingContext) ([]string, error) {
		var ret []string
		err2 := req.Tool.DoAfterFieldsResolved(ctx, -1, false, func() error {
			ret = req.Tool.GetCmd()
			return nil
		}, "BaseTool.cmd")
		if err2 != nil {
			return nil, err2
		}

		return ret, nil
	}

	// resolve tool to set tool specific context env
	err = req.Tool.DoAfterFieldsResolved(req.Context, -1, true, func() error {
		return nil
	}, "BaseTool.env")
	if err != nil {
		return fmt.Errorf("resolving tool specific env: %w", err)
	}

	// resolve hooks for whole task

	req.Context.SetTask(req.Tool.Key(), req.Task.Key())

	wg := &sync.WaitGroup{}

	unstoppableTaskCtx := req.Context.WithCustomParent(context.Background())
	// ensure hook `after` always run
	defer func() {
		// TODO: handle hook error
		hookAfter, err2 := req.Task.GetHookExecSpecs(
			unstoppableTaskCtx, dukkha.StageAfter,
		)

		if err2 != nil {
			appendErrorResult(nil, err2)
		} else {
			err2 = doRun(unstoppableTaskCtx, toolCmd, hookAfter, nil)
			if err2 != nil {
				appendErrorResult(nil, err2)
			}
		}

		if len(errCollection) != 0 {
			err2 := fmt.Errorf("%v", errCollection)
			if err != nil {
				err = multierr.Append(err, err2)
			} else {
				err = err2
			}
		}
	}()

	// run hook `before`
	hookBefore, err := req.Task.GetHookExecSpecs(
		unstoppableTaskCtx, dukkha.StageBefore,
	)
	if err != nil {
		// cancel task execution
		return err
	}

	err = doRun(unstoppableTaskCtx, toolCmd, hookBefore, nil)
	if err != nil {
		// cancel task execution
		return err
	}

	matrixSpecs, err := req.Task.GetMatrixSpecs(req.Context)
	if err != nil {
		return fmt.Errorf("creating execution matrix: %w", err)
	}

	if len(matrixSpecs) == 0 {
		// TODO: write warning and ignore error
		return fmt.Errorf("no matrix spec match")
	}

	// TODO: do real global limit
	workerCount := req.Context.ClaimWorkers(len(matrixSpecs))
	waitCh := make(chan struct{}, workerCount)
	for i := 0; i < workerCount; i++ {
		waitCh <- struct{}{}
	}

	// TODO: alloc real task exec id
	opts := dukkha.CreateTaskExecOptions(0, len(matrixSpecs))
matrixRun:
	for _, ms := range matrixSpecs {
		mCtx, options, err2 := CreateTaskMatrixContext(req, ms, opts)

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

		output.WriteTaskStart(mCtx.PrefixColor(),
			mCtx.CurrentTool(), mCtx.CurrentTask(), ms,
		)

		wg.Add(1)

		unstoppableMatrixCtx := mCtx.WithCustomParent(context.Background())
		go func(ms matrix.Entry) {
			var (
				err3 error
			)

			toolMatrixCmd := func(ctx dukkha.RenderingContext) ([]string, error) {
				var ret []string
				err4 := req.Tool.DoAfterFieldsResolved(ctx, -1, false, func() error {
					ret = req.Tool.GetCmd()
					return nil
				}, "BaseTool.cmd")
				if err4 != nil {
					return nil, err4
				}

				return ret, nil
			}

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
					unstoppableMatrixCtx, dukkha.StageAfterMatrix,
				)
				if err4 != nil {
					appendErrorResult(ms, err4)
				} else {
					// TODO: handle hook error
					err4 = doRun(unstoppableMatrixCtx, toolCmd, hookAfterMatrix, nil)
					if err4 != nil {
						appendErrorResult(ms, err4)
					}
				}
			}()

			hookBeofreMatrix, err3 := req.Task.GetHookExecSpecs(
				unstoppableMatrixCtx, dukkha.StageBeforeMatrix,
			)
			if err3 != nil {
				appendErrorResult(ms, err3)
				return
			}

			err3 = doRun(unstoppableMatrixCtx, toolCmd, hookBeofreMatrix, nil)
			if err3 != nil {
				appendErrorResult(ms, err3)
				return
			}

			// produce a snapshot of what to do
			execSpecs, err3 := req.Task.GetExecSpecs(mCtx, options)
			if err3 != nil {
				appendErrorResult(
					ms,
					fmt.Errorf("generating task exec specs: %w", err3),
				)
				return
			}

			err3 = doRun(mCtx, toolMatrixCmd, execSpecs, nil)

			output.WriteExecResult(mCtx.PrefixColor(),
				mCtx.CurrentTool(), mCtx.CurrentTask(),
				ms.String(), err3,
			)

			if err3 != nil {
				// cancel other tasks if in fail-fast mode
				if req.Context.FailFast() {
					req.Context.Cancel()
				}

				appendErrorResult(ms, err3)

				hookAfterMatrixFailure, err4 := req.Task.GetHookExecSpecs(
					unstoppableMatrixCtx, dukkha.StageAfterMatrixFailure,
				)
				if err4 != nil {
					appendErrorResult(ms, err4)
				} else {
					err4 = doRun(unstoppableMatrixCtx, toolMatrixCmd, hookAfterMatrixFailure, nil)
					if err4 != nil {
						appendErrorResult(ms, err4)
					}
				}

				return
			}

			hookAfterMatrixSuccess, err3 := req.Task.GetHookExecSpecs(
				unstoppableMatrixCtx, dukkha.StageAfterMatrixSuccess,
			)
			if err3 != nil {
				appendErrorResult(ms, err3)
			} else {
				err3 = doRun(unstoppableMatrixCtx, toolMatrixCmd, hookAfterMatrixSuccess, nil)
				if err3 != nil {
					appendErrorResult(ms, err3)
				}
			}
		}(ms)
	}

	wg.Wait()

	if len(errCollection) != 0 {
		hookAfterFailure, err2 := req.Task.GetHookExecSpecs(
			unstoppableTaskCtx, dukkha.StageAfterFailure,
		)
		if err2 != nil {
			appendErrorResult(nil, err2)
			return
		}

		err2 = doRun(unstoppableTaskCtx, toolCmd, hookAfterFailure, nil)
		if err2 != nil {
			appendErrorResult(nil, err2)
		}

		return
	}

	hookAfterSuccess, err := req.Task.GetHookExecSpecs(
		unstoppableTaskCtx, dukkha.StageAfterSuccess,
	)
	if err != nil {
		appendErrorResult(nil, err)
		return
	}

	err = doRun(unstoppableTaskCtx, toolCmd, hookAfterSuccess, nil)
	if err != nil {
		appendErrorResult(nil, err)
		return
	}

	return
}

// CreateTaskMatrixContext creates a per matrix entry task exec options
// with context resolved
func CreateTaskMatrixContext(
	req *TaskExecRequest,
	ms matrix.Entry,
	opts dukkha.TaskExecOptions,
) (dukkha.TaskExecContext, dukkha.TaskMatrixExecOptions, error) {
	mCtx := req.Context.DeriveNew()

	// set default matrix filter for referenced hook tasks
	mFilter := matrix.NewFilter(nil)
	for k, v := range ms {
		mFilter.AddMatch(k, v)
	}

	mCtx.SetMatrixFilter(mFilter)

	for k, v := range ms {
		name := "MATRIX_" + strings.ToUpper(k)
		mCtx.AddEnv(true, &dukkha.EnvEntry{
			Name:  name,
			Value: v,
		})

		if name == constant.ENV_MATRIX_ARCH {
			mCtx.AddEnv(true, &dukkha.EnvEntry{
				Name:  constant.ENV_MATRIX_ARCH_SIMPLE,
				Value: constant.SimpleArch(v),
			})
		}
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

	// tool may have reference to MATRIX_ values
	// but MUST not have reference to task specific env

	// now everything prepared for the tool, resolve all of it

	var options dukkha.TaskMatrixExecOptions
	err := req.Tool.DoAfterFieldsResolved(mCtx, -1, true, func() error {
		options = opts.NextMatrixExecOptions()

		mCtx.SetTaskColors(output.PickColor(options.Seq()))

		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("add tool env: %w", err)
	}

	return mCtx, options, nil
}
