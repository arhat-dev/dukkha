package cmd

import (
	"bufio"
	"context"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"arhat.dev/pkg/exechelper"
	"arhat.dev/tlang"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sysinfo"
)

// TODO(all): Update docs/environment-variables.md when updating this file

func createGlobalEnv(ctx context.Context, cwd string) *dukkha.GlobalEnvSet {
	now := time.Now().Local()
	zone, offset := now.Zone()

	osNameAndVersion := tlang.LazyValue[string]{Create: func() string {
		name, version := getOSNameAndVersion()
		return name + "," + version
	}}

	cacheDir, hasCacheDir := os.LookupEnv(constant.EnvName_DUKKHA_CACHE_DIR)
	if !hasCacheDir {
		cacheDir = filepath.Join(cwd, constant.DefaultCacheDir)
	}

	hostArch := &tlang.LazyValue[string]{Create: sysinfo.Arch}

	return &dukkha.GlobalEnvSet{
		constant.GlobalEnv_DUKKHA_WORKDIR: tlang.ImmediateString(cwd),
		// cache dir can be set in config
		constant.GlobalEnv_DUKKHA_CACHE_DIR: tlang.ImmediateString(cacheDir),

		constant.GlobalEnv_TIME_ZONE:        tlang.ImmediateString(zone),
		constant.GlobalEnv_TIME_ZONE_OFFSET: tlang.ImmediateString(strconv.FormatInt(int64(offset), 10)),
		constant.GlobalEnv_TIME_YEAR:        tlang.ImmediateString(strconv.FormatInt(int64(now.Year()), 10)),
		constant.GlobalEnv_TIME_MONTH:       tlang.ImmediateString(strconv.FormatInt(int64(now.Month()), 10)),
		constant.GlobalEnv_TIME_DAY:         tlang.ImmediateString(strconv.FormatInt(int64(now.Day()), 10)),
		constant.GlobalEnv_TIME_HOUR:        tlang.ImmediateString(strconv.FormatInt(int64(now.Hour()), 10)),
		constant.GlobalEnv_TIME_MINUTE:      tlang.ImmediateString(strconv.FormatInt(int64(now.Minute()), 10)),
		constant.GlobalEnv_TIME_SECOND:      tlang.ImmediateString(strconv.FormatInt(int64(now.Second()), 10)),

		constant.GlobalEnv_HOST_KERNEL:         tlang.ImmediateString(runtime.GOOS),
		constant.GlobalEnv_HOST_KERNEL_VERSION: &tlang.LazyValue[string]{Create: sysinfo.KernelVersion},

		constant.GlobalEnv_HOST_OS: &tlang.LazyValue[string]{Create: func() string {
			nameAndVer := osNameAndVersion.GetLazyValue()
			return nameAndVer[:strings.IndexByte(nameAndVer, ',')]
		}},

		constant.GlobalEnv_HOST_OS_VERSION: &tlang.LazyValue[string]{Create: func() string {
			nameAndVer := osNameAndVersion.GetLazyValue()
			return nameAndVer[strings.IndexByte(nameAndVer, ',')+1:]
		}},

		constant.GlobalEnv_HOST_ARCH: hostArch,
		constant.GlobalEnv_HOST_ARCH_SIMPLE: &tlang.LazyValue[string]{Create: func() string {
			return constant.SimpleArch(hostArch.GetLazyValue())
		}},
		constant.GlobalEnv_GIT_BRANCH: GitBranch(ctx, cwd),
		constant.GlobalEnv_GIT_COMMIT: GitCommit(ctx, cwd),
		constant.GlobalEnv_GIT_TAG:    GitTag(ctx, cwd),

		constant.GlobalEnv_GIT_WORKTREE_CLEAN: GitWorkTreeClean(ctx, cwd),
		constant.GlobalEnv_GIT_DEFAULT_BRANCH: GitDefaultBranch(ctx, cwd),
	}
}

func getOSNameAndVersion() (osName, osVersion string) {
	switch runtime.GOOS {
	case constant.KERNEL_Linux:
		osReleaseFile, err2 := os.Open("/etc/os-release")
		if err2 != nil {
			break
		}
		defer func() { _ = osReleaseFile.Close() }()

		s := bufio.NewScanner(osReleaseFile)
		s.Split(bufio.ScanLines)

		for s.Scan() {
			line := s.Text()
			switch {
			case strings.HasPrefix(line, "ID="):
				// TODO: ubuntu has ID_LIKE=debian, check other platforms

				osName = strings.TrimPrefix(line, "ID=")
				osName = strings.TrimRight(strings.TrimLeft(osName, `"`), `"`)
			case strings.HasPrefix(line, "VERSION_ID="):
				osVersion = strings.TrimPrefix(line, "VERSION_ID=")
				osVersion = strings.TrimRight(strings.TrimLeft(osVersion, `"`), `"`)
			}
		}
	default:
		// TODO: support other os
	}

	return
}

func newLazyExecVal(
	ctx context.Context,
	dir string,
	command []string,
	onError func() string,
	onSuccess func(string) string,
) *tlang.LazyValue[string] {
	return &tlang.LazyValue[string]{
		Create: func() string {
			var buf strings.Builder
			cmd, err2 := exechelper.Do(exechelper.Spec{
				Context: ctx,
				Dir:     dir,
				Command: command,
				Stdout:  &buf,
				Stderr:  io.Discard,
			})
			if err2 != nil {
				return onError()
			}

			_, err2 = cmd.Wait()
			if err2 != nil {
				return onError()
			}

			return onSuccess(buf.String())
		},
	}
}
