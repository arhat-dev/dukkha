package buildah

import (
	"bytes"
	"fmt"
	"sort"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/templateutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindPush = "push"

func init() {
	dukkha.RegisterTask(ToolKind, TaskKindPush, tools.NewTask[TaskPush, *TaskPush])
}

type TaskPush struct {
	tools.BaseTask[BuildahPush, *BuildahPush]
}

type BuildahPush struct {
	ImageNames []ImageNameSpec `yaml:"image_names"`

	manifestCache map[manifestCacheKey]manifestCacheValue

	parent tools.BaseTaskType
}

func (w *BuildahPush) ToolKind() dukkha.ToolKind       { return ToolKind }
func (w *BuildahPush) Kind() dukkha.TaskKind           { return TaskKindPush }
func (w *BuildahPush) LinkParent(p tools.BaseTaskType) { w.parent = p }

func (c *BuildahPush) GetExecSpecs(
	rc dukkha.TaskExecContext,
	opts dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var result []dukkha.TaskExecSpec

	err := c.parent.DoAfterFieldsResolved(rc, -1, true, func() error {
		if c.manifestCache == nil {
			c.manifestCache = make(map[manifestCacheKey]manifestCacheValue)
		}

		targets := c.ImageNames
		if len(targets) == 0 {
			targets = []ImageNameSpec{
				{
					Image:    string(c.parent.Name()),
					Manifest: "",
				},
			}
		}

		for i, spec := range targets {
			if len(spec.Image) != 0 {
				imageName := templateutils.GetFullImageName_UseDefault_IfIfNoTagSet(rc, spec.Image, true)
				imageIDFile, err := GetImageIDFileForImageName(rc, imageName, false)
				if err != nil {
					return err
				}

				imageIDBytes, err := rc.FS().ReadFile(imageIDFile)
				if err != nil {
					return fmt.Errorf("image id file not found: %w", err)
				}

				result = append(result, dukkha.TaskExecSpec{
					Command: []string{constant.DUKKHA_TOOL_CMD, "push",
						string(bytes.TrimSpace(imageIDBytes)),
						// TODO: support other destination
						"docker://" + imageName,
					},
					IgnoreError: false,
				})
			}

			if len(spec.Manifest) == 0 {
				continue
			}

			manifestName := templateutils.GetFullManifestName_UseDefault_IfNoTagSet(rc, spec.Manifest)
			c.cacheManifestPushSpec(i, opts, manifestName)
		}

		// push all manifests at last
		if opts.IsLast() {
			result = append(result,
				c.createManifestPushSpecsFromCache(opts.ID())...,
			)
		}

		return nil
	})

	return result, err
}

type manifestCacheKey struct {
	execID int
	name   string
}

type manifestCacheValue struct {
	subIndex int
	name     string

	opts dukkha.TaskMatrixExecOptions
}

func (c *BuildahPush) cacheManifestPushSpec(
	index int,
	opts dukkha.TaskMatrixExecOptions,
	manifestName string,
) {
	key := manifestCacheKey{
		execID: opts.ID(),
		name:   manifestName,
	}

	c.manifestCache[key] = manifestCacheValue{
		subIndex: index,

		name: manifestName,
		opts: opts,
	}
}

func (c *BuildahPush) createManifestPushSpecsFromCache(execID int) []dukkha.TaskExecSpec {
	var (
		values []manifestCacheValue
	)

	// filter manifests belong to this exec
	for k, v := range c.manifestCache {
		if k.execID != execID {
			continue
		}

		values = append(values, v)
	}

	// restore original order
	sort.Slice(values, func(i, j int) bool {
		less := values[i].opts.Seq() < values[j].opts.Seq()
		if less {
			return true
		}

		return values[i].subIndex < values[j].subIndex
	})

	// generate specs using original options
	var ret []dukkha.TaskExecSpec
	for _, v := range values {
		delete(c.manifestCache, manifestCacheKey{
			execID: v.opts.ID(),
			name:   v.name,
		})

		// buildah manifest push --all \
		//   <manifest-list-name> <transport>:<transport-details>
		ret = append(ret, dukkha.TaskExecSpec{
			Command: []string{constant.DUKKHA_TOOL_CMD, "manifest", "push", "--all",
				getLocalManifestName(v.name),
				// TODO: support other destination
				"docker://" + v.name,
			},
			IgnoreError: false,
		})
	}

	return ret
}
