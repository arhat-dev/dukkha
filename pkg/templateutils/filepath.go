package templateutils

import (
	"path/filepath"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/third_party/gomplate/conv"
)

type filepathNS struct {
	rc dukkha.RenderingContext
}

func createFilePathNS(rc dukkha.RenderingContext) *filepathNS {
	return &filepathNS{
		rc: rc,
	}
}

func (f *filepathNS) Base(in interface{}) string {
	return filepath.Base(conv.ToString(in))
}

func (f *filepathNS) Clean(in interface{}) string {
	return filepath.Clean(conv.ToString(in))
}

func (f *filepathNS) Dir(in interface{}) string {
	return filepath.Dir(conv.ToString(in))
}

func (f *filepathNS) Ext(in interface{}) string {
	return filepath.Ext(conv.ToString(in))
}

func (f *filepathNS) FromSlash(in interface{}) string {
	return filepath.FromSlash(conv.ToString(in))
}

func (f *filepathNS) IsAbs(in interface{}) bool {
	return filepath.IsAbs(conv.ToString(in))
}

func (f *filepathNS) Join(elem ...interface{}) string {
	s := conv.ToStrings(elem...)
	return filepath.Join(s...)
}

func (f *filepathNS) Match(pattern, name interface{}) (matched bool, err error) {
	return filepath.Match(conv.ToString(pattern), conv.ToString(name))
}

func (f *filepathNS) Rel(basepath, targpath interface{}) (string, error) {
	return filepath.Rel(conv.ToString(basepath), conv.ToString(targpath))
}

func (f *filepathNS) Split(in interface{}) []string {
	dir, file := filepath.Split(conv.ToString(in))
	return []string{dir, file}
}

func (f *filepathNS) ToSlash(in interface{}) string {
	return filepath.ToSlash(conv.ToString(in))
}

func (f *filepathNS) VolumeName(in interface{}) string {
	return filepath.VolumeName(conv.ToString(in))
}

func (f *filepathNS) Glob(in interface{}) ([]string, error) {
	return f.rc.FS().Glob(conv.ToString(in))
}

func (f *filepathNS) Abs(in interface{}) (string, error) {
	return f.rc.FS().Abs(conv.ToString(in))
}
