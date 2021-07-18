package buildah

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"arhat.dev/pkg/hashhelper"
	"arhat.dev/pkg/textquery"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindBuild = "build"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindBuild,
		func(toolName string) dukkha.Task {
			t := &TaskBuild{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindBuild, t)
			return t
		},
	)
}

type TaskBuild struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	Context    string          `yaml:"context"`
	ImageNames []ImageNameSpec `yaml:"image_names"`
	File       string          `yaml:"file"`

	// --build-arg
	BuildArgs []string `yaml:"build_args"`

	ExtraArgs []string `yaml:"extra_args"`
}

type ImageNameSpec struct {
	field.BaseField

	Image    string `yaml:"image"`
	Manifest string `yaml:"manifest"`
}

func (c *TaskBuild) GetExecSpecs(
	rc dukkha.TaskExecContext, options *dukkha.TaskExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var steps []dukkha.TaskExecSpec
	err := c.DoAfterFieldsResolved(rc, -1, func() error {
		ret, err := c.createExecSpecs(rc, options)
		steps = ret
		return err
	})

	return steps, err
}

func (c *TaskBuild) createExecSpecs(
	rc dukkha.TaskExecContext, options *dukkha.TaskExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	// create an image id file
	dukkhaCacheDir := rc.CacheDir()
	tmpImageIDFile, err := ioutil.TempFile(dukkhaCacheDir, "buildah-build-image-id-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create a temp file for image id: %w", err)
	}
	tmpImageIDFilePath := tmpImageIDFile.Name()
	_ = tmpImageIDFile.Close()

	budCmd := sliceutils.NewStrings(options.ToolCmd, "bud", "--iidfile", tmpImageIDFilePath)
	if len(c.File) != 0 {
		budCmd = append(budCmd, "-f", c.File)
	}

	for _, bArg := range c.BuildArgs {
		budCmd = append(budCmd, "--build-arg", bArg)
	}

	budCmd = append(budCmd, c.ExtraArgs...)

	targets := c.ImageNames
	if len(targets) == 0 {
		targets = []ImageNameSpec{{
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

		imageName := SetDefaultImageTagIfNoTagSet(rc, spec.Image)

		// local image name is to handle bud regression bugs related to
		// FQDN image names
		budCmd = append(budCmd, "-t", imageName)

		filePath := GetImageIDFileForImageName(
			dukkhaCacheDir, imageName,
		)
		err = os.MkdirAll(filepath.Dir(filePath), 0750)
		if err != nil && !os.IsExist(err) {
			return nil, fmt.Errorf("failed to ensure image id dir exists")
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
		Env:         sliceutils.NewStrings(c.Env),
		Command:     append(budCmd, context),
		IgnoreError: false,
		UseShell:    options.UseShell,
		ShellName:   options.ShellName,
	})

	// read image id file to get image id
	const replaceTargetImageID = "<IMAGE_ID>"
	steps = append(steps, dukkha.TaskExecSpec{
		OutputAsReplace:     replaceTargetImageID,
		FixOutputForReplace: bytes.TrimSpace,
		Env:                 sliceutils.NewStrings(c.Env),
		AlterExecFunc: func(
			replace map[string][]byte,
			stdin io.Reader, stdout, stderr io.Writer,
		) (dukkha.RunTaskOrRunCmd, error) {
			imageIDBytes, err := os.ReadFile(tmpImageIDFilePath)
			if err != nil {
				return nil, err
			}

			for _, f := range imageIDFiles {
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
	steps = append(steps, dukkha.TaskExecSpec{
		OutputAsReplace:     replaceTargetImageDigest,
		FixOutputForReplace: bytes.TrimSpace,
		Env:                 sliceutils.NewStrings(c.Env),
		Command: sliceutils.NewStrings(
			options.ToolCmd, "inspect", "--type", "image",
			"--format", `"{{ .FromImageDigest }}"`,
			replaceTargetImageID,
		),
		IgnoreError: false,
		UseShell:    options.UseShell,
		ShellName:   options.ShellName,
	})

	// add to manifest, ensure same os/arch/variant only one exist
	mArch := rc.MatrixArch()
	arch := constant.GetOciArch(mArch)
	os := constant.GetOciOS(rc.MatrixKernel())
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

		manifestName := SetDefaultManifestTagIfNoTagSet(rc, spec.Manifest)
		localManifestName := getLocalManifestName(manifestName)
		// ensure local manifest exists
		steps = append(steps, dukkha.TaskExecSpec{
			Env: sliceutils.NewStrings(c.Env),
			Command: sliceutils.NewStrings(
				options.ToolCmd, "manifest", "create", localManifestName,
			),
			IgnoreError: true,
			UseShell:    options.UseShell,
			ShellName:   options.ShellName,
		})

		const replaceTargetManifestSpec = "<MANIFEST_SPEC>"
		steps = append(steps, dukkha.TaskExecSpec{
			OutputAsReplace:     replaceTargetManifestSpec,
			FixOutputForReplace: nil,
			Env:                 sliceutils.NewStrings(c.Env),
			Command: sliceutils.NewStrings(
				options.ToolCmd, "manifest", "inspect", localManifestName,
			),
			// manifest may not exist
			IgnoreError: true,
			UseShell:    options.UseShell,
			ShellName:   options.ShellName,
		})

		manifestAddCmd := sliceutils.NewStrings(options.ToolCmd, "manifest", "add")
		manifestAddCmd = append(manifestAddCmd, osArchVariantArgs...)
		manifestAddCmd = append(manifestAddCmd, localManifestName, replaceTargetImageID)

		// find existing manifest entries with same os/arch/variant
		steps = append(steps, dukkha.TaskExecSpec{
			IgnoreError: false,
			Env:         sliceutils.NewStrings(c.Env),
			AlterExecFunc: func(
				replace map[string][]byte,
				stdin io.Reader, stdout, stderr io.Writer,
			) (dukkha.RunTaskOrRunCmd, error) {
				manifestSpec, ok := replace[replaceTargetManifestSpec]
				if !ok {
					// manifest not created, usually should not happen since we just created before
					return []dukkha.TaskExecSpec{
						{
							// do not ignore manifest create error this time
							Command: sliceutils.NewStrings(
								options.ToolCmd, "manifest", "create", localManifestName,
							),
							IgnoreError: false,
							UseShell:    options.UseShell,
							ShellName:   options.ShellName,
						},
						{
							Command:     sliceutils.NewStrings(manifestAddCmd),
							IgnoreError: false,
							UseShell:    options.UseShell,
							ShellName:   options.ShellName,
						},
					}, nil
				}

				// manifest already created, query to get all matching digests
				digestLines, err := textquery.JQBytes(manifestOsArchVariantQueryForDigest, manifestSpec)
				if err != nil {
					// no manifests entries, add this image directly
					return []dukkha.TaskExecSpec{{
						Command:     sliceutils.NewStrings(manifestAddCmd),
						IgnoreError: false,
						UseShell:    options.UseShell,
						ShellName:   options.ShellName,
					}}, nil
				}

				var subSteps []dukkha.TaskExecSpec

				// remove existing entries with same os/arch/variant
				for _, digest := range strings.Split(digestLines, "\n") {
					digest = strings.TrimSpace(digest)
					if len(digest) == 0 {
						continue
					}

					subSteps = append(subSteps, dukkha.TaskExecSpec{
						Command: sliceutils.NewStrings(
							options.ToolCmd, "manifest", "remove", localManifestName, digest,
						),
						IgnoreError: false,
						UseShell:    options.UseShell,
						ShellName:   options.ShellName,
					})
				}

				// add this image to manifest with correct os/arch/variant
				subSteps = append(subSteps, dukkha.TaskExecSpec{
					Command:     sliceutils.NewStrings(manifestAddCmd),
					IgnoreError: false,
					UseShell:    options.UseShell,
					ShellName:   options.ShellName,
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

func GetImageIDFileForImageName(dukkhaCacheDir, imageName string) string {
	return filepath.Join(
		dukkhaCacheDir,
		"buildah",
		fmt.Sprintf(
			"image-id-%s", getLocalImageName(imageName),
		),
	)
}
