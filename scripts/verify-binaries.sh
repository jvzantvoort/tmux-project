#!/bin/bash
#
#  verify-binaries.sh <goos> <goarch> <directory>
#
#  Fails when a Go binary in <directory> is not built for <goos>/<goarch>.
#

goos="$1"
goarch="$2"
directory="$3"

if [[ -z "${goos}" || -z "${goarch}" || ! -d "${directory}" ]]; then
  echo "USAGE: $0 <goos> <goarch> <directory>" >&2
  exit 2
fi

build_setting() {
  local binary="$1"
  local name="$2"
  go version -m "${binary}" | awk -v key="${name}=" \
    '$1 == "build" && index($2, key) == 1 { print substr($2, length(key) + 1) }'
}

fail=0
found=0
for binary in "${directory}"/*; do
  [[ -f "${binary}" ]] || continue
  found=1
  bin_goos="$(build_setting "${binary}" GOOS)"
  bin_goarch="$(build_setting "${binary}" GOARCH)"
  echo "${binary}: ${bin_goos}/${bin_goarch}"
  if [[ "${bin_goos}" != "${goos}" || "${bin_goarch}" != "${goarch}" ]]; then
    echo "::error file=${binary}::built for ${bin_goos}/${bin_goarch}, expected ${goos}/${goarch}"
    fail=1
  fi
done

if [[ "${found}" == "0" ]]; then
  echo "::error::no binaries found in ${directory}"
  exit 1
fi
exit "${fail}"
