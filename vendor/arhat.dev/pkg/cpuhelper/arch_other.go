//go:build js || illumos || ios || plan9
// +build js illumos ios plan9

package cpuhelper

import "arhat.dev/pkg/archconst"

func Arch(cpu CPU) archconst.ArchValue {
	return ArchByCPUFeatures(cpu)
}
