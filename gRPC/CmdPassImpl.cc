#include "ffmpeg_pipeline_generator.hpp"
#include <sstream>
#include "CmdPassImpl.h"

// Helper function to convert string to any type
template <typename T>
T from_string(const std::string& str) {
    std::istringstream iss(str);
    T value;
    iss >> value;
    return value;
}

// Function to convert vector of string pairs to FrameRate
static FrameRate stringPairsToFrameRate(const std::vector<std::pair<std::string, std::string>>& pairs) {
    FrameRate frameRate;

    for (const auto& pair : pairs) {
        if (pair.first == "numerator") {
            frameRate.numerator = from_string<int>(pair.second);
        } else if (pair.first == "denominator") {
            frameRate.denominator = from_string<int>(pair.second);
        }
    }

    return frameRate;
}

// Function to convert vector of string pairs to Video
static Video stringPairsToVideo(const std::vector<std::pair<std::string, std::string>>& pairs) {
    Video video;

    for (const auto& pair : pairs) {
        if (pair.first == "frame_width") {
            video.frame_width = from_string<int>(pair.second);
        } else if (pair.first == "frame_height") {
            video.frame_height = from_string<int>(pair.second);
        } else if (pair.first == "pixel_format") {
            video.pixel_format = pair.second;
        } else if (pair.first == "video_type") {
            video.video_type = pair.second;
        } else if (pair.first == "numerator" || pair.first == "denominator") {
            video.frame_rate = stringPairsToFrameRate(pairs);
        }
    }

    return video;
}

// Function to convert vector of string pairs to Audio
static Audio stringPairsToAudio(const std::vector<std::pair<std::string, std::string>>& pairs) {
    Audio audio;

    for (const auto& pair : pairs) {
        if (pair.first == "channels") {
            audio.channels = from_string<int>(pair.second);
        } else if (pair.first == "sample_rate") {
            audio.sample_rate = from_string<int>(pair.second);
        } else if (pair.first == "format") {
            audio.format = pair.second;
        } else if (pair.first == "packet_time") {
            audio.packet_time = pair.second;
        }
    }

    return audio;
}

// Function to convert vector of string pairs to File
static File stringPairsToFile(const std::vector<std::pair<std::string, std::string>>& pairs) {
    File file;

    for (const auto& pair : pairs) {
        if (pair.first == "path") {
            file.path = pair.second;
        } else if (pair.first == "filename") {
            file.filename = pair.second;
        }
    }

    return file;
}

// Function to convert vector of string pairs to ST2110
static ST2110 stringPairsToST2110(const std::vector<std::pair<std::string, std::string>>& pairs) {
    ST2110 st2110;

    for (const auto& pair : pairs) {
        if (pair.first == "network_interface") {
            st2110.network_interface = pair.second;
        } else if (pair.first == "local_ip") {
            st2110.local_ip = pair.second;
        } else if (pair.first == "remote_ip") {
            st2110.remote_ip = pair.second;
        } else if (pair.first == "transport") {
            st2110.transport = pair.second;
        } else if (pair.first == "remote_port") {
            st2110.remote_port = from_string<int>(pair.second);
        } else if (pair.first == "payload_type") {
            st2110.payload_type = from_string<int>(pair.second);
        }
    }

    return st2110;
}

// Function to convert vector of string pairs to MCM
static MCM stringPairsToMCM(const std::vector<std::pair<std::string, std::string>>& pairs) {
    MCM mcm;

    for (const auto& pair : pairs) {
        if (pair.first == "conn_type") {
            mcm.conn_type = pair.second;
        } else if (pair.first == "transport") {
            mcm.transport = pair.second;
        } else if (pair.first == "transport_pixel_format") {
            mcm.transport_pixel_format = pair.second;
        } else if (pair.first == "ip") {
            mcm.ip = pair.second;
        } else if (pair.first == "port") {
            mcm.port = from_string<int>(pair.second);
        } else if (pair.first == "urn") {
            mcm.urn = pair.second;
        }
    }

    return mcm;
}

