#!/bin/bash -x

BASH_FILES=$(find . -name "*.sh")
OPTIONS="-s bash"
OUT_FILE="shellcheck_output.log"

rm -rf shellcheck_logs
mkdir -p shellcheck_logs

shellcheck ${OPTIONS} ${BASH_FILES} | tee  shellcheck_logs/${OUT_FILE}
