# Environment

- Ubuntu 22.04.4 LTS
- kernel: 5.15.0-133-generic
- minikube: v1.34.0 (1 node)
- docker: 27.3.1
- kubectl: v1.31.2
- NMOS supports standards IS-04 and IS-05: <https://specs.amwa.tv/is-04/>, <https://specs.amwa.tv/is-05/>

## Configuration note for NMOS client in the context of Intel Broadcast Suite

`nmos-cpp` repository has been simplified to **IS-04** & **IS-05** implementation only.
The key change is in configuration of senders and receivers for BCS pipeline.

BCS Pipeline is a NMOS client that is treated as one node that has 1 device and has x senders and y receivers that are provided from the level of JSON config `node.json`.
Here is sample config `node.json` (treated as transmitter node with 1 device and 1 sender):

```json
{
  "logging_level": 10,
  "http_port": 90,
  "label": "intel-broadcast-suite",
  "device_tags": {
    "pipeline": ["tx"]
  },
  "color_sampling": "YCbCr-4:2:2",
  "function": "tx",
  "gpu_hw_acceleration": "none",
  "domain": "local",
  "ffmpeg_grpc_server_address": "localhost",
  "ffmpeg_grpc_server_port": "50051",
  "sender_payload_type":112,
  "frame_rate": { "numerator": 60000, "denominator": 1001 },
  "sender": [{
    "stream_payload": {
      "video": {
        "frame_width": 1920,
        "frame_height": 1080,
        "frame_rate": { "numerator": 60, "denominator": 1 },
        "pixel_format": "yuv422p10le",
        "video_type": "rawvideo"
      },
      "audio": {
        "channels": 2,
        "sampleRate": 48000,
        "format": "pcm_s24be",
        "packetTime": "1ms"
      }
    },
    "stream_type": {
      "st2110": {
        "transport": "st2110-20",
        "payloadType" : 112
      }
    }
  }],
  "receiver": [{
    "stream_payload": {
      "video": {
        "frame_width": 1920,
        "frame_height": 1080,
        "frame_rate": { "numerator": 60, "denominator": 1 },
        "pixel_format": "yuv422p10le",
        "video_type": "rawvideo"
      },
      "audio": {
        "channels": 2,
        "sampleRate": 48000,
        "format": "pcm_s24be",
        "packetTime": "1ms"
      }
    },
    "stream_type": {
      "file": {
        "path": "/root",
        "filename": "1920x1080p10le_1.yuv"
      }
    }
  }]
}
```

Curretly only video mode is supported. The audio support is under development and will be relesed too.
- `logging_level`: The level of logging detail.
- `http_port`: The port number for HTTP communication (90 in this case).
- `label`: A label or identifier for the configuration ("intel-broadcast-suite").
- `color_sampling`: The color sampling format ("YCbCr-4:2:2").
- `function`: The function of the device, here indicating the pipeline type ("tx" for transmit).
- `gpu_hw_acceleration`: Indicates if GPU hardware acceleration is used ("none").
- `domain`: The domain of the device ("local").
- `ffmpeg_grpc_server_address`: The address of the FFmpeg gRPC server ("localhost").
- `ffmpeg_grpc_server_port`: The port of the FFmpeg gRPC server (50051).
- `sender_payload_type`: The payload type for the sender (112).
- `frame_rate`: The frame rate for the video, given as a fraction ({"numerator": 60000, "denominator": 1001}).
- `sender`: An array of sender configurations:
  - `stream_payload`: Contains details about the video and audio streams:
    - `video`: Details about the video stream
      - `frame_width`: Width of the video frame (1920).
      - `frame_height`: Height of the video frame (1080).
      - `frame_rate`: Frame rate of the video ({"numerator": 60, "denominator": 1}).
      - `pixel_format`: Pixel format of the video ("yuv422p10le").
      - `video_type`: Type of video ("rawvideo").
    - `audio`: Details about the audio stream:
      - `channels`: Number of audio channels (2).
      - `sampleRate`: Sample rate of the audio (48000).
      - `format`: Audio format ("pcm_s24be").
      - `packetTime`: Packet time for the audio ("1ms").
    - `stream_type`: Type of stream:
      - `st2110`: Details for ST 2110 transport:
        - `transport`: Transport type ("st2110-20").
        - `payloadType`: Payload type (112).
- `receiver`: An array of receiver configurations:
  - `stream_payload`: Contains details about the video stream that acts as ffmpeg receiver. Just to indicate the ffmpeg pipeline the source of the video.
  - `stream_type`: Type of stream:
    - `file`: Details for file-based stream:
      - `path`: Path to the file ("/root").
      - `filename`: Filename ("1920x1080p10le_1.yuv").

For testing purposes there are also NMOS sample cotroller, NMOS registry pod and NMOS testing tool for validation of features.

## Installation

### Docker option

Go to project root directory and run:

```bash
./build.sh
```
```bash
cd scripts/
```
```bash
./first_run.sh
```

### For development purposes

#### 1. Local build

```bash
cd src/
```
```bash
./build_local.sh
```

Binaries related to `nmos-node` will be located in the `src/nmos/nmos-node/build` directory.

### License

```text
SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation

SPDX-License-Identifier: BSD-3-Clause
```
