// Package af (archivefile) provides a renderer generating value by extracting
// file content from archive directly
package af

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"github.com/h2non/filetype"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/cache"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer"
)

const (
	DefaultName = "af"
)

func init() {
	dukkha.RegisterRenderer(DefaultName, NewDefault)
}

func NewDefault(name string) dukkha.Renderer {
	return &Driver{
		name:        name,
		CacheConfig: renderer.CacheConfig{Enabled: false},
	}
}

var _ dukkha.Renderer = (*Driver)(nil)

type Driver struct {
	rs.BaseField `yaml:"-"`

	RendererAlias string `yaml:"alias"`

	name string

	CacheConfig renderer.CacheConfig `yaml:"cache"`

	cache *cache.TwoTierCache
}

func (d *Driver) Alias() string { return d.RendererAlias }

func (d *Driver) Init(cacheFS *fshelper.OSFS) error {
	if d.CacheConfig.Enabled {
		d.cache = cache.NewTwoTierCache(
			cacheFS,
			int64(d.CacheConfig.MaxItemSize),
			int64(d.CacheConfig.Size),
			int64(d.CacheConfig.Timeout.Seconds()),
		)
	} else {
		d.cache = cache.NewTwoTierCache(cacheFS, 0, 0, -1)
	}

	return nil
}

func (d *Driver) RenderYaml(
	rc dukkha.RenderingContext, rawData interface{}, attributes []dukkha.RendererAttribute,
) ([]byte, error) {
	rawData, err := rs.NormalizeRawData(rawData)
	if err != nil {
		return nil, err
	}

	var (
		onlineSpec string
		spec       *inputSpec
	)

	switch t := rawData.(type) {
	case string:
		onlineSpec = t
	case []byte:
		onlineSpec = string(t)
	default:
		var rawDataBytes []byte
		rawDataBytes, err = yamlhelper.ToYamlBytes(t)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", d.name, err)
		}

		spec = rs.Init(&inputSpec{}, nil).(*inputSpec)
		err = yaml.Unmarshal(rawDataBytes, spec)
		if err != nil {
			return nil, fmt.Errorf("%s: invalid input spec: %w", d.name, err)
		}

		err = spec.ResolveFields(rc, -1)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to resolve input spec: %w", d.name, err)
		}
	}

	if spec == nil {
		spec = parseOneLineSpec(onlineSpec)
	} else {
		spec.Path = path.Clean(spec.Path)
	}

	data, err := renderer.HandleRenderingRequestWithRemoteFetch(
		d.cache,
		spec,
		func(obj cache.IdentifiableObject) (io.ReadCloser, error) {
			return extractFileFromArchive(rc.FS(), obj)
		},
		attributes,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", d.name, err)
	}

	return data, err
}

func parseOneLineSpec(onelineSpec string) *inputSpec {
	onelineSpec = strings.TrimSpace(onelineSpec)

	ret := &inputSpec{
		// there is no way to set password in one line spec
		Password: "",
	}

	idx := strings.LastIndex(onelineSpec, ":")
	if idx == -1 {
		ret.Archive = onelineSpec
	} else {
		ret.Archive = onelineSpec[:idx]
		ret.Path = path.Clean(onelineSpec[idx+1:])
	}

	return ret
}

func extractFileFromArchive(ofs *fshelper.OSFS, obj cache.IdentifiableObject) (io.ReadCloser, error) {
	spec := obj.(*inputSpec)
	info, err := ofs.Stat(spec.Archive)
	if err != nil {
		return nil, err
	}

	typ, err := filetype.MatchFile(spec.Archive)
	if err != nil {
		return nil, err
	}

	f, err := ofs.Open(spec.Archive)
	if err != nil {
		return nil, err
	}

	type src struct {
		sizeIface
		*os.File
	}

	return unarchive(&src{info, f.(*os.File)}, typ, spec.Path, spec.Password)
}
