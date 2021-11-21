package golang

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
)

func TestCreateBuildEnv(t *testing.T) {
	goosTests := []struct {
		mKernel string

		goos string
	}{
		{mKernel: "some-custom-goos", goos: "some-custom-goos"},
	}

	for _, test := range goosTests {
		t.Run(test.mKernel, func(t *testing.T) {
			expected := dukkha.Env{
				{Name: "GOOS", Value: test.goos},
				{Name: "GOARCH", Value: "amd64"},
				{Name: "CGO_ENABLED", Value: "0"},
			}

			rc := dukkha_test.NewTestContext(context.TODO())
			rc.(di.CacheDirSetter).SetCacheDir(t.TempDir())

			rc.AddEnv(true, &dukkha.EnvEntry{
				Name:  constant.ENV_MATRIX_KERNEL,
				Value: test.mKernel,
			}, &dukkha.EnvEntry{
				Name:  constant.ENV_MATRIX_ARCH,
				Value: constant.ARCH_AMD64,
			})

			assert.EqualValues(t, expected, createBuildEnv(rc, CGOSepc{}))
		})
	}

	goarchTests := []struct {
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

		{mArch: "some-custom-goarch", goarch: "some-custom-goarch"},

		// {mArch: constant.ARCH_IA64, goarch: "ia64"},
	}

	for _, test := range goarchTests {
		t.Run(test.mArch, func(t *testing.T) {
			expected := dukkha.Env{
				{
					Name:  "GOOS",
					Value: constant.KERNEL_LINUX,
				},
				{
					Name:  "GOARCH",
					Value: test.goarch,
				},
			}

			if len(test.goarm) != 0 {
				expected = append(expected, &dukkha.EnvEntry{
					Name:  "GOARM",
					Value: test.goarm,
				})
			}
			if len(test.gomips) != 0 {
				expected = append(expected, &dukkha.EnvEntry{
					Name:  "GOMIPS",
					Value: test.gomips,
				}, &dukkha.EnvEntry{
					Name:  "GOMIPS64",
					Value: test.gomips,
				})
			}

			expected = append(expected, &dukkha.EnvEntry{
				Name:  "CGO_ENABLED",
				Value: "0",
			})

			rc := dukkha_test.NewTestContext(context.TODO())
			rc.(di.CacheDirSetter).SetCacheDir(t.TempDir())
			rc.AddEnv(true, &dukkha.EnvEntry{
				Name:  constant.ENV_MATRIX_KERNEL,
				Value: constant.KERNEL_LINUX,
			}, &dukkha.EnvEntry{
				Name:  constant.ENV_MATRIX_ARCH,
				Value: test.mArch,
			})

			assert.Equal(t, expected, createBuildEnv(
				rc, CGOSepc{},
			))
		})
	}
}
