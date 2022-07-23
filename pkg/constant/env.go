package constant

// TODO(all): Update docs/environment-variables.md when updating this file

type GlobalEnv int32

const (
	GlobalEnv_DUKKHA_WORKDIR GlobalEnv = iota
	GlobalEnv_DUKKHA_CACHE_DIR

	GlobalEnv_GIT_BRANCH
	GlobalEnv_GIT_COMMIT
	GlobalEnv_GIT_TAG
	GlobalEnv_GIT_WORKTREE_CLEAN
	GlobalEnv_GIT_DEFAULT_BRANCH

	GlobalEnv_TIME_ZONE
	GlobalEnv_TIME_ZONE_OFFSET
	GlobalEnv_TIME_YEAR
	GlobalEnv_TIME_MONTH
	GlobalEnv_TIME_DAY
	GlobalEnv_TIME_HOUR
	GlobalEnv_TIME_MINUTE
	GlobalEnv_TIME_SECOND

	GlobalEnv_HOST_OS
	GlobalEnv_HOST_OS_VERSION
	GlobalEnv_HOST_KERNEL
	GlobalEnv_HOST_KERNEL_VERSION
	GlobalEnv_HOST_ARCH
	GlobalEnv_HOST_ARCH_SIMPLE

	GlobalEnv_Count
)

func (f GlobalEnv) String() string {
	if f < 0 || f >= GlobalEnv_Count {
		return ""
	}

	return GlobalEnvNames[f]
}

// GetGlobalEnvIDByName return -1 when name is not a global env name
func GetGlobalEnvIDByName(name string) GlobalEnv {
	switch name {
	case EnvName_DUKKHA_CACHE_DIR:
		return GlobalEnv_DUKKHA_CACHE_DIR
	case EnvName_DUKKHA_WORKDIR:
		return GlobalEnv_DUKKHA_WORKDIR

	case EnvName_GIT_BRANCH:
		return GlobalEnv_GIT_BRANCH
	case EnvName_GIT_COMMIT:
		return GlobalEnv_GIT_COMMIT
	case EnvName_GIT_TAG:
		return GlobalEnv_GIT_TAG
	case EnvName_GIT_WORKTREE_CLEAN:
		return GlobalEnv_GIT_WORKTREE_CLEAN
	case EnvName_GIT_DEFAULT_BRANCH:
		return GlobalEnv_GIT_DEFAULT_BRANCH

	case EnvName_TIME_ZONE:
		return GlobalEnv_TIME_ZONE
	case EnvName_TIME_ZONE_OFFSET:
		return GlobalEnv_TIME_ZONE_OFFSET
	case EnvName_TIME_YEAR:
		return GlobalEnv_TIME_YEAR
	case EnvName_TIME_MONTH:
		return GlobalEnv_TIME_MONTH
	case EnvName_TIME_DAY:
		return GlobalEnv_TIME_DAY
	case EnvName_TIME_HOUR:
		return GlobalEnv_TIME_HOUR
	case EnvName_TIME_MINUTE:
		return GlobalEnv_TIME_MINUTE
	case EnvName_TIME_SECOND:
		return GlobalEnv_TIME_SECOND

	case EnvName_HOST_OS:
		return GlobalEnv_HOST_OS
	case EnvName_HOST_OS_VERSION:
		return GlobalEnv_HOST_OS_VERSION
	case EnvName_HOST_KERNEL:
		return GlobalEnv_HOST_KERNEL
	case EnvName_HOST_KERNEL_VERSION:
		return GlobalEnv_HOST_KERNEL_VERSION
	case EnvName_HOST_ARCH:
		return GlobalEnv_HOST_ARCH
	case EnvName_HOST_ARCH_SIMPLE:
		return GlobalEnv_HOST_ARCH_SIMPLE

	default:
		return -1
	}
}

