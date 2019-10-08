#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

if [[ ${#@} < 2 ]]; then
    echo "Usage: $0 generate "
    echo "* version: semver-formatted version for this package"
    echo "* overlay: directory relative to the deploy directory containing a kustomize overlay"
    exit 1
fi

version=$1
overlay=$2

tmp=$(mktemp -d)
trap 'rm -rf "${tmp}"' ERR EXIT

cp -R deploy/ "${tmp}"
# tree "${tmp}"
tmp_overlay="${tmp}/${overlay}"

kubectl kustomize "${tmp_overlay}" > "${tmp_overlay}/manifests/${version}/olm.yaml"

cp -R "${tmp}" "deploy"