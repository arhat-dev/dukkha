package debug

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/itchyny/gojq"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"arhat.dev/pkg/textquery"

	"arhat.dev/dukkha/pkg/cmd/utils"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

func NewDebugTaskSpecCmd(ctx *dukkha.Context) *cobra.Command {
	var (
		matrixFilter   []string
		headerToStderr bool
		queryStr       string
	)

	debugTaskSpecCmd := &cobra.Command{
		Use:   "spec",
		Short: "Show task spec at runtime",
		Long: "Tasks may contain step by step jobs (e.g. actions in hooks), " +
			"if one such step is using rendering suffix, its value will not be printed",

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

			var query *gojq.Query
			if len(queryStr) != 0 {
				var err error
				query, err = gojq.Parse(queryStr)
				if err != nil {
					return fmt.Errorf("invalid query %q: %w", queryStr, err)
				}
			}

			return debugTasks(appCtx, args,
				func(appCtx dukkha.Context, tool dukkha.Tool, task dukkha.Task) error {
					matrixSpecs, err := task.GetMatrixSpecs(appCtx)
					if err != nil {
						return fmt.Errorf("failed to get task matrix specs: %w", err)
					}

					execOpts := dukkha.CreateTaskExecOptions(0, len(matrixSpecs))
					tskCtx := appCtx.DeriveNew()
					tskCtx.SetTask(tool.Key(), task.Key())

					for _, ms := range matrixSpecs {
						mCtx, mExecOpts, err2 := tools.CreateTaskMatrixContext(
							&tools.TaskExecRequest{
								Context: tskCtx,
								Tool:    tool,
								Task:    task,
							},
							ms, execOpts,
						)
						_ = mExecOpts
						if err2 != nil {
							return fmt.Errorf("failed to create task matrix context: %w", err2)
						}

						err2 = task.DoAfterFieldsResolved(mCtx, -1, true, func() error {
							buf := &bytes.Buffer{}
							enc := yaml.NewEncoder(buf)
							enc.SetIndent(2)
							defer func() { _ = enc.Close() }()

							dec := yaml.NewDecoder(buf)

							err = enc.Encode(task)
							if err != nil {
								return err
							}

							var data interface{}
							err = dec.Decode(&data)
							if err != nil {
								return err
							}

							if query != nil {
								var ret []interface{}
								ret, _, err = textquery.RunQuery(query, data, nil)
								if err != nil {
									return err
								}

								switch len(ret) {
								case 0:
									data = nil
								case 1:
									data = ret[0]
								default:
									data = ret
								}
							}

							jenc := json.NewEncoder(buf)
							jenc.SetIndent("", "  ")
							err = jenc.Encode(data)
							if err != nil {
								return err
							}

							headerOut := os.Stdout
							if headerToStderr {
								headerOut = os.Stderr
							}
							fmt.Fprintln(headerOut, `--- # { "task": "`+
								tool.Key().String()+":"+task.Key().String()+
								`", "matrix": { "`+
								strings.Join(sliceutils.FormatStringMap(ms, `": "`, false), `", "`)+
								`" } }`,
							)

							_, err = os.Stdout.ReadFrom(buf)
							return err
						})
						if err2 != nil {
							return fmt.Errorf("failed to generate task spec: %w", err2)
						}
					}

					return nil
				},
			)
		},
	}

	flags := debugTaskSpecCmd.Flags()
	utils.RegisterMatrixFilterFlag(flags, &matrixFilter)

	flags.BoolVarP(&headerToStderr, "header-to-stderr", "H", false,
		"write yaml doc separator (`--- # { \"name\":...`) to stderr (helpful for json parsing)",
	)

	flags.StringVarP(&queryStr, "query", "q", "",
		"use jq query to filter output",
	)

	err := utils.SetupTaskAndTaskMatrixCompletion(ctx, debugTaskSpecCmd)
	if err != nil {
		panic(err)
	}

	debugTaskSpecCmd.SetHelpCommand(&cobra.Command{
		SilenceUsage: true,
		Hidden:       true,
	})

	return debugTaskSpecCmd
}
