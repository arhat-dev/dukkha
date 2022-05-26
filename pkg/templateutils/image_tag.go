package templateutils

import (
	"strings"

	"github.com/huandu/xstrings"

	"arhat.dev/dukkha/pkg/dukkha"
)

func GetImageTag(
	rc dukkha.RenderingContext, imageName string, keepKernelInfo bool,
) string {
	return GetDefaultTag(rc, imageName, false, keepKernelInfo)
}

func GetManifestTag(
	rc dukkha.RenderingContext, manifestName string,
) string {
	return GetDefaultTag(rc, manifestName, true, false)
}

func GetFullImageName_UseDefault_IfIfNoTagSet(
	rc dukkha.RenderingContext, imageName string, keepKernelInfo bool,
) string {
	if hasTag(imageName) {
		return imageName
	}

	return imageName + ":" + GetDefaultTag(rc, imageName, false, keepKernelInfo)
}

func GetFullManifestName_UseDefault_IfNoTagSet(
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
				// image tag, <branch>-<commit> or <commit> when branch missing
				if len(tag) == 0 {
					// no branch info (can happen in github actions)
					// TODO: add test for this case
					tag = rc.GitCommit()
				} else {
					tag += "-" + rc.GitCommit()
				}
			} else {
				// manifest tag, <branch> or <commit> when branch missing
				if len(tag) == 0 {
					// no branch info (can happen in github actions)
					// TODO: add test for this case
					tag = rc.GitCommit()
				}
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

	// has tag if there is no more slash after (to handle port number in host)
	return strings.IndexByte(name[tagIndex+1:], '/') < 0
}
