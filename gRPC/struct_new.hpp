#include <string>
#include <vector>
#include <map>

struct FrameRate {
    int numerator;
    int denominator;
};
 
struct Video {
    int frame_width;
    int frame_height;
    FrameRate frame_rate;
    std::string pixel_format;
    std::string video_type; 
    // in case of rawvideo then ffmpeg param = "-f rawwideo"
    // otherwise -c:v <video_type> e.g. -c:v x264
};
 
 // Audio struct is a placeholder for future implementation
struct Audio {
    int channels;
    int sample_rate;
    std::string format;
    std::string packet_time;
};
 
struct File {
    std::string path;
    std::string filename;
};
 
struct ST2110 {
    std::string network_interface; //VFIO port address 0000:00:00.0; ffmpeg param name: -p_port
    std::string local_ip; // ffmpeg param name: -p_sip
    std::string remote_ip; // ffmpeg param name: -p_rx_ip / -p_tx_ip
    std::string transport; 
    int remote_port; // ffmpeg param name: -udp_port
    int payload_type;
};
 
struct MCM {
    std::string conn_type;
    std::string transport;
    std::string transport_pixel_format;
    std::string ip;
    int port;
    std::string urn;
};
 
enum payload_type {
    video = 0,
    audio
 };

struct Payload {
    payload_type type;
    Video video;
    Audio audio;
};
 
enum stream_type {
    file = 0,
    st2110,
    mcm
 };

struct StreamType {
    stream_type type; 
    File file;
    ST2110 st2110;
    MCM mcm;
};
 
struct Stream {
    Payload payload;
    StreamType stream_type;
};
 
struct Config {
    std::vector<Stream> senders;
    std::vector<Stream> receivers;

    std::string function; //multiviewer, upscale, replay, recorder, jpegxs, rx, tx
    std::string gpu_hw_acceleration; //intel, nvidia, none
    int logging_level;
};