package tools

import (
	"context"
	"reflect"
)

type (
	ToolFactoryFunc       func(ctx context.Context, config interface{}) (Interface, error)
	ToolConfigFactoryFunc func() interface{}

	// nolint:structcheck,unused
	toolFactory struct {
		createTool ToolFactoryFunc
		newConfig  ToolConfigFactoryFunc
	}
)

// nolint:deadcode,unused,varcheck
var (
	supportedTools = make(map[ToolKey]toolFactory)
)

func Register(
	name string,
	createTool ToolFactoryFunc,
	createToolConfig ToolConfigFactoryFunc,
) {
	// supportedTools[toolKey{name: name}] = toolFactory{
	// 	f:  createTool,
	// 	cf: createToolConfig,
	// }
}

func GetToolConfigType(toolName string) (reflect.Type, error) {
	return reflect.TypeOf(nil), nil
}

func GetTaskConfigType(toolName, taskType string) (reflect.Type, error) {
	return reflect.TypeOf(nil), nil
}
