#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

BUILD_SCRIPTS=($(ls .github/tools/coverity/*build_cmd.sh))

for SCRIPT in "${BUILD_SCRIPTS[@]}"
do 
  # cov-build --dir . $SCRIPT
done