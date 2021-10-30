# Template Functions

{{ range $_, $v := . -}}
- `{{ $v.Name }}` (`{{ $v.Func }}`)
{{ end -}}