var GlobalEnvNames = &[GlobalEnv_Count]string{
	GlobalEnv_DUKKHA_CACHE_DIR: EnvName_DUKKHA_CACHE_DIR,
	GlobalEnv_DUKKHA_WORKDIR:   EnvName_DUKKHA_WORKDIR,

	GlobalEnv_GIT_BRANCH:         EnvName_GIT_BRANCH,
	GlobalEnv_GIT_COMMIT:         EnvName_GIT_COMMIT,
	GlobalEnv_GIT_TAG:            EnvName_GIT_TAG,
	GlobalEnv_GIT_WORKTREE_CLEAN: EnvName_GIT_WORKTREE_CLEAN,
	GlobalEnv_GIT_DEFAULT_BRANCH: EnvName_GIT_DEFAULT_BRANCH,

	GlobalEnv_TIME_ZONE:        EnvName_TIME_ZONE,
	GlobalEnv_TIME_ZONE_OFFSET: EnvName_TIME_ZONE_OFFSET,
	GlobalEnv_TIME_YEAR:        EnvName_TIME_YEAR,
	GlobalEnv_TIME_MONTH:       EnvName_TIME_MONTH,
	GlobalEnv_TIME_DAY:         EnvName_TIME_DAY,
	GlobalEnv_TIME_HOUR:        EnvName_TIME_HOUR,
	GlobalEnv_TIME_MINUTE:      EnvName_TIME_MINUTE,
	GlobalEnv_TIME_SECOND:      EnvName_TIME_SECOND,

	GlobalEnv_HOST_OS:             EnvName_HOST_OS,
	GlobalEnv_HOST_OS_VERSION:     EnvName_HOST_OS_VERSION,
	GlobalEnv_HOST_KERNEL:         EnvName_HOST_KERNEL,
	GlobalEnv_HOST_KERNEL_VERSION: EnvName_HOST_KERNEL_VERSION,
	GlobalEnv_HOST_ARCH:           EnvName_HOST_ARCH,
	GlobalEnv_HOST_ARCH_SIMPLE:    EnvName_HOST_ARCH_SIMPLE,
}

// Environment variables for all tasks

const (
	EnvName_DUKKHA_CACHE_DIR = "DUKKHA_CACHE_DIR"
	EnvName_DUKKHA_WORKDIR   = "DUKKHA_WORKDIR"

	EnvName_GIT_BRANCH         = "GIT_BRANCH"
	EnvName_GIT_COMMIT         = "GIT_COMMIT"
	EnvName_GIT_TAG            = "GIT_TAG"
	EnvName_GIT_WORKTREE_CLEAN = "GIT_WORKTREE_CLEAN"
	EnvName_GIT_DEFAULT_BRANCH = "GIT_DEFAULT_BRANCH"

	EnvName_TIME_ZONE        = "TIME_ZONE"
	EnvName_TIME_ZONE_OFFSET = "TIME_ZONE_OFFSET"
	EnvName_TIME_YEAR        = "TIME_YEAR"
	EnvName_TIME_MONTH       = "TIME_MONTH"
	EnvName_TIME_DAY         = "TIME_DAY"
	EnvName_TIME_HOUR        = "TIME_HOUR"
	EnvName_TIME_MINUTE      = "TIME_MINUTE"
	EnvName_TIME_SECOND      = "TIME_SECOND"

	// for linux: ID value in /etc/os-release
	EnvName_HOST_OS = "HOST_OS"

	// for linux: VERSION_ID value in /etc/os-release
	EnvName_HOST_OS_VERSION = "HOST_OS_VERSION"

	// value of runtime.GOOS
	EnvName_HOST_KERNEL = "HOST_KERNEL"

	// value of uname -r syscall
	EnvName_HOST_KERNEL_VERSION = "HOST_KERNEL_VERSION"

	// arch value
	EnvName_HOST_ARCH        = "HOST_ARCH"
	EnvName_HOST_ARCH_SIMPLE = "HOST_ARCH_SIMPLE"

	// triple name parts
	EnvName_MATRIX_KERNEL      = "MATRIX_KERNEL"
	EnvName_MATRIX_ARCH        = "MATRIX_ARCH"
	EnvName_MATRIX_ARCH_SIMPLE = "MATRIX_ARCH_SIMPLE"
	EnvName_MATRIX_LIBC        = "MATRIX_LIBC"
)
