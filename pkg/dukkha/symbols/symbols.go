package dukkha_symbols

import (
	"reflect"

	"arhat.dev/dukkha/pkg/dukkha"
	"github.com/traefik/yaegi/interp"
)

var Symbols = interp.Exports{}

func init() {
	Symbols["arhat.dev/dukkha/pkg/dukkha/dukkha"] = map[string]reflect.Value{
		// Global variable definitions
		"GlobalInterfaceTypeHandler": reflect.ValueOf(&dukkha.GlobalInterfaceTypeHandler).Elem(),
	}
}
