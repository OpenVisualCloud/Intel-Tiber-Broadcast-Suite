/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 */

#include <gtest/gtest.h>
#include "intel_config_parser.h"

class ConfigManagerTest : public ConfigManager {
public:
    Stream parse_stream_public(const web::json::value& stream_data) const {
        return parse_stream(stream_data);
    }

    int parse_json_string(const std::string& str){
        try {
            web::json::value json_value = web::json::value::parse(str);

            // Fill the Config struct
            config.logging_level = json_value.at(U("logging_level")).as_integer();
            if (json_value.has_field(U("stream_loop"))) {
              config.stream_loop = json_value.at(U("stream_loop")).as_integer();
            } else {
              config.stream_loop = 0;
            }
            config.function = json_value.at(U("function")).as_string();
            config.gpu_hw_acceleration = json_value.at(U("gpu_hw_acceleration")).as_string();
            if (config.function == "multiviewer") {
              config.multiviewer_columns = json_value.at(U("multiviewer_columns")).as_integer();
            }
            else{
                config.multiviewer_columns = 3;
            }
  
            if (config.gpu_hw_acceleration != "none") {
                config.gpu_hw_acceleration_device = json_value.at(U("gpu_hw_acceleration_device")).as_string();
            }
            else{
                config.gpu_hw_acceleration_device = "";
            }

            for (const auto& sender : json_value.at(U("sender")).as_array()) {
                config.senders.push_back(parse_stream(sender));
            }

            for (const auto& receiver : json_value.at(U("receiver")).as_array()) {
                config.receivers.push_back(parse_stream(receiver));
            }

        } catch (const std::exception& e) {
            std::cerr << "Error parsing JSON file: " << e.what() << std::endl;
            return -1;
        }
        return 0;
    }
};

const std::string json_str = R"json(
    {
      "logging_level": 0,
      "http_port": 95,
      "label": "intel-broadcast-suite",
      "activate_senders": false,
      "device_tags": {
        "pipeline": ["rx"]
      },
      "function": "rx",
      "gpu_hw_acceleration": "none",
      "color_sampling": "YCbCr-4:2:2",
      "domain": "local",
      "ffmpeg_grpc_server_address": "localhost",
      "ffmpeg_grpc_server_port": "50052",
      "receiver_payload_type": 112,
      "frame_rate": { "numerator": 60000, "denominator": 1001 },
      "sender": [{
        "stream_payload": {
          "video": {
            "frame_width": 1920,
            "frame_height": 1080,
            "frame_rate": { "numerator": 60, "denominator": 1 },
            "pixel_format": "yuv422p10le",
            "video_type": "rawvideo"
          },
          "audio": {
            "channels": 2,
            "sampleRate": 48000,
            "format": "pcm_s24be",
            "packetTime": "1ms"
          }
        },
        "stream_type": {
          "file": {
            "path": "/root/recv",
            "filename": "1920x1080p10le_2.yuv"
          }
        }
      }],
      "receiver": [{
        "stream_payload": {
          "video": {
            "frame_width": 1920,
            "frame_height": 1080,
            "frame_rate": { "numerator": 60, "denominator": 1 },
            "pixel_format": "yuv422p10le",
            "video_type": "rawvideo"
          },
          "audio": {
            "channels": 2,
            "sampleRate": 48000,
            "format": "pcm_s24be",
            "packetTime": "1ms"
          }
        },
        "stream_type": {
          "st2110": {
            "transport": "st2110-20",
            "payloadType": 112
          }
        }
      }]
    }
    )json";


