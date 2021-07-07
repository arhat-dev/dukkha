package buildah

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
	"arhat.dev/dukkha/pkg/types"
)

const TaskKindPush = "push"

func init() {
	field.RegisterInterfaceField(
		dukkha.TaskType,
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

var _ dukkha.Task = (*TaskPush)(nil)

type TaskPush struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	ImageNames []ImageNameSpec `yaml:"image_names"`
}

func (c *TaskPush) ToolKind() dukkha.ToolKind { return ToolKind }
func (c *TaskPush) Kind() dukkha.TaskKind     { return TaskKindPush }

func (c *TaskPush) GetExecSpecs(rc types.RenderingContext, buildahCmd []string) ([]dukkha.TaskExecSpec, error) {
	targets := c.ImageNames
	if len(targets) == 0 {
		targets = []ImageNameSpec{
			{
				Image:    c.TaskName,
				Manifest: "",
			},
		}
	}

	dukkhaCacheDir := rc.CacheDir()

	var result []dukkha.TaskExecSpec
	for _, spec := range targets {
		if len(spec.Image) != 0 {
			imageName := SetDefaultImageTagIfNoTagSet(rc, spec.Image)
			imageIDFile := GetImageIDFileForImageName(
				dukkhaCacheDir, imageName,
			)
			imageIDBytes, err := os.ReadFile(imageIDFile)
			if err != nil {
				return nil, fmt.Errorf("image id file not found: %w", err)
			}

			result = append(result, dukkha.TaskExecSpec{
				Command: sliceutils.NewStrings(
					buildahCmd, "push",
					string(bytes.TrimSpace(imageIDBytes)),
					// TODO: support other destination
					"docker://"+imageName,
				),
				IgnoreError: false,
			})
		}

		if len(spec.Manifest) == 0 {
			continue
		}

		// buildah manifest push --all \
		//   <manifest-list-name> <transport>:<transport-details>
		manifestName := SetDefaultManifestTagIfNoTagSet(rc, spec.Manifest)
		result = append(result, dukkha.TaskExecSpec{
			Command: sliceutils.NewStrings(
				buildahCmd, "manifest", "push", "--all",
				getLocalManifestName(manifestName),
				// TODO: support other destination
				"docker://"+manifestName,
			),
			IgnoreError: false,
		})
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
