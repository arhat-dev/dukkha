package helm

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindIndex = "index"

func init() {
	field.RegisterInterfaceField(
		tools.TaskType,
		regexp.MustCompile(`^helm(:.+){0,1}:index$`),
		func(subMatches []string) interface{} {
			t := &TaskIndex{}
			if len(subMatches) != 0 {
				t.SetToolName(strings.TrimPrefix(subMatches[0], ":"))
			}
			return t
		},
	)
}

var _ tools.Task = (*TaskIndex)(nil)

type TaskIndex struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	RepoURL     string `yaml:"repo_url"`
	PackagesDir string `yaml:"packages_dir"`
	Merge       string `yaml:"merge"`
}

func (c *TaskIndex) ToolKind() string { return ToolKind }
func (c *TaskIndex) TaskKind() string { return TaskKindIndex }

func (c *TaskIndex) GetExecSpecs(ctx *field.RenderingContext, helmCmd []string) ([]tools.TaskExecSpec, error) {
	indexCmd := sliceutils.NewStrings(helmCmd, "repo", "index")

	if len(c.RepoURL) != 0 {
		indexCmd = append(indexCmd, "--url", c.RepoURL)
	}

	dukkhaWorkingDir := ctx.Values().Env[constant.ENV_DUKKHA_WORKING_DIR]
	if len(c.PackagesDir) != 0 {
		pkgDir, err2 := filepath.Abs(c.PackagesDir)
		if err2 != nil {
			return nil, fmt.Errorf("failed to determine absolute path of package_dir: %w", err2)
		}

		indexCmd = append(indexCmd, pkgDir)
	} else {
		indexCmd = append(indexCmd, dukkhaWorkingDir)
	}

	if len(c.Merge) != 0 {
		indexCmd = append(indexCmd, "--merge", c.Merge)
	}

	return []tools.TaskExecSpec{
		{
			Command: indexCmd,
		},
	}, nil
}
