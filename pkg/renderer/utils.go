package renderer

import (
	"path/filepath"
	"time"

	"arhat.dev/dukkha/pkg/cache"
	"arhat.dev/dukkha/pkg/dukkha"
)

func FormatCacheDir(dukkhaCacheDir, rendererName string) string {
	return filepath.Join(dukkhaCacheDir, "renderer-"+rendererName)
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
		case "allow-expired":
			allowExpired = true
		case "cached-file":
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
