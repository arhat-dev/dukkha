# shellcheck shell=bash

set -x

command -v cygpath || echo "cygpath not found"
command -v dukkha || echo "dukkha not found"

make dukkha

./build/dukkha render <<EOF
tpl@tpl: |-
  {{- eval.Shell "command -v dukkha" -}}
shell@shell: command -v dukkha
EOF
