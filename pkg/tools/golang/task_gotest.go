package golang

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
	"arhat.dev/pkg/hashhelper"
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

		env := sliceutils.NewStrings(
			c.CGO.getEnv(
				rc.HostKernel() != mKernel || rc.HostArch() != mArch,
				mKernel, mArch,
				rc.HostOS(),
				rc.MatrixLibc(),
			),
			createBuildEnv(mKernel, mArch)...,
		)

		// compile the test
		builtTestExecutable := filepath.Join(
			rc.CacheDir(), "golang-test",
			c.TaskName+"-"+hex.EncodeToString(hashhelper.MD5Sum([]byte(c.Path)))+".test",
		)

		// remove previously built test executable if any
		steps = append(steps, dukkha.TaskExecSpec{
			AlterExecFunc: func(
				replace map[string][]byte,
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
			options.ToolCmd(), "test", "-c", "-o", builtTestExecutable,
		)

		compileCmd = append(compileCmd, c.Build.generateArgs(options.UseShell())...)
		compileCmd = append(compileCmd, c.Test.generateArgs(true)...)
		compileCmd = append(compileCmd, c.Benchmark.generateArgs(true)...)
		compileCmd = append(compileCmd, c.Profile.generateArgs(rc.WorkingDir(), true)...)

		steps = append(steps, dukkha.TaskExecSpec{
			Env:         sliceutils.NewStrings(env, c.Env...),
			Chdir:       c.Chdir,
			Command:     append(compileCmd, c.Path),
			UseShell:    options.UseShell(),
			ShellName:   options.ShellName(),
			IgnoreError: false,
		})

		workdir := c.Test.WorkDir
		if len(workdir) == 0 {
			// use same default workdir as go test
			if filepath.IsAbs(c.Path) {
				workdir = c.Path
			} else {
				if filepath.IsAbs(c.Chdir) {
					workdir = filepath.Join(c.Chdir, c.Path)
				} else {
					workdir = filepath.Join(rc.WorkingDir(), c.Chdir, c.Path)
				}
			}
		}

		testEnv := sliceutils.NewStrings(c.Env)

		runCmd := []string{builtTestExecutable}
		runCmd = append(runCmd, c.Test.generateArgs(false)...)
		runCmd = append(runCmd, c.Benchmark.generateArgs(false)...)
		runCmd = append(runCmd, c.Profile.generateArgs(rc.WorkingDir(), false)...)
		if len(c.CustomArgs) != 0 {
			runCmd = append(runCmd, "--")
			runCmd = append(runCmd, c.CustomArgs...)
		}

		// check if compiled test file exists
		// can be missing if no test was found in the package
		steps = append(steps, dukkha.TaskExecSpec{
			AlterExecFunc: func(
				replace map[string][]byte,
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

				return []dukkha.TaskExecSpec{
					{
						Env:       testEnv,
						Chdir:     workdir,
						Command:   runCmd,
						UseShell:  options.UseShell(),
						ShellName: options.ShellName(),
					},
				}, nil
			},
		})

		return nil
	})

	return steps, err
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
