package templateutils

func funcNameToFuncID(name string) funcID {
	switch name {

	// start of static funcs

	{{- range $_, $v := .StaticFuncs }}
	case funcName_{{- $v.Ident -}}:
		return funcID_{{- $v.Ident }}
	{{- end }}

	// end of static funcs

	// start of contextual funcs

	{{- range $_, $v := .ContextualFuncs }}
	case funcName_{{- $v.Ident -}}:
		return funcID_{{- $v.Ident }}
	{{- end }}

	// end of contextual funcs

	// start of placeholder funcs

	{{- range $_, $v := .PlaceholderFuncs }}
	case funcName_{{- $v.Ident -}}:
		return funcID_{{- $v.Ident }}
	{{- end }}

	// end of placeholder funcs

	default:
		return _unknown_template_func
	}
}

func (id funcID) String() string {
	switch id {

	// start of static funcs

	{{- range $_, $v := .StaticFuncs }}
	case funcID_{{- $v.Ident -}}:
		return funcName_{{- $v.Ident }}
	{{- end }}

	// end of static funcs

	// start of contextual funcs

	{{- range $_, $v := .ContextualFuncs }}
	case funcID_{{- $v.Ident -}}:
		return funcName_{{- $v.Ident }}
	{{- end }}

	// end of contextual funcs

	// start of placeholder funcs

	{{- range $_, $v := .PlaceholderFuncs }}
	case funcID_{{- $v.Ident -}}:
		return funcName_{{- $v.Ident }}
	{{- end }}

	// end of placeholder funcs
	default:
		return ""
	}
}

const (
	_unknown_template_func funcID = iota

	// start of static funcs

	{{- range $_, $v := .StaticFuncs }}
	funcID_{{- $v.Ident }}
	{{- end }}

	// end of static funcs

	// start of contextual funcs

	{{- range $_, $v := .ContextualFuncs }}
	funcID_{{- $v.Ident }}
	{{- end }}

	// end of contextual funcs

	// start of placeholder funcs

	{{- range $_, $v := .PlaceholderFuncs }}
	funcID_{{- $v.Ident }}
	{{- end }}

	// end of placeholder funcs

	funcID_COUNT
)

const (
	funcID_LAST_STATIC_FUNC = funcID_{{ .LastStaticFunc.Ident }}
	funcID_LAST_CONTEXTUAL_FUNC = funcID_{{ .LastContextualFunc.Ident }}
	funcID_LAST_Placeholder_FUNC = funcID_{{ .LastPlaceholderFunc.Ident }}
)

const (
	// start of static funcs

	{{- range $_, $v := .StaticFuncs }}
	funcName_{{- $v.Ident }} = "{{- $v.Name -}}"
	{{- end }}

	// end of static funcs

	// start of contextual funcs

	{{- range $_, $v := .ContextualFuncs }}
	funcName_{{- $v.Ident }} = "{{- $v.Name -}}"
	{{- end }}

	// end of contextual funcs

	// start of placeholder funcs

	{{- range $_, $v := .PlaceholderFuncs }}
	funcName_{{- $v.Ident }} = "{{- $v.Name -}}"
	{{- end }}

	// end of placeholder funcs
)
