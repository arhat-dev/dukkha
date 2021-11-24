package debug

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"arhat.dev/pkg/textquery"
	"github.com/spf13/cobra"

	"arhat.dev/dukkha/pkg/cmd/utils"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/matrix"
	"arhat.dev/dukkha/pkg/sliceutils"
)

func NewDebugTaskMatrixCmd(ctx *dukkha.Context, opts *Options) *cobra.Command {
	var (
		matrixFilter []string
	)

	debugTaskMatrixCmd := &cobra.Command{
		Use:   "matrix",
		Short: "List task matrix entries used at runtime",

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
			appCtx.SetMatrixFilter(matrix.ParseMatrixFilter(matrixFilter))

			query, err := opts.getQuery()
			if err != nil {
				return err
			}

			var enc *json.Encoder
			if query != nil {
				enc = json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
			}

			buf := &bytes.Buffer{}
			return forEachTask(appCtx, args,
				func(appCtx dukkha.Context, tool dukkha.Tool, task dukkha.Task, _, _ int) error {
					matrixSpecs, err := task.GetMatrixSpecs(appCtx)
					if err != nil {
						return fmt.Errorf("failed to get task matrix specs: %w", err)
					}

					err = opts.writeHeader(TaskHeaderLineData{
						ToolKind: tool.Kind(),
						ToolName: tool.Name(),
						TaskKind: task.Kind(),
						TaskName: task.Name(),
					}.json())
					if err != nil {
						return err
					}

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
