/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

#include "ffmpeg_pipeline_generator.hpp"

// Return 0 if payloads match, 1 otherwise:
int compare_payloads(Payload &p1, Payload &p2) {
    if (p1.type != p2.type) {
        std::cout << "Error: payloads types do not match" << std::endl;
        return 1;
    }

    if (p1.type == payload_type::video && (
        p1.video.frame_width != p2.video.frame_width ||
        p1.video.frame_height != p2.video.frame_height ||
        p1.video.frame_rate.numerator != p2.video.frame_rate.numerator ||
        p1.video.frame_rate.denominator != p2.video.frame_rate.denominator ||
        p1.video.pixel_format != p2.video.pixel_format ||
        p1.video.video_type != p2.video.video_type))
    {
        std::cout << "Error: video payloads do not match" << std::endl;
        return 1;
    }

    if (p1.type == payload_type::audio && (
        p1.audio.sample_rate != p2.audio.sample_rate ||
        p1.audio.format != p2.audio.format ||
        p1.audio.packet_time != p2.audio.packet_time))
    {
        std::cout << "Error: audio payloads do not match" << std::endl;
        return 1;
    }

    return 0;
}

int ffmpeg_append_payload(Payload &p, std::string &pipeline_string) {
    switch (p.type) {
    case video:
        pipeline_string += " -video_size " + std::to_string(p.video.frame_width) + "x" + std::to_string(p.video.frame_height);
        pipeline_string += " -pix_fmt " + p.video.pixel_format;
        pipeline_string += " -r " + std::to_string(p.video.frame_rate.numerator) + "/" + std::to_string(p.video.frame_rate.denominator);
        if(p.video.video_type == "rawvideo")
            pipeline_string += " -f rawvideo";
        else{
            pipeline_string += " -c:v " + p.video.video_type;
        }
        break;
    case audio:
        std::cout << "Error: audio not supported yet" << std::endl;
        return 1;
    default:
        std::cout << "Error: unknown payload type" << std::endl;
        break;
    }
    return 0;
}

int ffmpeg_append_st2110_transport(std::string &transport, std::string &pipeline_string) {
    if (transport == "st2110-20") {
        pipeline_string += " -f mtl_st20p";
    }
    else if (transport == "st2110-22") {
        pipeline_string += " -f mtl_st22p";
    }
    else if (transport == "st2110-30") {
        pipeline_string += " -f mtl_st30p";
    }
    else {
        std::cout << "Error: transport " << transport << "not supported yet" << std::endl;
        return 1;
    }
    return 0;
}

int ffmpeg_append_stream_type(StreamType &s, bool is_rx, int idx, std::string &pipeline_string) {
    switch (s.type) {
    case file:
    {
        pipeline_string += " ";
        if (is_rx) {
            pipeline_string += "-i ";
        }

        pipeline_string += s.file.path;
        if (!s.file.path.empty() && s.file.path.back() != '/') {
            pipeline_string += '/';
        }
        pipeline_string += s.file.filename;

        break;
    }
    case st2110:
        pipeline_string += " -p_port " + s.st2110.network_interface;
        pipeline_string += " -p_sip " + s.st2110.local_ip;
        pipeline_string += " -udp_port " + std::to_string(s.st2110.remote_port);
        pipeline_string += " -payload_type " + std::to_string(s.st2110.payload_type);
        if(is_rx) {
            pipeline_string += " -p_rx_ip " + s.st2110.remote_ip;
        }
        else {
            pipeline_string += " -p_tx_ip " + s.st2110.remote_ip;
        }

        if(ffmpeg_append_st2110_transport(s.st2110.transport, pipeline_string) != 0){
            pipeline_string.clear();
            std::cout << "Error appending st2110 transport" << std::endl;
            return 1;
        }
        if(is_rx) {
            pipeline_string += " -i \"" + std::to_string(idx) + "\"";
        }
        else {
            pipeline_string += " -";
        }


        break;
    case mcm:
        std::cout << "Error: mcm not supported yet" << std::endl;
        return 1;
    default:
        break;
    }

    return 0;
}

int ffmpeg_combine_rx_tx(Stream &rx, Stream &tx, int idx, std::string &pipeline_string){
    if(compare_payloads(rx.payload, tx.payload) != 0){
        std::cout << "Error: payloads do not match" << std::endl;
        return 1;
    }

    if(rx.payload.video.video_type != "rawvideo"){
        std::cout << "Error: video type not supported" << rx.payload.video.video_type << std::endl;
        return 1;
    }


    if(ffmpeg_append_payload(rx.payload,  pipeline_string) != 0){
        pipeline_string.clear();
        std::cout << "Error appending rx payload" << std::endl;
        return 1;
    }

    if(ffmpeg_append_stream_type(rx.stream_type, true/*is_rx*/, idx, pipeline_string) != 0){
        pipeline_string.clear();
        std::cout << "Error appending rx stream type" << std::endl;
        return 1;
    }

    if(ffmpeg_append_stream_type(tx.stream_type, false/*is_rx*/, idx, pipeline_string) != 0){
        pipeline_string.clear();
        std::cout << "Error appending tx stream type" << std::endl;
        return 1;
    }

    return 0;
}

int ffmpeg_generate_pipeline(Config &config, std::string &pipeline_string) {
    if (config.logging_level > 0) {
        // TODO: add more logging level
        pipeline_string += " -v debug";
    }

    if (config.gpu_hw_acceleration == "none") {
        pipeline_string += " -y";
    }
    else if (config.gpu_hw_acceleration == "intel") {
        pipeline_string += " -y -qsv_device /dev/dri/renderD128 -hwaccel qsv";
    }
    else if (config.gpu_hw_acceleration == "nvidia") {
        pipeline_string += " -y -hwaccel cuda -hwaccel_output_format cuda";
    }
    else {
        std::cout << "Unsupported GPU acceleration" << config.gpu_hw_acceleration << std::endl;
        return 1;
    }

    if (config.function == "tx" || config.function == "rx") {
        if ((config.receivers.size() == 0) || (config.receivers.size() != config.senders.size())) {
            std::cout << "Error: function " << config.function << " requires equal number of receivers and senders, greater than 0" << std::endl;
            return 1;
        }
        else {
            for (int i = 0; i < config.receivers.size(); i++) {
                if (ffmpeg_combine_rx_tx(config.receivers[i], config.senders[i], i, pipeline_string) != 0) {
                    pipeline_string.clear();
                    std::cout << "Error combining rx and tx" << std::endl;
                    return 1;
                }
            }
        }
    }
    return 0;
}
