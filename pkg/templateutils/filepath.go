package templateutils

import (
	"path/filepath"

	"arhat.dev/dukkha/pkg/dukkha"
)

func createFilePathNS(rc dukkha.RenderingContext) filepathNS { return filepathNS{rc: rc} }

type filepathNS struct{ rc dukkha.RenderingContext }

func (filepathNS) Base(p String) string               { return filepath.Base(toString(p)) }
func (filepathNS) Clean(p String) string              { return filepath.Clean(toString(p)) }
func (filepathNS) Dir(p String) string                { return filepath.Dir(toString(p)) }
func (filepathNS) Ext(p String) string                { return filepath.Ext(toString(p)) }
func (filepathNS) FromSlash(p String) string          { return filepath.FromSlash(toString(p)) }
func (filepathNS) IsAbs(p String) bool                { return filepath.IsAbs(toString(p)) }
func (filepathNS) Join(elem ...String) string         { return filepath.Join(toStrings(elem...)...) }
func (filepathNS) Split(p String) (dir, file string)  { return filepath.Split(toString(p)) }
func (filepathNS) ToSlash(p String) string            { return filepath.ToSlash(toString(p)) }
func (filepathNS) VolumeName(p String) string         { return filepath.VolumeName(toString(p)) }
func (ns filepathNS) Glob(p String) ([]string, error) { return ns.rc.FS().Glob(toString(p)) }
func (ns filepathNS) Abs(p String) (string, error)    { return ns.rc.FS().Abs(toString(p)) }

func (filepathNS) Match(pattern, name String) (matched bool, err error) {
	return filepath.Match(toString(pattern), toString(name))
}

func (filepathNS) Rel(basepath, targpath String) (string, error) {
	return filepath.Rel(toString(basepath), toString(targpath))
}
