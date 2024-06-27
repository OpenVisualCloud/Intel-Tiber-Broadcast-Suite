# Run guide

> ⚠️ Make sure that all of the hosts used are set up according to the [host setup](build.md).

> **Note:** This instruction regards running the predefined scripts from `pipelines` folder present in the root of the repository. For more information on how to prepare an own pipeline, see:
- [Docker command breakdown](run-know-how.md)
- [FFmpeg Intel® Media Communications Mesh Muxer Parameters Table](plugins/media-communications-mesh.md)
- [Intel® Media Transport Library](plugins/media-transport-library.md)
- [FFmpeg Intel® JPEG XS Parameters Table](plugins/svt-jpeg-xs.md)
- [Raisr FFmpeg Filter Plugin Parameters Table](plugins/video-super-resolution.md)


## Run sample pipelines

The Intel® Tiber™ Broadcast Suite is a package designed for creation of high-performance and high-quality solutions used in live video production.

Video pipelines described below are built using Intel-optimized version of FFmpeg and combine: media transport protocols compliant with SMPTE ST 2110, JPEG XS encoder and decoder, GPU media processing and rendering.

`session A`, `session B` etc. mark separate shell (terminal) sessions. As the Suite is a containerized solution, those sessions can be opened on a single server or multiple servers - on systems connected with each other, after the ports are exposed and IP addresses aligned in pipeline commands.

---

### Multiviewer

Input streams from eight ST 2110-20 cameras are scaled down and composed into a tiled 4x2 multi-view of all inputs on a single frame.

![Multiviewer tile composition](images/multiviewer-process.png)

Scaling and composition are examplary operations that can be replaced by customers with their own visualization apps, for example OpenGL- or Vulcan-based.

Pipeline outputs a single ST 2110 stream.

The example also shows how to use GPU capture to encode a secondary AVC/HEVC stream that can be transmitted with WebRTC for preview.

![Multiviewer](images/multiviewer.png)

Execute a following set of scripts in according terminal sessions to run the Multiviewer pipeline:
```text
session A > multiviewer_tx.sh
session B > multiviewer_process.sh
session C > multiviewer_rx.sh
```


### Recorder

Input streams from ST 2110-20 camera is split to two streams with different resolution 1/4 and 1/16. Scaled outputs are stored on local drive.

![Recorder process](images/recorder-process.png)
![Recorder](images/recorder.png)

Execute a following set of scripts in according terminal sessions to run the Recorder pipeline:
```text
session A > recorder_tx.sh
session B > recorder_rx.sh
```


### Replay

Input streams from two ST 2110-20 camera and are blended together. Blended output is send out via ST 2110 stream.

![Replay process](images/replay-process.png)
![Replay](images/replay.png)

Execute a following set of scripts in according terminal sessions to run the Replay pipeline:
```text
session A > replay_tx.sh
session B > replay_process.sh
session C > replay_rx.sh
```


### Upscale

Input streams from ST 2110-20 camera is scaled up using Video Super Resolution from FullHD to the 4K resolution. Output is send out via ST 2110-20 stream.

![Upscale process](images/upscale-process.png)
![Upscale](images/upscale.png)

Execute a following set of scripts in according terminal sessions to run the Upscale pipeline:
```text
session A > upscale_tx.sh
session B > upscale_process.sh
session C > upscale_rx.sh
```


### JPEG XS

Two input streams from local drive are encoded using JPEG XS codec and send out using ST 2110-22 streams.
Input streams from two ST 2110-22 camera are decoded using JPEG XS codec stored on local drive.

![JPEG XS process](images/jpeg_xs-process.png)
![JPEG XS](images/jpeg_xs.png)

Execute a following set of scripts in according terminal sessions to run the JPEG XS pipeline:
```text
session A > jpeg_xs_tx.sh
session B > jpeg_xs_rx.sh
```

<!-- Temporarily hidden
### Video production pipeline
This pipeline does not have its equivalent in code at the moment, but shows a production-ready solution that could be built using Intel® Tiber™ Broadcast Suite.

![Video production pipeline](images/production-pipeline-example.png)

Two 8K cameras capable of sending ST 2110 stream with video encoded using JPEG XS codec, send their streams using UDP multicast.

Server A receives the streams by two Virtual Functions of Intel® E810 Series Ethernet Adapter card used within a single Intel® Tiber™ Broadcast Suite container. Both streams are decoded with low latency using accelerated SVT JPEG XS on Intel® Xeon® Scalable Processor. One stream is downscaled to 1/4th of the size (to 4K), and the other is downscaled to 1/4th and 1/64th of the size (to 4K and 1080p).

Both 4K streams are sent with the same Virtual Functions they were received with to the next container running on Server B. 1080p stream is also sent to a Recorder/Instant replay machine for archival and replay possibility.

Server B receives three streams, two 4K (close to) real-time ones, and one delayed 1080p stream used for replays. The smallest one is later upscaled with Video Super Resolution on Intel® Data Center GPU Flex Series card to match 4K output.

All of the streams are blended and mixed based on predefined instructions. The output is then compressed and sent using RTP protocol (TCP) as a 4K stream.
-->