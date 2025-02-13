cmake_minimum_required(VERSION 3.24 FATAL_ERROR)

# project name
project(nmos-node)

include_directories("${CMAKE_CURRENT_BINARY_DIR}"
                    ../../common
                    ../../gRPC/build
                    /usr/local/include
                    /usr/include)

link_directories(../../gRPC/build
                /usr/lib
                /usr/local/lib)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -pthread ")

#-Wall -Wstrict-aliasing -fstrict-aliasing -Wextra -Wno-unused-parameter -pedantic -Wno-long-long

# nmos-cpp-node executable
set(NMOS_CPP_NODE_SOURCES
    main.cpp
    node_implementation.cpp
    intel_config_parser.cpp
    )
set(NMOS_CPP_NODE_HEADERS
    node_implementation.h
    intel_config_parser.h
    )

find_package(OpenSSL REQUIRED)
find_package(Boost REQUIRED COMPONENTS thread)
find_package(nlohmann_json REQUIRED)

# Find gRPC installation
# Looks for gRPCConfig.cmake file installed by gRPC's cmake installation.
set(protobuf_MODULE_COMPATIBLE TRUE)
find_package(Protobuf CONFIG REQUIRED)
message(STATUS "Using protobuf ${Protobuf_VERSION}")
set(_PROTOBUF_LIBPROTOBUF protobuf::libprotobuf)
set(_REFLECTION gRPC::grpc++_reflection)
find_package(gRPC CONFIG REQUIRED)
message(STATUS "Using gRPC ${gRPC_VERSION}")

add_executable(
    nmos-cpp-node
    ${NMOS_CPP_NODE_SOURCES}
    ${NMOS_CPP_NODE_HEADERS}
    ../../common/config_params.hpp
    ../../common/FFmpeg_wrapper_client.h
)

source_group("Source Files" FILES ${NMOS_CPP_NODE_SOURCES})
source_group("Header Files" FILES ${NMOS_CPP_NODE_HEADERS})

target_link_libraries(
    nmos-cpp-node
    nmos-cpp
    cpprest
    mdns
    OpenSSL::SSL 
    OpenSSL::Crypto
    Boost::thread
    nmos_is04_schemas
    nmos_is05_schemas
    nmos_is08_schemas
    nmos_is09_schemas
    nmos_is10_schemas
    nmos_is12_schemas
    nlohmann_json::nlohmann_json
    FFmpeg_wrapper_client
    protobuf::libprotobuf
    hw_grpc_proto
    gRPC::grpc++
    )
# root directory to find e.g. nmos-cpp-node/node_implementation.h
target_include_directories(nmos-cpp-node PRIVATE
    ${CMAKE_CURRENT_SOURCE_DIR}
    )

#list(APPEND NMOS_CPP_TARGETS nmos-cpp-node)
