package field

import (
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

func (f *BaseField) ResolveFields(ctx *RenderingContext, render RenderingFunc, depth int) error {
	if atomic.LoadUint32(&f._initialized) == 0 {
		return fmt.Errorf("field resolve: struct not intialized with Init()")
	}

	if depth == 0 {
		return nil
	}

	logger := log.Log.WithName("BaseField").
		WithFields(
			log.String("func", "Render"),
			log.String("struct", f._parentValue.Type().String()),
		)

	logger.D("resolving",
		log.Int("unresolved_fields", len(f.unresolvedFields)),
	)

	for k, v := range f.unresolvedFields {
		logger := logger.WithFields(
			log.String("field", k.fieldName),
			log.String("type", v.fieldValue.Type().String()),
			log.String("yaml_field", v.yamlFieldName),
		)

		logger.V("resolving", log.Any("values", ctx.Values()))

		var target interface{}
		switch v.fieldValue.Kind() {
		case reflect.Ptr:
			target = v.fieldValue.Interface()
		default:
			target = v.fieldValue.Addr().Interface()
		}

		for _, rawData := range v.rawData {
			resolvedValue, err := render(ctx, k.renderer, rawData)
			if err != nil {
				return fmt.Errorf("field: failed to render value of this base field: %w", err)
			}

			err = yaml.Unmarshal([]byte(resolvedValue), target)
			if err != nil {
				return fmt.Errorf("field: failed to unmarshal resolved value: %w", err)
			}

			logger.V("resolved field", log.Any("value", target))
		}

		if depth > 1 || depth < 0 {
			innerF, canCallResolve := target.(Interface)
			if !canCallResolve {
				continue
			}

			err := innerF.ResolveFields(ctx, render, depth-1)
			if err != nil {
				return fmt.Errorf("failed to resolve inner field: %w", err)
			}
		}
	}

	return nil
}

func (f *BaseField) addUnresolvedField(
	fieldName string,
	fieldValue reflect.Value,
	yamlKey string,
	renderer, rawData string,
) error {
	if f.unresolvedFields == nil {
		f.unresolvedFields = make(map[unresolvedFieldKey]*unresolvedFieldValue)
	}

	key := unresolvedFieldKey{
		fieldName: fieldName,
		renderer:  renderer,
	}

	oe := fieldValue
	for {
		switch oe.Kind() {
		case reflect.Slice:
			oe.Set(reflect.MakeSlice(oe.Type(), 0, 0))
		case reflect.Map:
			oe.Set(reflect.MakeMap(oe.Type()))
		case reflect.Interface:
			fVal, err := CreateInterfaceField(oe.Type(), yamlKey)
			if err != nil {
				return fmt.Errorf("failed to create interface field: %w", err)
			}

			oe.Set(reflect.ValueOf(fVal))
		case reflect.Ptr:
			// process later
		default:
			// scalar types or struct/array/func/chan/unsafe.Pointer
			// hand it to go-yaml
		}

		if oe.Kind() != reflect.Ptr {
			break
		}

		if oe.IsZero() {
			oe.Set(reflect.New(oe.Type().Elem()))
		}

		oe = oe.Elem()
	}

	var iface interface{}
	switch fieldValue.Kind() {
	case reflect.Ptr:
		iface = fieldValue.Interface()
	default:
		iface = fieldValue.Addr().Interface()
	}

	fVal, canCallInit := iface.(Interface)
	if canCallInit {
		_ = Init(fVal)
	}

	if old, exists := f.unresolvedFields[key]; exists {
		old.rawData = append(old.rawData, rawData)
		return nil
	}

	f.unresolvedFields[key] = &unresolvedFieldValue{
		fieldValue:    fieldValue,
		yamlFieldName: yamlKey,
		rawData:       []string{rawData},
	}

	return nil
}

// UnmarshalYAML handles renderer suffix
// nolint:gocyclo,revive
func (self *BaseField) UnmarshalYAML(n *yaml.Node) error {
	if atomic.LoadUint32(&self._initialized) == 0 {
		return fmt.Errorf("field unmarshal: struct not intialized with Init()")
	}

	type fieldKey struct {
		yamlKey string
	}

	type fieldSpec struct {
		fieldName  string
		fieldValue reflect.Value
		base       *BaseField
	}

	fields := make(map[fieldKey]*fieldSpec)
	pt := self._parentValue.Type().Elem()

	addField := func(
		yamlKey, fieldName string,
		fieldValue reflect.Value,
		base *BaseField,
	) bool {
		key := fieldKey{yamlKey: yamlKey}
		if _, exists := fields[key]; exists {
			return false
		}

		fields[fieldKey{yamlKey: yamlKey}] = &fieldSpec{
			fieldName:  fieldName,
			fieldValue: fieldValue,
			base:       base,
		}
		return true
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
		fieldType := pt.Field(i)
		yTags := strings.Split(fieldType.Tag.Get("yaml"), ",")

		// check if ignored
		for _, t := range yTags[1:] {
			if t == "-" {
				continue fieldLoop
			}
		}

		// get yaml field name
		yamlKey := yTags[0]
		if len(yamlKey) != 0 {
			if !addField(yamlKey, fieldType.Name, self._parentValue.Elem().Field(i), self) {
				return fmt.Errorf(
					"field: duplicate yaml key %q in %s",
					yamlKey, pt.String(),
				)
			}
		}

		// process yaml tag flags
		for _, t := range yTags[1:] {
			switch t {
			case "inline":
				kind := fieldType.Type.Kind()
				switch {
				case kind == reflect.Struct:
				case kind == reflect.Ptr && fieldType.Type.Elem().Kind() == reflect.Struct:
				default:
					return fmt.Errorf(
						"field: non struct nor struct pointer field %s.%s has inline tag",
						pt.String(), fieldType.Name,
					)
				}

				logger.V("inspecting inline fields", log.String("field", fieldType.Name))

				inlineFv := self._parentValue.Elem().Field(i)
				inlineFt := self._parentValue.Type().Elem().Field(i).Type

				var iface interface{}
				switch inlineFv.Kind() {
				case reflect.Ptr:
					iface = inlineFv.Interface()
				default:
					iface = inlineFv.Addr().Interface()
				}

				base := self
				fVal, canCallInit := iface.(Interface)
				if canCallInit {
					innerBaseF := reflect.ValueOf(Init(fVal)).Elem().Field(0)

					if innerBaseF.Kind() == reflect.Struct {
						if innerBaseF.Addr().Type() == baseFieldPtrType {
							base = innerBaseF.Addr().Interface().(*BaseField)
						}
					} else {
						if innerBaseF.Type() == baseFieldPtrType {
							base = innerBaseF.Interface().(*BaseField)
						}
					}
				}

				for j := 0; j < inlineFv.NumField(); j++ {
					innerFv := inlineFv.Field(j)
					innerFt := inlineFt.Field(j)

					innerYamlKey := strings.Split(innerFt.Tag.Get("yaml"), ",")[0]
					if len(innerYamlKey) == 0 {
						// already in a inline field, do not check inline anymore
						continue
					}

					if !addField(innerYamlKey, innerFt.Name, innerFv, base) {
						return fmt.Errorf(
							"field: duplicate yaml key %q in inline field %s",
							innerYamlKey, pt.String(),
						)
					}
				}
			default:
				// TODO: handle other yaml tag flags
			}
		}

		// dukkha tag is used to extend yaml tag
		dTags := strings.Split(fieldType.Tag.Get("dukkha"), ",")
		for _, t := range dTags {
			// nolint:gocritic
			switch t {
			case "other":
				// other is used to match unhandled values
				// only supports map[string]Any

				if catchOtherField != nil {
					return fmt.Errorf(
						"field: bad field tags in %s: only one map in a struct can have `dukkha:\"other\"` tag",
						pt.String(),
					)
				}

				logger.V("found catch other field", log.String("field", fieldType.Name))
				catchOtherField = &fieldSpec{
					fieldName:  fieldType.Name,
					fieldValue: self._parentValue.Elem().Field(i),
					base:       self,
				}
			case "":
			default:
				return fmt.Errorf("field: unknown dukkha tag value %q", t)
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
				return fmt.Errorf(
					"field: failed to unmarshal yaml field %q to struct field %q: %w",
					yamlKey, fSpec.fieldName, err,
				)
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

			fSpec = catchOtherField
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

		err = fSpec.base.addUnresolvedField(
			fSpec.fieldName, fSpec.fieldValue,
			yamlKey,
			renderer, rawData,
		)
		if err != nil {
			return fmt.Errorf("field: failed to add unresolved field: %w", err)
		}
	}

	for k := range handledYamlValues {
		delete(m, k)
	}

	if len(m) == 0 {
		// all values consumed
		return nil
	}

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

func unmarshal(yamlKey string, in interface{}, outField reflect.Value) error {
	oe := outField

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
			// map key MUST be string

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
			// hand it to go-yaml
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
		return fmt.Errorf("field: failed to marshal back yaml field %q: %w", yamlKey, err)
	}

	return yaml.Unmarshal(dataBytes, outPtr.Interface())
}
