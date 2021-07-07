package helm

import (
	"path/filepath"
	"strings"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
	"arhat.dev/dukkha/pkg/types"
)

const TaskKindPackage = "package"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindPackage,
		func(toolName string) dukkha.Task {
			t := &TaskPackage{}
			t.SetToolName(toolName)
			return t
		},
	)
}

var _ dukkha.Task = (*TaskPackage)(nil)

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

func (c *TaskPackage) ToolKind() dukkha.ToolKind { return ToolKind }
func (c *TaskPackage) Kind() dukkha.TaskKind     { return TaskKindPackage }

func (c *TaskPackage) GetExecSpecs(rc types.RenderingContext, helmCmd []string) ([]dukkha.TaskExecSpec, error) {
	pkgStep := &dukkha.TaskExecSpec{
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

	return []dukkha.TaskExecSpec{*pkgStep}, nil
}
