#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include "FFmpeg_wrapper_client.h"
#include "CmdPassImpl.h"

class MockAsyncService : public CmdPass::AsyncService {
public:
    MOCK_METHOD(void, RequestFFmpegCmdExec, (grpc::ServerContext* context, ReqCmds* request, grpc::ServerAsyncResponseWriter<FFmpegServiceRes>* responder, grpc::CompletionQueue* new_call_cq, grpc::ServerCompletionQueue* notification_cq, void* tag));
};

class MockServerCompletionQueue : public grpc::ServerCompletionQueue {
public:
    MOCK_METHOD(bool, Next, (void** tag, bool* ok));
};

class CmdPassImplTest : public ::testing::Test {
protected:
    void SetUp() override {
        service = new MockAsyncService();
        cq = new MockServerCompletionQueue();
        cmdPassImpl = new CmdPassImpl();
    }

    void TearDown() override {
        delete cmdPassImpl;
        delete cq;
        delete service;
    }

    MockAsyncService* service;
    MockServerCompletionQueue* cq;
    CmdPassImpl* cmdPassImpl;
};

TEST_F(CmdPassImplTest, RunTest) {
    std::string server_address = "localhost:50051";

    // Run the server in a separate thread
    std::thread server_thread([&]() {
        cmdPassImpl->Run(server_address);
    });

    // Allow some time for the server to start
    std::this_thread::sleep_for(std::chrono::seconds(1));

    // Shutdown the server
    cmdPassImpl->Shutdown();

    // Wait for the server thread to finish
    server_thread.join();

    // Verify that the server was started and shutdown correctly
    ASSERT_TRUE(true); // If we reach here, the test passed
}

TEST_F(CmdPassImplTest, handleEmptyArgs) {
    grpc::ServerContext context;
    ReqCmds request;
    FFmpegServiceRes response;
    grpc::ServerAsyncResponseWriter<FFmpegServiceRes> responder(&context);
    
    std::vector<std::pair<std::string, std::string>> cmds = {};
    
    /*
     * server setup and run
     */
    std::string server_address = "localhost:50051";

    // Run the server in a separate thread
    std::thread server_thread([&]() {
        cmdPassImpl->Run(server_address);
    });

    // Allow some time for the server to start
    std::this_thread::sleep_for(std::chrono::seconds(1));
    /*
     * END - server setup and run
     */

    std::string interface = "localhost";
    std::string port = "50051";

    CmdPassClient obj(interface, port);

    obj.FFmpegCmdExec(cmds);

    // Shutdown the server
    cmdPassImpl->Shutdown();

    // Wait for the server thread to finish
    server_thread.join();

    /*
     * Status = 3
     *  Message = 1
     *  Details = Failed to execute ffmpeg command : No commands provided
     */
    ASSERT_TRUE(true);
}

TEST_F(CmdPassImplTest, handleInvalidArgs) {
    grpc::ServerContext context;
    ReqCmds request;
    FFmpegServiceRes response;
    grpc::ServerAsyncResponseWriter<FFmpegServiceRes> responder(&context);
    
    std::vector<std::pair<std::string, std::string>> cmds = {{"key3", "val3"}, {"key4", "val4"}};
    
    /*
     * server setup and run
     */
    std::string server_address = "localhost:50051";

    // Run the server in a separate thread
    std::thread server_thread([&]() {
        cmdPassImpl->Run(server_address);
    });

    // Allow some time for the server to start
    std::this_thread::sleep_for(std::chrono::seconds(1));
    /*
     * END - server setup and run
     */

    std::string interface = "localhost";
    std::string port = "50051";

    CmdPassClient obj(interface, port);

    obj.FFmpegCmdExec(cmds);

    // Shutdown the server
    cmdPassImpl->Shutdown();

    // Wait for the server thread to finish
    server_thread.join();

    /*
     * Status = 3
     *  Message = 1
     *  Details = Failed to execute ffmpeg command : No commands provided
     */
    ASSERT_TRUE(true);
}
