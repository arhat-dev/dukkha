package buildah

import (
	"fmt"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
)

// step is structured `buildah <subcmd>` for image building
type step struct {
	rs.BaseField `yaml:"-"`

	// ID of this step, if not set, will be the array index of this step
	ID string `yaml:"id,omitempty"`

	// Record to add flag --add-history
	Record *bool `yaml:"record,omitempty"`

	// Commit this step as a new layer after this step finished
	//
	// this is set to true by default when:
	// - at last step
	// - switching to different container at next step (next step is a FROM statement)
	Commit *bool `yaml:"commit,omitempty"`

	// CommitAs set the image name the container committed as
	CommitAs        string   `yaml:"commit_as,omitempty"`
	ExtraCommitArgs []string `yaml:"extra_commit_args,omitempty"`

	// Compress when commit, defaults to true
	Compress *bool `yaml:"compress,omitempty"`

	// Skip this step when set to true
	Skip bool `yaml:"skip"`

	//
	// Step spec
	//

	// Set default options for all following steps
	Set *stepSet `yaml:"set,omitempty"`

	// From some rootfs
	From *stepFrom `yaml:"from,omitempty"`

	// Run some command
	Run *stepRun `yaml:"run,omitempty"`

	// Copy files to somewhere
	Copy *stepCopy `yaml:"copy,omitempty"`
}

func (s *step) genSpec(
	rc dukkha.TaskExecContext,
	cacheFS *fshelper.OSFS,
) ([]dukkha.TaskExecSpec, error) {
	record := true
	if s.Record != nil {
		record = *s.Record
	}

	switch {
	case s.Set != nil:
		return s.Set.genSpec(rc, record)
	case s.From != nil:
		return s.From.genSpec(rc)
	case s.Run != nil:
		return s.Run.genSpec(rc, cacheFS, record)
	case s.Copy != nil:
		return s.Copy.genSpec(rc, cacheFS, record)
	default:
		return nil, fmt.Errorf("unknown step")
	}
}
