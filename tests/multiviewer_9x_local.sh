#!/bin/bash
CAM_NIC_PORT="0000:b1:01.1"
CAM_LOCAL_IP_ADDRESS="192.168.2.1"
CAM_DEST_IP_ADDRESS="192.168.2.2"

docker run \
  --user root\
  --privileged \
  --device=/dev/vfio:/dev/vfio \
  --device=/dev/dri:/dev/dri \
  --cap-add ALL \
  -v $(pwd):/videos \
  -v /usr/lib/x86_64-linux-gnu/dri \
  -v /tmp/kahawai_lcore.lock:/tmp/kahawai_lcore.lock \
  -v /dev/null:/dev/null \
  -v /tmp/hugepages:/tmp/hugepages \
  -v /hugepages:/hugepages \
  --network=my_net_801f0 \
  --ip=192.168.2.1 \
  --expose=20000-20170 \
  --ipc=host -v /dev/shm:/dev/shm \
  --cpuset-cpus="28-55" \
  -e MTL_PARAM_LCORES="28-55" \
  -e MTL_PARAM_DATA_QUOTA=10356 \
  my_ffmpeg \
  -y \
  -hwaccel qsv -hwaccel_output_format qsv \
  -qsv_device /dev/dri/renderD128 \
  -an \
  -framerate 50 -pixel_format y210le -width 3840 -height 2160 -port 0000:b1:01.2 -local_addr 192.168.2.2 -src_addr 192.168.2.1 -udp_port 20000 -total_sessions 9 -nullsrc 2 -ext_frames_mode 1 -f kahawai -i "0" \
  -framerate 50 -pixel_format y210le -width 3840 -height 2160 -port 0000:b1:01.2 -local_addr 192.168.2.2 -src_addr 192.168.2.1 -udp_port 20001 -total_sessions 9 -nullsrc 2 -ext_frames_mode 1 -f kahawai -i "1" \
  -framerate 50 -pixel_format y210le -width 3840 -height 2160 -port 0000:b1:01.2 -local_addr 192.168.2.2 -src_addr 192.168.2.1 -udp_port 20002 -total_sessions 9 -nullsrc 2 -ext_frames_mode 1 -f kahawai -i "2" \
  -framerate 50 -pixel_format y210le -width 3840 -height 2160 -port 0000:b1:01.2 -local_addr 192.168.2.2 -src_addr 192.168.2.1 -udp_port 20003 -total_sessions 9 -nullsrc 2 -ext_frames_mode 1 -f kahawai -i "3" \
  -framerate 50 -pixel_format y210le -width 3840 -height 2160 -port 0000:b1:01.2 -local_addr 192.168.2.2 -src_addr 192.168.2.1 -udp_port 20004 -total_sessions 9 -nullsrc 2 -ext_frames_mode 1 -f kahawai -i "4" \
  -framerate 50 -pixel_format y210le -width 3840 -height 2160 -port 0000:b1:01.2 -local_addr 192.168.2.2 -src_addr 192.168.2.1 -udp_port 20005 -total_sessions 9 -nullsrc 2 -ext_frames_mode 1 -f kahawai -i "5" \
  -framerate 50 -pixel_format y210le -width 3840 -height 2160 -port 0000:b1:01.2 -local_addr 192.168.2.2 -src_addr 192.168.2.1 -udp_port 20006 -total_sessions 9 -nullsrc 2 -ext_frames_mode 1 -f kahawai -i "6" \
  -framerate 50 -pixel_format y210le -width 3840 -height 2160 -port 0000:b1:01.2 -local_addr 192.168.2.2 -src_addr 192.168.2.1 -udp_port 20007 -total_sessions 9 -nullsrc 2 -ext_frames_mode 1 -f kahawai -i "7" \
  -framerate 50 -pixel_format y210le -width 3840 -height 2160 -port 0000:b1:01.2 -local_addr 192.168.2.2 -src_addr 192.168.2.1 -udp_port 20008 -total_sessions 9 -nullsrc 2 -ext_frames_mode 1 -f kahawai -i "8" \
  -noauto_conversion_filters \
  -filter_complex_frames 2 \
  -filter_complex_policy 1 \
  -filter_complex "\
    [0:v]hwupload=extra_hw_frames=1,scale_qsv=w=iw/3:h=ih/3:mode=compute:async_depth=1[t0];\
    [1:v]hwupload=extra_hw_frames=1,scale_qsv=w=iw/3:h=ih/3:mode=compute:async_depth=1[t1];\
    [2:v]hwupload=extra_hw_frames=1,scale_qsv=w=iw/3:h=ih/3:mode=compute:async_depth=1[t2];\
    [3:v]hwupload=extra_hw_frames=1,scale_qsv=w=iw/3:h=ih/3:mode=compute:async_depth=1[t3];\
    [4:v]hwupload=extra_hw_frames=1,scale_qsv=w=iw/3:h=ih/3:mode=compute:async_depth=1[t4];\
    [5:v]hwupload=extra_hw_frames=1,scale_qsv=w=iw/3:h=ih/3:mode=compute:async_depth=1[t5];\
    [6:v]hwupload=extra_hw_frames=1,scale_qsv=w=iw/3:h=ih/3:mode=compute:async_depth=1[t6];\
    [7:v]hwupload=extra_hw_frames=1,scale_qsv=w=iw/3:h=ih/3:mode=compute:async_depth=1[t7];\
    [8:v]hwupload=extra_hw_frames=1,scale_qsv=w=iw/3:h=ih/3:mode=compute:async_depth=1[t8];\
    [t0][t1][t2][t3][t4][t5][t6][t7][t8]xstack_qsv=inputs=9:layout=0_0|w0_0|w0+w1_0|0_h0|w0_h0|w0+w1_h0|0_h0+h1|w0_h0+h1|w0+w1_h0+h1[out];\
    [out]hwdownload,format=y210[multiview]" \
  -map [multiview] -f rawvideo -pix_fmt y210le /dev/null  
