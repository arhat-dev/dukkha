package plugin

import (
	"reflect"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unrestricted"
	"github.com/traefik/yaegi/stdlib/unsafe"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

var (
	dukkhaSymbols = interp.Exports{
		"arhat.dev/dukkha/pkg/dukkha": map[string]reflect.Value{
			"GlobalInterfaceTypeHandler": reflect.ValueOf(dukkha.GlobalInterfaceTypeHandler),

			"ToolName":        reflect.ValueOf((*dukkha.ToolName)(nil)).Elem(),
			"TaskName":        reflect.ValueOf((*dukkha.TaskName)(nil)).Elem(),
			"ArbitraryValues": reflect.ValueOf((*dukkha.ArbitraryValues)(nil)).Elem(),
		},
		"arhat.dev/dukkha/pkg/tools": map[string]reflect.Value{
			"BaseTask":        reflect.ValueOf((*tools.BaseTask)(nil)).Elem(),
			"BaseTool":        reflect.ValueOf((*tools.BaseTool)(nil)).Elem(),
			"TaskExecRequest": reflect.ValueOf((*tools.TaskExecRequest)(nil)).Elem(),
			"Action":          reflect.ValueOf((*tools.Action)(nil)).Elem(),
			"TaskHooks":       reflect.ValueOf((*tools.TaskHooks)(nil)).Elem(),
		},
		"arhat.dev/dukkha/pkg/sliceutils": map[string]reflect.Value{
			"NewStrings":      reflect.ValueOf(sliceutils.NewStrings),
			"FormatStringMap": reflect.ValueOf(sliceutils.FormatStringMap),
		},
		"arhat.dev/dukkha/pkg/constant": map[string]reflect.Value{
			"GetOciOS":   reflect.ValueOf(constant.GetOciOS),
			"ARCH_AMD64": reflect.ValueOf(constant.ARCH_AMD64),
		},
	}
)

func init() {
	dukkhaSymbols["arhat.dev/dukkha/pkg/plugin"] = map[string]reflect.Value{
		"Symbols": reflect.ValueOf(dukkhaSymbols),
	}

	t := interp.New(interp.Options{})
	t.Use(stdlib.Symbols)
	t.Use(syscall.Symbols)
	t.Use(unsafe.Symbols)
	t.Use(unrestricted.Symbols)
	t.Use(dukkhaSymbols)

	t.ImportUsed()

	t.EvalPath("")

	_ = t
}

type ToolPlugin struct {
}

func NewToolPlugin() {

}
