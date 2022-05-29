package templateutils

import (
	"reflect"

	"arhat.dev/dukkha/pkg/dukkha"
)

// nolint:gocyclo
func FuncNameToFuncID(name string) FuncID {
	switch name {

	// start of static funcs

	{{- range $_, $v := .StaticFuncs }}
	case FuncName_{{- $v.Ident -}}:
		return FuncID_{{- $v.Ident }}
	{{- end }}

	// end of static funcs

	// start of contextual funcs

	{{- range $_, $v := .ContextualFuncs }}
	case FuncName_{{- $v.Ident -}}:
		return FuncID_{{- $v.Ident }}
	{{- end }}

	// end of contextual funcs

	// start of placeholder funcs

	{{- range $_, $v := .PlaceholderFuncs }}
	case FuncName_{{- $v.Ident -}}:
		return FuncID_{{- $v.Ident }}
	{{- end }}

	// end of placeholder funcs

	default:
		return _unknown_template_func
	}
}

// nolint:gocyclo
func (id FuncID) String() string {
	switch id {

	// start of static funcs

	{{- range $_, $v := .StaticFuncs }}
	case FuncID_{{- $v.Ident -}}:
		return FuncName_{{- $v.Ident }}
	{{- end }}

	// end of static funcs

	// start of contextual funcs

	{{- range $_, $v := .ContextualFuncs }}
	case FuncID_{{- $v.Ident -}}:
		return FuncName_{{- $v.Ident }}
	{{- end }}

	// end of contextual funcs

	// start of placeholder funcs

	{{- range $_, $v := .PlaceholderFuncs }}
	case FuncID_{{- $v.Ident -}}:
		return FuncName_{{- $v.Ident }}
	{{- end }}

	// end of placeholder funcs
	default:
		return ""
	}
}

const (
	_unknown_template_func FuncID = iota

	// start of static funcs

	{{- range $_, $v := .StaticFuncs }}
	FuncID_{{- $v.Ident }} // {{ $v.FuncType }}
	{{- end }}

	// end of static funcs

	// start of contextual funcs

	{{- range $_, $v := .ContextualFuncs }}
	FuncID_{{- $v.Ident }} // {{ $v.FuncType }}
	{{- end }}

	// end of contextual funcs

	// start of placeholder funcs

	{{- range $_, $v := .PlaceholderFuncs }}
	FuncID_{{- $v.Ident }} // {{ $v.FuncType }}
	{{- end }}

	// end of placeholder funcs

	FuncID_COUNT
)

const (
	FuncID_LAST_Static_FUNC = FuncID_{{ .LastStaticFunc.Ident }}
	FuncID_LAST_Contextual_FUNC = FuncID_{{ .LastContextualFunc.Ident }}
	FuncID_LAST_Placeholder_FUNC = FuncID_{{ .LastPlaceholderFunc.Ident }}
)

const (
	// start of static funcs

	{{- range $_, $v := .StaticFuncs }}
	FuncName_{{- $v.Ident }} = "{{- $v.UserCallHandle -}}"
	{{- end }}

	// end of static funcs

	// start of contextual funcs

	{{- range $_, $v := .ContextualFuncs }}
	FuncName_{{- $v.Ident }} = "{{- $v.UserCallHandle -}}"
	{{- end }}

	// end of contextual funcs

	// start of placeholder funcs

	{{- range $_, $v := .PlaceholderFuncs }}
	FuncName_{{- $v.Ident }} = "{{- $v.UserCallHandle -}}"
	{{- end }}

	// end of placeholder funcs
)

var staticFuncs = [FuncID_LAST_Static_FUNC]any{
	{{- range $_, $v := .StaticFuncs }}
	FuncID_{{- $v.Ident }} - 1: {{- $v.CodeCallHandle -}},
	{{- end }}
}

func createContextualFuncs(rc dukkha.RenderingContext) *ContextualFuncs {
	var (
		ns_dukkha = createDukkhaNS(rc)
		ns_fs     = createFSNS(rc)
		ns_os     = createOSNS(rc)
		ns_eval   = createEvalNS(rc)
		ns_tag    = createTagNS(rc)
		ns_state  = createStateNS(rc)
		ns_misc   = createMiscNS(rc)
	)

	get_ns_dukkha := func() dukkhaNS { return ns_dukkha }
	get_ns_fs := func() fsNS { return ns_fs }
	get_ns_os := func() osNS { return ns_os }
	get_ns_eval := func() evalNS { return ns_eval }
	get_ns_tag := func() tagNS { return ns_tag }
	get_ns_state := func() stateNS { return ns_state }

	return &ContextualFuncs{
		{{- range $_, $v := .ContextualFuncs }}
		FuncID_{{- $v.Ident }} - FuncID_LAST_Static_FUNC - 1: reflect.ValueOf({{- $v.CodeCallHandle -}}),
		{{- end }}
	}
}
