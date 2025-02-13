#include "intel_config_parser.h"
#include <cpprest/filestream.h>
#include <iostream>
#include <fstream>

const Config& ConfigManager::get_config() const {
    return config;
}

std::pair<int, int> ConfigManager::get_framerate(const Stream& stream) const {
    return {stream.payload.video.frame_rate.numerator, stream.payload.video.frame_rate.denominator};
}

void ConfigManager::parse_json_file(const std::string& file_path) {
    try {
        // Open the file stream
        std::ifstream file(file_path);
        if (!file.is_open()) {
            throw std::runtime_error("Could not open file");
        }

        // Read the file into a string
        std::string json_str((std::istreambuf_iterator<char>(file)), std::istreambuf_iterator<char>());

        // Parse the JSON string
        web::json::value json_value = web::json::value::parse(json_str);

        // Fill the Config struct
        config.logging_level = json_value.at(U("logging_level")).as_integer();
        config.function = json_value.at(U("function")).as_string();
        config.gpu_hw_acceleration = json_value.at(U("gpu_hw_acceleration")).as_string();

        for (const auto& sender : json_value.at(U("sender")).as_array()) {
            std::cout<<"Sender: "<<std::endl;
            config.senders.push_back(parse_stream(sender));
        }

        for (const auto& receiver : json_value.at(U("receiver")).as_array()) {
            config.receivers.push_back(parse_stream(receiver));
        }

    } catch (const std::exception& e) {
        std::cerr << "Error parsing JSON file: " << e.what() << std::endl;
    }
}

void ConfigManager::print_config() const {
    std::cout << "Logging Level: " << config.logging_level << std::endl;
    std::cout << "Function: " << config.function << std::endl;
    std::cout << "GPU HW Acceleration: " << config.gpu_hw_acceleration << std::endl;

    for (const auto& sender : config.senders) {
        std::cout << "Sender Video Frame Width: " << sender.payload.video.frame_width << std::endl;
        std::cout << "Sender Video Frame Height: " << sender.payload.video.frame_height << std::endl;
        std::cout << "Sender Video Frame Rate: " << sender.payload.video.frame_rate.numerator << "/" << sender.payload.video.frame_rate.denominator << std::endl;
        std::cout << "Sender Video Pixel Format: " << sender.payload.video.pixel_format << std::endl;
        std::cout << "Sender Video Type: " << sender.payload.video.video_type << std::endl;
        if (sender.stream_type.type == stream_type::file) {
            std::cout<<"Sender Stream type: File"<<std::endl;
            std::cout << "|_ File Path: " << sender.stream_type.file.path << std::endl;
            std::cout << "|_ File Name: " << sender.stream_type.file.filename << std::endl;
        } else if (sender.stream_type.type == stream_type::st2110) {
            std::cout<<"Sender Stream type: st2110"<<std::endl;
            std::cout << "|_ Transport: " << sender.stream_type.st2110.transport << std::endl;
            std::cout << "|_ Payload Type: " << sender.stream_type.st2110.payload_type << std::endl;
        } else if (sender.stream_type.type == stream_type::mcm) {
            std::cout<<"Sender Stream type: mcm"<<std::endl;
            std::cout << "|_ Connection Type: " << sender.stream_type.mcm.conn_type << std::endl;
            std::cout << "|_ Transport: " << sender.stream_type.mcm.transport << std::endl;
            std::cout << "|_ URN: " << sender.stream_type.mcm.urn << std::endl;
            std::cout << "|_ Transport Pixel Format: " << sender.stream_type.mcm.transport_pixel_format << std::endl;
        }
    }

    for (const auto& receiver : config.receivers) {
        std::cout << "Receiver Video Frame Width: " << receiver.payload.video.frame_width << std::endl;
        std::cout << "Receiver Video Frame Height: " << receiver.payload.video.frame_height << std::endl;
        std::cout << "Receiver Video Frame Rate: " << receiver.payload.video.frame_rate.numerator << "/" << receiver.payload.video.frame_rate.denominator << std::endl;
        std::cout << "Receiver Video Pixel Format: " << receiver.payload.video.pixel_format << std::endl;
        std::cout << "Receiver Video Type: " << receiver.payload.video.video_type << std::endl;
        if (receiver.stream_type.type == stream_type::file) {
            std::cout<<"Receiver Stream type: File"<<std::endl;
            std::cout << "|_ File Path: " << receiver.stream_type.file.path << std::endl;
            std::cout << "|_ File Name: " << receiver.stream_type.file.filename << std::endl;
        } else if (receiver.stream_type.type == stream_type::st2110) {
            std::cout<<"Receiver Stream type: st2110"<<std::endl;
            std::cout << "|_ Transport: " << receiver.stream_type.st2110.transport << std::endl;
            std::cout << "|_ Payload Type: " << receiver.stream_type.st2110.payload_type << std::endl;
        } else if (receiver.stream_type.type == stream_type::mcm) {
            std::cout<<"Receiver Stream type: mcm"<<std::endl;
            std::cout << "|_ Connection Type: " << receiver.stream_type.mcm.conn_type << std::endl;
            std::cout << "|_ Transport: " << receiver.stream_type.mcm.transport << std::endl;
            std::cout << "|_ URN: " << receiver.stream_type.mcm.urn << std::endl;
            std::cout << "|_ Transport Pixel Format: " << receiver.stream_type.mcm.transport_pixel_format << std::endl;
        }
    }
}

