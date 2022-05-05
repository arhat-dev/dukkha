package golang

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"arhat.dev/pkg/archconst"

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
)

func TestCreateBuildEnv(t *testing.T) {
	t.Parallel()

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
				Value: archconst.ARCH_AMD64,
			})

			assert.EqualValues(t, expected, createBuildEnv(rc, CGOSepc{}))
		})
	}

	goarchTests := []struct {
		mArch string

		goarch  string
		goarm   string
		gomips  string
		goamd64 string
		goppc64 string
	}{
		{mArch: archconst.ARCH_X86, goarch: "386"},
		{mArch: archconst.ARCH_AMD64, goarch: "amd64"},

		{mArch: archconst.ARCH_AMD64_V1, goarch: "amd64", goamd64: "v1"},
		{mArch: archconst.ARCH_AMD64_V2, goarch: "amd64", goamd64: "v2"},
		{mArch: archconst.ARCH_AMD64_V3, goarch: "amd64", goamd64: "v3"},
		{mArch: archconst.ARCH_AMD64_V4, goarch: "amd64", goamd64: "v4"},

		{mArch: archconst.ARCH_ARM64, goarch: "arm64"},

		{mArch: archconst.ARCH_ARM_V5, goarch: "arm", goarm: "5"},
		{mArch: archconst.ARCH_ARM_V6, goarch: "arm", goarm: "6"},
		{mArch: archconst.ARCH_ARM_V7, goarch: "arm", goarm: "7"},

		{mArch: archconst.ARCH_MIPS, goarch: "mips", gomips: "hardfloat"},
		{mArch: archconst.ARCH_MIPS_SF, goarch: "mips", gomips: "softfloat"},
		{mArch: archconst.ARCH_MIPS_LE, goarch: "mipsle", gomips: "hardfloat"},
		{mArch: archconst.ARCH_MIPS_LE_SF, goarch: "mipsle", gomips: "softfloat"},

		{mArch: archconst.ARCH_MIPS64, goarch: "mips64", gomips: "hardfloat"},
		{mArch: archconst.ARCH_MIPS64_SF, goarch: "mips64", gomips: "softfloat"},
		{mArch: archconst.ARCH_MIPS64_LE, goarch: "mips64le", gomips: "hardfloat"},
		{mArch: archconst.ARCH_MIPS64_LE_SF, goarch: "mips64le", gomips: "softfloat"},

		// ppc not supported
		// {mArch: archconst.ARCH_PPC, goarch: ""},
		// {mArch: archconst.ARCH_PPC_SF, goarch: ""},
		// {mArch: archconst.ARCH_PPC_LE, goarch: ""},
		// {mArch: archconst.ARCH_PPC_LE_SF, goarch: ""},

		{mArch: archconst.ARCH_PPC64, goarch: "ppc64"},
		{mArch: archconst.ARCH_PPC64_LE, goarch: "ppc64le"},
		{mArch: archconst.ARCH_PPC64_V8, goarch: "ppc64", goppc64: "power8"},
		{mArch: archconst.ARCH_PPC64_LE_V8, goarch: "ppc64le", goppc64: "power8"},
		{mArch: archconst.ARCH_PPC64_V9, goarch: "ppc64", goppc64: "power9"},
		{mArch: archconst.ARCH_PPC64_LE_V9, goarch: "ppc64le", goppc64: "power9"},

		{mArch: archconst.ARCH_RISCV_64, goarch: "riscv64"},

		{mArch: archconst.ARCH_S390X, goarch: "s390x"},

		{mArch: "some-custom-goarch", goarch: "some-custom-goarch"},

		// {mArch: archconst.ARCH_IA64, goarch: "ia64"},
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

			if len(test.goamd64) != 0 {
				expected = append(expected, &dukkha.EnvEntry{
					Name:  "GOAMD64",
					Value: test.goamd64,
				})
			}

			if len(test.goppc64) != 0 {
				expected = append(expected, &dukkha.EnvEntry{
					Name:  "GOPPC64",
					Value: test.goppc64,
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
