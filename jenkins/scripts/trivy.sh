#!/bin/sh -x

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd)
cd $SCRIPT_DIR

rm -rf Trivy
mkdir Trivy

# get the latest video_production_image 
IMAGE_TAR=${1}
# vulnerability scans
trivy image --no-progress \
            --exit-code 0 \
            --format json \
            -o Trivy/Trivy_vulnerability.Dockerfile.json \
            --input ${IMAGE_TAR%%.tar}
# spdx licence scans 
trivy image --no-progress \
            --exit-code 0 \
            --format spdx-json \
            -o Trivy/Trivy_spdx.Dockerfile.json \
            --input ${IMAGE_TAR%%.tar}
