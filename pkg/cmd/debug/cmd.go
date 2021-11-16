package debug

import (
	"fmt"
	"os"

	"github.com/itchyny/gojq"
	"github.com/spf13/cobra"

	"arhat.dev/dukkha/pkg/dukkha"
)

const (
	defaultHeaderPrefix = "--- # "
)

// Options cli options available to all debug commands
type Options struct {
	headerToStderr *bool
	headerPrefix   *string
	query          *string
}

func (opts *Options) writeHeader(data string) error {
	out := os.Stdout
	if *opts.headerToStderr {
		out = os.Stderr
	}

	if len(*opts.headerPrefix) != 0 {
		_, err := out.Write([]byte(*opts.headerPrefix))
		if err != nil {
			return err
		}
	}

	_, err := out.WriteString(data + "\n")
	return err
}

func (opts *Options) getQuery() (*gojq.Query, error) {
	if len(*opts.query) != 0 {
		var err error
		query, err := gojq.Parse(*opts.query)
		if err != nil {
			return nil, fmt.Errorf("invalid query %q: %w", *opts.query, err)
		}

		return query, err
	}

	return nil, nil
}

func NewDebugCmd(ctx *dukkha.Context) (*cobra.Command, *Options) {
	debugCmd := &cobra.Command{
		Use:           "debug",
		Short:         "Debug config and task definitions",
		SilenceErrors: true,
		SilenceUsage:  true,

		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   false,
			DisableNoDescFlag:   false,
			DisableDescriptions: true,
		},
	}

	return debugCmd, &Options{
		headerToStderr: debugCmd.PersistentFlags().BoolP(
			"header-to-stderr", "H",
			false,
			"write document header (`--- # { \"name\":...`) to stderr (helpful for json parsing)",
		),
		headerPrefix: debugCmd.PersistentFlags().StringP(
			"header-prefix", "P",
			defaultHeaderPrefix,
			"set custom prefix to header line",
		),
		query: debugCmd.PersistentFlags().StringP(
			"query", "q",
			"",
			"use jq query to filter output",
		),
	}
}
