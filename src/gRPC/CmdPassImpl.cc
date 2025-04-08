
#include "ffmpeg_pipeline_generator.hpp"
#include "config_serialize_deserialize.hpp"
#include <sstream>
#include "CmdPassImpl.h"

// Function to convert vector of string pairs to Config
static Config stringPairsToConfig(const std::vector<std::pair<std::string, std::string>>& pairs) {
    Config config;

    if (deserialize_config_json(config, pairs.front().second) != 0) {
        std::cout << "Error deserializing Config from json, trying previous solution" << std::endl;
    }

    return config;
}

void CmdPassImpl::Run(std::string server_address) {
    ServerBuilder builder;
    stop = false;

    builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
    builder.RegisterService(&service_);

    cq_ = builder.AddCompletionQueue();

    server_ = builder.BuildAndStart();

    std::cout << "[*] Server run and listening on " << server_address << std::endl;

    std::jthread th([&]() {
        stop.wait(false);
        std::cout << "[*] Shutting down the server" << std::endl;

        server_->Shutdown();
        cq_->Shutdown();

        void *ignored_tag;
        bool ignored_ok;
        while (cq_->Next(&ignored_tag, &ignored_ok)) {
        }
    });

    HandleRpcs();
}

void CmdPassImpl::Shutdown() {
    stop = true;
    stop.notify_one();
}

CmdPassImpl::CallData::CallData(CmdPass::AsyncService *service, grpc::ServerCompletionQueue *cq)
    : service_(service), cq_(cq), responder_(&ctx_), status_(CREATE) {
    Proceed();
}

void CmdPassImpl::CallData::Proceed() {
    if (status_ == CREATE) {
        status_ = PROCESS;
        service_->RequestFFmpegCmdExec(&ctx_, &request_, &responder_, cq_, cq_, this);
    } else if (status_ == PROCESS) { /* FFmpeg command processing starts */
        new CallData(service_, cq_);

        std::stringstream ss;
        std::string pipeline_string;
        std::vector<std::pair<std::string, std::string>> committed_config;

        if (request_.obj().empty()) {
            responder_.Finish(response_,
                              Status(grpc::INVALID_ARGUMENT,
                                     FFMPEG_INVALID_COMMAND_STATUS,
                                     FFMPEG_INVALID_COMMAND_MSG),
                              this);

            status_ = FINISH;
            return;
        }

        for (const auto &cmd : request_.obj()) {
            committed_config.push_back(std::make_pair(cmd.cmd_key(), cmd.cmd_val()));
        }

        Config recieved_config = stringPairsToConfig(committed_config);

        if (ffmpeg_generate_pipeline(recieved_config, pipeline_string) != 0) {
            pipeline_string.clear();
            std::cout << "Error generating pipeline" << std::endl; //TODO : need to return as response error code
            //return 1;
        }

        ss << "ffmpeg ";
        ss << pipeline_string;

        pipeline_string = ss.str();

        std::array<char, 128> buffer;
        std::string result;

        FILE *pipe = popen(pipeline_string.c_str(), "r");
        if (!pipe) { /* FFmpeg pipeline/execution failed i.e memory allocation */
            responder_.Finish(response_,
                              Status(grpc::INTERNAL, FFMPEG_APP_EXEC_FAIL_STATUS,
                                     FFMPEG_APP_EXEC_FAIL_MSG),
                              this);

            status_ = FINISH;
            return;
        }

        /* FFmpeg output */
        while (fgets(buffer.data(), buffer.size(), pipe) != nullptr) {
            result += buffer.data();
        }

        std::cout << result << std::endl;
        /* FFmpeg output */

        /* FFmpeg pipeline/execution ended */
        if (pclose(pipe) != 0) {
            /* FFmpeg app fail : returns the exit status of FFmpeg app itself */
            responder_.Finish(response_,
                              Status(grpc::UNKNOWN, FFMPEG_COMMAND_FAIL_STATUS,
                                     FFMPEG_COMMAND_FAIL_MSG),
                              this);

            status_ = FINISH;
            return;
            /* FFmpeg app fail : returns the exit status of FFmpeg app itself */
        }
        /* FFmpeg pipeline/execution ended */

        /* FFmpeg successfully executed the commands */
        responder_.Finish(
            response_,
            Status(grpc::OK, FFMPEG_EXEC_OK_STATUS, FFMPEG_EXEC_OK_MSG), this);

        status_ = FINISH;
        return;
    } /* FFmpeg command processing ends */
    else {
        // Once in the FINISH state, deallocate ourselves (CallData).
        delete this;
    }
}

void CmdPassImpl::HandleRpcs() {
    new CallData(&service_, cq_.get());

    void *tag;
    bool ok;
    while (true) {
        cq_->Next(&tag, &ok);

        if (ok) {
            /*
             * calling cq_->Shutdown will render "ok" to be false
             * hence we should terminate the main loop
             */
            static_cast<CallData *>(tag)->Proceed();
        } else
            break;
    }
}
