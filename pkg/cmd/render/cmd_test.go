package render

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"arhat.dev/pkg/testhelper"
	"arhat.dev/rs"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	dt "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/renderer/env"
	"arhat.dev/dukkha/pkg/renderer/file"
	"arhat.dev/dukkha/pkg/renderer/shell"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/utils"
)

func TestCmd(t *testing.T) {

	type TestSpec struct {
		rs.BaseField

		Options    []string `yaml:"options"`
		BadOptions bool     `yaml:"bad_options"`

		Sources   []string `yaml:"sources"`
		BadSource bool     `yaml:"bad_source"`

		OutputFiles []string `yaml:"output_files"`
	}

	type CheckSpec struct {
		rs.BaseField

		Expected string `yaml:"expected"`
	}

	testhelper.TestFixtures(t, "./fixtures",
		func() interface{} { return rs.Init(&TestSpec{}, nil) },
		func() interface{} { return rs.Init(&CheckSpec{}, nil) },
		func(t *testing.T, in, c interface{}) {
			spec := in.(*TestSpec)
			check := c.(*CheckSpec)

			cwd, err := os.Getwd()
			if !assert.NoError(t, err) {
				return
			}

			defer t.Cleanup(func() {
				if !assert.NoError(t, os.Chdir(cwd)) {
					return
				}
			})

			ctx := dt.NewTestContextWithGlobalEnv(context.TODO(), map[string]utils.LazyValue{
				constant.ENV_DUKKHA_WORKING_DIR: utils.ImmediateString(cwd),
			})
			ctx.(di.CacheDirSetter).SetCacheDir(t.TempDir())

			ctx.AddRenderer("file", file.NewDefault("file"))
			ctx.AddRenderer("env", env.NewDefault("env"))
			ctx.AddRenderer("shell", shell.NewDefault("shell"))
			outputDir := filepath.Join(t.TempDir(), "output")
			ctx.AddEnv(true, &dukkha.EnvEntry{
				Name:  "out_dir",
				Value: outputDir,
			})

			assert.NoError(t, os.MkdirAll(outputDir, 0755))

			appCtx := dukkha.Context(ctx)
			cmd := &cobra.Command{}
			opts := &Options{}
			createOptionsFlags(cmd, opts)

			assert.NoError(t, spec.ResolveFields(ctx, -1))
			err = cmd.ParseFlags(sliceutils.NewStrings(spec.Options, spec.Sources...))
			if spec.BadOptions {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			stdout := &bytes.Buffer{}
			err = run(appCtx, opts, spec.Sources, stdout)

			if spec.BadSource {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if !assert.NoError(t, os.Chdir(cwd)) {
				return
			}

			ctx.(di.WorkingDirOverrider).OverrideWorkingDir(cwd)

			// have a look at the output dir
			entries, err := os.ReadDir(outputDir)
			assert.NoError(t, err)
			t.Log(len(entries))
			for _, v := range entries {
				t.Log(v.Name())
			}

			var actual []*rs.AnyObject
			for _, file := range spec.OutputFiles {
				switch file {
				case "-":
					actualPart, err2 := parseYaml(yaml.NewDecoder(stdout))

					assert.NoError(t, err2)

					actual = append(actual, actualPart...)
				default:
					f, err2 := os.Open(file)
					if !assert.NoError(t, err2) {
						return
					}

					actualPart, err2 := parseYaml(yaml.NewDecoder(f))
					_ = f.Close()

					assert.NoError(t, err2)

					actual = append(actual, actualPart...)
				}
			}

			expected, err := parseYaml(yaml.NewDecoder(strings.NewReader(check.Expected)))
			assert.NoError(t, err)

			if !assert.Equal(t, len(expected), len(actual)) {
				return
			}

			for i, v := range expected {
				assert.NoError(t, v.ResolveFields(ctx, -1))

				assert.EqualValues(t, v.NormalizedValue(), actual[i].NormalizedValue())
			}
		},
	)
}

// func testCmd(t *testing.T) {
// 	testhelper.TestCmdFixtures(t, "./fixtures",
// 		map[string][]string{},
// 		generateNewSpec,
// 		prepareCmd,
// 	)
// }
//
// func generateNewSpec(
// 	flagSets [][]string,
// 	baseSpec *testhelper.CmdTestCase,
// 	baseCheck *testhelper.CmdTestCheckSpec,
// ) (*testhelper.CmdTestCase, *testhelper.CmdTestCheckSpec) {
// 	return nil, nil
// }
//
// func prepareCmd(flags []string) (checkFlags func() error, runCmd func() error, _ error) {
// 	return func() error { return nil }, func() error { return nil }, nil
// }
