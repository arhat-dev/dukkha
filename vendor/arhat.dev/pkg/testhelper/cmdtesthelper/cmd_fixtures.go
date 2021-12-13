package cmdtesthelper

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"

	"arhat.dev/pkg/matrixhelper"
	"arhat.dev/pkg/testhelper"
)

const (
	OmitThisFlag  = "<OMIT_THIS_FLAG>"
	OmitFlagValue = "<OMIT_FLAG_VALUE>"
)

type CmdTestCase struct {
	rs.BaseField

	// Flags is the command line flags to cmd, including positional args
	Flags []string `yaml:"flags"`

	// Stdin is the input to cmd stdin
	Stdin string `yaml:"stdin"`
}

type CmdTestCheckSpec struct {
	rs.BaseField

	// BadFlags expects optoins and/or args are invalid
	BadFlags bool `yaml:"bad_flags"`

	// BadCmd expects cmd run error
	//
	// after parsing all flags and before checking output
	BadCmd bool `yaml:"bad_cmd"`

	// Stdout is the expected output to stdout
	Stdout string `yaml:"stdout"`

	// Stderr is the expected output to stderr
	Stderr string `yaml:"stderr"`
}

func TestCmdFixtures(t *testing.T,
	dir string,
	flagMatrix map[string][]string,
	genNewSpec func(
		flagSets [][]string, baseSpec *CmdTestCase, baseCheck *CmdTestCheckSpec,
	) (*CmdTestCase, *CmdTestCheckSpec),
	prepareRun func(flags []string) (checkFlags func() error, runCmd func() error, _ error),
) {
	testhelper.TestFixtures(t, dir,
		func() interface{} { return rs.Init(&CmdTestCase{}, nil) },
		func() interface{} { return rs.Init(&CmdTestCheckSpec{}, nil) },
		func(t *testing.T, spec, exp interface{}) {
			flagMats := matrixhelper.CartesianProduct(flagMatrix)
			for _, m := range flagMats {
				var flagSets [][]string
				for k, v := range m {
					switch v {
					case OmitFlagValue:
						flagSets = append(flagSets, []string{k})
					case OmitThisFlag:
						continue
					default:
						flagSets = append(flagSets, []string{k, v})
					}
				}

				specNew, expNew := genNewSpec(flagSets, spec.(*CmdTestCase), exp.(*CmdTestCheckSpec))

				var data []string
				for _, o := range flagSets {
					data = append(data, strings.Join(o, "="))
				}

				t.Run(strings.Join(data, ","), func(t *testing.T) {
					checkFlags, runCmd, err := prepareRun(specNew.Flags)
					if !assert.NoError(t, err) {
						return
					}

					err = checkFlags()
					if expNew.BadFlags {
						assert.Error(t, err)
						return
					}

					if !assert.NoError(t, err) {
						return
					}

					readStdout, stdout, err := os.Pipe()
					if !assert.NoError(t, err) {
						return
					}

					readStderr, stderr, err := os.Pipe()
					if !assert.NoError(t, err) {
						return
					}

					stdin, writeStdin, err := os.Pipe()
					if !assert.NoError(t, err) {
						return
					}

					var (
						stdoutData []byte
						stderrData []byte

						wg = new(sync.WaitGroup)
					)

					wg.Add(3)

					go func() {
						defer func() {
							wg.Done()
							readStdout.Close()
						}()

						var err2 error
						stdoutData, err2 = io.ReadAll(readStdout)
						assert.NoError(t, err2)
					}()

					go func() {
						defer func() {
							wg.Done()
							readStderr.Close()
						}()

						var err2 error
						stderrData, err2 = io.ReadAll(readStderr)
						assert.NoError(t, err2)
					}()

					go func() {
						defer func() {
							wg.Done()
							writeStdin.Close()
						}()

						n, err2 := io.Copy(writeStdin, bytes.NewReader([]byte(specNew.Stdin)))
						assert.NoError(t, err2)
						assert.EqualValues(t, len(specNew.Stdin), n)
					}()

					testhelper.HijackStandardStreams(stdin, stdout, stderr, func() {
						defer func() {
							_ = stdout.Close()
							_ = stderr.Close()
							_ = stdin.Close()
						}()

						err = runCmd()
					})

					if expNew.BadCmd {
						assert.Error(t, err)
					} else {
						assert.NoError(t, err)
					}

					wg.Wait()

					assert.Equal(t, expNew.Stdout, string(stdoutData))
					assert.Equal(t, expNew.Stderr, string(stderrData))
				})
			}
		},
	)
}
