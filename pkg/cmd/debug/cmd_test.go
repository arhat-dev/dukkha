package debug

import (
	"context"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"

	"arhat.dev/pkg/testhelper"
	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/conf"
	"arhat.dev/dukkha/pkg/dukkha"
	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/matrix"
	"arhat.dev/dukkha/pkg/sliceutils"

	_ "arhat.dev/dukkha/cmd/dukkha/addon"
)

type TestSpec struct {
	Options []string `yaml:"options"`
	Args    []string `yaml:"args"`
}

type CheckSpec struct {
	Stdout    string `yaml:"stdout"`
	Stderr    string `yaml:"stderr"`
	BadFlags  bool   `yaml:"bad_flags"`
	ExpectErr bool   `yaml:"expect_err"`
}

const (
	OmitThisOption  = "<OMIT_THIS_OPTION>"
	OmitOptionValue = "<OMIT_OPTION_VALUE>"
)

func testCmdFlags(
	t *testing.T,
	dir string,
	createTestSpec func() interface{},
	createExpected func() interface{},
	optMat map[string][]string,
	genNewSpec func(opt [][]string, baseSpec, baseCheck interface{}) (_, _ interface{}),
	check func(t *testing.T, spec interface{}, exp interface{}),
) {
	testhelper.TestFixtures(t, dir, createTestSpec, createExpected, func(t *testing.T, spec, exp interface{}) {
		optMats := matrix.CartesianProduct(optMat)
		for _, m := range optMats {
			var ent [][]string
			for k, v := range m {
				switch v {
				case OmitOptionValue:
					ent = append(ent, []string{k})
				case OmitThisOption:
					continue
				default:
					ent = append(ent, []string{k, v})
				}
			}

			specNew, expNew := genNewSpec(ent, spec, exp)

			var data []string
			for _, o := range ent {
				data = append(data, strings.Join(o, "="))
			}
			t.Run(strings.Join(data, ","), func(t *testing.T) {
				check(t, specNew, expNew)
			})
		}
	})
}

func TestDebugCmd(t *testing.T) {
	var (
		matchDefaultHeaderLine = regexp.MustCompile(`(?m:^--- # .*\n)`)
	)

	testCmdFlags(t, "./fixtures",
		func() interface{} { return &TestSpec{} },
		func() interface{} { return &CheckSpec{} },
		map[string][]string{
			"-H": {
				OmitThisOption,
				OmitOptionValue,
			},
			"-P": {
				OmitThisOption,
				"TEST_PREFIX",
			},
		},
		func(extraPairs [][]string, baseSpec, baseCheck interface{}) (_ interface{}, _ interface{}) {
			var (
				in       = baseSpec.(*TestSpec)
				expected = baseCheck.(*CheckSpec)

				stdout    = expected.Stdout
				stderr    = expected.Stderr
				badFlags  = false
				expectErr = expected.ExpectErr
			)

			var opts []string

			newHeaderPrefix := ""
			for _, p := range extraPairs {
				optName, optVal := p[0], ""
				if len(p) == 2 {
					optVal = p[1]
				}

				switch optName {
				case "-H":
					switch optVal {
					case "":
						count := len(matchDefaultHeaderLine.FindAllStringIndex(stdout, -1))
						if len(stderr) == 0 {
							stderr = strings.Join(
								matchDefaultHeaderLine.FindAllString(stdout, -1),
								"",
							)

							for i := 0; i < count; i++ {

							}
						}

						stdout = matchDefaultHeaderLine.ReplaceAllLiteralString(stdout, "")
					default:
					}
				case "-P":
					switch optVal {
					case "":
						badFlags = true
					default:
						newHeaderPrefix = optVal
					}
				}

				opts = append(opts, p...)
			}

			if len(newHeaderPrefix) != 0 {
				stdout = strings.ReplaceAll(stdout, defaultHeaderPrefix, newHeaderPrefix)
				stderr = strings.ReplaceAll(stderr, defaultHeaderPrefix, newHeaderPrefix)
			}

			return &TestSpec{
					Options: sliceutils.NewStrings(in.Options, opts...),
					Args:    sliceutils.NewStrings(in.Args),
				}, &CheckSpec{
					Stdout:    stdout,
					Stderr:    stderr,
					BadFlags:  badFlags,
					ExpectErr: expectErr,
				}
		},
		func(t *testing.T, spec, exp interface{}) {
			testDebugCmd(t, spec.(*TestSpec), exp.(*CheckSpec))
		},
	)
}

func testDebugCmd(t *testing.T, test *TestSpec, check *CheckSpec) {
	ctx := dukkha_test.NewTestContext(t, context.TODO())

	config := conf.NewConfig()
	err := conf.ReadConfigRecursively(
		os.DirFS("./testdata"),
		[]string{"."},
		false,
		&map[string]struct{}{},
		config,
	)
	if !assert.NoError(t, err) {
		return
	}

	if !assert.NoError(t, config.Resolve(ctx, true)) {
		return
	}

	appCtx := ctx.(dukkha.Context)
	debugCmd, opts := NewDebugCmd(&appCtx)

	debugTaskCmd := NewDebugTaskCmd(&appCtx, opts)
	debugTaskCmd.AddCommand(
		NewDebugTaskListCmd(&appCtx, opts),
		NewDebugTaskMatrixCmd(&appCtx, opts),
		NewDebugTaskSpecCmd(&appCtx, opts),
	)

	debugCmd.AddCommand(debugTaskCmd)

	err = debugCmd.ParseFlags(sliceutils.NewStrings(test.Options))
	if check.BadFlags {
		assert.Error(t, err)
		return
	}

	debugCmd.SetArgs(sliceutils.NewStrings(test.Options, test.Args...))
	assert.NoError(t, err)

	readStdout, stdout, err := os.Pipe()
	if !assert.NoError(t, err) {
		return
	}

	readStderr, stderr, err := os.Pipe()
	if !assert.NoError(t, err) {
		return
	}

	wg := new(sync.WaitGroup)

	var (
		stdoutData []byte
		stderrData []byte
	)

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			readStdout.Close()
		}()

		var err2 error
		stdoutData, err2 = io.ReadAll(readStdout)
		assert.NoError(t, err2)
	}()

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			readStderr.Close()
		}()

		var err2 error
		stderrData, err2 = io.ReadAll(readStderr)
		assert.NoError(t, err2)
	}()

	testhelper.HijackStandardStreams(nil, stdout, stderr, func() {
		defer func() {
			_ = stdout.Close()
			_ = stderr.Close()
		}()

		err = debugCmd.Execute()
	})

	if check.ExpectErr {
		assert.NoError(t, err)
		return
	}

	assert.NoError(t, err)

	wg.Wait()

	assert.Equal(t, check.Stdout, string(stdoutData))
	assert.Equal(t, check.Stderr, string(stderrData))
}
