package golang

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"arhat.dev/pkg/archconst"

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
			expected := dukkha.NameValueList{
				{Name: "GOOS", Value: test.goos},
				{Name: "GOARCH", Value: "s390x"},
				{Name: "CGO_ENABLED", Value: "0"},
			}

			rc := dukkha_test.NewTestContext(context.TODO(), t.TempDir())

			rc.AddEnv(true, &dukkha.NameValueEntry{
				Name:  constant.EnvName_MATRIX_KERNEL,
				Value: test.mKernel,
			}, &dukkha.NameValueEntry{
				Name:  constant.EnvName_MATRIX_ARCH,
				Value: archconst.ARCH_S390X,
			})

			assert.EqualValues(t, expected, createBuildEnv(rc, buildOptions{}, CGOSepc{}))
		})
	}

	goarchTests := []struct {
		mArch string

		goarch string

		go386    string
		goamd64  string
		goarm64  string
		goarm    string
		gomips   string
		gomips64 string
		goppc64  string
	}{
		{mArch: archconst.ARCH_X86, goarch: "386", go386: "sse2"},
		{mArch: archconst.ARCH_X86_SF, goarch: "386", go386: "softfloat"},

		{mArch: archconst.ARCH_AMD64, goarch: "amd64", goamd64: "v1"},
		{mArch: archconst.ARCH_AMD64_V1, goarch: "amd64", goamd64: "v1"},
		{mArch: archconst.ARCH_AMD64_V2, goarch: "amd64", goamd64: "v2"},
		{mArch: archconst.ARCH_AMD64_V3, goarch: "amd64", goamd64: "v3"},
		{mArch: archconst.ARCH_AMD64_V4, goarch: "amd64", goamd64: "v4"},

		{mArch: archconst.ARCH_ARM64, goarch: "arm64", goarm64: "8"},
		{mArch: archconst.ARCH_ARM64_V8, goarch: "arm64", goarm64: "8"},
		{mArch: archconst.ARCH_ARM64_V9, goarch: "arm64", goarm64: "9"},

		{mArch: archconst.ARCH_ARM, goarch: "arm", goarm: "7"},
		{mArch: archconst.ARCH_ARM_V5, goarch: "arm", goarm: "5"},
		{mArch: archconst.ARCH_ARM_V6, goarch: "arm", goarm: "6"},
		{mArch: archconst.ARCH_ARM_V7, goarch: "arm", goarm: "7"},

		{mArch: archconst.ARCH_MIPS, goarch: "mips", gomips: "hardfloat"},
		{mArch: archconst.ARCH_MIPS_SF, goarch: "mips", gomips: "softfloat"},
		{mArch: archconst.ARCH_MIPS_LE, goarch: "mipsle", gomips: "hardfloat"},
		{mArch: archconst.ARCH_MIPS_LE_SF, goarch: "mipsle", gomips: "softfloat"},

		{mArch: archconst.ARCH_MIPS64, goarch: "mips64", gomips64: "hardfloat"},
		{mArch: archconst.ARCH_MIPS64_SF, goarch: "mips64", gomips64: "softfloat"},
		{mArch: archconst.ARCH_MIPS64_LE, goarch: "mips64le", gomips64: "hardfloat"},
		{mArch: archconst.ARCH_MIPS64_LE_SF, goarch: "mips64le", gomips64: "softfloat"},

		// ppc not supported
		// {mArch: archconst.ARCH_PPC, goarch: ""},
		// {mArch: archconst.ARCH_PPC_SF, goarch: ""},
		// {mArch: archconst.ARCH_PPC_LE, goarch: ""},
		// {mArch: archconst.ARCH_PPC_LE_SF, goarch: ""},

		{mArch: archconst.ARCH_PPC64, goarch: "ppc64", goppc64: "power8"},
		{mArch: archconst.ARCH_PPC64_LE, goarch: "ppc64le", goppc64: "power8"},
		{mArch: archconst.ARCH_PPC64_V8, goarch: "ppc64", goppc64: "power8"},
		{mArch: archconst.ARCH_PPC64_LE_V8, goarch: "ppc64le", goppc64: "power8"},
		{mArch: archconst.ARCH_PPC64_V9, goarch: "ppc64", goppc64: "power9"},
		{mArch: archconst.ARCH_PPC64_LE_V9, goarch: "ppc64le", goppc64: "power9"},

		{mArch: archconst.ARCH_RISCV64, goarch: "riscv64"},

		{mArch: archconst.ARCH_S390X, goarch: "s390x"},

		{mArch: "some-custom-goarch", goarch: "some-custom-goarch"},

		// {mArch: archconst.ARCH_IA64, goarch: "ia64"},
	}

	for _, test := range goarchTests {
		t.Run(test.mArch, func(t *testing.T) {
			expected := dukkha.NameValueList{
				{
					Name:  "GOOS",
					Value: constant.KERNEL_Linux,
				},
				{
					Name:  "GOARCH",
					Value: test.goarch,
				},
			}

			if len(test.go386) != 0 {
				expected = append(expected, &dukkha.NameValueEntry{
					Name:  "GO386",
					Value: test.go386,
				})
			}

			if len(test.goamd64) != 0 {
				expected = append(expected, &dukkha.NameValueEntry{
					Name:  "GOAMD64",
					Value: test.goamd64,
				})
			}

			if len(test.goarm64) != 0 {
				expected = append(expected, &dukkha.NameValueEntry{
					Name:  "GOARM64",
					Value: test.goarm64,
				})
			}

			if len(test.goarm) != 0 {
				expected = append(expected, &dukkha.NameValueEntry{
					Name:  "GOARM",
					Value: test.goarm,
				})
			}

			if len(test.gomips) != 0 {
				expected = append(expected, &dukkha.NameValueEntry{
					Name:  "GOMIPS",
					Value: test.gomips,
				})
			}

			if len(test.gomips64) != 0 {
				expected = append(expected, &dukkha.NameValueEntry{
					Name:  "GOMIPS64",
					Value: test.gomips64,
				})
			}

			if len(test.goppc64) != 0 {
				expected = append(expected, &dukkha.NameValueEntry{
					Name:  "GOPPC64",
					Value: test.goppc64,
				})
			}

			expected = append(expected, &dukkha.NameValueEntry{
				Name:  "CGO_ENABLED",
				Value: "0",
			})

			rc := dukkha_test.NewTestContext(context.TODO(), t.TempDir())
			rc.AddEnv(true, &dukkha.NameValueEntry{
				Name:  constant.EnvName_MATRIX_KERNEL,
				Value: constant.KERNEL_Linux,
			}, &dukkha.NameValueEntry{
				Name:  constant.EnvName_MATRIX_ARCH,
				Value: test.mArch,
			})

			assert.Equal(t, expected, createBuildEnv(
				rc, buildOptions{}, CGOSepc{},
			))
		})
	}
}
