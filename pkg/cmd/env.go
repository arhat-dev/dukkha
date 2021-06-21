package cmd

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"arhat.dev/pkg/exechelper"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/sysinfo"
)

func populateGlobalEnv(ctx context.Context) {
	envs := []struct {
		name      string
		command   []string
		onError   func() string
		onSuccess func(result string) string
	}{
		{
			name: constant.EnvGIT_BRANCH,
			command: []string{
				"git", "symbolic-ref", "--short", "-q", "HEAD",
			},
			onError:   func() string { return "" },
			onSuccess: strings.TrimSpace,
		},
		{
			name: constant.EnvGIT_COMMIT,
			command: []string{
				"git", "rev-parse", "HEAD",
			},
			onError:   func() string { return "" },
			onSuccess: strings.TrimSpace,
		},
		{
			name: constant.EnvGIT_TAG,
			command: []string{
				"git", "describe", "--tags",
			},
			onError:   func() string { return "" },
			onSuccess: strings.TrimSpace,
		},
		{
			name: constant.EnvGIT_WORKSPACE_CLEAN,
			command: []string{
				"git", "diff-index", "--quiet", "HEAD",
			},
			onError:   func() string { return "false" },
			onSuccess: func(result string) string { return "true" },
		},
		{
			name: constant.EnvGIT_DEFAULT_BRANCH,
			command: []string{
				"git", "symbolic-ref", "refs/remotes/origin/HEAD",
			},
			onError: func() string { return os.Getenv(constant.EnvGIT_DEFAULT_BRANCH) },
			onSuccess: func(result string) string {
				ret := strings.TrimSpace(
					strings.TrimPrefix(result, "refs/remotes/origin/"),
				)
				if len(ret) != 0 {
					return ret
				}

				return os.Getenv(constant.EnvGIT_DEFAULT_BRANCH)
			},
		},
	}

	for _, e := range envs {
		buf := &bytes.Buffer{}
		cmd, err := exechelper.Do(exechelper.Spec{
			Context: ctx,
			Command: e.command,
			Stdout:  buf,
			Stderr:  ioutil.Discard,
		})
		if err != nil {
			os.Setenv(e.name, e.onError())
			continue
		}

		_, err = cmd.Wait()
		if err != nil {
			os.Setenv(e.name, e.onError())
			continue
		}

		os.Setenv(e.name, e.onSuccess(buf.String()))
	}

	now := time.Now()
	for k, v := range map[string]string{
		constant.EnvTIME_YEAR:   strconv.FormatInt(int64(now.Year()), 10),
		constant.EnvTIME_MONTH:  strconv.FormatInt(int64(now.Month()), 10),
		constant.EnvTIME_DAY:    strconv.FormatInt(int64(now.Day()), 10),
		constant.EnvTIME_HOUR:   strconv.FormatInt(int64(now.Hour()), 10),
		constant.EnvTIME_MINUTE: strconv.FormatInt(int64(now.Minute()), 10),
		constant.EnvTIME_SECOND: strconv.FormatInt(int64(now.Second()), 10),
		constant.EnvHOST_OS:     runtime.GOOS,
		constant.EnvHOST_ARCH:   sysinfo.Arch(),
	} {
		os.Setenv(k, v)
	}

	// check ci platform specific settings

	switch {
	case os.Getenv("GITHUB_ACTIONS") == "true":
		// github actions

		// https://docs.github.com/en/actions/reference/environment-variables
		commit := strings.TrimSpace(os.Getenv("GITHUB_SHA"))
		if len(commit) != 0 {
			os.Setenv(constant.EnvGIT_COMMIT, commit)
		}

		branch := strings.TrimSpace(strings.TrimPrefix(os.Getenv("GITHUB_REF"), "refs/heads/"))
		if len(branch) != 0 {
			os.Setenv(constant.EnvGIT_BRANCH, branch)
		}
	case os.Getenv("GITLAB_CI") == "true":
		// gitlab-ci

		// https://docs.gitlab.com/ee/ci/variables/predefined_variables.html

		commit := strings.TrimSpace(os.Getenv("CI_COMMIT_SHA"))
		if len(commit) != 0 {
			os.Setenv(constant.EnvGIT_COMMIT, commit)
		}

		branch := strings.TrimSpace(os.Getenv("CI_COMMIT_BRANCH"))
		if len(branch) != 0 {
			os.Setenv(constant.EnvGIT_BRANCH, branch)
		}

		tag := strings.TrimSpace(os.Getenv("CI_COMMIT_TAG"))
		if len(tag) != 0 {
			os.Setenv(constant.EnvGIT_TAG, tag)
		}
	default:
	}
}
