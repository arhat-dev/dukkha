package helm

import (
	"path/filepath"
	"regexp"
	"strings"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindPackage = "package"

func init() {
	field.RegisterInterfaceField(
		tools.TaskType,
		regexp.MustCompile(`^helm(:.+){0,1}:package$`),
		func(subMatches []string) interface{} {
			t := &TaskPackage{}
			if len(subMatches) != 0 {
				t.SetToolName(strings.TrimPrefix(subMatches[0], ":"))
			}
			return t
		},
	)
}

var _ tools.Task = (*TaskPackage)(nil)

type TaskPackage struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	Chart       string `yaml:"chart"`
	PackagesDir string `yaml:"packages_dir"`

	Signing PackageSigningSpec `yaml:"signing"`
}

type PackageSigningSpec struct {
	Enabled          bool   `yaml:"enabled"`
	GPGKeyring       string `yaml:"gpg_keyring"`
	GPGKeyName       string `yaml:"gpg_key_name"`
	GPGKeyPassphrase string `yaml:"gpg_key_passphrase"`
}

func (c *TaskPackage) ToolKind() string { return ToolKind }
func (c *TaskPackage) TaskKind() string { return TaskKindPackage }

func (c *TaskPackage) GetExecSpecs(ctx *field.RenderingContext, helmCmd []string) ([]tools.TaskExecSpec, error) {
	pkgStep := &tools.TaskExecSpec{
		Command: sliceutils.NewStrings(helmCmd, "package"),
	}

	matches, err := filepath.Glob(c.Chart)
	if err != nil {
		pkgStep.Command = append(pkgStep.Command, c.Chart)
	} else {
		pkgStep.Command = append(pkgStep.Command, matches...)
	}

	if len(c.PackagesDir) != 0 {
		pkgStep.Command = append(pkgStep.Command,
			"--destination", c.PackagesDir,
		)
	}

	if c.Signing.Enabled {
		pkgStep.Command = append(pkgStep.Command, "--sign",
			"--key", c.Signing.GPGKeyName,
		)

		if len(c.Signing.GPGKeyring) != 0 {
			pkgStep.Command = append(pkgStep.Command,
				"--keyring", c.Signing.GPGKeyring,
			)
		}

		if len(c.Signing.GPGKeyPassphrase) != 0 {
			pkgStep.Command = append(pkgStep.Command, "--passphrase-file", "-")

			input := c.Signing.GPGKeyPassphrase
			pkgStep.Stdin = strings.NewReader(input)
		}
	}

	return []tools.TaskExecSpec{*pkgStep}, nil
}
