#ifndef FFMPEG_CONFIG_PROCESSOR_H
#define FFMPEG_CONFIG_PROCESSOR_H

#include <string>
#include <vector>
#include <nlohmann/json.hpp>

// Alias for convenience
using json = nlohmann::json;

// Function declarations
std::string processAppParams(const json& data);
std::string processConfigFile(const std::string& configFile);

#endif // FFMPEG_CONFIG_PROCESSOR_H
