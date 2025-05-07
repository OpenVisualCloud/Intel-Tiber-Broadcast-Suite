#include <gtest/gtest.h>
#include "ffmpeg_pipeline_generator.hpp"
#include "config_serialize_deserialize.hpp"

static Payload get_video_payload_common(){
    Payload p;
    p.type = payload_type::video;
    p.video.frame_width = 1920;
    p.video.frame_height = 1080;
    p.video.frame_rate = {30, 1};
    p.video.pixel_format = "yuv422p10le";
    p.video.video_type = "rawvideo";
    return p;
}

void fill_conf_sender(Config &config) {
    config.function = "tx";
    config.gpu_hw_acceleration = "none";
    config.logging_level = 0;
    config.stream_loop = -1;

    Payload p = get_video_payload_common();
    {
        Stream s;

        s.payload = p;
        s.stream_type.type = stream_type::file;
        s.stream_type.file.path = "/home/test";
        s.stream_type.file.filename = "1920x1080p10le_1.yuv";
        config.receivers.push_back(s);

        s.stream_type.file.path = "/home/test/";
        s.stream_type.file.filename = "1920x1080p10le_2.yuv";
        config.receivers.push_back(s);
    }

    {
        Stream s;

        s.payload = p;
        s.stream_type.type = stream_type::st2110;
        s.stream_type.st2110.network_interface = "0000:4b:11.0";
        s.stream_type.st2110.local_ip = "192.168.2.1";
        s.stream_type.st2110.remote_ip = "192.168.2.2";
        s.stream_type.st2110.transport = "st2110-20";
        s.stream_type.st2110.remote_port = 20000;
        s.stream_type.st2110.payload_type = 112;
        s.stream_type.st2110.queues_cnt = 0;
        config.senders.push_back(s);

        s.stream_type.st2110.remote_port = 20001;
        config.senders.push_back(s);
    }
}

void fill_conf_receiver(Config &config) {
    config.function = "rx";
    config.gpu_hw_acceleration = "none";
    config.logging_level = 0;
    config.stream_loop = 0;

    Payload p = get_video_payload_common();
    {
        Stream s;

        s.payload = p;
        s.stream_type.type = stream_type::file;
        s.stream_type.file.path = "/home/test/recv";
        s.stream_type.file.filename = "1920x1080p10le_1.yuv";
        config.senders.push_back(s);

        s.stream_type.file.path = "";
        s.stream_type.file.filename = "1920x1080p10le_2.yuv";
        config.senders.push_back(s);
    }

    {
        Stream s;

        s.payload = p;
        s.stream_type.type = stream_type::st2110;
        s.stream_type.st2110.network_interface = "0000:4b:11.1";
        s.stream_type.st2110.local_ip = "192.168.2.2";
        s.stream_type.st2110.remote_ip = "192.168.2.1";
        s.stream_type.st2110.transport = "st2110-20";
        s.stream_type.st2110.remote_port = 20000;
        s.stream_type.st2110.payload_type = 112;
        s.stream_type.st2110.queues_cnt = 0;

        config.receivers.push_back(s);

        s.stream_type.st2110.remote_port = 20001;
        config.receivers.push_back(s);
    }
}

void fill_conf_sender_mcm(Config &config) {
    config.function = "tx";
    config.gpu_hw_acceleration = "none";
    config.logging_level = 0;
    config.stream_loop = 2;

    Payload p = get_video_payload_common();
    {
        Stream s;

        s.payload = p;
        s.stream_type.type = stream_type::file;
        s.stream_type.file.path = "/home/test";
        s.stream_type.file.filename = "1920x1080p10le_1.yuv";
        config.receivers.push_back(s);
    }

    {
        Stream s;

        s.payload = p;
        s.stream_type.type = stream_type::mcm;
        s.stream_type.mcm.conn_type = "st2110";
        s.stream_type.mcm.transport = "st2110-20";
        s.stream_type.mcm.transport_pixel_format = "yuv422p10rfc4175";
        s.stream_type.mcm.ip = "192.168.96.11";
        s.stream_type.mcm.port = 9002;
        s.stream_type.mcm.urn = "abc";

        config.senders.push_back(s);
    }
}

