#!/bin/sh -ex

SCRIPT_DIR="jenkins/scripts"
mkdir -p Trivy
mkdir -p Trivy/dockerfile
chmod a+w Trivy
chmod a+w Trivy/dockerfile
touch Trivy/dockerfile/trivy_clean_reports_source_code
touch Trivy/dockerfile/trivy_clean_reports_source_code_sbom
TRIVY_PASS=true

COMPONENTS_PASS=true
trivy filesystem --severity HIGH,CRITICAL --ignore-unfixed --no-progress --exit-code 1 --scanners vuln --format json -o Trivy/dockerfile/Source_code_Trivy.json --timeout 30m ./Dockerfile  || COMPONENTS_PASS=false
trivy convert --format template --template ${SCRIPT_DIR}/trivy_report_template.tmpl -o Trivy/dockerfile/Source_code_Trivy.txt Trivy/dockerfile/Source_code_Trivy.json
if [ "$COMPONENTS_PASS" = "true" ]; then
    echo Trivy/dockerfile/Trivy_idc.txt >> Trivy/dockerfile/trivy_clean_reports_source_code
else
    TRIVY_PASS=false
fi

SPDX_PASS=true
trivy filesystem --no-progress --exit-code 1 --format spdx -o Trivy/dockerfile/Trivy_sbom.spdx ./Dockerfile || SPDX_PASS=false
if [ "$SPDX_PASS" = "true" ]; then
    echo Trivy/dockerfile/Trivy_sbom.spdx >> Trivy/dockerfile/trivy_clean_reports_source_code_sbom
else
    TRIVY_PASS=false
fi

echo $TRIVY_PASS>Trivy/dockerfile/trivy_source_code_pass

echo "creating summary ..."
python3 jenkins/scripts/trivy_dockerfile_summary.py Trivy/dockerfile/Source_code_Trivy.json Trivy/dockerfile_scan_summary.csv
column -t -s, Trivy/dockerfile_scan_summary.csv > Trivy/dockerfile_scan_summary.txt

echo "images scanning done"
