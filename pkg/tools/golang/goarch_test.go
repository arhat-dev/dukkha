package golang

import (
	"testing"

	"arhat.dev/dukkha/pkg/constant"
	"github.com/stretchr/testify/assert"
)

func TestCreateBuildEnv(t *testing.T) {
	tests := []struct {
		mArch string

		goarch string
		goarm  string
		gomips string
	}{
		{mArch: constant.ARCH_X86, goarch: "386"},
		{mArch: constant.ARCH_AMD64, goarch: "amd64"},
		{mArch: constant.ARCH_ARM64, goarch: "arm64"},

		{mArch: constant.ARCH_ARM_V5, goarch: "arm", goarm: "5"},
		{mArch: constant.ARCH_ARM_V6, goarch: "arm", goarm: "6"},
		{mArch: constant.ARCH_ARM_V7, goarch: "arm", goarm: "7"},

		{mArch: constant.ARCH_MIPS, goarch: "mips", gomips: "hardfloat"},
		{mArch: constant.ARCH_MIPS_SF, goarch: "mips", gomips: "softfloat"},
		{mArch: constant.ARCH_MIPS_LE, goarch: "mipsle", gomips: "hardfloat"},
		{mArch: constant.ARCH_MIPS_LE_SF, goarch: "mipsle", gomips: "softfloat"},

		{mArch: constant.ARCH_MIPS64, goarch: "mips64", gomips: "hardfloat"},
		{mArch: constant.ARCH_MIPS64_SF, goarch: "mips64", gomips: "softfloat"},
		{mArch: constant.ARCH_MIPS64_LE, goarch: "mips64le", gomips: "hardfloat"},
		{mArch: constant.ARCH_MIPS64_LE_SF, goarch: "mips64le", gomips: "softfloat"},

		// ppc not supported
		// {mArch: constant.ARCH_PPC, goarch: ""},
		// {mArch: constant.ARCH_PPC_SF, goarch: ""},
		// {mArch: constant.ARCH_PPC_LE, goarch: ""},
		// {mArch: constant.ARCH_PPC_LE_SF, goarch: ""},

		{mArch: constant.ARCH_PPC64, goarch: "ppc64"},
		{mArch: constant.ARCH_PPC64_LE, goarch: "ppc64le"},

		{mArch: constant.ARCH_RISCV_64, goarch: "riscv64"},

		{mArch: constant.ARCH_S390X, goarch: "s390x"},

		// {mArch: constant.ARCH_IA64, goarch: "ia64"},
	}

	for _, test := range tests {
		t.Run(test.mArch, func(t *testing.T) {

			expected := []string{
				"GOOS=linux",
				"GOARCH=" + test.goarch,
			}

			if len(test.goarm) != 0 {
				expected = append(expected, "GOARM="+test.goarm)
			}
			if len(test.gomips) != 0 {
				expected = append(expected, "GOMIPS="+test.gomips, "GOMIPS64="+test.gomips)
			}

			assert.Equal(t, expected, createBuildEnv("linux", test.mArch))
		})
	}
}
