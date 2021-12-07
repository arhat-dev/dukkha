//go:build js || illumos || ios || plan9
// +build js illumos ios plan9

package cpuhelper

func Arch() string {
	return ArchByCPUFeatures()
}
