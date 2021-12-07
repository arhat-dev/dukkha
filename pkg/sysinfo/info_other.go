//go:build js || illumos || ios || plan9
// +build js illumos ios plan9

package sysinfo

func OSName() string {
	return ""
}

func OSVersion() string {
	return ""
}

func KernelVersion() string {
	return ""
}
