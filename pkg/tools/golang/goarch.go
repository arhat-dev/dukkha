package golang

import (
	"strings"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
)

func createBuildEnv(mKernel, mArch string) dukkha.Env {
	var env dukkha.Env
	goos, _ := constant.GetGolangOS(mKernel)
	if len(goos) != 0 {
		env = append(env, dukkha.EnvEntry{
			Name:  "GOOS",
			Value: goos,
		})
	}

	goarch, _ := constant.GetGolangArch(mArch)
	if len(goarch) != 0 {
		env = append(env, dukkha.EnvEntry{
			Name:  "GOARCH",
			Value: goarch,
		})
	}

	if gomips := getGOMIPS(mArch); len(gomips) != 0 {
		env = append(env, dukkha.EnvEntry{
			Name:  "GOMIPS",
			Value: gomips,
		}, dukkha.EnvEntry{
			Name:  "GOMIPS64",
			Value: gomips,
		})
	}

	if goarm := getGOARM(mArch); len(goarm) != 0 {
		env = append(env, dukkha.EnvEntry{
			Name:  "GOARM",
			Value: goarm,
		})
	}

	return env
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
