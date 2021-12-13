# shellcheck shell=bash

set -x

wine64 ./build/dukkha.windows.amd64.exe render <<EOF
tpl@tpl: |-
  {{- eval.Shell "command -v winepath" -}}
shell@shell: command -v winepath
EOF
