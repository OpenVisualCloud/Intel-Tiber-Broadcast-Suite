/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

#ifndef _CMD_PASS_CLIENT_H_
#define _CMD_PASS_CLIENT_H_

#include <iostream>
#include <memory>
#include <ostream>
#include <random>
#include <string>
#include <thread>
#include <utility>
#include <vector>
#include <atomic>
#include "build/ffmpeg_cmd_wrap.pb.h"

#include <grpc/grpc.h>
#include <grpcpp/channel.h>
#include <grpcpp/client_context.h>
#include <grpcpp/create_channel.h>
#include <grpcpp/security/credentials.h>
#include "build/ffmpeg_cmd_wrap.grpc.pb.h"

using grpc::Channel;
using grpc::ClientContext;
using grpc::CompletionQueue;
using grpc::Status;

std::vector<std::pair<std::string, std::string>> commitConfigs(const Config& config);

class CmdPassClient {
public:
    CmdPassClient(std::string interface, std::string port);
    ~CmdPassClient();

    void FFmpegCmdExec(std::vector<std::pair<std::string, std::string>>& cmd_pairs);
    void WaitForAllRequests();

private:
    void AsyncCompleteRpc();

	// Create a new call object
    struct AsyncClientCall {
        FFmpegServiceRes response;
        ClientContext context;
        Status status;
        std::unique_ptr<grpc::ClientAsyncResponseReader<FFmpegServiceRes>> response_reader;
    };

    std::unique_ptr<CmdPass::Stub> stub_;
    CompletionQueue cq_;
    std::thread cq_thread_;
    std::atomic<int> pending_requests_{0};
    std::atomic<bool> all_tasks_completed{false};
};

#endif
