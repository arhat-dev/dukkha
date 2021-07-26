package golang

import (
	"strings"

	"arhat.dev/dukkha/pkg/constant"
)

func createBuildEnv(mKernel, mArch string) []string {
	var env []string
	goos, _ := constant.GetGolangOS(mKernel)
	if len(goos) != 0 {
		env = append(env, "GOOS="+goos)
	}

	goarch, _ := constant.GetGolangArch(mArch)
	if len(goarch) != 0 {
		env = append(env, "GOARCH="+goarch)
	}

	if gomips := getGOMIPS(mArch); len(gomips) != 0 {
		env = append(env, "GOMIPS="+gomips, "GOMIPS64="+gomips)
	}

	if goarm := getGOARM(mArch); len(goarm) != 0 {
		env = append(env, "GOARM="+goarm)
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
