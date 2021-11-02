package diff

import (
	"fmt"
	"strings"

	"arhat.dev/dukkha/pkg/diff"
)

// ReasonDiff (WIP) try to reason how to update src doc to generate these diff entries
//
// src is the yaml doc with rendering suffix unmarshaled as a trie node
// diffEntries are calculated by comparing yaml doc generated from src and
// actual state of that generated doc
func ReasonDiff(src *diff.Node, diffEntries []*diff.Entry) []*diff.Entry {
	for _, d := range diffEntries {
		node, exact := src.Get(d.Key)
		if !exact {
			fmt.Printf("src is not compatible with key %v\n", d.Key)
			continue
		}

		switch d.Kind {
		case diff.KindAdded:
			fmt.Printf("ADD: %s\n", strings.Join(d.Key, ""))
		case diff.KindDeleted:
			fmt.Printf("DEL: %s\n", strings.Join(d.Key, ""))
		case diff.KindUpdated:
			fmt.Printf("UPD: %s\n", strings.Join(d.Key, ""))
		}
		_ = node
	}

	return diffEntries
}
