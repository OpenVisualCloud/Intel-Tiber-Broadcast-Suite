/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

#include "ffmpeg_pipeline_generator.hpp"
#include "FFmpeg_wrapper_client.h"
#include "CmdPassImpl.h"
#include <iostream>
#include <utility>
#include <vector>

void fill_conf_sender(Config &config) {
    config.function = "tx";
    config.gpu_hw_acceleration = "none";
    config.logging_level = 0;

    Payload p;
    p.type = payload_type::video;
    p.video.frame_width = 1920;
    p.video.frame_height = 1080;
    p.video.frame_rate = {30, 1};
    p.video.pixel_format = "yuv422p10le";
    p.video.video_type = "rawvideo";

    {
        Stream s;

        s.payload = p;
        s.stream_type.type = stream_type::file;
        s.stream_type.file.path = "/home/tszumski";
        s.stream_type.file.filename = "1920x1080p10le_1.yuv";
        config.receivers.push_back(s);

        s.stream_type.file.filename = "1920x1080p10le_2.yuv";
        config.receivers.push_back(s);
    }

    {
        Stream s;

        s.payload = p;
        s.stream_type.type = stream_type::st2110;
        s.stream_type.st2110.network_interface = "0000:4b:11.0";
        s.stream_type.st2110.local_ip = "192.168.2.1";
        s.stream_type.st2110.remote_ip = "192.168.2.2";
        s.stream_type.st2110.transport = "st2110-20";
        s.stream_type.st2110.remote_port = 20000;
        s.stream_type.st2110.payload_type = 112;
        config.senders.push_back(s);

        s.stream_type.st2110.remote_port = 20001;
        config.senders.push_back(s);
    }
}

int main(int argc, char* argv[]) {
        Config conf;
        fill_conf_sender(conf);

    auto aaa = commitConfigs(conf);

    Config recv = stringPairsToConfig(aaa);
    
    // if (argc != 5) {
    //     std::cout << "client sample app requires the following arguments: 1) interface, 2) port, 3) source_ip, 4) destination_port" << std::endl;
    //     return 1;
    // }

    // std::string interface = argv[1];//"localhost";
    // std::string port = argv[2];//"50051";
    // std::string source_ip = argv[3];
    // std::string destination_port = argv[4];

    // // Populate the connection_info vector with the provided values
    // std::vector<std::pair<std::string, std::string>> connection_info = {{"ip_addr", source_ip}, {"port", destination_port}};

    // CmdPassClient obj(interface, port);

    // // Send multiple asynchronous requests
    // obj.FFmpegCmdExec(connection_info);

    // // Wait for all asynchronous operations to complete
    // obj.WaitForAllRequests();

    return 0;
}
