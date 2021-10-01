package buildah

import (
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/rs"
)

type stepSet struct {
	rs.BaseField

	// Commit every following step as new layer
	Commit *bool `yaml:"commit"`

	// Workdir
	Workdir *string `yaml:"workdir"`

	// User for command running in run step
	User *string `yaml:"user"`

	// Shell command to interpreter scripts in run step
	Shell []string `yaml:"shell"`
}

func (s *stepSet) genCtx(stepCtx *xbuildContext) *xbuildContext {
	ctx := stepCtx.clone()

	if s.Commit != nil {
		ctx.Commit = *s.Commit
	}

	if s.Workdir != nil {
		ctx.WorkDir = *s.Workdir
	}

	if s.User != nil {
		ctx.User = *s.User
	}

	if len(s.Shell) != 0 {
		ctx.Shell = sliceutils.NewStrings(s.Shell)
	}

	return ctx
}
