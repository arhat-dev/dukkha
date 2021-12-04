package buildah

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"arhat.dev/pkg/md5helper"
	"arhat.dev/pkg/textquery"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/templateutils"
	"arhat.dev/dukkha/pkg/tools"
)

const (
	TaskKindBuild = "build"
)

func init() {
	dukkha.RegisterTask(ToolKind, TaskKindBuild, newTaskBuild)

	templateutils.RegisterTemplateFuncs(map[string]templateutils.TemplateFuncFactory{
		"getBuildahImageIDFile": func(rc dukkha.RenderingContext) interface{} {
			return func(imageName string) (string, error) {
				return GetImageIDFileForImageName(
					rc,
					templateutils.SetDefaultImageTagIfNoTagSet(rc, imageName, true),
				)
			}
		},
	})
}

func newTaskBuild(toolName string) dukkha.Task {
	t := &TaskBuild{}
	t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), t)
	return t
}

type TaskBuild struct {
	rs.BaseField `yaml:"-"`

	TaskName string `yaml:"name"`

	tools.BaseTask `yaml:",inline"`

	Context    string           `yaml:"context"`
	ImageNames []*ImageNameSpec `yaml:"image_names"`
	File       string           `yaml:"file"`

	// --build-arg
	BuildArgs []string `yaml:"build_args"`

	ExtraArgs []string `yaml:"extra_args"`
}

type ImageNameSpec struct {
	rs.BaseField `yaml:"-"`

	Image    string `yaml:"image"`
	Manifest string `yaml:"manifest"`
}

func (c *TaskBuild) Kind() dukkha.TaskKind { return TaskKindBuild }
func (c *TaskBuild) Name() dukkha.TaskName { return dukkha.TaskName(c.TaskName) }
func (c *TaskBuild) Key() dukkha.TaskKey {
	return dukkha.TaskKey{Kind: c.Kind(), Name: c.Name()}
}

func (c *TaskBuild) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var steps []dukkha.TaskExecSpec
	err := c.DoAfterFieldsResolved(rc, -1, true, func() error {
		ret, err := c.createExecSpecs(rc, options)
		steps = ret
		return err
	})

	return steps, err
}

