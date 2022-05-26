package cosign

import (
	"fmt"
	"strings"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/templateutils"
	"arhat.dev/dukkha/pkg/tools"
	"arhat.dev/dukkha/pkg/tools/buildah"
)

const TaskKindUpload = "upload"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindUpload,
		func(toolName string) dukkha.Task {
			t := &TaskUpload{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), t)
			return t
		},
	)
}

type TaskUpload struct {
	rs.BaseField `yaml:"-"`

	TaskName string `yaml:"name"`

	tools.BaseTask `yaml:",inline"`

	// Kind is either blob or wasm
	//
	// Defaults to `"blob"`
	UploadKind string `yaml:"kind"`

	// Files to upload at one batch
	Files []FileSpec `yaml:"files"`

	// Signing sign uploaded images
	Signing signingSpec `yaml:"signing"`

	// ImageNames
	ImageNames []buildah.ImageNameSpec `yaml:"image_names"`
}

type FileSpec struct {
	rs.BaseField `yaml:"-"`

	// Path to local blob/wasm file
	Path string `yaml:"path"`

	// ContentType of the local file
	ContentType string `yaml:"content_type"`
}

type signingSpec struct {
	rs.BaseField `yaml:"-"`

	// Enable signing
	Enabled bool `yaml:"enabled"`

	Options imageSigningOptions `yaml:",inline"`
}

func (c *TaskUpload) Kind() dukkha.TaskKind { return TaskKindUpload }
func (c *TaskUpload) Name() dukkha.TaskName { return dukkha.TaskName(c.TaskName) }
func (c *TaskUpload) Key() dukkha.TaskKey {
	return dukkha.TaskKey{Kind: c.Kind(), Name: c.Name()}
}

func (c *TaskUpload) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var steps []dukkha.TaskExecSpec
	err := c.DoAfterFieldsResolved(rc, -1, true, func() error {
		var keyFile string
		if c.Signing.Enabled {
			var err error
			keyFile, err = c.Signing.Options.Options.ensurePrivateKey(c.CacheFS)
			if err != nil {
				return fmt.Errorf("ensuring private key: %w", err)
			}
		}

		// cosign
		kind := c.UploadKind
		if len(kind) == 0 {
			kind = "blob"
		}

		ociOS, ok := constant.GetOciOS(rc.MatrixKernel())
		if !ok {
			ociOS = rc.MatrixKernel()
		}

		ociArch, ok := constant.GetOciArch(rc.MatrixArch())
		if !ok {
			ociArch = string(rc.MatrixArch())
		}

		ociVariant, _ := constant.GetOciArchVariant(rc.MatrixArch())

		var ociPlatformParts []string
		if len(ociOS) != 0 {
			ociPlatformParts = append(ociPlatformParts, ociOS)
		}

		if len(ociArch) != 0 {
			ociPlatformParts = append(ociPlatformParts, ociArch)
		}

		if len(ociVariant) != 0 {
			ociPlatformParts = append(ociPlatformParts, ociVariant)
		}

		ociPlatform := strings.Join(ociPlatformParts, "/")

		uploadCmd := []string{constant.DUKKHA_TOOL_CMD, "upload", kind}
		for _, fSpec := range c.Files {
			path := fSpec.Path

			if kind != "blob" {
				uploadCmd = append(uploadCmd, "--files", path)
				continue
			}

			if len(ociPlatform) != 0 {
				path += ":" + ociPlatform
			}

			uploadCmd = append(uploadCmd, "--files", path)

			if len(fSpec.ContentType) != 0 {
				uploadCmd = append(uploadCmd, "--ct", fSpec.ContentType)
			}
		}

		for _, spec := range c.ImageNames {
			if len(spec.Image) == 0 {
				continue
			}

			imageName := templateutils.GetFullImageName_UseDefault_IfIfNoTagSet(
				rc, spec.Image, true,
			)

			steps = append(steps, dukkha.TaskExecSpec{
				Command: sliceutils.NewStrings(uploadCmd, imageName),
			})

			if c.Signing.Enabled {
				steps = append(steps,
					c.Signing.Options.genSignAndVerifySpec(
						keyFile,
						imageName,
					)...,
				)
			}
		}

		return nil
	})

	return steps, err
}
