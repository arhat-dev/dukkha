package golang

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"arhat.dev/pkg/byteshelper"
	"arhat.dev/pkg/fshelper"
	"arhat.dev/pkg/md5helper"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindTest = "test"

func init() {
	dukkha.RegisterTask(ToolKind, TaskKindTest, tools.NewTask[TaskTest, *TaskTest])
}

type TaskTest struct {
	tools.BaseTask[GolangTest, *GolangTest]
}

// nolint:revive
type GolangTest struct {
	CGO CGOSepc `yaml:"cgo"`

	Path  string `yaml:"path"`
	Chdir string `yaml:"chdir"`

	BuildOptions buildOptions `yaml:",inline"`
	Test         testSpec     `yaml:",inline"`

	Benchmark testBenchmarkSpec `yaml:"benchmark"`
	Profile   testProfileSpec   `yaml:"profile"`

	// ExtraArgs for go test (inserted before `Path`)
	ExtraArgs []string `yaml:"extra_args"`

	// CustomCmdPrefix to run compiled test file with this cmd prefix
	// e.g. built xxx.test, usually will run in local host as ./xxx.test
	//      but with `custom_cmd_prefx=[ssh, testsrv]`, will run as `ssh testsrv xxx.test`
	CustomCmdPrefix []string `yaml:"custom_cmd_prefix"`

	// CustomArgs appended when running the test
	CustomArgs []string `yaml:"custom_args"`

	parent tools.BaseTaskType
}

func (c *GolangTest) ToolKind() dukkha.ToolKind       { return ToolKind }
func (c *GolangTest) Kind() dukkha.TaskKind           { return TaskKindTest }
func (c *GolangTest) LinkParent(p tools.BaseTaskType) { c.parent = p }

func (c *GolangTest) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var steps []dukkha.TaskExecSpec
	err := c.parent.DoAfterFieldsResolved(rc, -1, true, func() error {
		// get a list of packages to be tested
		const (
			targetReplaceABSPackageDirs    = "<GO_ABS_PACKAGE_DIRS>"
			targetReplaceGoListErrorResult = "<GO_LIST_ERROR_RESULT>"
		)
		steps = append(steps, dukkha.TaskExecSpec{
			StdoutAsReplace: targetReplaceABSPackageDirs,
			StderrAsReplace: targetReplaceGoListErrorResult,

			Chdir: c.Chdir,
			Command: []string{
				constant.DUKKHA_TOOL_CMD, "list", "-f", `{{ .Dir }}`, c.Path,
			},
			IgnoreError: true,
		})

		// copy values and do not reference task fields as they are generated dynamically
		toolCmd := []string{constant.DUKKHA_TOOL_CMD}

		_fs, err := rc.FS().Sub(c.Chdir)
		if err != nil {
			return err
		}
		cwdFS := _fs.(*fshelper.OSFS)

		workDir := c.Test.WorkDir
		jsonOutputFile := c.Test.JSONOutputFile

		var compileArgs []string

		buildEnv := createBuildEnv(rc, c.BuildOptions, c.CGO)

		compileArgs = append(compileArgs, c.BuildOptions.generateArgs()...)
		compileArgs = append(compileArgs, c.Test.generateArgs(true)...)
		compileArgs = append(compileArgs, c.Benchmark.generateArgs(true)...)
		compileArgs = append(compileArgs, c.Profile.generateArgs(cwdFS, true)...)
		compileArgs = append(compileArgs, c.ExtraArgs...)

		runCmdPrefix := sliceutils.NewStrings(c.CustomCmdPrefix)
		var runArgs []string
		runArgs = append(runArgs, c.Test.generateArgs(false)...)
		runArgs = append(runArgs, c.Benchmark.generateArgs(false)...)
		runArgs = append(runArgs, c.Profile.generateArgs(cwdFS, false)...)
		if len(c.CustomArgs) != 0 {
			runArgs = append(runArgs, "--")
			runArgs = append(runArgs, c.CustomArgs...)
		}

		steps = append(steps, dukkha.TaskExecSpec{
			AlterExecFunc: func(
				replace dukkha.ReplaceEntries,
				stdin io.Reader, stdout, stderr io.Writer,
			) (dukkha.RunTaskOrRunCmd, error) {
				stdoutResult, ok := replace[targetReplaceABSPackageDirs]
				if ok && stdoutResult.Err == nil {
					// found packages to be tested, test these packages
					var (
						compileSteps []dukkha.TaskExecSpec
						runSteps     []dukkha.TaskExecSpec
					)

					for _, absPkgDir := range strings.Split(byteshelper.ToString(stdoutResult.Data), "\n") {
						absPkgDir = strings.TrimSpace(absPkgDir)
						if len(absPkgDir) == 0 {
							continue
						}

						builtTestExecutable, subCompileSteps := generateCompileSpecs(
							c.parent.CacheFS(),
							cwdFS,
							buildEnv, compileArgs, absPkgDir,
							toolCmd,
						)

						compileSteps = append(compileSteps, subCompileSteps...)

						subRunSpecs := generateRunSpecs(
							cwdFS,
							builtTestExecutable,
							workDir,

							toolCmd,
							runCmdPrefix,
							runArgs,
							absPkgDir,
							jsonOutputFile,
						)

						runSteps = append(runSteps, subRunSpecs...)
					}

					return append(compileSteps, runSteps...), nil
				}

				// no stdout data, check stderr data
				stderrResult, ok := replace[targetReplaceGoListErrorResult]
				if ok && stderrResult.Err != nil {
					if bytes.Contains(stderrResult.Data, []byte("no Go files")) {
						// no go source file in this package, skip
						return nil, nil
					}

					return nil, stderrResult.Err
				}

				return nil, fmt.Errorf("unable to determine which packages to test")
			},
			IgnoreError: false,
		})

		return nil
	})

	return steps, err
}

