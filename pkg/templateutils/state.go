package templateutils

import "arhat.dev/dukkha/pkg/dukkha"

func createStateNS(rc dukkha.RenderingContext) *stateNS {
	return &stateNS{ctx: rc}
}

type stateNS struct {
	ctx dukkha.RenderingContext
}

func (s *stateNS) Succeeded() bool {
	return s.ctx.(dukkha.TaskExecContext).State() == dukkha.TaskExecSucceeded
}

func (s *stateNS) Failed() bool {
	return s.ctx.(dukkha.TaskExecContext).State() == dukkha.TaskExecFailed
}
