package buildah

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"arhat.dev/pkg/textquery"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/templateutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindXBuild = "xbuild"

const (
	replace_XBUILD_CURRENT_CONTAINER_ID = "<XBUILD_CURRENT_CONTAINER_ID>"
	// replace_XBUILD_CONTEXT_DIR          = "<XBUILD_CONTEXT_DIR>"
)

// nolint:revive
func replace_XBUILD_STEP_CONTAINER_ID(stepID string) string {
	return "<XBUILD_STEP_CONTAINER_ID_" + stepID + ">"
}

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindXBuild,
		func(toolName string) dukkha.Task {
			t := &TaskXBuild{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), t)
			return t
		},
	)
}

type TaskXBuild struct {
	rs.BaseField `yaml:"-"`

	TaskName string `yaml:"name"`

	tools.BaseTask `yaml:",inline"`

	// Context string  `yaml:"context"`
	Steps []*step `yaml:"steps"`

	ImageNames []*ImageNameSpec `yaml:"image_names"`
}

func (w *TaskXBuild) Kind() dukkha.TaskKind { return TaskKindXBuild }
func (w *TaskXBuild) Name() dukkha.TaskName { return dukkha.TaskName(w.TaskName) }
func (w *TaskXBuild) Key() dukkha.TaskKey {
	return dukkha.TaskKey{Kind: w.Kind(), Name: w.Name()}
}

