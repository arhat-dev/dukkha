package cmd

import (
	"os"
	"strings"
)

// CI environment variables
// Refs:
// 		github: https://docs.github.com/en/actions/reference/environment-variables
// 		gitlab: https://docs.gitlab.com/ee/ci/variables/predefined_variables.html

func isGithubActions() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}

func isGitlabCI() bool {
	return os.Getenv("GITLAB_CI") == "true"
}

// GitCommitFromCI find git commit sha from ci env
func GitCommitFromCI() string {
	switch {
	case isGithubActions():
		return strings.TrimSpace(os.Getenv("GITHUB_SHA"))
	case isGitlabCI():
		return strings.TrimSpace(os.Getenv("CI_COMMIT_SHA"))
	default:
		return ""
	}
}

// GitBranchFromCI find git branch name from ci env
func GitBranchFromCI() string {
	switch {
	case isGithubActions():
		ghRef := strings.TrimSpace(os.Getenv("GITHUB_REF"))
		if len(ghRef) == 0 {
			ghRef = strings.TrimSpace(os.Getenv("GITHUB_HEAD_REF"))
		}

		switch {
		case strings.HasPrefix(ghRef, "refs/heads/"):
			return strings.TrimPrefix(ghRef, "refs/heads/")
		default:
			return ""
		}
	case isGitlabCI():
		return strings.TrimSpace(os.Getenv("CI_COMMIT_BRANCH"))
	default:
		return ""
	}
}

// GitTagFromCI find git tag name from ci env
func GitTagFromCI() string {
	switch {
	case isGithubActions():
		ghRef := strings.TrimSpace(os.Getenv("GITHUB_REF"))
		if len(ghRef) == 0 {
			ghRef = strings.TrimSpace(os.Getenv("GITHUB_HEAD_REF"))
		}

		switch {
		case strings.HasPrefix(ghRef, "refs/tags/"):
			return strings.TrimPrefix(ghRef, "refs/tags/")
		default:
			return ""
		}
	case isGitlabCI():
		return strings.TrimSpace(os.Getenv("CI_COMMIT_TAG"))
	default:
		return ""
	}
}
