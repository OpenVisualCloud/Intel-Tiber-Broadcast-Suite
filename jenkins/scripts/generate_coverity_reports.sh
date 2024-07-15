#!/bin/bash -x


echo "generating reports ..."

export JAVA_HOME=${HOME}/cov-reports-2023.3.4/jre/bin
export PATH="$PATH:${HOME}/cov-reports-2023.3.4/bin"
export no_proxy=${no_proxy},*.intel.com

CONFIG_FILES=($(find . -name "cov_report*.yaml"  -maxdepth 1))

for CONFIG_FILE in "${CONFIG_FILES[@]}"
do
  echo "generating  security report from ${CONFIG_FILE} ..."
  STREAM=${CONFIG_FILE/cov_report_/}
  STREAM=${CONFIG_FILE/.yaml/}
  cov-generate-security-report ${CONFIG_FILE} \
                       --output ${STREAM}-security-report.pdf \
                       --user ${USERNAME} --password env:PASSWORD
  echo "generating cvss report from ${CONFIG_FILE} ..."
  cov-generate-cvss-report ${CONFIG_FILE} \
                       --output ${STREAM}-cvss-report.pdf \
                       --report \
                       --user ${USERNAME} --password env:PASSWORD

done
