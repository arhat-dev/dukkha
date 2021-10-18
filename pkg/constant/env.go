package constant

// TODO(all): Update docs/environment-variables.md when updating this file

// Environment variables for all tasks
// nolint:revive
const (
	ENV_DUKKHA_CACHE_DIR   = "DUKKHA_CACHE_DIR"
	ENV_DUKKHA_WORKING_DIR = "DUKKHA_WORKING_DIR"

	ENV_GIT_BRANCH         = "GIT_BRANCH"
	ENV_GIT_COMMIT         = "GIT_COMMIT"
	ENV_GIT_TAG            = "GIT_TAG"
	ENV_GIT_WORKTREE_CLEAN = "GIT_WORKTREE_CLEAN"
	ENV_GIT_DEFAULT_BRANCH = "GIT_DEFAULT_BRANCH"

	ENV_TIME_ZONE        = "TIME_ZONE"
	ENV_TIME_ZONE_OFFSET = "TIME_ZONE_OFFSET"
	ENV_TIME_YEAR        = "TIME_YEAR"
	ENV_TIME_MONTH       = "TIME_MONTH"
	ENV_TIME_DAY         = "TIME_DAY"
	ENV_TIME_HOUR        = "TIME_HOUR"
	ENV_TIME_MINUTE      = "TIME_MINUTE"
	ENV_TIME_SECOND      = "TIME_SECOND"

	// for linux: ID value in /etc/os-release
	ENV_HOST_OS = "HOST_OS"

	// for linux: VERSION_ID value in /etc/os-release
	ENV_HOST_OS_VERSION = "HOST_OS_VERSION"

	// value of runtime.GOOS
	ENV_HOST_KERNEL = "HOST_KERNEL"

	// value of uname -r syscall
	ENV_HOST_KERNEL_VERSION = "HOST_KERNEL_VERSION"

	// arch value
	ENV_HOST_ARCH = "HOST_ARCH"

	// triple name parts
	ENV_MATRIX_KERNEL = "MATRIX_KERNEL"
	ENV_MATRIX_ARCH   = "MATRIX_ARCH"
	ENV_MATRIX_LIBC   = "MATRIX_LIBC"
)
