#!/bin/sh

set -eux

local_package_names="dukkha constant matrix renderer sliceutils templateutils tools"

bin_sed="sed"
if command -v gsed >/dev/null 2>&1 ; then
  bin_sed="gsed"
fi

for pkg in ${local_package_names}; do
  dest_dir="pkg/${pkg}/symbols"
  mkdir -p "${dest_dir}"
  cd "${dest_dir}"
  yaegi extract -name "${pkg}_symbols" "arhat.dev/dukkha/pkg/${pkg}"

  case "${pkg}" in
  constant)
    ${bin_sed} -i 's#"go/constant"#goconst "go/constant"#g' ./*
    ${bin_sed} -i 's#constant\.MakeFromLiteral#goconst\.MakeFromLiteral#g' ./*
    ;;
  esac

  cd -
done

third_party_package_names="arhat.dev/rs gopkg.in/yaml.v3"

for pkg in ${third_party_package_names}; do
  name="$(basename "${pkg}")"

  dest_dir="third_party/${name}/symbols"
  mkdir -p "${dest_dir}"

  cd "${dest_dir}"
  yaegi extract -name "${name%.*}_symbols" "${pkg}"
  cd -
done
