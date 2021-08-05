package field

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

	Data map[string]*AnyObject `dukkha:"other"`
}

func (md *mapData) MarshalYAML() (interface{}, error) { return md.Data, nil }
func (md *mapData) MarshalJSON() ([]byte, error)      { return json.Marshal(md.Data) }

type arrayData struct {
	data []*AnyObject
}

func (ad *arrayData) MarshalYAML() (interface{}, error) { return ad.data, nil }
func (ad *arrayData) MarshalJSON() ([]byte, error)      { return json.Marshal(ad.data) }

func (ad *arrayData) UnmarshalYAML(n *yaml.Node) error {
	return n.Decode(&ad.data)
}

// AnyObject is a `interface{}` equivalent with rendering suffix support
type AnyObject struct {
	mapData   *mapData
	arrayData *arrayData

	scalarData interface{}
}

func (o *AnyObject) MarshalYAML() (interface{}, error) {
	switch {
	case o.mapData != nil:
		return o.mapData, nil
	case o.arrayData != nil:
		return o.arrayData, nil
	default:
		return o.scalarData, nil
	}
}

func (o *AnyObject) MarshalJSON() ([]byte, error) {
	switch {
	case o.mapData != nil:
		return json.Marshal(o.mapData)
	case o.arrayData != nil:
		return json.Marshal(o.arrayData)
	default:
		return json.Marshal(o.scalarData)
	}
}

func (o *AnyObject) UnmarshalYAML(n *yaml.Node) error {
	switch n.Kind {
	case yaml.SequenceNode:
		o.arrayData = &arrayData{}
		return n.Decode(o.arrayData)
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

	if o.arrayData != nil {
		for _, v := range o.arrayData.data {
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
