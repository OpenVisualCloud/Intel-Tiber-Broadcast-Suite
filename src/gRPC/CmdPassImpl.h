#ifndef CMD_PASS_IMPL_H
#define CMD_PASS_IMPL_H

#include <cerrno>
#include <csignal>
#include <cstdio>
#include <iostream>
#include <sstream>
#include <thread>

#include "ffmpeg_cmd_wrap.grpc.pb.h"
#include "ffmpeg_cmd_wrap.pb.h"
#include <grpc/grpc.h>
#include <grpc/support/log.h>
#include <grpcpp/security/server_credentials.h>
#include <grpcpp/server.h>
#include <grpcpp/server_builder.h>
#include <grpcpp/server_context.h>

using grpc::Server;
using grpc::ServerAsyncResponseWriter;
using grpc::ServerBuilder;
using grpc::ServerCompletionQueue;
using grpc::ServerContext;
using grpc::Status;

#define FFMPEG_INVALID_COMMAND_STATUS "1"
#define FFMPEG_INVALID_COMMAND_MSG std::string("Failed to execute ffmpeg command : No commands provided")

#define FFMPEG_APP_EXEC_FAIL_STATUS "-1"
#define FFMPEG_APP_EXEC_FAIL_MSG std::string("Failed to execute ffmpeg pipeline (popen) : ") + std::strerror(errno)

#define FFMPEG_COMMAND_FAIL_STATUS "2"
#define FFMPEG_COMMAND_FAIL_MSG std::string("FFmpeg command failed : ")

#define FFMPEG_EXEC_OK_STATUS "0"
#define FFMPEG_EXEC_OK_MSG std::string("FFmpeg command : ") + std::strerror(errno)

class CmdPassImpl final {
public:
    std::atomic<bool> stop;

    void Run(std::string server_address);
    void Shutdown();

private:
    class CallData {
    public:
        CallData(CmdPass::AsyncService *service, grpc::ServerCompletionQueue *cq);
        void Proceed();

    private:
        CmdPass::AsyncService *service_;
        grpc::ServerCompletionQueue *cq_;
        grpc::ServerContext ctx_;
        ReqCmds request_;
        FFmpegServiceRes response_;
        grpc::ServerAsyncResponseWriter<FFmpegServiceRes> responder_;
        enum CallStatus { CREATE, PROCESS, FINISH };
        CallStatus status_;
    };

    void HandleRpcs();

    std::unique_ptr<ServerCompletionQueue> cq_;
    CmdPass::AsyncService service_;
    std::unique_ptr<Server> server_;
};

#endif // CMD_PASS_IMPL_H
