package dukkha

import (
	"strings"
	"sync"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/matrix"
)

// This file describes runtime values derived from env

type GlobalEnvValues interface {
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

type EnvValues interface {
	GlobalEnvValues

	SetMatrixFilter(*matrix.Filter)
	MatrixFilter() *matrix.Filter

	MatrixArch() string
	MatrixKernel() string
	MatrixLibc() string

	AddEnv(override bool, env ...*EnvEntry)
	AddListEnv(env ...string)
}

func newEnvValues(globalEnv map[string]string) *envValues {
	ret := &envValues{
		matrixFilter: nil,

		globalEnv: globalEnv,

		env: make(map[string]string),
		mu:  new(sync.RWMutex),
	}

	return ret
}

var _ EnvValues = (*envValues)(nil)

type envValues struct {
	matrixFilter *matrix.Filter

	globalEnv map[string]string

	env map[string]string
	mu  *sync.RWMutex
}

func (c *envValues) clone() *envValues {
	newValues := &envValues{
		matrixFilter: nil,
		globalEnv:    c.globalEnv,
		env:          make(map[string]string),
		mu:           new(sync.RWMutex),
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.matrixFilter != nil {
		newValues.matrixFilter = c.matrixFilter.Clone()
	}

	for k, v := range c.env {
		newValues.env[k] = v
	}

	return newValues
}

func (c *envValues) SetMatrixFilter(f *matrix.Filter) {
	c.matrixFilter = f
}

func (c *envValues) MatrixFilter() *matrix.Filter {
	return c.matrixFilter
}

func (c *envValues) MatrixArch() string {
	return c.env[constant.ENV_MATRIX_ARCH]
}

func (c *envValues) MatrixKernel() string {
	return c.env[constant.ENV_MATRIX_KERNEL]
}

func (c *envValues) MatrixLibc() string {
	return c.env[constant.ENV_MATRIX_LIBC]
}

func (c *envValues) AddEnv(override bool, entries ...*EnvEntry) {
	for _, e := range entries {
		if _, ok := c.env[e.Name]; ok && !override {
			continue
		}

		c.env[e.Name] = e.Value
	}
}

func (c *envValues) AddListEnv(env ...string) {
	for _, entry := range env {
		parts := strings.SplitN(entry, "=", 2)
		key, value := parts[0], ""
		if len(parts) == 2 {
			value = parts[1]
		}

		c.env[key] = value
	}
}

func (c *envValues) SetCacheDir(dir string) {
	c.globalEnv[constant.ENV_DUKKHA_CACHE_DIR] = dir
}

func (c *envValues) OverrideDefaultGitBranch(branch string) {
	c.globalEnv[constant.ENV_GIT_DEFAULT_BRANCH] = branch
}

// OverrideWorkingDir set DUKKHA_WORKING_DIR to cwd
// should not be exposed by any interface type in this package
func (c *envValues) OverrideWorkingDir(cwd string) {
	c.globalEnv[constant.ENV_DUKKHA_WORKING_DIR] = cwd
}

func (c *envValues) WorkingDir() string {
	return c.globalEnv[constant.ENV_DUKKHA_WORKING_DIR]
}

func (c *envValues) CacheDir() string {
	return c.globalEnv[constant.ENV_DUKKHA_CACHE_DIR]
}

func (c *envValues) GitBranch() string {
	return c.globalEnv[constant.ENV_GIT_BRANCH]
}

func (c *envValues) GitWorkTreeClean() bool {
	return c.globalEnv[constant.ENV_GIT_WORKTREE_CLEAN] == "true"
}

func (c *envValues) GitTag() string {
	return c.globalEnv[constant.ENV_GIT_TAG]
}

func (c *envValues) GitDefaultBranch() string {
	return c.globalEnv[constant.ENV_GIT_DEFAULT_BRANCH]
}

func (c *envValues) GitCommit() string {
	return c.globalEnv[constant.ENV_GIT_COMMIT]
}

func (c *envValues) HostArch() string {
	return c.globalEnv[constant.ENV_HOST_ARCH]
}

func (c *envValues) HostKernel() string {
	return c.globalEnv[constant.ENV_HOST_KERNEL]
}

func (c *envValues) HostKernelVersion() string {
	return c.globalEnv[constant.ENV_HOST_KERNEL_VERSION]
}

func (c *envValues) HostOS() string {
	return c.globalEnv[constant.ENV_HOST_OS]
}

func (c *envValues) HostOSVersion() string {
	return c.globalEnv[constant.ENV_HOST_OS_VERSION]
}
