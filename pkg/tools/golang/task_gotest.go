package golang

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"arhat.dev/pkg/md5helper"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindTest = "test"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindTest,
		func(toolName string) dukkha.Task {
			t := &TaskTest{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindTest, t)
			return t
		},
	)
}

type TaskTest struct {
	rs.BaseField

	tools.BaseTask `yaml:",inline"`

	CGO CGOSepc `yaml:"cgo"`

	Path  string `yaml:"path"`
	Chdir string `yaml:"chdir"`

	Build buildOptions `yaml:",inline"`
	Test  testSpec     `yaml:",inline"`

	Benchmark testBenchmarkSpec `yaml:"benchmark"`
	Profile   testProfileSpec   `yaml:"profile"`

	// CustomCmdPrefix to run compiled test file with this cmd prefix
	CustomCmdPrefix []string `yaml:"custom_cmd_prefix"`

	// custom args only used when running the test
	CustomArgs []string `yaml:"custom_args"`
}

func (c *TaskTest) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var steps []dukkha.TaskExecSpec
	err := c.DoAfterFieldsResolved(rc, -1, func() error {
		// get package prefix to be trimed
		const targetReplaceModuleName = "<MODULE_NAME>"
		steps = append(steps, dukkha.TaskExecSpec{
			StdoutAsReplace:          targetReplaceModuleName,
			FixStdoutValueForReplace: bytes.TrimSpace,

			Chdir:       c.Chdir,
			Command:     sliceutils.NewStrings(options.ToolCmd(), "list", "-m"),
			IgnoreError: false,
		})

		// get a list of packages to be tested
		listFormat := `{{ .ImportPath }}`
		if options.UseShell() {
			listFormat = `'` + listFormat + `'`
		}
		const (
			targetReplacePackages          = "<GO_PACKAGES>"
			targetReplaceGoListErrorResult = "<GO_LIST_ERROR_RESULT>"
		)
		steps = append(steps, dukkha.TaskExecSpec{
			StdoutAsReplace: targetReplacePackages,
			StderrAsReplace: targetReplaceGoListErrorResult,

			Chdir: c.Chdir,
			Command: sliceutils.NewStrings(
				options.ToolCmd(), "list", "-f", listFormat, c.Path,
			),
			IgnoreError: true,
		})

		// copy values and do not reference task fields
		// since they are generated dynamically
		taskName := c.TaskName
		dukkhaCacheDir := rc.CacheDir()
		dukkhaWorkingDir := rc.WorkingDir()
		toolCmd := sliceutils.NewStrings(options.ToolCmd())
		useShell := options.UseShell()
		shellName := options.ShellName()
		chdir := c.Chdir
		workDir := c.Test.WorkDir
		jsonOutputFile := c.Test.JSONOutputFile

		var compileArgs []string
		compileArgs = append(compileArgs, c.Build.generateArgs(useShell)...)
		compileArgs = append(compileArgs, c.Test.generateArgs(true)...)
		compileArgs = append(compileArgs, c.Benchmark.generateArgs(true)...)
		compileArgs = append(compileArgs, c.Profile.generateArgs(dukkhaWorkingDir, true)...)

		runCmdPrefix := sliceutils.NewStrings(c.CustomCmdPrefix)
		var runArgs []string
		runArgs = append(runArgs, c.Test.generateArgs(false)...)
		runArgs = append(runArgs, c.Benchmark.generateArgs(false)...)
		runArgs = append(runArgs, c.Profile.generateArgs(dukkhaWorkingDir, false)...)
		if len(c.CustomArgs) != 0 {
			runArgs = append(runArgs, "--")
			runArgs = append(runArgs, c.CustomArgs...)
		}
		buildEnv := createBuildEnv(rc, c.CGO)

		steps = append(steps, dukkha.TaskExecSpec{
			AlterExecFunc: func(
				replace dukkha.ReplaceEntries,
				stdin io.Reader, stdout, stderr io.Writer,
			) (dukkha.RunTaskOrRunCmd, error) {
				moduleNameBytes, ok := replace[targetReplaceModuleName]
				if !ok {
					return nil, fmt.Errorf("unexpected no module name set")
				}
				moduleName := string(moduleNameBytes.Data)

				stdoutResult, ok := replace[targetReplacePackages]
				if ok && stdoutResult.Err == nil {
					// found packages to be tested, test these packages
					var (
						compileSteps []dukkha.TaskExecSpec
						runSteps     []dukkha.TaskExecSpec
					)

					pkgsToTest := strings.Split(string(stdoutResult.Data), "\n")
					for _, pkg := range pkgsToTest {
						pkg = strings.TrimSpace(pkg)
						if len(pkg) == 0 {
							continue
						}
						pkgRelPath := strings.TrimPrefix(pkg, moduleName)
						if strings.HasPrefix(pkgRelPath, "/") {
							pkgRelPath = "." + pkgRelPath
						} else {
							pkgRelPath = "./" + pkgRelPath
						}

						builtTestExecutable, subCompileSteps := generateCompileSpecs(
							taskName,
							dukkhaCacheDir,
							chdir,
							buildEnv, compileArgs, pkgRelPath,
							toolCmd, useShell, shellName,
						)

						compileSteps = append(compileSteps, subCompileSteps...)

						subRunSpecs := generateRunSpecs(
							dukkhaWorkingDir,
							builtTestExecutable,
							chdir,
							workDir,

							toolCmd,
							runCmdPrefix,
							runArgs,
							pkgRelPath,
							jsonOutputFile,

							useShell, shellName,
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

				return nil, fmt.Errorf("failed to determine which packages to test")
			},
			IgnoreError: false,
		})

		return nil
	})

	return steps, err
}

func getGoTestCompileResultReplaceKey(pkgRelPath string) string {
	return fmt.Sprintf("<GO_TEST_COMPILE_RESULT:%s>", hex.EncodeToString(md5helper.Sum([]byte(pkgRelPath))))
}

func getBuiltTestExecutablePath(dukkhaCacheDir, taskName, pkgRelPath string) string {
	return filepath.Join(
		dukkhaCacheDir, "golang-test",
		taskName+"-"+hex.EncodeToString(md5helper.Sum([]byte(pkgRelPath)))+".test",
	)
}

// compile one package for testing at a time
func generateCompileSpecs(
	taskName string,
	dukkhaCacheDir string,
	chdir string,
	buildEnv dukkha.Env,
	args []string,
	pkgRelPath string,

	// options
	toolCmd []string,
	useShell bool,
	shellName string,
) (string, []dukkha.TaskExecSpec) {
	var steps []dukkha.TaskExecSpec

	builtTestExecutable := getBuiltTestExecutablePath(
		dukkhaCacheDir, taskName, pkgRelPath,
	)

	// remove previously built test executable if any
	steps = append(steps, dukkha.TaskExecSpec{
		AlterExecFunc: func(
			replace dukkha.ReplaceEntries,
			stdin io.Reader,
			stdout, stderr io.Writer,
		) (dukkha.RunTaskOrRunCmd, error) {
			err := os.Remove(builtTestExecutable)
			if err != nil && !os.IsNotExist(err) {
				return nil, fmt.Errorf("failed to remove previously built test executable: %w", err)
			}

			return nil, nil
		},
	})

	compileCmd := sliceutils.NewStrings(
		toolCmd, "test", "-c", "-o", builtTestExecutable,
	)

	compileCmd = append(compileCmd, args...)

	steps = append(steps, dukkha.TaskExecSpec{
		StdoutAsReplace: getGoTestCompileResultReplaceKey(pkgRelPath),
		ShowStdout:      true,

		EnvSuggest:  buildEnv,
		Chdir:       chdir,
		Command:     append(compileCmd, pkgRelPath),
		UseShell:    useShell,
		ShellName:   shellName,
		IgnoreError: false,
	})

	return builtTestExecutable, steps
}

func getTestRunResultReplaceKey(pkgRelPath string) string {
	return fmt.Sprintf("<GO_TEST_RUN_RESULT:%s>", hex.EncodeToString(md5helper.Sum([]byte(pkgRelPath))))
}

func getGoToolTest2JsonResultReplaceKey(pkgRelPath string) string {
	return fmt.Sprintf("<GO_TOOL_TEST2JSON_RESULT:%s>", hex.EncodeToString(md5helper.Sum([]byte(pkgRelPath))))
}

func generateRunSpecs(
	dukkhaWorkingDir string,
	builtTestExecutable string,
	chdir string,
	_workdir string,

	toolCmd []string,

	cmdPrefix []string,
	args []string,
	pkgRelPath string,
	jsonOutputFile string,

	// options
	useShell bool,
	shellName string,
) []dukkha.TaskExecSpec {
	var steps []dukkha.TaskExecSpec

	workdir := _workdir
	switch {
	case len(workdir) == 0:
		// use same default workdir as go test (the pakcage dir)
		if filepath.IsAbs(chdir) {
			workdir = filepath.Join(chdir, pkgRelPath)
		} else {
			workdir = filepath.Join(dukkhaWorkingDir, chdir, pkgRelPath)
		}
	case filepath.IsAbs(workdir):
		// just use it
	default:
		// workdir is a relative path
		if filepath.IsAbs(chdir) {
			workdir = filepath.Join(chdir, workdir)
		} else {
			workdir = filepath.Join(dukkhaWorkingDir, chdir, workdir)
		}
	}

	runCmd := sliceutils.NewStrings(cmdPrefix, builtTestExecutable)
	runCmd = append(runCmd, args...)

	stdoutReplaceKey := ""
	if len(jsonOutputFile) != 0 {
		stdoutReplaceKey = getTestRunResultReplaceKey(pkgRelPath)
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

			compileResult, ok := replace[getGoTestCompileResultReplaceKey(pkgRelPath)]
			if ok && bytes.Contains(compileResult.Data, []byte("no test files")) {
				// no test
				return nil, nil
			}

			subSteps := []dukkha.TaskExecSpec{
				{
					StdoutAsReplace: stdoutReplaceKey,
					ShowStdout:      true,

					Chdir:     workdir,
					Command:   runCmd,
					UseShell:  useShell,
					ShellName: shellName,
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
						return nil, fmt.Errorf("failed to get test output")
					}

					resultKey := getGoToolTest2JsonResultReplaceKey(pkgRelPath)
					return []dukkha.TaskExecSpec{
						{
							StdoutAsReplace: resultKey,
							Stdin:           bytes.NewReader(testOutput.Data),
							Command:         sliceutils.NewStrings(toolCmd, "tool", "test2json"),
							UseShell:        useShell,
							ShellName:       shellName,
						},
						{
							AlterExecFunc: func(
								replace dukkha.ReplaceEntries,
								stdin io.Reader,
								stdout, stderr io.Writer,
							) (dukkha.RunTaskOrRunCmd, error) {
								jsonOutput, ok := replace[resultKey]
								if !ok {
									return nil, fmt.Errorf("failed to get json result of test")
								}

								err := os.WriteFile(jsonOutputFile, jsonOutput.Data, 0644)
								if err != nil {
									return nil, fmt.Errorf("failed to save test json output: %w", err)
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
	rs.BaseField

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

	// Panic on call to os.Exit(0)
	PanicOnExit0 bool `yaml:"panic_on_exit_0"`

	// WorkDir to run test, defaults to DUKKHA_WORKING_DIR
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

func getTestFlagPrefix(compileTime bool) string {
	if compileTime {
		return "-"
	}

	return "-test."
}
