/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

#include "FFmpeg_wrapper_client.h"
#include <iostream>
#include <utility>
#include <vector>

int main(int argc, char* argv[]) {

    if (argc != 5) {
        std::cout << "client sample app only takes two arguments: 1) interface 2) port 3)source_ip 4) destination_port" << std::endl;
        return 1;
    }

    std::string interface = argv[1];//"localhost";
    std::string port = argv[2];//"50051";
    std::string source_ip = argv[3];
    std::string destination_port = argv[4];

    // Populate the connection_info vector with the provided values
    std::vector<std::pair<std::string, std::string>> connection_info = {{"source_ip", source_ip}, {"destination_port", destination_port}};

    CmdPassClient obj(interface, port);

    // Send multiple asynchronous requests
    obj.FFmpegCmdExec(connection_info);

    // Wait for all asynchronous operations to complete
    obj.WaitForAllRequests();

    return 0;
}