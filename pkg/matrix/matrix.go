package matrix

import (
	"arhat.dev/rs"
)

// specItem is a helper type to support rendering suffix
// for list of maps, used in Include/Exclude
type specItem struct {
	rs.BaseField `yaml:"-"`

	Data map[string]*Vector `yaml:",inline"`
}

type Spec struct {
	rs.BaseField `yaml:"-"`

	Include []*specItem `yaml:"include,omitempty"`
	Exclude []*specItem `yaml:"exclude,omitempty"`

	// TODO: validate kernel and arch values to ensure
	// 		 tools get expected value set
	Kernel *Vector `yaml:"kernel,omitempty"`
	Arch   *Vector `yaml:"arch,omitempty"`

	// catch other matrix fields
	Custom map[string]*Vector `yaml:",inline,omitempty"`
}

func defaultSpecs(hostKernel, hostArch string) []Entry {
	return []Entry{
		{
			"kernel": hostKernel,
			"arch":   hostArch,
		},
	}
}

// nolint:gocyclo
func (mc *Spec) GenerateEntries(
	hostKernel, hostArch string,
	filter Filter,
) []Entry {
	if mc == nil {
		return defaultSpecs(hostKernel, hostArch)
	}

	hasUserValue := len(mc.Include) != 0 || len(mc.Exclude) != 0
	hasUserValue = hasUserValue || !mc.Kernel.Empty() || !mc.Arch.Empty() || len(mc.Custom) != 0

	if !hasUserValue {
		return defaultSpecs(hostKernel, hostArch)
	}

	all := make(map[string][]string)

	if !mc.Kernel.Empty() {
		all["kernel"] = mc.Kernel.Vector
	}

	if !mc.Arch.Empty() {
		all["arch"] = mc.Arch.Vector
	}

	for name := range mc.Custom {
		all[name] = mc.Custom[name].Vector
	}

	// remove excluded
	var removeMatchList []map[string]string
	for _, ex := range mc.Exclude {
		removeMatchList = append(
			removeMatchList,
			CartesianProduct(flattenVectorMap(ex.Data))...,
		)
	}

	var result []Entry

	var (
		matchFilter  []map[string]string
		ignoreFilter = filter.ignore
	)
	if len(filter.match) != 0 {
		matchFilter = CartesianProduct(flattenVectorMap(filter.match))
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

		for _, f := range ignoreFilter {
			if spec.MatchKV(f[0], f[1]) {
				continue loop
			}
		}

		if len(matchFilter) == 0 {
			// no filter, add it
			result = append(result, spec)
			continue
		}

		for _, f := range matchFilter {
			if spec.Match(f) {
				result = append(result, spec)
				continue loop
			}
		}
	}

	// add included
	for _, inc := range mc.Include {
		mat := CartesianProduct(flattenVectorMap(inc.Data))
	addInclude:
		for i := range mat {
			includeEntry := Entry(mat[i])

			for _, spec := range result {
				if spec.Equals(includeEntry) {
					// already included
					continue addInclude
				}
			}

			for _, f := range ignoreFilter {
				if includeEntry.MatchKV(f[0], f[1]) {
					continue addInclude
				}
			}

			if len(matchFilter) == 0 {
				result = append(result, includeEntry)
				continue
			}

			for _, f := range matchFilter {
				if includeEntry.Match(f) {
					result = append(result, includeEntry)
					continue addInclude
				}
			}
		}
	}

	return result
}
