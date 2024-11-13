#!/bin/bash -x

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

rm -rf Hadolint
mkdir -p Hadolint

hadolint -v -V --config ${SCRIPT_DIR}/hadolint_config.yaml > Hadolint/hadolint-Dockerfile
hadolint --config ${SCRIPT_DIR}/hadolint_config.yaml --no-color Dockerfile 2>&1 > Hadolint/Dockerfile.log
