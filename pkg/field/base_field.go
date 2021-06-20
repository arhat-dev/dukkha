package field

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync/atomic"

	"arhat.dev/pkg/log"
	"gopkg.in/yaml.v3"
)

type (
	unresolvedFieldKey struct {
		fieldName string
		renderer  string
	}

	unresolvedFieldValue struct {
		fieldValue    reflect.Value
		yamlFieldName string
		rawData       []string
	}
)

type BaseField struct {
	_initialized uint32

	// _parentValue is always a pointer type with .Elem() to the struct
	// when initialized
	_parentValue reflect.Value

	unresolvedFields map[unresolvedFieldKey]*unresolvedFieldValue
}

func (f *BaseField) addUnresolvedField(
	fieldName string,
	fieldValue reflect.Value,
	yamlFieldName string,
	renderer, rawData string,
) {
	if f.unresolvedFields == nil {
		f.unresolvedFields = make(map[unresolvedFieldKey]*unresolvedFieldValue)
	}

	key := unresolvedFieldKey{
		fieldName: fieldName,
		renderer:  renderer,
	}

	if old, exists := f.unresolvedFields[key]; exists {
		old.rawData = append(old.rawData, rawData)
		return
	}

	f.unresolvedFields[key] = &unresolvedFieldValue{
		fieldValue:    fieldValue,
		yamlFieldName: yamlFieldName,
		rawData:       []string{rawData},
	}
}

func (f *BaseField) Resolve(ctx context.Context, render RenderingFunc, depth int) error {
	logger := log.Log.WithName("BaseField").WithFields(log.String("func", "Render"))

	var toRemove []unresolvedFieldKey
	for k, v := range f.unresolvedFields {
		logger := logger.WithFields(log.String("field", k.fieldName))

		logger.V("rendering")

		out := v.fieldValue.Interface()

		for _, rawData := range v.rawData {
			resolvedValue, err := render(ctx, k.renderer, rawData)
			if err != nil {
				return fmt.Errorf("field: failed to render value of this base field: %w", err)
			}

			err = yaml.Unmarshal([]byte(resolvedValue), out)
			if err != nil {
				return fmt.Errorf("field: failed to unmarshal resolved value: %w", err)
			}

			toRemove = append(toRemove, k)
		}
	}

	for _, k := range toRemove {
		delete(f.unresolvedFields, k)
	}

	return nil
}

