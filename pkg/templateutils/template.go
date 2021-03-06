package templateutils

import (
	"arhat.dev/tlang"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/third_party/golang/text/template"
)

func CreateTextTemplate(rc dukkha.RenderingContext) *template.Template {
	tfs := CreateTemplateFuncs(rc)
	return template.New("tmpl").Funcs(&tfs)
}

func CreateTLangTemplate(rc dukkha.RenderingContext) *tlang.Template {
	tfs := CreateTemplateFuncs(rc)
	return tlang.New("tlang").Funcs(&tfs)
}
