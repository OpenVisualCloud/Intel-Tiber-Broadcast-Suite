#!/bin/bash
set -x



# get the latest video_production_image.tar.gz
SDB_DOCKER_IMAGE="${1}"
IMAGE_LOG="Trivy_video_production_image"

mkdir -p "Trivy/image/"
chmod -R a+w "${REPO_DIR}/Trivy"

trivy image --exit-code 0 --timeout 15m \
    --severity HIGH,CRITICAL \
    --ignore-unfixed \
    --no-progress    \
    --scanners vuln  \
    --format table    \
    --input "${SDB_DOCKER_IMAGE}" \
    -o "Trivy/image/${IMAGE_LOG}.txt"

trivy image --exit-code 0 \
    --no-progress    \
    --format spdx    \
    --input "${SDB_DOCKER_IMAGE}" \
    -o Trivy/image/${IMAGE_LOG}.spdx"

# prompt "Creating Intel--Tiber-Broadcast-Suite summary."

# python3 "${REPO_DIR}/jenkins/scripts/trivy_images_summary.py" "${REPO_DIR}/Trivy/image/${IMAGE_LOG}.json" "${REPO_DIR}/Trivy/images_scan_summary.csv"
# column -t -s, "${REPO_DIR}/Trivy/images_scan_summary.csv" > "${REPO_DIR}/Trivy/images_scan_summary.txt"

prompt "Trivy Scanning of Intel--Tiber-Broadcast-Suite done." 
chmod -R a+rw "${REPO_DIR}/Trivy"
