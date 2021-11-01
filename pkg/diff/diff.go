package diff

type Kind string

const (
	// missing in base doc
	KindAdded Kind = "added"
	// updated in the new doc
	KindUpdated Kind = "updated"
	// deleted in the new doc
	KindDeleted Kind = "deleted"
)

type Entry struct {
	Key  []string `yaml:"key"`
	Kind Kind     `yaml:"kind"`

	// depending on kind
	//  - added: set to the added node of new doc
	// 	- updated: set to node of base doc that changed
	//  - deleted: set to node of base doc that deleted
	DivertAt *Node `yaml:"divert_at"`
}

// Diff
func Diff(base, other *Node, visitingKey []string) []*Entry {
	switch {
	case base == nil && other == nil:
		return nil

		// not both nil
	case base == nil || other == nil:
		// only one is nil

		if base == nil {
			// base: null
			// other: some value
			// => kind = added
			return []*Entry{{Key: visitingKey, DivertAt: other, Kind: KindAdded}}
		}

		// base: some value
		// other: null
		// => kind = deleted
		return []*Entry{{Key: visitingKey, DivertAt: base, Kind: KindDeleted}}

		// both not nil
	case base.scalarData == nil && other.scalarData == nil:
		// all non scalar, compare children

		// not both are map/sclice
	case base.scalarData == nil || other.scalarData == nil:
		// only one is scalar

		if base.scalarData == nil {
			// other: scalar or null
			if len(base.childIndex) == 0 {
				// base: empty => added
				return []*Entry{{Key: visitingKey, DivertAt: other, Kind: KindAdded}}
			}

			// base: map/slice => updated

			if other.raw.ShortTag() == "!!null" {
				return []*Entry{{Key: visitingKey, DivertAt: base, Kind: KindDeleted}}
			}

			return []*Entry{{Key: visitingKey, DivertAt: base, Kind: KindUpdated}}
		}

		// base: scalar
		// other: map or null

		if len(other.childIndex) == 0 {
			// other: empty => deleted
			return []*Entry{{Key: visitingKey, DivertAt: base, Kind: KindDeleted}}
		}

		if base.raw.ShortTag() == "!!null" {
			return []*Entry{{Key: visitingKey, DivertAt: other, Kind: KindAdded}}
		}

		// other: map/slice => updated
		return []*Entry{{Key: visitingKey, DivertAt: base, Kind: KindUpdated}}

		// both scalar
	case base.scalarData.Value == other.scalarData.Value:
		// same scalar
		// TODO: check scalarData.Style?
		return nil
	default:
		// different scalar value
		return []*Entry{{Key: visitingKey, DivertAt: base, Kind: KindUpdated}}
	}

	var (
		ret     []*Entry
		visited = make(map[string]struct{})
	)

	// iterate by children (slice) rather than childIdx (map)
	// to generate deterministic result

	for _, child := range base.children {
		visited[child.elemKey] = struct{}{}

		if len(other.childIndex) == 0 {
			ret = append(ret, &Entry{
				Key:      append(visitingKey, child.elemKey),
				Kind:     KindDeleted,
				DivertAt: child,
			})

			continue
		}

		j, ok := other.childIndex[child.elemKey]
		if !ok {
			ret = append(ret, &Entry{
				Key:      append(visitingKey, child.elemKey),
				Kind:     KindDeleted,
				DivertAt: child,
			})

			continue
		}

		ret = append(ret,
			Diff(
				child,
				other.children[j],
				append(visitingKey, child.elemKey),
			)...,
		)
	}

	for _, child := range other.children {
		if _, skip := visited[child.elemKey]; skip {
			continue
		}

		// only can be missing
		ret = append(ret, &Entry{
			Key:      append(visitingKey, child.elemKey),
			Kind:     KindAdded,
			DivertAt: child,
		})
	}

	return ret
}
