package matrix

import (
	"arhat.dev/dukkha/pkg/field"
)

type Spec struct {
	field.BaseField

	Include []map[string][]string `yaml:"include"`
	Exclude []map[string][]string `yaml:"exclude"`

	// TODO: validate kernel and arch values to ensure
	// 		 tools get expected value set
	Kernel []string `yaml:"kernel"`
	Arch   []string `yaml:"arch"`

	// catch other matrix fields
	Custom map[string][]string `dukkha:"other"`
}

func defaultSpecs(hostKernel, hostArch string) []Entry {
	return []Entry{
		{
			"kernel": hostKernel,
			"arch":   hostArch,
		},
	}
}

func (mc *Spec) GetSpecs(
	matchFilter map[string][]string,
	hostKernel, hostArch string,
) []Entry {
	if mc == nil {
		return defaultSpecs(hostKernel, hostArch)
	}

	hasUserValue := len(mc.Include) != 0 || len(mc.Exclude) != 0
	hasUserValue = hasUserValue || len(mc.Kernel) != 0 || len(mc.Arch) != 0 || len(mc.Custom) != 0

	if !hasUserValue {
		return defaultSpecs(hostKernel, hostArch)
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

	var result []Entry

	var mf []map[string]string
	if len(matchFilter) != 0 {
		mf = CartesianProduct(matchFilter)
	}
	mat := CartesianProduct(all)
loop:
	for i := range mat {
		spec := Entry(mat[i])

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
			includeEntry := Entry(mat[i])

			for _, spec := range result {
				if spec.Equals(includeEntry) {
					// already included
					continue addInclude
				}
			}

			if len(mf) == 0 {
				// no filter, add it
				result = append(result, includeEntry)
				continue
			}

			for _, f := range mf {
				if includeEntry.Match(f) {
					result = append(result, includeEntry)
					continue addInclude
				}
			}
		}
	}

	return result
}
