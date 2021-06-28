package buildah

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
	"arhat.dev/pkg/textquery"
)

const TaskKindBud = "bud"

func init() {
	field.RegisterInterfaceField(
		tools.TaskType,
		regexp.MustCompile(`^buildah(:.+){0,1}:bud$`),
		func(params []string) interface{} {
			t := &TaskBud{}
			if len(params) != 0 {
				t.SetToolName(strings.TrimPrefix(params[0], ":"))
			}
			return t
		},
	)
}

var _ tools.Task = (*TaskBud)(nil)

type TaskBud struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	Context    string          `yaml:"context"`
	ImageNames []ImageNameSpec `yaml:"image_names"`
	Dockerfile string          `yaml:"dockerfile"`
	ExtraArgs  []string        `yaml:"extraArgs"`
}

type ImageNameSpec struct {
	Image    string `yaml:"image"`
	Manifest string `yaml:"manifest"`
}

func (c *TaskBud) ToolKind() string { return ToolKind }
func (c *TaskBud) TaskKind() string { return TaskKindBud }

func (c *TaskBud) GetExecSpecs(ctx *field.RenderingContext, toolCmd []string) ([]tools.TaskExecSpec, error) {
	var steps []tools.TaskExecSpec

	// create an image id file
	imageIDFile, err := ioutil.TempFile(
		ctx.Values().Env[constant.ENV_DUKKHA_CACHE_DIR], "buildah-bud-image-id-*",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create a temp file for image id: %w", err)
	}
	imageIDFilePath := imageIDFile.Name()
	_ = imageIDFile.Close()

	budCmd := sliceutils.NewStringSlice(toolCmd, "bud", "--iidfile", imageIDFilePath)
	if len(c.Dockerfile) != 0 {
		budCmd = append(budCmd, "-f", c.Dockerfile)
	}

	budCmd = append(budCmd, c.ExtraArgs...)

	targets := c.ImageNames
	if len(targets) == 0 {
		targets = []ImageNameSpec{{
			Image:    c.Name,
			Manifest: "",
		}}
	}

	// set image names
	for _, spec := range targets {
		if len(spec.Image) == 0 {
			continue
		}

		budCmd = append(budCmd, "-t", spec.Image)
	}

	context := c.Context
	if len(context) == 0 {
		context = "."
	}

	// buildah bud
	steps = append(steps, tools.TaskExecSpec{
		Command:     append(budCmd, context),
		IgnoreError: false,
	})

	// read image id file to get image id
	const replaceTargetImageID = "<IMAGE_ID>"
	steps = append(steps, tools.TaskExecSpec{
		OutputAsReplace:     replaceTargetImageID,
		FixOutputForReplace: strings.TrimSpace,

		AlterExecFunc: func(
			replace map[string]string,
			stdin io.Reader, stdout, stderr io.Writer,
		) ([]tools.TaskExecSpec, error) {
			imageIDBytes, err := os.ReadFile(imageIDFilePath)
			if err != nil {
				return nil, err
			}

			_, err = stdout.Write(imageIDBytes)
			return nil, err
		},
		IgnoreError: false,
	})

	// buildah inspect --type image to get image digest from image id
	const replaceTargetImageDigest = "<IMAGE_DIGEST>"
	steps = append(steps, tools.TaskExecSpec{
		OutputAsReplace:     replaceTargetImageDigest,
		FixOutputForReplace: strings.TrimSpace,

		Command: sliceutils.NewStringSlice(
			toolCmd, "inspect", "--type", "image",
			"--format", `"{{ .FromImageDigest }}"`,
			replaceTargetImageID,
		),
		IgnoreError: false,
	})

	// add to manifest, ensure same os/arch/variant only one exist
	mArch := ctx.Values().Env[constant.ENV_MATRIX_ARCH]
	arch := constant.GetOciArch(mArch)
	os := constant.GetOciOS(ctx.Values().Env[constant.ENV_MATRIX_KERNEL])
	variant := constant.GetOciArchVariant(mArch)

	manifestOsArchVariantQueryForDigest := fmt.Sprintf(
		`.manifests[] | select( .platform.os == "%s")`+
			` | select (.platform.architecture == "%s")`,
		os, arch,
	)
	osArchVariantArgs := []string{
		"--os", os, "--arch", arch,
	}
	if len(variant) != 0 {
		manifestOsArchVariantQueryForDigest += fmt.Sprintf(
			` | select (.platform.vairant == "%s")`, variant,
		)
		osArchVariantArgs = append(osArchVariantArgs, "--variant", variant)
	}
	manifestOsArchVariantQueryForDigest += ` | .digest`

	for _, spec := range targets {
		if len(spec.Manifest) == 0 {
			continue
		}

		const replaceTargetManifestSpec = "<MANIFEST_SPEC>"
		steps = append(steps, tools.TaskExecSpec{
			OutputAsReplace:     replaceTargetManifestSpec,
			FixOutputForReplace: nil,

			Command: sliceutils.NewStringSlice(
				toolCmd, "manifest", "inspect", spec.Manifest,
			),
			// manifest may not exist
			IgnoreError: true,
		})

		manifestAddCmd := sliceutils.NewStringSlice(toolCmd, "manifest", "add")
		manifestAddCmd = append(manifestAddCmd, osArchVariantArgs...)

		// find existing manifest entries with same os/arch/variant
		steps = append(steps, tools.TaskExecSpec{
			AlterExecFunc: func(
				replace map[string]string,
				stdin io.Reader, stdout, stderr io.Writer,
			) ([]tools.TaskExecSpec, error) {
				var subSteps []tools.TaskExecSpec

				manifestSpec, ok := replace[replaceTargetManifestSpec]
				if !ok {
					// manifest not created, create and add this image
					subSteps = append(subSteps)
					return []tools.TaskExecSpec{
						{
							Command: sliceutils.NewStringSlice(
								toolCmd, "manifest", "create", spec.Manifest,
							),
							IgnoreError: false,
						},
						{
							Command: sliceutils.NewStringSlice(
								manifestAddCmd, spec.Manifest, replaceTargetImageID,
							),
							IgnoreError: false,
						},
					}, nil
				}

				// manifest already created
				digestLines, err := textquery.JQ(manifestOsArchVariantQueryForDigest, manifestSpec)
				if err != nil {
					return nil, fmt.Errorf("failed to lookup entries in manifest spec: %w", err)
				}

				// remove existing entries with same os/arch/variant
				for _, digest := range strings.Split(digestLines, "\n") {
					digest = strings.TrimSpace(digest)
					if len(digest) == 0 {
						continue
					}

					subSteps = append(subSteps, tools.TaskExecSpec{
						Command:     sliceutils.NewStringSlice(toolCmd, "manifest", "remove", spec.Manifest, digest),
						IgnoreError: false,
					})
				}

				// add this image to manifest with correct os/arch/variant
				subSteps = append(subSteps, tools.TaskExecSpec{
					Command:     sliceutils.NewStringSlice(manifestAddCmd, spec.Manifest, replaceTargetImageID),
					IgnoreError: false,
				})

				return subSteps, nil
			},
			IgnoreError: false,
		})
	}

	return steps, nil
}
