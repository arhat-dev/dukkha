//go:build !rs_noyamlmarshaler
// +build !rs_noyamlmarshaler

package rs

import (
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"
)

var _ yaml.Marshaler = (*BaseField)(nil)

func (f *BaseField) MarshalYAML() (interface{}, error) {
	if !f.initialized() {
		return nil, fmt.Errorf("rs: struct not intialized before marshaling")
	}

	ret := make(map[string]interface{})

	// handle catch other fields first, so they can be overridden
	// by normal fields
	if f.catchOtherField != nil && f.catchOtherField.fieldValue.IsValid() {
		// NOTE: we MUST not use catch other cache since user can
		// 		 access map directly without updating our cache
		iter := f.catchOtherField.fieldValue.MapRange()
		for iter.Next() {
			// always include it no matter what value it has
			//
			// if not, unmarshal using our output will give
			// different result
			k, v := iter.Key().String(), iter.Value()

			if v.Kind() == reflect.Ptr && v.IsNil() {
				ret[k] = nil
				continue
			}

			ret[iter.Key().String()] = v.Interface()
		}
	}

	for k, v := range f.normalFields {
		if v.omitempty {
			if !v.fieldValue.IsValid() || v.fieldValue.IsZero() {
				// value not set or zero value, just ignore it
				continue
			}

			// value already set and not zero value
			switch v.fieldValue.Kind() {
			case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
				if v.fieldValue.Len() == 0 {
					continue
				}
			}
		}

		if v.fieldValue.Kind() == reflect.Ptr && v.fieldValue.IsNil() {
			ret[k] = nil
			continue
		}

		ret[k] = v.fieldValue.Interface()
	}

	return ret, nil
}
