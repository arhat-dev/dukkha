package templateutils

import (
	"path/filepath"

	"arhat.dev/dukkha/third_party/gomplate/conv"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/afero"
)

type _filepathNS struct{}

var filepathNS = &_filepathNS{}

func (f *_filepathNS) Base(in interface{}) string {
	return filepath.Base(conv.ToString(in))
}

func (f *_filepathNS) Clean(in interface{}) string {
	return filepath.Clean(conv.ToString(in))
}

func (f *_filepathNS) Dir(in interface{}) string {
	return filepath.Dir(conv.ToString(in))
}

func (f *_filepathNS) Ext(in interface{}) string {
	return filepath.Ext(conv.ToString(in))
}

func (f *_filepathNS) FromSlash(in interface{}) string {
	return filepath.FromSlash(conv.ToString(in))
}

func (f *_filepathNS) IsAbs(in interface{}) bool {
	return filepath.IsAbs(conv.ToString(in))
}

func (f *_filepathNS) Join(elem ...interface{}) string {
	s := conv.ToStrings(elem...)
	return filepath.Join(s...)
}

func (f *_filepathNS) Match(pattern, name interface{}) (matched bool, err error) {
	return filepath.Match(conv.ToString(pattern), conv.ToString(name))
}

func (f *_filepathNS) Rel(basepath, targpath interface{}) (string, error) {
	return filepath.Rel(conv.ToString(basepath), conv.ToString(targpath))
}

func (f *_filepathNS) Split(in interface{}) []string {
	dir, file := filepath.Split(conv.ToString(in))
	return []string{dir, file}
}

func (f *_filepathNS) ToSlash(in interface{}) string {
	return filepath.ToSlash(conv.ToString(in))
}

func (f *_filepathNS) VolumeName(in interface{}) string {
	return filepath.VolumeName(conv.ToString(in))
}

func (f *_filepathNS) Glob(in interface{}) ([]string, error) {
	return doublestar.Glob(afero.NewIOFS(afero.NewOsFs()), conv.ToString(in))
}
