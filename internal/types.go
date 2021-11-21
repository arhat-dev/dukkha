package dukkha_internal

// Internal interfaces made implicit to make sure called with caution
type (
	DefaultGitBranchOverrider interface {
		OverrideDefaultGitBranch(branch string)
	}

	WorkingDirOverrider interface {
		OverrideWorkingDir(cwd string)
	}

	CacheDirSetter interface {
		SetCacheDir(dir string)
	}

	VALUEGetter interface {
		VALUE() interface{}
	}

	VALUESetter interface {
		SetVALUE(v interface{})
	}
)
