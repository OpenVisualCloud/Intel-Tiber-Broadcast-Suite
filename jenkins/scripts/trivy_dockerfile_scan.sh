#!/bin/sh -ex

SCRIPT_DIR="jenkins/scripts"
mkdir -p Trivy
mkdir -p Trivy/dockerfile
chmod a+w Trivy
chmod a+w Trivy/dockerfile
touch Trivy/dockerfile/trivy_clean_reports_source_code
touch Trivy/dockerfile/trivy_clean_reports_source_code_sbom

DOCKER_FILE=${1}

trivy filesystem \
      --severity HIGH,CRITICAL \
      --ignore-unfixed \
      --no-progress \
      --exit-code 0 \
      --scanners vuln \
      --format table \
      -o Trivy/dockerfile/Source_code_Trivy.txt --timeout 30m ${DOCKER_FILE}  


trivy filesystem \
      --no-progress \
      --exit-code 0 \
      --format spdx \
      -o Trivy/dockerfile/Trivy_sbom.spdx ${DOCKER_FILE}

# echo "creating summary ..."
# python3 jenkins/scripts/trivy_dockerfile_summary.py Trivy/dockerfile/Source_code_Trivy.json Trivy/dockerfile_scan_summary.csv
# column -t -s, Trivy/dockerfile_scan_summary.csv > Trivy/dockerfile_scan_summary.txt

echo "images scanning done"
