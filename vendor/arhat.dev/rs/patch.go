package rs

import (
	"encoding/json"
	"fmt"
	"reflect"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/itchyny/gojq"
	"gopkg.in/yaml.v3"
)

type MergeSource struct {
	BaseField `yaml:"-" json:"-"`

	// Value for the source
	Value *AnyObject `yaml:"value,omitempty"`

	// Select some data from the source
	Select string `yaml:"select,omitempty"`
}

type renderingPatchSpec struct {
	BaseField `yaml:"-" json:"-"`

	// Value for the renderer
	//
	// 	say we have a yaml list (`[bar]`) stored at https://example.com/bar.yaml
	//
	// 		foo@http!:
	// 		  value: https://example.com/bar.yaml
	// 		  merge: { value: [foo] }
	//
	// then the resolve value of foo will be `[bar, foo]`
	Value *AnyObject `yaml:"value"`

	// Merge additional data into Value
	//
	// this action happens first
	Merge []MergeSource `yaml:"merge,omitempty"`

	// Patch Value using standard rfc6902 json-patch
	//
	// this action happens after merge
	Patch []JSONPatchSpec `yaml:"patch"`

	// Select part of the data as final result
	//
	// this action happens after merge and patch
	Select string `yaml:"select"`

	// Unique to make sure elements in the sequence is unique
	//
	// only effective when Value is yaml sequence
	Unique bool `yaml:"unique"`

	// MapListItemUnique to ensure items are unique in all merged lists respectively
	// lists with no merge data input are untouched
	MapListItemUnique bool `yaml:"map_list_item_unique"`

	// MapListAppend to append lists instead of replacing existing list
	MapListAppend bool `yaml:"map_list_append"`
}

func runJQ(query string, data interface{}) (interface{}, error) {
	q, err := gojq.Parse(query)
	if err != nil {
		return nil, fmt.Errorf("invalid jq query: %w", err)
	}

	var (
		ret interface{}
		idx int
	)

	iter := q.Run(data)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		if err, ok := v.(error); ok {
			return nil, fmt.Errorf("jq query failed: %w", err)
		}

		switch idx {
		case 0:
			ret = v
		case 1:
			ret = []interface{}{ret, v}
		default:
			ret = append(ret.([]interface{}), v)
		}

		idx++
	}

	return ret, nil
}

func (s *renderingPatchSpec) merge(resolvedValueData []byte) (interface{}, error) {
	type resolvedMergeSource struct {
		Value interface{} `yaml:"value,omitempty"`
	}

	var (
		valueData interface{}
		err       error
	)
	if len(resolvedValueData) != 0 {
		err = yaml.Unmarshal(resolvedValueData, &valueData)
		if err != nil {
			return nil, err
		}
	}

	var mergeSrc []*resolvedMergeSource
	for _, m := range s.Merge {
		mBytes, err := yaml.Marshal(m.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid non yaml merge value: %w", err)
		}

		var v interface{}
		err = yaml.Unmarshal(mBytes, &v)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal merge value as plain text: %w", err)
		}

		if len(m.Select) != 0 {
			v, err = runJQ(m.Select, v)
			if err != nil {
				return nil, err
			}
		}

		mergeSrc = append(mergeSrc, &resolvedMergeSource{
			Value: v,
		})
	}

doMerge:
	switch dt := valueData.(type) {
	case []interface{}:
		for _, merge := range mergeSrc {
			switch mt := merge.Value.(type) {
			case []interface{}:
				dt = append(dt, mt...)

				if s.Unique {
					dt = UniqueList(dt)
				}
			case nil:
				// no value to merge, skip
			default:
				// invalid type, not able to merge
				return nil, fmt.Errorf("unexpected non list value of merge, got %T", mt)
			}
		}

		return dt, nil
	case map[string]interface{}:
		var err error
		for _, merge := range mergeSrc {
			switch mt := merge.Value.(type) {
			case map[string]interface{}:
				dt, err = MergeMap(dt, mt, s.MapListAppend, s.MapListItemUnique)
				if err != nil {
					return nil, fmt.Errorf("failed to merge map value: %w", err)
				}
			case nil:
				// no value to merge, skip
			default:
				// invalid type, not able to merge
				return nil, fmt.Errorf("unexpected non map value of merge, got %T", mt)
			}
		}

		return dt, nil
	case nil:
		switch len(mergeSrc) {
		case 0:
			return nil, nil
		case 1:
			return mergeSrc[0].Value, nil
		default:
			valueData = mergeSrc[0].Value
			mergeSrc = mergeSrc[1:]
			goto doMerge
		}
	default:
		// scalar types, not supported
		return nil, fmt.Errorf("mergering scalar type value is not supported")
	}
}

