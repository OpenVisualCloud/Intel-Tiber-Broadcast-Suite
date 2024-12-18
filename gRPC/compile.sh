#!/bin/bash

#SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
#
#SPDX-License-Identifier: BSD-3-Clause

COMPILE_DIR="$(readlink -f "$(dirname -- "${BASH_SOURCE[0]}")")"

cmake -S "${COMPILE_DIR}" -B "${COMPILE_DIR}/build" && \
make -C "${COMPILE_DIR}/build"
