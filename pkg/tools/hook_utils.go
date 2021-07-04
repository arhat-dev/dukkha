package tools

import "fmt"

func HandleHookRunError(stage TaskExecStage, err error) error {
	if err != nil {
		return fmt.Errorf("hook `%s` failed: %w", stage.String(), err)
	}

	return nil
}
