package render

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

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
			return run(*ctx, opts, args)
		},
	}

	flags := renderCmd.Flags()
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

	return renderCmd
}

func run(appCtx dukkha.Context, opts *Options, args []string) error {
	resolvedOpts, err := opts.Resolve(args, os.Stdout)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		// render yaml from stdin only
		return renderYamlReader(
			appCtx,
			os.Stdin,
			resolvedOpts.outputMapping["-"],
			0664,
			resolvedOpts,
		)
	}

	originalWorkDir := appCtx.WorkingDir()
	for _, src := range args {
		if src == "-" {
			err = renderYamlReader(
				appCtx,
				os.Stdin,
				resolvedOpts.outputMapping["-"],
				0664,
				resolvedOpts,
			)
			if err != nil {
				return err
			}

			continue
		}

		err = func() error {
			var info fs.FileInfo
			info, err = os.Stat(src)
			if err != nil {
				return err
			}

			// default values when src is a dir
			var (
				chdir      string
				entrypoint string
			)
			if info.IsDir() {
				// we are going to chdir into src dicrectory, or we
				// have already been there
				//
				// so the entrypoint is always current dir
				entrypoint = "."
				chdir = src
			} else {
				// regular file
				entrypoint = src
				chdir = filepath.Dir(chdir)
			}

			chdir, err = filepath.Abs(chdir)
			if err != nil {
				return err
			}

			// chdir at the entrypoint (root of the source yaml)
			// make relative paths in that dir happy

			if chdir != originalWorkDir {
				err = os.Chdir(chdir)
				if err != nil {
					return fmt.Errorf(
						"chdir: going to source root %q: %w",
						src, err,
					)
				}

				// change DUKKHA_WORKING_DIR to make renderers like `shell`, `env` happy
				appCtx.(interface {
					OverrideWorkingDir(cwd string)
				}).OverrideWorkingDir(chdir)

				// always chdir back to original working dir, since other
				// input source can be relative path to the original working dir
				defer func() {
					appCtx.(interface {
						OverrideWorkingDir(cwd string)
					}).OverrideWorkingDir(originalWorkDir)

					err = os.Chdir(originalWorkDir)
					if err != nil {
						panic(fmt.Errorf(
							"failed to go back to dukkha working dir: %w", err,
						))
					}
				}()
			}

			return renderYamlFile(
				appCtx,
				entrypoint,
				resolvedOpts.outputMapping[src],
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
