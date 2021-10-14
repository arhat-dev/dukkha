package rs

import (
	"reflect"
	"strings"
)

type BaseField struct {
	_initialized uint32

	// _parentValue represents the parent struct value after being initialized
	_parentValue reflect.Value

	unresolvedFields map[unresolvedFieldKey]*unresolvedFieldValue

	ifaceTypeHandler InterfaceTypeHandler

	// yamlKey -> map value
	// TODO: separate static data and runtime generated data
	catchOtherCache map[string]reflect.Value

	catchOtherFields map[string]struct{}
}

type (
	// TODO: make field name as key
	unresolvedFieldKey struct {
		// NOTE: put `suffix` and `yamlKey` in key is to support fields with
		// 		 `rs:"other"` field tag, each item should be able
		// 		 to have its own renderer
		yamlKey string
		suffix  string
	}

	unresolvedFieldValue struct {
		fieldName   string
		fieldValue  reflect.Value
		rawDataList []*alterInterface
		renderers   []*suffixSpec

		isCatchOtherField bool
	}
)

func (f *BaseField) addUnresolvedField(
	// key part
	yamlKey string,
	suffix string,

	// value part
	fieldName string,
	fieldValue reflect.Value,
	isCatchOtherField bool,
	rawData *alterInterface,
) {
	if f.unresolvedFields == nil {
		f.unresolvedFields = make(map[unresolvedFieldKey]*unresolvedFieldValue)
	}

	key := unresolvedFieldKey{
		// yamlKey@suffix: ...
		yamlKey: yamlKey,
		suffix:  suffix,
	}

	if isCatchOtherField {
		if f.catchOtherFields == nil {
			f.catchOtherFields = make(map[string]struct{})
		}

		f.catchOtherFields[yamlKey] = struct{}{}
	}

	if old, exists := f.unresolvedFields[key]; exists {
		old.rawDataList = append(old.rawDataList, rawData)
		return
	}

	f.unresolvedFields[key] = &unresolvedFieldValue{
		fieldName:   fieldName,
		fieldValue:  fieldValue,
		rawDataList: []*alterInterface{rawData},
		renderers:   parseRenderingSuffix(suffix),

		isCatchOtherField: isCatchOtherField,
	}
}

type suffixSpec struct {
	name string

	patchSpec bool
	typeHint  TypeHint
}

func parseRenderingSuffix(rs string) []*suffixSpec {
	var ret []*suffixSpec
	for _, part := range strings.Split(rs, "|") {
		size := len(part)
		if size == 0 {
			continue
		}

		spec := &suffixSpec{
			patchSpec: part[size-1] == '!',
		}

		if spec.patchSpec {
			part = part[:size-1]
			// size-- // size not used any more
		}

		if idx := strings.LastIndexByte(part, '?'); idx >= 0 {
			// TODO: do we really want to panic when type hint is not valid?
			var err error
			spec.typeHint, err = ParseTypeHint(part[idx+1:])
			if err != nil {
				panic(err)
			}

			part = part[:idx]
		}

		spec.name = part

		ret = append(ret, spec)
	}

	return ret
}
