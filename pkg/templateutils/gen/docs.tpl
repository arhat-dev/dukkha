# Template Functions

{{- strings.NIndent 0 "" -}}

{{- range $_, $v := . }}
  {{- if not ($v.FuncType | strings.HasSuffix "NS") }}
- `{{ $v.UserCallHandle }}{{ $v.FuncType | strings.TrimPrefix "func" }}`
  {{- end }}
{{- end }}
