#!/bin/bash -x

COV_VER="$(ls /opt/coverity/bin/)"
COVERITY_BIN_DIR="/opt/coverity/bin/$COV_VER/bin"
export JAVA_HOME="/opt/coverity/bin/$COV_VER/jre/bin/java"
# Array to hold the PIDs of background jobs
BUILD_SCRIPTS=($(ls .github/tools/coverity/*build_cmd.sh))

PIDS=()
mkdir -p cov
rm -rf cov/*

for SCRIPT in "${BUILD_SCRIPTS[@]}"; do 
  # Run cov-build in the background and store the PID
  cov-build --dir "cov/${NAME/_build_cmd.sh/}" "$SCRIPT" > ${SCRIPT/.sh/.log} &
  PIDS+=($!)
done

# Wait for all background jobs to complete
for PID in "${PIDS[@]}"; do wait $PID; done

echo "All cov-build processes have completed."
