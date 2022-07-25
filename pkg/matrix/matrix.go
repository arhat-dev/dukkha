package matrix

import (
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/pkg/matrixhelper"
	"arhat.dev/rs"
)

// SpecItem is a helper type to support rendering suffix
// for list of maps, used in Include/Exclude
type SpecItem struct {
	rs.BaseField `yaml:"-"`

	Data map[string]*Vector `yaml:",inline"`
}

type Spec struct {
	rs.BaseField `yaml:"-"`

	// Exclude matched entries
	Exclude []*SpecItem `yaml:"exclude,omitempty"`

	// Include
	//
	// NOTE: included entries will be excluded by Exclude entires as well
	Include []*SpecItem `yaml:"include,omitempty"`

	// Values to catch all matrix values
	Values map[string]*Vector `yaml:",inline,omitempty"`
}

// IsEmpty returns true when this is no value in s
func (s *Spec) IsEmpty() bool {
	if s == nil {
		return true
	}

	return len(s.Include) == 0 && len(s.Exclude) == 0 && len(s.Values) == 0
}

// GenerateEntries generates a set of matrix entries from the spec
func (s *Spec) GenerateEntries(filter Filter) (ret []Entry) {
	if s.IsEmpty() {
		return
	}

	all := make(map[string][]string)

	for name, vec := range s.Values {
		all[name] = vec.Vec
	}

	// remove excluded
	var removeMatchList []map[string]string
	for _, ex := range s.Exclude {
		removeMatchList = append(
			removeMatchList,
			matrixhelper.CartesianProduct(flattenVectorMap(ex.Data), sliceutils.SortByKernelCmdArchLibcOther)...,
		)
	}

	var (
		matchFilter  []map[string]string
		ignoreFilter = filter.ignore
	)
	if len(filter.match) != 0 {
		matchFilter = matrixhelper.CartesianProduct(flattenVectorMap(filter.match), sliceutils.SortByKernelCmdArchLibcOther)
	}

	mat := matrixhelper.CartesianProduct(all, sliceutils.SortByKernelCmdArchLibcOther)
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
			ret = append(ret, spec)
			continue
		}

		for _, f := range matchFilter {
			if spec.Match(f) {
				ret = append(ret, spec)
				continue loop
			}
		}
	}

	// add included
	for _, inc := range s.Include {
		mat := matrixhelper.CartesianProduct(flattenVectorMap(inc.Data), sliceutils.SortByKernelCmdArchLibcOther)
	addInclude:
		for i := range mat {
			includeEntry := Entry(mat[i])

			for _, spec := range ret {
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
				ret = append(ret, includeEntry)
				continue
			}

			for _, f := range matchFilter {
				if includeEntry.Match(f) {
					ret = append(ret, includeEntry)
					continue addInclude
				}
			}
		}
	}

	return
}

// AsFilter converts this spec to a Filter
//
// s.Values and s.Include will become match rules
//
// s.Exclude will become ignore rules
func (s *Spec) AsFilter() (ret Filter) {
	if s.IsEmpty() {
		return
	}

	entries := s.GenerateEntries(Filter{})
	for _, ent := range entries {
		for k, v := range ent {
			ret.AddMatch(k, v)
		}
	}

	for _, ex := range s.Exclude {
		if ex == nil {
			continue
		}

		for k, v := range ex.Data {
			for _, vv := range v.Vec {
				ret.AddIgnore(k, vv)
			}
		}
	}

	return
}
