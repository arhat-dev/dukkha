package buildah

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"arhat.dev/pkg/hashhelper"
	"arhat.dev/pkg/textquery"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
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
	ExtraArgs  []string        `yaml:"extra_args"`
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
	dukkhaCacheDir := ctx.Values().Env[constant.ENV_DUKKHA_CACHE_DIR]
	imageIDFile, err := ioutil.TempFile(dukkhaCacheDir, "buildah-bud-image-id-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create a temp file for image id: %w", err)
	}
	imageIDFilePath := imageIDFile.Name()
	_ = imageIDFile.Close()

	budCmd := sliceutils.NewStrings(toolCmd, "bud", "--iidfile", imageIDFilePath)
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
	var localImageIDFiles []string
	for _, spec := range targets {
		if len(spec.Image) == 0 {
			continue
		}

		imageName := SetDefaultImageTagIfNoTagSet(ctx, spec.Image)

		// local image name is to handle bud regression bugs related to
		// FQDN image names
		budCmd = append(budCmd, "-t", imageName)

		filePath := getImageIDFilePathForImageName(
			dukkhaCacheDir, imageName,
		)
		err = os.MkdirAll(filepath.Dir(filePath), 0750)
		if err != nil && !os.IsExist(err) {
			return nil, fmt.Errorf("failed to ensure image id dir exists")
		}

		localImageIDFiles = append(localImageIDFiles, filePath)
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
		FixOutputForReplace: bytes.TrimSpace,

		AlterExecFunc: func(
			replace map[string][]byte,
			stdin io.Reader, stdout, stderr io.Writer,
		) ([]tools.TaskExecSpec, error) {
			imageIDBytes, err := os.ReadFile(imageIDFilePath)
			if err != nil {
				return nil, err
			}

			for _, f := range localImageIDFiles {
				err = os.WriteFile(f, imageIDBytes, 0750)
				if err != nil {
					return nil, err
				}
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
		FixOutputForReplace: bytes.TrimSpace,

		Command: sliceutils.NewStrings(
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

		manifestName := SetDefaultManifestTagIfNoTagSet(ctx, spec.Manifest)
		localManifestName := getLocalManifestName(manifestName)
		// ensure local manifest exists
		steps = append(steps, tools.TaskExecSpec{
			Command: sliceutils.NewStrings(
				toolCmd, "manifest", "create", localManifestName,
			),
			IgnoreError: true,
		})

		const replaceTargetManifestSpec = "<MANIFEST_SPEC>"
		steps = append(steps, tools.TaskExecSpec{
			OutputAsReplace:     replaceTargetManifestSpec,
			FixOutputForReplace: nil,

			Command: sliceutils.NewStrings(
				toolCmd, "manifest", "inspect", localManifestName,
			),
			// manifest may not exist
			IgnoreError: true,
		})

		manifestAddCmd := sliceutils.NewStrings(toolCmd, "manifest", "add")
		manifestAddCmd = append(manifestAddCmd, osArchVariantArgs...)
		manifestAddCmd = append(manifestAddCmd, localManifestName, replaceTargetImageID)

		// find existing manifest entries with same os/arch/variant
		steps = append(steps, tools.TaskExecSpec{
			IgnoreError: false,
			AlterExecFunc: func(
				replace map[string][]byte,
				stdin io.Reader, stdout, stderr io.Writer,
			) ([]tools.TaskExecSpec, error) {
				manifestSpec, ok := replace[replaceTargetManifestSpec]
				if !ok {
					// manifest not created, usually should not happen since we just created before
					return []tools.TaskExecSpec{
						{
							// do not ignore manifest create error this time
							Command: sliceutils.NewStrings(
								toolCmd, "manifest", "create", localManifestName,
							),
							IgnoreError: false,
						},
						{
							Command:     sliceutils.NewStrings(manifestAddCmd),
							IgnoreError: false,
						},
					}, nil
				}

				// manifest already created, query to get all matching digests
				digestLines, err := textquery.JQBytes(manifestOsArchVariantQueryForDigest, manifestSpec)
				if err != nil {
					// no manifests entries, add this image directly
					return []tools.TaskExecSpec{{
						Command:     sliceutils.NewStrings(manifestAddCmd),
						IgnoreError: false,
					}}, nil
				}

				var subSteps []tools.TaskExecSpec

				// remove existing entries with same os/arch/variant
				for _, digest := range strings.Split(digestLines, "\n") {
					digest = strings.TrimSpace(digest)
					if len(digest) == 0 {
						continue
					}

					subSteps = append(subSteps, tools.TaskExecSpec{
						Command: sliceutils.NewStrings(
							toolCmd, "manifest", "remove", localManifestName, digest,
						),
						IgnoreError: false,
					})
				}

				// add this image to manifest with correct os/arch/variant
				subSteps = append(subSteps, tools.TaskExecSpec{
					Command:     sliceutils.NewStrings(manifestAddCmd),
					IgnoreError: false,
				})

				return subSteps, nil
			},
		})
	}

	return steps, nil
}

func getLocalImageName(imageName string) string {
	return hex.EncodeToString(hashhelper.MD5Sum([]byte(imageName)))
}

func getLocalManifestName(manifestName string) string {
	return hex.EncodeToString(hashhelper.MD5Sum([]byte(manifestName)))
}

func getImageIDFilePathForImageName(dukkhaCacheDir, imageName string) string {
	return filepath.Join(
		dukkhaCacheDir,
		"buildah",
		fmt.Sprintf(
			"image-id-%s", getLocalImageName(imageName),
		),
	)
}