void fill_conf_receiver_mcm(Config &config) {
    config.function = "rx";
    config.gpu_hw_acceleration = "none";
    config.logging_level = 0;
    config.stream_loop = 0;

    Payload p = get_video_payload_common();
    {
        Stream s;

        s.payload = p;
        s.stream_type.type = stream_type::file;
        s.stream_type.file.path = "/home/test/recv";
        s.stream_type.file.filename = "1920x1080p10le_1.yuv";
        config.senders.push_back(s);
    }

    {
        Stream s;

        s.payload = p;
        s.stream_type.type = stream_type::mcm;
        s.stream_type.mcm.conn_type = "st2110";
        s.stream_type.mcm.transport = "st2110-20";
        s.stream_type.mcm.transport_pixel_format = "yuv422p10rfc4175";
        s.stream_type.mcm.ip = "192.168.96.10";
        s.stream_type.mcm.port = 9002;
        s.stream_type.mcm.urn = "abc";

        config.receivers.push_back(s);
    }
}

void fill_conf_multiviewer(Config &config) {
    config.function = "multiviewer";
    config.gpu_hw_acceleration = "intel";
    config.gpu_hw_acceleration_device = "/dev/dri/renderD128";
    config.multiviewer_columns = 3;
    config.logging_level = 0;
    config.stream_loop = 0;

    Payload p = get_video_payload_common();
    {
        Stream s;

        s.payload = p;
        s.stream_type.type = stream_type::file;
        s.stream_type.file.path = "/videos";
        s.stream_type.file.filename = "1920x1080p10le_1.yuv";
        config.receivers.push_back(s);

        s.stream_type.file.filename = "1920x1080p10le_2.yuv";
        config.receivers.push_back(s);

        s.stream_type.file.filename = "1920x1080p10le_1.yuv";
        config.receivers.push_back(s);
        s.stream_type.file.filename = "1920x1080p10le_2.yuv";
        config.receivers.push_back(s);
        s.stream_type.file.filename = "1920x1080p10le_1.yuv";
        config.receivers.push_back(s);
        s.stream_type.file.filename = "1920x1080p10le_2.yuv";
        config.receivers.push_back(s);
        s.stream_type.file.filename = "1920x1080p10le_1.yuv";
        config.receivers.push_back(s);
    }

    {
        Stream s;

        s.payload = p;
        s.stream_type.type = stream_type::file;
        s.stream_type.file.path = "/videos/recv";
        s.stream_type.file.filename = "1920x1080p10le_1.yuv";
        config.senders.push_back(s);
    }
}

void fill_conf_convert(Config &config) {
    config.function = "rx";
    config.gpu_hw_acceleration = "none";
    config.logging_level = 0;
    config.stream_loop = 0;

    {
        Stream s;

        s.payload = get_video_payload_common();
        s.stream_type.type = stream_type::file;
        s.stream_type.file.path = "/videos/";
        s.stream_type.file.filename = "1920x1080p10le_1.yuv";
        config.receivers.push_back(s);
    }

    {
        Stream s;

        s.payload.type = payload_type::video;
        s.payload.video.frame_width = 1280;
        s.payload.video.frame_height = 720;
        s.payload.video.pixel_format = "yuv422p";
        s.payload.video.video_type = "x264";
        s.payload.video.frame_rate = {5, 1};

        s.stream_type.type = stream_type::file;
        s.stream_type.file.path = "/videos/recv";
        s.stream_type.file.filename = "1920x1080p10le_1.mp4";
        config.senders.push_back(s);
    }
}

