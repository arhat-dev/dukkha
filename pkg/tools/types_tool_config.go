package tools

import (
	"reflect"

	"arhat.dev/dukkha/pkg/field"
)

// ToolConfigType for tools.Config interface type registration
var ToolConfigType = reflect.TypeOf((*ToolConfig)(nil)).Elem()

type ToolConfig interface {
	field.Interface

	// Kind of the tool, e.g. golang, docker
	Kind() string
}
