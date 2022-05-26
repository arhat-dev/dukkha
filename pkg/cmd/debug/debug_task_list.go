package debug

import (
	"encoding/json"
	"strings"

	"arhat.dev/pkg/textquery"
	"github.com/spf13/cobra"

	"arhat.dev/dukkha/pkg/cmd/utils"
	"arhat.dev/dukkha/pkg/dukkha"
)

func NewDebugTaskListCmd(ctx *dukkha.Context, opts *Options) *cobra.Command {
	debugTaskListCmd := &cobra.Command{
		Use:   "list",
		Short: "List task names",

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

			query, err := opts.getQuery()
			if err != nil {
				return err
			}

			stdout, stderr := appCtx.Stdout(), appCtx.Stderr()

			var enc *json.Encoder
			if query != nil {
				enc = json.NewEncoder(stdout)
				enc.SetIndent("", "  ")
			}

			var buf []string

			var (
				lastToolName dukkha.ToolName
				showToolName bool
			)

			return forEachTask(appCtx, args,
				func(appCtx dukkha.Context, tool dukkha.Tool, task dukkha.Task, i, count int) error {
					if i == 0 {
						showToolName = true
						lastToolName = ""
					}

					buf = append(buf, string(task.Name()))
					if len(lastToolName) != 0 && lastToolName != tool.Name() {
						showToolName = false
					}

					lastToolName = tool.Name()

					if i != count-1 {
						return nil
					}

					actualToolName := lastToolName
					if !showToolName {
						actualToolName = ""
					}
					err := opts.writeHeader(stdout, stderr, TaskHeaderLineData{
						ToolKind: tool.Kind(),
						ToolName: actualToolName,
						TaskKind: task.Kind(),
					}.json())
					if err != nil {
						return err
					}

					if query != nil {
						var ret []any

						tmp := make([]any, len(buf))
						for idx, v := range buf {
							tmp[idx] = v
						}

						ret, err = textquery.RunQuery(query, tmp, nil)
						if err != nil {
							return err
						}

						var data any
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

					out := `["` + strings.Join(buf, `", "`) + `"]` + "\n"
					_, err = stdout.Write([]byte(out))
					buf = make([]string, 0)
					return err
				},
			)
		},
	}

	utils.SetupTaskCompletion(ctx, debugTaskListCmd)
	debugTaskListCmd.SetHelpCommand(&cobra.Command{
		SilenceUsage: true,
		Hidden:       true,
	})

	return debugTaskListCmd
}
