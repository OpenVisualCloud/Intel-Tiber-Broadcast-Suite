/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

#include "config_params.hpp"
#include "FFmpeg_wrapper_client.h"
#include "build/ffmpeg_cmd_wrap.pb.h"
#include "config_serialize_deserialize.hpp"
#include <sstream>
#include <utility>
#include <string>
#include <iostream>

CmdPassClient::CmdPassClient(std::string interface, std::string port) : pending_requests_(0) {
    std::cout << "------ [start] initiate channel --------" << std::endl;

    std::stringstream ss;
    std::string channel_config;

    ss << interface << ":" << port;

    channel_config = ss.str();

    std::shared_ptr<Channel> channel = grpc::CreateChannel(channel_config, grpc::InsecureChannelCredentials());

    stub_ = CmdPass::NewStub(channel);

    std::cout << "------ [done] initiate channel --------" << std::endl;

    // Start the thread to process the completion queue
    cq_thread_ = std::thread(&CmdPassClient::AsyncCompleteRpc, this);
}

CmdPassClient::~CmdPassClient() {
    cq_.Shutdown();
    cq_thread_.join();
}

void CmdPassClient::FFmpegCmdExec(std::vector<std::pair<std::string, std::string>>& cmd_pairs) {

    ReqCmds req_obj;
    CmdMsg *cmd_msg;
    
    for (const auto& cmd_pair : cmd_pairs) {
        cmd_msg = req_obj.add_obj();
        cmd_msg->set_cmd_key(cmd_pair.first);
        cmd_msg->set_cmd_val(cmd_pair.second);
    }

    auto* call = new AsyncClientCall;

    ++pending_requests_;

    // Initiate the asynchronous RPC call
    call->response_reader = stub_->PrepareAsyncFFmpegCmdExec(&call->context, req_obj, &cq_);
    call->response_reader->StartCall();
    call->response_reader->Finish(&call->response, &call->status, (void*)call);
}

void CmdPassClient::AsyncCompleteRpc() {
    void* got_tag;
    bool ok = false;

    while (cq_.Next(&got_tag, &ok)) {
        auto* call = static_cast<AsyncClientCall*>(got_tag);

        if (ok) {
            if (call->status.ok()) {
                std::cout << "FFmpeg command executed successfully: " << call->status.error_code() << std::endl;
            }
            else {
                std::cout << "FFmpeg command execution failed:" << std::endl;
                std::cout << "Status = " << call->status.error_code() << std::endl;
                std::cout << "Message = " << call->status.error_message() << std::endl;
                std::cout << "Details = " << call->status.error_details() << std::endl;
            }
        }
        else {
            std::cout << "RPC failed" << std::endl;
        }

        delete call;

        if (--pending_requests_ == 0) {
            all_tasks_completed = true;
            all_tasks_completed.notify_one();
        }
    }
}

void CmdPassClient::WaitForAllRequests() {
    all_tasks_completed.wait(false);
}


// Function to convert Config to vector of string pairs
std::vector<std::pair<std::string, std::string>> commitConfigs(const Config& config) {
    std::vector<std::pair<std::string, std::string>> result;

    std::string json_str;
    if (serialize_config_json(config, json_str) == 0) {
        result.push_back({"json", json_str});
    }
    else {
        std::cout << "Error serializing Config to json, trying previos solution" << std::endl;
    };

    return result;
}
