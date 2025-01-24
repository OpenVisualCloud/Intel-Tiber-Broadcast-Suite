cmake_minimum_required(VERSION 3.8)

# Project name
project(FFmpegWrapperServiceTests)

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
# For Windows: Prevent overriding the parent project's compiler/linker settings
set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
FetchContent_MakeAvailable(googletest)

FetchContent_Declare(json URL https://github.com/nlohmann/json/releases/download/v3.11.3/json.tar.xz)
FetchContent_MakeAvailable(json)

# Ensure GoogleTest is built with C++20
set_target_properties(gtest PROPERTIES CXX_STANDARD 20)
set_target_properties(gtest_main PROPERTIES CXX_STANDARD 20)

# Include directories
include_directories(${CMAKE_CURRENT_SOURCE_DIR})
include_directories(${googletest_SOURCE_DIR}/googletest/include)
include_directories(${googletest_SOURCE_DIR}/googlemock/include)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../build)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/..)

# Source files
set(SOURCE_FILES
    ${CMAKE_CURRENT_BINARY_DIR}/../../build/ffmpeg_cmd_wrap.grpc.pb.cc
    ${CMAKE_CURRENT_BINARY_DIR}/../../build/ffmpeg_cmd_wrap.pb.cc
    ${CMAKE_CURRENT_BINARY_DIR}/../../CmdPassImpl.cc
    ${CMAKE_CURRENT_BINARY_DIR}/../../FFmpeg_wrapper_client.cc
    ${CMAKE_CURRENT_BINARY_DIR}/../../ffmpeg_pipeline_generator.cpp
    ${CMAKE_CURRENT_BINARY_DIR}/../FFmpeg_wrapper_service_UT.cc
)

# Add the executable
add_executable(FFmpegWrapperServiceTests ${SOURCE_FILES})

# Link libraries
target_link_libraries(FFmpegWrapperServiceTests
    gRPC::grpc++
    protobuf::libprotobuf
    absl::check
    absl::string_view
    absl::log
    absl::flags
    absl::flags_parse
    gtest
    gtest_main
    gmock
    gmock_main
    nlohmann_json::nlohmann_json
)

# Enable testing
enable_testing()

# Add test
add_test(NAME RunTest COMMAND FFmpegWrapperServiceTests)
add_test(NAME handleEmptyArgs COMMAND FFmpegWrapperServiceTests)
add_test(NAME handleInvalidArgs COMMAND FFmpegWrapperServiceTests)
