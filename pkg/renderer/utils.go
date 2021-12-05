package renderer

import (
	"time"

	"arhat.dev/rs"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/cache"
	"arhat.dev/dukkha/pkg/dukkha"
)

func ResolveInputSpec(rc rs.RenderingHandler, rawBytes []byte, out rs.Field) (rs.Field, error) {
	out = rs.Init(out, nil)
	err := yaml.Unmarshal(rawBytes, out)
	if err != nil {
		return nil, err
	}

	err = out.ResolveFields(rc, -1)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func HandleRenderingRequestWithRemoteFetch(
	cache *cache.TwoTierCache,
	obj cache.IdentifiableObject,
	fetchRemote cache.RemoteCacheRefreshFunc,
	attributes []dukkha.RendererAttribute,
) ([]byte, error) {
	var (
		allowExpired         = false
		returnCachedFilePath = false
	)

	for _, attr := range attributes {
		switch attr {
		case AttrAllowExpired:
			allowExpired = true
		case AttrCachedFile:
			returnCachedFilePath = true
		}
	}

	var (
		data []byte
		err  error
	)

	if returnCachedFilePath {
		var path string
		path, _, err = cache.GetPath(
			obj,
			time.Now().Unix(),
			allowExpired,
			fetchRemote,
		)
		data = []byte(path)
	} else {
		data, _, err = cache.Get(
			obj,
			time.Now().Unix(),
			allowExpired,
			fetchRemote,
		)
	}

	return data, err
}
