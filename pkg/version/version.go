/*
Copyright 2020 The arhat.dev Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package version

import (
	"fmt"
	"runtime"
	"time"
)

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

func Arch() string {
	return arch
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
