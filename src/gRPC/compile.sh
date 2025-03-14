#!/bin/bash

#SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
#
#SPDX-License-Identifier: BSD-3-Clause

COMPILE_DIR="$(readlink -f "$(dirname -- "${BASH_SOURCE[0]}")")"
num_proc=$(nproc)
BUILD_TYPE="Release"
UNIT_TESTING=false

# Parse input parameters
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --build_type) BUILD_TYPE="$2"; shift 2 ;;
        --unit_testing) UNIT_TESTING=true; shift ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
done

cmake -S "${COMPILE_DIR}" -B "${COMPILE_DIR}/build" -DCMAKE_BUILD_TYPE=$BUILD_TYPE && \
make -j"$num_proc" -C "${COMPILE_DIR}/build"

if [[ "$UNIT_TESTING" == true ]]; then
    cmake -S "${COMPILE_DIR}/unit_test" -B "${COMPILE_DIR}/unit_test/build" -DCMAKE_BUILD_TYPE=$BUILD_TYPE && \
    make -j"$num_proc" -C "${COMPILE_DIR}/unit_test/build" && ctest --test-dir "${COMPILE_DIR}/unit_test/build"
fi