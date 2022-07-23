package dukkha_internal

// Internal interfaces made implicit to make sure called with caution
type (
	DefaultGitBranchOverrider interface {
		OverrideDefaultGitBranch(branch string)
	}

	WorkDirOverrider interface {
		OverrideWorkDir(cwd string)
	}

	VALUEGetter interface {
		VALUE() interface{}
	}

	VALUESetter interface {
		SetVALUE(v interface{})
	}
)
