#!/bin/bash

# SPDX-License-Identifier: BSD-3-Clause
# SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation

help() {
    echo "Pattern: $0 --source-dir <source_dir> --build-dir <build_dir> --patch-dir <patch_dir> --run-dir <RUN_DIR> [--prepare-only] [--apply-patches] [--build-images] [--run-docker-compose] [--update-submodules]"
    echo "  --source-dir <source_dir>  : Absolute path to directory with source code of repository nmos-cpp 3rd party submodule"
    echo "  --build-dir <build_dir>    : Absolute path to directory with dockerfile and other build files in build-nmos-cpp 3rd party submodule"
    echo "  --patch-dir <patch_dir>    : Absolute path to directory with patches for both 3rd party submodules"
    echo "  --run-dir <run_dir>        : Absolute path to directory with run.sh and docker-compose.yaml"
    echo "  --prepare-only             : Run steps in script that prepares images for nmos but option with running docker containers is not applicable"
    echo "  --build-images             : Build docker images for nmos-client and nmos-registry"
    echo "  --update-submodules        : Update git submodules"
    echo "  --apply-patches            : Apply patches for 3rd party submodules"
    echo "  --run-docker-compose       : Run Docker Compose (nmos-client + nmos-registry + nmos-testing)"
    echo "                               in customized network of bridge type."
    echo "                               Else, by default the <docker run> command will run:"
    echo "                               (nmos-client + nmos-registry without nmos-testing tool container)"
    echo "                               in host network"
    exit 1
}

# Flags:
# Enable/disable to update 3rd party submodules nmos-cpp and build-nmos-cpp 
UPDATE_SUBMODULES=false
# Enable/disable to run containers via docker compose, else it runs using docker run command
RUN_DOCKER_COMPOSE=false
# Enable/disable to build images for nmos-client and nmos-registry
BUILD=false
# Enable/disable to apply patches to 3rd party submodules nmos-cpp and build-nmos-cpp
APPLY_PATCHES=false
# Run script but without running containers
PREPARE=true

while [[ "$#" -gt 0 ]]; do
    case $1 in
        --source-dir)
            # path to source code of repository nmos-cpp 3rd party submodule
            SOURCE_DIR="$2"
            shift 2
            ;;
        --build-dir)
            # path to dockerfile and other build files in build-nmos-cpp 3rd party submodule
            BUILD_DIR="$2"
            shift 2
            ;;
        --patch-dir)
            # path to patches for 3rd party submodules
            PATCH_DIR="$2"
            shift 2
            ;;
        --run-dir)
            # path for run.sh and docker-compose.yaml
            RUN_DIR="$2"
            shift 2
            ;;
        --update-submodules)
            # enabling git submodule updates
            UPDATE_SUBMODULES=true
            shift
            ;;
        --apply-patches)
            # enabling git submodule updates
            APPLY_PATCHES=true
            shift
            ;;
        --run-docker-compose)
            # enabling docker compose run process
            RUN_DOCKER_COMPOSE=true
            shift
            ;;
        --build-images)
            # build process of images for nmos-client and nmos-registry
            BUILD=true
            shift
            ;;
        --prepare-only)
            # build process of images for nmos-client and nmos-registry
            PREPARE=true
            shift
            ;;
        *)
            help
            ;;
    esac
done

# Check provided input paths
if [[ -z "${SOURCE_DIR}" || -z "${BUILD_DIR}" || -z "${PATCH_DIR}" || -z "${RUN_DIR}" ]]; then
    help
fi

# Update 3rd party git submodules
if [[ "${UPDATE_SUBMODULES}" == true ]]; then
    echo "Updating git submodules..."
    git submodule update --init --recursive
fi

# Apply patches
if [[ "${APPLY_PATCHES}" == true ]]; then
    echo "Applying patches for nmos-cpp..."
    cd "${SOURCE_DIR}" || exit
    git apply "${PATCH_DIR}/nmos-cpp.patch"

    echo "Applying patches for build-nmos-cpp..."
    cd "${BUILD_DIR}" || exit
    git apply "${PATCH_DIR}/build-nmos-cpp.patch"
fi

# Copy files from nmos-cpp to build-nmos-cpp
echo "Copying files to build directory..."
cp "${SOURCE_DIR}/Development/nmos-cpp-node/node_implementation.cpp" "${BUILD_DIR}/"
cp "${SOURCE_DIR}/Development/nmos-cpp-node/main.cpp" "${BUILD_DIR}/"
cp -r ../../gRPC/ "${BUILD_DIR}/"

# Build grpc client that is necessary for NMOS node and ffmpeg pipeline communication
echo "Building grpc client that is necessary for NMOS node and ffmpeg pipeline communication..."
cd "${BUILD_DIR}/gRPC" || exit
./compile.sh

if [ $? -ne 0 ]; then
  echo "Error: ./compile.sh failed. Check if all necesarry packages are installed"
  exit 1
fi

# Build NMOS registry and controller
if [[ "${BUILD}" == true ]]; then
    echo "Building NMOS registry and controller..."
    cd "${BUILD_DIR}" || exit
    make build
    make buildnode
fi

if [[ "${PREPARE}" == false ]]; then
    if [[ "${RUN_DOCKER_COMPOSE}" == true ]]; then
        echo "Running Docker Compose from configuration file docker-compose.yaml..."
        cd "${RUN_DIR}" || exit
        docker compose up
    else
        echo "Running Docker Containers from command line..."
        cd "${BUILD_DIR}" || exit
        make run
        make runnode
    fi
fi