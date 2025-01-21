# Summary of Changes for Intel® Tiber™ Broadcast Suite - 24.10 Release: 

What's Changed
- Support Nvidia GPU (CUDA filters)
- Uploaded Intel® Tiber™ Broadcast Suite to the Docker Hub
- Enabled build of Intel® Tiber™ Broadcast Suite from Debian packages
- Updated supported FFmpeg version to 7.0
- Updated documentation
- Updated build.sh


# Summary of Changes for Intel® Tiber™ Broadcast Suite - 24.07 Release: 

What's Changed
- Move the lcore management to mtl-manager
- Remove unnecessary MTL patch and adjust tag
- Minor update to build.md


# Summary of Changes for Intel® Tiber™ Broadcast Suite - Initial Release: 

## New Features:

- Comprehensive Integration with FFmpeg: Enhanced version of FFmpeg tailored with Intel's patches for optimized media handling. 
- Advanced Media Processing: Integration of Intel® QSV and OneVPL for GPU-accelerated processing. 
- SMPTE ST 2110 Compliance: Full support for media transport protocols ensuring high compatibility and performance in professional environments. 
- JPEG XS Support: Includes encoding and decoding capabilities for JPEG XS, optimizing bandwidth and storage. 
- GPU Media Rendering: Utilizes Intel Flex GPU cards for efficient media rendering tasks. 
- Intel® QSV and OneVPL libraries: These libraries support hardware-accelerated media processing with Intel Flex GPU cards, enhancing the performance of the Suite. 
- DPC++ kernels: Although not supported in this release, the Suite is designed to enable custom effect filters used in video production in future releases. 
- OpenGL/Vulcan feature integration. 
- Performance features: The Suite includes several performance features such as memory management optimizations for page-aligned surface allocations, asynchronous execution of video pipeline filters to maximize GPU utilization, and high-throughput GPU-CPU memory data transfers. 
- Mesh Plugin: This plugin allows single or multiple instances of FFmpeg with Mesh Plugin to connect to selected Media Proxy instance, enhancing the efficiency of the Suite.
- Support for PTP Time synchronization: This feature uses Media Transport Library PTP Time synchronization feature, enhancing the synchronization of the Suite. 
- Support for changing input and output streams Payload ID: This feature adds flexibility to the deployment of the Suite. 


## Enhancements:

- Optimized Memory Management: Improvements in memory allocation that enhance performance and reduce latency. 
- Asynchronous Execution: Video pipeline filters operate asynchronously to maximize GPU utilization. 
- High-Throughput Transfers: Enhanced GPU-CPU memory data transfers for high-performance media streaming. 
