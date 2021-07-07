package field

import (
	"fmt"
	"reflect"
	"sync/atomic"

	"arhat.dev/pkg/log"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/types"
)

func (f *BaseField) HasUnresolvedField() bool {
	return len(f.unresolvedFields) != 0
}

func (f *BaseField) ResolveFields(rc types.RenderingContext, depth int, fieldName string) error {
	if atomic.LoadUint32(&f._initialized) == 0 {
		return fmt.Errorf("field resolve: struct not intialized with Init()")
	}

	if depth == 0 {
		return nil
	}

	structName := f._parentValue.Type().String()
	logger := log.Log.WithName("BaseField").
		WithFields(
			log.String("func", "Render"),
			log.String("struct", structName),
		)

	if len(fieldName) != 0 {
		// has target field
		logger = logger.WithFields(
			log.String("target", fieldName),
		)

		logger.D("trying to resolve specified single field")
		for k, v := range f.unresolvedFields {
			logger.V("looking up unresolved fields",
				log.String("met", k.fieldName),
			)

			if k.fieldName != fieldName {
				continue
			}

			logger = logger.WithFields(
				log.String("field", k.fieldName),
				log.String("type", v.fieldValue.Type().String()),
				log.String("yaml_field", v.yamlFieldName),
			)

			logger.D("resolving specified single field")

			return f.resolveSingleField(
				rc,
				logger,
				depth,
				structName,

				k.fieldName,
				k.renderer,

				v,
			)
		}

		logger.V("no such unresolved target single field")

		return nil
	}

	logger.D("resolving all fields",
		log.Int("count", len(f.unresolvedFields)),
	)

	return f.resolveAllFields(
		rc,
		logger,
		depth,
		structName,
	)
}

func (f *BaseField) resolveSingleField(
	rc types.RenderingContext,
	logger log.Interface,
	depth int,

	structName string, // to make error message helpful
	fieldName string, // to make error message helpful

	renderer string,
	v *unresolvedFieldValue,
) error {
	var target reflect.Value
	switch v.fieldValue.Kind() {
	case reflect.Ptr:
		target = v.fieldValue
	default:
		target = v.fieldValue.Addr()
	}

	for i, rawData := range v.rawData {
		resolvedValue, err := rc.RenderYaml(renderer, rawData)
		if err != nil {
			input, ok := rawData.(string)
			if !ok {
				inputBytes, err2 := yaml.Marshal(rawData)
				if err2 == nil {
					input = string(inputBytes)
				} else {
					input = fmt.Sprint(rawData)
				}
			}

			return fmt.Errorf(
				"field: failed to render value of %s.%s from\n\n%s\nerror: %w",
				structName, fieldName, input, err,
			)
		}

		if target.Type() == stringPtrType {
			// resolved value is the target value
			target.Elem().SetString(resolvedValue)
			continue
		}

		var tmp interface{}
		err = yaml.Unmarshal([]byte(resolvedValue), &tmp)
		if err != nil {
			logger.V("failed to unmarshal resolved value as interface",
				log.String("value", resolvedValue),
			)
			return fmt.Errorf("field: failed to unmarshal resolved value to interface: %w", err)
		}

		err = f.unmarshal(v.yamlFieldName, tmp, target, i != 0)
		if err != nil {
			return fmt.Errorf("field: failed to unmarshal resolved value %T: %w", target, err)
		}

		logger.V("resolved field", log.Any("value", target))
	}

	if depth > 1 || depth < 0 {
		innerF, canCallResolve := target.Interface().(types.Field)
		if !canCallResolve {
			return nil
		}

		err := innerF.ResolveFields(
			rc, depth-1, "",
		)
		if err != nil {
			return fmt.Errorf("failed to resolve inner field: %w", err)
		}
	}

	return nil
}

func (f *BaseField) resolveAllFields(
	rc types.RenderingContext,
	logger log.Interface,
	depth int,
	structName string, // to make error message helpful
) error {
	for k, v := range f.unresolvedFields {
		logger := logger.WithFields(
			log.String("field", k.fieldName),
			log.String("type", v.fieldValue.Type().String()),
			log.String("yaml_field", v.yamlFieldName),
		)

		logger.V("resolving single field", log.Any("values", rc))

		err := f.resolveSingleField(
			rc,
			logger,
			depth,

			structName,
			k.fieldName,

			k.renderer,
			v,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *BaseField) addUnresolvedField(
	fieldName string,
	fieldValue reflect.Value,
	yamlKey string,
	renderer string,
	rawData interface{},
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
			fVal, err := f.ifaceTypeHandler.Create(oe.Type(), yamlKey)
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

	fVal, canCallInit := iface.(types.Field)
	if canCallInit {
		_ = Init(fVal, f.ifaceTypeHandler)
	}

	if old, exists := f.unresolvedFields[key]; exists {
		old.rawData = append(old.rawData, rawData)
		return nil
	}

	f.unresolvedFields[key] = &unresolvedFieldValue{
		fieldValue:    fieldValue,
		yamlFieldName: yamlKey,
		rawData:       []interface{}{rawData},
	}

	return nil
}
