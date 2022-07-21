package cmd

import (
	"bufio"
	"context"
	"os"
	"strconv"
	"strings"

	"arhat.dev/tlang"

	"arhat.dev/dukkha/pkg/constant"
)

// GitBranch find git branch name of working dir wd
func GitBranch(ctx context.Context, wd string) *tlang.LazyValue[string] {
	return newLazyExecVal(
		ctx,
		wd,
		[]string{
			"git", "symbolic-ref", "--short", "-q", "HEAD",
		},
		GitBranchFromCI,
		func(s string) string {
			s = strings.TrimSpace(s)
			if len(s) == 0 {
				return GitBranchFromCI()
			}

			return s
		},
	)
}

// GitCommit get git commit sha of working dir wd
func GitCommit(ctx context.Context, dir string) *tlang.LazyValue[string] {
	return newLazyExecVal(
		ctx,
		dir,
		[]string{
			"git", "rev-parse", "HEAD",
		},
		GitCommitFromCI,
		func(s string) string {
			s = strings.TrimSpace(s)
			if len(s) == 0 {
				return GitCommitFromCI()
			}

			return s
		},
	)
}

// GitTag find current git tag name of working dir wd
func GitTag(ctx context.Context, dir string) *tlang.LazyValue[string] {
	gitTagList := newLazyExecVal(
		ctx,
		dir,
		[]string{
			// get git tags in latest first order
			"git", "tag", "--list", "--sort", "-version:refname",
		},
		func() string { return "" },
		func(s string) string { return s },
	)

	return newLazyExecVal(
		ctx,
		dir,
		[]string{
			"git", "describe", "--tags",
		},
		GitTagFromCI,
		func(result string) string {
			result = strings.TrimSpace(strings.SplitN(result, " ", 2)[0])
			if len(result) == 0 {
				return GitTagFromCI()
			}

			s := bufio.NewScanner(strings.NewReader(gitTagList.GetLazyValue()))
			s.Split(bufio.ScanLines)
			for s.Scan() {
				if strings.Contains(s.Text(), result) {
					return result
				}
			}

			return ""
		},
	)
}

// GitWorkTreeClean check whether current dir contains no new file
// or uncommitted changes
func GitWorkTreeClean(ctx context.Context, dir string) *tlang.LazyValue[string] {
	gitDiffIndex := newLazyExecVal(
		ctx,
		dir,
		[]string{"git", "diff-index", "--quiet", "HEAD"},
		func() string { return "false" },
		func(s string) string { return "true" },
	)

	return newLazyExecVal(
		ctx,
		dir,
		[]string{
			"git", "clean", "--dry-run",
		},
		func() string { return "false" },
		func(s string) string {
			return strconv.FormatBool(
				// no output means no new files
				len(strings.TrimSpace(s)) == 0 &&
					// git diff index exit 0 means no files modified
					gitDiffIndex.GetLazyValue() == "true",
			)
		},
	)
}

// GitDefaultBranch find default git branch of dir
func GitDefaultBranch(ctx context.Context, dir string) *tlang.LazyValue[string] {
	return newLazyExecVal(
		ctx,
		dir,
		[]string{
			"git", "remote", "show", "origin",
		},
		func() string { return os.Getenv(constant.ENV_GIT_DEFAULT_BRANCH) },
		func(result string) string {
			s := bufio.NewScanner(strings.NewReader(result))
			s.Split(bufio.ScanLines)
			for s.Scan() {
				line := s.Text()
				const prefix = "HEAD branch: "
				if idx := strings.Index(line, prefix); idx != -1 {
					return line[idx+len(prefix):]
				}
			}

			return os.Getenv(constant.ENV_GIT_DEFAULT_BRANCH)
		},
	)
}
