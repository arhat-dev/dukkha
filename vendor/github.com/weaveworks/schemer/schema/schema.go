/*
Copyright 2019 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Modifications:
  Copyright 2021 Weaveworks
*/

package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/weaveworks/schemer/definition"
	"github.com/weaveworks/schemer/importer"
)

const (
	version7 = "http://json-schema.org/draft-07/schema#"
)

// Schema represents a JSON Schema
type Schema struct {
	*definition.Definition
	Version     string                            `json:"$schema,omitempty"`
	Definitions map[string]*definition.Definition `json:"definitions,omitempty"`
}

func defaultFormatRefName(pkg, name string) string {
	return strings.ReplaceAll(pkg, "/", "|") + "." + name
}

// GenerateSchema is the entrypoint for schema generation
func GenerateSchema(
	pkgPath string,
	rootRef string,
	tagNamespace string,
	formatRefName definition.RefNameFormatFunc,
	strict bool,
) (Schema, error) {
	if len(tagNamespace) == 0 {
		tagNamespace = "json"
	}

	if formatRefName == nil {
		formatRefName = defaultFormatRefName
	}

	definitions := make(map[string]*definition.Definition)

	importer, err := importer.NewImporter(pkgPath)
	if err != nil {
		return Schema{}, err
	}
	dg := definition.Generator{
		Strict:        strict,
		Definitions:   definitions,
		Importer:      importer,
		TagNamespace:  tagNamespace,
		FromatRefName: formatRefName,
	}

	dg.CollectDefinitionsFromStruct(pkgPath, rootRef)

	s := Schema{
		Version: version7,
		Definition: &definition.Definition{
			Type: "object",
			Ref:  definition.DefPrefix + rootRef,
		},
		Definitions: dg.Definitions,
	}
	if _, ok := dg.Definitions[rootRef]; !ok {
		return s, fmt.Errorf("Couldn't find ref %s in definitions", rootRef)
	}
	return s, nil
}

// ToJSON serializes and makes sure HTML description are not escaped
func ToJSON(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(v); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