void fill_conf_recorder(Config &config) {
    config.function = "recorder";
    config.gpu_hw_acceleration = "none";
    config.logging_level = 0;
    config.stream_loop = 0;

    Payload p = get_video_payload_common();
    {
        Stream s;

        s.payload = p;
        s.stream_type.type = stream_type::file;
        s.stream_type.file.path = "/videos";
        s.stream_type.file.filename = "1920x1080p10le_1.yuv";
        config.receivers.push_back(s);
    }

    {
        Stream s;
        s.payload = p;

        s.payload.video.frame_width = 640;
        s.payload.video.frame_height = 360;
        s.stream_type.type = stream_type::file;
        s.stream_type.file.path = "/videos/recv";
        s.stream_type.file.filename = "recv_1.yuv";
        config.senders.push_back(s);

        s.payload.video.frame_width = 1280;
        s.payload.video.frame_height = 720;
        s.payload.video.video_type = "h263p";
        s.stream_type.type = stream_type::file;
        s.stream_type.file.path = "/videos/recv";
        s.stream_type.file.filename = "recv_2.mov";
        config.senders.push_back(s);
    }
}

void fill_conf_upscale(Config &config) {
    config.function = "upscale";
    config.gpu_hw_acceleration = "none";
    config.logging_level = 0;
    config.stream_loop = 0;

    {
        Stream s;

        s.payload = get_video_payload_common();
        s.stream_type.type = stream_type::file;
        s.stream_type.file.path = "/videos/";
        s.stream_type.file.filename = "1920x1080p10le_1.yuv";
        config.receivers.push_back(s);
    }

    {
        Stream s;

        s.payload.type = payload_type::video;
        s.payload.video.frame_width = 3840;
        s.payload.video.frame_height = 2160;
        s.payload.video.pixel_format = "yuv422p10le";
        s.payload.video.frame_rate = {30, 1};

        s.stream_type.type = stream_type::file;
        s.stream_type.file.path = "/videos/recv";
        s.stream_type.file.filename = "3840x2160p10le_1.yuv";
        config.senders.push_back(s);
    }
}


TEST(FFmpegPipelineGeneratorTest, test_sender) {
    Config conf;
    fill_conf_sender(conf);

    std::string pipeline_string;

    if (ffmpeg_generate_pipeline(conf, pipeline_string) != 0) {
            ASSERT_EQ(1, 0) << "Error generating sender pipeline" << std::endl;
    }
    std::string expected_string = " -stream_loop -1 -y -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /home/test/1920x1080p10le_1.yuv -p_port 0000:4b:11.0 -p_sip 192.168.2.1 -udp_port 20000 -payload_type 112 -p_tx_ip 192.168.2.2 -f mtl_st20p - -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /home/test/1920x1080p10le_2.yuv -p_port 0000:4b:11.0 -p_sip 192.168.2.1 -udp_port 20001 -payload_type 112 -p_tx_ip 192.168.2.2 -f mtl_st20p -";
    ASSERT_EQ(pipeline_string.compare(expected_string) == 0, 1) << "Expected: " << std::endl << expected_string << std::endl << " Got: " << std::endl << pipeline_string << std::endl;
}

TEST(FFmpegPipelineGeneratorTest, test_receiver) {
    Config conf;
    fill_conf_receiver(conf);

    std::string pipeline_string;

    if (ffmpeg_generate_pipeline(conf, pipeline_string) != 0) {
            ASSERT_EQ(1, 0) << "Error generating receiver pipeline" << std::endl;
    }
    std::string expected_string = " -y -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -p_port 0000:4b:11.1 -p_sip 192.168.2.2 -udp_port 20000 -payload_type 112 -p_rx_ip 192.168.2.1 -f mtl_st20p -i \"0\" -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo /home/test/recv/1920x1080p10le_1.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -p_port 0000:4b:11.1 -p_sip 192.168.2.2 -udp_port 20001 -payload_type 112 -p_rx_ip 192.168.2.1 -f mtl_st20p -i \"1\" -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo 1920x1080p10le_2.yuv";
    ASSERT_EQ(pipeline_string.compare(expected_string) == 0, 1) << "Expected: " << std::endl << expected_string << std::endl << " Got: " << std::endl << pipeline_string << std::endl;
}

