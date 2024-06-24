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