func (c *TaskBuild) createExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	// create an image id file
	dukkhaCacheDir := rc.CacheDir()
	tmpImageIDFile, err := os.CreateTemp(dukkhaCacheDir, "buildah-build-image-id-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create a temp file for image id: %w", err)
	}
	tmpImageIDFilePath := tmpImageIDFile.Name()
	_ = tmpImageIDFile.Close()

	budCmd := []string{constant.DUKKHA_TOOL_CMD, "bud", "--iidfile", tmpImageIDFilePath}
	if len(c.File) != 0 {
		budCmd = append(budCmd, "-f", c.File)
	}

	for _, bArg := range c.BuildArgs {
		budCmd = append(budCmd, "--build-arg", bArg)
	}

	budCmd = append(budCmd, c.ExtraArgs...)

	targets := c.ImageNames
	if len(targets) == 0 {
		targets = []*ImageNameSpec{{
			Image:    c.TaskName,
			Manifest: "",
		}}
	}

	// set image names
	var imageIDFiles []string
	for _, spec := range targets {
		if len(spec.Image) == 0 {
			continue
		}

		imageName := templateutils.SetDefaultImageTagIfNoTagSet(
			rc, spec.Image, true,
		)

		// local image name is to handle bud regression bugs related to
		// FQDN image names
		budCmd = append(budCmd, "-t", imageName)

		filePath, err := GetImageIDFileForImageName(rc, imageName)
		if err != nil {
			return nil, err
		}

		imageIDFiles = append(imageIDFiles, filePath)
	}

	context := c.Context
	if len(context) == 0 {
		context = "."
	}

	var steps []dukkha.TaskExecSpec

	// buildah bud
	steps = append(steps, dukkha.TaskExecSpec{
		Command:     append(budCmd, context),
		IgnoreError: false,
	})

	// read image id file to get image id
	const replace_TARGET_IMAGE_ID = "<IMAGE_ID>"
	steps = append(steps, dukkha.TaskExecSpec{
		StdoutAsReplace:          replace_TARGET_IMAGE_ID,
		FixStdoutValueForReplace: bytes.TrimSpace,
		AlterExecFunc: func(
			replace dukkha.ReplaceEntries,
			stdin io.Reader, stdout, stderr io.Writer,
		) (dukkha.RunTaskOrRunCmd, error) {
			imageIDBytes, err := rc.FS().ReadFile(tmpImageIDFilePath)
			if err != nil {
				return nil, err
			}

			for _, f := range imageIDFiles {
				err = rc.FS().WriteFile(f, imageIDBytes, 0750)
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
	const replace_TARGET_IMAGE_DIGEST = "<IMAGE_DIGEST>"
	steps = append(steps, dukkha.TaskExecSpec{
		StdoutAsReplace: replace_TARGET_IMAGE_DIGEST,
		ShowStdout:      true,

		FixStdoutValueForReplace: bytes.TrimSpace,
		Command: []string{
			constant.DUKKHA_TOOL_CMD, "inspect", "--type", "image",
			"--format", `"{{ .FromImageDigest }}"`,
			replace_TARGET_IMAGE_ID,
		},
		IgnoreError: false,
	})

	// add to manifest, ensure same os/arch/variant only one exist
	mArch := rc.MatrixArch()
	variant, _ := constant.GetOciArchVariant(mArch)
	os, _ := constant.GetOciOS(rc.MatrixKernel())
	arch, _ := constant.GetOciArch(mArch)

	osArchVariantArgs := []string{"--os", os, "--arch", arch}
	if len(variant) != 0 {
		osArchVariantArgs = append(osArchVariantArgs, "--variant", variant)
	}

	manifestOsArchVariantQueryForDigest := createManifestPlatformQueryForDigest(
		rc.MatrixKernel(), mArch,
	)

	for _, spec := range targets {
		if len(spec.Manifest) == 0 {
			continue
		}

		manifestName := templateutils.SetDefaultManifestTagIfNoTagSet(rc, spec.Manifest)
		localManifestName := getLocalManifestName(manifestName)
		// ensure local manifest exists
		steps = append(steps, dukkha.TaskExecSpec{
			Command:     []string{constant.DUKKHA_TOOL_CMD, "manifest", "create", localManifestName},
			IgnoreError: true,
		})

		const replaceTargetManifestSpec = "<MANIFEST_SPEC>"
		steps = append(steps, dukkha.TaskExecSpec{
			StdoutAsReplace:          replaceTargetManifestSpec,
			FixStdoutValueForReplace: nil,
			Command:                  []string{constant.DUKKHA_TOOL_CMD, "manifest", "inspect", localManifestName},
			// manifest may not exist
			IgnoreError: true,
		})

		manifestAddCmd := []string{constant.DUKKHA_TOOL_CMD, "manifest", "add"}
		manifestAddCmd = append(manifestAddCmd, osArchVariantArgs...)
		manifestAddCmd = append(manifestAddCmd, localManifestName, replace_TARGET_IMAGE_ID)

		// find existing manifest entries with same os/arch/variant
		steps = append(steps, dukkha.TaskExecSpec{
			IgnoreError: false,
			AlterExecFunc: func(
				replace dukkha.ReplaceEntries,
				stdin io.Reader, stdout, stderr io.Writer,
			) (dukkha.RunTaskOrRunCmd, error) {
				manifestSpec, ok := replace[replaceTargetManifestSpec]
				if !ok {
					// manifest not created, usually should not happen since we just created before
					return []dukkha.TaskExecSpec{
						{
							// do not ignore manifest create error this time
							Command:     []string{constant.DUKKHA_TOOL_CMD, "manifest", "create", localManifestName},
							IgnoreError: false,
						},
						{
							Command:     sliceutils.NewStrings(manifestAddCmd),
							IgnoreError: false,
						},
					}, nil
				}

				// manifest already created, query to get all matching digests
				digestResult, err := textquery.JQBytes(
					manifestOsArchVariantQueryForDigest, manifestSpec.Data,
				)
				if err != nil {
					// no manifests entries, add this image directly
					return []dukkha.TaskExecSpec{{
						Command:     sliceutils.NewStrings(manifestAddCmd),
						IgnoreError: false,
					}}, nil
				}

				digests, err := parseManifestOsArchVariantQueryResult(digestResult)
				if err != nil {
					return nil, fmt.Errorf("failed to parse digest result: %w", err)
				}

				var subSteps []dukkha.TaskExecSpec

				// remove existing entries with same os/arch/variant
				for _, digest := range digests {
					digest = strings.TrimSpace(digest)
					if len(digest) == 0 {
						continue
					}

					subSteps = append(subSteps, dukkha.TaskExecSpec{
						Command: []string{
							constant.DUKKHA_TOOL_CMD, "manifest", "remove",
							localManifestName, digest,
						},
						IgnoreError: false,
					})
				}

				// add this image to manifest with correct os/arch/variant
				subSteps = append(subSteps, dukkha.TaskExecSpec{
					Command:     sliceutils.NewStrings(manifestAddCmd),
					IgnoreError: false,
				})

				return subSteps, nil
			},
		})

		// check manifests in last matrix execution
		if options.IsLast() {
			steps = append(steps, dukkha.TaskExecSpec{
				Command:     []string{constant.DUKKHA_TOOL_CMD, "manifest", "inspect", localManifestName},
				IgnoreError: false,
			})
		}
	}

	return steps, nil
}

func parseManifestOsArchVariantQueryResult(result string) ([]string, error) {
	var data interface{}

	err := json.Unmarshal([]byte(result), &data)
	if err != nil {
		// plain text
		return []string{result}, nil
	}

	switch t := data.(type) {
	case []interface{}:
		var ret []string
		for _, v := range t {
			if r, ok := v.(string); ok {
				ret = append(ret, r)
			} else {
				return nil, fmt.Errorf("unexpected non string digest %T", v)
			}
		}

		return ret, nil
	default:
		return nil, fmt.Errorf("unexpected result type %T, want []interface{}", t)
	}
}

func createManifestPlatformQueryForDigest(mKernel, mArch string) string {
	os, _ := constant.GetOciOS(mKernel)
	arch, _ := constant.GetOciArch(mArch)
	variant, _ := constant.GetOciArchVariant(mArch)

	manifestOsArchVariantQueryForDigest := fmt.Sprintf(
		`.manifests[] | select((.platform.os == "%s")`+
			` and (.platform.architecture == "%s")`,
		os,
		arch,
	)

	if len(variant) != 0 {
		manifestOsArchVariantQueryForDigest += fmt.Sprintf(
			` and (.platform.variant == "%s")`, variant,
		)
	} else {
		manifestOsArchVariantQueryForDigest += ` and select(.platform.variant == "" or  .platform.variant == null)`
	}

	return manifestOsArchVariantQueryForDigest + `) | .digest`
}

func getLocalImageName(imageName string) string {
	return hex.EncodeToString(md5helper.Sum([]byte(imageName)))
}

func getLocalManifestName(manifestName string) string {
	return hex.EncodeToString(md5helper.Sum([]byte(manifestName)))
}

func GetImageIDFileForImageName(rc dukkha.RenderingContext, imageName string) (string, error) {
	const imageIDCacheDir = "buildah/image-id"

	ret, err := rc.GlobalCacheFS(imageIDCacheDir).
		Abs(getLocalImageName(imageName))
	if err != nil {
		return "", err
	}

	return ret, err
}
