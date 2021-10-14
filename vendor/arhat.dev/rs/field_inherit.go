package rs

import (
	"fmt"
	"reflect"
)

func (f *BaseField) Inherit(b *BaseField) error {
	if len(b.unresolvedFields) == 0 {
		return nil
	}

	if f.unresolvedFields == nil {
		f.unresolvedFields = make(map[unresolvedFieldKey]*unresolvedFieldValue)
	}

	for k, v := range b.unresolvedFields {
		existingV, ok := f.unresolvedFields[k]
		if !ok {
			f.unresolvedFields[k] = &unresolvedFieldValue{
				fieldName:         v.fieldName,
				fieldValue:        f._parentValue.FieldByName(v.fieldName),
				isCatchOtherField: v.isCatchOtherField,
				rawDataList:       v.rawDataList,
				renderers:         v.renderers,
			}

			continue
		}

		switch {
		case existingV.fieldName != v.fieldName,
			existingV.isCatchOtherField != v.isCatchOtherField:
			return fmt.Errorf(
				"rs: invalid field not match, want %q, got %q",
				existingV.fieldName, v.fieldName,
			)
		}

		existingV.rawDataList = append(existingV.rawDataList, v.rawDataList...)
	}

	if len(b.catchOtherCache) != 0 {
		if f.catchOtherCache == nil {
			f.catchOtherCache = make(map[string]reflect.Value)
		}

		for k, v := range b.catchOtherCache {
			f.catchOtherCache[k] = v
		}
	}

	if len(b.catchOtherFields) != 0 {
		if f.catchOtherFields == nil {
			f.catchOtherFields = make(map[string]struct{})
		}

		for k, v := range b.catchOtherFields {
			f.catchOtherFields[k] = v
		}
	}

	return nil
}