func getGoTestCompileResultReplaceKey(pkgDir string) string {
	return fmt.Sprintf("<GO_TEST_COMPILE_RESULT:%s>", hex.EncodeToString(md5helper.Sum([]byte(pkgDir))))
}

func getBuiltTestExecutablePath(pkgDir string) string {
	return hex.EncodeToString(md5helper.Sum([]byte(pkgDir))) + ".test"
}

// compile one package for testing at a time
func generateCompileSpecs(
	cacheFS *fshelper.OSFS, // cachefs to store compiled data
	cwdFS *fshelper.OSFS, // cwdfs with current workdir in target path
	buildEnv dukkha.NameValueList,
	args []string,
	pkgDir string,

	// options
	toolCmd []string,
) (string, []dukkha.TaskExecSpec) {
	var steps []dukkha.TaskExecSpec

	builtTestExecutable, err := cacheFS.Abs(
		getBuiltTestExecutablePath(pkgDir),
	)
	if err != nil {
		panic(err)
	}

	// remove previously built test executable if any
	steps = append(steps, dukkha.TaskExecSpec{
		AlterExecFunc: func(
			replace dukkha.ReplaceEntries,
			stdin io.Reader,
			stdout, stderr io.Writer,
		) (dukkha.RunTaskOrRunCmd, error) {
			err2 := cacheFS.Remove(builtTestExecutable)
			if err2 != nil && !errors.Is(err2, fs.ErrNotExist) {
				return nil, fmt.Errorf("removing previously built test executable: %w", err2)
			}

			return nil, nil
		},
	})

	compileCmd := sliceutils.NewStrings(
		toolCmd, "test", "-c", "-o", builtTestExecutable,
	)

	compileCmd = append(compileCmd, args...)

	chdir, err := cwdFS.Abs(".")
	if err != nil {
		panic(err)
	}

	steps = append(steps, dukkha.TaskExecSpec{
		StdoutAsReplace: getGoTestCompileResultReplaceKey(pkgDir),
		ShowStdout:      true,

		EnvSuggest:  buildEnv,
		Chdir:       chdir,
		Command:     append(compileCmd, pkgDir),
		IgnoreError: false,
	},
		dukkha.TaskExecSpec{
			AlterExecFunc: func(
				replace dukkha.ReplaceEntries,
				stdin io.Reader,
				stdout, stderr io.Writer,
			) (dukkha.RunTaskOrRunCmd, error) {
				return nil, cacheFS.Chmod(builtTestExecutable, 0750)
			},
			IgnoreError: true,
		},
	)

	return builtTestExecutable, steps
}

func getTestRunResultReplaceKey(pkgRelPath string) string {
	return fmt.Sprintf("<GO_TEST_RUN_RESULT:%s>", hex.EncodeToString(md5helper.Sum([]byte(pkgRelPath))))
}

func getGoToolTest2JsonResultReplaceKey(pkgRelPath string) string {
	return fmt.Sprintf("<GO_TOOL_TEST2JSON_RESULT:%s>", hex.EncodeToString(md5helper.Sum([]byte(pkgRelPath))))
}

