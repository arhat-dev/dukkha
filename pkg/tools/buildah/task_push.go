package buildah

import (
	"regexp"
	"strings"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindPush = "push"

func init() {
	field.RegisterInterfaceField(
		tools.TaskType,
		regexp.MustCompile(`^buildah(:.+){0,1}:push$`),
		func(params []string) interface{} {
			t := &TaskPush{}
			if len(params) != 0 {
				t.SetToolName(strings.TrimPrefix(params[0], ":"))
			}
			return t
		},
	)
}

var _ tools.Task = (*TaskPush)(nil)

type TaskPush struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	ImageNames []ImageNameSpec `yaml:"image_names"`
	ExtraArgs  []string        `yaml:"extraArgs"`
}

func (c *TaskPush) ToolKind() string { return ToolKind }
func (c *TaskPush) TaskKind() string { return TaskKindPush }

func (c *TaskPush) GetExecSpecs(ctx *field.RenderingContext, toolCmd []string) ([]tools.TaskExecSpec, error) {
	targets := c.ImageNames
	if len(targets) == 0 {
		targets = []ImageNameSpec{
			{
				Image:    c.Name,
				Manifest: "",
			},
		}
	}

	var result []tools.TaskExecSpec
	for _, spec := range targets {
		if len(spec.Image) == 0 {
			continue
		}

		result = append(result, tools.TaskExecSpec{
			Command:     sliceutils.NewStringSlice(toolCmd, "push", spec.Image),
			IgnoreError: false,
		})

		if len(spec.Manifest) == 0 {
			continue
		}

		manifestCmd := sliceutils.NewStringSlice(toolCmd, "manifest")

		// buildah manifest push <manifest-list-name> {<transport>:<transport-details>}
		result = append(result, tools.TaskExecSpec{
			Command:     sliceutils.NewStringSlice(manifestCmd, "push", spec.Manifest),
			IgnoreError: false,
		})
	}

	return result, nil
}

func (c *TaskPush) getManifestOS(os string) string {
	return strings.ToLower(os)
}
