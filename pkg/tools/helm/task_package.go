package helm

import (
	"path/filepath"
	"strings"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindPackage = "package"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindPackage,
		func(toolName string) dukkha.Task {
			t := &TaskPackage{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindPackage, t)
			return t
		},
	)
}

type TaskPackage struct {
	rs.BaseField

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

func (c *TaskPackage) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	pkgStep := &dukkha.TaskExecSpec{
		Command: sliceutils.NewStrings(options.ToolCmd(), "package"),

		UseShell:  options.UseShell(),
		ShellName: options.ShellName(),
	}

	err := c.DoAfterFieldsResolved(rc, -1, func() error {
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

				pkgStep.Stdin = strings.NewReader(c.Signing.GPGKeyPassphrase)
			}
		}

		return nil
	})

	return []dukkha.TaskExecSpec{*pkgStep}, err
}
