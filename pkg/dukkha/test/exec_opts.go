package dukkha_test

import "arhat.dev/dukkha/pkg/dukkha"

func CreateTaskMatrixExecOptions() dukkha.TaskMatrixExecOptions {
	opts := dukkha.CreateTaskExecOptions(1, 1)
	return opts.NextMatrixExecOptions(false, "")
}
