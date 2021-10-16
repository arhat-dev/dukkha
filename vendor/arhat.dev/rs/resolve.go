package rs

import (
	"fmt"
	"reflect"
	"sync/atomic"

	"gopkg.in/yaml.v3"
)

func (f *BaseField) HasUnresolvedField() bool {
	return len(f.unresolvedFields) != 0
}

func (f *BaseField) ResolveFields(rc RenderingHandler, depth int, fieldNames ...string) error {
	if atomic.LoadUint32(&f._initialized) == 0 {
		return fmt.Errorf("rs: struct not intialized before resolving")
	}

	if depth == 0 {
		return nil
	}

	if len(fieldNames) == 0 {
		// resolve all

		for i := 1; i < f._parentValue.NumField(); i++ {
			sf := f._parentType.Field(i)
			if len(sf.PkgPath) != 0 {
				// not exported
				continue
			}

			err := f.resolveSingleField(
				rc, depth, sf.Name, f._parentValue.Field(i),
			)
			if err != nil {
				return fmt.Errorf(
					"rs: failed to resolve field %s.%s: %w",
					f._parentType.String(), sf.Name, err,
				)
			}
		}

		return nil
	}

	for _, name := range fieldNames {
		fv := f._parentValue.FieldByName(name)
		if !fv.IsValid() {
			return fmt.Errorf(
				"rs: no such field %q in struct %q",
				name, f._parentType.String(),
			)
		}

		err := f.resolveSingleField(rc, depth, name, fv)
		if err != nil {
			return fmt.Errorf(
				"rs: failed to resolve requested field %s.%s: %w",
				f._parentType.String(), name, err,
			)
		}
	}

	return nil
}

func (f *BaseField) resolveSingleField(
	rc RenderingHandler,
	depth int,
	fieldName string,
	targetField reflect.Value,
) error {
	keepOld := false

	for k, v := range f.unresolvedFields {
		if v.fieldName == fieldName {
			err := f.handleUnResolvedField(
				rc, depth, k, v, keepOld,
			)
			if err != nil {
				return err
			}

			keepOld = true
		}
	}

	return f.handleResolvedField(rc, depth, targetField)
}

// nolint:gocyclo
func (f *BaseField) handleResolvedField(
	rc RenderingHandler,
	depth int,
	field reflect.Value,
) error {
	if depth == 0 {
		return nil
	}

	switch fk := field.Kind(); fk {
	case reflect.Map:
		if field.IsNil() {
			return nil
		}

		iter := field.MapRange()
		for iter.Next() {
			err := f.handleResolvedField(rc, depth-1, iter.Value())
			if err != nil {
				return err
			}
		}
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		if fk == reflect.Slice && field.IsNil() {
			// this is a resolved field, slice/array empty means no value
			return nil
		}

		for i := 0; i < field.Len(); i++ {
			tt := field.Index(i)
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
		if !field.IsValid() || field.IsZero() || field.IsNil() {
			return nil
		}
		// handled after switch
	default:
		// scalar types, no action required
		return nil
	}

	return tryResolve(rc, depth, field)
}

func tryResolve(rc RenderingHandler, depth int, targetField reflect.Value) error {
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
	yamlKey string,
	v *unresolvedFieldSpec,
	keepOld bool,
) error {
	target := v.fieldValue

	for i, rawData := range v.rawDataList {
		toResolve := rawData
		if v.isCatchOtherField {
			// unwrap map data for resolving
			toResolve = rawData.mapData[yamlKey]
		}

		var (
			err error
		)

		for _, renderer := range v.renderers {

			// a patch is implied when the renderer has a `!` suffix
			var patchSpec *renderingPatchSpec
			if renderer.patchSpec {
				patchSpec, toResolve, err = f.resolvePatchSpec(rc, toResolve)
				if err != nil {
					return fmt.Errorf(
						"failed to resolve patch spec for renderer %q: %w",
						renderer.name, err,
					)
				}
			}

			if len(renderer.name) != 0 {
				var tmp interface{}

				// toResolve can be nil when no patch value is set
				if toResolve != nil {
					tmp = toResolve.NormalizedValue()
				}

				tmp, err = rc.RenderYaml(renderer.name, tmp)
				if err != nil {
					return fmt.Errorf(
						"renderer %q failed to render value: %w",
						renderer.name, err,
					)
				}

				toResolve = &alterInterface{
					scalarData: tmp,
				}
			}

			// apply patch if set
			if patchSpec != nil {
				// apply type hint before patching to make sure value
				// being patched is correctly typed
				hint := renderer.typeHint
				toResolve, err = applyTypeHint(hint, toResolve)
				if err != nil {
					return fmt.Errorf(
						"failed to ensure type hint %q on yaml key %q: %w",
						hint, yamlKey, err,
					)
				}

				var tmp interface{}
				tmp, err = patchSpec.ApplyTo(toResolve.NormalizedValue())
				if err != nil {
					return fmt.Errorf(
						"failed to apply patches: %w",
						err,
					)
				}

				toResolve = &alterInterface{
					scalarData: tmp,
				}
			}

			// apply hint after resolving (rendering)
			hint := renderer.typeHint
			toResolve, err = applyTypeHint(hint, toResolve)
			if err != nil {
				return fmt.Errorf(
					"failed to ensure type hint %q on yaml key %q: %w",
					hint, yamlKey, err,
				)
			}
		}

		resolved := toResolve.NormalizedValue()
		if v.isCatchOtherField {
			// wrap back for catch other filed
			resolved = map[string]interface{}{
				yamlKey: resolved,
			}
		}

		// TODO: currently we always keepOld when the field has tag
		// 		 `rs:"other"`, need to ensure this behavior won't
		// 	     leave inconsistent data

		actualKeepOld := keepOld || v.isCatchOtherField || i != 0
		err = f.unmarshal(yamlKey, reflect.ValueOf(resolved), target, actualKeepOld)
		if err != nil {
			return fmt.Errorf(
				"failed to unmarshal resolved value of yaml key %q to field %q: %w",
				yamlKey, v.fieldName, err,
			)
		}
	}

	return tryResolve(rc, depth-1, target)
}

// resolve user provided data as patch spec
func (f *BaseField) resolvePatchSpec(
	rc RenderingHandler,
	toResolve interface{},
) (
	patchSpec *renderingPatchSpec,
	value *alterInterface,
	err error,
) {
	var patchSpecBytes []byte
	switch t := toResolve.(type) {
	case string:
		patchSpecBytes = []byte(t)
	case []byte:
		patchSpecBytes = t
	default:
		patchSpecBytes, err = yaml.Marshal(toResolve)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"failed to marshal renderer data for patch spec: %w",
				err,
			)
		}
	}

	patchSpec = Init(&renderingPatchSpec{}, f.ifaceTypeHandler).(*renderingPatchSpec)
	err = yaml.Unmarshal(patchSpecBytes, patchSpec)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to unmarshal patch spec %q: %w",
			string(patchSpecBytes), err,
		)
	}

	err = patchSpec.ResolveFields(rc, -1)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to resolve patch spec %q: %w",
			string(patchSpecBytes), err,
		)
	}

	value = convertAnyObjectToAlterInterface(&patchSpec.Value)
	return patchSpec, value, nil
}

func (f *BaseField) isCatchOtherField(yamlKey string) bool {
	if f.catchOtherFields == nil {
		return false
	}

	_, ok := f.catchOtherFields[yamlKey]
	return ok
}
