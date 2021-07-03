package cmd

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"arhat.dev/pkg/exechelper"
	"go.uber.org/multierr"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/sysinfo"
)

// TODO(all): Update docs/environment-variables.md when updating this file

func populateGlobalEnv(ctx context.Context) error {
	envs := []struct {
		name      string
		command   []string
		onError   func() string
		onSuccess func(result string) string
	}{
		{
			name: constant.ENV_GIT_BRANCH,
			command: []string{
				"git", "symbolic-ref", "--short", "-q", "HEAD",
			},
			onError:   func() string { return "" },
			onSuccess: strings.TrimSpace,
		},
		{
			name: constant.ENV_GIT_COMMIT,
			command: []string{
				"git", "rev-parse", "HEAD",
			},
			onError:   func() string { return "" },
			onSuccess: strings.TrimSpace,
		},
		{
			name: constant.ENV_GIT_TAG,
			command: []string{
				"git", "describe", "--tags",
			},
			onError: func() string { return "" },
			onSuccess: func(result string) string {
				return strings.TrimSpace(strings.SplitN(result, " ", 2)[0])
			},
		},
		{
			name: constant.ENV_GIT_WORKTREE_CLEAN,
			command: []string{
				"git", "diff-index", "--quiet", "HEAD",
			},
			onError:   func() string { return "false" },
			onSuccess: func(result string) string { return "true" },
		},
		{
			name: constant.ENV_GIT_DEFAULT_BRANCH,
			command: []string{
				"git", "symbolic-ref", "refs/remotes/origin/HEAD",
			},
			onError: func() string { return os.Getenv(constant.ENV_GIT_DEFAULT_BRANCH) },
			onSuccess: func(result string) string {
				ret := strings.TrimSpace(
					strings.TrimPrefix(result, "refs/remotes/origin/"),
				)
				if len(ret) != 0 {
					return ret
				}

				return os.Getenv(constant.ENV_GIT_DEFAULT_BRANCH)
			},
		},
	}

	var err error
	for _, e := range envs {
		buf := &bytes.Buffer{}
		cmd, err2 := exechelper.Do(exechelper.Spec{
			Context: ctx,
			Command: e.command,
			Stdout:  buf,
			Stderr:  ioutil.Discard,
		})
		if err2 != nil {
			err = multierr.Append(err, os.Setenv(e.name, e.onError()))
			continue
		}

		_, err2 = cmd.Wait()
		if err2 != nil {
			err = multierr.Append(err, os.Setenv(e.name, e.onError()))
			continue
		}

		err = multierr.Append(err, os.Setenv(e.name, e.onSuccess(buf.String())))
	}

	now := time.Now()
	for k, v := range map[string]string{
		constant.ENV_DUKKHA_WORKING_DIR: func() string {
			pwd, err2 := os.Getwd()
			if err2 != nil {
				return ""
			}

			pwd, err2 = filepath.Abs(pwd)
			if err2 != nil {
				panic(fmt.Errorf("failed to get dukkha working dir: %w", err2))
			}

			return pwd
		}(),

		constant.ENV_TIME_YEAR:   strconv.FormatInt(int64(now.Year()), 10),
		constant.ENV_TIME_MONTH:  strconv.FormatInt(int64(now.Month()), 10),
		constant.ENV_TIME_DAY:    strconv.FormatInt(int64(now.Day()), 10),
		constant.ENV_TIME_HOUR:   strconv.FormatInt(int64(now.Hour()), 10),
		constant.ENV_TIME_MINUTE: strconv.FormatInt(int64(now.Minute()), 10),
		constant.ENV_TIME_SECOND: strconv.FormatInt(int64(now.Second()), 10),

		constant.ENV_HOST_KERNEL:         runtime.GOOS,
		constant.ENV_HOST_KERNEL_VERSION: sysinfo.KernelVersion(),

		constant.ENV_HOST_OS:         "",
		constant.ENV_HOST_OS_VERSION: "",

		constant.ENV_HOST_ARCH: sysinfo.Arch(),
	} {
		err = multierr.Append(err, os.Setenv(k, v))
	}

	// set host os name and version
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
				osName := strings.TrimPrefix(line, "ID=")
				osName = strings.TrimRight(strings.TrimLeft(osName, `"`), `"`)

				// TODO: ubuntu has ID_LIKE=debian, check other platforms
				err = multierr.Append(err, os.Setenv(constant.ENV_HOST_OS, osName))
			case strings.HasPrefix(line, "VERSION_ID="):
				osVersion := strings.TrimPrefix(line, "VERSION_ID=")
				osVersion = strings.TrimRight(strings.TrimLeft(osVersion, `"`), `"`)

				err = multierr.Append(err, os.Setenv(constant.ENV_HOST_OS_VERSION, osVersion))
			}
		}
	default:
	}

	// check ci platform specific settings

	switch {
	case os.Getenv("GITHUB_ACTIONS") == "true":
		// github actions

		// https://docs.github.com/en/actions/reference/environment-variables
		commit := strings.TrimSpace(os.Getenv("GITHUB_SHA"))
		if len(commit) != 0 {
			err = multierr.Append(err, os.Setenv(constant.ENV_GIT_COMMIT, commit))
		}

		branch := strings.TrimSpace(strings.TrimPrefix(os.Getenv("GITHUB_REF"), "refs/heads/"))
		if len(branch) != 0 {
			err = multierr.Append(err, os.Setenv(constant.ENV_GIT_BRANCH, branch))
		}
	case os.Getenv("GITLAB_CI") == "true":
		// gitlab-ci

		// https://docs.gitlab.com/ee/ci/variables/predefined_variables.html

		commit := strings.TrimSpace(os.Getenv("CI_COMMIT_SHA"))
		if len(commit) != 0 {
			err = multierr.Append(err, os.Setenv(constant.ENV_GIT_COMMIT, commit))
		}

		branch := strings.TrimSpace(os.Getenv("CI_COMMIT_BRANCH"))
		if len(branch) != 0 {
			err = multierr.Append(err, os.Setenv(constant.ENV_GIT_BRANCH, branch))
		}

		tag := strings.TrimSpace(os.Getenv("CI_COMMIT_TAG"))
		if len(tag) != 0 {
			err = multierr.Append(err, os.Setenv(constant.ENV_GIT_TAG, tag))
		}
	default:
	}

	return err
}
