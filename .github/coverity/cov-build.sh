#!/bin/bash -xe



ROOT_DIR="$(git rev-parse --show-toplevel)"
source ${ROOT_DIR}/.github/coverity/enviroment.sh


function coverity_build(){
  local NAME=${1/.sh/}
  ${COV_BUILD} --dir "cov/${NAME}" "$1" | tee  ${NAME}.log
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
  cov-build prepare-nmos-cpp.sh
}

function build_grpc(){
  echo "building gRPC"
  cd ${ROOT_DIR}/gRPC
  cov-build compile.sh
}

function build_launcher(){
  echo "building launcher"
  cd ${ROOT_DIR}/launcher
  cov-build "go build -a -o manager cmd/main.go"
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