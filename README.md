# Video Production Pipeline

## 1. Overview

The Intel® Video Production Pipeline is a software based solution designed for creation of video processing pipeliens useed in live video production. 
The pipelines are built using Intel-optimized version of FFMpeg and combine: media transport protocols (SMPTE ST 2110-compliant), JPEG-XS encoding/decoding, GPU media processing and rendering. 

The diagram below illustrates an example pipeline created with this software package:  
[alt text](https://github.com/intel-innersource/applications.services.cloud.visualcloud.vcdp.video-production-pipeline/blob/main/doc/png/multiviewer.png)

## 2. Software Architecture 

The Intel® Video Production Pipeline uses open-source FFMpeg framework as a baseline, and enhances it with: 
- Intel® Media Transport Library (IMTL) with SMPTE 2110 transport protocols and yuv422p10le and y210le pixel formats. 
- Intel® QSV and OneVPL libraries to support hardware-accelerated media processing with Intel Flex GPU cards. 
- DPC++ kernels with custom color-space-conversion filters suitable for video production (not supported in this release).
- OpenGL/Vulkan integration to allow custom rendering effects (not supported in this release).

The software package includes performance features on to of regular Intel® FFMpeg-Carthweel releases:
- memory management optimizations for page-aligned surface allocations
- asynchronous execution of video pipeline filters to maximise GPU utilization
- high-throuthput GPU-CPU memory copy 

[Architecture](https://github.com/intel-innersource/applications.services.cloud.visualcloud.vcdp.video-production-pipeline/blob/main/doc/png/architecture.png)

## 3. Build Instructions 

Before using IMTL in Docker please install required packages on host machine:
[Build_instruction](https://github.com/OpenVisualCloud/Media-Transport-Library/blob/main/doc/build.md)

To run ffmpeg in Docker please run command which creates docker image:

```
docker build -t my_ffmpeg .
```

Step 1. If IMTL plugin support is needed then please run commands on host as a root:

```
./first_run.sh
```

Steps 1 is required each time host is restarted and IMTL is needed.

Step 2. Run .sh script with ffmpeg parameters. Examples are in [test_scripts](./test_scripts) directory.


