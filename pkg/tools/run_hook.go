package tools

import (
	"fmt"

	"arhat.dev/dukkha/pkg/dukkha"
)

func runHook(
	ctx dukkha.TaskExecContext,
	stage dukkha.TaskExecStage,
	specs []dukkha.RunTaskOrRunShell,
) error {
	for i, execSpec := range specs {

		var err error
		switch t := execSpec.(type) {
		case *CompleteTaskExecSpecs:
			err = RunTask(t)
		case []dukkha.TaskExecSpec:
			err = doRun(ctx, t, nil)
		default:
			return fmt.Errorf("unexpected hook run, unknown type: %T", t)
		}

		if err != nil {
			return fmt.Errorf(
				"hook %q action #%d failed: %w",
				stage.String(), i, err,
			)
		}
	}

	return nil
}
