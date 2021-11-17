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

func createGlobalEnv(ctx context.Context) map[string]utils.LazyValue {
	now := time.Now().Local()
	zone, offset := now.Zone()

	// ref: https://docs.github.com/en/actions/reference/environment-variables
	isGithubActions := os.Getenv("GITHUB_ACTIONS") == "true"

	// ref: https://docs.gitlab.com/ee/ci/variables/predefined_variables.html
	isGitlabCI := os.Getenv("GITLAB_CI") == "true"

	osNameAndVersion := utils.NewLazyValue(func() string {
		name, version := getOSNameAndVersion()
		return name + "," + version
	})

	return map[string]utils.LazyValue{
		constant.ENV_DUKKHA_WORKING_DIR: utils.ImmediateString(func() string {
			pwd, err2 := os.Getwd()
			if err2 != nil {
				return ""
			}

			pwd, err2 = filepath.Abs(pwd)
			if err2 != nil {
				panic(fmt.Errorf("failed to get dukkha working dir: %w", err2))
			}

			return pwd
		}()),

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

		constant.ENV_HOST_ARCH: utils.NewLazyValue(sysinfo.Arch),
		constant.ENV_GIT_BRANCH: newLazyExecVal(
			ctx,
			[]string{
				"git", "symbolic-ref", "--short", "-q", "HEAD",
			},
			func() string {
				switch {
				case isGithubActions:
					ghRef := strings.TrimSpace(os.Getenv("GITHUB_REF"))
					if len(ghRef) == 0 {
						ghRef = strings.TrimSpace(os.Getenv("GITHUB_HEAD_REF"))
					}

					switch {
					case strings.HasPrefix(ghRef, "refs/heads/"):
						return strings.TrimPrefix(ghRef, "refs/heads/")
					default:
						return ""
					}
				case isGitlabCI:
					return strings.TrimSpace(os.Getenv("CI_COMMIT_BRANCH"))
				default:
					return ""
				}
			},
			strings.TrimSpace,
		),
		constant.ENV_GIT_COMMIT: newLazyExecVal(
			ctx,
			[]string{
				"git", "rev-parse", "HEAD",
			},
			func() string {
				switch {
				case isGithubActions:
					return strings.TrimSpace(os.Getenv("GITHUB_SHA"))
				case isGitlabCI:
					return strings.TrimSpace(os.Getenv("CI_COMMIT_SHA"))
				default:
					return ""
				}
			},
			strings.TrimSpace,
		),

		constant.ENV_GIT_TAG: newLazyExecVal(
			ctx,
			[]string{
				"git", "describe", "--tags",
			},
			func() string {
				switch {
				case isGithubActions:
					ghRef := strings.TrimSpace(os.Getenv("GITHUB_REF"))
					if len(ghRef) == 0 {
						ghRef = strings.TrimSpace(os.Getenv("GITHUB_HEAD_REF"))
					}

					switch {
					case strings.HasPrefix(ghRef, "refs/tags/"):
						return strings.TrimPrefix(ghRef, "refs/tags/")
					default:
						return ""
					}
				case isGitlabCI:
					return strings.TrimSpace(os.Getenv("CI_COMMIT_TAG"))
				default:
					return ""
				}
			},
			func(result string) string {
				return strings.TrimSpace(strings.SplitN(result, " ", 2)[0])
			},
		),

		constant.ENV_GIT_WORKTREE_CLEAN: newLazyExecVal(
			ctx,
			[]string{
				"git", "diff-index", "--quiet", "HEAD",
			},
			func() string { return "false" },
			func(_ string) string { return "true" },
		),

		constant.ENV_GIT_DEFAULT_BRANCH: newLazyExecVal(
			ctx,
			[]string{
				"git", "remote", "show", "origin",
			},
			func() string { return os.Getenv(constant.ENV_GIT_DEFAULT_BRANCH) },
			func(result string) string {
				s := bufio.NewScanner(strings.NewReader(result))
				s.Split(bufio.ScanLines)
				for s.Scan() {
					line := s.Text()
					const prefix = "HEAD branch: "
					if idx := strings.Index(line, prefix); idx != -1 {
						return line[idx+len(prefix):]
					}
				}

				return os.Getenv(constant.ENV_GIT_DEFAULT_BRANCH)
			},
		),
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
	}

	return
}

func newLazyExecVal(
	ctx context.Context,
	command []string,
	onError func() string,
	onSuccess func(string) string,
) utils.LazyValue {
	buf := &bytes.Buffer{}
	return utils.NewLazyValue(func() string {
		cmd, err2 := exechelper.Do(exechelper.Spec{
			Context: ctx,
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
