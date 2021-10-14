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

type AnyObjectMap struct {
	BaseField `yaml:"-" json:"-"`

	Data map[string]*AnyObject `rs:"other"`
}

// NormalizedValue returns value of the AnyObjectMap with type map[string]interface{}
func (aom *AnyObjectMap) NormalizedValue() interface{} {
	if aom.Data == nil {
		return map[string]interface{}(nil)
	}

	ret := make(map[string]interface{}, len(aom.Data))
	for k, v := range aom.Data {
		ret[k] = v.NormalizedValue()
	}

	return ret
}

func (aom *AnyObjectMap) MarshalYAML() (interface{}, error) { return aom.Data, nil }
func (aom *AnyObjectMap) MarshalJSON() ([]byte, error)      { return json.Marshal(aom.Data) }

// AnyObject is a `interface{}` equivalent with rendering suffix support
type AnyObject struct {
	mapData *AnyObjectMap

	// TODO: currently there is no way to support rendering suffix
	// 	     for array object
	sliceData []*AnyObject

	scalarData interface{}

	originalNode *yaml.Node
}

// NormalizedValue returns value with primitive types of the AnyObject
// that is:
// 	for maps: map[string]interface{}
// 	for slices: []interface{}
//  and primitive types for scalar types
//
// should be called after fields being resolved
func (o *AnyObject) NormalizedValue() interface{} {
	switch {
	case o == nil:
		return nil
	case o.mapData != nil:
		return o.mapData.NormalizedValue()
	case o.sliceData != nil:
		ret := make([]interface{}, len(o.sliceData))
		for i := range o.sliceData {
			ret[i] = o.sliceData[i].NormalizedValue()
		}
		return ret
	default:
		return o.scalarData
	}
}

func (o *AnyObject) MarshalYAML() (interface{}, error) { return o.NormalizedValue(), nil }
func (o *AnyObject) MarshalJSON() ([]byte, error)      { return json.Marshal(o.NormalizedValue()) }

func (o *AnyObject) UnmarshalYAML(n *yaml.Node) error {
	switch n.Kind {
	case yaml.SequenceNode:
		return n.Decode(&o.sliceData)
	case yaml.MappingNode:
		o.mapData = Init(&AnyObjectMap{}, nil).(*AnyObjectMap)
		return n.Decode(o.mapData)
	default:
		switch n.ShortTag() {
		case "!!str":
			o.scalarData = n.Value
		case "!!binary":
			o.scalarData = n.Value
		default:
			o.originalNode = n
			return n.Decode(&o.scalarData)
		}

		return nil
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
