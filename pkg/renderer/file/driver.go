package file

import (
	"encoding/hex"
	"fmt"
	"io/fs"
	"strings"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/pkg/sha256helper"
	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/cache"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer"
)

const (
	DefaultName = "file"
)

func init() {
	dukkha.RegisterRenderer(DefaultName, NewDefault)
}

func NewDefault(name string) dukkha.Renderer {
	return &Driver{name: name}
}

type Driver struct {
	rs.BaseField `yaml:"-"`

	renderer.BaseInMemCachedRenderer `yaml:",inline"`

	name string

	// BasePath to be used instead of current working dir
	BasePath string `yaml:"base_path"`
}

func (d *Driver) RenderYaml(
	rc dukkha.RenderingContext,
	rawData interface{},
	attributes []dukkha.RendererAttribute,
) ([]byte, error) {
	rawData, err := rs.NormalizeRawData(rawData)
	if err != nil {
		return nil, err
	}

	dataBytes, err := yamlhelper.ToYamlBytes(rawData)
	if err != nil {
		return nil, err
	}

	var (
		cacheData  bool
		cachedFile bool
	)
	for _, attr := range d.Attributes(attributes) {
		switch attr {
		case renderer.AttrCacheData:
			cacheData = true
		case renderer.AttrCachedFile:
			cachedFile = true
		default:
		}
	}

	if cacheData {
		return d.cacheData(dataBytes)
	}

	return d.readFile(
		rc.FS(),
		strings.TrimSpace(string(dataBytes)),
		cachedFile,
	)
}

func (d *Driver) cacheData(data []byte) ([]byte, error) {
	filename := hex.EncodeToString(sha256helper.Sum(data))
	err := d.CacheFS.WriteFile(filename, data, 0400)
	if err != nil {
		return nil, fmt.Errorf("renderer.%s: %w", d.name, err)
	}

	path, err := d.CacheFS.Abs(filename)
	if err != nil {
		return nil, fmt.Errorf("renderer.%s: %w", d.name, err)
	}

	return []byte(path), nil
}

func (d *Driver) readFile(ofs *fshelper.OSFS, target string, getPath bool) ([]byte, error) {
	var (
		data []byte
		err  error
	)

	if len(d.BasePath) != 0 {
		var fs2 fs.FS
		fs2, err = ofs.Sub(d.BasePath)
		if err != nil {
			return nil, err
		}

		ofs = fs2.(*fshelper.OSFS)
	}

	if getPath {
		var ret string
		ret, err = ofs.Abs(target)
		if err != nil {
			return nil, err
		}

		return []byte(ret), nil
	}

	if d.Cache != nil {
		data, err = d.Cache.Get(
			cache.IdentifiableString(target),
			func(_ cache.IdentifiableObject) ([]byte, error) {
				return ofs.ReadFile(target)
			},
		)
	} else {
		data, err = ofs.ReadFile(target)
	}

	if err != nil {
		return nil, fmt.Errorf("renderer.%s: %w", d.name, err)
	}

	return data, err
}
