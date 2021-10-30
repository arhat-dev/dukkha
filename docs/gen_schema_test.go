package docs_test

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	_ "unsafe" // for go:linkname

	"arhat.dev/dukkha/pkg/dukkha"
	"github.com/huandu/xstrings"
	"github.com/stretchr/testify/assert"
	"github.com/weaveworks/schemer/definition"
	"github.com/weaveworks/schemer/schema"
	"golang.org/x/tools/imports"

	_ "arhat.dev/dukkha/cmd/dukkha/addon" // add types disabled by default
	_ "arhat.dev/dukkha/pkg/conf"         // add types enabled by default
)

func TestGenerateSchema(t *testing.T) {
	_renderers, _tools, _tasks := collectSchema()

	rdrs := make([]reflect.StructField, len(_renderers))
	for i, r := range _renderers {
		rdrs[i] = reflect.StructField{
			Name: xstrings.FirstRuneToUpper(xstrings.ToCamelCase(r.Name)),
			Type: r.Typ,
			Tag:  reflect.StructTag(fmt.Sprintf(`yaml:"%s"`, r.Name)),
		}
	}

	tools := make([]reflect.StructField, len(_tools))
	for i, t := range _tools {
		tools[i] = reflect.StructField{
			Name: xstrings.FirstRuneToUpper(xstrings.ToCamelCase(t.Name)),
			Type: reflect.SliceOf(t.Typ),
			Tag:  reflect.StructTag(fmt.Sprintf(`yaml:"%s"`, t.Name)),
		}
	}

	rdrsTyp := reflect.StructOf(rdrs)
	configModelFields := make([]reflect.StructField, len(_tasks)+2)
	configModelFields[0] = reflect.StructField{
		Name: "Renderers",
		Type: rdrsTyp,
		Tag:  reflect.StructTag(`yaml:"renderers"`),
	}

	toolsTyp := reflect.StructOf(tools)
	configModelFields[1] = reflect.StructField{
		Name: "Tools",
		Type: toolsTyp,
		Tag:  reflect.StructTag(`yaml:"tools"`),
	}

	for i, t := range _tasks {
		configModelFields[i+2] = reflect.StructField{
			Name: xstrings.FirstRuneToUpper(xstrings.ToCamelCase(strings.ReplaceAll(t.Name, ":", "_"))),
			Type: reflect.SliceOf(t.Typ),
			Tag:  reflect.StructTag(fmt.Sprintf(`yaml:"%s"`, t.Name)),
		}
	}

	structStr := strings.ReplaceAll(
		reflect.StructOf(configModelFields).String(), `"yaml:\"`, "`yaml:\"")
	structStr = strings.ReplaceAll(structStr, `\""`, "\"`")

	structStr = `package generated

type Schema ` + structStr

	// imports.Debug = true
	imports.LocalPrefix = "arhat.dev/dukkha"
	ret, err := imports.Process("", []byte(structStr), &imports.Options{})
	if !assert.NoError(t, err) {
		return
	}

	err = os.WriteFile("generated/schema.go", ret, 0644)
	if !assert.NoError(t, err) {
		return
	}

	schemaBytes, err := generateSchemaJSON("./generated", "Schema")
	if !assert.NoError(t, err) {
		return
	}

	err = os.WriteFile("generated/schema.json", schemaBytes, 0644)
	if !assert.NoError(t, err) {
		return
	}
}

func formatRefName(pkg, name string) string {
	return strings.ReplaceAll(pkg, "/", ".") + "." + name
}

func generateSchemaJSON(pkgPath, topLevelStructName string) ([]byte, error) {
	scm, err := schema.GenerateSchema(
		pkgPath, topLevelStructName, "yaml", formatRefName, false,
	)
	if err != nil {
		return nil, err
	}

	topLevelStruct := scm.Definitions[topLevelStructName]
	for k, def := range topLevelStruct.Properties {
		switch {
		case strings.Contains(k, ":"):
			// tasks can be tool specific
			parts := strings.SplitN(k, ":", 2)
			taskPattern := fmt.Sprintf(`^%s(:.+){0,1}:%s$`, parts[0], parts[1])

			if topLevelStruct.PatternProperties == nil {
				topLevelStruct.PatternProperties = make(map[string]*definition.Definition)
			}

			topLevelStruct.PatternProperties[taskPattern] = def
		}
	}

	psScm, err := schema.GenerateSchema("arhat.dev/rs", "PatchSpec", "yaml", formatRefName, false)
	if err != nil {
		return nil, err
	}

	psDef := psScm.Ref
	for k, def := range psScm.Definitions {
		scm.Definitions[k] = def
	}

	// TODO: add rendering suffix aware patterns
	// 		 currently only have patch spec support
	for _, def := range scm.Definitions {
		if len(def.Properties) == 0 {
			continue
		}

		if def.PatternProperties == nil {
			def.PatternProperties = make(map[string]*definition.Definition, len(def.Properties))
		}

		for name := range def.Properties {
			if name == "name" {
				continue
			}

			patchSpecPattern := fmt.Sprintf(`^%s@((.*\?.+\|?)+)?!$`, name)
			def.PatternProperties[patchSpecPattern] = &definition.Definition{
				Ref: psDef,
			}

			// pattern := fmt.Sprintf(`^%s(@(.*\?.+\|?)+)?$`, name)
			// patternProps[pattern] = prop
		}
	}

	return schema.ToJSON(scm)
}

var (
	rendererType = reflect.TypeOf((*dukkha.Renderer)(nil)).Elem()
	toolType     = reflect.TypeOf((*dukkha.Tool)(nil)).Elem()
	taskType     = reflect.TypeOf((*dukkha.Task)(nil)).Elem()
)

type typeSchema struct {
	Name string
	Typ  reflect.Type
}

//go:linkname gm arhat.dev/dukkha/pkg/dukkha.globalTypeManager
var gm *dukkha.TypeManager

func collectSchema() (renderers, tools, tasks []*typeSchema) {
	renderers = createSchema(
		gm.Types()[dukkha.IfaceTypeKey{
			Typ: rendererType,
		}].Factories,
	)

	tools = createSchema(
		gm.Types()[dukkha.IfaceTypeKey{
			Typ: toolType,
		}].Factories,
	)

	tasks = createSchema(
		gm.Types()[dukkha.IfaceTypeKey{
			Typ: taskType,
		}].Factories,
	)

	return
}

func createSchema(impl []*dukkha.IfaceFactoryImpl) []*typeSchema {
	ret := make([]*typeSchema, len(impl))
	for i, f := range impl {
		instance := f.Create(nil)
		val := reflect.ValueOf(instance)

		for val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		switch {
		case val.Kind() == reflect.Interface:
			val = val.Elem()

			for val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
		case val.Kind() == reflect.Struct:
		default:
			panic("unexpected non struct value")
		}

		println("TYP:", f.Name, val.Type().String())
		ret[i] = &typeSchema{
			Name: f.Name,
			Typ:  val.Type(),
		}
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Name < ret[j].Name
	})

	return ret
}
