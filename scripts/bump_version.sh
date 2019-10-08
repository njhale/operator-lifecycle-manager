#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

yq w -i deploy/$(target)/kustomization.yaml 'images.[0].digest' $(olmsha)
yq w -i deploy/bases/default/crs.yaml 'metadata.labels."olm.version"' $(ver)
yq w -i deploy/$(target)/crs.patch.yaml 'spec.version' $(ver)