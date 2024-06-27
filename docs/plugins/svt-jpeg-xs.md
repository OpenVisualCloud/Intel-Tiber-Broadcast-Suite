# Intel® JPEG XS Codec

FFmpeg Intel® JPEG XS Codec Plugin Documentation

> 💡 _**Tip:** For up to date documentation refer to: [https://github.com/OpenVisualCloud/SVT-JPEG-XS](https://github.com/OpenVisualCloud/SVT-JPEG-XS)_

## FFmpeg Intel® JPEG XS Parameters Table

The FFmpeg JPEG XS codec plugin provides encoding and decoding functionality for the JPEG XS format. Below are the parameters that can be configured for both the libsvtjpegxs encoder and decoder. To use the plugin, specify the codec with `-codec jpegxs` for encoding or `-codec:v jpegxs` for decoding.

### Encoder available params

Name | mandatory/optional | Accepted values | description
  --         |     --    |                           --                                     |  --
bpp          | mandatory | any integer/float greater than 0 (example: 0.5, 3, 3.75, 5 etc.) | Bits Per Pixel
decomp_v     | optional  | 0, 1, 2(default)                                                 | Number of Vertical decompositions
decomp_h     | optional  | 0, 1, 2, 3, 4, 5(default)                                        | Number of Horizontal decompositions, have to be greater or equal to decomp_v
threads      | optional  | Any integer in range< 1;64>                                      | Number of threads encoder can create
slice_height | optional  | (default:16), Any integer in range <1;source_height>, also it have to be multiple of 2^(decomp_v) | Coding feature: Specify slice height in units of picture luma pixels
quantization | optional  | (default:deadzone), 0(deadzone), 1(uniform)                    | Coding feature: Quantization method
coding-signs | optional  | (default:off), 0(off), 1(fast), 2(full)                        | Coding feature: Sign handling strategy
coding-sigf  | optional  | (default:on), 0(off), 1(on)                                    | Coding feature: Significance coding
coding-vpred | optional  | (default:off), 0(off), 1(on)                                   | Coding feature: Vertical-prediction

### Decoder available params

Name | mandatory/optional | Accepted values | description
  --     |     --    |               --                                                | --
threads  | optional  | Any integer in range< 1;64>                                     | Number of threads decoder can create

## FFmpeg Intel® JPEG XS Encoder Parameters

### Input Parameters

#### Mandatory Arguments
- **`bpp`**: Bits Per Pixel.
  - **Type**: Float/Integer
  - **Default**: `N/A`
  - **Flags**: Encoding parameter
  - **Range**: `> 0`

#### Optional Arguments
- **`decomp_v`**: Number of Vertical decompositions.
  - **Type**: Integer
  - **Default**: `2`
  - **Flags**: Encoding parameter
  - **Range**: `0`, `1`, `2`

- **`decomp_h`**: Number of Horizontal decompositions.
  - **Type**: Integer
  - **Default**: `5`
  - **Flags**: Encoding parameter
  - **Range**: `0` to `5`

- **`threads`**: Number of threads encoder can create.
  - **Type**: Integer
  - **Default**: `N/A`
  - **Flags**: Encoding parameter
  - **Range**: `1` to `64`

- **`slice_height`**: Coding feature: Specify slice height in units of picture luma pixels.
  - **Type**: Integer
  - **Default**: `16`
  - **Flags**: Encoding parameter
  - **Range**: `1` to `source_height`
  - **Note**: Must be a multiple of `2^(decomp_v)`

- **`quantization`**: Coding feature: Quantization method.
  - **Type**: Enum
  - **Default**: `deadzone`
  - **Flags**: Encoding parameter
  - **Possible Values**: `0` (deadzone), `1` (uniform)

- **`coding-signs`**: Coding feature: Sign handling strategy.
  - **Type**: Enum
  - **Default**: `off`
  - **Flags**: Encoding parameter
  - **Possible Values**: `0` (off), `1` (fast), `2` (full)

- **`coding-sigf`**: Coding feature: Significance coding.
  - **Type**: Enum
  - **Default**: `on`
  - **Flags**: Encoding parameter
  - **Possible Values**: `0` (off), `1` (on)

- **`coding-vpred`**: Coding feature: Vertical-prediction.
  - **Type**: Enum
  - **Default**: `off`
  - **Flags**: Encoding parameter
  - **Possible Values**: `0` (off), `1` (on)

## FFmpeg Intel® JPEG XS Decoder Parameters

### Input Parameters

#### Optional Arguments
- **`threads`**: Number of threads decoder can create.
  - **Type**: Integer
  - **Default**: `N/A`
  - **Flags**: Decoding parameter
  - **Range**: `1` to `64`

### Example Usage

#### Encoding raw video

```bash
ffmpeg -y -s:v 1920x1080 -c:v rawvideo -pix_fmt yuv420p -i raw_stream.yuv -codec jpegxs -bpp 1.25 -threads 5 encoded_file.mov
```

#### Decoding JPEG XS streams to raw video

```bash
ffmpeg -threads 10 -i jpegxs-file.mov output.yuv
```

#### Transcoding from any format to JPEG XS

```bash
ffmpeg -i input.mov -c:v jpegxs -bpp 2 -threads 15 encoder.mov
```
