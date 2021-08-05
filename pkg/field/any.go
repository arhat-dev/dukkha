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
	dataBytes, err := yaml.Marshal(n)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(dataBytes, &ad.data)
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
	switch n.ShortTag() {
	case "!!seq":
		o.arrayData = &arrayData{}

		return o.arrayData.UnmarshalYAML(n)
	case "!!map":
		md := Init(&mapData{}, nil).(*mapData)
		err := md.UnmarshalYAML(n)
		o.mapData = md

		return err
	default:
		dataBytes, err := yaml.Marshal(n)
		if err != nil {
			return err
		}

		return yaml.Unmarshal(dataBytes, &o.scalarData)
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