// UnmarshalYAML handles renderer suffix
// nolint:gocyclo
func (f *BaseField) UnmarshalYAML(n *yaml.Node) error {
	if atomic.LoadUint32(&f._initialized) == 0 {
		return fmt.Errorf("field: struct not intialized with Init()")
	}

	type fieldKey struct {
		yamlKey string
	}

	type fieldSpec struct {
		fieldName  string
		fieldValue reflect.Value
	}

	fields := make(map[fieldKey]*fieldSpec)
	pt := f._parentValue.Type().Elem()

	addField := func(yamlKey, fieldName string, fieldValue reflect.Value) bool {
		key := fieldKey{yamlKey: yamlKey}
		if _, exists := fields[key]; exists {
			return false
		}

		fields[fieldKey{yamlKey: yamlKey}] = &fieldSpec{
			fieldName:  fieldName,
			fieldValue: fieldValue,
		}
		return true
	}

	ignoreField := func(yamlKey string) {
		key := fieldKey{yamlKey: yamlKey}
		delete(fields, key)
	}

	getField := func(yamlKey string) *fieldSpec {
		return fields[fieldKey{
			yamlKey: yamlKey,
		}]
	}

	logger := log.Log.WithName("field.BaseField").WithFields(
		log.String("func", "UnmarshalYAML"),
		log.String("struct", pt.String()),
	)

	var catchOtherField *fieldSpec
	// get expected fields first, the first field (myself)
fieldLoop:
	for i := 1; i < pt.NumField(); i++ {
		field := pt.Field(i)
		yTags := strings.Split(field.Tag.Get("yaml"), ",")

		// get yaml field name
		yamlKey := yTags[0]
		if len(yamlKey) != 0 {
			if !addField(yamlKey, pt.Field(i).Name, f._parentValue.Elem().Field(i)) {
				return fmt.Errorf(
					"field: duplicate yaml key %q in %s",
					yamlKey, pt.String(),
				)
			}
		}

		// process yaml tag flags
		for _, t := range yTags[1:] {
			switch t {
			case "-":
				ignoreField(yamlKey)
				continue fieldLoop
			case "inline":
				kind := field.Type.Kind()
				switch {
				case kind == reflect.Struct:
				case kind == reflect.Slice:
				case kind == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct:
				default:
					return fmt.Errorf(
						"field: non struct nor struct pointer field %s.%s has inline tag",
						pt.String(), field.Name,
					)
				}

				logger.V("inspecting inline fields", log.String("field", field.Name))

				inlineFv := f._parentValue.Elem().Field(i)
				inlineFt := f._parentValue.Type().Elem().Field(i).Type

				for j := 0; j < inlineFv.NumField(); j++ {
					innerFv := inlineFv.Field(j)
					innerFt := inlineFt.Field(j)

					innerYamlKey := strings.Split(innerFt.Tag.Get("yaml"), ",")[0]
					if len(innerYamlKey) == 0 {
						// already in a inline field, do not check inline anymore
						continue
					}

					if !addField(innerYamlKey, innerFt.Name, innerFv) {
						return fmt.Errorf(
							"field: duplicate yaml key %q in inline field %s",
							innerYamlKey, pt.String(),
						)
					}
				}
			}
		}

		// dukkha tag is used to extend yaml tag
		dTags := strings.Split(field.Tag.Get("dukkha"), ",")
		for _, t := range dTags {
			// nolint:gocritic
			switch t {
			case "other":
				// match unhandled values
				if catchOtherField != nil {
					return fmt.Errorf(
						"field: bad field tags in %s: only one struct field can have `dukkha:\"other\"` tag",
						pt.String(),
					)
				}

				logger.V("found catch other field", log.String("field", field.Name))
				catchOtherField = &fieldSpec{
					fieldName:  pt.Field(i).Name,
					fieldValue: f._parentValue.Elem().Field(i),
				}
			}
		}
	}

	switch n.ShortTag() {
	case "!!map":
	default:
		return fmt.Errorf("field: unsupported yaml tag %q when handling %s", n.Tag, pt.String())
	}

	dataBytes, err := yaml.Marshal(n)
	if err != nil {
		return fmt.Errorf("field: data marshal back failed for %s: %w", pt.String(), err)
	}

	m := make(map[string]interface{})
	err = yaml.Unmarshal(dataBytes, &m)
	if err != nil {
		return fmt.Errorf("field: data unmarshal failed for %s: %w", pt.String(), err)
	}

	handledYamlValues := make(map[string]struct{})
	// handle rendering suffix
	for rawYamlKey, v := range m {
		yamlKey := rawYamlKey

		logger := logger.WithFields(log.String("raw_yaml_field", rawYamlKey))

		logger.V("inspecting yaml field")

		parts := strings.SplitN(rawYamlKey, "@", 2)
		if len(parts) == 1 {
			if _, ok := handledYamlValues[yamlKey]; ok {
				return fmt.Errorf(
					"field: duplicate yaml field name %q",
					yamlKey,
				)
			}

			// no rendering suffix, fill value

			handledYamlValues[yamlKey] = struct{}{}

			fSpec := getField(yamlKey)
			if fSpec == nil {
				if catchOtherField == nil {
					return fmt.Errorf("field: unknown yaml field %q for %s", yamlKey, pt.String())
				}

				fSpec = catchOtherField
				v = map[string]interface{}{
					yamlKey: v,
				}
			}

			logger := logger.WithFields(log.String("field", fSpec.fieldName))

			logger.V("working on plain field")

			err = unmarshal(yamlKey, v, fSpec.fieldValue)
			if err != nil {
				panic(err)
			}

			continue
		}

		// has rendering suffix

		yamlKey, renderer := parts[0], parts[1]

		logger = logger.WithFields(
			log.String("yaml_field", yamlKey),
			log.String("renderer", renderer),
		)

		if _, ok := handledYamlValues[yamlKey]; ok {
			return fmt.Errorf(
				"field: duplicate yaml field name %q, rendering suffix won't change the field name",
				yamlKey,
			)
		}

		rawData, ok := v.(string)
		if !ok {
			return fmt.Errorf(
				"field: expecting string value for field %q (using rendering suffix), got %T",
				rawYamlKey, v,
			)
		}

		fSpec := getField(yamlKey)
		if fSpec == nil {
			if catchOtherField == nil {
				return fmt.Errorf("field: unknown yaml field %q for %s", yamlKey, pt.String())
			}

			// TODO: handle catch all
			fSpec = catchOtherField
			v = map[string]interface{}{
				yamlKey: v,
			}
			_ = v
		}

		logger = logger.WithFields(log.String("field", fSpec.fieldName))

		// do not unmarshal now, we need to evaluate value and unmarshal
		//
		// 		err = unmarshal(v, fSpec.fieldValue)
		//

		handledYamlValues[yamlKey] = struct{}{}
		// don't forget the raw name with rendering suffix
		handledYamlValues[rawYamlKey] = struct{}{}

		logger.V("found field to be rendered")

		f.addUnresolvedField(
			fSpec.fieldName, fSpec.fieldValue,
			yamlKey,
			renderer, rawData,
		)
	}

	for k := range handledYamlValues {
		delete(m, k)
	}

	if len(m) == 0 {
		// all values consumed
		return nil
	}

	if catchOtherField == nil {
		var unknownFields []string
		for k := range m {
			unknownFields = append(unknownFields, k)
		}
		sort.Strings(unknownFields)

		return fmt.Errorf(
			"field: unknown yaml fields for %s: %s",
			pt.String(), strings.Join(unknownFields, ", "),
		)
	}

	for k, v := range m {
		// TODO: fill values to catchOtherField
		_, _ = k, v
	}

	return nil
}

