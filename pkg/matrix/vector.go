package matrix

import (
	"arhat.dev/rs"
	"gopkg.in/yaml.v3"
)

func flattenVectorMap(m map[string]*Vector) map[string][]string {
	ret := make(map[string][]string, len(m))
	for k, v := range m {
		ret[k] = v.Vector
	}

	return ret
}

func newVector(elems ...string) *Vector {
	return rs.Init(&Vector{Vector: elems}, nil).(*Vector)
}

type Vector struct {
	rs.BaseField

	Vector []string `yaml:"__"`
}

func (v *Vector) IsEmpty() bool {
	if v == nil {
		return true
	}

	return len(v.Vector) == 0
}

func (v *Vector) UnmarshalYAML(value *yaml.Node) error {
	// fake a map for vector
	return v.BaseField.UnmarshalYAML(&yaml.Node{
		Kind:  yaml.MappingNode,
		Value: "",
		Content: []*yaml.Node{
			{
				Kind:  yaml.ScalarNode,
				Value: "__@",
				Tag:   "!!str",
			},
			value,
		},
	})
}

func (v *Vector) ResolveFields(rc rs.RenderingHandler, depth int, names ...string) error {
	_ = names
	return v.BaseField.ResolveFields(rc, depth, "__")
}
