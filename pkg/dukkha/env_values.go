package dukkha

import (
	"strings"
	"sync"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/matrix"
	"arhat.dev/dukkha/pkg/utils"
)

// This file describes runtime values derived from env

type GlobalEnvValues interface {
	WorkDir() string
	CacheDir() string

	GitBranch() string
	GitWorkTreeClean() bool
	GitTag() string
	GitDefaultBranch() string
	GitCommit() string
	GitValues() map[string]utils.LazyValue

	HostKernel() string
	HostKernelVersion() string
	HostArch() string
	HostOS() string
	HostOSVersion() string
	HostValues() map[string]utils.LazyValue
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

func newEnvValues(globalEnv map[string]utils.LazyValue) *envValues {
	ret := &envValues{
		matrixFilter: nil,

		globalEnv: globalEnv,

		env: make(map[string]utils.LazyValue),
		mu:  new(sync.RWMutex),
	}

	return ret
}

var _ EnvValues = (*envValues)(nil)

type envValues struct {
	matrixFilter *matrix.Filter

	globalEnv map[string]utils.LazyValue

	env map[string]utils.LazyValue
	mu  *sync.RWMutex
}

func (c *envValues) clone() *envValues {
	newValues := &envValues{
		matrixFilter: nil,
		globalEnv:    c.globalEnv,
		env:          make(map[string]utils.LazyValue),
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
	return getValueOrDefault(c.env[constant.ENV_MATRIX_ARCH])
}

func (c *envValues) MatrixKernel() string {
	return getValueOrDefault(c.env[constant.ENV_MATRIX_KERNEL])
}

func (c *envValues) MatrixLibc() string {
	return getValueOrDefault(c.env[constant.ENV_MATRIX_LIBC])
}

func (c *envValues) AddEnv(override bool, entries ...*EnvEntry) {
	for _, e := range entries {
		if _, ok := c.env[e.Name]; ok && !override {
			continue
		}

		c.env[e.Name] = utils.ImmediateString(e.Value)
	}
}

func (c *envValues) AddListEnv(env ...string) {
	for _, entry := range env {
		parts := strings.SplitN(entry, "=", 2)
		key, value := parts[0], ""
		if len(parts) == 2 {
			value = parts[1]
		}

		c.env[key] = utils.ImmediateString(value)
	}
}

func (c *envValues) WorkDir() string {
	return getValueOrDefault(c.globalEnv[constant.ENV_DUKKHA_WORKDIR])
}

func (c *envValues) CacheDir() string {
	return getValueOrDefault(c.globalEnv[constant.ENV_DUKKHA_CACHE_DIR])
}

func (c *envValues) GitValues() map[string]utils.LazyValue {
	return map[string]utils.LazyValue{
		"branch":         c.globalEnv[constant.ENV_GIT_BRANCH],
		"worktree_clean": c.globalEnv[constant.ENV_GIT_WORKTREE_CLEAN],
		"tag":            c.globalEnv[constant.ENV_GIT_TAG],
		"default_branch": c.globalEnv[constant.ENV_GIT_DEFAULT_BRANCH],
		"commit":         c.globalEnv[constant.ENV_GIT_COMMIT],
	}
}

func (c *envValues) GitBranch() string {
	return getValueOrDefault(c.globalEnv[constant.ENV_GIT_BRANCH])
}

func (c *envValues) GitWorkTreeClean() bool {
	return getValueOrDefault(c.globalEnv[constant.ENV_GIT_WORKTREE_CLEAN]) == "true"
}

func (c *envValues) GitTag() string {
	return getValueOrDefault(c.globalEnv[constant.ENV_GIT_TAG])
}

func (c *envValues) GitDefaultBranch() string {
	return getValueOrDefault(c.globalEnv[constant.ENV_GIT_DEFAULT_BRANCH])
}

func (c *envValues) GitCommit() string {
	return getValueOrDefault(c.globalEnv[constant.ENV_GIT_COMMIT])
}

func (c *envValues) HostValues() map[string]utils.LazyValue {
	return map[string]utils.LazyValue{
		"arch":           c.globalEnv[constant.ENV_HOST_ARCH],
		"arch_simple":    c.globalEnv[constant.ENV_HOST_ARCH_SIMPLE],
		"kernel":         c.globalEnv[constant.ENV_HOST_KERNEL],
		"kernel_version": c.globalEnv[constant.ENV_HOST_KERNEL_VERSION],
		"os":             c.globalEnv[constant.ENV_HOST_OS],
		"os_version":     c.globalEnv[constant.ENV_HOST_OS_VERSION],
	}
}

func (c *envValues) HostArch() string {
	return getValueOrDefault(c.globalEnv[constant.ENV_HOST_ARCH])
}

func (c *envValues) HostKernel() string {
	return getValueOrDefault(c.globalEnv[constant.ENV_HOST_KERNEL])
}

func (c *envValues) HostKernelVersion() string {
	return getValueOrDefault(c.globalEnv[constant.ENV_HOST_KERNEL_VERSION])
}

func (c *envValues) HostOS() string {
	return getValueOrDefault(c.globalEnv[constant.ENV_HOST_OS])
}

func (c *envValues) HostOSVersion() string {
	return getValueOrDefault(c.globalEnv[constant.ENV_HOST_OS_VERSION])
}

func getValueOrDefault(v utils.LazyValue) string {
	if v == nil {
		return ""
	}

	return v.Get()
}
