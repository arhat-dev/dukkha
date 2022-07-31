package transform

import (
	"arhat.dev/pkg/stringhelper"

	"arhat.dev/dukkha/pkg/renderer/tmpl"
	"arhat.dev/dukkha/pkg/templateutils"
)

type tmplSpec tmpl.InputSpec

func (s *tmplSpec) Run(rc extendedUserFacingRenderContext) (ret string, err error) {
	tfs := templateutils.CreateTemplateFuncs(rc)
	retBytes, err := tmpl.RenderTemplate(rc, &tfs, s.Config.Include, s.Config.Variables.NormalizedValue(), s.Template)
	ret = stringhelper.Convert[string, byte](retBytes)
	return
}
