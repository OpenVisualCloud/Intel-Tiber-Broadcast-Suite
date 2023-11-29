docker run --rm -it \
  --privileged \
  --device=/dev/dri:/dev/dri \
  -v $(pwd):/config \
  -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
  my_ffmpeg \
  -y \
  -qsv_device /dev/dri/renderD128 \
  -an \
  -hwaccel qsv -i /config/test0_4k.mkv \
  -hwaccel qsv -i /config/test1_4k.mkv \
  -hwaccel qsv -i /config/test2_4k.mkv \
  -hwaccel qsv -i /config/test3_4k.mkv \
  -hwaccel qsv -i /config/test4_4k.mkv \
  -hwaccel qsv -i /config/test5_4k.mkv \
  -hwaccel qsv -i /config/test6_4k.mkv \
  -hwaccel qsv -i /config/test7_4k.mkv \
  -hwaccel qsv -i /config/test8_4k.mkv \
  -hwaccel qsv -i /config/test9_4k.mkv \
  -hwaccel qsv -i /config/test10_4k.mkv \
  -hwaccel qsv -i /config/test11_4k.mkv \
  -hwaccel qsv -i /config/test12_4k.mkv \
  -hwaccel qsv -i /config/test13_4k.mkv \
  -hwaccel qsv -i /config/test14_4k.mkv \
  -hwaccel qsv -i /config/test15_4k.mkv \
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
      layout=0_0|0_540|0_1080|0_1620|\
             960_0|960_540|960_1080|960_1620|\ 
             1920_0|1920_540|1920_1080|1920_1620|\
             2880_0|2880_540|2880_1080|2880_1620[out] \
  " \
  -map "[out]" \
  -c:v hevc_qsv /config/out_4k.mkv
