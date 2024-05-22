#!/bin/sh -x

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd)
cd $SCRIPT_DIR

rm -rf Trivy
mkdir Trivy

# get the latest video_production_image 
IMAGE_TAR=${1}
trivy image --no-progress \
            --exit-code 1 \
            --format spdx \
            -o Trivy/Trivy_Dockerfile.spdx \
            --input ${IMAGE_TAR%%.tar}
