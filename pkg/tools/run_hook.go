package tools

import (
	"fmt"

	"arhat.dev/dukkha/pkg/dukkha"
)

func runHook(
	ctx dukkha.TaskExecContext,
	stage dukkha.TaskExecStage,
	specs [][]dukkha.TaskExecSpec,
) error {
	for i, execSpec := range specs {
		err := doRun(ctx, execSpec, nil)
		if err != nil {
			return fmt.Errorf(
				"hook %q action #%d failed: %w",
				stage.String(), i, err,
			)
		}
	}

	return nil
}