TEST(FFmpegPipelineGeneratorTest, test_sender_queues_cnt) {
    Config conf;
    fill_conf_sender(conf);
    conf.senders[0].stream_type.st2110.queues_cnt = 2;

    std::string pipeline_string;

    if (ffmpeg_generate_pipeline(conf, pipeline_string) != 0) {
            ASSERT_EQ(1, 0) << "Error generating sender pipeline" << std::endl;
    }
    std::string expected_string = " -stream_loop -1 -y -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /home/test/1920x1080p10le_1.yuv -p_port 0000:4b:11.0 -p_sip 192.168.2.1 -udp_port 20000 -payload_type 112 -p_tx_ip 192.168.2.2 -tx_queues 2 -f mtl_st20p - -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /home/test/1920x1080p10le_2.yuv -p_port 0000:4b:11.0 -p_sip 192.168.2.1 -udp_port 20001 -payload_type 112 -p_tx_ip 192.168.2.2 -f mtl_st20p -";
    ASSERT_EQ(pipeline_string.compare(expected_string) == 0, 1) << "Expected: " << std::endl << expected_string << std::endl << " Got: " << std::endl << pipeline_string << std::endl;
}

TEST(FFmpegPipelineGeneratorTest, test_receiver_queues_cnt) {
    Config conf;
    fill_conf_receiver(conf);
    conf.receivers[0].stream_type.st2110.queues_cnt = 4;

    std::string pipeline_string;

    if (ffmpeg_generate_pipeline(conf, pipeline_string) != 0) {
            ASSERT_EQ(1, 0) << "Error generating receiver pipeline" << std::endl;
    }
    std::string expected_string = " -y -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -p_port 0000:4b:11.1 -p_sip 192.168.2.2 -udp_port 20000 -payload_type 112 -p_rx_ip 192.168.2.1 -rx_queues 4 -f mtl_st20p -i \"0\" -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo /home/test/recv/1920x1080p10le_1.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -p_port 0000:4b:11.1 -p_sip 192.168.2.2 -udp_port 20001 -payload_type 112 -p_rx_ip 192.168.2.1 -f mtl_st20p -i \"1\" -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo 1920x1080p10le_2.yuv";
    ASSERT_EQ(pipeline_string.compare(expected_string) == 0, 1) << "Expected: " << std::endl << expected_string << std::endl << " Got: " << std::endl << pipeline_string << std::endl;
}

TEST(FFmpegPipelineGeneratorTest, test_multiviewer) {
    Config conf;
    fill_conf_multiviewer(conf);

    std::string pipeline_string;

    if (ffmpeg_generate_pipeline(conf, pipeline_string) != 0) {
            ASSERT_EQ(1, 0) << "Error generating multiviewer pipeline" << std::endl;
    }
    std::string expected_string = " -y -qsv_device /dev/dri/renderD128 -hwaccel qsv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_2.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_2.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_2.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -filter_complex \"[0:v]hwupload=extra_hw_frames=1,scale_qsv=640:360[out0];[1:v]hwupload=extra_hw_frames=1,scale_qsv=640:360[out1];[2:v]hwupload=extra_hw_frames=1,scale_qsv=640:360[out2];[3:v]hwupload=extra_hw_frames=1,scale_qsv=640:360[out3];[4:v]hwupload=extra_hw_frames=1,scale_qsv=640:360[out4];[5:v]hwupload=extra_hw_frames=1,scale_qsv=640:360[out5];[6:v]hwupload=extra_hw_frames=1,scale_qsv=640:360[out6];[out0][out1][out2][out3][out4][out5][out6]xstack_qsv=inputs=7:layout=0_0|640_0|1280_0|0_360|640_360|1280_360|0_720,hwdownload,format=y210le,format=yuv422p10le\" -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo /videos/recv/1920x1080p10le_1.yuv";
    ASSERT_EQ(pipeline_string.compare(expected_string) == 0, 1) << "Expected: " << std::endl << expected_string << std::endl << " Got: " << std::endl << pipeline_string << std::endl;
}

