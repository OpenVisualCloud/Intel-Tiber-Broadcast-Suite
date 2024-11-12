/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

#include "FFmpeg_wrapper_client.h"
#include <utility>
#include <vector>

int main(int argc, char* argv[]) {

    if(argc != 3) {
        std::cout << "client sample app only takes two argument 1) interface 2) port" << std::endl;
        return 1;
    }

    std::vector<std::pair<std::string, std::string>> vec1 = {};
    std::vector<std::pair<std::string, std::string>> vec2 = {{"key3", "val3"}, {"key4", "val4"}};
    std::vector<std::pair<std::string, std::string>> vec3 = {{"key5", "val5"}, {"key6", "val6"}};

    std::string interface = "localhost";
    std::string port = "50051";

    CmdPassClient obj(interface, port);

    // Send multiple asynchronous requests
    obj.FFmpegCmdExec(vec1);
    obj.FFmpegCmdExec(vec2);
    obj.FFmpegCmdExec(vec3);

    // Wait for all asynchronous operations to complete
    obj.WaitForAllRequests();

    return 0;
}
