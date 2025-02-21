#
# SPDX-FileCopyrightText: Copyright (c) 2025 Intel Corporation
#
# SPDX-License-Identifier: BSD-3-Clause
#

cmake_minimum_required(VERSION 3.8)

# Project name
project(NmosNodeTests)

# Set C++ standard
set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED True)

# Add gRPC and Protocol Buffers
find_package(Protobuf CONFIG REQUIRED)
find_package(gRPC CONFIG REQUIRED)

# Include FetchContent module
include(FetchContent)

# Download and install Google Test
FetchContent_Declare(
  googletest
  URL https://github.com/google/googletest/archive/refs/tags/v1.15.0.tar.gz
)
FetchContent_MakeAvailable(googletest)

# Ensure GoogleTest is built with C++20
set_target_properties(gtest PROPERTIES CXX_STANDARD 20)
set_target_properties(gtest_main PROPERTIES CXX_STANDARD 20)

# Include directories
include_directories(
    ${CMAKE_CURRENT_SOURCE_DIR}/../../app
    ${CMAKE_CURRENT_SOURCE_DIR}/../../common
    ${CMAKE_CURRENT_SOURCE_DIR}/..
)

# Source files
set(NMOS_NODE_SOURCES
    ${CMAKE_CURRENT_SOURCE_DIR}/config_parser_tests.cc
)

# Add the executable
add_executable(NmosNodeTests ${NMOS_NODE_SOURCES})

# Link libraries
target_link_libraries(NmosNodeTests
  gtest
  gtest_main
  bcs-nmos-node-lib
)
