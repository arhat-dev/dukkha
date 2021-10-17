package git

import (
	"fmt"
	"strings"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/renderer/ssh"
	"arhat.dev/pkg/rshelper"
	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"gopkg.in/yaml.v3"
)

// nolint:revive
const (
	DefaultName = "git"
)

func init() { dukkha.RegisterRenderer(DefaultName, NewDefault) }

func NewDefault(name string) dukkha.Renderer {
	return &driver{
		name:        name,
		CacheConfig: renderer.CacheConfig{EnableCache: false},
		FetchConfig: FetchSpec{},
	}
}

var _ dukkha.Renderer = (*driver)(nil)

type driver struct {
	rs.BaseField `yaml:"-"`
	name         string

	renderer.CacheConfig `yaml:",inline"`

	SSHConfig ssh.Spec `yaml:",inline"`

	FetchConfig FetchSpec `yaml:",inline"`

	cache *renderer.Cache
}

func (d *driver) Init(ctx dukkha.ConfigResolvingContext) error {
	if d.EnableCache {
		d.cache = renderer.NewCache(int64(d.CacheSizeLimit), d.CacheMaxAge)
	}

	return nil
}

func (d *driver) RenderYaml(
	rc dukkha.RenderingContext, rawData interface{},
) (interface{}, error) {
	var (
		// reqURL format: <repo-name>.git/<path-in-repo>[@ref]
		reqURL      string
		sshConfig   *ssh.Spec
		fetchConfig *FetchSpec
	)

	switch t := rawData.(type) {
	case string:
		reqURL = t
		sshConfig = &d.SSHConfig
		fetchConfig = &d.FetchConfig
	case []byte:
		reqURL = string(t)
		sshConfig = &d.SSHConfig
		fetchConfig = &d.FetchConfig
	default:
		rawBytes, err := yamlhelper.ToYamlBytes(rawData)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: unexpected non yaml input: %w",
				d.name, err,
			)
		}

		spec := rshelper.InitAll(&inputFetchSpec{}, &rs.Options{
			InterfaceTypeHandler: rc,
		}).(*inputFetchSpec)

		err = yaml.Unmarshal(rawBytes, spec)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: failed to unmarshal input as config: %w",
				d.name, err,
			)
		}

		err = spec.ResolveFields(rc, -1)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: failed to resolve input config: %w",
				d.name, err,
			)
		}

		sshConfig = &spec.Spec
		fetchConfig = &spec.FetchSpec
	}

	if len(reqURL) != 0 {
		fetchConfig = &FetchSpec{}

		if idx := strings.LastIndexByte(reqURL, '@'); idx > 0 {
			fetchConfig.Ref = reqURL[idx+1:]
			reqURL = reqURL[:idx]
		}

		parts := strings.SplitAfterN(reqURL, ".git", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf(
				"invalid request url %q: no `.git` found",
				reqURL,
			)
		}

		fetchConfig.Repo, fetchConfig.Path = parts[0], parts[1]
		fetchConfig.Path = strings.TrimPrefix(fetchConfig.Path, "/")
	}

	var (
		data []byte
		err  error
	)
	if d.cache != nil {
		data, err = d.cache.Get(reqURL,
			renderer.CreateRefreshFuncForRemote(
				renderer.FormatCacheDir(rc.CacheDir(), d.name),
				d.CacheMaxAge,
				func(key string) ([]byte, error) {
					// key is the url we passed in
					return fetchConfig.fetchRemote(sshConfig)
				},
			),
		)
	} else {
		data, err = fetchConfig.fetchRemote(sshConfig)
	}

	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s failed to fetch http content: %w",
			d.name, err,
		)
	}

	return data, err
}