// Function to convert vector of string pairs to Payload
static Payload stringPairsToPayload(const std::vector<std::pair<std::string, std::string>>& pairs) {
    Payload payload;

    for (const auto& pair : pairs) {
        if (pair.first == "type") {
            payload.type = static_cast<payload_type>(from_string<int>(pair.second));
        } else if (pair.first == "frame_width" || pair.first == "frame_height" ||
                   pair.first == "pixel_format" || pair.first == "video_type" ||
                   pair.first == "numerator" || pair.first == "denominator") {
            payload.video = stringPairsToVideo(pairs);

        } else if (pair.first == "channels" || pair.first == "sample_rate" ||
                   pair.first == "format" || pair.first == "packet_time") {
            payload.audio = stringPairsToAudio(pairs);
        }
    }

    return payload;
}

// Function to convert vector of string pairs to StreamType
static StreamType stringPairsToStreamType(const std::vector<std::pair<std::string, std::string>>& pairs) {
    StreamType streamType;

    for (const auto& pair : pairs) {
        if (pair.first == "type") {
            streamType.type = static_cast<stream_type>(from_string<int>(pair.second));
        } else if (pair.first == "path" || pair.first == "filename") {
            streamType.file = stringPairsToFile(pairs);
        } else if (pair.first == "network_interface" || pair.first == "local_ip" || pair.first == "remote_ip" || pair.first == "transport" || pair.first == "remote_port" || pair.first == "payload_type") {
            streamType.st2110 = stringPairsToST2110(pairs);
        } else if (pair.first == "conn_type" || pair.first == "transport" || pair.first == "transport_pixel_format" || pair.first == "ip" || pair.first == "port" || pair.first == "urn") {
            streamType.mcm = stringPairsToMCM(pairs);
        }
    }

    return streamType;
}

// Function to convert vector of string pairs to Stream
static Stream stringPairsToStream(const std::vector<std::pair<std::string, std::string>>& pairs) {
    Stream stream;

    stream.payload = stringPairsToPayload(pairs);
    stream.stream_type = stringPairsToStreamType(pairs);

    return stream;
}

// Function to convert vector of string pairs to Config
Config stringPairsToConfig(const std::vector<std::pair<std::string, std::string>>& pairs) {
    Config config;

    for (const auto& pair : pairs) {
        if (pair.first == "function") {
            config.function = pair.second;
        } else if (pair.first == "gpu_hw_acceleration") {
            config.gpu_hw_acceleration = pair.second;
        } else if (pair.first == "logging_level") {
            config.logging_level = from_string<int>(pair.second);
        } else if (pair.first == "type" || pair.first == "path" || pair.first == "filename" ||
                   pair.first == "network_interface" || pair.first == "local_ip" || pair.first == "remote_ip" ||
                   pair.first == "transport" || pair.first == "remote_port" || pair.first == "payload_type" ||
                   pair.first == "conn_type" || pair.first == "transport_pixel_format" || pair.first == "ip" ||
                   pair.first == "port" || pair.first == "urn" || pair.first == "frame_width" || pair.first == "frame_height" ||
                   pair.first == "pixel_format" || pair.first == "video_type" || pair.first == "numerator" || pair.first == "denominator" ||
                   pair.first == "channels" || pair.first == "sample_rate" || pair.first == "format" || pair.first == "packet_time") {
            Stream stream = stringPairsToStream(pairs);
            if (pair.first == "type" && pair.second == "0") {
                config.senders.push_back(stream);
            } else {
                config.receivers.push_back(stream);
            }
        }
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
        std::string pipelinie_string;
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

        if (ffmpeg_generate_pipeline(recieved_config, pipelinie_string) != 0) {
            pipelinie_string.clear();
            std::cout << "Error generating pipeline" << std::endl; //TODO : need to return as response error code
            //return 1;
        }

        ss << "ffmpeg ";
        ss << pipelinie_string;

        pipelinie_string = ss.str();

        std::array<char, 128> buffer;
        std::string result;

        FILE *pipe = popen(pipelinie_string.c_str(), "r");
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
