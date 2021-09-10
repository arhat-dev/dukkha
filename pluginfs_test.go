package dukkha

import (
	"io/fs"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPluginFS_packages(t *testing.T) {
	const (
		PluginGOPATH = "goPath"
	)

	pfs := NewPluginFS(PluginGOPATH, "plugin")

	tests := []struct {
		name string
		path string
	}{
		{
			name: "package rs",
			path: path.Join(PluginGOPATH, "src", "plugin", "vendor", "arhat.dev/rs"),
		},
		{
			name: "package yaml.v3",
			path: path.Join(PluginGOPATH, "src", "plugin", "vendor", "gopkg.in/yaml.v3"),
		},
		{
			name: "package dukkha",
			path: path.Join(PluginGOPATH, "src", "plugin", "vendor", "arhat.dev/dukkha/pkg"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, err := pfs.Open(test.path)
			assert.NoError(t, err)

			info, err := f.Stat()
			assert.NoError(t, err)
			assert.True(t, info.IsDir())

			entries, err := fs.ReadDir(pfs, test.path)
			assert.NoError(t, err)
			for _, ent := range entries {
				t.Log(ent.Name())
			}
		})
	}
}
