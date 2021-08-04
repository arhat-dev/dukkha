package field

import (
	"encoding/json"
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"

	jsonpatch "github.com/evanphx/json-patch/v5"
)

type MergeSource struct {
	BaseField `yaml:"-" json:"-"`

	Data interface{} `yaml:"data,omitempty"`
}

type PatchSpec struct {
	BaseField

	// Value for the renderer
	//
	// 	say we have a yaml list ([bar]) stored at https://example.com/dukkha.yaml
	//
	// 		foo@http!:
	// 		  value: https://example.com/dukkha.yaml
	// 		  merge: [foo]
	//
	// then the resolve value of foo will be [bar, foo]
	Value interface{} `yaml:"value"`

	// Merge additional data into Value
	Merge []MergeSource `yaml:"merge,omitempty"`

	// Patches Value using standard rfc6902 json-patch
	Patches []JSONPatchSpec `yaml:"patches"`

	// Unique to make sure elements in the sequence is unique
	//
	// only effective when Value is yaml sequence
	Unique bool `yaml:"unique"`

	// MapListItemUnique to ensure items are unique in all merged lists respectively
	// lists with no merge data input are untouched
	MapListItemUnique bool `yaml:"map_list_item_unique"`
}

func (s *PatchSpec) merge(yamlData []byte) (interface{}, error) {
	var data interface{}
	if len(yamlData) != 0 {
		err := yaml.Unmarshal(yamlData, &data)
		if err != nil {
			return nil, err
		}
	}

	mergeSrc := s.Merge

doMerge:
	switch dt := data.(type) {
	case []interface{}:
		for _, merge := range mergeSrc {
			switch mt := merge.Data.(type) {
			case []interface{}:
				dt = append(dt, mt...)

				if s.Unique {
					dt = uniqueList(dt)
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
			switch mt := merge.Data.(type) {
			case map[string]interface{}:
				dt, err = mergeMap(dt, mt, s.Unique)
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
		// TODO: do we really want to allow empty value?
		// 		 could it be some kind of error that should be checked during merging?
		switch len(mergeSrc) {
		case 0:
			return nil, nil
		case 1:
			return mergeSrc[0].Data, nil
		default:
			data = mergeSrc[0].Data
			mergeSrc = mergeSrc[1:]
			goto doMerge
		}
	default:
		// scalar types, not supported
		return nil, fmt.Errorf("mergering scalar type value is not supported")
	}
}

// Apply Merge and Patch to Value, Unique is ensured if set to true
func (s *PatchSpec) ApplyTo(yamlData []byte) ([]byte, error) {
	data, err := s.merge(yamlData)
	if err != nil {
		return nil, err
	}

	if len(s.Patches) == 0 {
		return yaml.Marshal(data)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	patchData, err := json.Marshal(s.Patches)
	if err != nil {
		return nil, err
	}

	patch, err := jsonpatch.DecodePatch(patchData)
	if err != nil {
		return nil, err
	}

	return patch.ApplyIndentWithOptions(jsonData, "", &jsonpatch.ApplyOptions{
		SupportNegativeIndices:   true,
		EnsurePathExistsOnAdd:    true,
		AccumulatedCopySizeLimit: 0,
		AllowMissingPathOnRemove: false,
	})
}

func mergeMap(original, additional map[string]interface{}, unique bool) (map[string]interface{}, error) {
	out := make(map[string]interface{}, len(original))
	for k, v := range original {
		out[k] = v
	}

	var err error
	for k, v := range additional {
		switch v := v.(type) {
		case map[string]interface{}:
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k], err = mergeMap(bv, v, unique)
					if err != nil {
						return nil, err
					}

					continue
				} else {
					return nil, fmt.Errorf("unexpected non map data %q: %v", k, bv)
				}
			} else {
				out[k] = v
			}
		case []interface{}:
			if bv, ok := out[k]; ok {
				if bv, ok := bv.([]interface{}); ok {
					if unique {
						out[k] = uniqueList(append(bv, v...))
					} else {
						out[k] = append(bv, v...)
					}

					continue
				} else {
					return nil, fmt.Errorf("unexpected non list data %q: %v", k, bv)
				}
			} else {
				out[k] = v
			}
		default:
			out[k] = v
		}
	}

	return out, nil
}

func uniqueList(dt []interface{}) []interface{} {
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

	Value interface{} `yaml:"value,omitempty" json:"value,omitempty"`
}
