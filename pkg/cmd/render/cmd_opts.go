package render

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/itchyny/gojq"
)

type Options struct {
	outputFormat string
	indentSize   int
	indentStyle  string
	recursive    bool
	resultQuery  string

	chdir       []string
	outputDests []string
}

func (opts *Options) Resolve(args []string, defaultOutputDest io.Writer) (*ResolvedOptions, error) {
	ret := &ResolvedOptions{
		_specs: make(map[string]*renderingSpec),

		outputFormat: opts.outputFormat,
		recursive:    opts.recursive,
	}

	// validate input source, and prepare rendering specs for each of them
	if len(args) == 0 {
		// stdin input without using `-`, generalize this case
		args = []string{"-"}
		ret._specs["-"] = &renderingSpec{}
	} else {
		foundStdin := false
		for _, v := range args {
			switch v {
			case "-", "":
				if foundStdin {
					return nil, fmt.Errorf("too many stdin source, only one allowed")
				}

				foundStdin = true
				ret._specs["-"] = &renderingSpec{}
			default:
				ret._specs[v] = &renderingSpec{}
			}
		}
	}

	// length of args is at least 1

	// validate chdir options
	if len(opts.chdir) > len(args) {
		return nil, fmt.Errorf(
			"too many chdir options: more than the count of input source (%d)",
			len(args),
		)
	}

	// count of output destinations
	// can either be 0 or equals count of input sources
	if destCount := len(opts.outputDests); destCount != 0 && destCount != len(args) {
		return nil, fmt.Errorf(
			"destination count and source count not match: source = %d, dest = %d",
			len(args), len(opts.outputDests),
		)
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

	for i, dst := range opts.outputDests {
		src := args[i]

		path, err := filepath.Abs(dst)
		if err != nil {
			return nil, err
		}

		ret._specs[src].outputPath = &path
	}

	for i, chdir := range opts.chdir {
		if chdir == "" {
			// no path provided, use default, if wish to stay at current dir
			// it should be "."

			continue
		}

		targetChdir, err := filepath.Abs(chdir)
		if err != nil {
			return nil, fmt.Errorf("invalid chdir option: %w", err)
		}

		src := args[i]
		ret._specs[src].chdir = targetChdir
	}

	// set deafault chdir options for all input source without one

	for _, src := range args {
		switch src {
		case "-", "":
			continue
		}

		// do not follow symlink
		info, err := os.Lstat(src)
		if err != nil {
			return nil, fmt.Errorf("invalid input source: %w", err)
		}

		var chdir string
		if info.IsDir() {
			// we are going to chdir into src dicrectory, or we
			// have already been there
			//
			// so the entrypoint is always the current dir
			ret._specs[src].entrypoint = "."
			chdir = src
		} else {
			// regular file
			ret._specs[src].entrypoint, err = filepath.Abs(src)
			if err != nil {
				return nil, err
			}

			chdir = filepath.Dir(src)
		}

		// only set default, do not override
		if len(ret._specs[src].chdir) == 0 {
			ret._specs[src].chdir, err = filepath.Abs(chdir)
			if err != nil {
				return nil, err
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

type renderingSpec struct {
	outputPath *string

	chdir      string
	entrypoint string
}

type ResolvedOptions struct {
	indentStr string
	query     *gojq.Query

	_specs        map[string]*renderingSpec
	createEncoder func(w io.Writer) (encoder, error)

	outputFormat string
	recursive    bool
}

func (opts *ResolvedOptions) getSpecFor(src string) *renderingSpec {
	spec, ok := opts._specs[src]
	if !ok {
		panic("unknown source")
	}

	return spec
}

func (opts *ResolvedOptions) OutputPathFor(src string) *string {
	return opts.getSpecFor(src).outputPath
}

func (opts *ResolvedOptions) ChdirFor(src string) string {
	return opts.getSpecFor(src).chdir
}

func (opts *ResolvedOptions) EntrypointFor(src string) string {
	return opts.getSpecFor(src).entrypoint
}
