package templateutils

import "arhat.dev/dukkha/pkg/dukkha"

func createStateNS(rc dukkha.RenderingContext) *_stateNS {
	return &_stateNS{ctx: rc}
}

type _stateNS struct {
	ctx dukkha.RenderingContext
}

func (s *_stateNS) Succeeded() bool {
	return s.ctx.(dukkha.TaskExecContext).State() == dukkha.TaskExecSucceeded
}

func (s *_stateNS) Failed() bool {
	return s.ctx.(dukkha.TaskExecContext).State() == dukkha.TaskExecFailed
}
