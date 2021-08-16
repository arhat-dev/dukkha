package rshelper

import (
	"reflect"

	"arhat.dev/rs"
)

func InitAll(f rs.Field, h rs.InterfaceTypeHandler) rs.Field {
	rs.InitRecursively(reflect.ValueOf(f), h)
	return f
}
