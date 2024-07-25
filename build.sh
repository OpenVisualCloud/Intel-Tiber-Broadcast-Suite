#!/bin/bash

#
# SPDX-License-Identifier: BSD-3-Clause
# Copyright(©) 2024 Intel Corporation
# Intel® Tiber™ Broadcast Suite
#

set -eo pipefail
SCRIPT_DIR="$(readlink -f "$(dirname -- "${BASH_SOURCE[0]}")")"

ENV_PROXY_ARGS=()
while IFS='' read -r line; do ENV_PROXY_ARGS+=("--build-arg"); ENV_PROXY_ARGS+=("$line"); done < <(compgen -e | grep -E "_(proxy|PROXY)")

IMAGE_CACHE_REGISTRY="${IMAGE_CACHE_REGISTRY:-docker.io}"
IMAGE_REGISTRY="${IMAGE_REGISTRY:-docker.io}"
IMAGE_TAG="${IMAGE_TAG:-latest}"

docker buildx build "${ENV_PROXY_ARGS[@]}" "$@" \
    --progress=plain \
    --network=host \
    --build-arg IMAGE_CACHE_REGISTRY="${IMAGE_CACHE_REGISTRY}" \
    -t "${IMAGE_REGISTRY}/tiber-broadcast-suite:${IMAGE_TAG}" \
    -f "${SCRIPT_DIR}/Dockerfile" \
    "${SCRIPT_DIR}"

docker tag "${IMAGE_REGISTRY}/tiber-broadcast-suite:${IMAGE_TAG}" video_production_image

if docker images | grep -q "mtl-manager\s*latest"; then
     echo -e '\e[32mmtl-manager:latest image exists. skipping the build\e[0m'
     exit 0
fi

if [ -z "$mtl_source_code" ]; then
    mtl_source_code=$(find "$HOME" -type d -name "Media-Transport-Library" -print -quit)

    if [ -n "$mtl_source_code" ]  && [ -f "$mtl_source_code/manager/Dockerfile" ]; then
        echo -e '\e[33mmtl_source_code variable is set to '"$mtl_source_code"'.\e[0m'
    else
        echo -e '\e[31m'"Media-Transport-Library manager directory not found in $HOME directory."'\e[0m'
        echo -e '\e[33mMTL manager not installed\e[0m'
        exit 0
    fi
fi

cd "$mtl_source_code"/manager
docker build -t mtl-manager:latest .
