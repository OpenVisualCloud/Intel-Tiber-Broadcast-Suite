/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

#include "ffmpeg_pipeline_generator.hpp"

#include <cmath>
#include <cstdint>

// Return 0 if payloads match,
// Return -1 if payload types are incompatible,
// Return 1 otherwise
int compare_payloads(Payload &p1, Payload &p2) {
    if (p1.type != p2.type) {
        std::cout << "Error: payloads types do not match" << std::endl;
        return -1;
    }

    if (p1.type == payload_type::video && (
        p1.video.frame_width != p2.video.frame_width ||
        p1.video.frame_height != p2.video.frame_height ||
        p1.video.frame_rate.numerator != p2.video.frame_rate.numerator ||
        p1.video.frame_rate.denominator != p2.video.frame_rate.denominator ||
        p1.video.pixel_format != p2.video.pixel_format ||
        p1.video.video_type != p2.video.video_type))
    {
        return 1;
    }

    if (p1.type == payload_type::audio && (
        p1.audio.sample_rate != p2.audio.sample_rate ||
        p1.audio.format != p2.audio.format ||
        p1.audio.packet_time != p2.audio.packet_time))
    {
        return 1;
    }

    return 0;
}

int ffmpeg_append_payload(Payload &p, std::string &pipeline_string) {
    switch (p.type) {
    case video:
        if(p.video.video_type == "rawvideo") {
        pipeline_string += " -video_size " + std::to_string(p.video.frame_width) + "x" + std::to_string(p.video.frame_height);
        pipeline_string += " -pix_fmt " + p.video.pixel_format;
        pipeline_string += " -r " + std::to_string(p.video.frame_rate.numerator) + "/" + std::to_string(p.video.frame_rate.denominator);
        pipeline_string += " -f rawvideo";
        }
        else if (!p.video.video_type.empty()) {
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

int ffmpeg_append_stream_convert(Payload &rx, Payload &tx, std::string &pipeline_string) {
    if (rx.type != tx.type) {
        std::cout << "Error: payloads types do not match" << std::endl;
        return 1;
    }

    if (rx.type != payload_type::video) {
        std::cout << "Error: audio payload conversion not supported yet" << std::endl;
        return 1;
    }

    if (rx.video.pixel_format != tx.video.pixel_format) {
        pipeline_string += " -pix_fmt " + tx.video.pixel_format;
    }
    if (rx.video.frame_width != tx.video.frame_width || rx.video.frame_height != tx.video.frame_height) {
        pipeline_string += " -vf scale=" + std::to_string(tx.video.frame_width) + "x" + std::to_string(tx.video.frame_height);
    }
    if (rx.video.frame_rate.numerator != tx.video.frame_rate.numerator || rx.video.frame_rate.denominator != tx.video.frame_rate.denominator) {
        pipeline_string += " -r " + std::to_string(tx.video.frame_rate.numerator) + "/" + std::to_string(tx.video.frame_rate.denominator);
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
        std::cout << "Error: transport " << transport << " not supported yet" << std::endl;
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
    if (ffmpeg_append_payload(rx.payload,  pipeline_string) != 0) {
        pipeline_string.clear();
        std::cout << "Error appending rx payload" << std::endl;
        return 1;
    }

    if (ffmpeg_append_stream_type(rx.stream_type, true/*is_rx*/, idx, pipeline_string) != 0) {
        pipeline_string.clear();
        std::cout << "Error appending rx stream type" << std::endl;
        return 1;
    }

    if (compare_payloads(rx.payload, tx.payload) > 0) {
        if (ffmpeg_append_stream_convert(rx.payload, tx.payload, pipeline_string)) {
            pipeline_string.clear();
            std::cout << "Error appending tx payload" << std::endl;
            return 1;
        }
    }

    if (ffmpeg_append_stream_type(tx.stream_type, false/*is_rx*/, idx, pipeline_string) != 0) {
        pipeline_string.clear();
        std::cout << "Error appending tx stream type" << std::endl;
        return 1;
    }

    return 0;
}

int ffmpeg_append_multiviewer_input(Stream &s, int idx, std::string &pipeline_string){
    if(s.payload.type != payload_type::video){
        std::cout << "Error: multiviewer requires video payload all receivers" << std::endl;
        return 1;
    }

    if(ffmpeg_append_payload(s.payload,  pipeline_string) != 0){
        pipeline_string.clear();
        std::cout << "Error appending rx payload" << std::endl;
        return 1;
    }

    if(ffmpeg_append_stream_type(s.stream_type, true/*is_rx*/, idx, pipeline_string) != 0){
        pipeline_string.clear();
        std::cout << "Error appending rx stream type" << std::endl;
        return 1;
    }

    return 0;
}

int ffmpeg_append_multiviewer_process(std::vector<Stream> &receivers, Video &output_video, uint32_t columns, uint32_t intel_gpu, std::string &pipeline_string) {
    pipeline_string += " -filter_complex \"";

    const uint rows = std::ceil((float)receivers.size() / columns);

    const uint single_screen_height = output_video.frame_height / rows;
    const uint single_screen_width = output_video.frame_width / columns;

    for (int i = 0; i < receivers.size(); i++) {
        pipeline_string += "[" + std::to_string(i) + ":v]";
        if(intel_gpu){
            pipeline_string += "hwupload,scale_qsv=";
        }
        else{
            pipeline_string += "scale=";
        }
        pipeline_string += std::to_string(single_screen_width) + ":" + std::to_string(single_screen_height) + "[out" + std::to_string(i) + "];";
    }
    for (int i = 0; i < receivers.size(); i++) {
        pipeline_string += "[out" + std::to_string(i) + "]";
    }
    //xstack_qsv=inputs=Z:layout=
    if(intel_gpu){
        pipeline_string += "xstack_qsv";
    }
    else{
        pipeline_string += "xstack";
    }
    pipeline_string += "=inputs=" + std::to_string(receivers.size()) + ":layout=";
    for (int i = 0; i < receivers.size(); i++) {
        const uint x_cord = (i % columns) * single_screen_width;
        const uint y_cord = (i / columns) * single_screen_height;
        pipeline_string += std::to_string(x_cord) + "_" + std::to_string(y_cord);
        if(i != receivers.size() - 1) {
            pipeline_string += "|";
        }
    }
    pipeline_string += ",format=y210le,format=yuv422p10le\"";

    return 0;
}

int ffmpeg_append_split_process(std::vector<Stream> &senders, uint32_t intel_gpu, std::string &pipeline_string) {
    pipeline_string += " -filter_complex \"split=" + std::to_string(senders.size());
    for (int i = 0; i < senders.size(); i++) {
        pipeline_string += "[in" + std::to_string(i) + "]";
    }
    pipeline_string += ";";
    for (int i = 0; i < senders.size(); i++) {
        if(intel_gpu){
            pipeline_string += "[in" + std::to_string(i) + "]hwupload,scale_qsv=";
        }
        else{
            pipeline_string += "[in" + std::to_string(i) + "]scale=";
        }
        pipeline_string += std::to_string(senders[i].payload.video.frame_width) + ":" + std::to_string(senders[i].payload.video.frame_height) + "[out" + std::to_string(i) + "];";
    }
    pipeline_string += "\"";

    for (int i = 0; i < senders.size(); i++) {
        pipeline_string += " -map \"[out" + std::to_string(i) + "]\"";
        if(senders[i].payload.video.video_type != "rawvideo") {
            pipeline_string += " -c:v " + senders[i].payload.video.video_type;
        }
        ffmpeg_append_stream_type(senders[i].stream_type, false /*is_rx*/, i, pipeline_string);
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
        std::cout << "Unsupported GPU acceleration " << config.gpu_hw_acceleration << std::endl;
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
        return 0;
    }
    else if (config.function == "multiviewer") {
        if (config.senders.size() != 1) {
            std::cout << "Error: multiviewer requires exactly 1 sender" << std::endl;
            return 1;
        }
        if(config.senders[0].payload.type != payload_type::video) {
            std::cout << "Error: multiviewer requires video payload" << std::endl;
            return 1;
        }
        if (config.receivers.size() < 2) {
            std::cout << "Error: multiviewer requires at least 2 receivers" << std::endl;
            return 1;
        }

        for (int i = 0; i < config.receivers.size(); i++) {
            if(ffmpeg_append_multiviewer_input(config.receivers[i], i, pipeline_string) != 0) {
                pipeline_string.clear();
                std::cout << "Error appending multiviewer input" << std::endl;
                return 1;
            }
        }
        if(ffmpeg_append_multiviewer_process(config.receivers, config.senders[0].payload.video, 3/*columns*/, config.gpu_hw_acceleration == "intel", pipeline_string) != 0) {
            pipeline_string.clear();
            std::cout << "Error appending multiviewer process" << std::endl;
            return 1;
        }
        if (ffmpeg_append_stream_type(config.senders[0].stream_type, false /*is_rx*/, 0, pipeline_string) != 0) {
            pipeline_string.clear();
            std::cout << "Error appending multiviewer tx stream" << std::endl;
            return 1;
        }

        return 0;
    }
    else if (config.function == "recorder") {
        if (config.receivers.size() != 1) {
            std::cout << "Error: recorder requires exactly 1 receiver" << std::endl;
            return 1;
        }
        if(config.receivers[0].payload.type != payload_type::video) {
            std::cout << "Error: recorder requires video payload" << std::endl;
            return 1;
        }
        if (config.senders.size() < 2) {
            std::cout << "Error: recorder requires at least 2 senders" << std::endl;
        }

        if(ffmpeg_append_payload(config.receivers[0].payload,  pipeline_string) != 0){
            pipeline_string.clear();
            std::cout << "Error appending recorder rx payload" << std::endl;
            return 1;
        }
        if (ffmpeg_append_stream_type(config.receivers[0].stream_type, true /*is_rx*/, 0, pipeline_string) != 0) {
            pipeline_string.clear();
            std::cout << "Error appending recorder rx stream" << std::endl;
            return 1;
        }
        if(ffmpeg_append_split_process(config.senders, config.gpu_hw_acceleration == "intel", pipeline_string) != 0) {
            pipeline_string.clear();
            std::cout << "Error appending recorder process" << std::endl;
            return 1;
        }

    }
    else {
        std::cout << "Error: function " << config.function << " not supported yet" << std::endl;
        return 1;
    }
    return 0;
}
