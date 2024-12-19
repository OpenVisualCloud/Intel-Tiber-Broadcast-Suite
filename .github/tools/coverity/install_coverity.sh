#!/bin/bash

PLATFORM="linux64"

PROJECT_NAME="Intel-Tiber-Broadcast-Suite"
DOWNLOAD_DIR="/tmp/cov-${LANGUAGE}"

mkdir -p /tmp/cov
echo "download coverity binaries"
HASH=$(curl https://scan.coverity.com/download/${LANGUAGE}/${PLATFORM} --no-progress-meter --data "token=${TOKEN}&project=${PROJECT_NAME}&md5=1")


INSTALLER_FILE="${DOWNLOAD_DIR}/cov-analysis.tar.gz-${HASH}"
if [ ! -f "${INSTALLER_FILE}" ]; then
  echo " installer not found, downloading..."
  curl https://scan.coverity.com/download/${LANGUAGE}/${PLATFORM} \
    --output  cov-analysis.tar.gz \
    --data "token=${TOKEN}&project=${PROJECT_NAME}"

  echo "extracting installer"
  mkdir -p /tmp/cov-analysis
  tar -xzf ${INSTALLER_FILE} --strip 1 -C /tmp/cov-${LANGUAGE}-analysis
  mv cov-analysis.tar.gz INSTALLER_FILE
  echo "export PATH=\$PATH:/tmp/cov-${LANGUAGE}-analysis/bin" >> ~/.bashrc
else
  echo "installer already downloaded, nothing to do"
fi
