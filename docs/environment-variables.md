# Environment Variables

Available environment variables when running `dukkha run`

__NOTE:__ This doc should be synced with [pkg/cmd/env.go](../pkg/cmd/env.go) and [pkg/constant/env.go](../pkg/constant/env.go)

## `dukkha` Runtime Information

- `DUKKHA_WORKING_DIR`
  - Description: The absolute directory path in which you invoked `dukkha`
  - Default Value: `$(pwd)`
  - Customization: Not Supported
  - Potential Use Cases:
    - To mount proper working dir for containerized tools when `chdir` is used in your task

- `DUKKHA_CACHE_DIR`
  - Description: The absolute path of the cache directory used for shell script caching and temporary data storage for task execution
  - Default Value: `$(pwd)/.dukkha/cache`
  - Customization: Set `bootstrap.cache_dir` in your config file

## `git` Repo Information

- `GIT_BRANCH`
  - Description: Active branch name when invoking `dukkha`
  - Default Value: `$(git symbolic-ref --short -q HEAD)`
    - Example Values: `master`, `test/foo`
  - Customization: Not Supported

- `GIT_COMMIT`
  - Description: Current commit sha when invoking `dukkha`
    - Example Values: `46a0cbe436971d66e79f4d03745ce9f61acb282f`
  - Default Value: `$(git rev-parse HEAD)`
  - Customization: Not Supported

- `GIT_TAG`
  - Description: Current git tag value when invoking `dukkha`
  - Default Value: first value in `$(git describe --tags)`
    - Example Values: `v0.0.1`, `1.0.2`
  - Customization: Not Supported

- `GIT_WORKTREE_CLEAN`
  - Description: Indicate whether there is file not committed when invoking `dukkha`
  - Default Value: `true` if `git diff-index --quiet HEAD` exited with 0, otherwise `false`
    - Example Values: `true` or `false`
  - Customization: Not Supported

- `GIT_DEFAULT_BRANCH`
  - Description: Default remote branch of this repo
  - Default Value: `$(git symbolic-ref refs/remotes/origin/HEAD)` with prefix `refs/remotes/origin/` trimed
    - Example Values: `master`, `main`
  - Customization: Not Supported

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
