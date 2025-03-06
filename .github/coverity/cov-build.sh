#!/bin/bash -e



ROOT_DIR="$(git rev-parse --show-toplevel)"
. ${ROOT_DIR}/.github/coverity/enviroment.sh


function coverity_cpp_build(){
  local SCRIPT=${1}
  ${COVERITY_CPP_BIN_DIR}/cov-build "--dir" "${ROOT_DIR}/cov-int/" "--append-log" "$(pwd)/${SCRIPT}" > "$(basename $(pwd)).log"
  log_info "cov-build ${FOLDER} done"
}

function coverity_other_build(){
  local SCRIPT=${1}
  rm -rf cov/*
  ${COVERITY_OTHER_BIN_DIR}/cov-build "--dir" "${ROOT_DIR}/cov-int/" "--append-log" "$(pwd)/${SCRIPT}" > "$(basename $(pwd)).log"
  log_info "cov-build ${FOLDER} done"
}

function usage(){
  echo " Usage : $0 <BUILD_TYPE>"
  echo " BUILD_TYPE : all | grpc | launcher"
}


function build_grpc(){
  log_info "building gRPC"
  cd ${ROOT_DIR}/src/gRPC
  sed -i 's/make -C "${COMPILE_DIR}\/build"/make -B -C "${COMPILE_DIR}\/build"/' compile.sh
  coverity_cpp_build compile.sh
}

function build_launcher(){
  log_info "building launcher"
  cd ${ROOT_DIR}/launcher
  echo "go build -a -o manager cmd/main.go" > build.sh
  chmod +x build.sh
  coverity_other_build build.sh
}

function build_all(){
  log_info "starting cov-build"
  build_grpc 
  build_launcher 
  log_info "All builds have completed"
}

if [ $# -ne 1 ]; then
  usage
  exit 1
fi

case $1 in
  grpc)
    build_grpc
    ;;
  launcher)
    build_launcher
    ;;
  all)
    build_all
    ;;
  *)
    usage
    exit 1
    ;;
esac
