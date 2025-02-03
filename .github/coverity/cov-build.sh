#!/bin/bash -xe



ROOT_DIR="$(git rev-parse --show-toplevel)"
. ${ROOT_DIR}/.github/coverity/enviroment.sh


function coverity_build(){
  local FOLDER=${1}
  local SCRIPT=${2}
  ${COV_BUILD} "--dir" "cov/" "${ROOT_DIR}/${FOLDER}/${SCRIPT}" >  ${NAME}.log
  echo "cov build ${FOLDER} done"
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
  coverity_build nmos prepare-nmos-cpp.sh
}

function build_grpc(){
  echo "building gRPC"
  cd ${ROOT_DIR}/gRPC
  sed -i 's/make -C "${COMPILE_DIR}\/build"/make -B -C "${COMPILE_DIR}\/build"/' compile.sh
  coverity_build gRPC compile.sh
}

function build_launcher(){
  echo "building launcher"
  cd ${ROOT_DIR}/launcher
  echo "go build -a -o manager cmd/main.go" > build.sh
  chmod +x build.sh
  coverity_build launcher build.sh
}

function build_all(){
  build_nmos &
  build_nmos_cpp &
  build_grpc &
  build_launcher &
  echo "waiting for all builds to complete"
  wait
  echo "All builds have completed"
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
  all)
    build_all
    ;;
  *)
    usage
    exit 1
    ;;
esac