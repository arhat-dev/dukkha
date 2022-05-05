package golang

import (
	"strings"

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

	goarch, _ := constant.GetGolangArch(v.MatrixArch())
	switch {
	case len(goarch) != 0:
	case len(v.MatrixArch()) != 0:
		goarch = v.MatrixArch()
	}

	if len(goarch) != 0 {
		env = append(env, &dukkha.EnvEntry{
			Name:  "GOARCH",
			Value: goarch,
		})
	}

	if gomips := getGOMIPS(v.MatrixArch()); len(gomips) != 0 {
		env = append(env, &dukkha.EnvEntry{
			Name:  "GOMIPS",
			Value: gomips,
		}, &dukkha.EnvEntry{
			Name:  "GOMIPS64",
			Value: gomips,
		})
	} else if goarm := getGOARM(v.MatrixArch()); len(goarm) != 0 {
		env = append(env, &dukkha.EnvEntry{
			Name:  "GOARM",
			Value: goarm,
		})
	} else if goamd64 := getGOAMD64(v.MatrixArch()); len(goamd64) != 0 {
		env = append(env, &dukkha.EnvEntry{
			Name:  "GOAMD64",
			Value: goamd64,
		})
	} else if goppc64 := getGOPPC64(v.MatrixArch()); len(goppc64) != 0 {
		env = append(env, &dukkha.EnvEntry{
			Name:  "GOPPC64",
			Value: goppc64,
		})
	}

	return append(env, cgoSpec.getEnv(
		v.HostKernel() != v.MatrixKernel() || v.HostArch() != v.MatrixArch(),
		v.MatrixKernel(), v.MatrixArch(),
		v.HostOS(),
		v.MatrixLibc(),
	)...)
}

func getGOAMD64(mArch string) string {
	if strings.HasPrefix(mArch, "amd64v") {
		return strings.TrimPrefix(mArch, "amd64")
	}

	return ""
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

func getGOPPC64(mArch string) string {
	// ppc64{, le}{, v8, v9}

	if len(mArch) == 0 ||
		!strings.HasPrefix(mArch, "ppc64") ||
		!strings.HasSuffix(mArch[:len(mArch)-1], "v") {
		return ""
	}

	// power8 or power9
	return "power" + mArch[len(mArch)-1:]
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
