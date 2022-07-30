package debug

import (
	"context"
	"os"
	"regexp"
	"strings"
	"testing"

	"arhat.dev/pkg/synchain"
	"arhat.dev/pkg/testhelper/cmdtesthelper"

	"arhat.dev/dukkha/pkg/conf"
	"arhat.dev/dukkha/pkg/dukkha"
	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/sliceutils"

	_ "arhat.dev/dukkha/cmd/dukkha/addon"
)

func TestDebugCmd(t *testing.T) {
	t.Parallel()

	tmpdir := t.TempDir()
	cmdtesthelper.TestCmdFixtures(t, "./fixtures", map[string][]string{
		"-H": {
			cmdtesthelper.OmitThisFlag,
			cmdtesthelper.OmitFlagValue,
		},
		"-P": {
			cmdtesthelper.OmitThisFlag,
			"TEST_PREFIX",
		},
	}, genNewDebugCmdFlags, func(flags []string) (checkFlags func() error, runCmd func() error, _ error) {
		return prepareDebugCmd(flags, tmpdir)
	})
}

var (
	matchDefaultHeaderLine = regexp.MustCompile(`(?m:^--- # .*\n)`)
)

func genNewDebugCmdFlags(
	flagSets [][]string,
	baseSpec *cmdtesthelper.CmdTestCase,
	baseCheck *cmdtesthelper.CmdTestCheckSpec,
) (*cmdtesthelper.CmdTestCase, *cmdtesthelper.CmdTestCheckSpec) {
	var (
		stdout   = baseCheck.Stdout
		stderr   = baseCheck.Stderr
		badFlags = false
		badCmd   = baseCheck.BadCmd
	)

	var opts []string

	newHeaderPrefix := ""
	for _, p := range flagSets {
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

	return &cmdtesthelper.CmdTestCase{
			Flags: sliceutils.NewStrings(baseSpec.Flags, opts...),
		}, &cmdtesthelper.CmdTestCheckSpec{
			Stdout:   stdout,
			Stderr:   stderr,
			BadFlags: badFlags,
			BadCmd:   badCmd,
		}
}

func prepareDebugCmd(flags []string, cacheDir string) (checkFlags func() error, runCmd func() error, _ error) {
	ctx := dukkha_test.NewTestContext(context.TODO(), cacheDir)

	config := conf.NewConfig()
	err := conf.Read(
		ctx,
		&conf.ReadSpec{
			Flags:        conf.ReadFlag_Full,
			ConfFS:       os.DirFS("./testdata"),
			VisitedPaths: &map[string]struct{}{},
			MergedConfig: config,
		},
		synchain.NewSynchain(),
		[]string{"."},
		false,
	)
	if err != nil {
		panic(err)
	}

	err = config.Resolve(ctx, true)
	if err != nil {
		panic(err)
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
	debugCmd.SetArgs(flags)
	return func() error {
		// TODO: test bad flags, currently always return nil due to we want
		// 		 sub command Flags() has PersistentFlags() from debugCmd
		return nil
	}, debugCmd.Execute, nil
}
