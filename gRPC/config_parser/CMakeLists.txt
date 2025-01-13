cmake_minimum_required(VERSION 3.14)
project(MyFFmpegApp LANGUAGES CXX)

# Set C++ standard to C++20
set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Include the JSON for Modern C++ library
include(FetchContent)

FetchContent_Declare(
  json
  GIT_REPOSITORY https://github.com/nlohmann/json.git
  GIT_TAG v3.10.5
)

FetchContent_MakeAvailable(json)

# Add the source files for the static library
set(SOURCES
    handlers.cpp
    ffmpeg_config_processor.cpp
)

# Add the header files for the static library
set(HEADERS
    handlers.h
    ffmpeg_config_processor.h
)

# Create the static library
add_library(ffmpegconfig STATIC ${SOURCES} ${HEADERS})

# Link the JSON library to the static library
target_link_libraries(ffmpegconfig PRIVATE nlohmann_json::nlohmann_json)

# Specify include directories
target_include_directories(ffmpegconfig PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})
