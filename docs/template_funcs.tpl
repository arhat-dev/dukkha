# Template Functions

{{ range $_, $v := . -}}
- `{{ $v.Name }}{{ $v.Func | trimPrefix "func" }}`
{{ end -}}
