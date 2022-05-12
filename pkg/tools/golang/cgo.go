package golang

import (
	"strings"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
)

type CGOSepc struct {
	rs.BaseField `yaml:"-"`

	// Enable cgo
	Enabled bool `yaml:"enabled"`

	// CGO_CPPFLAGS
	CPPFlags []string `yaml:"cppflags"`

	// CGO_CFLAGS
	CFlags []string `yaml:"cflags"`

	// CGO_CXXFLAGS
	CXXFlags []string `yaml:"cxxflags"`

	// CGO_FFLAGS
	FFlags []string `yaml:"fflags"`

	// CGO_LDFLAGS
	LDFlags []string `yaml:"ldflags"`

	CC string `yaml:"cc"`

	CXX string `yaml:"cxx"`

	FC string `yaml:"fc"`
}

func (c CGOSepc) getEnv(
	doingCrossCompiling bool,
	mKernel, mArch, hostOS, targetLibc string,
) dukkha.Env {
	if !c.Enabled {
		return dukkha.Env{
			{
				Name:  "CGO_ENABLED",
				Value: "0",
			},
		}
	}

	var ret dukkha.Env
	ret = append(ret, &dukkha.EnvEntry{
		Name:  "CGO_ENABLED",
		Value: "1",
	})

	appendListEnv := func(name string, values, defVals []string) {
		actual := values
		if len(values) == 0 {
			actual = defVals
		}

		if len(actual) != 0 {
			ret = append(ret, &dukkha.EnvEntry{
				Name:  name,
				Value: strings.Join(actual, " "),
			})
		}
	}

	appendEnv := func(name, value, defVal string) {
		actual := value
		if len(value) == 0 {
			actual = defVal
		}

		if len(actual) != 0 {
			ret = append(ret, &dukkha.EnvEntry{
				Name:  name,
				Value: actual,
			})
		}
	}

	var (
		cppflags []string
		cflags   []string
		cxxflags []string
		fflags   []string
		ldflags  []string

		cc, cxx string
	)

	if doingCrossCompiling {
		switch hostOS {
		case constant.Platform_Debian, constant.Platform_Ubuntu:
			var tripleName string
			switch mKernel {
			case constant.KERNEL_Linux:
				tripleName, _ = constant.GetDebianTripleName(mArch, mKernel, targetLibc)
			case constant.KERNEL_Darwin:
				// TODO: set darwin version
				tripleName, _ = constant.GetAppleTripleName(mArch, "")
			case constant.KERNEL_Windows:
				tripleName, _ = constant.GetDebianTripleName(mArch, mKernel, targetLibc)
			default:
				// not set
			}

			if len(tripleName) != 0 {
				cc = tripleName + "-gcc"
				cxx = tripleName + "-g++"
			}
		case constant.Platform_Alpine:
			tripleName, _ := constant.GetAlpineTripleName(mArch)

			if len(tripleName) != 0 {
				cc = tripleName + "-gcc"
				cxx = tripleName + "-g++"
			}
		case constant.Platform_MacOS:
			switch mKernel {
			case constant.KERNEL_Darwin:
				// cc = "clang"
				// cxx = "clang++"
			case constant.KERNEL_Linux:
				// TODO
			case constant.KERNEL_iOS:
				// TODO
			}
		}
	}

	// TODO: generate suitable flags
	appendListEnv("CGO_CPPFLAGS", c.CPPFlags, cppflags)
	appendListEnv("CGO_CFLAGS", c.CFlags, cflags)
	appendListEnv("CGO_CXXFLAGS", c.CXXFlags, cxxflags)
	appendListEnv("CGO_FFLAGS", c.FFlags, fflags)
	appendListEnv("CGO_LDFLAGS", c.LDFlags, ldflags)

	appendEnv("CC", c.CC, cc)
	appendEnv("CXX", c.CXX, cxx)
	appendEnv("FC", c.FC, "")

	return ret
}
