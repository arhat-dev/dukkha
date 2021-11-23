{{- define "artifacts.bin-name" -}}

dukkha.{{- matrix.kernel -}}.{{- matrix.arch -}}

  {{- if eq matrix.kernel "windows" -}}
    .exe
  {{- end -}}

{{- end -}}
