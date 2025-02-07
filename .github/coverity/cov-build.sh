#!/bin/bash -e



ROOT_DIR="$(git rev-parse --show-toplevel)"
. ${ROOT_DIR}/.github/coverity/enviroment.sh


function coverity_cpp_build(){
  local FOLDER=${1}
  local SCRIPT=${2}
  ${COVERITY_CPP_BIN_DIR}/cov-build "--dir""${ROOT_DIR}/cov-int/"  "--append-log" "${ROOT_DIR}/${FOLDER}/${SCRIPT}" >  ${FOLDER}.log
  log_info "cov-build ${FOLDER} done"
}

function coverity_other_build(){
  local FOLDER=${1}
  local SCRIPT=${2}
  rm -rf cov/*
  ${COVERITY_OTHER_BIN_DIR}/cov-build "--dir" "${ROOT_DIR}/cov-int/" "--append-log" "${ROOT_DIR}/${FOLDER}/${SCRIPT}" >  ${FOLDER}.log
  log_info "cov-build ${FOLDER} done"
}

function usage(){
  echo " Usage : $0 <BUILD_TYPE>"
  echo " BUILD_TYPE : nmos | nmos-cpp | grpc | launcher"
}

function build_nmos(){
  log_info "building nmos"
}

function build_nmos(){
  log_info "building nmos-cpp"
  cd ${ROOT_DIR}/nmos
  # ./prepare-nmos-cpp.sh
  # echo " cd docker \
  #  ./run.sh\ 
  #  --source-dir <source_dir> \
  #  --build-dir <build_dir> \
  #  --patch-dir <patch_dir> \
  #  --run-dir <RUN_DIR> \
  #  --update-submodules \
  #  --apply-patches \
  #  --build-images \
  #  --prepare-only

  # coverity_cpp_build nmos prepare-nmos-cpp.sh
}

function build_grpc(){
  log_info "building gRPC"
  cd ${ROOT_DIR}/gRPC
  sed -i 's/make -C "${COMPILE_DIR}\/build"/make -B -C "${COMPILE_DIR}\/build"/' compile.sh
  coverity_cpp_build gRPC compile.sh
}

function build_launcher(){
  log_info "building launcher"
  cd ${ROOT_DIR}/launcher
  echo "go build -a -o manager cmd/main.go" > build.sh
  chmod +x build.sh
  coverity_other_build launcher build.sh
}

function build_all(){
  log_info "starting cov-build"
  build_grpc 
  build_launcher 
  build_nmos 
  log_info "All builds have completed"
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