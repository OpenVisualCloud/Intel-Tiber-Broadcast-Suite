/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

#ifndef INTEL_CONFIG_PARSER_H
#define INTEL_CONFIG_PARSER_H

#include "config_params.hpp"
#include <cpprest/json.h>

class ConfigManager {
public:
  ConfigManager() = default;
  int parse_json_file(const std::string &file_path);
  void print_config() const;
  const Config &get_config() const;
  std::pair<int, int> get_framerate(const Stream &stream) const;

protected:
  Config config;
  Video parse_video(const web::json::value &video_data) const;
  Audio parse_audio(const web::json::value &audio_data) const;
  StreamType parse_stream_type(const web::json::value &stream_type_data) const;
  Stream parse_stream(const web::json::value &stream_data) const;
};

#endif // INTEL_CONFIG_PARSER_H