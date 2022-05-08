package golang

import (
	"strings"

	"arhat.dev/pkg/archconst"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
)

func createBuildEnv(v dukkha.EnvValues, cgoSpec CGOSepc) dukkha.Env {
	var env dukkha.Env

	goos, _ := constant.GetGolangOS(v.MatrixKernel())
	switch {
	case len(goos) != 0:
	case len(v.MatrixKernel()) != 0:
		goos = v.MatrixKernel()
	}

	if len(goos) != 0 {
		env = append(env, &dukkha.EnvEntry{
			Name:  "GOOS",
			Value: goos,
		})
	}

	mArch := v.MatrixArch()
	goarch, _ := constant.GetGolangArch(mArch)
	switch {
	case len(goarch) != 0:
	case len(mArch) != 0:
		goarch = mArch
	}

	if len(goarch) != 0 {
		env = append(env, &dukkha.EnvEntry{
			Name:  "GOARCH",
			Value: goarch,
		})
	}

	switch {
	case strings.HasPrefix(mArch, archconst.ARCH_AMD64):
		env = append(env, &dukkha.EnvEntry{
			Name:  "GOAMD64",
			Value: getGOAMD64(mArch),
		})
	case strings.HasPrefix(mArch, archconst.ARCH_X86):
		env = append(env, &dukkha.EnvEntry{
			Name:  "GO386",
			Value: getGO386(mArch),
		})
	case strings.HasPrefix(mArch, archconst.ARCH_ARM64): // MUST be prior to ARCH_ARM
		env = append(env, &dukkha.EnvEntry{
			Name:  "GOARM64",
			Value: getGOARM64(mArch),
		})
	case strings.HasPrefix(mArch, archconst.ARCH_ARM):
		env = append(env, &dukkha.EnvEntry{
			Name:  "GOARM",
			Value: getGOARM(mArch),
		})
	case strings.HasPrefix(mArch, archconst.ARCH_MIPS64): // MUST be prior to ARCH_MIPS
		env = append(env, &dukkha.EnvEntry{
			Name:  "GOMIPS64",
			Value: getGOMIPS64(mArch),
		})
	case strings.HasPrefix(mArch, archconst.ARCH_MIPS):
		env = append(env, &dukkha.EnvEntry{
			Name:  "GOMIPS",
			Value: getGOMIPS(mArch),
		})
	case strings.HasPrefix(mArch, archconst.ARCH_PPC64):
		env = append(env, &dukkha.EnvEntry{
			Name:  "GOPPC64",
			Value: getGOPPC64(mArch),
		})
	}

	return append(env, cgoSpec.getEnv(
		v.HostKernel() != v.MatrixKernel() || v.HostArch() != mArch, /* doing cross compile */
		v.MatrixKernel(), /* target kernel */
		mArch,            /* target arch */
		v.HostOS(),       /* host os */
		v.MatrixLibc(),   /* target libc */
	)...)
}

func getGO386(mArch string) string {
	if !strings.HasPrefix(mArch, "x86") {
		return ""
	}

	if strings.HasSuffix(mArch, "sf") {
		return "softfloat"
	}

	return "sse2"
}

func getGOAMD64(mArch string) string {
	microArch := strings.TrimPrefix(mArch, "amd64")
	if len(microArch) == 0 {
		return "v1"
	}

	return microArch
}

func getGOARM(mArch string) string {
	level := strings.TrimPrefix(strings.TrimPrefix(mArch, "arm"), "v")
	if len(level) == 0 {
		return "7"
	}

	return level
}

func getGOARM64(mArch string) string {
	level := strings.TrimPrefix(strings.TrimPrefix(mArch, "arm64"), "v")
	if len(level) == 0 {
		return "8"
	}

	return level
}

func getGOMIPS(mArch string) string {
	if strings.HasSuffix(mArch, "sf") {
		return "softfloat"
	}

	return "hardfloat"
}

func getGOMIPS64(mArch string) string {
	if strings.HasSuffix(mArch, "sf") {
		return "softfloat"
	}

	return "hardfloat"
}

func getGOPPC64(mArch string) string {
	// ppc64{, le}{, v8, v9}

	isa := strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(mArch, "ppc64"), "le"), "v")
	if len(isa) == 0 {
		return "power8"
	}

	// power8 or power9
	return "power" + isa
}

type buildOptions struct {
	rs.BaseField `yaml:"-"`

	Race    bool     `yaml:"race"`
	LDFlags []string `yaml:"ldflags"`
	Tags    []string `yaml:"tags"`
}

func (opts buildOptions) generateArgs() []string {
	var args []string
	if opts.Race {
		args = append(args, "-race")
	}

	if len(opts.LDFlags) != 0 {
		args = append(args, "-ldflags", strings.Join(opts.LDFlags, " "))
	}

	if len(opts.Tags) != 0 {
		args = append(args, "-tags",
			// ref: https://golang.org/doc/go1.13#go-command
			// The go build flag -tags now takes a comma-separated list of build tags,
			// to allow for multiple tags in GOFLAGS. The space-separated form is
			// deprecated but still recognized and will be maintained.
			strings.Join(opts.Tags, ","),
		)
	}

	return args
}