TEST(ConfigManagerTest, ParseStreamValidData) {
    ConfigManagerTest c_mgr;

    c_mgr.parse_json_string(json_str);

    Config conf = c_mgr.get_config();

    // Validate the parsed senders
    EXPECT_EQ(conf.logging_level, 0);                     // "logging_level": 0,
    EXPECT_EQ(conf.function, "rx");                       // "function": "rx",
    EXPECT_EQ(conf.gpu_hw_acceleration, "none");          // "gpu_hw_acceleration": "none",
    EXPECT_EQ(conf.gpu_hw_acceleration_device, "");        // "gpu_hw_acceleration_device": "",
    //EXPECT_EQ(color_sampling, "YCbCr-4:2:2");           // "color_sampling": "YCbCr-4:2:2",
    //EXPECT_EQ(domain, "local");                         // "domain": "local",
    //EXPECT_EQ(ffmpeg_grpc_server_address, "localhost"); // "ffmpeg_grpc_server_address": "localhost",
    //EXPECT_EQ(ffmpeg_grpc_server_port, 50052);          // "ffmpeg_grpc_server_port": "50052",
    //EXPECT_EQ(receiver_payload_type, 112);              // "receiver_payload_type": 112,
    //EXPECT_EQ(frame_rate.numerator, 60000);             // "frame_rate": { "numerator": 60000,
    //EXPECT_EQ(frame_rate.denominator, 1001);            // "denominator": 1001 },

    // Validate the parsed senders
    ASSERT_EQ(conf.senders.size(), 1);
    EXPECT_EQ(conf.senders[0].payload.video.frame_width, 1920);
    EXPECT_EQ(conf.senders[0].payload.video.frame_height, 1080);
    EXPECT_EQ(conf.senders[0].payload.video.frame_rate.numerator, 60);
    EXPECT_EQ(conf.senders[0].payload.video.frame_rate.denominator, 1);
    EXPECT_EQ(c_mgr.get_framerate(conf.senders[0]).first, 60);
    EXPECT_EQ(c_mgr.get_framerate(conf.senders[0]).second, 1);
    EXPECT_EQ(conf.senders[0].payload.video.pixel_format, "yuv422p10le");
    EXPECT_EQ(conf.senders[0].payload.video.video_type, "rawvideo");
    EXPECT_EQ(conf.senders[0].payload.audio.channels, 2);
    EXPECT_EQ(conf.senders[0].payload.audio.sample_rate, 48000);
    EXPECT_EQ(conf.senders[0].payload.audio.format, "pcm_s24be");
    EXPECT_EQ(conf.senders[0].payload.audio.packet_time, "1ms");
    EXPECT_EQ(conf.senders[0].stream_type.file.path, "/root/recv");
    EXPECT_EQ(conf.senders[0].stream_type.file.filename, "1920x1080p10le_2.yuv");

    // Validate the parsed receivers
    ASSERT_EQ(conf.receivers.size(), 1);
    EXPECT_EQ(conf.receivers[0].payload.video.frame_width, 1920);
    EXPECT_EQ(conf.receivers[0].payload.video.frame_height, 1080);
    EXPECT_EQ(conf.receivers[0].payload.video.frame_rate.numerator, 60);
    EXPECT_EQ(conf.receivers[0].payload.video.frame_rate.denominator, 1);
    EXPECT_EQ(conf.receivers[0].payload.video.pixel_format, "yuv422p10le");
    EXPECT_EQ(conf.receivers[0].payload.video.video_type, "rawvideo");
    EXPECT_EQ(conf.receivers[0].payload.audio.channels, 2);
    EXPECT_EQ(conf.receivers[0].payload.audio.sample_rate, 48000);
    EXPECT_EQ(conf.receivers[0].payload.audio.format, "pcm_s24be");
    EXPECT_EQ(conf.receivers[0].payload.audio.packet_time, "1ms");
    EXPECT_EQ(conf.receivers[0].stream_type.st2110.transport, "st2110-20");
    EXPECT_EQ(conf.receivers[0].stream_type.st2110.payload_type, 112);
}


