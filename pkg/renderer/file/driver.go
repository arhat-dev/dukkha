package file

import (
	"encoding/hex"
	"fmt"
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
		cacheData bool
	)
	for _, attr := range d.Attributes(attributes) {
		switch attr {
		case renderer.AttrCacheData:
			cacheData = true
		default:
		}
	}

	if cacheData {
		return d.cacheData(dataBytes)
	}

	return d.readFile(rc.FS(), strings.TrimSpace(string(dataBytes)))
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

func (d *Driver) readFile(fs *fshelper.OSFS, path string) ([]byte, error) {
	var (
		data []byte
		err  error
	)

	if d.Cache != nil {
		data, err = d.Cache.Get(
			cache.IdentifiableString(path),
			func(_ cache.IdentifiableObject) ([]byte, error) {
				return fs.ReadFile(path)
			},
		)
	} else {
		data, err = fs.ReadFile(path)
	}

	if err != nil {
		return nil, fmt.Errorf("renderer.%s: %w", d.name, err)
	}

	return data, err
}
