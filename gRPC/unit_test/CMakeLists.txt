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
  URL https://github.com/google/googletest/archive/refs/tags/release-1.11.0.zip
)
# For Windows: Prevent overriding the parent project's compiler/linker settings
set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
FetchContent_MakeAvailable(googletest)

# Include directories
include_directories(${CMAKE_CURRENT_SOURCE_DIR})
include_directories(${googletest_SOURCE_DIR}/googletest/include)
include_directories(${googletest_SOURCE_DIR}/googlemock/include)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../build)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/..)

# Source files
set(SOURCE_FILES
    ${CMAKE_CURRENT_BINARY_DIR}/../FFmpeg_wrapper_service_UT.cc
    ${CMAKE_CURRENT_BINARY_DIR}/../../build/ffmpeg_cmd_wrap.grpc.pb.cc
    ${CMAKE_CURRENT_BINARY_DIR}/../../build/ffmpeg_cmd_wrap.pb.cc
    ${CMAKE_CURRENT_BINARY_DIR}/../../FFmpeg_wrapper_service.cc
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
)

# Enable testing
enable_testing()

# Add test
add_test(NAME FFmpegWrapperServiceTests COMMAND FFmpegWrapperServiceTests)
