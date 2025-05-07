/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

#include <string>
#include <vector>

#ifndef CONFIG_PARAMS_H
#define CONFIG_PARAMS_H

struct FrameRate {
    int numerator;
    int denominator;
};

struct Video {
    int frame_width;
    int frame_height;
    FrameRate frame_rate;
    std::string pixel_format;
    std::string video_type;
    std::string preset; //optional
    std::string profile; //optional
    // in case of rawvideo then ffmpeg param = "-f rawwideo"
    // otherwise -c:v <video_type> <preset> <profile> e.g. -c:v x264 -preset veryfast -profile main

};

 // Audio struct is a placeholder for future implementation
struct Audio {
    int channels;
    int sample_rate;
    std::string format;
    std::string packet_time;
};

struct File {
    std::string path;
    std::string filename;
};

struct ST2110 {
    std::string network_interface; //VFIO port address 0000:00:00.0; ffmpeg param name: -p_port
    std::string local_ip; // ffmpeg param name: -p_sip
    std::string remote_ip; // ffmpeg param name: -p_rx_ip / -p_tx_ip
    std::string transport;
    int remote_port; // ffmpeg param name: -udp_port
    int payload_type;
    int queues_cnt; // ffmpeg param name: -rx_queues / -tx_queues ; 0 mean use default ffmpeg plugin value
};

struct MCM {
    std::string conn_type;
    std::string transport;
    std::string transport_pixel_format;
    std::string ip;
    int port;
    std::string urn;
};

enum payload_type {
    video = 0,
    audio
 };

struct Payload {
    payload_type type;
    Video video;
    Audio audio;
};

enum stream_type {
    file = 0,
    st2110,
    mcm
 };

struct StreamType {
    stream_type type;
    File file;
    ST2110 st2110;
    MCM mcm;
};

struct Stream {
    Payload payload;
    StreamType stream_type;
};

struct Config {
    std::vector<Stream> senders;
    std::vector<Stream> receivers;

    std::string function; //multiviewer, upscale, replay, recorder, jpegxs, rx, tx
    int multiviewer_columns; //number of streams in a row

    std::string gpu_hw_acceleration; //intel, nvidia, none
    std::string gpu_hw_acceleration_device; // /dev/dri/renderD128

    int stream_loop; // number of times to loop the input stream

    int logging_level;
};

#endif // CONFIG_PARAMS_H
