package dukkha

import (
	"arhat.dev/dukkha/pkg/constant"
)

type ImmutableValues interface {
	WorkingDir() string
	CacheDir() string

	GitBranch() string
	GitWorkTreeClean() bool
	GitTag() string
	GitDefaultBranch() string
	GitCommit() string

	HostKernel() string
	HostKernelVersion() string
	HostArch() string
	HostOS() string
	HostOSVersion() string
}

func newContextImmutableValues(globalEnv map[string]string) *immutableValues {
	return &immutableValues{
		globalEnv: globalEnv,
	}
}

var _ ImmutableValues = (*immutableValues)(nil)

type immutableValues struct {
	// pre-defined environment variables and bootstrap env
	globalEnv map[string]string
}

func (c *immutableValues) WorkingDir() string {
	return c.globalEnv[constant.ENV_DUKKHA_WORKING_DIR]
}

func (c *immutableValues) CacheDir() string {
	return c.globalEnv[constant.ENV_DUKKHA_CACHE_DIR]
}

func (c *immutableValues) GitBranch() string {
	return c.globalEnv[constant.ENV_GIT_BRANCH]
}

func (c *immutableValues) GitWorkTreeClean() bool {
	return c.globalEnv[constant.ENV_GIT_WORKTREE_CLEAN] == "true"
}

func (c *immutableValues) GitTag() string {
	return c.globalEnv[constant.ENV_GIT_TAG]
}

func (c *immutableValues) GitDefaultBranch() string {
	return c.globalEnv[constant.ENV_GIT_DEFAULT_BRANCH]
}

func (c *immutableValues) GitCommit() string {
	return c.globalEnv[constant.ENV_GIT_COMMIT]
}

func (c *immutableValues) HostArch() string {
	return c.globalEnv[constant.ENV_HOST_ARCH]
}

func (c *immutableValues) HostKernel() string {
	return c.globalEnv[constant.ENV_HOST_KERNEL]
}

func (c *immutableValues) HostKernelVersion() string {
	return c.globalEnv[constant.ENV_HOST_KERNEL_VERSION]
}

func (c *immutableValues) HostOS() string {
	return c.globalEnv[constant.ENV_HOST_OS]
}

func (c *immutableValues) HostOSVersion() string {
	return c.globalEnv[constant.ENV_HOST_OS_VERSION]
}
