package buildah

import (
	"strings"

	"github.com/huandu/xstrings"

	"arhat.dev/dukkha/pkg/types"
)

func SetDefaultImageTagIfNoTagSet(rc types.RenderingContext, name string) string {
	return setDefaultTagIfNoTagSet(rc, name, false)
}

func SetDefaultManifestTagIfNoTagSet(rc types.RenderingContext, manfiestName string) string {
	return setDefaultTagIfNoTagSet(rc, manfiestName, true)
}

func setDefaultTagIfNoTagSet(
	rc types.RenderingContext,
	name string,
	isManifest bool,
) string {
	if hasTag(name) {
		return name
	}

	rawBranch := rc.GitBranch()
	branch := xstrings.ToKebabCase(strings.ReplaceAll(rawBranch, "/", "-"))

	workTreeClean := rc.GitWorkTreeClean()
	arch := rc.MatrixArch()

	var tag string
	if workTreeClean {
		gitTag := rc.GitTag()
		switch {
		case len(gitTag) != 0:
			tag = gitTag
		case rawBranch == rc.GitDefaultBranch():
			tag = "latest"
		default:
			tag = branch
			if !isManifest {
				tag += "-" + rc.GitCommit()
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
