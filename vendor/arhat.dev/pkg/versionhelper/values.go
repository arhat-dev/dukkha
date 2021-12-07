package versionhelper

import (
	"fmt"
	"runtime"
	"time"

	"arhat.dev/pkg/archconst"
)

// values should be set at build time using ldflags
// -X arhat.dev/pkg/versionhelper.{branch,commit, ...}=...
var (
	branch, commit, tag, arch string

	worktreeClean string
	buildTime     string

	goCompilerPlatform string
)

var version string

func init() {
	version = fmt.Sprintf(`branch: %s
commit: %s
tag: %s
arch: %s
goVersion: %s
buildTime: %s
worktreeClean: %s
goCompilerPlatform: %s
`,
		Branch(),
		Commit(),
		Tag(),
		Arch(),
		GoVersion(),
		buildTime,
		worktreeClean,
		GoCompilerPlatform(),
	)
}

func Version() string {
	return version
}

// Branch name of the source code
func Branch() string {
	return branch
}

// Commit hash of the source code
func Commit() string {
	return commit
}

// Tag the tag name of the source code
func Tag() string {
	return tag
}

// Arch returns cpu arch with default micro arch applied if missing
func Arch() string {
	switch arch {
	case archconst.ARCH_AMD64:
		return archconst.ARCH_AMD64_V1
	case archconst.ARCH_PPC64:
		return archconst.ARCH_PPC64_V8
	case archconst.ARCH_PPC64_LE:
		return archconst.ARCH_PPC64_V8_LE
	default:
		return arch
	}
}

func GoVersion() string {
	return runtime.Version()
}

func BuildTime() time.Time {
	ret, err := time.Parse(time.RFC3339, buildTime)
	if err != nil {
		return time.Time{}
	}

	return ret
}

func WorktreeClean() bool {
	return worktreeClean == "true"
}

func GoCompilerPlatform() string {
	return goCompilerPlatform
}
