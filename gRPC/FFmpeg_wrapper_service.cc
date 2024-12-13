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

  void Run(std::string server_address) {
    ServerBuilder builder;
    stop = false;

    builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
    builder.RegisterService(&service_);

    cq_ = builder.AddCompletionQueue();

    server_ = builder.BuildAndStart();

    std::cout << "[*] Server run and listening on " << server_address
              << std::endl;

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

  void Shutdown() {
    stop = true;
    stop.notify_one();
  }

private:
  class CallData {
  public:
    CallData(CmdPass::AsyncService *service, grpc::ServerCompletionQueue *cq)
        : service_(service), cq_(cq), responder_(&ctx_), status_(CREATE) {
      Proceed();
    }

    void Proceed() {
      if (status_ == CREATE) {
        status_ = PROCESS;
        service_->RequestFFmpegCmdExec(&ctx_, &request_, &responder_, cq_, cq_,
                                       this);
      } else if (status_ == PROCESS) { /* FFmpeg command processing starts */
        new CallData(service_, cq_);

        std::stringstream ss;
        ss << "ffmpeg";

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
          ss << " -" << cmd.cmd_key() << " " << cmd.cmd_val();
        }

        std::string ffmpeg_full_cmd = ss.str();

        std::array<char, 128> buffer;
        std::string result;

        FILE *pipe = popen(ffmpeg_full_cmd.c_str(), "r");
        if (!pipe) { /* FFmpeg pipeline/execution failed i.e memory allocation
                      */
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

  void HandleRpcs() {
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

  std::unique_ptr<ServerCompletionQueue> cq_;
  CmdPass::AsyncService service_;
  std::unique_ptr<Server> server_;
};

CmdPassImpl manager;

void wrapperSignalHandler(int signal) {
  std::cout << "[*] wrapperSignalHandler called SIG : " << signal << std::endl;
  manager.Shutdown();
}

int main(int argc, char *argv[]) {
  if (argc != 3) {
    std::cout
        << "[*] FFmpeg wrapper service takes 2 arguments: interface and port"
        << std::endl;
    return 1;
  }

  std::stringstream ss;
  ss << argv[1] << ":" << argv[2];

  std::signal(SIGTERM, wrapperSignalHandler);
  std::signal(SIGINT, wrapperSignalHandler);

  manager.Run(ss.str());

  std::cout << "[*] Service exited gracefully" << std::endl;

  return 0;
}