TEST(FFmpegPipelineGeneratorTest, test_multiviewer_2) {
    Config conf;
    fill_conf_multiviewer(conf);

    conf.receivers[0].payload.video.video_type = "";
    conf.receivers[0].stream_type.file.filename = "compressed720p.mp4";

    std::string pipeline_string;

    if (ffmpeg_generate_pipeline(conf, pipeline_string) != 0) {
            ASSERT_EQ(1, 0) << "Error generating multiviewer pipeline" << std::endl;
    }
    std::string expected_string = " -y -qsv_device /dev/dri/renderD128 -hwaccel qsv -i /videos/compressed720p.mp4 -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_2.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_2.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_2.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -filter_complex \"[0:v]hwupload=extra_hw_frames=1,scale_qsv=640:360[out0];[1:v]hwupload=extra_hw_frames=1,scale_qsv=640:360[out1];[2:v]hwupload=extra_hw_frames=1,scale_qsv=640:360[out2];[3:v]hwupload=extra_hw_frames=1,scale_qsv=640:360[out3];[4:v]hwupload=extra_hw_frames=1,scale_qsv=640:360[out4];[5:v]hwupload=extra_hw_frames=1,scale_qsv=640:360[out5];[6:v]hwupload=extra_hw_frames=1,scale_qsv=640:360[out6];[out0][out1][out2][out3][out4][out5][out6]xstack_qsv=inputs=7:layout=0_0|640_0|1280_0|0_360|640_360|1280_360|0_720,hwdownload,format=y210le,format=yuv422p10le\" -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo /videos/recv/1920x1080p10le_1.yuv";
    ASSERT_EQ(pipeline_string.compare(expected_string) == 0, 1) << "Expected: " << std::endl << expected_string << std::endl << " Got: " << std::endl << pipeline_string << std::endl;
}

TEST(FFmpegPipelineGeneratorTest, test_multiviewer_3) {
    Config conf;
    fill_conf_multiviewer(conf);
    conf.gpu_hw_acceleration = "none";

    std::string pipeline_string;

    if (ffmpeg_generate_pipeline(conf, pipeline_string) != 0) {
            ASSERT_EQ(1, 0) << "Error generating multiviewer pipeline" << std::endl;
    }
    std::string expected_string = " -y -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_2.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_2.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_2.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -filter_complex \"[0:v]scale=640:360[out0];[1:v]scale=640:360[out1];[2:v]scale=640:360[out2];[3:v]scale=640:360[out3];[4:v]scale=640:360[out4];[5:v]scale=640:360[out5];[6:v]scale=640:360[out6];[out0][out1][out2][out3][out4][out5][out6]xstack=inputs=7:layout=0_0|640_0|1280_0|0_360|640_360|1280_360|0_720,format=yuv422p10le\" -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo /videos/recv/1920x1080p10le_1.yuv";
    ASSERT_EQ(pipeline_string.compare(expected_string) == 0, 1) << "Expected: " << std::endl << expected_string << std::endl << " Got: " << std::endl << pipeline_string << std::endl;
}

TEST(FFmpegPipelineGeneratorTest, test_convert) {
    Config conf;
    fill_conf_convert(conf);

    std::string pipeline_string;

    if (ffmpeg_generate_pipeline(conf, pipeline_string) != 0) {
            ASSERT_EQ(1, 0) << "Error generating convert pipeline" << std::endl;
    }
    std::string expected_string = " -y -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -pix_fmt yuv422p -vf scale=1280x720 -r 5/1 -c:v x264 /videos/recv/1920x1080p10le_1.mp4";
    ASSERT_EQ(pipeline_string.compare(expected_string) == 0, 1) << "Expected: " << std::endl << expected_string << std::endl << " Got: " << std::endl << pipeline_string << std::endl;
}

