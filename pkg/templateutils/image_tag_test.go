package templateutils

import (
	"context"
	"sort"
	"testing"

	"arhat.dev/pkg/matrixhelper"
	"arhat.dev/pkg/sorthelper"
	"arhat.dev/tlang"
	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	dt "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/matrix"
	"arhat.dev/dukkha/pkg/sliceutils"
)

func TestSetDefaultImageTag(t *testing.T) {
	t.Parallel()

	testMatrix := map[string][]string{
		constant.EnvName_GIT_BRANCH:         {"eXtream/branch"},
		constant.EnvName_GIT_DEFAULT_BRANCH: {"eXtream/branch", "different-branch"},
		constant.EnvName_GIT_WORKTREE_CLEAN: {"true", "false"},
		constant.EnvName_GIT_COMMIT:         {"commit"},
		constant.EnvName_GIT_TAG:            {"tag", ""},
		constant.EnvName_MATRIX_ARCH:        {"amd64"},
		constant.EnvName_MATRIX_KERNEL:      {"linux"},
	}

	caseWorkTreeCleanTagPresent := map[string]string{
		constant.EnvName_GIT_TAG:            "tag",
		constant.EnvName_GIT_WORKTREE_CLEAN: "true",
	}

	caseWorkTreeCleanIsDefaultBranch := map[string]string{
		constant.EnvName_GIT_BRANCH:         "eXtream/branch",
		constant.EnvName_GIT_DEFAULT_BRANCH: "eXtream/branch",
		constant.EnvName_GIT_WORKTREE_CLEAN: "true",
	}

	caseWorkTreeCleanNotDefaultBranch := map[string]string{
		constant.EnvName_GIT_DEFAULT_BRANCH: "different-branch",
		constant.EnvName_GIT_WORKTREE_CLEAN: "true",
	}

	caseWorkTreeDirty := map[string]string{
		constant.EnvName_GIT_WORKTREE_CLEAN: "false",
	}

	tests := matrixhelper.CartesianProduct(testMatrix, func(names []string, mat [][]string) {
		sort.Sort(sorthelper.NewCustomSortable(
			func(i, j int) {
				names[i], names[j] = names[j], names[i]
				mat[i], mat[j] = mat[j], mat[i]
			},
			func(i, j int) bool { return names[i] < names[j] },
			func() int { return len(names) },
		))
	})

	for _, mat := range tests {
		spec := matrix.Entry(mat)

		genv := &dukkha.GlobalEnvSet{
			constant.GlobalEnv_DUKKHA_CACHE_DIR: tlang.ImmediateString(t.TempDir()),
		}
		for k, v := range mat {
			id := constant.GetGlobalEnvIDByName(k)
			if id == -1 {
				continue
			}

			genv[id] = tlang.ImmediateString(v)
		}

		rc := dt.NewTestContextWithGlobalEnv(context.TODO(), genv)
		for k, v := range mat {
			if constant.GetGlobalEnvIDByName(k) == -1 {
				rc.AddEnv(false, &dukkha.NameValueEntry{
					Name:  k,
					Value: v,
				})
			}
		}

		rc.AddListEnv(sliceutils.FormatStringMap(mat, "=", false)...)

		t.Run(spec.BriefString()+"_image_no_kernel_info", func(t *testing.T) {
			actual := GetFullImageName_UseDefault_IfIfNoTagSet(
				rc, GetFullImageName_UseDefault_IfIfNoTagSet(rc, "foo", false), false,
			)
			switch {
			case spec.Match(caseWorkTreeCleanTagPresent):
				assert.Equal(t, "foo:tag-amd64", actual)
			case spec.Match(caseWorkTreeCleanIsDefaultBranch):
				assert.Equal(t, "foo:latest-amd64", actual)
			case spec.Match(caseWorkTreeCleanNotDefaultBranch):
				assert.Equal(t, "foo:e-xtream-branch-commit-amd64", actual)
			case spec.Match(caseWorkTreeDirty):
				assert.Equal(t, "foo:dev-e-xtream-branch-amd64", actual)
			default:
				assert.Fail(t, "unmatched condition")
			}
		})

		t.Run(spec.BriefString()+"_image_with_kernel_info", func(t *testing.T) {
			actual := GetFullImageName_UseDefault_IfIfNoTagSet(
				rc, GetFullImageName_UseDefault_IfIfNoTagSet(rc, "foo", true), true,
			)
			switch {
			case spec.Match(caseWorkTreeCleanTagPresent):
				assert.Equal(t, "foo:tag-linux-amd64", actual)
			case spec.Match(caseWorkTreeCleanIsDefaultBranch):
				assert.Equal(t, "foo:latest-linux-amd64", actual)
			case spec.Match(caseWorkTreeCleanNotDefaultBranch):
				assert.Equal(t, "foo:e-xtream-branch-commit-linux-amd64", actual)
			case spec.Match(caseWorkTreeDirty):
				assert.Equal(t, "foo:dev-e-xtream-branch-linux-amd64", actual)
			default:
				assert.Fail(t, "unmatched condition")
			}
		})

		t.Run(spec.BriefString()+"_manifest", func(t *testing.T) {
			actual := GetFullManifestName_UseDefault_IfNoTagSet(
				rc, GetFullManifestName_UseDefault_IfNoTagSet(rc, "foo"),
			)
			switch {
			case spec.Match(caseWorkTreeCleanTagPresent):
				assert.Equal(t, "foo:tag", actual)
			case spec.Match(caseWorkTreeCleanIsDefaultBranch):
				assert.Equal(t, "foo:latest", actual)
			case spec.Match(caseWorkTreeCleanNotDefaultBranch):
				assert.Equal(t, "foo:e-xtream-branch", actual)
			case spec.Match(caseWorkTreeDirty):
				assert.Equal(t, "foo:dev-e-xtream-branch", actual)
			default:
				assert.Fail(t, "unmatched condition")
			}
		})
	}
}

func TestHasTag(t *testing.T) {
	t.Parallel()

	t.Run("HasTag", func(t *testing.T) {
		for _, s := range []string{
			"local:latest",
			":no-name",
			":",
			"local:8080/image:tag",
		} {
			assert.True(t, hasTag(s), s)
		}
	})

	t.Run("NoTag", func(t *testing.T) {
		for _, s := range []string{
			"local",
			"local:8080/image",
		} {
			assert.False(t, hasTag(s), s)
		}
	})
}
