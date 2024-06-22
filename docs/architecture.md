# Architecture

## 1. Overview
The Intel® Video Production Pipeline is a software-based package designed for creation of high-performance and high-quality solutions used in live video production. The video pipelines are built using Intel-optimized version of FFmpeg and combine: media transport protocols (SMPTE ST 2110-compliant), JPEG-XS encoding/decoding, GPU media processing and rendering.

## 2. Software Architecture

The Intel® Video Production Pipeline uses open-source FFmpeg framework as a baseline, and enhances it with:
- Intel® Media Transport Library (MTL) with SMPTE 2110 transport protocols and yuv422p10le and y210le pixel formats.
- Intel® QSV and OneVPL libraries to support hardware-accelerated media processing with Intel Flex GPU cards.
- DPC++ kernels to enable custom effect filters used in video production (not supported in this release).
- OpenGL/Vulkan integration to display rendering effects (not supported in this release).

The software package includes several performance features on to of regular Intel® FFMpeg-Carthweel releases:
- memory management optimizations for page-aligned surface allocations
- asynchronous execution of video pipeline filters to maximize GPU utilization
- high-throughput GPU-CPU memory data transfers

![Architecture](images/architecture.png)


## 3. Components inside SDB docker image

Component               |   Vsrsion     |   Source
---                     |   ---         |   ---
FFmpeg                  |   6.1.1       |   [FFmpeg ](https://github.com/FFmpeg/FFmpeg)
Intel® FFmpeg patches   |   6.1         |   [Intel® FFmpeg patches](https://github.com/intel/cartwheel-ffmpeg)
Media Transport Library |   commitID: b210f1a85f571507f317d156b105dbe5690a234d   |   [Media Transport Library](https://github.com/OpenVisualCloud/Media-Transport-Library)
Media Communications Mesh| __TBD__      |   [Media Communications Mesh](https://github.com/OpenVisualCloud/Media-Communications-Mesh)
Data Plane Development Kit (DPDK)   |    23.11   |   [DPDK](https://github.com/DPDK/dpdk)
SVT JPEG-XS             |  0.9      |   [SVT JPEG-XS](https://github.com/OpenVisualCloud/SVT-JPEG-XS)
SVT AV1                 |  1.7.0        |   [SVT AV1](https://gitlab.com/AOMediaCodec/SVT-AV1)
Intel® Integrated Performance Primitives    |  2021.10.1.16    |	[IPP](https://www.intel.com/content/www/us/en/developer/articles/tool/oneapi-standalone-components.html#ipp)
Video Super Resolution  |   23.11       |   [VSR](https://github.com/OpenVisualCloud/Video-Super-Resolution-Library)
VMAF                    |   2.3.1       |   [VMAF](https://github.com/Netflix/vmaf)
oneVPL                  |   23.3.4      |   [oneVPL](https://github.com/intel/vpl-gpu-rt)
LIBVPL                  |   2023.3.1    |   [LIBVPL](https://github.com/intel/libvpl)
Intel Media Driver (IHD)|   23.3.5      |   [IHD](https://github.com/intel/media-driver)
GMMLIB                  |   22.3.12     |   [GMMLIB](https://github.com/intel/gmmlib)
Vulkan-Headers          |   1.3.268.0   |   [Vulkan](https://github.com/KhronosGroup/Vulkan-Headers)
LIBVA                   |   2.20.0      |   [LIBVA](https://github.com/intel/libva)
