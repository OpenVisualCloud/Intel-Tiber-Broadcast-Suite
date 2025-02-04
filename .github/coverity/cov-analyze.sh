#!/bin/bash 

ROOT_DIR="$(git rev-parse --show-toplevel)"
. ${ROOT_DIR}/.github/coverity/enviroment.sh

cd ${ROOT_DIR}

COV_BUILD_OUT=($(find . -name cov))
MERGE_DIR="cov-all"
mkdir -p ${MERGE_DIR}


for OUT in "${COV_BUILD_OUT[@]}"; do
   log_info "Merging ${OUT} to ${MERGE_DIR}"
   ${COVERITY_OTHER_BIN_DIR}/cov-manage-emit --dir ${ROOT_DIR}/${MERGE_DIR} add ${ROOT_DIR}/${OUT}

done