// Apply Merge and Patch to Value, Unique is ensured if set to true
func (s *renderingPatchSpec) ApplyTo(resolvedValueData []byte) ([]byte, error) {
	data, err := s.merge(resolvedValueData)
	if err != nil {
		return nil, err
	}

	// apply select action to patches
	var patchSrc []*resolvedJSONPatchSpec
	for i, p := range s.Patch {
		var vBytes []byte
		vBytes, err = json.Marshal(p.Value)
		if err != nil {
			return nil, err
		}

		var v interface{}
		err = json.Unmarshal(vBytes, &v)
		if err != nil {
			return nil, err
		}

		patchSrc = append(patchSrc, &resolvedJSONPatchSpec{
			Path:      p.Path,
			Operation: p.Operation,
			Value:     v,
		})

		if len(p.Select) == 0 {
			continue
		}

		patchSrc[i].Value, err = runJQ(p.Select, v)
		if err != nil {
			return nil, err
		}
	}

	if len(patchSrc) == 0 {
		if len(s.Select) != 0 {
			data, err = runJQ(s.Select, data)
			if err != nil {
				return nil, err
			}
		}

		return yaml.Marshal(data)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	patchData, err := json.Marshal(patchSrc)
	if err != nil {
		return nil, err
	}

	patch, err := jsonpatch.DecodePatch(patchData)
	if err != nil {
		return nil, err
	}

	patchedDoc, err := patch.ApplyIndentWithOptions(jsonData, "", &jsonpatch.ApplyOptions{
		SupportNegativeIndices:   true,
		EnsurePathExistsOnAdd:    false,
		AccumulatedCopySizeLimit: 0,
		AllowMissingPathOnRemove: true,
	})
	if err != nil {
		return nil, err
	}

	if len(s.Select) == 0 {
		return patchedDoc, nil
	}

	var ret interface{}
	err = json.Unmarshal(patchedDoc, &ret)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal patched value: %w", err)
	}

	data, err = runJQ(s.Select, ret)
	if err != nil {
		return nil, err
	}

	return yaml.Marshal(data)
}

func MergeMap(
	original, additional map[string]interface{},

	// options
	appendList bool,
	uniqueInListItems bool,
) (map[string]interface{}, error) {
	out := make(map[string]interface{}, len(original))
	for k, v := range original {
		out[k] = v
	}

	var err error
	for k, v := range additional {
		switch newVal := v.(type) {
		case map[string]interface{}:
			if originalVal, ok := out[k]; ok {
				if orignalMap, ok := originalVal.(map[string]interface{}); ok {
					out[k], err = MergeMap(orignalMap, newVal, appendList, uniqueInListItems)
					if err != nil {
						return nil, err
					}

					continue
				} else {
					return nil, fmt.Errorf("unexpected non map data %q: %v", k, orignalMap)
				}
			} else {
				out[k] = newVal
			}
		case []interface{}:
			if originalVal, ok := out[k]; ok {
				if originalList, ok := originalVal.([]interface{}); ok {
					if appendList {
						originalList = append(originalList, newVal...)
					} else {
						originalList = newVal
					}

					if uniqueInListItems {
						originalList = UniqueList(originalList)
					}

					out[k] = originalList

					continue
				} else {
					return nil, fmt.Errorf("unexpected non list data %q: %v", k, originalList)
				}
			} else {
				out[k] = newVal
			}
		default:
			out[k] = newVal
		}
	}

	return out, nil
}

func UniqueList(dt []interface{}) []interface{} {
	var ret []interface{}
	dupAt := make(map[int]struct{})
	for i := range dt {
		_, isDup := dupAt[i]
		if isDup {
			continue
		}

		for j := i; j < len(dt); j++ {
			if reflect.DeepEqual(dt[i], dt[j]) {
				dupAt[j] = struct{}{}
			}
		}

		ret = append(ret, dt[i])
	}

	return ret
}

// JSONPatchSpec per rfc6902
type JSONPatchSpec struct {
	BaseField `yaml:"-" json:"-"`

	Operation string `yaml:"op" json:"op"`

	Path string `yaml:"path" json:"path"`

	Value *AnyObject `yaml:"value,omitempty" json:"value,omitempty"`

	// Select part of the value for patching
	//
	// this action happens before patching
	Select string `yaml:"select" json:"-"`
}

type resolvedJSONPatchSpec struct {
	Operation string      `yaml:"op" json:"op"`
	Path      string      `yaml:"path" json:"path"`
	Value     interface{} `yaml:"value,omitempty" json:"value,omitempty"`
}
