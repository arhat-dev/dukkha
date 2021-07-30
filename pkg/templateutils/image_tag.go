package templateutils

import (
	"strings"

	"github.com/huandu/xstrings"

	"arhat.dev/dukkha/pkg/dukkha"
)

func SetDefaultImageTagIfNoTagSet(
	rc dukkha.RenderingContext, imageName string, keepKernelInfo bool,
) string {
	return setDefaultTagIfNoTagSet(rc, imageName, false, keepKernelInfo)
}

func SetDefaultManifestTagIfNoTagSet(
	rc dukkha.RenderingContext, manfiestName string,
) string {
	return setDefaultTagIfNoTagSet(rc, manfiestName, true, false)
}

func setDefaultTagIfNoTagSet(
	rc dukkha.RenderingContext,
	name string,
	isManifest bool,
	keepKernelInfo bool,
) string {
	if hasTag(name) {
		return name
	}

	rawBranch := rc.GitBranch()
	branch := xstrings.ToKebabCase(strings.ReplaceAll(rawBranch, "/", "-"))

	workTreeClean := rc.GitWorkTreeClean()
	mArch := rc.MatrixArch()
	mKernel := rc.MatrixKernel()

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

	if !isManifest {
		if keepKernelInfo && len(mKernel) != 0 {
			tag += "-" + mKernel
		}

		if len(mArch) != 0 {
			tag += "-" + mArch
		}
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
