package docs_test

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
	_ "unsafe" // for go:linkname

	"github.com/huandu/xstrings"
	"github.com/stretchr/testify/assert"
	"github.com/weaveworks/schemer/definition"
	"github.com/weaveworks/schemer/schema"
	"golang.org/x/tools/imports"

	"arhat.dev/dukkha/pkg/conf"
	"arhat.dev/dukkha/pkg/dukkha"

	_ "arhat.dev/dukkha/cmd/dukkha/addon" // add types disabled by default
	_ "arhat.dev/dukkha/pkg/cmd"          // add types enabled by default
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

	var configModelFields []reflect.StructField
	dukkhaConfigType := reflect.TypeOf((*conf.Config)(nil)).Elem()
	for i := 1; i < dukkhaConfigType.NumField(); i++ {
		fi := dukkhaConfigType.Field(i)
		if fi.PkgPath != "" {
			// unexported
			continue
		}

		yamlKey := strings.Split(fi.Tag.Get("yaml"), ",")[0]
		switch yamlKey {
		case "-", "", "renderers", "tools":
			continue
		}

		configModelFields = append(configModelFields, fi)
	}

	configModelFields = append(configModelFields, reflect.StructField{
		Name: "Renderers",
		Type: reflect.SliceOf(reflect.StructOf(rdrs)),
		Tag:  reflect.StructTag(`yaml:"renderers"`),
	}, reflect.StructField{
		Name: "Tools",
		Type: reflect.StructOf(tools),
		Tag:  reflect.StructTag(`yaml:"tools"`),
	})

	for _, t := range _tasks {
		configModelFields = append(configModelFields, reflect.StructField{
			Name: xstrings.FirstRuneToUpper(xstrings.ToCamelCase(strings.ReplaceAll(t.Name, ":", "_"))),
			Type: reflect.SliceOf(t.Typ),
			Tag:  reflect.StructTag(fmt.Sprintf(`yaml:"%s"`, t.Name)),
		})
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

const (
	// only match when first renderer with a patch suffix
	patchSuffixPatternFormat = `^%s@[^\|]*!`

	renderingSuffixPatternFormat = `^%s@.*`
)

func TestPatchSpecPattern(t *testing.T) {
	tests := []struct {
		str   string
		match bool
	}{
		{"foo@!", true},
		{"foo@test!", true},
		{"foo@test!|env", true},
		{"foo@test?int!", true},
		{"foo@test?int!|env", true},

		// should not match
		{"foo", false},
		{"foo@", false},
		{"foo@?int", false},
		{"foo@test|env!", false},
		{"foo@env|int!", false},
	}
	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			match, err := regexp.MatchString(fmt.Sprintf(patchSuffixPatternFormat, "foo"), test.str)
			assert.NoError(t, err)
			assert.Equal(t, test.match, match)
		})
	}
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
		case strings.Contains(k, ":"): // is task schema
			// tasks can be tool specific
			parts := strings.SplitN(k, ":", 2)
			taskPattern := fmt.Sprintf(`^%s(:.+){0,1}:%s$`, parts[0], parts[1])

			if topLevelStruct.PatternProperties == nil {
				topLevelStruct.PatternProperties = make(map[string]*definition.Definition)
			}

			topLevelStruct.PatternProperties[taskPattern] = def
		case k == "renderers": // is renderers
			def.Items.PatternProperties = make(map[string]*definition.Definition, len(def.Items.Properties))
			for rdr, rdrDef := range def.Items.Properties {
				rdrPattern := fmt.Sprintf("^%s(:.+){0,1}$", rdr)
				def.Items.PatternProperties[rdrPattern] = rdrDef
			}
		default:
		}
	}

	// include PatchSpec, so we can add patch spec for all definitions
	// using patternProperties
	psScm, err := schema.GenerateSchema("arhat.dev/rs", "PatchSpec", "yaml", formatRefName, false)
	if err != nil {
		return nil, err
	}

	psDef := psScm.Ref
	for k, def := range psScm.Definitions {
		scm.Definitions[k] = def
	}

	for kind, def := range scm.Definitions {
		if len(def.Properties) == 0 {
			continue
		}

		if def.PatternProperties == nil {
			def.PatternProperties = make(map[string]*definition.Definition, len(def.Properties))
		}

		// Add additional pattern properties while keep plain properties to make
		// autocompletion happy
		for name, prop := range def.Properties {
			switch {
			case name == "name":
				// name property exists in tasks and tools
				// which usually doesn't allow rendering suffix
				continue
			case kind == "Schema" && name == "include":
				// top-level `include` doesn't have rendering suffix support
				continue
			}

			// TODO: add rendering suffix aware patterns
			// 		 currently only have patch spec support

			patchPattern := fmt.Sprintf(patchSuffixPatternFormat, name)
			def.PatternProperties[patchPattern] = &definition.Definition{
				Ref: psDef,
			}

			rsPattern := fmt.Sprintf(renderingSuffixPatternFormat, name)
			def.PatternProperties[rsPattern] = prop

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
