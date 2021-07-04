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
}

func (c *TaskPush) ToolKind() string { return ToolKind }
func (c *TaskPush) TaskKind() string { return TaskKindPush }

func (c *TaskPush) GetExecSpecs(ctx *field.RenderingContext, buildahCmd []string) ([]tools.TaskExecSpec, error) {
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
		if ImageOrManifestHasFQDN(spec.Image) {
			result = append(result, tools.TaskExecSpec{
				Command: sliceutils.NewStrings(
					buildahCmd, "push",
					SetDefaultImageTagIfNoTagSet(ctx, spec.Image),
				),
				IgnoreError: false,
			})
		}

		if ImageOrManifestHasFQDN(spec.Manifest) {
			// buildah manifest push --all \
			//   <manifest-list-name> <transport>:<transport-details>
			manifestName := SetDefaultManifestTagIfNoTagSet(ctx, spec.Manifest)
			result = append(result, tools.TaskExecSpec{
				Command: sliceutils.NewStrings(
					buildahCmd, "manifest", "push", "--all",
					getLocalManifestName(manifestName),
					"docker://"+manifestName,
				),
				IgnoreError: false,
			})
		}
	}

	return result, nil
}

func ImageOrManifestHasFQDN(s string) bool {
	parts := strings.SplitN(s, "/", 2)
	if len(parts) == 1 {
		return false
	}

	return strings.Contains(parts[0], ".")
}