func generateRunSpecs(
	cwdFS *fshelper.OSFS,
	builtTestExecutable string,
	_workdir string,

	toolCmd []string,

	cmdPrefix []string,
	args []string,
	absPkgDir string,
	jsonOutputFile string,
) []dukkha.TaskExecSpec {
	var steps []dukkha.TaskExecSpec

	workdir := _workdir
	switch {
	case len(workdir) == 0:
		// use same default workdir as go test (the pakcage dir)
		workdir = absPkgDir
	case filepath.IsAbs(workdir):
		// just use it
	default:
		// workdir is a relative path
		var err error
		workdir, err = cwdFS.Abs(workdir)
		if err != nil {
			panic(err)
		}
	}

	runCmd := sliceutils.NewStrings(cmdPrefix, builtTestExecutable)
	runCmd = append(runCmd, args...)

	stdoutReplaceKey := ""
	if len(jsonOutputFile) != 0 {
		stdoutReplaceKey = getTestRunResultReplaceKey(absPkgDir)
	}

	// check if compiled test file exists
	// can be missing if no test was found in the package
	steps = append(steps, dukkha.TaskExecSpec{
		AlterExecFunc: func(
			replace dukkha.ReplaceEntries,
			stdin io.Reader,
			stdout, stderr io.Writer,
		) (dukkha.RunTaskOrRunCmd, error) {
			// we only compile a package for testing at a time
			//
			// so we can check if the package contains tests:
			// 		if the output of `go test -c` contains "no test file"
			// 		go test will not produce a test executable
			//
			// and we need to skip this package

			compileResult, ok := replace[getGoTestCompileResultReplaceKey(absPkgDir)]
			if ok && bytes.Contains(compileResult.Data, []byte("no test files")) {
				// no test
				return nil, nil
			}

			subSteps := []dukkha.TaskExecSpec{
				{
					StdoutAsReplace: stdoutReplaceKey,
					ShowStdout:      true,

					Chdir:   workdir,
					Command: runCmd,
				},
			}

			if len(jsonOutputFile) == 0 {
				return subSteps, nil
			}

			subSteps = append(subSteps, dukkha.TaskExecSpec{
				AlterExecFunc: func(
					replace dukkha.ReplaceEntries,
					stdin io.Reader,
					stdout, stderr io.Writer,
				) (dukkha.RunTaskOrRunCmd, error) {
					testOutput, ok := replace[stdoutReplaceKey]
					if !ok {
						return nil, fmt.Errorf("test output not found")
					}

					resultKey := getGoToolTest2JsonResultReplaceKey(absPkgDir)
					return []dukkha.TaskExecSpec{
						{
							StdoutAsReplace: resultKey,
							Stdin:           bytes.NewReader(testOutput.Data),
							Command:         sliceutils.NewStrings(toolCmd, "tool", "test2json"),
						},
						{
							AlterExecFunc: func(
								replace dukkha.ReplaceEntries,
								stdin io.Reader,
								stdout, stderr io.Writer,
							) (dukkha.RunTaskOrRunCmd, error) {
								jsonOutput, ok := replace[resultKey]
								if !ok {
									return nil, fmt.Errorf("json of test result not found")
								}

								err := cwdFS.WriteFile(jsonOutputFile, jsonOutput.Data, 0644)
								if err != nil {
									return nil, fmt.Errorf("saving test json output: %w", err)
								}

								return nil, nil
							},
						},
					}, nil
				},
			})

			return subSteps, nil
		},
	})

	return steps
}

type testSpec struct {
	rs.BaseField `yaml:"-"`

	LogFile string `yaml:"log_file"`

	// go test -count
	Count int `yaml:"count"`

	// go test -cpu
	CPU []int `yaml:"cpu"`

	// go test -parallel
	Parallel int `yaml:"parallel"`

	// go test -failfast
	FailFast bool `yaml:"failfast"`

	// go test -short
	Short bool `yaml:"short"`

	// go test -timeout
	Timeout time.Duration `yaml:"timeout"`

	// go test -run
	Match string `yaml:"match"`

	// go test -v
	Verbose bool `yaml:"verbose"`

	// JSONOutputFile
	JSONOutputFile string `yaml:"json_output_file"`

	// Panic on calling os.Exit(0)
	PanicOnExit0 bool `yaml:"panic_on_exit_0"`

	// WorkDir to run test, defaults to DUKKHA_WORKDIR
	WorkDir string `yaml:"work_dir"`
}

func (s testSpec) generateArgs(compileTime bool) []string {
	var args []string

	prefix := getTestFlagPrefix(compileTime)

	if s.Count != 0 {
		args = append(args, prefix+"count", strconv.FormatInt(int64(s.Count), 10))
	} else {
		// disables test caching
		args = append(args, prefix+"count", "1")
	}

	if len(s.CPU) != 0 {
		var cpu []string
		for _, c := range s.CPU {
			cpu = append(cpu, strconv.FormatInt(int64(c), 10))
		}

		args = append(args, prefix+"cpu", strings.Join(cpu, ","))
	}

	if s.Parallel != 0 {
		args = append(args, prefix+"parallel", strconv.FormatInt(int64(s.Parallel), 10))
	}

	if s.FailFast {
		args = append(args, prefix+"failfast")
	}

	if s.Short {
		args = append(args, prefix+"short")
	}

	if s.Timeout != 0 {
		args = append(args, prefix+"timeout", s.Timeout.String())
	}

	if len(s.Match) != 0 {
		args = append(args, prefix+"run", s.Match)
	}

	if s.Verbose || len(s.JSONOutputFile) != 0 {
		args = append(args, prefix+"v")
	}

	if len(s.LogFile) != 0 && !compileTime {
		args = append(args, "-test.testlogfile", s.LogFile)
	}

	if s.PanicOnExit0 && !compileTime {
		args = append(args, "-test.paniconexit0")
	}

	return args
}

func getTestFlagPrefix(isCompileTime bool) string {
	if isCompileTime {
		return "-"
	}

	return "-test."
}
