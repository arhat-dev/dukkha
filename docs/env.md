# Environment Variables

Available environment variables when running `dukkha`

__NOTE:__ This doc should be synced with [pkg/cmd/env.go](../pkg/cmd/env.go), [pkg/constant/env.go](../pkg/constant/env.go) and [pkg/dukkha/rendering.go](../pkg/dukkha/rendering.go)

## Usage

For `env` renderer: Use like unix shell env (e.g. `${SOME_ENV}`)
For `template` renderer: Available under `env` object (e.g. `{{ env.SOME_ENV }}`)

## `dukkha` Runtime Information

- `DUKKHA_WORKING_DIR`
  - Description: The absolute directory path in which you invoked `dukkha`
  - Default Value: `$(pwd)` value in the directory you run dukkha
  - Customization: Not Supported
  - Potential Use Cases:
    - Mount proper working dir for containerized tools when `chdir` used in your task

- `DUKKHA_CACHE_DIR`
  - Description: The absolute path of the cache directory used for task intermediate output caching
  - Default Value: `${DUKKHA_WORKING_DIR}/.dukkha/cache`
  - Customization: Set `bootstrap.cache_dir` in your config file

## `git` Repo Information

- `GIT_BRANCH`
  - Description: Current working branch name
  - Default Value: `$(git symbolic-ref --short -q HEAD)`, fallback to value from CI system env (github: `GITHUB_REF`, `GITHUB_HEAD_REF`, gitlab: `CI_COMMIT_BRANCH`)
    - Example Values: `master`, `test/foo`
  - Customization: Not Supported

- `GIT_COMMIT`
  - Description: Current commit sha
  - Default Value: `$(git rev-parse HEAD)`, fallback to value from CI system env (github: `GITHUB_SHA`, gitlab: `CI_COMMIT_SHA`)
    - Example Values: `46a0cbe436971d66e79f4d03745ce9f61acb282f`
  - Customization: Not Supported

- `GIT_TAG`
  - Description: Current git tag
  - Default Value: First value in `$(git describe --tags)`, if it's not listed in `$(git tag --list --sort -version:refname)`, fallback to value provided by CI systems (github: `GITHUB_REF`, `GITHUB_HEAD_REF`, gitlab: `CI_COMMIT_TAG`)
    - Example Values: `v0.0.1`, `1.0.2`
  - Customization: Not Supported

- `GIT_WORKTREE_CLEAN`
  - Description: Indicate whether there is file not committed in current working tree
  - Default Value: `true` if `git clean --dry-run` writes no output, otherwise `false`
    - Example Values: `true` or `false`
  - Customization: Not Supported

- `GIT_DEFAULT_BRANCH`
  - Description: Default remote branch of current repo
  - Default Value: Value of `HEAD branch` from output of `git remote show origin`
    - Example Values: `master`, `main`
  - Customization:
    - Set `GIT_DEFAULT_BRANCH` env before running dukkha to provide default value when `git remote show origin` doesn't work properly
    - Set `global.default_git_branch` to force override.

__NOTE:__ These git related values are evaluated at the first time when used.

## Time Information

All time related values are based on local time

- `TIME_ZONE`
  - Description: Name of local timezone
  - Default Value: golang `time.Now().Local().Zone()` value
  - Customization: Not Supported

- `TIME_ZONE_OFFSET`
  - Description: Local timezone offset to UTC
  - Default Value: golang `time.Now().Local().Zone()` value
  - Customization: Not Supported

- `TIME_YEAR`, `TIME_MONTH`, `TIME_DAY`, `TIME_HOUR`, `TIME_MINUTE`, `TIME_SECOND`
  - Description: Current year, month, day, hour, minute, second number when invoking `dukkha`
  - Default Value: golang `time.Now().Local()` values
  - Customization: Not Supported

## Host System Information

- `HOST_KERNEL`
  - Description: Kernel name of the host system running `dukkha`
  - Default Value: value of golang `runtime.GOOS`
    - Example Values: `linux`, `darwin`
  - Customization: Not Supported

- `HOST_KERNEL_VERSION`
  - Description: Kernel version of the host system running `dukkha`
  - Default Value: `$(uname -r)`
    - Example Values: `5.12.12-300.fc34.x86_64` (on fedora 34), `20.5.0` (on macOS 11.4)
  - Customization: Not Supported

- `HOST_OS`
  - Description: OS name of the host system running `dukkha`
  - Default Value:
    - linux: value of `ID` field in `/etc/os-release`
      - Example Values: `ubuntu`, `debian`, `fedora`
    - other: value of golang `runtime.GOOS`
  - Customization: Not Supported

- `HOST_OS_VERSION`
  - Description: OS version of the host system running `dukkha`
  - Default Value:
    - linux: value of `VERSION_ID` field in `/etc/os-release`
      - Example Values: `34` (on fedora 34), `20.04` (on ubuntu 20.04)
  - Customization: Not Supported

- `HOST_ARCH`
  - Description: CPU arch value of the host system running `dukkha`
  - Default Value: `dukkha` defined mapped value of `$(uname -m)`
    - Example Values: see [`System Arch` section in docs/constants.md](./constants.md#system-arch)
  - Customization: Not Supported

## Task Execution Information

__NOTE:__ Environment variables in this section are only available for your tasks and tools.

- `MATRIX_<upper-case-matrix-spec-key>`
  - Description: Matrix value
  - Example Names: `MATRIX_KERNEL` for `matrix.kernel`, `MATRIX_FOO_DATA` for `matrix.foo_data`