TEST(FFmpegPipelineGeneratorTest, test_recorder) {
    Config conf;
    fill_conf_recorder(conf);
    conf.senders[1].payload.video.preset = "veryfast";
    conf.senders[1].payload.video.profile = "main";

    std::string pipeline_string;

    if (ffmpeg_generate_pipeline(conf, pipeline_string) != 0) {
            ASSERT_EQ(1, 0) << "Error generating convert pipeline" << std::endl;
    }
    std::string expected_string = " -y -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -filter_complex \"split=2[in0][in1];[in0]scale=640:360[out0];[in1]scale=1280:720[out1];\" -map \"[out0]\" -video_size 640x360 -pix_fmt yuv422p10le -r 30/1 -f rawvideo /videos/recv/recv_1.yuv -map \"[out1]\" -c:v h263p  -preset veryfast  -profile main /videos/recv/recv_2.mov";
    ASSERT_EQ(pipeline_string.compare(expected_string) == 0, 1) << "Expected: " << std::endl << expected_string << std::endl << " Got: " << std::endl << pipeline_string << std::endl;
}

TEST(FFmpegPipelineGeneratorTest, test_upscale) {
    Config conf;
    fill_conf_upscale(conf);

    std::string pipeline_string;

    if (ffmpeg_generate_pipeline(conf, pipeline_string) != 0) {
            ASSERT_EQ(1, 0) << "Error generating convert pipeline" << std::endl;
    }
    std::string expected_string = " -y -init_hw_device vaapi=va -init_hw_device opencl@va -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -vf \"format=yuv420p,hwupload,raisr_opencl,hwdownload,format=yuv420p,format=yuv422p10le\" /videos/recv/3840x2160p10le_1.yuv";
    ASSERT_EQ(pipeline_string.compare(expected_string) == 0, 1) << "Expected: " << std::endl << expected_string << std::endl << " Got: " << std::endl << pipeline_string << std::endl;
}

TEST(FFmpegPipelineGeneratorTest, test_mcm_sender) {
    Config conf;
    fill_conf_sender_mcm(conf);

    std::string pipeline_string;

    if (ffmpeg_generate_pipeline(conf, pipeline_string) != 0) {
            ASSERT_EQ(1, 0) << "Error generating sender pipeline" << std::endl;
    }
    std::string expected_string = " -stream_loop 2 -y -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /home/test/1920x1080p10le_1.yuv -f mcm -conn_type st2110 -transport st2110-20 -transport_pixel_format yuv422p10rfc4175 -ip_addr 192.168.96.11 -port 9002 -";
    ASSERT_EQ(pipeline_string.compare(expected_string) == 0, 1) << "Expected: " << std::endl << expected_string << std::endl << " Got: " << std::endl << pipeline_string << std::endl;
}

TEST(FFmpegPipelineGeneratorTest, test_mcm_receiver) {
    Config conf;
    fill_conf_receiver_mcm(conf);

    std::string pipeline_string;

    if (ffmpeg_generate_pipeline(conf, pipeline_string) != 0) {
            ASSERT_EQ(1, 0) << "Error generating receiver pipeline" << std::endl;
    }
    std::string expected_string = " -y -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -f mcm -conn_type st2110 -transport st2110-20 -transport_pixel_format yuv422p10rfc4175 -ip_addr 192.168.96.10 -port 9002 -i \"0\" -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo /home/test/recv/1920x1080p10le_1.yuv";
    ASSERT_EQ(pipeline_string.compare(expected_string) == 0, 1) << "Expected: " << std::endl << expected_string << std::endl << " Got: " << std::endl << pipeline_string << std::endl;
}

