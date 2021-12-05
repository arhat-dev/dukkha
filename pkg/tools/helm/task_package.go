package helm

import (
	"strings"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindPackage = "package"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindPackage,
		func(toolName string) dukkha.Task {
			t := &TaskPackage{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), t)
			return t
		},
	)
}

type TaskPackage struct {
	rs.BaseField `yaml:"-"`

	TaskName string `yaml:"name"`

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

func (c *TaskPackage) Kind() dukkha.TaskKind { return TaskKindPackage }
func (c *TaskPackage) Name() dukkha.TaskName { return dukkha.TaskName(c.TaskName) }
func (c *TaskPackage) Key() dukkha.TaskKey {
	return dukkha.TaskKey{Kind: c.Kind(), Name: c.Name()}
}

func (c *TaskPackage) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	pkgStep := &dukkha.TaskExecSpec{
		Command: []string{constant.DUKKHA_TOOL_CMD, "package"},
	}

	err := c.DoAfterFieldsResolved(rc, -1, true, func() error {
		matches, err := rc.FS().Glob(c.Chart)
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