Video ConfigManager::parse_video(const web::json::value& video_data) const {
    Video video;
    video.frame_width = video_data.at(U("frame_width")).as_integer();
    video.frame_height = video_data.at(U("frame_height")).as_integer();
    video.frame_rate.numerator = video_data.at(U("frame_rate")).at(U("numerator")).as_integer();
    video.frame_rate.denominator = video_data.at(U("frame_rate")).at(U("denominator")).as_integer();
    video.pixel_format = video_data.at(U("pixel_format")).as_string();
    video.video_type = video_data.at(U("video_type")).as_string();
    return video;
}

Audio ConfigManager::parse_audio(const web::json::value& audio_data) const {
    Audio audio;
    audio.channels = audio_data.at(U("channels")).as_integer();
    audio.sample_rate = audio_data.at(U("sampleRate")).as_integer();
    audio.format = audio_data.at(U("format")).as_string();
    audio.packet_time = audio_data.at(U("packetTime")).as_string();
    return audio;
}

StreamType ConfigManager::parse_stream_type(const web::json::value& stream_type_data) const {
    StreamType stream_type;
    if (stream_type_data.has_field(U("file"))) {
        stream_type.type = stream_type::file;
        stream_type.file.path = stream_type_data.at(U("file")).at(U("path")).as_string();
        stream_type.file.filename = stream_type_data.at(U("file")).at(U("filename")).as_string();
    } else if (stream_type_data.has_field(U("st2110"))) {
        stream_type.type = stream_type::st2110;
        stream_type.st2110.transport = stream_type_data.at(U("st2110")).at(U("transport")).as_string();
        stream_type.st2110.payload_type = stream_type_data.at(U("st2110")).at(U("payloadType")).as_integer();
    } else if (stream_type_data.has_field(U("mcm"))) {
        stream_type.type = stream_type::mcm;
        stream_type.mcm.conn_type = stream_type_data.at(U("mcm")).at(U("conn_type")).as_string();
        stream_type.mcm.transport = stream_type_data.at(U("mcm")).at(U("transport")).as_string();
        stream_type.mcm.urn = stream_type_data.at(U("mcm")).at(U("urn")).as_string();
        stream_type.mcm.transport_pixel_format = stream_type_data.at(U("mcm")).at(U("transportPixelFormat")).as_string();
    }
    return stream_type;
}

Stream ConfigManager::parse_stream(const web::json::value& stream_data) const {
    Stream stream;
    stream.payload.video = parse_video(stream_data.at(U("stream_payload")).at(U("video")));
    stream.payload.audio = parse_audio(stream_data.at(U("stream_payload")).at(U("audio")));
    stream.stream_type = parse_stream_type(stream_data.at(U("stream_type")));
    return stream;
}