TEST(FFmpegPipelineConfigTest, serialize_deserialize_multiviewer) {
    Config conf_reference;
    fill_conf_multiviewer(conf_reference);
    conf_reference.senders[0].payload.video.video_type = "hevc_qsv";
    conf_reference.senders[0].payload.video.preset = "veryfast";
    conf_reference.senders[0].payload.video.profile = "main";
    conf_reference.senders[0].stream_type.file.filename = "1920x1080p10le_1.mp4";

    std::string pipeline_string_reference;

    if (ffmpeg_generate_pipeline(conf_reference, pipeline_string_reference) != 0) {
            ASSERT_EQ(1, 0) << "Error generating multiviewer pipeline" << std::endl;
    }

    std::string json_conf_serialized;
    if(serialize_config_json(conf_reference, json_conf_serialized) != 0) {
        ASSERT_EQ(1, 0) << "Error serializing config" << std::endl;
    }

    Config conf_deserialized;
    std::string pipeline_string_deserialized;
    if(deserialize_config_json(conf_deserialized, json_conf_serialized) != 0) {
        ASSERT_EQ(1, 0) << "Error serializing config" << std::endl;
    }
    if (ffmpeg_generate_pipeline(conf_deserialized, pipeline_string_deserialized) != 0) {
        ASSERT_EQ(1, 0) << "Error generating convert pipeline after deserialization" << std::endl;
    }

    ASSERT_EQ(pipeline_string_reference.compare(pipeline_string_deserialized) == 0, 1) << "Expected: " << std::endl << pipeline_string_reference 
    << std::endl << " Got: " << std::endl << pipeline_string_deserialized << std::endl;
}

TEST(FFmpegPipelineConfigTest, serialize_deserialize_upscale) {
    Config conf_reference;
    fill_conf_upscale(conf_reference);

    std::string pipeline_string_reference;

    if (ffmpeg_generate_pipeline(conf_reference, pipeline_string_reference) != 0) {
        ASSERT_EQ(1, 0) << "Error generating convert pipeline" << std::endl;
    }

    std::string json_conf_serialized;
    if(serialize_config_json(conf_reference, json_conf_serialized) != 0) {
        ASSERT_EQ(1, 0) << "Error serializing config" << std::endl;
    }

    Config conf_deserialized;
    std::string pipeline_string_deserialized;
    if(deserialize_config_json(conf_deserialized, json_conf_serialized) != 0) {
        ASSERT_EQ(1, 0) << "Error serializing config" << std::endl;
    }
    if (ffmpeg_generate_pipeline(conf_deserialized, pipeline_string_deserialized) != 0) {
        ASSERT_EQ(1, 0) << "Error generating convert pipeline after deserialization" << std::endl;
    }

    ASSERT_EQ(pipeline_string_reference.compare(pipeline_string_deserialized) == 0, 1) << "Expected: " << std::endl << pipeline_string_reference 
    << std::endl << " Got: " << std::endl << pipeline_string_deserialized << std::endl;
}

TEST(FFmpegPipelineConfigTest, serialize_deserialize_sender_st2110) {
    Config conf_reference;
    fill_conf_sender(conf_reference);
    conf_reference.senders[0].stream_type.st2110.queues_cnt = 8;

    std::string pipeline_string_reference;

    if (ffmpeg_generate_pipeline(conf_reference, pipeline_string_reference) != 0) {
        ASSERT_EQ(1, 0) << "Error generating convert pipeline" << std::endl;
    }

    std::string json_conf_serialized;
    if(serialize_config_json(conf_reference, json_conf_serialized) != 0) {
        ASSERT_EQ(1, 0) << "Error serializing config" << std::endl;
    }

    Config conf_deserialized;
    std::string pipeline_string_deserialized;
    if(deserialize_config_json(conf_deserialized, json_conf_serialized) != 0) {
        ASSERT_EQ(1, 0) << "Error serializing config" << std::endl;
    }
    if (ffmpeg_generate_pipeline(conf_deserialized, pipeline_string_deserialized) != 0) {
        ASSERT_EQ(1, 0) << "Error generating convert pipeline after deserialization" << std::endl;
    }

    ASSERT_EQ(pipeline_string_reference.compare(pipeline_string_deserialized) == 0, 1) << "Expected: " << std::endl << pipeline_string_reference 
    << std::endl << " Got: " << std::endl << pipeline_string_deserialized << std::endl;
}