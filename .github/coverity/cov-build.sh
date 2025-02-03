#!/bin/bash -e



ROOT_DIR="$(git rev-parse --show-toplevel)"
. ${ROOT_DIR}/.github/coverity/enviroment.sh


function coverity_build(){
  local NAME=${1}
  local SCRIPT=${2}
  ${COV_BUILD} "--dir cov/${NAME}" "${SCRIPT}" | tee  ${NAME}.log
}


function usage(){
  echo " Usage : $0 <BUILD_TYPE>"
  echo " BUILD_TYPE : nmos | nmos-cpp | grpc | launcher"
}

function build_nmos(){
  echo "building nmos"
}

function build_nmos_cpp(){
  echo "building nmos-cpp"
  cd ${ROOT_DIR}/nmos
  coverity_build nmos-cpp prepare-nmos-cpp.sh
}

function build_grpc(){
  echo "building gRPC"
  cd ${ROOT_DIR}/gRPC
  coverity_build grpc compile.sh
}

function build_launcher(){
  echo "building launcher"
  cd ${ROOT_DIR}/launcher
  coverity_build launcher "go build -a -o manager cmd/main.go"
}

if [ $# -ne 1 ]; then
  usage
  exit 1
fi

case $1 in
  nmos)
    build_nmos
    ;;
  nmos-cpp)
    build_nmos_cpp
    ;;
  grpc)
    build_grpc
    ;;
  launcher)
    build_launcher
    ;;
  *)
    usage
    exit 1
    ;;
esac