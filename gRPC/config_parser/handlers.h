#ifndef HANDLERS_H
#define HANDLERS_H

#include <iostream>
#include <string>
#include <unordered_map>
#include <functional>
#include <stdexcept>
#include <sstream>
#include <nlohmann/json.hpp>

// Alias for convenience
using json = nlohmann::json;

// Handler function type
using HandlerFunc = std::function<std::string(const std::string&, const json&)>;

// Handlers map
extern std::unordered_map<std::string, HandlerFunc> handlers;

// Function declarations
std::string handleKeyValue(const std::string& key, const json& value);
std::string handleFrameRate(const std::string& key, const json& value);
std::string handleFormatKeyValue(const std::string& key, const json& value);
std::string handleCodecKeyValue(const std::string& key, const json& value);

#endif // HANDLERS_H
