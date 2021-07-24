package constant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// TODO: update this map when changing kernel.go
	requiredKernelMappingValues = map[string]string{
		KERNEL_WINDOWS:    "windows",
		KERNEL_LINUX:      "linux",
		KERNEL_DARWIN:     "darwin",
		KERNEL_FREEBSD:    "freebsd",
		KERNEL_NETBSD:     "netbsd",
		KERNEL_OPENBSD:    "openbsd",
		KERNEL_SOLARIS:    "solaris",
		KERNEL_ILLUMOS:    "illumos",
		KERNEL_JAVASCRIPT: "js",
		KERNEL_AIX:        "aix",
		KERNEL_ANDROID:    "android",
		KERNEL_IOS:        "ios",
		KERNEL_PLAN9:      "plan9",
	}
)

func TestKernelMapping(t *testing.T) {
	tests := []struct {
		name        string
		mappingFunc func(mKernel string) (string, bool)
	}{
		{
			name:        "DockerOS",
			mappingFunc: GetDockerOS,
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
