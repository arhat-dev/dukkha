package templateutils

import "arhat.dev/dukkha/pkg/dukkha"

func createStateNS(rc dukkha.RenderingContext) stateNS { return stateNS{rc: rc} }

type stateNS struct{ rc dukkha.RenderingContext }

func (ns stateNS) Succeeded() bool {
	return ns.rc.(dukkha.TaskExecContext).State() == dukkha.TaskExecSucceeded
}

func (ns stateNS) Failed() bool {
	return ns.rc.(dukkha.TaskExecContext).State() == dukkha.TaskExecFailed
}
