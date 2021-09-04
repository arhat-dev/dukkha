// Code generated by 'yaegi extract arhat.dev/dukkha/pkg/templateutils'. DO NOT EDIT.

package templateutils_symbols

import (
	"arhat.dev/dukkha/pkg/templateutils"
	"reflect"
)

func init() {
	Symbols["arhat.dev/dukkha/pkg/templateutils/templateutils"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"AddPrefix":                       reflect.ValueOf(templateutils.AddPrefix),
		"AddSuffix":                       reflect.ValueOf(templateutils.AddSuffix),
		"CreateEmbeddedShellRunner":       reflect.ValueOf(templateutils.CreateEmbeddedShellRunner),
		"CreateTemplate":                  reflect.ValueOf(templateutils.CreateTemplate),
		"ExecCmdAsTemplateFuncCall":       reflect.ValueOf(templateutils.ExecCmdAsTemplateFuncCall),
		"GetDefaultImageTag":              reflect.ValueOf(templateutils.GetDefaultImageTag),
		"GetDefaultManifestTag":           reflect.ValueOf(templateutils.GetDefaultManifestTag),
		"GetDefaultTag":                   reflect.ValueOf(templateutils.GetDefaultTag),
		"RegisterTemplateFuncs":           reflect.ValueOf(templateutils.RegisterTemplateFuncs),
		"RemovePrefix":                    reflect.ValueOf(templateutils.RemovePrefix),
		"RemoveSuffix":                    reflect.ValueOf(templateutils.RemoveSuffix),
		"RunScriptInEmbeddedShell":        reflect.ValueOf(templateutils.RunScriptInEmbeddedShell),
		"SetDefaultImageTagIfNoTagSet":    reflect.ValueOf(templateutils.SetDefaultImageTagIfNoTagSet),
		"SetDefaultManifestTagIfNoTagSet": reflect.ValueOf(templateutils.SetDefaultManifestTagIfNoTagSet),

		// type definitions
		"TemplateFuncFactory": reflect.ValueOf((*templateutils.TemplateFuncFactory)(nil)),
	}
}
