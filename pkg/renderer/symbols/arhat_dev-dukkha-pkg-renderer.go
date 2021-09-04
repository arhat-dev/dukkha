// Code generated by 'yaegi extract arhat.dev/dukkha/pkg/renderer'. DO NOT EDIT.

package renderer_symbols

import (
	"arhat.dev/dukkha/pkg/renderer"
	"reflect"
)

func init() {
	Symbols["arhat.dev/dukkha/pkg/renderer/renderer"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"CreateRefreshFuncForRemote": reflect.ValueOf(renderer.CreateRefreshFuncForRemote),
		"FormatCacheDir":             reflect.ValueOf(renderer.FormatCacheDir),
		"NewCache":                   reflect.ValueOf(renderer.NewCache),

		// type definitions
		"Cache":            reflect.ValueOf((*renderer.Cache)(nil)),
		"CacheConfig":      reflect.ValueOf((*renderer.CacheConfig)(nil)),
		"CacheRefreshFunc": reflect.ValueOf((*renderer.CacheRefreshFunc)(nil)),
	}
}
