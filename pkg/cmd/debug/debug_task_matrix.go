package debug

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/itchyny/gojq"
	"github.com/spf13/cobra"

	"arhat.dev/pkg/textquery"

	"arhat.dev/dukkha/pkg/cmd/utils"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
)

func NewDebugTaskMatrixCmd(ctx *dukkha.Context) *cobra.Command {
	var (
		matrixFilter   []string
		headerToStderr bool
		queryStr       string
	)

	debugTaskMatrixCmd := &cobra.Command{
		Use:   "matrix",
		Short: "Show task matrix at runtime",

		Args:          cobra.RangeArgs(0, 4),
		SilenceErrors: true,
		SilenceUsage:  true,

		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   false,
			DisableNoDescFlag:   false,
			DisableDescriptions: true,
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			appCtx := *ctx
			appCtx = appCtx.DeriveNew()
			appCtx.SetMatrixFilter(utils.ParseMatrixFilter(matrixFilter))

			buf := &bytes.Buffer{}
			var (
				query *gojq.Query
				enc   *json.Encoder
			)

			if len(queryStr) != 0 {
				var err error
				query, err = gojq.Parse(queryStr)
				if err != nil {
					return fmt.Errorf("invalid query %q: %w", queryStr, err)
				}

				enc = json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
			}

			return debugTasks(appCtx, args,
				func(appCtx dukkha.Context, tool dukkha.Tool, task dukkha.Task) error {
					matrixSpecs, err := task.GetMatrixSpecs(appCtx)
					if err != nil {
						return fmt.Errorf("failed to get task matrix specs: %w", err)
					}

					headerOut := os.Stdout
					if headerToStderr {
						headerOut = os.Stderr
					}

					fmt.Fprintln(headerOut, `--- # { "task": "`+
						tool.Key().String()+":"+task.Key().String()+`" }`,
					)

					if query != nil {
						var ret []interface{}
						for _, ms := range matrixSpecs {
							ent := make(map[string]interface{})
							for k, v := range ms {
								ent[k] = v
							}
							ret = append(ret, ent)
						}

						ret, _, err = textquery.RunQuery(query, ret, nil)
						if err != nil {
							return err
						}

						var data interface{}
						switch len(ret) {
						case 0:
							data = nil
						case 1:
							data = ret[0]
						default:
							data = ret
						}

						return enc.Encode(data)
					}

					// write matrix directly

					buf.WriteString("[")
					for i, ms := range matrixSpecs {
						if i != 0 {
							buf.WriteString(",")
						}
						buf.WriteString("\n" + `  { "`)
						buf.WriteString(strings.Join(sliceutils.FormatStringMap(ms, `": "`, false), `", "`))
						buf.WriteString(`" }`)
					}

					if len(matrixSpecs) != 0 {
						buf.WriteString("\n]\n")
					} else {
						buf.WriteString("]\n")
					}

					_, err = os.Stdout.ReadFrom(buf)
					return err
				},
			)
		},
	}

	flags := debugTaskMatrixCmd.Flags()
	utils.RegisterMatrixFilterFlag(flags, &matrixFilter)

	flags.BoolVarP(&headerToStderr, "header-to-stderr", "H", false,
		"write document header (`--- # { \"name\":...`) to stderr (helpful for json parsing)",
	)

	flags.StringVarP(&queryStr, "query", "q", "",
		"use jq query to filter output",
	)

	err := utils.SetupTaskAndTaskMatrixCompletion(ctx, debugTaskMatrixCmd)
	if err != nil {
		panic(err)
	}

	debugTaskMatrixCmd.SetHelpCommand(&cobra.Command{
		SilenceUsage: true,
		Hidden:       true,
	})

	return debugTaskMatrixCmd
}