// nolint:gocyclo
func (w *TaskXBuild) GetExecSpecs(
	rc dukkha.TaskExecContext,
	options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var ret []dukkha.TaskExecSpec

	err := w.DoAfterFieldsResolved(rc, -1, true, func() error {
		tmpImageIDFile, err := os.CreateTemp(rc.CacheDir(), "buildah-xbuild-image-id-*")
		if err != nil {
			return fmt.Errorf("failed to create a temp file for image id: %w", err)
		}
		tmpImageIDFilePath := tmpImageIDFile.Name()
		_ = tmpImageIDFile.Close()

		var (
			stepIDs      []string
			imageIDFiles []string
		)

		var realImageNames []string

		nameSum := sha256.New()
		for _, spec := range w.ImageNames {
			if len(spec.Image) == 0 {
				continue
			}

			imageName := templateutils.SetDefaultImageTagIfNoTagSet(
				rc, spec.Image, true,
			)

			_, err := nameSum.Write([]byte(imageName))
			if err != nil {
				return fmt.Errorf("failed to write image name to name sum: %w", err)
			}

			realImageNames = append(realImageNames, imageName)

			filePath, err := GetImageIDFileForImageName(rc, imageName)
			if err != nil {
				return err
			}

			imageIDFiles = append(imageIDFiles, filePath)
		}

		// generate deterministic image name
		finalImageName := "buildah-xbuild-" + hex.EncodeToString(nameSum.Sum(nil))

		// 		// set context dir
		// 		contextDir, err := filepath.Abs(w.Context)
		// 		if err != nil {
		// 			return fmt.Errorf("failed to get absolute path of context dir: %w", err)
		// 		}
		//
		// 		ret = append(ret, dukkha.TaskExecSpec{
		// 			StdoutAsReplace: replace_XBUILD_CONTEXT_DIR,
		// 			AlterExecFunc: func(
		// 				replace dukkha.ReplaceEntries,
		// 				stdin io.Reader,
		// 				stdout, stderr io.Writer,
		// 			) (dukkha.RunTaskOrRunCmd, error) {
		// 				_, err := stdout.Write([]byte(contextDir))
		// 				return nil, err
		// 			},
		// 		})

		for i, step := range w.Steps {
			if step.Skip {
				continue
			}

			stepID := step.ID
			if len(stepID) == 0 {
				stepID = strconv.FormatInt(int64(i), 10)
			}

			stepIDs = append(stepIDs, stepID)

			// set default container id of this step
			ret = append(ret, dukkha.TaskExecSpec{
				StdoutAsReplace:          replace_XBUILD_STEP_CONTAINER_ID(stepID),
				FixStdoutValueForReplace: bytes.TrimSpace,

				AlterExecFunc: func(
					replace dukkha.ReplaceEntries,
					stdin io.Reader,
					stdout, stderr io.Writer,
				) (dukkha.RunTaskOrRunCmd, error) {
					v, ok := replace[replace_XBUILD_CURRENT_CONTAINER_ID]
					if !ok {
						return nil, nil
					}

					// TODO: remove additional \n for ansi translation flush
					_, err := stdout.Write(append(v.Data, '\n'))
					return nil, err
				},
			})

			// add this step to global step index

			stepRet, err := step.genSpec(rc, w.CacheFS)
			if err != nil {
				return fmt.Errorf("failed to generate #%d step spec: %w", i, err)
			}

			ret = append(ret, stepRet...)

			// update container id when switching image
			if step.From != nil {
				ret = append(ret, dukkha.TaskExecSpec{
					StdoutAsReplace:          replace_XBUILD_STEP_CONTAINER_ID(stepID),
					FixStdoutValueForReplace: bytes.TrimSpace,

					AlterExecFunc: func(
						replace dukkha.ReplaceEntries,
						stdin io.Reader,
						stdout, stderr io.Writer,
					) (dukkha.RunTaskOrRunCmd, error) {
						v, ok := replace[replace_XBUILD_CURRENT_CONTAINER_ID]
						if !ok {
							return nil, nil
						}

						// TODO: remove additional \n for ansi translation flush
						_, err := stdout.Write(append(v.Data, '\n'))
						return nil, err
					},
				})
			}

			isLastStep := i == len(w.Steps)-1
			if !isLastStep {
				// is last step when all following steps are skipped
				isLastStep = true
				for _, s := range w.Steps[i+1:] {
					if !s.Skip {
						isLastStep = false
						break
					}
				}
			}

			// commit this container as image
			var imageName string
			commitCmd := []string{constant.DUKKHA_TOOL_CMD, "commit"}
			if step.Compress != nil && !*step.Compress {
				commitCmd = append(commitCmd, "--disable-compression")
			} else {
				commitCmd = append(commitCmd, "--disable-compression=false")
			}

			switch {
			case isLastStep:
				// at last step, commit with image names
				imageName = finalImageName
				if len(step.CommitAs) != 0 {
					realImageNames = append(realImageNames, step.CommitAs)
				}

				commitCmd = append(commitCmd, "--iidfile", tmpImageIDFilePath)
			case w.Steps[i+1].From != nil:
				// next step is a from statement
				imageName = step.CommitAs
				if len(imageName) == 0 {
					// TODO: generate image name for multi-step build, currently we do not commit if not set
					continue
				}
			case step.Commit != nil && *step.Commit:
				// set commit=true explicitly
				imageName = step.CommitAs
				if len(imageName) == 0 {
					// TODO: generate image name for intermediate layer
					continue
				}
			default:
				continue
			}

			commitCmd = append(commitCmd, step.ExtraCommitArgs...)
			commitCmd = append(commitCmd, replace_XBUILD_CURRENT_CONTAINER_ID, imageName)

			ret = append(ret, dukkha.TaskExecSpec{
				IgnoreError: false,
				Command:     commitCmd,
			})
		}

		// delete running containers
		ret = append(ret, dukkha.TaskExecSpec{
			AlterExecFunc: func(
				replace dukkha.ReplaceEntries,
				stdin io.Reader,
				stdout, stderr io.Writer,
			) (dukkha.RunTaskOrRunCmd, error) {
				var delActions []dukkha.TaskExecSpec

				visitedCtrIDs := make(map[string]struct{})
				for i := len(stepIDs) - 1; i >= 0; i-- {
					stepID := stepIDs[i]
					v, ok := replace[replace_XBUILD_STEP_CONTAINER_ID(stepID)]
					if !ok {
						return nil, fmt.Errorf("unexpected missing container id of step %q", stepID)
					}
					ctrID := string(v.Data)
					_, visited := visitedCtrIDs[ctrID]
					if visited {
						continue
					}

					visitedCtrIDs[ctrID] = struct{}{}
					delActions = append(delActions, dukkha.TaskExecSpec{
						IgnoreError: false,
						Command:     []string{constant.DUKKHA_TOOL_CMD, "rm", ctrID},
					})
				}

				return delActions, nil
			},
		})

		// retrieve and write image id

		const (
			replace_XBUILD_IMAGE_ID = "<XBUILD_IMAGE_ID>"
		)
		ret = append(ret, dukkha.TaskExecSpec{
			StdoutAsReplace:          replace_XBUILD_IMAGE_ID,
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

		// create tags
		for _, imageName := range realImageNames {
			ret = append(ret, dukkha.TaskExecSpec{
				IgnoreError: false,
				Command:     []string{constant.DUKKHA_TOOL_CMD, "tag", replace_XBUILD_IMAGE_ID, imageName},
			})
		}

		// update manifests
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

		for _, spec := range w.ImageNames {
			if len(spec.Manifest) == 0 {
				continue
			}

			manifestName := templateutils.SetDefaultManifestTagIfNoTagSet(rc, spec.Manifest)
			localManifestName := getLocalManifestName(manifestName)
			// ensure local manifest exists
			ret = append(ret, dukkha.TaskExecSpec{
				Command:     []string{constant.DUKKHA_TOOL_CMD, "manifest", "create", localManifestName},
				IgnoreError: true,
			})

			const replaceTargetManifestSpec = "<MANIFEST_SPEC>"
			ret = append(ret, dukkha.TaskExecSpec{
				StdoutAsReplace:          replaceTargetManifestSpec,
				FixStdoutValueForReplace: nil,

				Command: []string{constant.DUKKHA_TOOL_CMD, "manifest", "inspect", localManifestName},
				// manifest may not exist
				IgnoreError: true,
			})

			manifestAddCmd := []string{constant.DUKKHA_TOOL_CMD, "manifest", "add"}
			manifestAddCmd = append(manifestAddCmd, osArchVariantArgs...)
			manifestAddCmd = append(manifestAddCmd, localManifestName, replace_XBUILD_IMAGE_ID)

			// find existing manifest entries with same os/arch/variant
			ret = append(ret, dukkha.TaskExecSpec{
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
								Command: []string{
									constant.DUKKHA_TOOL_CMD, "manifest", "create",
									localManifestName,
								},
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
				ret = append(ret, dukkha.TaskExecSpec{
					Command:     []string{constant.DUKKHA_TOOL_CMD, "manifest", "inspect", localManifestName},
					IgnoreError: false,
				})
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate build spec: %w", err)
	}

	return ret, nil
}
