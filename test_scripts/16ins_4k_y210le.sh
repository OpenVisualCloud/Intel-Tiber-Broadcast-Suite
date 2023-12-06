docker run --rm -it \
  --privileged \
  --device=/dev/dri:/dev/dri \
  -v $(pwd)/clips:/clips \
  -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
  my_ffmpeg \
  -y \
  -qsv_device /dev/dri/renderD128 \
  -an \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 3840x2160 -i clips/random_y210le.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 3840x2160 -i clips/random_y210le_1.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 3840x2160 -i clips/random_y210le_2.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 3840x2160 -i clips/random_y210le_3.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 3840x2160 -i clips/random_y210le_4.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 3840x2160 -i clips/random_y210le_5.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 3840x2160 -i clips/random_y210le_6.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 3840x2160 -i clips/random_y210le_7.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 3840x2160 -i clips/random_y210le_8.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 3840x2160 -i clips/random_y210le_9.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 3840x2160 -i clips/random_y210le_10.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 3840x2160 -i clips/random_y210le_11.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 3840x2160 -i clips/random_y210le_12.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 3840x2160 -i clips/random_y210le_13.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 3840x2160 -i clips/random_y210le_14.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 3840x2160 -i clips/random_y210le_15.yuv \
  -noauto_conversion_filters -filter_complex " \
    [0:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4:mode=compute[tile0]; \
    [1:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4:mode=compute[tile1]; \
    [2:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4:mode=compute[tile2]; \
    [3:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4:mode=compute[tile3]; \
    [4:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4:mode=compute[tile4]; \
    [5:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4:mode=compute[tile5]; \
    [6:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4:mode=compute[tile6]; \
    [7:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4:mode=compute[tile7]; \
    [8:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4:mode=compute[tile8]; \
    [9:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4:mode=compute[tile9]; \
    [10:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4:mode=compute[tile10]; \
    [11:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4:mode=compute[tile11]; \
    [12:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4:mode=compute[tile12]; \
    [13:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4:mode=compute[tile13]; \
    [14:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4:mode=compute[tile14]; \
    [15:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4:mode=compute[tile15]; \
    [tile0][tile1][tile2][tile3] \
    [tile4][tile5][tile6][tile7] \
    [tile8][tile9][tile10][tile11] \
    [tile12][tile13][tile14][tile15] xstack_qsv=inputs=16:\
      layout=0_0|0_540|0_1080|0_1620|\
             960_0|960_540|960_1080|960_1620|\
             1920_0|1920_540|1920_1080|1920_1620|\
             2880_0|2880_540|2880_1080|2880_1620[out];[out]hwdownload,format=y210[multiview]" \
  -map "[multiview]" -f rawvideo -pix_fmt y210le clips/out_y210le.yuv
