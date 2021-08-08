package field

import (
	"fmt"
	"reflect"
	"strings"
	"sync/atomic"

	"gopkg.in/yaml.v3"
)

func (f *BaseField) HasUnresolvedField() bool {
	return len(f.unresolvedFields) != 0
}

func (f *BaseField) ResolveFields(rc RenderingHandler, depth int, fieldNames ...string) error {
	if atomic.LoadUint32(&f._initialized) == 0 {
		return fmt.Errorf("field resolve: struct not intialized with Init()")
	}

	if depth == 0 {
		return nil
	}

	parentStruct := f._parentValue.Type().Elem()
	structName := parentStruct.String()

	if len(fieldNames) == 0 {
		// resolve all

		for i := 1; i < f._parentValue.Elem().NumField(); i++ {
			sf := parentStruct.Field(i)
			if !IsExported(sf.Name) {
				continue
			}

			err := f.resolveSingleField(
				rc, depth, structName, sf.Name, f._parentValue.Elem().Field(i),
			)
			if err != nil {
				return err
			}
		}

		return nil
	}

	for _, name := range fieldNames {
		fv := f._parentValue.Elem().FieldByName(name)
		if !fv.IsValid() {
			return fmt.Errorf("no such field %q in struct %q", name, parentStruct.String())
		}

		err := f.resolveSingleField(rc, depth, structName, name, fv)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *BaseField) resolveSingleField(
	rc RenderingHandler,
	depth int,

	structName string, // to make error message helpful
	fieldName string, // to make error message helpful

	targetField reflect.Value,
) error {
	handled := false
	for k, v := range f.unresolvedFields {
		if v.fieldName == fieldName {
			err := f.handleUnResolvedField(
				rc, depth, structName, fieldName, k, v, handled,
			)
			if err != nil {
				return err
			}

			handled = true
		}
	}

	return f.handleResolvedField(rc, depth, targetField)
}

// nolint:gocyclo
func (f *BaseField) handleResolvedField(
	rc RenderingHandler,
	depth int,
	targetField reflect.Value,
) error {
	if depth == 0 {
		return nil
	}

	switch targetField.Kind() {
	case reflect.Map:
		if targetField.IsNil() {
			return nil
		}

		iter := targetField.MapRange()
		for iter.Next() {
			err := f.handleResolvedField(rc, depth-1, iter.Value())
			if err != nil {
				return err
			}
		}
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		if targetField.IsNil() {
			// this is a resolved field, slice/array empty means no value
			return nil
		}

		for i := 0; i < targetField.Len(); i++ {
			tt := targetField.Index(i)
			err := f.handleResolvedField(rc, depth-1, tt)
			if err != nil {
				return err
			}
		}
	case reflect.Struct:
		// handled after switch
	case reflect.Ptr:
		fallthrough
	case reflect.Interface:
		if !targetField.IsValid() || targetField.IsZero() || targetField.IsNil() {
			return nil
		}
	default:
		// scalar types, no action required
		return nil
	}

	return tryResolve(rc, targetField, depth)
}

func tryResolve(rc RenderingHandler, targetField reflect.Value, depth int) error {
	if targetField.CanInterface() {
		fVal, canCallResolve := targetField.Interface().(Field)
		if canCallResolve {
			return fVal.ResolveFields(rc, depth)
		}
	}

	if targetField.CanAddr() && targetField.Addr().CanInterface() {
		fVal, canCallResolve := targetField.Addr().Interface().(Field)
		if canCallResolve {
			return fVal.ResolveFields(rc, depth)
		}
	}

	// no more field to resolve
	return nil
}

func (f *BaseField) handleUnResolvedField(
	rc RenderingHandler,
	depth int,

	structName string, // to make error message helpful
	fieldName string, // to make error message helpful

	key unresolvedFieldKey,
	v *unresolvedFieldValue,
	keepOld bool,
) error {
	target := v.fieldValue

	for i, rawData := range v.rawDataList {
		toResolve := rawData
		if v.isCatchOtherField {
			toResolve = rawData.mapData[key.yamlKey]
		}

		var (
			resolvedValue []byte
			err           error
		)

		for _, renderer := range v.renderers {
			// a patch is implied when renderer has a `!` suffix

			var patchSpec *PatchSpec
			if strings.HasSuffix(renderer, "!") {
				renderer = renderer[:len(renderer)-1]

				var patchSpecBytes []byte
				switch t := toResolve.Value().(type) {
				case string:
					patchSpecBytes = []byte(t)
				case []byte:
					patchSpecBytes = t
				default:
					patchSpecBytes, err = yaml.Marshal(toResolve.Value())
					if err != nil {
						return fmt.Errorf(
							"field: failed to marshal renderer data to patch spec bytes: %w",
							err,
						)
					}
				}

				patchSpec = Init(&PatchSpec{}, f.ifaceTypeHandler).(*PatchSpec)
				err = yaml.Unmarshal(patchSpecBytes, patchSpec)
				if err != nil {
					return fmt.Errorf(
						"field: failed to unmarshal patch spec\n\n%s\n\nfor renderer %q of %s.%s: %w",
						string(patchSpecBytes), renderer, structName, fieldName, err,
					)
				}

				err = patchSpec.ResolveFields(rc, -1)
				if err != nil {
					return fmt.Errorf(
						"field: failed to resolve patch spec\n\n%s\n\nfor renderer %q of %s.%s: %w",
						string(patchSpecBytes), renderer, structName, fieldName, err,
					)
				}

				toResolve = patchSpec.Value
			}

			resolvedValue, err = rc.RenderYaml(renderer, toResolve.Value())
			if err != nil {
				return fmt.Errorf(
					"field: failed to render value of %s.%s: %w",
					structName, fieldName, err,
				)
			}

			if patchSpec != nil {
				resolvedValue, err = patchSpec.ApplyTo(resolvedValue)
				if err != nil {
					return fmt.Errorf(
						"field: failed to apply patches to %s.%s: %w",
						structName, fieldName, err,
					)
				}
			}

			toResolve = &alterInterface{
				scalarData: resolvedValue,
			}
		}

		tmp := &alterInterface{}
		switch {
		case target.Kind() == reflect.String:
			tmp.scalarData = string(resolvedValue)
		case target.Kind() == reflect.Slice && target.Type().Elem().Kind() == reflect.Uint8:
			tmp.scalarData = resolvedValue
		default:
			err = yaml.Unmarshal(resolvedValue, tmp)
		}

		if err != nil {
			switch {
			case target.Type() == rawInterfaceType:
				tmp.scalarData = string(resolvedValue)
			default:
				return fmt.Errorf(
					"field: failed to unmarshal resolved value %q to interface: %w",
					resolvedValue, err,
				)
			}
		} else {
			// sometimes go-yaml will parse the input as string when it is not yaml
			//
			// revert that change by checking and resetting scalarData to resolvedValue
			switch tmp.scalarData.(type) {
			case string:
				tmp.scalarData = string(resolvedValue)
			case []byte:
				tmp.scalarData = resolvedValue
			case nil, bool, uintptr,
				complex64, complex128,
				float32, float64,
				int, int8, int16, int32, int64,
				uint, uint8, uint16, uint32, uint64:
				// unmarshaled scalar types, do nothing
			case interface{}:
				// TODO: narrow down the interface{} match
				// 		 this case will match all other types

				// map/struct and array/slice types are handled in arrayData and mapData
				// so we don't have to worry about these cases here
				tmp.scalarData = string(resolvedValue)
			}
		}

		if v.isCatchOtherField {
			tmp = &alterInterface{
				mapData: map[string]*alterInterface{
					key.yamlKey: tmp,
				},
			}
		}

		// TODO: currently we always keepOld when the field has tag
		// 		 `dukkha:"other"`, need to ensure this behavior won't
		// 	     leave inconsistent data

		actualKeepOld := keepOld || v.isCatchOtherField || i != 0
		err = f.unmarshal(key.yamlKey, tmp, target, actualKeepOld)
		if err != nil {
			return fmt.Errorf("field: failed to unmarshal resolved value to %T: %w", tmp.Value(), err)
		}
	}

	return tryResolve(rc, target, depth-1)
}

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
		renderers:   strings.Split(suffix, "|"),

		isCatchOtherField: isCatchOtherField,
	}
}

func (f *BaseField) isCatchOtherField(yamlKey string) bool {
	if f.catchOtherFields == nil {
		return false
	}

	_, ok := f.catchOtherFields[yamlKey]
	return ok
}
