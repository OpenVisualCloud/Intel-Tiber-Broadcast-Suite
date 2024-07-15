#!/bin/sh -ex

# get the latest video_production_image.tar.gz
IMAGE=${1}
IMAGE_LOG="Trivy_video_production_image"
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd)

mkdir -p Trivy
mkdir -p Trivy/image/
chmod a+w Trivy
touch Trivy/image/trivy_clean_reports_images
touch Trivy/image/trivy_clean_reports_images_sbom
TRIVY_PASS=true

IMAGE_PASS=true
trivy image --severity HIGH,CRITICAL --ignore-unfixed --no-progress --exit-code 1 --scanners vuln --format json -o Trivy/image/${IMAGE_LOG}.json --input ${IMAGE} --timeout 15m || IMAGE_PASS=false
trivy convert --format template --template jenkins/scripts/trivy_report_template.tmpl -o Trivy/image/${IMAGE_LOG}.txt Trivy/image/${IMAGE_LOG}.json
if [ "$IMAGE_PASS" = "true" ]; then
    echo Trivy/${IMAGE_LOG}.txt >> Trivy/image/trivy_clean_reports_images
else
    TRIVY_PASS=false
fi
SPDX_PASS=true
trivy image --no-progress --exit-code 1 --format spdx -o Trivy/image/${IMAGE_LOG}.spdx --input ${IMAGE} || SPDX_PASS=false
if [ "$SPDX_PASS" = "true" ]; then
    echo Trivy/${IMAGE_LOG}.spdx >> Trivy/image/trivy_clean_reports_images_sbom
else
    TRIVY_PASS=false
fi

echo $TRIVY_PASS>Trivy/image/trivy_images_pass

echo "creating summary ..."
python3 jenkins/scripts/trivy_images_summary.py Trivy/image/${IMAGE_LOG}.json Trivy/images_scan_summary.csv
column -t -s, Trivy/images_scan_summary.csv > Trivy/images_scan_summary.txt

echo "images scanning done"
