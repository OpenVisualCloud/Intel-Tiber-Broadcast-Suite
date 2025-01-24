#!/bin/bash

# Function to handle errors
handle_error() {
    echo "Error: $1"
    exit 1
}

cd ../gRPC || handle_error "Failed to change directory to gRPC"
./compile.sh

# Change to the target directory
cd ../nmos/nmos-cpp/Development || handle_error "Failed to change directory to nmos/nmos-cpp/Development"

# Create necessary directories
mkdir -p nmos-cpp-node/build || handle_error "Failed to create directory nmos-cpp-node/build"
mkdir -p build || handle_error "Failed to create directory build"

# Copy files with error checking
cp ../../../gRPC/build/libFFmpeg_wrapper_client.a build/ || handle_error "Failed to copy libFFmpeg_wrapper_client.a"
cp ../../../gRPC/build/libhw_grpc_proto.a build/ || handle_error "Failed to copy libhw_grpc_proto.a"
cp ../../../gRPC/config_params.hpp nmos-cpp-node/ || handle_error "Failed to copy config_params.hpp"
cp ../../../gRPC/FFmpeg_wrapper_client.h nmos-cpp-node/ || handle_error "Failed to copy FFmpeg_wrapper_client.h"
cp ../../../gRPC/build/ffmpeg_cmd_wrap.pb.h nmos-cpp-node/build/ || handle_error "Failed to copy ffmpeg_cmd_wrap.pb.h"
cp ../../../gRPC/build/ffmpeg_cmd_wrap.grpc.pb.h nmos-cpp-node/build/ || handle_error "Failed to copy ffmpeg_cmd_wrap.grpc.pb.h"

# Change to the build directory
cd build || handle_error "Failed to change directory to build"

# Set the LIBRARY_PATH environment variable
export LIBRARY_PATH=$(pwd)

# Run cmake with error checking
cmake .. -DCMAKE_BUILD_TYPE:STRING="Release" -DCMAKE_PROJECT_TOP_LEVEL_INCLUDES:STRING="third_party/cmake/conan_provider.cmake" || handle_error "CMake configuration failed"
make -j100
echo "Script completed successfully"