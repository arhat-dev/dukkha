package dukkha

import (
	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/utils"
)

// implementation for internal interfaces in arhat.dev/dukkha/internal

// SetVALUE for transform renderer
func (c *contextRendering) SetVALUE(value interface{}) { c._VALUE = value }

// VALUE for transform renderer
func (c *contextRendering) VALUE() interface{} { return c._VALUE }

// SetCacheDir set env DUKKHA_CACHE_DIR
//
// should not be exposed by any interface type in this package
func (c *envValues) SetCacheDir(dir string) {
	c.globalEnv[constant.ENV_DUKKHA_CACHE_DIR] = utils.ImmediateString(dir)
}

// OverrideDefaultGitBranch set env GIT_DEFAULT_BRANCH
//
// should not be exposed by any interface type in this package
func (c *envValues) OverrideDefaultGitBranch(branch string) {
	c.globalEnv[constant.ENV_GIT_DEFAULT_BRANCH] = utils.ImmediateString(branch)
}

// OverrideWorkDir set env DUKKHA_WORKDIR to cwd
//
// should not be exposed by any interface type in this package
func (c *envValues) OverrideWorkDir(cwd string) {
	c.globalEnv[constant.ENV_DUKKHA_WORKDIR] = utils.ImmediateString(cwd)
}
