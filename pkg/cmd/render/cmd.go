package render

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

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
			return run(*ctx, opts, args, os.Stdout)
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
	resolvedOpts, err := opts.Resolve(args, stdout)
	if err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	if len(args) == 0 {
		// render yaml from stdin only
		return renderYamlReader(
			appCtx,
			os.Stdin,
			resolvedOpts.OutputPathFor("-"),
			0664,
			resolvedOpts,
		)
	}

	lastWorkDir := appCtx.WorkingDir()
	for _, src := range args {
		if src == "-" {
			err = renderYamlReader(
				appCtx,
				os.Stdin,
				resolvedOpts.OutputPathFor("-"),
				0664,
				resolvedOpts,
			)
			if err != nil {
				return err
			}

			continue
		}

		err = func() error {
			// chdir at the entrypoint (root of the source yaml)
			// make relative paths in that dir happy

			chdir := resolvedOpts.ChdirFor(src)

			if chdir != lastWorkDir {
				err = os.Chdir(chdir)
				if err != nil {
					return fmt.Errorf(
						"chdir: going to source root %q: %w",
						src, err,
					)
				}

				// change DUKKHA_WORKING_DIR to make renderers like
				// `file`, `shell` and `env` work properly
				appCtx.(interface {
					OverrideWorkingDir(cwd string)
				}).OverrideWorkingDir(chdir)

				lastWorkDir = chdir
			}

			return renderYamlFile(
				appCtx,
				resolvedOpts.EntrypointFor(src),
				resolvedOpts.OutputPathFor(src),
				resolvedOpts,
				make(map[string]os.FileMode),
			)
		}()

		if err != nil {
			return err
		}
	}

	return nil
}
