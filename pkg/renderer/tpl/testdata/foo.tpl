{{- define "foo" -}}
foo: {{ var.foo | default "null" }}
{{- end -}}
