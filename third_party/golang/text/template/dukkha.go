package template

import "arhat.dev/dukkha/third_party/golang/text/template/parse"

// GetExecFuncs returns all template funcs accessible when executing
// the template
func (t *Template) GetExecFuncs() parse.TemplateFuncs {
	return t.funcs
}
