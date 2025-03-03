#!/bin/bash

#SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
#
#SPDX-License-Identifier: BSD-3-Clause

COMPILE_DIR="$(readlink -f "$(dirname -- "${BASH_SOURCE[0]}")")"
num_proc=$(nproc)

cmake -S "${COMPILE_DIR}" -B "${COMPILE_DIR}/build" && \
make -j"$num_proc" -C "${COMPILE_DIR}/build"

if [[ "$1" == "--unit_testing" ]]; then
cmake -S "${COMPILE_DIR}/unit_test" -B "${COMPILE_DIR}/unit_test/build" && \
make -j"$num_proc" -C "${COMPILE_DIR}/unit_test/build" && ctest --test-dir "${COMPILE_DIR}/unit_test/build"
fi
