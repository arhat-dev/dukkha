package rshelper

import (
	"bytes"
	"fmt"
	"text/template"

	"arhat.dev/rs"

	"arhat.dev/pkg/yamlhelper"
)

var _ rs.RenderingHandler = (*TemplateHandler)(nil)

// TemplateHandler execute raw data as text/template
type TemplateHandler struct {
	CreateFuncMap func() template.FuncMap
}

func (h *TemplateHandler) RenderYaml(_ string, rawData interface{}) ([]byte, error) {
	tplBytes, err := yamlhelper.ToYamlBytes(rawData)
	if err != nil {
		return nil, fmt.Errorf("failed to get data bytes of input: %w", err)
	}

	t := template.New("")
	if h.CreateFuncMap != nil {
		t = t.Funcs(h.CreateFuncMap())
	}

	t, err = t.Parse(string(tplBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %q: %w", string(tplBytes), err)
	}

	buf := &bytes.Buffer{}
	err = t.Execute(buf, nil)

	return buf.Bytes(), err
}
