#include "handlers.h"

// Handlers map definition
std::unordered_map<std::string, HandlerFunc> handlers = {
    // per stream options
    {"codec", handleCodecKeyValue},
    {"filter", handleKeyValue},
    {"format", handleFormatKeyValue},

    // video options
    {"height", handleKeyValue},
    {"width", handleKeyValue},

    // common options for MCM, MTL
    {"payload_type", handleKeyValue},
    {"video_size", handleKeyValue},
    {"pixel_format", handleKeyValue},

    // MCM Muxer/Demuxer
    {"ip_addr", handleKeyValue},
    {"port", handleKeyValue},
    {"protocol_type", handleKeyValue},
    {"frame_rate", handleFrameRate},
    {"socket_name", handleKeyValue},
    {"interface_id", handleKeyValue},

    // MTL Device Arguments
    {"p_port", handleKeyValue},
    {"p_sip", handleKeyValue},
    {"dma_dev", handleKeyValue},

    // Tx Port Encoding Arguments
    {"p_tx_ip", handleKeyValue},
    {"udp_port", handleKeyValue},

    // Rx Port Decoding Arguments
    {"p_rx_ip", handleKeyValue},
    {"udp_port", handleKeyValue},

    // MTL st20p Muxer/Demuxer
    {"fb_cnt", handleKeyValue},
    {"pix_fmt", handleKeyValue},
    {"fps", handleKeyValue},
    {"timeout_s", handleKeyValue},

    // MTL st22p Muxer/Demuxer
    {"bpp", handleKeyValue},
    {"codec_thread_cnt", handleKeyValue},
    {"st22_codec", handleKeyValue},

    // MTL st30p Muxer/Demuxer
    {"at", handleKeyValue},
    {"ar", handleKeyValue},
    {"ac", handleKeyValue},
    {"pcm_fmt", handleKeyValue},

    // JPEG XS encoder
    {"decomp_v", handleKeyValue},
    {"decomp_h", handleKeyValue},
    {"threads", handleKeyValue},
    {"slice_height", handleKeyValue},
    {"quantization", handleKeyValue},
    {"coding-signs", handleKeyValue},
    {"coding-sigf", handleKeyValue},
    {"coding-vpred", handleKeyValue},

    // JPEG XS decoder
    {"threads", handleKeyValue},
};

// Helper function to handle prefixed key-value pairs
std::string handlePrefixedKeyValue(const std::string& prefix, const std::string& key, const json& value) {
    if (value.is_null()) {
        throw std::invalid_argument("key " + key + " has no value");
    }
    std::string valueStr = value.get<std::string>();
    if (valueStr.empty()) {
        throw std::invalid_argument("key " + key + " has an empty value");
    }
    std::cout << prefix << ": " << valueStr << std::endl;
    return prefix + " " + valueStr;
}

// handleFormatKeyValue handles the key "format", the value passed to ffmpeg as -f <format value>
std::string handleFormatKeyValue(const std::string& key, const json& value) {
    return handlePrefixedKeyValue("-f", key, value);
}

// handleCodecKeyValue handles the key "codec", the value passed to ffmpeg as -c <codec value>
std::string handleCodecKeyValue(const std::string& key, const json& value) {
    return handlePrefixedKeyValue("-c", key, value);
}

// handleKeyValue handles the key and value, prints them, and returns a formatted string
std::string handleKeyValue(const std::string& key, const json& value) {
    return handlePrefixedKeyValue("-" + key, key, value);
}

// handleFrameRate handles the frame_rate key by calculating the frame rate and returning a formatted string
std::string handleFrameRate(const std::string& key, const json& value) {
    if (value.is_object()) {
        double numerator = value["numerator"].get<double>();
        double denominator = value["denominator"].get<double>();
        if (denominator != 0) {
            int frameRate = static_cast<int>(numerator / denominator);
            std::cout << key << ": " << frameRate << std::endl;
            return "-" + key + " " + std::to_string(frameRate);
        }
        throw std::invalid_argument("invalid frame_rate: denominator is zero");
    }
    throw std::invalid_argument("invalid frame_rate format");
}
