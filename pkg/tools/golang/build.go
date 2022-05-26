package golang

import (
	"strings"

	"arhat.dev/pkg/archconst"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
)

func createBuildEnv(v dukkha.EnvValues, buildSpec buildOptions, cgoSpec CGOSepc) (env dukkha.Env) {
	// set GOOS
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

	// set GOARCH
	mArch := v.MatrixArch()
	goarch, ok := constant.GetGolangArch(mArch)
	if !ok {
		goarch = string(mArch)
	}

	if len(goarch) != 0 {
		env = append(env, &dukkha.EnvEntry{
			Name:  "GOARCH",
			Value: goarch,
		})
	}

	// set GOARCH specific micro arch

	spec, ok := archconst.Parse[byte](mArch)
	switch spec.Name {
	case archconst.ARCH_AMD64:
		env = append(env, &dukkha.EnvEntry{
			Name:  "GOAMD64",
			Value: spec.MicroArch,
		})
	case archconst.ARCH_X86:
		var go386 string
		if spec.SoftFloat {
			go386 = "softfloat"
		} else {
			go386 = "sse2"
		}

		env = append(env, &dukkha.EnvEntry{
			Name:  "GO386",
			Value: go386,
		})
	case archconst.ARCH_ARM64:
		env = append(env, &dukkha.EnvEntry{
			Name:  "GOARM64",
			Value: strings.TrimPrefix(spec.MicroArch, "v"),
		})
	case archconst.ARCH_ARM:
		env = append(env, &dukkha.EnvEntry{
			Name:  "GOARM",
			Value: strings.TrimPrefix(spec.MicroArch, "v"),
		})
	case archconst.ARCH_MIPS64:
		var gomips64 string
		if spec.SoftFloat {
			gomips64 = "softfloat"
		} else {
			gomips64 = "hardfloat"
		}

		env = append(env, &dukkha.EnvEntry{
			Name:  "GOMIPS64",
			Value: gomips64,
		})
	case archconst.ARCH_MIPS:
		var gomips string
		if spec.SoftFloat {
			gomips = "softfloat"
		} else {
			gomips = "hardfloat"
		}

		env = append(env, &dukkha.EnvEntry{
			Name:  "GOMIPS",
			Value: gomips,
		})
	case archconst.ARCH_PPC64:
		env = append(env, &dukkha.EnvEntry{
			Name:  "GOPPC64",
			Value: "power" + strings.TrimPrefix(spec.MicroArch, "v"),
		})
	}

	// set GCCGO
	if len(buildSpec.GCCGo) != 0 {
		env = append(env, &dukkha.EnvEntry{
			Name:  "GCCGO",
			Value: buildSpec.GCCGo,
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

type buildOptions struct {
	rs.BaseField `yaml:"-"`

	// Race (-race)
	Race bool `yaml:"race"`

	// Compiler to use, gccgo or gc
	Compiler string `yaml:"compiler"`

	// ASMFlags (-asmflags) for `go tool asm`
	ASMFlags []string `yaml:"asm_flags"`

	// GoCompilerFlags (-gcflags) for `go tool compile`
	GoCompilerFlags []string `yaml:"go_compiler_flags"`

	// GCCGo command to run gccgo
	GCCGo string `yaml:"gccgo"`

	// GCCGoFlags (-gccgoflags) for CC (both compiler and linker)
	GCCGoFlags []string `yaml:"gccgo_flags"`

	// LDFlags (-ldflags)
	LDFlags []string `yaml:"ldflags"`

	Tags []string `yaml:"tags"`
}

func (opts buildOptions) generateArgs() []string {
	var args []string
	if opts.Race {
		args = append(args, "-race")
	}

	if len(opts.Compiler) != 0 {
		args = append(args, opts.Compiler)
	}

	if len(opts.ASMFlags) != 0 {
		args = append(args, "-asmflags", strings.Join(opts.ASMFlags, " "))
	}

	if len(opts.GoCompilerFlags) != 0 {
		args = append(args, "-gcflags", strings.Join(opts.GoCompilerFlags, " "))
	}

	if len(opts.GCCGoFlags) != 0 {
		args = append(args, "-gccgoflags", strings.Join(opts.GCCGoFlags, " "))
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
