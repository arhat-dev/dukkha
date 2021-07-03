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

	Kernel []string `yaml:"kernel"`
	Arch   []string `yaml:"arch"`

	// catch other matrix fields
	Custom map[string][]string `dukkha:"other"`
}

func (mc *MatrixConfig) GetSpecs(matchFilter map[string][]string) []MatrixSpec {
	if mc == nil {
		return []MatrixSpec{
			{
				"kernel": os.Getenv(constant.ENV_HOST_KERNEL),
				"arch":   os.Getenv(constant.ENV_HOST_ARCH),
			},
		}
	}

	hasUserValue := len(mc.Include) != 0 || len(mc.Exclude) != 0
	hasUserValue = hasUserValue || len(mc.Kernel) != 0 || len(mc.Arch) != 0 || len(mc.Custom) != 0

	if !hasUserValue {
		return []MatrixSpec{
			{
				"kernel": os.Getenv(constant.ENV_HOST_KERNEL),
				"arch":   os.Getenv(constant.ENV_HOST_ARCH),
			},
		}
	}

	all := make(map[string][]string)

	if len(mc.Kernel) != 0 {
		all["kernel"] = mc.Kernel
	}

	if len(mc.Arch) != 0 {
		all["arch"] = mc.Arch
	}

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
	if len(matchFilter) != 0 {
		mf = CartesianProduct(matchFilter)
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

// BriefString return all values concatenated with slash
func (m MatrixSpec) BriefString() string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		parts = append(parts, m[k])
	}

	return strings.Join(parts, "/")
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
