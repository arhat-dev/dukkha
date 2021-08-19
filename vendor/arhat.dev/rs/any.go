package rs

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

var (
	_ Field          = (*AnyObject)(nil)
	_ yaml.Marshaler = (*AnyObject)(nil)
	_ json.Marshaler = (*AnyObject)(nil)
)

type mapData struct {
	BaseField `yaml:"-"`

	Data map[string]*AnyObject `rs:"other"`
}

func (md *mapData) MarshalYAML() (interface{}, error) { return md.Data, nil }
func (md *mapData) MarshalJSON() ([]byte, error)      { return json.Marshal(md.Data) }

// AnyObject is a `interface{}` equivalent with rendering suffix support
type AnyObject struct {
	mapData *mapData

	// TODO: currently there is no way to support rendering suffix
	// 	     for array object
	sliceData []*AnyObject

	scalarData interface{}
}

func (o *AnyObject) Value() interface{} {
	switch {
	case o.mapData != nil:
		return o.mapData
	case o.sliceData != nil:
		return o.sliceData
	default:
		return o.scalarData
	}
}

func (o *AnyObject) MarshalYAML() (interface{}, error) { return o.Value(), nil }
func (o *AnyObject) MarshalJSON() ([]byte, error)      { return json.Marshal(o.Value()) }

func (o *AnyObject) UnmarshalYAML(n *yaml.Node) error {
	switch n.Kind {
	case yaml.SequenceNode:
		return n.Decode(&o.sliceData)
	case yaml.MappingNode:
		o.mapData = Init(&mapData{}, nil).(*mapData)
		return n.Decode(o.mapData)
	default:
		return n.Decode(&o.scalarData)
	}
}

func (o *AnyObject) ResolveFields(rc RenderingHandler, depth int, fieldNames ...string) error {
	if o.mapData != nil {
		return o.mapData.ResolveFields(rc, depth, fieldNames...)
	}

	if o.sliceData != nil {
		for _, v := range o.sliceData {
			err := v.ResolveFields(rc, depth, fieldNames...)
			if err != nil {
				return err
			}
		}

		return nil
	}

	// scalar type data doesn't need resolving
	return nil
}
