package templateutils

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/third_party/golang/text/template"
)

func CreateTemplate(rc dukkha.RenderingContext) *template.Template {
	tfs := CreateTemplateFuncs(rc)
	return template.New("tpl").Funcs(&tfs)
}
