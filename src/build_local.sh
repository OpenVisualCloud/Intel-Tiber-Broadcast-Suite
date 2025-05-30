#!/bin/bash

#SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
#
#SPDX-License-Identifier: BSD-3-Clause

set -eo pipefail

SCRIPT_DIR="$(readlink -f "$(dirname -- "${BASH_SOURCE[0]}")")"
GREEN='\033[0;32m'
NC='\033[0m'
NMOS_CPP_VERSION=f54971298c47a633969e9e9adac824b56fc08da7
MY_INSTALL_DIR=$HOME/.local
num_proc=$(nproc)

# Function to display help message
show_help() {
  echo "Usage: $0 [-ut] [--build_type <type>] [-h|--help]"
  echo ""
  echo "Options:"
  echo "  -ut                Build with unit tests enabled"
  echo "  --build_type <type> Specify the build type (e.g., Debug, Release)"
  echo "  -h, --help         Show this help message and exit"
}

UT_OPTION=$1
BUILD_TYPE="Release"

# Parse input parameters
while [[ "$#" -gt 0 ]]; do
    case $1 in
        -ut) UT_OPTION="$1"; shift ;;
        --build_type) BUILD_TYPE="$2"; shift 2 ;;
        -h|--help) show_help; exit 0 ;;
        *) echo "Unknown parameter passed: $1"; show_help; exit 1 ;;
    esac
done

# Function to handle errors
handle_error() {
    echo "Error: $1"
    exit 1
}

function build_grpc() {
    if [ ! -d "grpc" ]; then
        mkdir -p $MY_INSTALL_DIR
        git clone --recurse-submodules -b v1.58.0 --depth 1 --shallow-submodules https://github.com/grpc/grpc
        mkdir -p grpc/cmake/build
    fi

    export PATH="$MY_INSTALL_DIR/bin:$PATH"
    cd grpc

    pushd cmake/build
    cmake -DgRPC_INSTALL=ON \
        -DgRPC_BUILD_TESTS=OFF \
        -DCMAKE_CXX_STANDARD=17 \
        -DCMAKE_INSTALL_PREFIX=$MY_INSTALL_DIR \
        -DCMAKE_BUILD_TYPE=$BUILD_TYPE \
        ../..
    make -j"$num_proc"
    make install
    popd
}

function build_grpc_based_ffmpeg_app() {
    cd "${SCRIPT_DIR}"/gRPC || handle_error "Failed to change directory to gRPC"
    if [ "$UT_OPTION" == "-ut" ]; then
        ./compile.sh --unit_testing  --build_type $BUILD_TYPE
    else
        ./compile.sh  --build_type $BUILD_TYPE
    fi
}

function build_nmos_cpp_library () {
    cd "${SCRIPT_DIR}"/nmos

    if [ ! -d "nmos-cpp" ]; then
        curl --output - -s -k https://codeload.github.com/sony/nmos-cpp/tar.gz/${NMOS_CPP_VERSION} | tar zxvf - -C . && \
        mv ./nmos-cpp-${NMOS_CPP_VERSION} ./nmos-cpp
        mkdir -p nmos-cpp/Development/build
    fi

    cd nmos-cpp/Development/build

    cmake .. -DNMOS_CPP_USE_SUPPLIED_JSON_SCHEMA_VALIDATOR=ON \
    -DNMOS_CPP_USE_SUPPLIED_JWT_CPP=ON \
    -DNMOS_CPP_BUILD_EXAMPLES=OFF \
    -DNMOS_CPP_BUILD_TESTS=OFF \
    -DCMAKE_BUILD_TYPE=$BUILD_TYPE && \
    make -j"$num_proc" && \
    make install
}

function build_nmos_node() {
    cd "${SCRIPT_DIR}"/nmos/nmos-node
    mkdir -p build && cd build
    if [ "$UT_OPTION" == "-ut" ]; then
        cmake .. -DENABLE_UNIT_TESTS=ON -DCMAKE_BUILD_TYPE=$BUILD_TYPE
    else
        cmake .. -DENABLE_UNIT_TESTS=OFF -DCMAKE_BUILD_TYPE=$BUILD_TYPE
    fi
    make -j"$num_proc"
}

build_grpc
build_grpc_based_ffmpeg_app
build_nmos_cpp_library
build_nmos_node

echo
echo -e ${GREEN}Build finished sucessfuly ${NC}
echo