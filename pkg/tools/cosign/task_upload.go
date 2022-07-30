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
	dukkha.RegisterTask(ToolKind, TaskKindUpload, tools.NewTask[TaskUpload, *TaskUpload])
}

type TaskUpload struct {
	tools.BaseTask[CosignUpload, *CosignUpload]
}

type CosignUpload struct {
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

	parent tools.BaseTaskType
}

func (w *CosignUpload) ToolKind() dukkha.ToolKind       { return ToolKind }
func (w *CosignUpload) Kind() dukkha.TaskKind           { return TaskKindUpload }
func (w *CosignUpload) LinkParent(p tools.BaseTaskType) { w.parent = p }

func (c *CosignUpload) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var steps []dukkha.TaskExecSpec
	err := c.parent.DoAfterFieldsResolved(rc, -1, true, func() error {
		var keyFile string
		if c.Signing.Enabled {
			var err error
			keyFile, err = c.Signing.Options.Options.ensurePrivateKey(c.parent.CacheFS())
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
			ociArch = rc.MatrixArch()
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
