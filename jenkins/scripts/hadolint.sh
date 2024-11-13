#!/bin/bash -x

rm -rf Hadolint
mkdir -p Hadolint

hadolint -v -V --config jenkins/scripts/hadolint_config.yaml > Hadolint/hadolint-Dockerfile
hadolint --config jenkins/scripts/hadolint_config.yaml --no-color Dockerfile 2>&1 > Hadolint/Dockerfile.log 
exit 0
