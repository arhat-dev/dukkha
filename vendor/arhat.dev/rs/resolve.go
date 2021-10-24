package rs

import (
	"encoding/json"
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"
)

func (f *BaseField) HasUnresolvedField() bool {
	return len(f.unresolvedFields) != 0
}

func (f *BaseField) ResolveFields(rc RenderingHandler, depth int, fieldNames ...string) error {
	if !f.initialized() {
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
			// fVal can have underlying data nil, but with Field type
			if targetField.IsNil() {
				return nil
			}

			return fVal.ResolveFields(rc, depth)
		}
	}

	if !targetField.CanAddr() {
		return nil
	}

	targetField = targetField.Addr()
	if !targetField.CanInterface() {
		return nil
	}

	fVal, canCallResolve := targetField.Interface().(Field)
	if canCallResolve {
		if targetField.IsNil() {
			return nil
		}

		return fVal.ResolveFields(rc, depth)
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
			toResolve = rawData.Content[1]
		}

		for _, renderer := range v.renderers {
			var (
				patchSpec *PatchSpec
				patchSrc  interface{}
				err       error
			)

			if renderer.patchSpec {
				patchSpec, patchSrc, err = f.resolvePatchSpec(rc, toResolve)
				if err != nil {
					return fmt.Errorf(
						"failed to resolve patch spec for renderer %q: %w",
						renderer.name, err,
					)
				}

				// should use patchSrc instead of toResolve since we need to patch
				toResolve = nil
			}

			if len(renderer.name) != 0 {
				var renderedData []byte

				var input interface{} = toResolve
				if patchSpec != nil {
					// resolve patch value if was a patch
					input = patchSrc
				}

				renderedData, err = rc.RenderYaml(renderer.name, input)
				if err != nil {
					return fmt.Errorf(
						"renderer %q failed to render value: %w",
						renderer.name, err,
					)
				}

				// check type hinting before assuming it's valid yaml
				//
				// see: TestResolve_yaml_unmarshal_invalid_but_no_error in resolve_test.go

				// scalar types cannot be applied with patch spec
				// so the rendered data will be the final value for resolving

				switch renderer.typeHint.(type) {
				case TypeHintStr:
					toResolve = &yaml.Node{
						Kind:  yaml.ScalarNode,
						Tag:   strTag,
						Value: string(renderedData),
					}
				case TypeHintInt:
					toResolve = &yaml.Node{
						Kind:  yaml.ScalarNode,
						Tag:   intTag,
						Value: string(renderedData),
					}
				case TypeHintFloat:
					toResolve = &yaml.Node{
						Kind:  yaml.ScalarNode,
						Tag:   floatTag,
						Value: string(renderedData),
					}
				default:
					// assume rendered data as yaml for further processing
					func() {
						defer func() {
							// TODO: yaml.Unmarshal can panic on invalid but seemingly
							// 		 correct input (e.g. markdown)
							//
							// related upstream issue:
							// 	https://github.com/go-yaml/yaml/issues/665
							errX := recover()

							if errX != nil {
								toResolve = &yaml.Node{
									Kind:  yaml.ScalarNode,
									Tag:   strTag,
									Value: string(renderedData),
								}
							}
						}()

						toResolve = new(yaml.Node)
						err = yaml.Unmarshal(renderedData, toResolve)
						if err != nil {
							toResolve = &yaml.Node{
								Kind:  yaml.ScalarNode,
								Tag:   strTag,
								Value: string(renderedData),
							}
						} else {
							// unmarshal ok
							if prepared := prepareYamlNode(toResolve); prepared != nil {
								toResolve = prepared
							}

							switch {
							case isStrScalar(toResolve), isBinaryScalar(toResolve):
								// use original string instead of yaml unmarshaled string
								// yaml.Unmarshal may modify string content when it's not
								// valid yaml

								toResolve = &yaml.Node{
									Kind:  yaml.ScalarNode,
									Tag:   strTag,
									Value: string(renderedData),
								}
							}
						}
					}()
				}
			}

			// apply patch if set
			if patchSpec != nil {
				// apply type hint before patching to make sure value
				// being patched is correctly typed

				var tmp interface{}
				if toResolve == nil {
					// toResolve can be nill if this is a patch without
					// renderer (e.g. `foo@!: { ... }`)
					tmp = patchSrc
				} else {
					hint := renderer.typeHint
					if hint == nil {
						hint = TypeHintNone{}
					}

					var hintedToResolve *yaml.Node

					// hintedToResolve can be nil
					hintedToResolve, err = hint.apply(toResolve)
					if err != nil {
						return fmt.Errorf(
							"failed to ensure type hint %q on patch target of %q: %w",
							hint, yamlKey, err,
						)
					}

					if hintedToResolve != nil {
						toResolve = hintedToResolve
					}

					err = toResolve.Decode(&tmp)
					if err != nil {
						return fmt.Errorf(
							"failed to decode data as patch source: %w",
							err,
						)
					}
				}

				tmp, err = patchSpec.ApplyTo(tmp)
				if err != nil {
					return fmt.Errorf(
						"failed to apply patches: %w",
						err,
					)
				}

				// patch doc is generated from arbitrary yaml data
				// with built-in interface{}, so we are able to marshal it into
				// json, and parse as *yaml.Node for further processing

				var dataBytes []byte
				dataBytes, err = json.Marshal(tmp)
				if err != nil {
					return fmt.Errorf("failed to marshal patched data: %w", err)
				}

				// json data is deemed to be valid yaml, if not, that means we
				// have a big problem then. so we don't need to save yaml.Unmarshal
				// from panic here
				toResolve = new(yaml.Node)
				err = yaml.Unmarshal(dataBytes, toResolve)
				if err != nil {
					return fmt.Errorf("failed to prepare patched data: %w", err)
				}
			}

			// apply hint after resolving (rendering)
			hint := renderer.typeHint
			if hint == nil {
				hint = TypeHintNone{}
			}

			var hintedToResolve *yaml.Node

			// hintedToResolve can be nil
			hintedToResolve, err = hint.apply(toResolve)
			if err != nil {
				return fmt.Errorf(
					"failed to ensure type hint %q on yaml key %q: %w",
					hint, yamlKey, err,
				)
			}

			if hintedToResolve != nil {
				toResolve = hintedToResolve
			}
		}

		resolved := toResolve
		if v.isCatchOtherField {
			// wrap back for catch other filed
			resolved = fakeMap(rawData.Content[0], resolved)
		}

		actualKeepOld := keepOld || v.isCatchOtherField || i != 0
		err := f.unmarshal(yamlKey, resolved, target, actualKeepOld)
		if err != nil {
			return fmt.Errorf(
				"failed to unmarshal resolved value of %q to field %q: %w",
				yamlKey, v.fieldName, err,
			)
		}
	}

	return tryResolve(rc, depth-1, target)
}

// resolve user provided data as patch spec
func (f *BaseField) resolvePatchSpec(
	rc RenderingHandler,
	toResolve *yaml.Node,
) (
	patchSpec *PatchSpec,
	value interface{},
	err error,
) {

	patchSpec = Init(&PatchSpec{}, f._opts).(*PatchSpec)
	err = toResolve.Decode(patchSpec)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to decode patch spec: %w",
			err,
		)
	}

	err = patchSpec.ResolveFields(rc, -1)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to resolve patch spec: %w",
			err,
		)
	}

	value = patchSpec.Value.NormalizedValue()
	return patchSpec, value, nil
}

func (f *BaseField) isCatchOtherField(yamlKey string) bool {
	if f.catchOtherFields == nil {
		return false
	}

	_, ok := f.catchOtherFields[yamlKey]
	return ok
}
