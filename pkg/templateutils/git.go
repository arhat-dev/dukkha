package templateutils

import "arhat.dev/dukkha/pkg/dukkha"

func createGitNS(rc dukkha.RenderingContext) *gitNS {
	return &gitNS{rc: rc}
}

type gitNS struct {
	rc dukkha.RenderingContext
}

// Branch get GIT_BRANCH
func (ns *gitNS) Branch() string {
	return ns.rc.GitBranch()
}

// Commit get GIT_COMMIT
func (ns *gitNS) Commit() string {
	return ns.rc.GitCommit()
}

// Tag get GIT_TAG
func (ns *gitNS) Tag() string {
	return ns.rc.GitTag()
}

// DefaultBranch get GIT_DEFAULT_BRANCH
func (ns *gitNS) DefaultBranch() string {
	return ns.rc.GitDefaultBranch()
}

// WorktreeClean get GIT_WORKTREE_CLEAN as bool value
func (ns *gitNS) WorktreeClean() bool {
	return ns.rc.GitWorkTreeClean()
}
