# Intel® Video Super Resolution

FFmpeg Intel® Video Super Resolution Filter Plugin Documentation

> 💡 _**Tip:** For up to date documentation refer to: [https://github.com/OpenVisualCloud/Video-Super-Resolution-Library](https://github.com/OpenVisualCloud/Video-Super-Resolution-Library)_

## Raisr FFmpeg Filter Plugin

The Intel® Video Super Resolution (Raisr) filter plugin for FFmpeg provides super-resolution capabilities for video frames using the Raisr algorithm. It allows for upscaling video frames by a specified ratio with options for bit depth, color range, threading, and more. Below are the input parameters that can be configured for the Raisr filter plugin. To start using Raisr filter, use `-vf "raisr=<<PARAMS>>"`

### Raisr FFmpeg Filter Plugin Parameters Table


| Parameter     | Description                                      | Type         | Default                | Range / Possible Values                          |
|---------------|--------------------------------------------------|--------------|------------------------|--------------------------------------------------|
| `ratio`       | Set the ratio of the upscaling, between 1 and 2. | Float        | `DEFAULT_RATIO`        | `MIN_RATIO` to `MAX_RATIO`                       |
| `bits`        | Set the bit depth for processing.                | Integer      | `8`                    | `8` to `10`                                      |
| `range`       | Set the input color range.                       | String       | `"video"`              | `"video"`, `"full"`                              |
| `threadcount` | Set the number of threads to use.                | Integer      | `DEFAULT_THREADCOUNT`  | `MIN_THREADCOUNT` to `MAX_THREADCOUNT`           |
| `filterfolder`| Set the absolute filter folder path.             | String       | `"filters_2x/filters_lowres"` | N/A                                      |
| `blending`    | Set the CT blending mode.                        | Integer      | `BLENDING_COUNT_OF_BITS_CHANGED` | `BLENDING_RANDOMNESS` to `BLENDING_COUNT_OF_BITS_CHANGED` |
| `passes`      | Set the number of passes to run.                 | Integer      | `1`                    | `1` to `2`                                       |
| `mode`        | Set the mode for two-pass processing.            | Integer      | `1`                    | `1` (upscale in 1st pass), `2` (upscale in 2nd pass) |
| `asm`         | Set the x86 asm type.                            | String       | `"avx512fp16"`         | `"avx512fp16"`, `"avx512"`, `"avx2"`, `"opencl"` |
| `platform`    | Select the OpenCL platform ID.                   | Integer      | `0`                    | `0` to `INT_MAX`                                 |
| `device`      | Select the OpenCL device ID.                     | Integer      | `0`                    | `0` to `INT_MAX`                                 |

## Raisr OpenCL FFmpeg Filter Plugin

The Raisr OpenCL filter plugin for FFmpeg is designed to perform super-resolution on video frames using the Raisr algorithm with OpenCL acceleration. This document outlines the input parameters and configuration options available for the Raisr OpenCL filter plugin. To start using Raisr OpenCL filter, use `-vf "raisr_opencl=<<PARAMS>>"`

### Raisr OpenCL FFmpeg Filter Plugin Parameters Table

| Parameter      | Description                                                  | Type    | Default Value                   | Range / Possible Values       |
|----------------|--------------------------------------------------------------|---------|---------------------------------|-------------------------------|
| `ratio`        | Set the ratio of the upscaling, between 1 and 2.             | Float   | `DEFAULT_RATIO`                 | `MIN_RATIO` to `MAX_RATIO`    |
| `bits`         | Set the bit depth for processing.                            | Integer | `8`                             | `8` to `10`                   |
| `range`        | Set the input color range.                                   | Integer | `VideoRange`                    | `VideoRange`, `FullRange`     |
| `filterfolder` | Set the absolute filter folder path.                         | String  | `"filters_2x/filters_lowres"`   | N/A                           |
| `blending`     | Set the CT blending mode.                                    | Integer | `CountOfBitsChanged`            | `Randomness`, `CountOfBitsChanged` |
| `passes`       | Set the number of passes to run.                             | Integer | `1`                             | `1` to `2`                    |
| `mode`         | Set the mode for two-pass processing.                        | Integer | `1`                             | `1`, `2`                      |

## Raisr FFmpeg Filter Plugin Documentation

### Input Parameters

#### `ratio`
- **Description**: Set the ratio of the upscaling, between 1 and 2.
- **Type**: Float
- **Default**: `DEFAULT_RATIO`
- **Range**: `MIN_RATIO` to `MAX_RATIO`
- **Flags**: Filtering parameter, Video parameter

#### `bits`
- **Description**: Set the bit depth for processing.
- **Type**: Integer
- **Default**: `8`
- **Range**: `8` to `10`
- **Flags**: Filtering parameter, Video parameter

#### `range`
- **Description**: Set the input color range.
- **Type**: String
- **Default**: `"video"`
- **Flags**: Filtering parameter, Video parameter
- **Possible Values**: `"video"`, `"full"`

#### `threadcount`
- **Description**: Set the number of threads to use.
- **Type**: Integer
- **Default**: `DEFAULT_THREADCOUNT`
- **Range**: `MIN_THREADCOUNT` to `MAX_THREADCOUNT`
- **Flags**: Filtering parameter, Video parameter

#### `filterfolder`
- **Description**: Set the absolute filter folder path.
- **Type**: String
- **Default**: `"filters_2x/filters_lowres"`
- **Flags**: Filtering parameter, Video parameter

#### `blending`
- **Description**: Set the CT blending mode.
- **Type**: Integer
- **Default**: `BLENDING_COUNT_OF_BITS_CHANGED`
- **Range**: `BLENDING_RANDOMNESS` to `BLENDING_COUNT_OF_BITS_CHANGED`
- **Flags**: Filtering parameter, Video parameter
- **Possible Values**: `1` (Randomness), `2` (CountOfBitsChanged)

#### `passes`
- **Description**: Set the number of passes to run.
- **Type**: Integer
- **Default**: `1`
- **Range**: `1` to `2`
- **Flags**: Filtering parameter, Video parameter

#### `mode`
- **Description**: Set the mode for two-pass processing.
- **Type**: Integer
- **Default**: `1`
- **Range**: `1` to `2`
- **Flags**: Filtering parameter, Video parameter
- **Possible Values**: `1` (upscale in 1st pass), `2` (upscale in 2nd pass)

#### `asm`
- **Description**: Set the x86 asm type.
- **Type**: String
- **Default**: `"avx512fp16"`
- **Flags**: Filtering parameter, Video parameter
- **Possible Values**: `"avx512fp16"`, `"avx512"`, `"avx2"`, `"opencl"`

#### `platform`
- **Description**: Select the OpenCL platform ID.
- **Type**: Integer
- **Default**: `0`
- **Range**: `0` to `INT_MAX`
- **Flags**: Filtering parameter, Video parameter

#### `device`
- **Description**: Select the OpenCL device ID.
- **Type**: Integer
- **Default**: `0`
- **Range**: `0` to `INT_MAX`
- **Flags**: Filtering parameter, Video parameter

### Usage Example

```sh
ffmpeg -i input.mp4 -vf "raisr=ratio=2:bits=10:range=full" output.mp4
```

## Raisr OpenCL FFmpeg Filter Plugin Documentation

### Input Parameters

#### `ratio`
- **Description**: Set the ratio of the upscaling, between 1 and 2.
- **Type**: Float
- **Default**: `DEFAULT_RATIO`
- **Range**: `MIN_RATIO` to `MAX_RATIO`

#### `bits`
- **Description**: Set the bit depth for processing.
- **Type**: Integer
- **Default**: `8`
- **Range**: `8` to `10`

#### `range`
- **Description**: Set the input color range.
- **Type**: Integer
- **Default**: `VideoRange`
- **Possible Values**: `VideoRange`, `FullRange`

#### `filterfolder`
- **Description**: Set the absolute filter folder path.
- **Type**: String
- **Default**: `"filters_2x/filters_lowres"`

#### `blending`
- **Description**: Set the CT blending mode.
- **Type**: Integer
- **Default**: `CountOfBitsChanged`
- **Possible Values**: `Randomness`, `CountOfBitsChanged`

#### `passes`
- **Description**: Set the number of passes to run.
- **Type**: Integer
- **Default**: `1`
- **Range**: `1` to `2`

#### `mode`
- **Description**: Set the mode for two-pass processing.
- **Type**: Integer
- **Default**: `1`
- **Range**: `1` to `2`
- **Possible Values**: `1` (upscale in 1st pass), `2` (upscale in 2nd pass)

### Constants

- `DEFAULT_RATIO`: The default upscaling ratio.
- `MIN_RATIO`: The minimum upscaling ratio.
- `MAX_RATIO`: The maximum upscaling ratio.
- `VideoRange`: Represents the video color range.
- `FullRange`: Represents the full color range.
- `Randomness`: Represents the randomness blending mode.
- `CountOfBitsChanged`: Represents the count of bits changed blending mode.

### Usage Example

```sh
ffmpeg -init_hw_device vaapi=va -init_hw_device opencl@va -i input.mp4 -vf "format=yuv420p,hwupload,raisr_opencl,hwdownload,format=yuv420p,format=yuv422p10le" -c:v h264 output.mp4
```
