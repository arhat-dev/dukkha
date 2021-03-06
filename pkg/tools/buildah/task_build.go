package buildah

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
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

// buildahNS for template funcs related to buildah
type buildahNS struct{ rc dukkha.RenderingContext }

func (ns buildahNS) ImageIDFile(imageName string) (string, error) {
	return GetImageIDFileForImageName(
		ns.rc,
		templateutils.GetFullImageName_UseDefault_IfIfNoTagSet(ns.rc, imageName, true),
		false,
	)
}

func init() {
	dukkha.RegisterTask(ToolKind, TaskKindBuild, tools.NewTask[TaskBuild, *TaskBuild])

	templateutils.RegisterTemplateFuncs(map[string]templateutils.TemplateFuncFactory{
		"buildah": func(rc dukkha.RenderingContext) any { return buildahNS{rc: rc} },
	})
}

type TaskBuild struct {
	tools.BaseTask[BuildahBuild, *BuildahBuild]
}

// nolint:revive
type BuildahBuild struct {
	Context    string           `yaml:"context"`
	ImageNames []*ImageNameSpec `yaml:"image_names"`
	File       string           `yaml:"file"`

	// --build-arg
	BuildArgs []string `yaml:"build_args"`

	ExtraArgs []string `yaml:"extra_args"`

	parent tools.BaseTaskType
}

func (c *BuildahBuild) ToolKind() dukkha.ToolKind       { return ToolKind }
func (c *BuildahBuild) Kind() dukkha.TaskKind           { return TaskKindBuild }
func (c *BuildahBuild) LinkParent(p tools.BaseTaskType) { c.parent = p }

type ImageNameSpec struct {
	rs.BaseField `yaml:"-"`

	Image    string `yaml:"image"`
	Manifest string `yaml:"manifest"`
}

func (c *BuildahBuild) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var steps []dukkha.TaskExecSpec
	err := c.parent.DoAfterFieldsResolved(rc, -1, true, func() error {
		ret, err := c.createExecSpecs(rc, options)
		steps = ret
		return err
	})

	return steps, err
}

func (c *BuildahBuild) createExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	// create an image id file
	dukkhaCacheDir := rc.CacheDir()
	tmpImageIDFile, err := os.CreateTemp(dukkhaCacheDir, "buildah-build-image-id-*")
	if err != nil {
		return nil, fmt.Errorf("creating temp file for image id: %w", err)
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
			Image:    string(c.parent.Name()),
			Manifest: "",
		}}
	}

	// set image names
	var imageIDFiles []string
	for _, spec := range targets {
		if len(spec.Image) == 0 {
			continue
		}

		imageName := templateutils.GetFullImageName_UseDefault_IfIfNoTagSet(
			rc, spec.Image, true,
		)

		// local image name is to handle bud regression bugs related to
		// FQDN image names
		budCmd = append(budCmd, "-t", imageName)

		imageIDPath, err := GetImageIDFileForImageName(rc, imageName, true)
		if err != nil {
			return nil, err
		}

		imageIDFiles = append(imageIDFiles, imageIDPath)
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

		manifestName := templateutils.GetFullManifestName_UseDefault_IfNoTagSet(rc, spec.Manifest)
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
				digestResult, err := textquery.JQ[byte](
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
					return nil, fmt.Errorf("parsing digest result: %w", err)
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

func GetImageIDFileForImageName(rc dukkha.RenderingContext, imageName string, ensureDir bool) (string, error) {
	const imageIDCacheDir = "buildah/image-id"

	cfs := rc.GlobalCacheFS(imageIDCacheDir)

	ret, err := cfs.Abs(getLocalImageName(imageName))
	if err != nil {
		return "", err
	}

	if ensureDir {
		err = cfs.MkdirAll(".", 0755)
		if err != nil && !errors.Is(err, fs.ErrExist) {
			return "", err
		}
	}

	return ret, err
}
