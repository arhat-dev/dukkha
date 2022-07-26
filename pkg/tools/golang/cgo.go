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

	// CPPFlags (env CGO_CPPFLAGS) C preprocessor flags
	CPPFlags []string `yaml:"cppflags"`

	// CFlags (env CGO_CFLAGS) C flags
	CFlags []string `yaml:"cflags"`

	// CXXFlags (env CGO_CXXFLAGS)
	CXXFlags []string `yaml:"cxxflags"`

	// FFlags (env CGO_FFLAGS) Fortran flags
	FFlags []string `yaml:"fflags"`

	// LDFlags (env CGO_LDFLAGS) ldflags
	LDFlags []string `yaml:"ldflags"`

	// CC sets C compiler path or executable name
	CC string `yaml:"cc"`

	// CXX sets C++ compiler path or executable name
	CXX string `yaml:"cxx"`

	// FC sets Fortan compiler path or executable name
	FC string `yaml:"fc"`
}

func (c CGOSepc) getEnv(
	doingCrossCompiling bool,
	mKernel, mArch, hostOS, targetLibc string,
) dukkha.NameValueList {
	if !c.Enabled {
		return dukkha.NameValueList{
			{
				Name:  "CGO_ENABLED",
				Value: "0",
			},
		}
	}

	var ret dukkha.NameValueList
	ret = append(ret, &dukkha.NameValueEntry{
		Name:  "CGO_ENABLED",
		Value: "1",
	})

	appendListEnv := func(name string, values, defVals []string) {
		actual := values
		if len(values) == 0 {
			actual = defVals
		}

		if len(actual) != 0 {
			ret = append(ret, &dukkha.NameValueEntry{
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
			ret = append(ret, &dukkha.NameValueEntry{
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
