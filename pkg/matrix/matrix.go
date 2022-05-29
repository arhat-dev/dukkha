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

	// catch other matrix fields
	Data map[string]*Vector `yaml:",inline,omitempty"`
}

func (s *Spec) HasUserValue() bool {
	if s == nil {
		return false
	}

	return len(s.Include) != 0 || len(s.Exclude) != 0 || len(s.Data) != 0
}

func (s *Spec) AsFilter() (ret Filter) {
	if !s.HasUserValue() {
		return
	}

	entries := s.GenerateEntries(Filter{})
	for _, ent := range entries {
		for k, v := range ent {
			ret.AddMatch(k, v)
		}
	}

	return
}

// GenerateEntries generates a set of matrix entries from the spec

func (s *Spec) GenerateEntries(filter Filter) (ret []Entry) {
	if !s.HasUserValue() {
		return
	}

	all := make(map[string][]string)

	for name, vec := range s.Data {
		all[name] = vec.Vec
	}

	// remove excluded
	var removeMatchList []map[string]string
	for _, ex := range s.Exclude {
		removeMatchList = append(
			removeMatchList,
			CartesianProduct(flattenVectorMap(ex.Data))...,
		)
	}

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
		mat := CartesianProduct(flattenVectorMap(inc.Data))
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
