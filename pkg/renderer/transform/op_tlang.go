package transform

import (
	"arhat.dev/pkg/stringhelper"

	"arhat.dev/dukkha/pkg/renderer/tlang"
	"arhat.dev/dukkha/pkg/templateutils"
)

type tlangSpec tlang.InputSpec

func (s *tlangSpec) Run(rc extendedUserFacingRenderContext) (ret string, err error) {
	tfs := templateutils.CreateTemplateFuncs(rc)
	retBytes, err := tlang.RenderTlang(rc, &tfs, s.Config.Include, s.Config.Variables.NormalizedValue(), s.Script)
	ret = stringhelper.Convert[string, byte](retBytes)
	return
}
