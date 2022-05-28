package matrix

import (
	"arhat.dev/rs"
	"gopkg.in/yaml.v3"
)

func flattenVectorMap(m map[string]*Vector) map[string][]string {
	ret := make(map[string][]string, len(m))
	for k, v := range m {
		ret[k] = v.Vec
	}

	return ret
}

func NewVector(elems ...string) *Vector {
	return rs.InitAny(&Vector{Vec: elems}, nil).(*Vector)
}

type Vector struct {
	rs.BaseField

	Vec []string `yaml:"__"`
}

func (v *Vector) Equals(a *Vector) bool {
	if v == nil {
		return a == nil
	}

	if a == nil {
		return false
	}

	if len(a.Vec) != len(v.Vec) {
		return false
	}

	for i, el := range v.Vec {
		if a.Vec[i] != el {
			return false
		}
	}

	return true
}

func (v *Vector) Empty() bool {
	if v == nil {
		return true
	}

	return len(v.Vec) == 0
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

func (v *Vector) MarshalYAML() (interface{}, error) {
	if v == nil {
		return nil, nil
	}

	return v.Vec, nil
}
