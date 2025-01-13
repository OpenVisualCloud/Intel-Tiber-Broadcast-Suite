#include "ffmpeg_config_processor.h"
#include "handlers.h"
#include <fstream>
#include <sstream>
#include <iostream>

// Function to process appParams section and use handlers to print key-value pairs and collect formatted strings
std::string processAppParams(const json& data) {
    std::vector<std::string> commandChunks;
    std::ostringstream commandStream;
    std::string commandString;
    std::string formattedString;

    for (const auto& [key, value] : data.items()) {
        if (handlers.find(key) != handlers.end()) {
            try {
                formattedString = handlers[key](key, value);
                commandChunks.push_back(formattedString);
            } catch (const std::exception& e) {
                throw std::runtime_error("error processing key " + key + ": " + e.what());
            }
        } else {
            try {
                formattedString = handleKeyValue(key, value);
                commandChunks.push_back(formattedString);
            } catch (const std::exception& e) {
                throw std::runtime_error("error processing key " + key + ": " + e.what());
            }
        }
    }

    for (const auto& str : commandChunks) {
        commandStream << str << " ";
    }

    commandString = commandStream.str();
    std::cout << "Command String: " << commandString << std::endl;
    
    return commandString;
}

// Function to process the configuration file and return the command string
std::string processConfigFile(const std::string& configFile) {
    // Open the JSON configuration file
    std::ifstream file(configFile);
    if (!file.is_open()) {
        throw std::runtime_error("Could not open the configuration file: " + configFile);
    }

    // Read the JSON configuration file
    std::stringstream buffer;
    buffer << file.rdbuf();
    std::string jsonString = buffer.str();

    // Parse the JSON configuration file
    json config;
    try {
        config = json::parse(jsonString);
    } catch (const json::parse_error& e) {
        throw std::runtime_error("JSON parse error: " + std::string(e.what()));
    }

    // Extract and process the appParams section
    if (config.contains("ffmpegPipelineDefinition")) {
        const auto& ffmpegPipelineDefinition = config["ffmpegPipelineDefinition"];
        if (ffmpegPipelineDefinition.contains("appParams")) {
            const auto& appParams = ffmpegPipelineDefinition["appParams"];
            // Construct and return the command string
            return processAppParams(appParams);
        } else {
            throw std::runtime_error("appParams section not found in the configuration file.");
        }
    } else {
        throw std::runtime_error("ffmpegPipelineDefinition section not found in the configuration file.");
    }
}
