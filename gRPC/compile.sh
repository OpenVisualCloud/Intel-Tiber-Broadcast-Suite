#SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
#
#SPDX-License-Identifier: BSD-3-Clause

#!/bin/bash

cmake -S . -B build && cd build && make && cd ..

