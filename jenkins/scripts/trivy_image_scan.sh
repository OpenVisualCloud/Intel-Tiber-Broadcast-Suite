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
touch "${REPO_DIR}Trivy/image/trivy_clean_reports_images" "${REPO_DIR}/Trivy/image/trivy_clean_reports_images_sbom"
chmod -R a+w "${REPO_DIR}/Trivy"

trivy image --exit-code 1 --timeout 15m \
    --db-repository public.ecr.aws/aquasecurity/trivy-db:2 \
    --severity HIGH,CRITICAL \
    --ignore-unfixed \
    --no-progress    \
    --scanners vuln  \
    --format table    \
    -o "${REPO_DIR}/Trivy/image/${IMAGE_LOG}.txt" "${SDB_DOCKER_IMAGE}"



trivy image \
    --db-repository public.ecr.aws/aquasecurity/trivy-db:2 \
    --exit-code 2 \
    --no-progress    \
    --format spdx    \
    -o "${REPO_DIR}/Trivy/image/${IMAGE_LOG}.spdx" "${SDB_DOCKER_IMAGE}"

