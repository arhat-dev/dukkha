package definition

import (
	"go/ast"
	"reflect"
	"strings"
)

// GetFieldTag gets the StructTag for a field
func GetFieldTag(field *ast.Field) reflect.StructTag {
	if field.Tag == nil {
		return ""
	}
	tag := strings.Replace(field.Tag.Value, "`", "", -1)
	return reflect.StructTag(tag)
}
