```
The start_processing() function handles the main processing loop:
```

- Opens input files using the `mtl_st20p` plugin and sets the necessary options.
- Finds stream information.
- Sets up decoders.
- Sets up the filter graph.
- Sets up the output format context using the `mtl_st20p` plugin and sets the necessary options.
- Creates the output stream and sets up the encoder.
- Opens the output file.
- Writes the output file header.
- Allocates frames and packets.
- Reads, decodes, filters, encodes, and writes frames.
- Writes the output file trailer.

# -filter_complex:
This option specifies a complex filter graph. It allows you to define a series of filters and how they are connected, including multiple inputs and outputs.

```
The setup_filter_graph() function sets up the filter graph, including buffer source filters, hwupload filters, scale_qsv filters, xstack_qsv filter, and format filters.
```
command from [multiviewer_process.sh](https://github.com/OpenVisualCloud/Intel-Tiber-Broadcast-Suite/blob/main/pipelines/multiviewer_process.sh) : 

```
-filter_complex "[0:v]hwupload,scale_qsv=iw/4:ih/2[out0]; \
                       [1:v]hwupload,scale_qsv=iw/4:ih/2[out1]; \
                       [2:v]hwupload,scale_qsv=iw/4:ih/2[out2]; \
                       [3:v]hwupload,scale_qsv=iw/4:ih/2[out3]; \
                       [4:v]hwupload,scale_qsv=iw/4:ih/2[out4]; \
                       [5:v]hwupload,scale_qsv=iw/4:ih/2[out5]; \
                       [6:v]hwupload,scale_qsv=iw/4:ih/2[out6]; \
                       [7:v]hwupload,scale_qsv=iw/4:ih/2[out7]; \
                       [out0][out1][out2][out3] \
                       [out4][out5][out6][out7] \
                       xstack_qsv=inputs=8:\
                       layout=0_0|w0_0|0_h0|w0_h0|w0+w1_0|w0+w1+w2_0|w0+w1_h0|w0+w1+w2_h0, \
                       format=y210le,format=yuv422p10le" \
```

**Input Streams:**

- `[0:v], [1:v], [2:v], [3:v], [4:v], [5:v], [6:v], [7:v]:`
  These are the video streams from the input files. The numbers (0, 1, 2, etc.) refer to the input file indices, and v indicates that these are video streams.

**Filters:**

**hwupload:**

  - This filter uploads the video frames to the GPU for hardware acceleration. It is used to prepare the frames for further processing by hardware-accelerated filters.

**scale_qsv=iw/4:ih/2:**

  - This filter scales the video frames using Intel's Quick Sync Video (QSV) hardware acceleration. The iw/4 and ih/2 specify the new width and height of the frames, which are one-fourth and one-half of the original width and height, respectively.

**Output Labels:**

  - `[out0], [out1], [out2], [out3], [out4], [out5], [out6], [out7]:`
    These labels are used to name the outputs of the scale_qsv filters. They are used as inputs to the next filter in the chain.

**Stacking Filter:**

  - `xstack_qsv=inputs=8:layout=0_0|w0_0|0_h0|w0_h0|w0+w1_0|w0+w1+w2_0|w0+w1_h0|w0+w1+w2_h0:`
    This filter stacks multiple video frames together using Intel's QSV hardware acceleration. The inputs=8 specifies that there are 8 input streams. The layout parameter defines the layout of the stacked frames. The layout 
  
**positions the frames in a grid-like pattern:**
  - `0_0:` The first frame is placed at the top-left corner.
  - `w0_0:` The second frame is placed to the right of the first frame.
  - `0_h0:` The third frame is placed below the first frame.
  - `w0_h0:` The fourth frame is placed to the right of the third frame.
  - `w0+w1_0:` The fifth frame is placed to the right of the second frame.
  - `w0+w1+w2_0:` The sixth frame is placed to the right of the fifth frame.
  - `w0+w1_h0:` The seventh frame is placed below the fifth frame.
  - `w0+w1+w2_h0:` The eighth frame is placed to the right of the seventh frame.

**Format Conversion:**
  - `format=y210le:`
    This filter converts the pixel format of the video frames to y210le, which is a 10-bit YUV 4:2:2 format with little-endian byte order.

  - `format=yuv422p10le:`
    This filter converts the pixel format of the video frames to yuv422p10le, which is another 10-bit YUV 4:2:2 format with little-endian byte order.