package diff

import (
	"bytes"
	"fmt"
	"strings"
	_ "unsafe" // for go:linkname

	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/diff"
)

type Reason struct {
	// Renderer name of the decisive step to generate difference
	Renderer string
	// Input to the Renderer that, once its content updated
	// can cause expected differernce
	Input string

	Err error

	DiffEntries []*diff.Entry
}

func (r *Reason) String() string {
	if r.Err != nil {
		if len(r.Renderer) != 0 {
			return r.Renderer + ": " + r.Input + ": " + r.Err.Error() + "\n"
		}

		return r.Err.Error() + "\n"
	}

	ret := &strings.Builder{}
	for _, ent := range r.DiffEntries {
		if len(r.Renderer) == 0 {
			_, _ = ret.WriteString(string(ent.Kind) + " ")
			_, _ = ret.WriteString(strings.Join(ent.Key, "") + " ")
			_, _ = ret.Write(handleRawInput(ent.DivertAt.RawNode))
			continue
		}

		_, _ = ret.WriteString(
			fmt.Sprintln(
				r.Renderer+":",
				r.Input,
				ent.Kind,
				strings.Join(ent.Key, ""),
				"=>",
				string(handleRawInput(ent.DivertAt.RawNode)),
			),
		)
	}

	return ret.String()
}

// reasonDiff (WIP) try to reason how to update src doc to generate these diff entries
//
// src is the yaml doc with rendering suffix unmarshaled as a trie node
// diffEntries are calculated by comparing yaml doc generated from src and
// actual state of that generated doc
func reasonDiff(rc rs.RenderingHandler, src, current *diff.Node, diffEntries []*diff.Entry) []*Reason {
	var ret []*Reason
	for i, d := range diffEntries {
		node, tailKey := src.Get(d.Key)
		if len(tailKey) != 0 {
			// non exact match, there should be some rendering suffix to generate such
			// difference, or the entry was added manually

			if len(node.Renderers) == 0 {
				ret = append(ret, &Reason{
					Err: fmt.Errorf("src is not compatible with key: %s", strings.Join(d.Key, "")),
				})

				continue
			}

			node = node.Clone()

			rdrs := node.Renderers

			lastRdr := -1
		findLastRdr:
			for j := len(rdrs) - 1; j >= 0; j-- {
				// TODO: handle patch spec and more renderers
				// find last meaningful renderer
				rdr := rdrs[j]
				switch rdr.Name {
				case "", "echo":
					continue
				case "http", "file", "s3":
					// only renderers with predictable output are supported as final renderer
					//
					// other renderers like template/shell can generate very different value
					// with different environment variables & global value
					lastRdr = j
					break findLastRdr
				default:
					ret = append(ret, &Reason{
						Err: fmt.Errorf("unsupported renderer %q as final renderer", rdr.Name),
					})

					break findLastRdr
				}
			}

			if lastRdr == -1 {
				// no meaningful renderer exists, doc incompatible
				ret = append(ret, &Reason{
					Err: fmt.Errorf("src is not compatible with key: %s",
						strings.Join(d.Key, ""),
					),
				})

				continue
			}

			// render input for last meaningful renderer
			rawInput, err := tryRender(rc, node.RawNode, rdrs[:lastRdr])
			if err != nil {
				ret = append(ret, &Reason{
					Err: fmt.Errorf("failed to render input for last meaningful renderer %v", err),
				})

				continue
			}

			reason := &Reason{
				Renderer: rdrs[lastRdr].Name,
				Input:    string(handleRawInput(rawInput)),
			}

			ret = append(ret, reason)

			rawInput, err = tryRender(rc, rawInput, rdrs[lastRdr:])
			if err != nil {
				reason.Err = err
				continue
			}

			rn := new(diff.Node)
			err = rawInput.Decode(rn)
			if err != nil {
				reason.Err = err
				continue
			}

			target, tk := current.Get(d.Key[:len(d.Key)-len(tailKey)])
			if len(tk) != 0 {
				panic("invalid diff set")
			}

			reason.DiffEntries = diff.Diff(rn, target)

			continue
		}

		// TODO: exact match, but maybe we have rendering suffix in child node somewhere
		// 	    if that's the case, the document need manual inspection

		ret = append(ret, &Reason{
			DiffEntries: []*diff.Entry{diffEntries[i]},
		})
	}

	return ret
}

func handleRawInput(rawInput *yaml.Node) []byte {
	inputData, err := rs.NormalizeRawData(rawInput)
	if err != nil {
		panic(err)
	}

	inputBytes, err := yamlhelper.ToYamlBytes(inputData)
	if err != nil {
		panic(err)
	}

	return bytes.TrimSpace(inputBytes)
}

// tryRender maintain same behavior as arhat.dev/rs.handleUnresolvedField
//
// ref: https://github.com/arhat-dev/rs/blob/master/resolve.go#L206
func tryRender(rc rs.RenderingHandler, n *yaml.Node, renderers []*diff.RendererSpec) (*yaml.Node, error) {

	for _, v := range renderers {
		ht, err := rs.ParseTypeHint(v.TypeHint)
		if err != nil {
			return nil, err
		}

		n, err = renderOnce(rc, v.Name, ht, v.Patch, n)
		if err != nil {
			return nil, err
		}
	}

	return n, nil
}

//go:linkname renderOnce arhat.dev/rs.tryRender
func renderOnce(
	rc rs.RenderingHandler,
	rendererName string,
	typeHint rs.TypeHint,
	isPatchSpec bool,
	toResolve *yaml.Node,
) (*yaml.Node, error)
