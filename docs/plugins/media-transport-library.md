# Intel® Media Transport Library

FFmpeg Intel® Media Transport Library Plugins Documentation

> 💡 _**Tip:** For up to date documentation refer to: [https://github.com/OpenVisualCloud/Media-Transport-Library](https://github.com/OpenVisualCloud/Media-Transport-Library)_

## Global Parameters

### Input Parameters

#### MTL Device Arguments

- **`p_port`**: Mtl hardware device to use. Mtl p_port. Example: "eth0"
  - **Type**: String
  - **Default**: `NULL`
  - **Flags**: Encoding/Decoding parameter

- **`p_sip`**: Mtl local (source) IP address to set/use.
  - **Type**: String
  - **Default**: `NULL`
  - **Flags**: Encoding/Decoding parameter

- **`dma_dev`**: Mtl DMA device.
  - **Type**: String
  - **Default**: `NULL`
  - **Flags**: Encoding/Decoding parameter

#### Tx Port Encoding Arguments

- **`p_tx_ip`**: Transmit (Tx) to IP address.
  - **Type**: String
  - **Default**: `NULL`
  - **Flags**: Encoding parameter

- **`udp_port`**: Transmit (Tx) to UDP port.
  - **Type**: Integer
  - **Default**: `20000`
  - **Flags**: Encoding parameter
  - **Range**: `-1` to `INT_MAX`

- **`payload_type`**: Transmit (Tx) payload type.
  - **Type**: Integer
  - **Default**: `112`
  - **Flags**: Encoding parameter
  - **Range**: `-1` to `INT_MAX`

#### Rx Port Decoding Arguments

- **`p_rx_ip`**: Receive (Rx) from IP address.
  - **Type**: String
  - **Default**: `NULL`
  - **Flags**: Decoding parameter

- **`udp_port`**: Receive (Rx) from UDP port.
  - **Type**: Integer
  - **Default**: `20000`
  - **Flags**: Decoding parameter
  - **Range**: `-1` to `INT_MAX`

- **`payload_type`**: Receive (Rx) payload type.
  - **Type**: Integer
  - **Default**: `112`
  - **Flags**: Decoding parameter
  - **Range**: `-1` to `INT_MAX`

## St20p Muxer Plugin Documentation for FFmpeg Intel® Media Transport Library.

The Mtl St20p Muxer plugin for FFmpeg is designed to handle the transmission of ST 2110-20 video streams over a network. Below are the input parameters that can be configured for the Mtl St20p Muxer plugin. To use plugin prepend input parameters with `-f mtl_st20p`.

### Input Parameters

#### [Device Arguments](#mtl-device-arguments)

#### [Port Encoding Arguments](#tx-port-encoding-arguments)

#### Session Arguments
- **`fb_cnt`**: Frame buffer count.
  - **Type**: Integer
  - **Default**: `3`
  - **Flags**: Encoding parameter
  - **Range**: `3` to `8`

### Example Usage

To use the Mtl St20p Muxer plugin with FFmpeg, you can specify the input parameters using the `-option` flag. Here is an example command that sets some of the parameters:

```bash
ffmpeg -i input.mp4 -c:v rawvideo -f mtl_st20p -p_port "eth0" -p_sip "192.168.1.1" -dma_dev "dma0" -p_tx_ip "239.0.0.1" -udp_port 1234 -payload_type 112 -fb_cnt 4 output.mtl
```

This command takes an input file input.mp4, encodes the video as raw video, and uses the Mtl St20p Muxer to send the data to IP address 239.0.0.1 on UDP port 1234 with payload type 112 and a frame buffer count of 4.

## St20p Demuxer Plugin Documentation for FFmpeg Intel® Media Transport Library.

The Mtl St20p Demuxer plugin for FFmpeg is designed to handle the reception of ST 2110-20 video streams over a network. Below are the input parameters that can be configured for the Mtl St20p Demuxer plugin. To use plugin prepend input parameters with `-f mtl_st20p`.

### Input Parameters

#### [Device Arguments](#mtl-device-arguments)

#### [Port Decoding Arguments](#rx-port-decoding-arguments)

#### Session Arguments
- **`video_size`**: Video frame size.
  - **Type**: Image size (String)
  - **Default**: `"1920x1080"`
  - **Flags**: Decoding parameter

- **`pix_fmt`** / **`pixel_format`**: Pixel format for framebuffer.
  - **Type**: Pixel format (Enum)
  - **Default**: `AV_PIX_FMT_YUV422P10LE`
  - **Flags**: Decoding parameter
  - **Possible Values**: `AV_PIX_FMT_YUV422P10LE`, `AV_PIX_FMT_RGB24`

- **`fps`**: Video frame rate.
  - **Type**: Rational
  - **Default**: `59.94` (interpreted as 60000/1001)
  - **Flags**: Decoding parameter
  - **Range**: `0` to `1000`

- **`timeout_s`**: Frame get timeout in seconds.
  - **Type**: Integer
  - **Default**: `0` (no timeout)
  - **Flags**: Decoding parameter
  - **Range**: `0` to `600` (10 minutes)

- **`fb_cnt`**: Frame buffer count.
  - **Type**: Integer
  - **Default**: `3`
  - **Flags**: Decoding parameter
  - **Range**: `3` to `8`

### Example Usage

To use the Mtl St20p Demuxer plugin with FFmpeg, you can specify the input parameters using the `-option` flag. Here is an example command that sets some of the parameters:

```bash
ffmpeg -f mtl_st20p -p_port "eth0" -p_sip "192.168.1.1" -dma_dev "dma0" -p_rx_ip "239.0.0.1" -udp_port 1234 -payload_type 112 -video_size 1280x720 -pix_fmt yuv422p10le -fps 50 -timeout_s 2 -fb_cnt 4 -i mtl_input -c:v copy output.mp4
```

This command receives an ST 2110-20 video stream with the specified device and port configurations, a frame size of 1280x720, pixel format yuv422p10le, frame rate 50, a timeout of 2 seconds for frame retrieval, and a frame buffer count of 4. The video stream is then copied to an output file output.mp4.

## St22p Muxer Plugin Documentation for FFmpeg Intel® Media Transport Library.

The Mtl St22p Muxer plugin for FFmpeg is designed to handle the transmission of ST 2110-22 video streams over a network. Below are the input parameters that can be configured for the Mtl St22p Muxer plugin. To use plugin prepend input parameters with `-f mtl_st22p`.

### Input Parameters

#### [Device Arguments](#mtl-device-arguments)

#### [Port Encoding Arguments](#tx-port-encoding-arguments)

#### Session Arguments
- **`fb_cnt`**: Frame buffer count.
  - **Type**: Integer
  - **Default**: `3`
  - **Flags**: Encoding parameter
  - **Range**: `3` to `8`

- **`bpp`**: Bit per pixel.
  - **Type**: Float
  - **Default**: `3.0`
  - **Flags**: Encoding parameter
  - **Range**: `0.1` to `8.0`

- **`codec_thread_cnt`**: Codec threads count.
  - **Type**: Integer
  - **Default**: `0`
  - **Flags**: Encoding parameter
  - **Range**: `0` to `64`

- **`st22_codec`**: ST 2110-22 codec.
  - **Type**: String
  - **Default**: `NULL`
  - **Flags**: Encoding parameter

### Example Usage

The Mtl St22p Muxer plugin usage with FFmpeg example command that sets some of the parameters:

```bash
ffmpeg -i input.mp4 -c:v rawvideo -f mtl_st22p -p_port "eth0" -p_sip "192.168.1.1" -dma_dev "dma0" -p_tx_ip "239.0.0.1" -udp_port 1234 -payload_type 112 -fb_cnt 4 -bpp 2.0 -codec_thread_cnt 4 -st22_codec "jpegxs" output.mtl
```

## St22p Demuxer Plugin Documentation for FFmpeg Intel® Media Transport Library.

The Mtl St22p Demuxer plugin for FFmpeg is designed to handle the reception of ST 2110-22 video streams over a network. Below are the input parameters that can be configured for the Mtl St22p Demuxer plugin. To use plugin prepend input parameters with `-f mtl_st22p`.

### Input Parameters

#### [Device Arguments](#mtl-device-arguments)

#### [Port Decoding Arguments](#rx-port-decoding-arguments)

#### Session Arguments
- **`video_size`**: Video frame size.
  - **Type**: Image size (String)
  - **Default**: `"1920x1080"`
  - **Flags**: Decoding parameter

- **`pix_fmt`** / **`pixel_format`**: Pixel format for framebuffer.
  - **Type**: Pixel format (Enum)
  - **Default**: `YUV422P10LE`
  - **Flags**: Decoding parameter
  - **Possible Values**: `YUV422P10LE`, `RGB24`

- **`fps`**: Video frame rate.
  - **Type**: Rational
  - **Default**: `59.94` (interpreted as 60000/1001)
  - **Flags**: Decoding parameter
  - **Range**: `0` to `1000`

- **`timeout_s`**: Frame get timeout in seconds.
  - **Type**: Integer
  - **Default**: `0` (no timeout)
  - **Flags**: Decoding parameter
  - **Range**: `0` to `600` (10 minutes)

- **`fb_cnt`**: Frame buffer count.
  - **Type**: Integer
  - **Default**: `3`
  - **Flags**: Decoding parameter
  - **Range**: `3` to `8`

- **`codec_thread_cnt`**: Codec threads count.
  - **Type**: Integer
  - **Default**: `0`
  - **Flags**: Decoding parameter
  - **Range**: `0` to `64`

- **`st22_codec`**: ST 2110-22 codec.
  - **Type**: String
  - **Default**: `NULL`
  - **Flags**: Decoding parameter

### Example Usage

The Mtl St22p FFmpeg Demuxer plugin usage example command that sets some of the parameters:

```bash
ffmpeg -f mtl_st22p -p_port "eth0" -p_sip "192.168.1.1" -dma_dev "dma0" -p_rx_ip "239.0.0.1" -udp_port 1234 -payload_type 112 -video_size 1280x720 -pix_fmt yuv422p10le -fps 50 -timeout_s 2 -fb_cnt 4 -codec_thread_cnt 4 -st22_codec "jpegxs" -i mtl_input -c:v copy output.mp4
```

This command receives an ST 2110-22 video stream with the specified device and port configurations, a frame size of 1280x720, pixel format yuv422p10le, frame rate 50, a timeout of 2 seconds for frame retrieval, a frame buffer count of 4, codec threads count of 4, and codec jpegxs. The video stream is then copied to an output file output.mp4.

## St30p Muxer Plugin Documentation for FFmpeg Intel® Media Transport Library.

The Mtl St30p Muxer plugin for FFmpeg is designed to handle the transmission of ST 2110-30 audio streams over a network. Below are the input parameters that can be configured for the Mtl St30p Muxer plugin. To use plugin prepend input parameters with `-f mtl_st30p`.

### Input Parameters

#### [Device Arguments](#mtl-device-arguments)

#### [Port Encoding Arguments](#tx-port-encoding-arguments)

#### Session Arguments
- **`fb_cnt`**: Frame buffer count.
  - **Type**: Integer
  - **Default**: `3`
  - **Flags**: Encoding parameter
  - **Range**: `3` to `8000`

- **`at`**: Audio packet time.
  - **Type**: String
  - **Default**: `NULL`
  - **Flags**: Encoding parameter
  - **Possible Values**: `"1ms"`, `"125us"`

### Example Usage

The Mtl St30p FFmpeg Muxer plugin usage example command that sets some of the parameters:

```bash
ffmpeg -i input.wav -c:a pcm_s24be -f mtl_st30p -p_port "eth0" -p_sip "192.168.1.1" -dma_dev "dma0" -p_tx_ip "239.0.0.1" -udp_port 1234 -payload_type 112 -fb_cnt 4 -at "1ms" output.mtl
```

This command takes an input file input.wav, encodes the audio as PCM 24-bit big-endian, and uses the Mtl St30p Muxer to send the data to IP address 239.0.0.1 on UDP port 1234 with payload type 112, frame buffer count of 4, and audio packet time of 1ms.

## St30p Demuxer Plugin Documentation for FFmpeg Intel® Media Transport Library.

The Mtl St30p Demuxer plugin for FFmpeg is designed to handle the reception of ST 2110-30 audio streams over a network. Below are the input parameters that can be configured for the Mtl St30p Demuxer plugin.  To use plugin prepend input parameters with `-f mtl_st30p`.

### Input Parameters

#### [Device Arguments](#mtl-device-arguments)

#### [Port Decoding Arguments](#rx-port-decoding-arguments)

#### Session Arguments
- **`fb_cnt`**: Frame buffer count.
  - **Type**: Integer
  - **Default**: `3`
  - **Flags**: Decoding parameter
  - **Range**: `3` to `8`

- **`timeout_s`**: Frame get timeout in seconds.
  - **Type**: Integer
  - **Default**: `0` (no timeout)
  - **Flags**: Decoding parameter
  - **Range**: `0` to `600` (10 minutes)

- **`ar`**: Audio sampling rate.
  - **Type**: Integer
  - **Default**: `48000`
  - **Flags**: Decoding parameter
  - **Range**: `1` to `INT_MAX`

- **`ac`**: Audio channel.
  - **Type**: Integer
  - **Default**: `2`
  - **Flags**: Decoding parameter
  - **Range**: `1` to `INT_MAX`

- **`pcm_fmt`**: Audio PCM format.
  - **Type**: String
  - **Default**: `NULL`
  - **Flags**: Decoding parameter
  - **Possible Values**: `"pcm24"`, `"pcm16"`, `"pcm8"`

- **`at`**: Audio packet time.
  - **Type**: String
  - **Default**: `NULL`
  - **Flags**: Decoding parameter
  - **Possible Values**: `"1ms"`, `"125us"`

### Example Usage

The Mtl St30p FFmpeg Demuxer plugin usage example command that sets some of the parameters:

```bash
ffmpeg -f mtl_st30p -p_port "eth0" -p_sip "192.168.1.1" -dma_dev "dma0" -p_rx_ip "239.0.0.1" -udp_port 1234 -payload_type 112 -fb_cnt 4 -timeout_s 2 -ar 48000 -ac 2 -pcm_fmt "pcm24" -at "1ms" -i mtl_input -c:a copy output.wav
```

This command receives an ST 2110-30 audio stream with the specified device and port configurations, frame buffer count of 4, a timeout of 2 seconds for frame retrieval, audio sampling rate of 48000, 2 audio channels, PCM format pcm24, and audio packet time of 1ms. The audio stream is then copied to an output file output.wav.

## Appendix A: More Examples

### A.1 ST20P raw video run guide

The MTL ST20P plugin is implemented as an FFMpeg input/output device, enabling direct reading from or sending raw video via the ST2110-20 stream.

#### A.1.1 St20p input

Reading two st2110-20 10bit YUV422 stream, one on "239.168.85.20:20000" and the second on "239.168.85.20:20002":

```bash
ffmpeg -p_port 0000:af:01.0 -p_sip 192.168.96.2 -p_rx_ip 239.168.85.20 -udp_port 20000 -payload_type 112 -fps 59.94 -pix_fmt yuv422p10le -video_size 1920x1080 -f mtl_st20p -i "1" -p_port 0000:af:01.0 -p_rx_ip 239.168.85.20 -udp_port 20002 -payload_type 112 -fps 59.94 -pix_fmt yuv422p10le -video_size 1920x1080 -f mtl_st20p -i "2" -map 0:0 -f rawvideo /dev/null -y -map 1:0 -f rawvideo /dev/null -y
```

Reading a st2110-20 10bit YUV422 stream on "239.168.85.20:20000" with payload_type 112, and use libopenh264 to encode the stream to out.264 file:

```bash
ffmpeg -p_port 0000:af:01.0 -p_sip 192.168.96.2 -p_rx_ip 239.168.85.20 -udp_port 20000 -payload_type 112 -fps 59.94 -pix_fmt yuv422p10le -video_size 1920x1080 -f mtl_st20p -i "k" -c:v libopenh264 out.264 -y
```

#### A.1.2 St20p output

Reading from a yuv stream from a local file and sending a st2110-20 10bit YUV422 stream on "239.168.85.20:20000" with payload_type 112:

```bash
ffmpeg -stream_loop -1 -video_size 1920x1080 -f rawvideo -pix_fmt yuv422p10le -i yuv422p10le_1080p.yuv -filter:v fps=59.94 -p_port 0000:af:01.1 -p_sip 192.168.96.3 -p_tx_ip 239.168.85.20 -udp_port 20000 -payload_type 112 -f mtl_st20p -
```

### A.2 ST22 compressed video run guide

A typical workflow for processing an MTL ST22 compressed stream with FFMpeg is outlined in the following steps: Initially, FFMpeg reads a YUV frame from the input source, then forwards the frame to a codec to encode the raw video into a compressed codec stream. Finally, the codec stream is sent to the MTL ST22 plugin.
The MTL ST22 plugin constructs the codec stream and transmits it as ST2110-22 RTP packets, adhering to the standard. In addition to the JPEG XS stream, the MTL ST22 plugin is capable of supporting various other common compressed codecs, including H264, H265, and HEVC, among others.

#### A.2.1 St22 output

Reading from a yuv stream from local source file, encode with h264 codec and sending a st2110-22 codestream on "239.168.85.20:20000" with payload_type 112:

```bash
ffmpeg -stream_loop -1 -video_size 1920x1080 -f rawvideo -pix_fmt yuv420p -i yuv420p_1080p.yuv -filter:v fps=59.94 -c:v libopenh264 -p_port 0000:af:01.1 -p_sip 192.168.96.3 -p_tx_ip 239.168.85.20 -udp_port 20000 -payload_type 112 -f mtl_st22 -
```

#### A.2.2 St22 input

Reading a st2110-22 codestream on "239.168.85.20:20000" with payload_type 112, decode with ffmpeg h264 codec:

```bash
ffmpeg -p_port 0000:af:01.0 -p_sip 192.168.96.2 -p_rx_ip 239.168.85.20 -udp_port 20000 -payload_type 112 -fps 59.94 -video_size 1920x1080 -st22_codec h264 -f mtl_st22 -i "k" -f rawvideo /dev/null -y
```

#### A.2.3 SVT-JPEGXS

Make sure the FFMpeg is build with both MTL and SVT-JPEGXS plugin:

```bash
# start rx
ffmpeg -p_port 0000:af:01.0 -p_sip 192.168.96.2 -p_rx_ip 239.168.85.20 -udp_port 20000 -payload_type 112 -fps 59.94 -video_size 1920x1080 -st22_codec jpegxs -timeout_s 10 -f mtl_st22 -i "k" -vframes 10 -f rawvideo /dev/null -y
# start tx
ffmpeg -stream_loop -1 -video_size 1920x1080 -f rawvideo -pix_fmt yuv420p -i yuv420p_1080p.yuv -filter:v fps=59.94 -c:v libsvt_jpegxs -p_port 0000:af:01.1 -p_sip 192.168.96.3 -p_tx_ip 239.168.85.20 -udp_port 20000 -payload_type 112 -f mtl_st22 -
```

#### A.2.4 SVT-HEVC

Make sure the FFMpeg is build with both MTL and SVT-HEVC plugin:

```bash
# start rx
ffmpeg -p_port 0000:af:01.0 -p_sip 192.168.96.2 -p_rx_ip 239.168.85.20 -udp_port 20000 -payload_type 112 -fps 59.94 -video_size 1920x1080 -st22_codec h265 -timeout_s 10 -f mtl_st22 -i "k" -vframes 10 -f rawvideo /dev/null -y
# start tx
ffmpeg -stream_loop -1 -video_size 1920x1080 -f rawvideo -pix_fmt yuv420p -i yuv420p_1080p.yuv -filter:v fps=59.94 -c:v libsvt_hevc -p_port 0000:af:01.1 -p_sip 192.168.96.3 -p_tx_ip 239.168.85.20 -udp_port 20000 -payload_type 112 -f mtl_st22 -
```

#### A.2.5 St22p support

Another option involves utilizing the MTL built-in ST22 codec plugin, where FFmpeg can directly send or retrieve the YUV raw frame to/from the MTL ST22P plugin. MTL will then internally decode or encode the codec stream.

Reading a st2110-22 pipeline jpegxs codestream on "239.168.85.20:20000" with payload_type 112:

```bash
ffmpeg -p_port 0000:af:01.0 -p_sip 192.168.96.2 -p_rx_ip 239.168.85.20 -udp_port 20000 -payload_type 112 -st22_codec jpegxs -fps 59.94 -pix_fmt yuv422p10le -video_size 1920x1080 -f mtl_st22p -i "k" -f rawvideo /dev/null -y
```

Reading from a yuv file and sending a st2110-22 pipeline jpegxs codestream on "239.168.85.20:20000" with payload_type 112:

```bash
ffmpeg -stream_loop -1 -video_size 1920x1080 -f rawvideo -pix_fmt yuv422p10le -i yuv422p10le_1080p.yuv -filter:v fps=59.94 -p_port 0000:af:01.1 -p_sip 192.168.96.3 -p_tx_ip 239.168.85.20 -udp_port 20000 -payload_type 112 -st22_codec jpegxs -f mtl_st22p -
```

### A.3 ST30P audio run guide

#### A.3.1 St30p input

Reading a st2110-30 stream(pcm24,1ms packet time,2 channels) on "239.168.85.20:30000" with payload_type 111 and encoded to a wav file:

```bash
ffmpeg -p_port 0000:af:01.0 -p_sip 192.168.96.2 -p_rx_ip 239.168.85.20 -udp_port 30000 -payload_type 111 -pcm_fmt pcm24 -at 1ms -ac 2 -f mtl_st30p -i "0" dump.wav -y
```

#### A.3.2 St30p output

Reading from a wav file and sending a st2110-30 stream(pcm24,1ms packet time,2 channels) on "239.168.85.20:30000" with payload_type 111:

```bash
ffmpeg -stream_loop -1 -i test.wav -p_port 0000:af:01.1 -p_sip 192.168.96.3 -p_tx_ip 239.168.85.20 -udp_port 30000 -payload_type 111 -at 1ms -f mtl_st30p -
```

#### A.3.3 St30p pcm16 example

For pcm16 audio, use `mtl_st30p_pcm16` muxer, set `pcm_fmt` to `pcm16` for demuxer.

```bash
ffmpeg -stream_loop -1 -i test.wav -p_port 0000:af:01.1 -p_sip 192.168.96.3 -p_tx_ip 239.168.85.20 -udp_port 30000 -payload_type 111 -at 1ms -f mtl_st30p_pcm16 -

ffmpeg -p_port 0000:af:01.0 -p_sip 192.168.96.2 -p_rx_ip 239.168.85.20 -udp_port 30000 -payload_type 111 -pcm_fmt pcm16 -at 1ms -ac 2 -f mtl_st30p -i "0" dump_pcm16.wav -y
```
