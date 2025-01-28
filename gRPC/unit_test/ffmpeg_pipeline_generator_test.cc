#include <gtest/gtest.h>
#include "ffmpeg_pipeline_generator.hpp"

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
        config.senders.push_back(s);

        s.stream_type.st2110.remote_port = 20001;
        config.senders.push_back(s);
    }
}

void fill_conf_receiver(Config &config) {
    config.function = "rx";
    config.gpu_hw_acceleration = "none";
    config.logging_level = 0;

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

        config.receivers.push_back(s);

        s.stream_type.st2110.remote_port = 20001;
        config.receivers.push_back(s);
    }
}

void fill_conf_multiviewer(Config &config) {
    config.function = "multiviewer";
    config.gpu_hw_acceleration = "intel";
    config.logging_level = 0;

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
        s.payload.video.frame_rate = {5, 1};

        s.stream_type.type = stream_type::file;
        s.stream_type.file.path = "/videos/recv";
        s.stream_type.file.filename = "1920x1080p10le_1.mp4";
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
    std::string expected_string = " -y -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /home/test/1920x1080p10le_1.yuv -p_port 0000:4b:11.0 -p_sip 192.168.2.1 -udp_port 20000 -payload_type 112 -p_tx_ip 192.168.2.2 -f mtl_st20p - -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /home/test/1920x1080p10le_2.yuv -p_port 0000:4b:11.0 -p_sip 192.168.2.1 -udp_port 20001 -payload_type 112 -p_tx_ip 192.168.2.2 -f mtl_st20p -";
    ASSERT_EQ(pipeline_string.compare(expected_string) == 0, 1) << "Expected: " << std::endl << expected_string << std::endl << " Got: " << std::endl << pipeline_string << std::endl;
}

TEST(FFmpegPipelineGeneratorTest, test_receiver) {
    Config conf;
    fill_conf_receiver(conf);

    std::string pipeline_string;

    if (ffmpeg_generate_pipeline(conf, pipeline_string) != 0) {
            ASSERT_EQ(1, 0) << "Error generating receiver pipeline" << std::endl;
    }
    std::string expected_string = " -y -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -p_port 0000:4b:11.1 -p_sip 192.168.2.2 -udp_port 20000 -payload_type 112 -p_rx_ip 192.168.2.1 -f mtl_st20p -i \"0\" /home/test/recv/1920x1080p10le_1.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -p_port 0000:4b:11.1 -p_sip 192.168.2.2 -udp_port 20001 -payload_type 112 -p_rx_ip 192.168.2.1 -f mtl_st20p -i \"1\" 1920x1080p10le_2.yuv";
    ASSERT_EQ(pipeline_string.compare(expected_string) == 0, 1) << "Expected: " << std::endl << expected_string << std::endl << " Got: " << std::endl << pipeline_string << std::endl;
}

TEST(FFmpegPipelineGeneratorTest, test_multiviewer) {
    Config conf;
    fill_conf_multiviewer(conf);

    std::string pipeline_string;

    if (ffmpeg_generate_pipeline(conf, pipeline_string) != 0) {
            ASSERT_EQ(1, 0) << "Error generating multiviewer pipeline" << std::endl;
    }
    std::string expected_string = " -y -qsv_device /dev/dri/renderD128 -hwaccel qsv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_2.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_2.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_2.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -filter_complex \"[0:v]hwupload,scale_qsv=640:360[out0];[1:v]hwupload,scale_qsv=640:360[out1];[2:v]hwupload,scale_qsv=640:360[out2];[3:v]hwupload,scale_qsv=640:360[out3];[4:v]hwupload,scale_qsv=640:360[out4];[5:v]hwupload,scale_qsv=640:360[out5];[6:v]hwupload,scale_qsv=640:360[out6];[out0][out1][out2][out3][out4][out5][out6]xstack_qsv=inputs=7:layout=0_0|640_0|1280_0|0_360|640_360|1280_360|0_720,format=y210le,format=yuv422p10le\" /videos/recv/1920x1080p10le_1.yuv";
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
    std::string expected_string = " -y -qsv_device /dev/dri/renderD128 -hwaccel qsv -i /videos/compressed720p.mp4 -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_2.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_2.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_2.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -filter_complex \"[0:v]hwupload,scale_qsv=640:360[out0];[1:v]hwupload,scale_qsv=640:360[out1];[2:v]hwupload,scale_qsv=640:360[out2];[3:v]hwupload,scale_qsv=640:360[out3];[4:v]hwupload,scale_qsv=640:360[out4];[5:v]hwupload,scale_qsv=640:360[out5];[6:v]hwupload,scale_qsv=640:360[out6];[out0][out1][out2][out3][out4][out5][out6]xstack_qsv=inputs=7:layout=0_0|640_0|1280_0|0_360|640_360|1280_360|0_720,format=y210le,format=yuv422p10le\" /videos/recv/1920x1080p10le_1.yuv";
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
    std::string expected_string = " -y -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_2.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_2.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_2.yuv -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -filter_complex \"[0:v]scale=640:360[out0];[1:v]scale=640:360[out1];[2:v]scale=640:360[out2];[3:v]scale=640:360[out3];[4:v]scale=640:360[out4];[5:v]scale=640:360[out5];[6:v]scale=640:360[out6];[out0][out1][out2][out3][out4][out5][out6]xstack=inputs=7:layout=0_0|640_0|1280_0|0_360|640_360|1280_360|0_720,format=y210le,format=yuv422p10le\" /videos/recv/1920x1080p10le_1.yuv";
      ASSERT_EQ(pipeline_string.compare(expected_string) == 0, 1) << "Expected: " << std::endl << expected_string << std::endl << " Got: " << std::endl << pipeline_string << std::endl;
}

TEST(FFmpegPipelineGeneratorTest, test_convert) {
    Config conf;
    fill_conf_convert(conf);

    std::string pipeline_string;

    if (ffmpeg_generate_pipeline(conf, pipeline_string) != 0) {
            ASSERT_EQ(1, 0) << "Error generating convert pipeline" << std::endl;
    }
    std::string expected_string = " -y -video_size 1920x1080 -pix_fmt yuv422p10le -r 30/1 -f rawvideo -i /videos/1920x1080p10le_1.yuv -pix_fmt yuv422p -vf scale=1280x720 -r 5/1 /videos/recv/1920x1080p10le_1.mp4";   
    ASSERT_EQ(pipeline_string.compare(expected_string) == 0, 1) << "Expected: " << std::endl << expected_string << std::endl << " Got: " << std::endl << pipeline_string << std::endl;
}
