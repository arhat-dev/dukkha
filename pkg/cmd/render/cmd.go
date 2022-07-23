package render

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/spf13/cobra"

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/dukkha"
)

func NewRenderCmd(ctx *dukkha.Context) *cobra.Command {
	// cli options

	opts := &Options{}

	renderCmd := &cobra.Command{
		Use:           "render",
		Short:         "Render yaml docs using rendering suffix",
		Args:          cobra.ArbitraryArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   false,
			DisableNoDescFlag:   false,
			DisableDescriptions: true,
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return run(*ctx, opts, args, (*ctx).Stdout())
		},
	}

	createOptionsFlags(renderCmd, opts)
	return renderCmd
}

func createOptionsFlags(cmd *cobra.Command, opts *Options) {
	flags := cmd.Flags()
	flags.StringSliceVarP(&opts.outputDests, "output", "o", nil,
		"set output destionation for specified inputs (args)",
	)
	flags.StringVarP(&opts.outputFormat, "output-format", "f", "yaml",
		"set output format, one of [json, yaml]",
	)
	flags.IntVarP(&opts.indentSize, "indent-size", "n", 2,
		"set indent size",
	)
	flags.StringVarP(&opts.indentStyle, "indent-style", "s", "space",
		"set indent style, custom string or one of [space, tab]",
	)
	flags.BoolVarP(&opts.recursive, "recursive", "r", false,
		"render directories recursively",
	)
	flags.StringVarP(&opts.resultQuery, "query", "q", "",
		"run jq style query over generated yaml/json docs before writing",
	)
	flags.StringSliceVar(&opts.chdir, "chdir", nil,
		"set root of the soure for specified inputs (args) for relative path resovling, "+
			"useful when you are rendering single file inside some child directory of the source directory",
	)
}

func run(appCtx dukkha.Context, opts *Options, args []string, stdout io.Writer) error {
	if len(args) == 0 {
		// defaults to read from stdin
		args = append(args, "-")
	}

	resolvedOpts, err := opts.Resolve(appCtx.FS(), args, stdout)
	if err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	stdin := appCtx.Stdin()
	lastWorkDir := appCtx.WorkDir()
	for _, src := range args {
		if src == "-" {
			err = renderYamlReader(
				appCtx,
				stdin,
				resolvedOpts.OutputPathFor("-"),
				0664,
				resolvedOpts,
			)
			if err != nil {
				return err
			}

			continue
		}

		// chdir at the entrypoint (root of the source yaml)
		// make relative paths in that dir happy

		chdir := resolvedOpts.ChdirFor(src)

		if chdir != lastWorkDir {
			// change DUKKHA_WORKDIR to make renderers like
			// `file`, `shell` and `env` work properly
			appCtx.(di.WorkDirOverrider).OverrideWorkDir(chdir)

			lastWorkDir = chdir
		}

		err = renderYamlFile(
			appCtx,
			resolvedOpts.EntrypointFor(src),
			resolvedOpts.OutputPathFor(src),
			resolvedOpts,
			make(map[string]fs.FileMode),
		)

		if err != nil {
			return err
		}
	}

	return nil
}
