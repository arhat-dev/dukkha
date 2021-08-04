package field

import (
	"encoding/json"
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"

	jsonpatch "github.com/evanphx/json-patch/v5"
)

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
	Merge interface{} `yaml:"merge,omitempty"`

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

// Apply Merge and Patch to Value, Unique is ensured if set to true
func (s *PatchSpec) ApplyTo(yamlData []byte) ([]byte, error) {
	var data interface{}
	err := yaml.Unmarshal(yamlData, &data)
	if err != nil {
		return nil, err
	}

	switch dt := data.(type) {
	case []interface{}:
		switch mt := s.Merge.(type) {
		case []interface{}:
			dt = append(dt, mt...)

			if !s.Unique {
				data = dt
			} else {
				data = uniqueList(dt)
			}
		case nil:
			// no value to merge, skip if no patch
			if len(s.Patches) == 0 {
				return yamlData, nil
			}
		default:
			// invalid type, not able to merge
			return nil, fmt.Errorf("unexpected non list value of merge, got %T", mt)
		}
	case map[string]interface{}:
		switch mt := s.Merge.(type) {
		case map[string]interface{}:
			data, err = mergeMap(dt, mt)
			if err != nil {
				return nil, fmt.Errorf("failed to merge map value: %w", err)
			}
		case nil:
			// no value to merge, skip if no patch
			if len(s.Patches) == 0 {
				return yamlData, nil
			}
		default:
			// invalid type, not able to merge
			return nil, fmt.Errorf("unexpected non map value of merge, got %T", mt)
		}
	case nil:
		// TODO: do we really want to allow empty value?
		data = s.Merge
	default:
		// scalar types
		if s.Merge != nil {
			return nil, fmt.Errorf("patching scalar types are not supported")
		}
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

func mergeMap(original, additional map[string]interface{}) (map[string]interface{}, error) {
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
					out[k], err = mergeMap(bv, v)
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
					out[k] = append(bv, v...)
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
