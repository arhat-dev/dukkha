package buildah

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/rs"
)

type stepCopy struct {
	rs.BaseField

	From copyFromSpec `yaml:"from"`
	To   copyToSpec   `yaml:"to"`
}

func (s *stepCopy) genSpec(
	rc dukkha.TaskExecContext,
	options dukkha.TaskMatrixExecOptions,
	stepCtx *xbuildContext,
) ([]dukkha.TaskExecSpec, error) {
	return nil, nil
}

type copyFromSpec struct {
	rs.BaseField

	Local *copyFromLocalSpec `yaml:"local"`
	HTTP  *copyFromHTTPSpec  `yaml:"http"`
	Image *copyFromImageSpec `yaml:"image"`
}

type copyFromLocalSpec struct {
	rs.BaseField

	Path string `yaml:"path"`
}

func (s *copyFromLocalSpec) genSpec(
	rc dukkha.TaskExecContext,
	options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	sliceutils.NewStrings(options.ToolCmd(), "add")
	return nil, nil
}

type copyFromHTTPSpec struct {
	rs.BaseField

	URL string `yaml:"url"`
}

type copyFromImageSpec struct {
	rs.BaseField

	Name   *string `yaml:"name"`
	Digest *string `yaml:"digest"`

	Path string `yaml:"path"`
}

type copyToSpec struct {
	rs.BaseField

	Path string `yaml:"path"`

	Chmod []chmodSpec `yaml:"chmod"`
	Chown []chownSpec `yaml:"chown"`
}

type chmodSpec struct {
	rs.BaseField

	// Match glob pattern to match files
	Match string `yaml:"match"`

	// Ignore glob pattern to ignore files
	Ignore string `yaml:"ignore"`

	// Value for chmod (e.g. a+x, 0755)
	Value string `yaml:"value"`

	// Recursive run chmod on matched files
	Recursive bool `yaml:"recursive"`
}

type chownSpec struct {
	rs.BaseField

	// Match glob pattern to match files
	Match string `yaml:"match"`

	// Ignore glob pattern to ignore files
	Ignore string `yaml:"ignore"`

	// Value for chown (e.g. user:group, user, uid, uid:gid)
	Value string `yaml:"value"`

	// Recursive run chown on matched files
	Recursive bool `yaml:"recursive"`
}
