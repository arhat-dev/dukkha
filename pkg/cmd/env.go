package cmd

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"arhat.dev/pkg/exechelper"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/sysinfo"
	"arhat.dev/dukkha/pkg/utils"
)

// TODO(all): Update docs/environment-variables.md when updating this file

func createGlobalEnv(ctx context.Context, cwd string) map[string]utils.LazyValue {
	now := time.Now().Local()
	zone, offset := now.Zone()

	osNameAndVersion := utils.NewLazyValue(func() string {
		name, version := getOSNameAndVersion()
		return name + "," + version
	})

	cwd, err := filepath.Abs(cwd)
	if err != nil {
		panic(fmt.Errorf("failed to get dukkha working dir: %w", err))
	}

	return map[string]utils.LazyValue{
		constant.ENV_DUKKHA_WORKING_DIR: utils.ImmediateString(cwd),

		constant.ENV_TIME_ZONE:        utils.ImmediateString(zone),
		constant.ENV_TIME_ZONE_OFFSET: utils.ImmediateString(strconv.FormatInt(int64(offset), 10)),
		constant.ENV_TIME_YEAR:        utils.ImmediateString(strconv.FormatInt(int64(now.Year()), 10)),
		constant.ENV_TIME_MONTH:       utils.ImmediateString(strconv.FormatInt(int64(now.Month()), 10)),
		constant.ENV_TIME_DAY:         utils.ImmediateString(strconv.FormatInt(int64(now.Day()), 10)),
		constant.ENV_TIME_HOUR:        utils.ImmediateString(strconv.FormatInt(int64(now.Hour()), 10)),
		constant.ENV_TIME_MINUTE:      utils.ImmediateString(strconv.FormatInt(int64(now.Minute()), 10)),
		constant.ENV_TIME_SECOND:      utils.ImmediateString(strconv.FormatInt(int64(now.Second()), 10)),

		constant.ENV_HOST_KERNEL:         utils.ImmediateString(runtime.GOOS),
		constant.ENV_HOST_KERNEL_VERSION: utils.NewLazyValue(sysinfo.KernelVersion),

		constant.ENV_HOST_OS: utils.NewLazyValue(func() string {
			nameAndVer := osNameAndVersion.Get()
			return nameAndVer[:strings.IndexByte(nameAndVer, ',')]
		}),

		constant.ENV_HOST_OS_VERSION: utils.NewLazyValue(func() string {
			nameAndVer := osNameAndVersion.Get()
			return nameAndVer[strings.IndexByte(nameAndVer, ',')+1:]
		}),

		constant.ENV_HOST_ARCH:  utils.NewLazyValue(sysinfo.Arch),
		constant.ENV_GIT_BRANCH: GitBranch(ctx, cwd),
		constant.ENV_GIT_COMMIT: GitCommit(ctx, cwd),
		constant.ENV_GIT_TAG:    GitTag(ctx, cwd),

		constant.ENV_GIT_WORKTREE_CLEAN: GitWorkTreeClean(ctx, cwd),
		constant.ENV_GIT_DEFAULT_BRANCH: GitDefaultBranch(ctx, cwd),
	}
}

func getOSNameAndVersion() (osName, osVersion string) {
	switch runtime.GOOS {
	case constant.KERNEL_LINUX:
		data, err2 := os.ReadFile("/etc/os-release")
		if err2 != nil {
			break
		}

		s := bufio.NewScanner(bytes.NewReader(data))
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
) utils.LazyValue {
	buf := &bytes.Buffer{}
	return utils.NewLazyValue(func() string {
		cmd, err2 := exechelper.Do(exechelper.Spec{
			Context: ctx,
			Dir:     dir,
			Command: command,
			Stdout:  buf,
			Stderr:  io.Discard,
		})
		if err2 != nil {
			return onError()
		}

		_, err2 = cmd.Wait()
		if err2 != nil {
			return onError()
		}

		return onSuccess(string(buf.Next(buf.Len())))
	})
}
