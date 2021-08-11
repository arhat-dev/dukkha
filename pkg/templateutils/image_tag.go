package templateutils

import (
	"strings"

	"github.com/huandu/xstrings"

	"arhat.dev/dukkha/pkg/dukkha"
)

func GetDefaultImageTag(
	rc dukkha.RenderingContext, imageName string, keepKernelInfo bool,
) string {
	return GetDefaultTag(rc, imageName, false, keepKernelInfo)
}

func GetDefaultManifestTag(
	rc dukkha.RenderingContext, manifestName string,
) string {
	return GetDefaultTag(rc, manifestName, true, false)
}

func SetDefaultImageTagIfNoTagSet(
	rc dukkha.RenderingContext, imageName string, keepKernelInfo bool,
) string {
	if hasTag(imageName) {
		return imageName
	}

	return imageName + ":" + GetDefaultTag(rc, imageName, false, keepKernelInfo)
}

func SetDefaultManifestTagIfNoTagSet(
	rc dukkha.RenderingContext, manifestName string,
) string {
	if hasTag(manifestName) {
		return manifestName
	}

	return manifestName + ":" + GetDefaultTag(rc, manifestName, true, false)
}

func GetDefaultTag(
	rc dukkha.RenderingContext,
	name string,
	isManifest bool,
	keepKernelInfo bool,
) string {
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
			tag = strings.TrimPrefix(gitTag, "v")
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

	return tag
}

func hasTag(name string) bool {
	tagIndex := strings.LastIndexByte(name, ':')
	if tagIndex < 0 {
		return false
	}

	// has tag if there is no more slash after
	return strings.IndexByte(name[tagIndex+1:], '/') < 0
}
