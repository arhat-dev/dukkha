package rs

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func convertAnyObjectToAlterInterface(o *AnyObject) *alterInterface {
	switch {
	case o == nil:
		return nil
	case o.mapData != nil:
		if o.mapData.Data == nil {
			return &alterInterface{
				mapData:      (map[string]*alterInterface)(nil),
				originalNode: o.originalNode,
			}
		}

		md := make(map[string]*alterInterface, len(o.mapData.Data))
		for k, v := range o.mapData.Data {
			md[k] = convertAnyObjectToAlterInterface(v)
		}

		return &alterInterface{
			mapData:      md,
			originalNode: o.originalNode,
		}
	case o.sliceData != nil:
		sd := make([]*alterInterface, len(o.sliceData))
		for i := range o.sliceData {
			sd[i] = convertAnyObjectToAlterInterface(o.sliceData[i])
		}

		return &alterInterface{
			sliceData:    sd,
			originalNode: o.originalNode,
		}
	default:
		return &alterInterface{
			scalarData:   o.scalarData,
			originalNode: o.originalNode,
		}
	}
}

// alterInterface is a direct `interface{}` replacement for data
// unmarshaling with no rendering suffix support
//
// this struct exists mainly for preserving original string value of yaml
// field, since we need to keep track of all scalar fields, we have to handle
// map/slice data as well
type alterInterface struct {
	mapData   map[string]*alterInterface
	sliceData []*alterInterface

	scalarData interface{}

	originalNode *yaml.Node
}

func (f *alterInterface) HasValue() bool {
	return f.mapData != nil || f.sliceData != nil || f.scalarData != nil
}

func (f *alterInterface) Value() interface{} {
	switch {
	case f.mapData != nil:
		return f.mapData
	case f.sliceData != nil:
		return f.sliceData
	default:
		return f.scalarData
	}
}

func (f *alterInterface) NormalizedValue() interface{} {
	switch {
	case f == nil:
		return nil
	case f.mapData != nil:
		ret := make(map[string]interface{}, len(f.mapData))
		for k, v := range f.mapData {
			ret[k] = v.NormalizedValue()
		}

		return ret
	case f.sliceData != nil:
		ret := make([]interface{}, len(f.sliceData))
		for i := range f.sliceData {
			ret[i] = f.sliceData[i].NormalizedValue()
		}

		return ret
	default:
		return f.scalarData
	}
}

func (f *alterInterface) UnmarshalYAML(n *yaml.Node) error {
	switch n.Kind {
	case yaml.ScalarNode:
		switch n.ShortTag() {
		case "!!str":
			f.scalarData = n.Value
		case "!!binary":
			f.scalarData = n.Value
		default:
			f.originalNode = n
			return n.Decode(&f.scalarData)
		}

		return nil
	case yaml.MappingNode:
		f.mapData = make(map[string]*alterInterface)
		return n.Decode(&f.mapData)
	case yaml.SequenceNode:
		f.sliceData = make([]*alterInterface, 0)
		return n.Decode(&f.sliceData)
	default:
		return fmt.Errorf("unexpected node tag %q", n.ShortTag())
	}
}

func (f *alterInterface) MarshalYAML() (interface{}, error) {
	switch {
	case f.mapData != nil:
		return f.mapData, nil
	case f.sliceData != nil:
		return f.sliceData, nil
	case f.originalNode != nil:
		return f.originalNode, nil
	default:
		return f.scalarData, nil
	}
}
