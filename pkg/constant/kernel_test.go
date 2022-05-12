package constant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// TODO: update this map when changing kernel.go
	requiredKernelMappingValues = map[string]string{
		KERNEL_Windows:    "windows",
		KERNEL_Linux:      "linux",
		KERNEL_Darwin:     "darwin",
		KERNEL_FreeBSD:    "freebsd",
		KERNEL_NetBSD:     "netbsd",
		KERNEL_OpenBSD:    "openbsd",
		KERNEL_Solaris:    "solaris",
		KERNEL_Illumos:    "illumos",
		KERNEL_JavaScript: "js",
		KERNEL_Aix:        "aix",
		KERNEL_Android:    "android",
		KERNEL_iOS:        "ios",
		KERNEL_Plan9:      "plan9",
	}
)

func TestKernelMapping(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		mappingFunc func(mKernel string) (string, bool)
	}{
		{
			name:        "DockerOS",
			mappingFunc: GetDockerOS,
		},
		{
			name:        "Golang OS",
			mappingFunc: GetGolangOS,
		},
		{
			name:        "OCI OS",
			mappingFunc: GetOciOS,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for mKernel := range requiredKernelMappingValues {
				_, ok := test.mappingFunc(mKernel)
				assert.True(t, ok, mKernel)
			}
		})
	}
}
