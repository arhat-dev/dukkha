package tlang

import "arhat.dev/dukkha/pkg/renderer/tl/tlang/parse"

// GetExecFuncs returns all template funcs accessible when executing
// the template
func (t *Template) GetExecFuncs() parse.TemplateFuncs {
	return t.funcs
}
