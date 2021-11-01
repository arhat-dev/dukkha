package diff

import (
	"fmt"
	"strings"
)

// src is the yaml doc with rendering suffix unmarshaled as a trie node
// diffEntries are calculated by comparing yaml doc generated from src and
// actual state of that generated doc
func ReasonDiff(src *Node, diffEntries []*Entry) []*Entry {
	// now we have found some differences between original and current yaml doc
	// reason cause of these differences

	for _, d := range diffEntries {
		node, exact := src.Get(d.Key)
		if !exact {
			// TODO: check rendering suffix
			fmt.Printf("src is not compatible with key %v\n", d.Key)
			continue
		}

		switch d.Kind {
		case KindAdded:
			fmt.Printf("ADD: %s\n", strings.Join(d.Key, ""))
		case KindDeleted:
			fmt.Printf("DEL: %s\n", strings.Join(d.Key, ""))
		case KindUpdated:
			fmt.Printf("UPD: %s\n", strings.Join(d.Key, ""))
		}
		_ = node
	}

	return diffEntries
}
