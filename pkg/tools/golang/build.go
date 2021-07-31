package golang

import (
	"strings"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
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
		env = append(env, dukkha.EnvEntry{
			Name:  "GOOS",
			Value: goos,
		})
	}

	goarch, _ := constant.GetGolangArch(v.MatrixArch())
	switch {
	case len(goarch) != 0:
	case len(v.MatrixArch()) != 0:
		goarch = v.MatrixArch()
	}

	if len(goarch) != 0 {
		env = append(env, dukkha.EnvEntry{
			Name:  "GOARCH",
			Value: goarch,
		})
	}

	if gomips := getGOMIPS(v.MatrixArch()); len(gomips) != 0 {
		env = append(env, dukkha.EnvEntry{
			Name:  "GOMIPS",
			Value: gomips,
		}, dukkha.EnvEntry{
			Name:  "GOMIPS64",
			Value: gomips,
		})
	}

	if goarm := getGOARM(v.MatrixArch()); len(goarm) != 0 {
		env = append(env, dukkha.EnvEntry{
			Name:  "GOARM",
			Value: goarm,
		})
	}

	return append(env, cgoSpec.getEnv(
		v.HostKernel() != v.MatrixKernel() || v.HostArch() != v.MatrixArch(),
		v.MatrixKernel(), v.MatrixArch(),
		v.HostOS(),
		v.MatrixLibc(),
	)...)
}

func getGOARM(mArch string) string {
	if strings.HasPrefix(mArch, "armv") {
		return strings.TrimPrefix(mArch, "armv")
	}

	return ""
}

func getGOMIPS(mArch string) string {
	if !strings.HasPrefix(mArch, "mips") {
		return ""
	}

	if strings.HasSuffix(mArch, "sf") {
		return "softfloat"
	}

	return "hardfloat"
}

type buildOptions struct {
	field.BaseField

	Race    bool     `yaml:"race"`
	LDFlags []string `yaml:"ldflags"`
	Tags    []string `yaml:"tags"`
}

func (opts buildOptions) generateArgs(useShell bool) []string {
	var args []string
	if opts.Race {
		args = append(args, "-race")
	}

	if len(opts.LDFlags) != 0 {
		args = append(args, "-ldflags",
			formatArgs(opts.LDFlags, useShell),
		)
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
