package templateutils

import (
	"arhat.dev/tlang"

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/dukkha"
)

func createMiscNS(rc dukkha.RenderingContext) miscNS {
	return miscNS{rc: rc}
}

type miscNS struct {
	rc dukkha.RenderingContext
}

func (ns miscNS) Git() map[string]tlang.LazyValueType[string]  { return ns.rc.GitValues() }
func (ns miscNS) Host() map[string]tlang.LazyValueType[string] { return ns.rc.HostValues() }
func (ns miscNS) Env() map[string]tlang.LazyValueType[string]  { return ns.rc.Env() }
func (ns miscNS) Values() map[string]any                       { return ns.rc.Values() }

func (ns miscNS) Matrix() map[string]string {
	mf := ns.rc.MatrixFilter()
	return mf.AsEntry()
}

// for transform renderer
func (ns miscNS) VALUE() any {
	vg, ok := ns.rc.(di.VALUEGetter)
	if ok {
		return vg.VALUE()
	}

	return nil
}
