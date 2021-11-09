package render

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/itchyny/gojq"
)

type Options struct {
	outputFormat string
	indentSize   int
	indentStyle  string
	recursive    bool
	resultQuery  string

	outputDests []string
}

type ResolvedOptions struct {
	indentStr string
	query     *gojq.Query

	outputMapping map[string]*string
	createEncoder func(w io.Writer) (encoder, error)

	outputFormat string
	recursive    bool
}

func (opts *Options) Resolve(args []string, defaultOutputDest io.Writer) (*ResolvedOptions, error) {
	ret := &ResolvedOptions{
		outputMapping: make(map[string]*string),

		outputFormat: opts.outputFormat,
		recursive:    opts.recursive,
	}

	switch opts.indentStyle {
	case "space":
		ret.indentStr = " "
	case "tab":
		ret.indentStr = "\t"
	default:
		ret.indentStr = opts.indentStyle
	}

	if len(opts.resultQuery) != 0 {
		var err error
		ret.query, err = gojq.Parse(opts.resultQuery)
		if err != nil {
			return nil, fmt.Errorf("invalid result query: %w", err)
		}
	}

	if len(opts.outputDests) != 0 {
		switch {
		case len(args) == 0 && len(opts.outputDests) != 1:
			return nil, fmt.Errorf("only one destination can be set for stdin input")
		case len(opts.outputDests) != len(args):
			return nil, fmt.Errorf(
				"number of output destination not matching sources: want %d, got %d",
				len(args), len(opts.outputDests),
			)
		}

		for i := range opts.outputDests {
			src := args[i]

			path, err := filepath.Abs(opts.outputDests[i])
			if err != nil {
				return nil, err
			}

			ret.outputMapping[src] = &path
		}
	} else {
		if len(args) == 0 {
			ret.outputMapping["-"] = nil
		} else {
			for _, src := range args {
				ret.outputMapping[src] = nil
			}
		}
	}

	var stdoutEnc encoder

	ret.createEncoder = func(w io.Writer) (encoder, error) {
		if w == nil {
			var err error
			if stdoutEnc == nil {
				stdoutEnc, err = newEncoder(
					ret.query, defaultOutputDest,
					opts.outputFormat, ret.indentStr, opts.indentSize,
				)
			}

			return stdoutEnc, err
		}

		return newEncoder(ret.query, w, opts.outputFormat, ret.indentStr, opts.indentSize)
	}

	return ret, nil
}
