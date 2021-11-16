package templateutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/constant"
	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/matrix"
	"arhat.dev/dukkha/pkg/sliceutils"
)

func TestSetDefaultImageTag(t *testing.T) {
	testMatrix := map[string][]string{
		constant.ENV_GIT_BRANCH:         {"eXtream/branch"},
		constant.ENV_GIT_DEFAULT_BRANCH: {"eXtream/branch", "different-branch"},
		constant.ENV_GIT_WORKTREE_CLEAN: {"true", "false"},
		constant.ENV_GIT_COMMIT:         {"commit"},
		constant.ENV_GIT_TAG:            {"tag", ""},
		constant.ENV_MATRIX_ARCH:        {"amd64"},
		constant.ENV_MATRIX_KERNEL:      {"linux"},
	}

	caseWorkTreeCleanTagPresent := map[string]string{
		constant.ENV_GIT_TAG:            "tag",
		constant.ENV_GIT_WORKTREE_CLEAN: "true",
	}

	caseWorkTreeCleanIsDefaultBranch := map[string]string{
		constant.ENV_GIT_BRANCH:         "eXtream/branch",
		constant.ENV_GIT_DEFAULT_BRANCH: "eXtream/branch",
		constant.ENV_GIT_WORKTREE_CLEAN: "true",
	}

	caseWorkTreeCleanNotDefaultBranch := map[string]string{
		constant.ENV_GIT_DEFAULT_BRANCH: "different-branch",
		constant.ENV_GIT_WORKTREE_CLEAN: "true",
	}

	caseWorkTreeDirty := map[string]string{
		constant.ENV_GIT_WORKTREE_CLEAN: "false",
	}

	tests := matrix.CartesianProduct(testMatrix)
	for _, mat := range tests {
		spec := matrix.Entry(mat)

		rc := dukkha_test.NewTestContextWithGlobalEnv(context.TODO(), mat)
		rc.SetCacheDir(t.TempDir())
		rc.AddListEnv(sliceutils.FormatStringMap(mat, "=", false)...)

		t.Run(spec.BriefString()+"_image_no_kernel_info", func(t *testing.T) {
			actual := SetDefaultImageTagIfNoTagSet(
				rc, SetDefaultImageTagIfNoTagSet(rc, "foo", false), false,
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
			actual := SetDefaultImageTagIfNoTagSet(
				rc, SetDefaultImageTagIfNoTagSet(rc, "foo", true), true,
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
			actual := SetDefaultManifestTagIfNoTagSet(
				rc, SetDefaultManifestTagIfNoTagSet(rc, "foo"),
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
