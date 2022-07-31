package transform

import (
	"fmt"
	"strings"

	"arhat.dev/rs"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

	"arhat.dev/dukkha/pkg/templateutils"
)

type shellSpec struct {
	rs.BaseField

	// Script is the bash script to run
	Script string `yaml:"script"`

	// TODO: add other options and update shell renderer input spec
}

func (s *shellSpec) Run(rc extendedUserFacingRenderContext) (ret string, err error) {
	var (
		sb     strings.Builder
		runner *interp.Runner
	)
	runner, err = templateutils.CreateShellRunner(
		rc.WorkDir(), rc, nil, &sb, rc.Stderr(),
	)
	if err != nil {
		err = fmt.Errorf("create embedded shell: %w", err)
		return
	}

	parser := syntax.NewParser(
		syntax.Variant(syntax.LangBash),
	)

	err = templateutils.RunScript(rc, runner, parser, s.Script)
	if err != nil {
		err = fmt.Errorf("run script: %w", err)
		return
	}

	return sb.String(), nil
}
