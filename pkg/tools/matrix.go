package tools

import (
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

func (mc *MatrixConfig) GetSpecs() []MatrixSpec {
	if mc == nil {
		return []MatrixSpec{
			{
				"os":   GetHostOS(),
				"arch": GetHostArch(),
			},
		}
	}

	all := make(map[string][]string)

	osList := mc.OS
	if len(osList) == 0 {
		// add default host arch
		osList = []string{GetHostOS()}
	}
	all["os"] = osList

	archList := mc.Arch
	if len(archList) == 0 {
		archList = []string{GetHostArch()}
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

	mat := CartesianProduct(all)
loop:
	for i := range mat {
		spec := MatrixSpec(mat[i])

		for _, toRemove := range removeMatchList {
			if spec.Match(toRemove) {
				continue loop
			}
		}

		result = append(result, spec)
	}

	// add included
	for _, inc := range mc.Include {
		mat := CartesianProduct(inc)
	addInclude:
		for i := range mat {
			for _, spec := range result {
				if spec.Equals(mat[i]) {
					continue addInclude
				}
			}

			result = append(result, mat[i])
		}
	}

	return result
}

type MatrixSpec map[string]string

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
