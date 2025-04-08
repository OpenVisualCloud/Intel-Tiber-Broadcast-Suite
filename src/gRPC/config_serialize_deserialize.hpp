#include "config_params.hpp"
#include "nlohmann/json.hpp"
#include <iostream>
#include <string>

#ifndef CONFIG_SERIALIZE_DESERIALIZE_H
#define CONFIG_SERIALIZE_DESERIALIZE_H

NLOHMANN_DEFINE_TYPE_NON_INTRUSIVE(FrameRate, numerator, denominator)
NLOHMANN_DEFINE_TYPE_NON_INTRUSIVE(Video, frame_width, frame_height, frame_rate,
                                   pixel_format, video_type)
NLOHMANN_DEFINE_TYPE_NON_INTRUSIVE(Audio, channels, sample_rate, format,
                                   packet_time)
NLOHMANN_DEFINE_TYPE_NON_INTRUSIVE(File, path, filename)
NLOHMANN_DEFINE_TYPE_NON_INTRUSIVE(ST2110, network_interface, local_ip,
                                   remote_ip, transport, remote_port,
                                   payload_type)
NLOHMANN_DEFINE_TYPE_NON_INTRUSIVE(MCM, conn_type, transport,
                                   transport_pixel_format, ip, port, urn)
NLOHMANN_DEFINE_TYPE_NON_INTRUSIVE(Payload, type, video, audio)
NLOHMANN_DEFINE_TYPE_NON_INTRUSIVE(StreamType, type, file, st2110, mcm)
NLOHMANN_DEFINE_TYPE_NON_INTRUSIVE(Stream, payload, stream_type)
NLOHMANN_DEFINE_TYPE_NON_INTRUSIVE(Config, senders, receivers, function,
                                   multiviewer_columns, gpu_hw_acceleration,
                                   gpu_hw_acceleration_device, logging_level,
                                   stream_loop)

// TODO: move serialize_config_json and deserialize_config_json to a separate
// file

static int serialize_config_json(const Config &input_config,
                                 std::string &output_string) {
  try {
    Config new_config = input_config;
    new_config.receivers[0].payload.type = payload_type::video;
    new_config.senders[0].payload.type = payload_type::video;

    nlohmann::json config_json = new_config;
    // Dump json to string
    output_string = config_json.dump();
    std::cout << "JSON output string: " << std::endl
              << output_string << std::endl;
  } catch (const nlohmann::json::parse_error &e) {
    std::cout << "JSON parse error: " << e.what() << std::endl;
    return 1;
  } catch (const nlohmann::json::type_error &e) {
    std::cout << "JSON type error: " << e.what() << std::endl;
    return 1;
  } catch (const std::exception &e) {
    std::cout << "Exception: " << e.what() << std::endl;
    return 1;
  }
  return 0;
}

static int deserialize_config_json(Config &output_config,
                                   const std::string &input_string) {
  try {
    nlohmann::json config_json = nlohmann::json::parse(input_string);
    std::cout << "JSON input string: " << std::endl
              << input_string << std::endl;
    // Deserialize from json to Config
    output_config = config_json.get<Config>();
  } catch (const nlohmann::json::parse_error &e) {
    std::cout << "JSON parse error: " << e.what() << std::endl;
    return 1;
  } catch (const nlohmann::json::type_error &e) {
    std::cout << "JSON type error: " << e.what() << std::endl;
    return 1;
  } catch (const std::exception &e) {
    std::cout << "Exception: " << e.what() << std::endl;
    return 1;
  }
  return 0;
}

#endif // CONFIG_SERIALIZE_DESERIALIZE_H
