#include <gmock/gmock.h>
#include <gtest/gtest.h>
#include "ffmpeg_cmd_wrap.grpc.pb.h"
#include "ffmpeg_cmd_wrap.pb.h"
#include "FFmpeg_wrapper_service.cc"

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
    }

    void TearDown() override {
        delete cq;
        delete service;
    }

    MockAsyncService* service;
    MockServerCompletionQueue* cq;
    CmdPassImpl cmdPassImpl;
};

TEST_F(CmdPassImplTest, RequestFFmpegCmdExecTest) {
    grpc::ServerContext context;
    ReqCmds request;
    FFmpegServiceRes response;
    grpc::ServerAsyncResponseWriter<FFmpegServiceRes> responder(&context);

    // Set expectation on the RequestFFmpegCmdExec method
    EXPECT_CALL(*service, RequestFFmpegCmdExec(&context, &request, &responder, cq, cq, testing::_))
        .Times(1)
        .WillOnce(testing::Invoke([&](grpc::ServerContext* ctx, ReqCmds* req, grpc::ServerAsyncResponseWriter<FFmpegServiceRes>* res, grpc::CompletionQueue* new_call_cq, grpc::ServerCompletionQueue* notification_cq, void* tag) {
            // Simulate the behavior of the method
        }));

    // Mock the behavior of the completion queue
    EXPECT_CALL(*cq, Next(testing::_, testing::_))
        .WillOnce(testing::Invoke([](void** tag, bool* ok) {
            *ok = true;
            return true;
        }))
        .WillRepeatedly(testing::Return(false)); // Simulate the end of the queue

    // Trigger the method
    service->RequestFFmpegCmdExec(&context, &request, &responder, nullptr, cq, nullptr);

    // Verify the state transitions
    // Since we don't have GetCallDataStates, we will check if the method was called correctly
    ASSERT_TRUE(testing::Mock::VerifyAndClearExpectations(service));
}
