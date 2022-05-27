package templateutils

import (
	"path"

	"github.com/bmatcuk/doublestar/v4"
)

// pathNS for slash-separated paths
type pathNS struct{}

func (pathNS) Base(p String) string  { return path.Base(must(toString(p))) }
func (pathNS) Clean(p String) string { return path.Clean(must(toString(p))) }
func (pathNS) Dir(p String) string   { return path.Dir(must(toString(p))) }
func (pathNS) Ext(p String) string   { return path.Ext(must(toString(p))) }
func (pathNS) IsAbs(p String) bool   { return path.IsAbs(must(toString(p))) }

func (pathNS) Join(elem ...String) (_ string, err error) {
	parts, err := toStrings(elem)
	if err != nil {
		return
	}

	return path.Join(parts...), nil
}

func (pathNS) Match(pattern, name String) (matched bool, err error) {
	ptn, err := toString(pattern)
	if err != nil {
		return
	}

	tgt, err := toString(name)
	if err != nil {
		return
	}

	return doublestar.Match(ptn, tgt)
}

func (pathNS) Split(p String) []string {
	dir, file := path.Split(must(toString(p)))
	return []string{dir, file}
}
