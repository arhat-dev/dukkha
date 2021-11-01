package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/diff"
	"arhat.dev/dukkha/pkg/dukkha"
)

// TODO: support source doc (doc with rendering suffix which generates base doc)
func NewDiffCmd(ctx *dukkha.Context) *cobra.Command {
	_ = ctx

	var (
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
			return diffFile(args[0], args[1])
		},
	}

	flags := diffCmd.Flags()
	flags.BoolVarP(&recursive, "recursive", "r", false,
		"diff directories recursively",
	)

	return diffCmd
}

func diffFile(baseDocSrc, newDocSrc string) error {
	if baseDocSrc == newDocSrc {
		return nil
	}

	var (
		baseDoc, newDoc io.ReadCloser

		err error
	)

	switch "-" {
	case baseDocSrc:
		baseDoc = os.Stdin
		newDoc, err = os.Open(newDocSrc)
		if err != nil {
			return err
		}
	case newDocSrc:
		baseDoc, err = os.Open(baseDocSrc)
		newDoc = os.Stdin
		if err != nil {
			return err
		}
	default:
		baseDoc, err = os.Open(baseDocSrc)
		if err != nil {
			return err
		}
		newDoc, err = os.Open(newDocSrc)
		if err != nil {
			return err
		}
	}

	defer func() {
		_ = baseDoc.Close()
		_ = newDoc.Close()
	}()

	bd := yaml.NewDecoder(baseDoc)
	nd := yaml.NewDecoder(newDoc)

	var (
		baseDocDrained, newDocDrained bool
	)

	count := 0
	for {
		current := new(diff.Node)
		err = nd.Decode(current)
		if err != nil {
			if err != io.EOF {
				return err
			}

			if baseDocDrained {
				// both drained
				return nil
			}

			newDocDrained = true
		}

		base := new(diff.Node)
		err = bd.Decode(base)
		if err != nil {
			if err != io.EOF {
				return err
			}

			if newDocDrained {
				// both drained
				return nil
			}

			baseDocDrained = true
		}

		if count != 0 {
			fmt.Println("---")
		}

		diffEntries := diff.Diff(base, current, []string{})
		// TODO: use source input as reason source
		if len(diff.ReasonDiff(base, diffEntries)) == 0 {
			fmt.Println("# no difference")
		}

		count++
	}
}
