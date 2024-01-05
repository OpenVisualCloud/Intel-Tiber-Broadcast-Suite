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
  -framerate 50 -pixel_format y210le -width 3840 -height 2160 -port 0000:b1:01.2 -local_addr 192.168.2.2 -src_addr 192.168.2.1 -udp_port 20000 -total_sessions 4 -nullsrc 2 -ext_frames_mode 1 -f kahawai -i "0" \
  -framerate 50 -pixel_format y210le -width 3840 -height 2160 -port 0000:b1:01.2 -local_addr 192.168.2.2 -src_addr 192.168.2.1 -udp_port 20001 -total_sessions 4 -nullsrc 2 -ext_frames_mode 1 -f kahawai -i "1" \
  -framerate 50 -pixel_format y210le -width 3840 -height 2160 -port 0000:b1:01.2 -local_addr 192.168.2.2 -src_addr 192.168.2.1 -udp_port 20002 -total_sessions 4 -nullsrc 2 -ext_frames_mode 1 -f kahawai -i "2" \
  -framerate 50 -pixel_format y210le -width 3840 -height 2160 -port 0000:b1:01.2 -local_addr 192.168.2.2 -src_addr 192.168.2.1 -udp_port 20003 -total_sessions 4 -nullsrc 2 -ext_frames_mode 1 -f kahawai -i "3" \
  -noauto_conversion_filters \
  -filter_complex_frames 2 \
  -filter_complex_policy 1 \
  -filter_complex "\
    [0:v]hwupload=extra_hw_frames=1,scale_qsv=w=iw/2:h=ih/2:mode=compute:async_depth=1[tile0];\
    [1:v]hwupload=extra_hw_frames=1,scale_qsv=w=iw/2:h=ih/2:mode=compute:async_depth=1[tile1];\
    [2:v]hwupload=extra_hw_frames=1,scale_qsv=w=iw/2:h=ih/2:mode=compute:async_depth=1[tile2];\
    [3:v]hwupload=extra_hw_frames=1,scale_qsv=w=iw/2:h=ih/2:mode=compute:async_depth=1[tile3];\
    [tile0][tile1][tile2][tile3]xstack_qsv=inputs=4:layout=0_0|0_h0|w0_0|w0_h0[out];\
    [out]hwdownload,format=y210[multiview]" \
  -map [multiview] -f rawvideo -pix_fmt y210le /dev/null  
