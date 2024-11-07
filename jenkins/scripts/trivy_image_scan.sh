#!/bin/bash
set -x

SCRIPT_DIR="$(readlink -f "$(dirname -- "${BASH_SOURCE[0]}")")"
if [ -z "${SCRIPT_DIR}" ] || [ ! -d "${SCRIPT_DIR}" ]; then
    SCRIPT_DIR="$(pwd)"
fi

REPO_DIR="$(readlink -f "${SCRIPT_DIR}/../..")"
if [ -z "${REPO_DIR}" ] || [ ! -d "${REPO_DIR}" ]; then
    REPO_DIR="${SCRIPT_DIR}/../.."
fi

. "${REPO_DIR}/scripts/common.sh"

# get the latest video_production_image.tar.gz
SDB_DOCKER_IMAGE="${1}"
IMAGE_LOG="Trivy_video_production_image"

mkdir -p "${REPO_DIR}/Trivy/image/"
chmod -R a+w "${REPO_DIR}/Trivy"

trivy image \
    --exit-code 2 \
    --no-progress    \
    --format spdx    \
    --input "${SDB_DOCKER_IMAGE}" \
    -o "${REPO_DIR}/Trivy/image/${IMAGE_LOG}.spdx" 


trivy image --exit-code 1 --timeout 15m \
    --severity HIGH,CRITICAL \
    --ignore-unfixed \
    --no-progress    \
    --scanners vuln  \
    --format table    \
    --input "${SDB_DOCKER_IMAGE}" \
    -o "${REPO_DIR}/Trivy/image/${IMAGE_LOG}.txt" 