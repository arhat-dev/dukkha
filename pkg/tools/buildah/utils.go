package buildah

import (
	"strings"

	"github.com/huandu/xstrings"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/field"
)

func SetDefaultImageTagIfNoTagSet(rc *field.RenderingContext, name string) string {
	return setDefaultTagIfNoTagSet(rc, name, false)
}

func SetDefaultManifestTagIfNoTagSet(rc *field.RenderingContext, manfiestName string) string {
	return setDefaultTagIfNoTagSet(rc, manfiestName, true)
}

func setDefaultTagIfNoTagSet(
	rc *field.RenderingContext,
	name string,
	isManifest bool,
) string {
	if hasTag(name) {
		return name
	}

	rawBranch := rc.Values().Env[constant.ENV_GIT_BRANCH]
	branch := xstrings.ToKebabCase(strings.ReplaceAll(rawBranch, "/", "-"))

	workTreeClean := rc.Values().Env[constant.ENV_GIT_WORKTREE_CLEAN] == "true"
	arch := rc.Values().Env[constant.ENV_MATRIX_ARCH]

	var tag string
	if workTreeClean {
		gitTag := rc.Values().Env[constant.ENV_GIT_TAG]
		switch {
		case len(gitTag) != 0:
			tag = gitTag
		case rawBranch == rc.Values().Env[constant.ENV_GIT_DEFAULT_BRANCH]:
			tag = "latest"
		default:
			tag = branch
			if !isManifest {
				tag += "-" + rc.Values().Env[constant.ENV_GIT_COMMIT]
			}
		}
	} else {
		// is expected to always pull without knowing image digest
		tag = "dev-" + branch
	}

	if !isManifest && len(arch) != 0 {
		tag += "-" + arch
	}

	return name + ":" + tag
}

func hasTag(name string) bool {
	tagIndex := strings.LastIndexByte(name, ':')
	if tagIndex < 0 {
		return false
	}

	// has tag if there is no more slash after
	return strings.IndexByte(name[tagIndex+1:], '/') < 0
}
