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
	_parentValue reflect.Value `yaml:"-"`

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
		f.unresolvedFields[key].rawData = append(old.rawData, rawData)
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
		return fmt.Errorf("field: struct not intialized with field.New()")
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

		yamlKey := yTags[0]
		if len(yamlKey) != 0 {
			if !addField(yamlKey, pt.Field(i).Name, f._parentValue.Elem().Field(i)) {
				return fmt.Errorf(
					"field: duplicate yaml key %q in %s",
					yamlKey, pt.String(),
				)
			}
		}

		for _, t := range yTags[1:] {
			if t == "-" {
				ignoreField(yamlKey)
				continue fieldLoop
			}

			if t == "inline" {
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

				// TODO: add inner fields
				logger.V("inspecting inline fields", log.String("field", field.Name))
				f._parentValue.Elem().Field(i)
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

	dataBytes, err := yaml.Marshal(n)
	if err != nil {
		return fmt.Errorf("field: data marshal back failed for %s: %w", pt.String(), err)
	}

	switch n.Tag {
	case "!!seq":
		// TODO
		return nil
	case "!!map":
	default:
		return fmt.Errorf("field: unsupported yaml tag %q when handling %s", n.Tag, pt.String())
	}

	m := make(map[string]interface{})
	err = yaml.Unmarshal(dataBytes, &m)
	if err != nil {
		return fmt.Errorf("field: data unmarshal failed for %s: %w", pt.String(), err)
	}

	handledYamlValues := make(map[string]struct{})
	// handle rendering suffix
	for k, v := range m {
		yamlKey := k

		logger := logger.WithFields(log.String("yaml_field", yamlKey))

		logger.V("inspecting yaml field")

		parts := strings.SplitN(k, "@", 2)
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
			}

			logger := logger.WithFields(log.String("field", fSpec.fieldName))

			logger.V("working on plain field")

			initializeBaseField(logger, fSpec.fieldValue)

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
				"field: expecting string value for field %q using rendering suffix, got %T",
				k, v,
			)
		}

		fSpec := getField(yamlKey)
		if fSpec == nil {
			if catchOtherField == nil {
				return fmt.Errorf("field: unknown yaml field %q for %s", yamlKey, pt.String())
			}

			fSpec = catchOtherField
		}

		// TODO: initialize struct with field.New()

		logger = logger.WithFields(log.String("field", fSpec.fieldName))
		if fSpec.fieldValue.CanInterface() {
			fVal, hasBaseField := fSpec.fieldValue.Interface().(Interface)
			if hasBaseField {
				logger.V("initializing field with rendering suffix")
				_ = New(fVal)
			}
		}

		handledYamlValues[yamlKey] = struct{}{}
		// don't forget the raw name with rendering suffix
		handledYamlValues[k] = struct{}{}

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

func initializeBaseField(logger log.Interface, f reflect.Value) {
	if !f.CanInterface() {
		return
	}

	fVal, hasBaseField := f.Interface().(Interface)
	if !hasBaseField {
		return
	}

	logger.V("initializing field using BaseField",
		log.Any("kind", f.Kind()),
		log.Any("elem", f.Elem()),
	)

	switch f.Kind() {
	case reflect.Struct:
	case reflect.Ptr:
		f.Set(reflect.New(f.Type().Elem()))
	default:
		panic(fmt.Sprintf("unexpected type: %q", f.Kind()))
	}

	_ = New(fVal)
}
