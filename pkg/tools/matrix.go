package tools

import (
	"os"
	"sort"
	"strings"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/field"
)

type MatrixConfig struct {
	field.BaseField

	Include []map[string][]string `yaml:"include"`
	Exclude []map[string][]string `yaml:"exclude"`

	OS   []string `yaml:"os"`
	Arch []string `yaml:"arch"`

	// catch other matrix fields
	Custom map[string][]string `dukkha:"other"`
}

func (mc *MatrixConfig) GetSpecs(filter map[string][]string) []MatrixSpec {
	if mc == nil {
		return []MatrixSpec{
			{
				"os":   os.Getenv(constant.ENV_HOST_OS),
				"arch": os.Getenv(constant.ENV_HOST_ARCH),
			},
		}
	}

	all := make(map[string][]string)

	osList := mc.OS
	if len(osList) == 0 {
		// add default host arch
		osList = []string{os.Getenv(constant.ENV_HOST_OS)}
	}
	all["os"] = osList

	archList := mc.Arch
	if len(archList) == 0 {
		archList = []string{os.Getenv(constant.ENV_HOST_ARCH)}
	}
	all["arch"] = archList

	for name := range mc.Custom {
		all[name] = mc.Custom[name]
	}

	// remove excluded
	var removeMatchList []map[string]string
	for _, ex := range mc.Exclude {
		removeMatchList = append(removeMatchList, CartesianProduct(ex)...)
	}

	var result []MatrixSpec

	var mf []map[string]string
	if len(filter) != 0 {
		mf = CartesianProduct(filter)
	}
	mat := CartesianProduct(all)
loop:
	for i := range mat {
		spec := MatrixSpec(mat[i])

		for _, toRemove := range removeMatchList {
			if spec.Match(toRemove) {
				continue loop
			}
		}

		if len(mf) == 0 {
			// no filter, add it
			result = append(result, spec)
			continue
		}

		for _, f := range mf {
			if spec.Match(f) {
				result = append(result, spec)
				continue loop
			}
		}
	}

	// add included
	for _, inc := range mc.Include {
		mat := CartesianProduct(inc)
	addInclude:
		for i := range mat {
			includeSpec := MatrixSpec(mat[i])

			for _, spec := range result {
				if spec.Equals(includeSpec) {
					continue addInclude
				}
			}

			if len(mf) == 0 {
				// no filter, add it
				result = append(result, includeSpec)
				continue
			}

			for _, f := range mf {
				if includeSpec.Match(f) {
					result = append(result, includeSpec)
					continue addInclude
				}
			}
		}
	}

	return result
}

type MatrixSpec map[string]string

func (m MatrixSpec) String() string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var pairs []string
	for _, k := range keys {
		pairs = append(pairs, k+": "+m[k])
	}

	return strings.Join(pairs, ", ")
}

func (m MatrixSpec) Match(a map[string]string) bool {
	if len(a) == 0 {
		return len(m) == 0
	}

	for k, v := range a {
		if m[k] != v {
			return false
		}
	}

	return true
}

func (m MatrixSpec) Equals(a map[string]string) bool {
	if a == nil {
		return m == nil
	}

	if len(a) != len(m) {
		return false
	}

	for k, v := range a {
		mv, ok := m[k]
		if !ok || mv != v {
			return false
		}
	}

	return true
}
