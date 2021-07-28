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

	"arhat.dev/pkg/hashhelper"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
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
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	CGO CGOSepc `yaml:"cgo"`

	Path  string `yaml:"path"`
	Chdir string `yaml:"chdir"`

	Build buildOptions `yaml:",inline"`
	Test  testSpec     `yaml:",inline"`

	Benchmark testBenchmarkSpec `yaml:"benchmark"`
	Profile   testProfileSpec   `yaml:"profile"`

	// custom args only used when running the test
	CustomArgs []string `yaml:"custom_args"`
}

func (c *TaskTest) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var steps []dukkha.TaskExecSpec
	err := c.DoAfterFieldsResolved(rc, -1, func() error {
		mKernel := rc.MatrixKernel()
		mArch := rc.MatrixArch()

		buildEnv := sliceutils.NewStrings(
			c.CGO.getEnv(
				rc.HostKernel() != mKernel || rc.HostArch() != mArch,
				mKernel, mArch,
				rc.HostOS(),
				rc.MatrixLibc(),
			),
			createBuildEnv(mKernel, mArch)...,
		)

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

		var compileArgs []string
		compileArgs = append(compileArgs, c.Build.generateArgs(useShell)...)
		compileArgs = append(compileArgs, c.Test.generateArgs(true)...)
		compileArgs = append(compileArgs, c.Benchmark.generateArgs(true)...)
		compileArgs = append(compileArgs, c.Profile.generateArgs(dukkhaWorkingDir, true)...)

		var runArgs []string
		runArgs = append(runArgs, c.Test.generateArgs(false)...)
		runArgs = append(runArgs, c.Benchmark.generateArgs(false)...)
		runArgs = append(runArgs, c.Profile.generateArgs(dukkhaWorkingDir, false)...)
		if len(c.CustomArgs) != 0 {
			runArgs = append(runArgs, "--")
			runArgs = append(runArgs, c.CustomArgs...)
		}

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
							runArgs, pkgRelPath,
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

func generateCompileSpecs(
	taskName string,
	dukkhaCacheDir string,
	chdir string,
	env []string,
	args []string,
	pkgRelPath string,

	// options
	toolCmd []string,
	useShell bool,
	shellName string,
) (string, []dukkha.TaskExecSpec) {
	var steps []dukkha.TaskExecSpec

	builtTestExecutable := filepath.Join(
		dukkhaCacheDir, "golang-test",
		taskName+"-"+hex.EncodeToString(hashhelper.MD5Sum([]byte(pkgRelPath)))+".test",
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
		Env:         sliceutils.NewStrings(env),
		Chdir:       chdir,
		Command:     append(compileCmd, pkgRelPath),
		UseShell:    useShell,
		ShellName:   shellName,
		IgnoreError: false,
	})

	return builtTestExecutable, steps
}

func generateRunSpecs(
	dukkhaWorkingDir string,
	builtTestExecutable string,
	chdir string,
	args []string,
	pkgRelPath string,
	useShell bool,
	shellName string,
) []dukkha.TaskExecSpec {
	var steps []dukkha.TaskExecSpec

	workdir := dukkhaWorkingDir
	if len(workdir) == 0 {
		// use same default workdir as go test (the pakcage dir)
		if filepath.IsAbs(chdir) {
			workdir = filepath.Join(chdir, pkgRelPath)
		} else {
			workdir = filepath.Join(dukkhaWorkingDir, chdir, pkgRelPath)
		}
	}

	runCmd := append([]string{builtTestExecutable}, args...)

	// check if compiled test file exists
	// can be missing if no test was found in the package
	steps = append(steps, dukkha.TaskExecSpec{
		AlterExecFunc: func(
			replace dukkha.ReplaceEntries,
			stdin io.Reader,
			stdout, stderr io.Writer,
		) (dukkha.RunTaskOrRunCmd, error) {
			_, err := os.Stat(builtTestExecutable)
			if err != nil {
				if os.IsNotExist(err) {
					// no test in that package
					return nil, nil
				}

				return nil, fmt.Errorf("failed to check compiled test executable: %w", err)
			}

			// found, run it

			return []dukkha.TaskExecSpec{{
				Chdir:     workdir,
				Command:   runCmd,
				UseShell:  useShell,
				ShellName: shellName,
			}}, nil
		},
	})

	return steps
}

type testSpec struct {
	field.BaseField

	LogFile  string        `yaml:"log_file"`
	Count    int           `yaml:"count"`
	CPU      []int         `yaml:"cpu"`
	Parallel int           `yaml:"parallel"`
	FailFast bool          `yaml:"failfast"`
	Short    bool          `yaml:"short"`
	Timeout  time.Duration `yaml:"timeout"`
	Match    string        `yaml:"match"`
	Verbose  bool          `yaml:"verbose"`

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

	if s.Verbose {
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
