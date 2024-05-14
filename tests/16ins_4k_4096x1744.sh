  # -v $(pwd):/config \
docker run --rm -it \
  --privileged \
  --device=/dev/dri:/dev/dri \
  -v /clips/mkv:/clips \
  -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
  video_production_image \
  -y \
  -qsv_device /dev/dri/renderD128 \
  -an \
  -hwaccel qsv -i /clips/test0_4k.mkv \
  -hwaccel qsv -i /clips/test1_4k.mkv \
  -hwaccel qsv -i /clips/test2_4k.mkv \
  -hwaccel qsv -i /clips/test3_4k.mkv \
  -hwaccel qsv -i /clips/test4_4k.mkv \
  -hwaccel qsv -i /clips/test5_4k.mkv \
  -hwaccel qsv -i /clips/test6_4k.mkv \
  -hwaccel qsv -i /clips/test7_4k.mkv \
  -hwaccel qsv -i /clips/test8_4k.mkv \
  -hwaccel qsv -i /clips/test9_4k.mkv \
  -hwaccel qsv -i /clips/test10_4k.mkv \
  -hwaccel qsv -i /clips/test11_4k.mkv \
  -hwaccel qsv -i /clips/test12_4k.mkv \
  -hwaccel qsv -i /clips/test13_4k.mkv \
  -hwaccel qsv -i /clips/test14_4k.mkv \
  -hwaccel qsv -i /clips/test15_4k.mkv \
  -filter_complex " \
    [0:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4[tile0]; \
    [1:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4[tile1]; \
    [2:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4[tile2]; \
    [3:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4[tile3]; \
    [4:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4[tile4]; \
    [5:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4[tile5]; \
    [6:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4[tile6]; \
    [7:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4[tile7]; \
    [8:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4[tile8]; \
    [9:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4[tile9]; \
    [10:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4[tile10]; \
    [11:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4[tile11]; \
    [12:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4[tile12]; \
    [13:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4[tile13]; \
    [14:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4[tile14]; \
    [15:v]hwupload=extra_hw_frames=10,scale_qsv=w=iw/4:h=ih/4[tile15]; \
    [tile0][tile1][tile2][tile3] \
    [tile4][tile5][tile6][tile7] \
    [tile8][tile9][tile10][tile11] \
    [tile12][tile13][tile14][tile15] xstack_qsv=inputs=16:\
    layout=0_0|0_436|0_872|0_1308|\
             1024_0|1024_436|1024_872|1024_1308|\
             2048_0|2048_436|2048_872|2048_1308|\
             3072_0|3072_436|3072_872|3072_1308[out] \
  " \
  -map "[out]" \
  -c:v hevc_qsv /clips/out_4k.mkv
