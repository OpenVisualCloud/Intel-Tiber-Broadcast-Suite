
#include "ffmpeg_pipeline_generator.hpp"

void fill_conf_sender(Config &config) {
    config.function = "tx";
    config.gpu_hw_acceleration = "none";
    config.logging_level = 0;

    Payload p;
    p.type = payload_type::video;
    p.video.frame_width = 1920;
    p.video.frame_height = 1080;
    p.video.frame_rate = {30, 1};
    p.video.pixel_format = "yuv422p10le";
    p.video.video_type = "rawvideo";

    {
        Stream s;

        s.payload = p;
        s.stream_type.type = stream_type::file;
        s.stream_type.file.path = "/home/tszumski";
        s.stream_type.file.filename = "1920x1080p10le_1.yuv";
        config.receivers.push_back(s);

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

void fill_conf_receiver(Config &config)
{
    config.function = "rx";
    config.gpu_hw_acceleration = "none";
    config.logging_level = 0;

    Payload p;
    p.type = payload_type::video;
    p.video.frame_width = 1920;
    p.video.frame_height = 1080;
    p.video.frame_rate = {30, 1};
    p.video.pixel_format = "yuv422p10le";
    p.video.video_type = "rawvideo";

    {
        Stream s;

        s.payload = p;
        s.stream_type.type = stream_type::file;
        s.stream_type.file.path = "/home/tszumski/recv";
        s.stream_type.file.filename = "1920x1080p10le_1.yuv";
        config.senders.push_back(s);

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

int main(int argc, char *argv[]) {
    //sender
    {
        Config conf;
        fill_conf_sender(conf);

        std::string pipelinie_string;

        if (ffmpeg_generate_pipeline(conf, pipelinie_string) != 0) {
            pipelinie_string.clear();
            std::cout << "Error generating pipeline" << std::endl;
            return 1;
        }
        std::cout << "Generated sender pipeline: " << std::endl
                  << pipelinie_string << std::endl << std::endl;
    }

    //receiver
    {
        Config conf;
        fill_conf_receiver(conf);

        std::string pipelinie_string;

        if (ffmpeg_generate_pipeline(conf, pipelinie_string) != 0) {
            pipelinie_string.clear();
            std::cout << "Error generating pipeline" << std::endl;
            return 1;
        }
        std::cout << "Generated receiver pipeline: " << std::endl
                  << pipelinie_string << std::endl;
    }
}