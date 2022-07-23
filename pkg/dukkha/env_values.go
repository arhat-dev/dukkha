package dukkha

import (
	"strings"
	"sync"

	"arhat.dev/tlang"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/matrix"
)

// This file describes runtime values derived from env

type GlobalEnvSet [constant.GlobalEnv_Count]tlang.LazyValueType[string]

func (s *GlobalEnvSet) Get(name string) (ret tlang.LazyValueType[string], ok bool) {
	id := constant.GetGlobalEnvIDByName(name)
	if id == -1 {
		return
	}

	return s[id], true
}

type GlobalEnvValues interface {
	WorkDir() string
	CacheDir() string

	GitBranch() string
	GitWorkTreeClean() bool
	GitTag() string
	GitDefaultBranch() string
	GitCommit() string
	GitValues() map[string]tlang.LazyValueType[string]

	HostKernel() string
	HostKernelVersion() string
	HostArch() string
	HostOS() string
	HostOSVersion() string
	HostValues() map[string]tlang.LazyValueType[string]
}

type EnvValues interface {
	GlobalEnvValues

	SetMatrixFilter(matrix.Filter)
	MatrixFilter() matrix.Filter

	MatrixArch() string
	MatrixKernel() string
	MatrixLibc() string

	AddEnv(override bool, env ...*EnvEntry)
	AddListEnv(env ...string)
}

func newEnvValues(globalEnv *GlobalEnvSet) envValues {
	return envValues{
		globalEnv: globalEnv,

		gitValues: &tlang.LazyValue[map[string]tlang.LazyValueType[string]]{
			Create: func() map[string]tlang.LazyValueType[string] {
				return map[string]tlang.LazyValueType[string]{
					"branch":         globalEnv[constant.GlobalEnv_GIT_BRANCH],
					"worktree_clean": globalEnv[constant.GlobalEnv_GIT_WORKTREE_CLEAN],
					"tag":            globalEnv[constant.GlobalEnv_GIT_TAG],
					"default_branch": globalEnv[constant.GlobalEnv_GIT_DEFAULT_BRANCH],
					"commit":         globalEnv[constant.GlobalEnv_GIT_COMMIT],
				}
			},
		},

		hostValues: &tlang.LazyValue[map[string]tlang.LazyValueType[string]]{
			Create: func() map[string]tlang.LazyValueType[string] {
				return map[string]tlang.LazyValueType[string]{
					"arch":           globalEnv[constant.GlobalEnv_HOST_ARCH],
					"arch_simple":    globalEnv[constant.GlobalEnv_HOST_ARCH_SIMPLE],
					"kernel":         globalEnv[constant.GlobalEnv_HOST_KERNEL],
					"kernel_version": globalEnv[constant.GlobalEnv_HOST_KERNEL_VERSION],
					"os":             globalEnv[constant.GlobalEnv_HOST_OS],
					"os_version":     globalEnv[constant.GlobalEnv_HOST_OS_VERSION],
				}
			},
		},

		env: make(map[string]tlang.LazyValueType[string]),
		mu:  new(sync.RWMutex),
	}
}

var _ EnvValues = (*envValues)(nil)

type envValues struct {
	matrixFilter matrix.Filter

	globalEnv *GlobalEnvSet

	gitValues  *tlang.LazyValue[map[string]tlang.LazyValueType[string]]
	hostValues *tlang.LazyValue[map[string]tlang.LazyValueType[string]]

	env map[string]tlang.LazyValueType[string]
	mu  *sync.RWMutex
}

func (c *envValues) clone() envValues {
	newValues := envValues{
		globalEnv: c.globalEnv,

		gitValues:  c.gitValues,
		hostValues: c.hostValues,

		env: make(map[string]tlang.LazyValueType[string]),
		mu:  new(sync.RWMutex),
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.matrixFilter.Empty() {
		newValues.matrixFilter = c.matrixFilter.Clone()
	}

	for k, v := range c.env {
		newValues.env[k] = v
	}

	return newValues
}

func (c *envValues) SetMatrixFilter(f matrix.Filter) {
	c.matrixFilter = f
}

func (c *envValues) MatrixFilter() matrix.Filter {
	return c.matrixFilter
}

func (c *envValues) MatrixArch() string {
	return getValueOrEmpty(c.env[constant.EnvName_MATRIX_ARCH])
}

func (c *envValues) MatrixKernel() string {
	return getValueOrEmpty(c.env[constant.EnvName_MATRIX_KERNEL])
}

func (c *envValues) MatrixLibc() string {
	return getValueOrEmpty(c.env[constant.EnvName_MATRIX_LIBC])
}

func getValueOrEmpty(v tlang.LazyValueType[string]) string {
	if v == nil {
		return ""
	}

	return v.GetLazyValue()
}

func (c *envValues) AddEnv(override bool, entries ...*EnvEntry) {
	for _, e := range entries {
		if _, ok := c.env[e.Name]; ok && !override {
			continue
		}

		c.env[e.Name] = tlang.ImmediateString(e.Value)
	}
}

func (c *envValues) AddListEnv(env ...string) {
	for _, entry := range env {
		parts := strings.SplitN(entry, "=", 2)
		key, value := parts[0], ""
		if len(parts) == 2 {
			value = parts[1]
		}

		c.env[key] = tlang.ImmediateString(value)
	}
}

func (c *envValues) WorkDir() string {
	return c.globalEnv[constant.GlobalEnv_DUKKHA_WORKDIR].GetLazyValue()
}

func (c *envValues) CacheDir() string {
	return c.globalEnv[constant.GlobalEnv_DUKKHA_CACHE_DIR].GetLazyValue()
}

func (c *envValues) GitValues() map[string]tlang.LazyValueType[string] {
	return c.gitValues.GetLazyValue()
}

func (c *envValues) GitBranch() string {
	return c.globalEnv[constant.GlobalEnv_GIT_BRANCH].GetLazyValue()
}

func (c *envValues) GitWorkTreeClean() bool {
	return c.globalEnv[constant.GlobalEnv_GIT_WORKTREE_CLEAN].GetLazyValue() == "true"
}

func (c *envValues) GitTag() string {
	return c.globalEnv[constant.GlobalEnv_GIT_TAG].GetLazyValue()
}

func (c *envValues) GitDefaultBranch() string {
	return c.globalEnv[constant.GlobalEnv_GIT_DEFAULT_BRANCH].GetLazyValue()
}

func (c *envValues) GitCommit() string {
	return c.globalEnv[constant.GlobalEnv_GIT_COMMIT].GetLazyValue()
}

func (c *envValues) HostValues() map[string]tlang.LazyValueType[string] {
	return c.hostValues.GetLazyValue()
}

func (c *envValues) HostArch() string {
	return c.globalEnv[constant.GlobalEnv_HOST_ARCH].GetLazyValue()
}

func (c *envValues) HostKernel() string {
	return c.globalEnv[constant.GlobalEnv_HOST_KERNEL].GetLazyValue()
}

func (c *envValues) HostKernelVersion() string {
	return c.globalEnv[constant.GlobalEnv_HOST_KERNEL_VERSION].GetLazyValue()
}

func (c *envValues) HostOS() string {
	return c.globalEnv[constant.GlobalEnv_HOST_OS].GetLazyValue()
}

func (c *envValues) HostOSVersion() string {
	return c.globalEnv[constant.GlobalEnv_HOST_OS_VERSION].GetLazyValue()
}
