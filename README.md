# Video Production Pipeline

## 1. Overview

The Intel® Video Production Pipeline is a software-based package designed for creation of high-perofrmance and high-quality solutions used in live video production. 
The video pipelines are built using Intel-optimized version of FFmpeg and combine: media transport protocols (SMPTE ST 2110-compliant), JPEG-XS encoding/decoding, GPU media processing and rendering. 

The diagram below illustrates an example pipeline created with this software package:  
![Multiviewer](https://github.com/intel-innersource/applications.services.cloud.visualcloud.vcdp.video-production-pipeline/blob/main/doc/png/multiviewer.png)

## 2. Software Architecture 

The Intel® Video Production Pipeline uses open-source FFmpeg framework as a baseline, and enhances it with: 
- Intel® Media Transport Library (IMTL) with SMPTE 2110 transport protocols and yuv422p10le and y210le pixel formats. 
- Intel® QSV and OneVPL libraries to support hardware-accelerated media processing with Intel Flex GPU cards. 
- DPC++ kernels to enable custom effect filters used in video production (not supported in this release).
- OpenGL/Vulkan integration to display rendering effects (not supported in this release).

The software package includes several performance features on to of regular Intel® FFMpeg-Carthweel releases:
- memory management optimizations for page-aligned surface allocations
- asynchronous execution of video pipeline filters to maximise GPU utilization
- high-throuthput GPU-CPU memory data transfers 

![Architecture](https://github.com/intel-innersource/applications.services.cloud.visualcloud.vcdp.video-production-pipeline/blob/main/doc/png/architecture.png)

## 3. Build Instructions 

Step 1. Please install required IMTL packages on the host machine:
[IMTL build instruction](https://github.com/OpenVisualCloud/Media-Transport-Library/blob/main/doc/build.md)

Step 2. Pleaes create a docker image shared by all video pipelines:

```
docker build -t my_ffmpeg .
```

Step 3. Please setup IMTL environment running the following command on the host machine as a root user:

```
./first_run.sh
```

Step 3 is required each time host is restarted and IMTL is needed.

Step 4. Run .sh pipeline scripts to execute video pipelines. Reference pipeliines can be found in [pipelines](./pipelines) directory.


