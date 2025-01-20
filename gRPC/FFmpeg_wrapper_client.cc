/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

#include "new_struct.hpp"
#include "FFmpeg_wrapper_client.h"
#include "build/ffmpeg_cmd_wrap.pb.h"
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

// Function to convert FrameRate to vector of string pairs
static void frameRateToStringPairs(const FrameRate& frameRate,
								   std::vector<std::pair<std::string, std::string>>& result) {

    result.push_back({"numerator", std::to_string(frameRate.numerator)});
    result.push_back({"denominator", std::to_string(frameRate.denominator)});
}

// Function to convert Video to vector of string pairs
static void videoToStringPairs(const Video& video,
							   std::vector<std::pair<std::string, std::string>>& result) {

    result.push_back({"frame_width", std::to_string(video.frame_width)});
    result.push_back({"frame_height", std::to_string(video.frame_height)});
    result.push_back({"pixel_format", video.pixel_format});
    result.push_back({"video_type", video.video_type});
    frameRateToStringPairs(video.frame_rate, result);
}

// Function to convert Audio to vector of string pairs
static void audioToStringPairs(const Audio& audio,
							   std::vector<std::pair<std::string, std::string>>& result) {

    result.push_back({"channels", std::to_string(audio.channels)});
    result.push_back({"sample_rate", std::to_string(audio.sample_rate)});
    result.push_back({"format", audio.format});
    result.push_back({"packet_time", audio.packet_time});
}

// Function to convert File to vector of string pairs
static void fileToStringPairs(const File& file,
							  std::vector<std::pair<std::string, std::string>>& result) {

    result.push_back({"path", file.path});
    result.push_back({"filename", file.filename});
}

// Function to convert ST2110 to vector of string pairs
static void st2110ToStringPairs(const ST2110& st2110,
								std::vector<std::pair<std::string, std::string>>& result) {

    result.push_back({"network_interface", st2110.network_interface});
    result.push_back({"local_ip", st2110.local_ip});
    result.push_back({"remote_ip", st2110.remote_ip});
    result.push_back({"transport", st2110.transport});
    result.push_back({"remote_port", std::to_string(st2110.remote_port)});
    result.push_back({"payload_type", std::to_string(st2110.payload_type)});
}

// Function to convert MCM to vector of string pairs
static void mcmToStringPairs(const MCM& mcm,
							 std::vector<std::pair<std::string, std::string>>& result) {

    result.push_back({"conn_type", mcm.conn_type});
    result.push_back({"transport", mcm.transport});
    result.push_back({"transport_pixel_format", mcm.transport_pixel_format});
    result.push_back({"ip", mcm.ip});
    result.push_back({"port", std::to_string(mcm.port)});
    result.push_back({"urn", mcm.urn});
}

// Function to convert Payload to vector of string pairs
static void payloadToStringPairs(const Payload& payload,
								 std::vector<std::pair<std::string, std::string>>& result) {

    result.push_back({"type", std::to_string(payload.type)});

	if (payload.type == video) {
        videoToStringPairs(payload.video, result);
    } else if (payload.type == audio) {
        audioToStringPairs(payload.audio, result);
    }
}

// Function to convert StreamType to vector of string pairs
static void streamTypeToStringPairs(const StreamType& streamType,
									std::vector<std::pair<std::string, std::string>>& result) {

    result.push_back({"type", std::to_string(streamType.type)});

    if (streamType.type == file) {
        fileToStringPairs(streamType.file, result);
    } else if (streamType.type == st2110) {
        st2110ToStringPairs(streamType.st2110, result);
    } else if (streamType.type == mcm) {
        mcmToStringPairs(streamType.mcm, result);
    }
}

// Function to convert Stream to vector of string pairs
static void streamToStringPairs(const Stream& stream,
								std::vector<std::pair<std::string, std::string>>& result) {

    payloadToStringPairs(stream.payload, result);
    streamTypeToStringPairs(stream.stream_type, result);
}

// Function to convert Config to vector of string pairs
std::vector<std::pair<std::string, std::string>> commitConfigs(const Config& config) {
    std::vector<std::pair<std::string, std::string>> result;

    result.push_back({"function", config.function});
    result.push_back({"gpu_hw_acceleration", config.gpu_hw_acceleration});
    result.push_back({"logging_level", std::to_string(config.logging_level)});

	for (const auto& sender : config.senders) {
        streamToStringPairs(sender, result);
    }

    for (const auto& receiver : config.receivers) {
        streamToStringPairs(receiver, result);
    }

	return result;
}
