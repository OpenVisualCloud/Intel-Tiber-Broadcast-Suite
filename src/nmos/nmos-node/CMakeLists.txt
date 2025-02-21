#
# SPDX-FileCopyrightText: Copyright (c) 2025 Intel Corporation
#
# SPDX-License-Identifier: BSD-3-Clause
#

cmake_minimum_required(VERSION 3.17 FATAL_ERROR)

# project name
project(bcs-nmos-node)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -pthread ")

# Find gRPC installation
# Looks for gRPCConfig.cmake file installed by gRPC's cmake installation.
set(protobuf_MODULE_COMPATIBLE TRUE)
find_package(Protobuf CONFIG REQUIRED)
message(STATUS "Using protobuf ${Protobuf_VERSION}")
set(_PROTOBUF_LIBPROTOBUF protobuf::libprotobuf)
set(_REFLECTION gRPC::grpc++_reflection)
find_package(gRPC CONFIG REQUIRED)
message(STATUS "Using gRPC ${gRPC_VERSION}")

find_package(nmos-cpp REQUIRED)
find_package(websocketpp REQUIRED)

# nmos-cpp-node executable
set(NMOS_CPP_NODE_SOURCES
    node_implementation.cpp
    intel_config_parser.cpp
)

set(NMOS_CPP_NODE_HEADERS
    node_implementation.h
    intel_config_parser.h
)

include_directories(
    ../../common
    ../../gRPC/build
)

link_directories(../../gRPC/build)

add_library(bcs-nmos-node-lib ${NMOS_CPP_NODE_SOURCES} ${NMOS_CPP_NODE_HEADERS})

target_link_libraries(
    bcs-nmos-node-lib
    nmos-cpp::nmos-cpp
    FFmpeg_wrapper_client
    protobuf::libprotobuf
    hw_grpc_proto
    gRPC::grpc++
)

add_executable(
    bcs-nmos-node
    main.cpp
)

target_link_libraries(
    bcs-nmos-node
    bcs-nmos-node-lib
)

# enable or disable the unit test suite
set(BUILD_TESTING ON CACHE BOOL "Build test suite application")

if(BUILD_TESTING)
    add_subdirectory(unit_tests)
endif()
