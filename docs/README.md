# Video Production Pipeline

## 1. Overview

The Intel® Video Production Pipeline is a software-based package designed for creation of high-performance and high-quality solutions used in live video production.
The video pipelines are built using Intel-optimized version of FFmpeg and combine: media transport protocols (SMPTE ST 2110-compliant), JPEG-XS encoding/decoding, GPU media processing and rendering.

The diagram below illustrates an example pipeline created with this software package:
![Multiviewer](/docs/png/multiviewer.png)

## 2. Software Architecture

The Intel® Video Production Pipeline uses open-source FFmpeg framework as a baseline, and enhances it with:
- Intel® Media Transport Library (MTL) with SMPTE 2110 transport protocols and yuv422p10le and y210le pixel formats.
- Intel® QSV and OneVPL libraries to support hardware-accelerated media processing with Intel Flex GPU cards.
- DPC++ kernels to enable custom effect filters used in video production (not supported in this release).
- OpenGL/Vulcan integration to display rendering effects (not supported in this release).

The software package includes several performance features on to of regular Intel® FFMpeg-Cartwheel releases:
- memory management optimizations for page-aligned surface allocations
- asynchronous execution of video pipeline filters to maximize GPU utilization
- high-throughput GPU-CPU memory data transfers

![Architecture](/docs/png/architecture.png)

## 3. Build Instructions

Step 1. Please install required MTL packages on the host machine:
[MTL build instruction](https://github.com/OpenVisualCloud/Media-Transport-Library/blob/main/doc/build.md)

Step 2. Create a docker image shared by all video pipelines:

```
docker build -t video_production_image -f Dockerfile .
```

Step 3. Setup MTL environment running the following command on the host machine as a root user:

```
./first_run.sh
```

Step 3 is required each time host is restarted and MTL is needed.

Step 4. Run .sh pipeline scripts to execute video pipelines. Reference pipelines can be found in [pipelines](./pipelines) directory.