func unmarshal(yamlKey string, in interface{}, outField reflect.Value) error {
	oe := outField

fieldLoop:
	for {
		switch oe.Kind() {
		case reflect.Slice:
			inSlice := in.([]interface{})
			size := len(inSlice)

			sliceVal := reflect.MakeSlice(oe.Type(), size, size)

			for i := 0; i < size; i++ {
				itemVal := sliceVal.Index(i)

				err := unmarshal(yamlKey, inSlice[i], itemVal)
				if err != nil {
					return fmt.Errorf("failed to unmarshal slice item %s: %w", itemVal.Type().String(), err)
				}
			}

			if oe.IsZero() {
				oe.Set(sliceVal)
			} else {
				oe.Set(reflect.AppendSlice(oe, sliceVal))
			}

			return nil
		case reflect.Map:
			inMap := reflect.ValueOf(in)

			mapVal := reflect.MakeMap(oe.Type())
			if oe.IsZero() {
				oe.Set(mapVal)
			}

			valType := oe.Type().Elem()

			iter := inMap.MapRange()
			for iter.Next() {
				valVal := reflect.New(valType)
				err := unmarshal(iter.Key().String(), iter.Value().Interface(), valVal)
				if err != nil {
					return fmt.Errorf("failed to unmarshal map value %s: %w", valType.String(), err)
				}

				oe.SetMapIndex(iter.Key(), valVal.Elem())
			}

			return nil
		case reflect.Interface:
			fVal, err := CreateInterfaceField(oe.Type(), yamlKey)
			if err != nil {
				return fmt.Errorf("failed to create interface field: %w", err)
			}

			val := reflect.ValueOf(fVal)
			outField.Set(val)

			// DO NOT use outField directly, which will always match reflect.Interface
			return unmarshal(yamlKey, in, val)
		case reflect.Ptr:
			// process later
		default:
			// scalar types or struct/array/func/chan/unsafe.Pointer
			break fieldLoop
		}

		if oe.Kind() != reflect.Ptr {
			break
		}

		if oe.IsZero() {
			oe.Set(reflect.New(oe.Type().Elem()))
		}

		oe = oe.Elem()
	}

	var outPtr reflect.Value
	if outField.Kind() != reflect.Ptr {
		outPtr = outField.Addr()
	} else {
		outPtr = outField
	}

	fVal, canCallInit := outPtr.Interface().(Interface)
	if canCallInit {
		_ = Init(fVal)
	}

	dataBytes, err := yaml.Marshal(in)
	if err != nil {
		return fmt.Errorf("field: failed to marshal back for plain field: %w", err)
	}

	return yaml.Unmarshal(dataBytes, outPtr.Interface())
}
