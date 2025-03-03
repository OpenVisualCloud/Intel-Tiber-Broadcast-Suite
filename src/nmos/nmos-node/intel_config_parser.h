#ifndef INTEL_CONFIG_PARSER_H
#define INTEL_CONFIG_PARSER_H

#include <cpprest/json.h>
#include "config_params.hpp"

class ConfigManager {
public:
    ConfigManager() = default;
    void parse_json_file(const std::string& file_path);
    void print_config() const;
    const Config& get_config() const;
    std::pair<int, int> get_framerate(const Stream& stream) const;

private:
    Config config;
    Video parse_video(const web::json::value& video_data) const;
    Audio parse_audio(const web::json::value& audio_data) const;
    StreamType parse_stream_type(const web::json::value& stream_type_data) const;
    Stream parse_stream(const web::json::value& stream_data) const;
};

#endif // INTEL_CONFIG_PARSER_H