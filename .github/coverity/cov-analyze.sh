#!/bin/bash 

ROOT_DIR="$(git rev-parse --show-toplevel)"
. ${ROOT_DIR}/.github/coverity/enviroment.sh

cd ${ROOT_DIR}

MERGE_CPP_DIR="cov-cpp-all"
mkdir -p ${MERGE_CPP_DIR}
${COVERITY_CPP_BIN_DIR}/cov-manage-emit --dir ${ROOT_DIR}/${MERGE_CPP_DIR} add ${ROOT_DIR}/gRPC/cov
${COVERITY_CPP_BIN_DIR}/cov-manage-emit --dir ${ROOT_DIR}/${MERGE_CPP_DIR} add ${ROOT_DIR}/nmos/cov

MERGE_OTHER_DIR="cov-other-all"
mkdir -p ${MERGE_OTHER_DIR}
${COVERITY_OTHER_BIN_DIR}/cov-manage-emit --dir ${ROOT_DIR}/${MERGE_OTHER_DIR} add ${ROOT_DIR}/launcher/cov

