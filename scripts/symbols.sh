#!/bin/sh

set -eux

# local_package_names="dukkha"
#
# bin_sed="sed"
# if command -v gsed >/dev/null 2>&1 ; then
#   bin_sed="gsed"
# fi
#
# for pkg in ${local_package_names}; do
#   dest_dir="pkg/${pkg}/symbols"
#   mkdir -p "${dest_dir}"
#   cd "${dest_dir}"
#   yaegi extract -name "${pkg}_symbols" "arhat.dev/dukkha/pkg/${pkg}"
#
#   case "${pkg}" in
#   constant)
#     ${bin_sed} -i 's#"go/constant"#goconst "go/constant"#g' ./*
#     ${bin_sed} -i 's#constant\.MakeFromLiteral#goconst\.MakeFromLiteral#g' ./*
#     ;;
#   esac
#
#   cd -
# done

vendor_pkgs="arhat.dev/rs"
vendor_pkgs="${vendor_pkgs} gopkg.in/yaml.v3"
vendor_pkgs="${vendor_pkgs} github.com/evanphx/json-patch/v5"
vendor_pkgs="${vendor_pkgs} github.com/pkg/errors"
vendor_pkgs="${vendor_pkgs} mvdan.cc/sh/v3"
vendor_pkgs="${vendor_pkgs} arhat.dev/pkg"

for pkg in ${vendor_pkgs}; do
  dest_dir="third_party/${pkg}"

  # cleanup old data
  rm -rf "${dest_dir}"
  mkdir -p "$(dirname "${dest_dir}")"

  cp -a "vendor/${pkg}" "${dest_dir}"

  case "${pkg}" in
  "github.com/pkg/errors" | "mvdan.cc/sh/v3" | "arhat.dev/pkg")
    :
    ;;
  "golang.org/x/sys" | "golang.org/x/sync")
    :
    ;;
  *)
    mv "${dest_dir}/go.mod" "${dest_dir}/go.mod_"
    ;;
  esac
done

third_party_package_names="arhat.dev/rs gopkg.in/yaml.v3"

for pkg in ${third_party_package_names}; do
  name="$(basename "${pkg}")"

  dest_dir="third_party/${pkg}/symbols"
  mkdir -p "${dest_dir}"

  cd "${dest_dir}"

  pkg_name="${name%.*}_symbols"
  yaegi extract -name "${pkg_name}" "${pkg}"
  cat > "symbols.go" <<EOF
package ${pkg_name}

import "github.com/traefik/yaegi/interp"

var Symbols = interp.Exports{}
EOF
  cd -
done