TEST(ConfigManagerTest, ParseStreamGpuData) {
  ConfigManagerTest c_mgr;

  const std::string gpu_json_str = R"json(
  {
      "logging_level": 0,
      "function": "rx",
      "gpu_hw_acceleration": "intel",
      "gpu_hw_acceleration_device": "/dev/dri/renderD128",
      "sender": [],
      "receiver": []
  }
  )json";

  c_mgr.parse_json_string(gpu_json_str);

  Config conf = c_mgr.get_config();

  EXPECT_EQ(conf.gpu_hw_acceleration, "intel");
  EXPECT_EQ(conf.gpu_hw_acceleration_device, "/dev/dri/renderD128");
}

TEST(ConfigManagerTest, ParseStreamNoneGpuData) {
  ConfigManagerTest c_mgr;

  const std::string gpu_json_str = R"json(
  {
      "logging_level": 0,
      "function": "rx",
      "gpu_hw_acceleration": "none",
      "gpu_hw_acceleration_device": "",
      "sender": [],
      "receiver": []
  }
  )json";

  c_mgr.parse_json_string(gpu_json_str);

  Config conf = c_mgr.get_config();

  EXPECT_EQ(conf.gpu_hw_acceleration, "none");
  EXPECT_EQ(conf.gpu_hw_acceleration_device, "");
}

TEST(ConfigManagerTest, ParseStreamMultiviewer) {
  ConfigManagerTest c_mgr;

  const std::string gpu_json_str = R"json(
  {
      "logging_level": 0,
      "function": "multiviewer",
      "multiviewer_columns": 4,
      "gpu_hw_acceleration": "none",
      "gpu_hw_acceleration_device": "",
      "sender": [],
      "receiver": []
  }
  )json";

  c_mgr.parse_json_string(gpu_json_str);

  Config conf = c_mgr.get_config();

  EXPECT_EQ(conf.function, "multiviewer");
  EXPECT_EQ(conf.multiviewer_columns, 4);
}

TEST(ConfigManagerTest, ParseStreamMultiviewerDefault) {
  ConfigManagerTest c_mgr;

  const std::string gpu_json_str = R"json(
  {
      "logging_level": 0,
      "function": "tx",
      "gpu_hw_acceleration": "none",
      "gpu_hw_acceleration_device": "",
      "sender": [],
      "receiver": []
  }
  )json";

  c_mgr.parse_json_string(gpu_json_str);

  Config conf = c_mgr.get_config();

  EXPECT_EQ(conf.function, "tx");
  EXPECT_EQ(conf.multiviewer_columns, 3);
}

TEST(ConfigManagerTest, ParseStreamValidStreamLoop) {
  ConfigManagerTest c;

  const std::string stream_json = R"json(
  {
      "logging_level": 0,
      "function": "tx",
      "gpu_hw_acceleration": "none",
      "gpu_hw_acceleration_device": "",
      "stream_loop": -1,
      "sender": [],
      "receiver": []
  }
  )json";
  c.parse_json_string(stream_json);

  Config conf = c.get_config();

  EXPECT_EQ(conf.stream_loop, -1);
}

TEST(ConfigManagerTest, ParseStreamDefaultStreamLoop) {
  ConfigManagerTest c;

  const std::string stream_json = R"json(
  {
      "logging_level": 0,
      "function": "tx",
      "gpu_hw_acceleration": "none",
      "gpu_hw_acceleration_device": "",
      "sender": [],
      "receiver": []
  }
  )json";
  c.parse_json_string(stream_json);

  Config conf = c.get_config();

  EXPECT_EQ(conf.stream_loop, 0);
}

TEST(ConfigManagerTest, ParseStreamInvalidData) {
    ConfigManagerTest c;

    // Invalid JSON data
    const std::string invalid_json_str = R"json(
    {
        "logging_level": "invalid",
        "http_port": 95
    }
    )json";

    EXPECT_EQ(c.parse_json_string(invalid_json_str), -1);
}

TEST(ConfigManagerTest, ParseStreamInvalidFile) {
    ConfigManager c;

    EXPECT_EQ(c.parse_json_file(""), -1);
    EXPECT_EQ(c.parse_json_file("\\"), -1);
}
