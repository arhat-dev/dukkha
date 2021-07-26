package golang

import (
	"strings"
)

type buildOptions struct {
	Race    bool     `yaml:"race"`
	LDFlags []string `yaml:"ldflags"`
	Tags    []string `yaml:"tags"`
}

func (opts buildOptions) generateArgs(useShell bool) []string {
	var args []string
	if opts.Race {
		args = append(args, "-race")
	}

	if len(opts.LDFlags) != 0 {
		args = append(args, "-ldflags",
			formatArgs(opts.LDFlags, useShell),
		)
	}

	if len(opts.Tags) != 0 {
		args = append(args, "-tags",
			// ref: https://golang.org/doc/go1.13#go-command
			// The go build flag -tags now takes a comma-separated list of build tags,
			// to allow for multiple tags in GOFLAGS. The space-separated form is
			// deprecated but still recognized and will be maintained.
			strings.Join(opts.Tags, ","),
		)
	}

	return args
}
