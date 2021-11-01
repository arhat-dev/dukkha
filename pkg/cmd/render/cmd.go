package render

import (
	"fmt"
	"io"
	"os"

	"github.com/itchyny/gojq"
	"github.com/spf13/cobra"

	"arhat.dev/dukkha/pkg/dukkha"
)

func NewRenderCmd(ctx *dukkha.Context) *cobra.Command {
	var (
		outputFormat string
		indentSize   int
		indentStyle  string
		recursive    bool
		resultQuery  string

		outputDests []string
	)

	renderCmd := &cobra.Command{
		Use:           "render",
		Short:         "Render your yaml docs",
		Args:          cobra.ArbitraryArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   false,
			DisableNoDescFlag:   false,
			DisableDescriptions: true,
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			var indentStr string
			switch indentStyle {
			case "space":
				indentStr = " "
			case "tab":
				indentStr = "\t"
			default:
				indentStr = indentStyle
			}

			var query *gojq.Query
			if len(resultQuery) != 0 {
				var err error
				query, err = gojq.Parse(resultQuery)
				if err != nil {
					return fmt.Errorf("invalid result query: %w", err)
				}
			}

			var om map[string]*string
			if len(outputDests) != 0 {
				if len(outputDests) != len(args) {
					return fmt.Errorf(
						"number of output destination not matching sources: want %d, got %d",
						len(args), len(outputDests),
					)
				}

				om = make(map[string]*string)
				for i := range outputDests {
					src := args[i]

					om[src] = &outputDests[i]
				}
			} else {
				om = make(map[string]*string)
				for _, src := range args {
					om[src] = nil
				}
			}

			var stdoutEnc encoder

			createEncoder := func(w io.Writer) (encoder, error) {
				if w == os.Stdout {
					var err error
					if stdoutEnc == nil {
						stdoutEnc, err = newEncoder(
							query, os.Stdout,
							outputFormat, indentStr, indentSize,
						)
					}

					return stdoutEnc, err
				}

				return newEncoder(query, w, outputFormat, indentStr, indentSize)
			}

			for _, src := range args {
				err := renderYamlFileOrDir(
					*ctx, src, om[src], outputFormat,
					createEncoder,
					recursive,
					make(map[string]os.FileMode),
				)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	flags := renderCmd.Flags()
	flags.StringSliceVarP(&outputDests, "output", "o", nil,
		"set output destionation for specified inputs (args)",
	)
	flags.StringVarP(&outputFormat, "output-format", "f", "yaml",
		"set output format, one of [json, yaml]",
	)
	flags.IntVarP(&indentSize, "indent-size", "n", 2,
		"set indent size",
	)
	flags.StringVarP(&indentStyle, "indent-style", "s", "space",
		"set indent style, custom string or one of [space, tab]",
	)
	flags.BoolVarP(&recursive, "recursive", "r", false,
		"render directories recursively",
	)
	flags.StringVarP(&resultQuery, "query", "q", "",
		"run jq style query over generated yaml/json docs before writing",
	)

	return renderCmd
}
