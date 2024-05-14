  # -v $(pwd)/clips:/clips \
docker run --rm -it \
  --privileged \
  --device=/dev/dri:/dev/dri \
  -v /clips/y210le:/clips \
  -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
  video_production_image \
  -y \
  -v debug \
  -qsv_device /dev/dri/renderD128 \
  -an \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 4096x1744 -i clips/out200.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 4096x1744 -i clips/out300.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 4096x1744 -i clips/out400.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 4096x1744 -i clips/out500.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 4096x1744 -i clips/out600.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 4096x1744 -i clips/out700.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 4096x1744 -i clips/out800.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 4096x1744 -i clips/out900.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 4096x1744 -i clips/out1000.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 4096x1744 -i clips/out1100.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 4096x1744 -i clips/out1200.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 4096x1744 -i clips/out1300.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 4096x1744 -i clips/out1400.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 4096x1744 -i clips/out1500.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 4096x1744 -i clips/out1600.yuv \
  -hwaccel qsv -hwaccel_output_format qsv -f rawvideo -pix_fmt y210le -s:v 4096x1744 -i clips/out1700.yuv \
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
      layout=0_0|0_436|0_872|0_1308|\
             1024_0|1024_436|1024_872|1024_1308|\
             2048_0|2048_436|2048_872|2048_1308|\
             3072_0|3072_436|3072_872|3072_1308[out];[out]hwdownload,format=y210[multiview]" \
  -map "[multiview]" -f rawvideo -pix_fmt y210le clips/out_y210le.yuv
