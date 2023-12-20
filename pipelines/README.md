# Video Production Pipelines

The following reference pipelines are delivered as part of this project. 

## 1. Multiviewer pipeline 

Input streams from multiple ST2110 cameras are scaled down and composed into a tiled multi-view of all inputs on a single frame. Scaling and composition are example operations that will be replaced by customers with their visualization apps. Majority of customers use OpenGL as a visualization app. The live video streaming solution should allow integrating either OpenGL  or Vulkan applications, with OpenGL being target for example pipeline.These apps already exist and use OpenGL which we must also support. Pipeline output is a single ST2110 stream. The example also shows how to use GPU capture to encode a secondary AVC/HEVC stream that can be transmitted with WebRTC for preview. 

Note that both input and output are expected to use Y’CRCB 10b BT2110 HLG with BT2020 Colour Gamut, but there might be optional color space conversions to improve performance or for specific customer requirements. 

For development and integration ease, input/output to support either input from file or from network. 

There are separate docker files for camera and multiviewer pipelines. The pipelines should be run on separate nodes for performance testing. Alternatively, one can launch both pipelines on a single node using the following command line: 
```
docker_camera.sh & docker_multiviewer.sh
```

![Multiviewer](https://github.com/intel-innersource/applications.services.cloud.visualcloud.vcdp.video-production-pipeline/blob/main/doc/png/multiviewer.png)

Known restrictions in the first release:
- 4 input streams are supported for performance testing
- simultaneous receive and transmit is not supported in IMTL plugin
- no integration with visualization app

## 2. Replay pipeline 1

Input stream from an ST2110 camera is scaled down, tone mapped, encoded to AVC/HEVC and stored in a file for future use. There are two different resolutions stored. For performance it is allowed to cascade scale operations (scaling to 270p can use 540p as input instead of the original stream). Note that input is expected to use Y’CRCB 10b BT2110 HLG with BT2020 Colour Gamut. Output could use either 422 or 420 sampling. 

![Replay 1](https://github.com/intel-innersource/applications.services.cloud.visualcloud.vcdp.video-production-pipeline/blob/main/doc/png/replay1.png)

Known restrictions in the first release:
- custom 3D/LUT not available

## 3. Replay pipeline 2

Input stream is decoded from AVC/HEVC stored in a file and blended with a live input stream coming from an ST2110 camera. Note that live input is expected to use Y’CRCB 10b BT2110 HLG with BT2020 Colour Gamut. Replay input could use either 422 or 420 sampling. 

![Replay 2](https://github.com/intel-innersource/applications.services.cloud.visualcloud.vcdp.video-production-pipeline/blob/main/doc/png/replay2.png)

Known restrictions in the first release:
- custom blend filter not available
