package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/diff"
	"arhat.dev/dukkha/pkg/dukkha"
)

// TODO: support source doc (doc with rendering suffix which generates base doc)
func NewDiffCmd(ctx *dukkha.Context) *cobra.Command {
	_ = ctx

	var (
		source    string
		recursive bool
	)

	diffCmd := &cobra.Command{
		Use:           "diff [file-source] <file-original> <file-updated>",
		Short:         "Show yaml aware differences",
		Args:          cobra.ExactArgs(2),
		SilenceErrors: true,
		SilenceUsage:  true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   false,
			DisableNoDescFlag:   false,
			DisableDescriptions: true,
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return diffFile(*ctx, source, args[0], args[1])
		},
	}

	flags := diffCmd.Flags()
	flags.BoolVarP(&recursive, "recursive", "r", false,
		"diff directories recursively",
	)

	flags.StringVarP(&source, "source", "s", "",
		"set path to source doc which generated the original file",
	)

	return diffCmd
}

func diffFile(rc rs.RenderingHandler, srcDocSrc, baseDocSrc, newDocSrc string) error {
	if baseDocSrc == newDocSrc {
		return nil
	}

	if srcDocSrc == newDocSrc {
		return fmt.Errorf("invalid source doc should not be the same as new doc")
	}

	stdinUsed := false
	newDecoder := func(src string) (*yaml.Decoder, func(), error) {
		if src == "-" {
			if stdinUsed {
				return nil, nil, fmt.Errorf("only one of src/base/new doc can use stdin")
			}
			stdinUsed = true
			return yaml.NewDecoder(os.Stdin), func() {}, nil
		}

		rd, err := os.Open(src)
		if err != nil {
			return nil, nil, err
		}

		return yaml.NewDecoder(rd), func() { _ = rd.Close() }, nil
	}

	var sd *yaml.Decoder
	if srcDocSrc != "" && srcDocSrc != baseDocSrc {
		var (
			cleanupSD func()
			err       error
		)
		sd, cleanupSD, err = newDecoder(srcDocSrc)
		if err != nil {
			return fmt.Errorf("failed to open src doc: %w", err)
		}
		defer cleanupSD()
	}

	bd, cleanupBD, err := newDecoder(baseDocSrc)
	if err != nil {
		return fmt.Errorf("failed to open base doc: %w", err)
	}
	defer cleanupBD()

	nd, cleanupND, err := newDecoder(newDocSrc)
	if err != nil {
		return fmt.Errorf("failed to open target doc: %w", err)
	}
	defer cleanupND()

	var (
		sdDrain  = sd == nil
		bdDrain  = false
		ndDrain  = false
		printSep = false
	)

	for {
		current := new(diff.Node)
		err = nd.Decode(current)
		if err != nil {
			if err != io.EOF {
				return err
			}

			if bdDrain && sdDrain {
				// all drained
				return nil
			}

			ndDrain = true
		}

		base := new(diff.Node)
		err = bd.Decode(base)
		if err != nil {
			if err != io.EOF {
				return err
			}

			if ndDrain && sdDrain {
				// all drained
				return nil
			}

			bdDrain = true
		}

		src := base
		if sd != nil {
			src = new(diff.Node)
			err = sd.Decode(src)
			if err != nil {
				if err != io.EOF {
					return err
				}

				if ndDrain && bdDrain {
					// all drained
					return nil
				}

				sdDrain = true
			}
		}

		if printSep {
			fmt.Println("---")
		} else {
			printSep = true
		}

		diffEntries := diff.Diff(base, current)
		if len(diffEntries) == 0 {
			fmt.Println("# no difference")
			continue
		}

		reasons := reasonDiff(rc, src, current, diffEntries)
		for _, v := range reasons {
			fmt.Print(v)
		}
	}
}
