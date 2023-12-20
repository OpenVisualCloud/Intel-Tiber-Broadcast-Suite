#!/bin/bash
NIC_PORT="0000:b1:01.2"
LOCAL_IP_ADDRESS="192.168.2.2"
SOURCE_IP_ADDRESS="192.168.2.1"
DEST_IP_ADDRESS="192.168.2.3"
OUTPUT_FILE_NAME="output_test_4k.mkv"

docker run -it \
  --user root\
  --privileged \
  --device=/dev/vfio:/dev/vfio \
  --device=/dev/dri:/dev/dri \
  --cap-add ALL \
  -v $(pwd):/config \
  -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
  -v /tmp/kahawai_lcore.lock:/tmp/kahawai_lcore.lock \
  -v /dev/null:/dev/null \
  -v /tmp/hugepages:/tmp/hugepages \
  -v /hugepages:/hugepages \
  --network=my_net_801f0 \
  --ip=192.168.2.2 \
  --expose=20000-20170 \
  --ipc=host -v /dev/shm:/dev/shm \
  --cpuset-cpus="96-127" \
  -e MTL_PARAM_LCORES="123-127" \
  -e MTL_PARAM_DATA_QUOTA=10356 \
  my_ffmpeg \
  -y \
  -an \
  -qsv_device /dev/dri/renderD128 \
  -hwaccel qsv -hwaccel_output_format qsv \
  -framerate 50 -pixel_format y210le -width 3840 -height 2160 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -src_addr $SOURCE_IP_ADDRESS -udp_port 20000 -total_sessions 4 -ext_frames_mode 1 -f kahawai -i "0" \
  -framerate 50 -pixel_format y210le -width 3840 -height 2160 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -src_addr $SOURCE_IP_ADDRESS -udp_port 20001 -total_sessions 4 -ext_frames_mode 1 -f kahawai -i "1" \
  -framerate 50 -pixel_format y210le -width 3840 -height 2160 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -src_addr $SOURCE_IP_ADDRESS -udp_port 20002 -total_sessions 4 -ext_frames_mode 1 -f kahawai -i "2" \
  -framerate 50 -pixel_format y210le -width 3840 -height 2160 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -src_addr $SOURCE_IP_ADDRESS -udp_port 20003 -total_sessions 4 -ext_frames_mode 1 -f kahawai -i "3" \
  -noauto_conversion_filters \
  -filter_complex "\
    [0:v]hwupload=extra_hw_frames=4[t0];\
    [1:v]hwupload[t1];\
    [2:v]hwupload[t2];\
    [3:v]hwupload[t3];\
    [t0]scale_qsv=w=iw/2:h=ih/2:mode=compute:async_depth=1[tile0];\
    [t1]scale_qsv=w=iw/2:h=ih/2:mode=compute:async_depth=1[tile1];\
    [t2]scale_qsv=w=iw/2:h=ih/2:mode=compute:async_depth=1[tile2];\
    [t3]scale_qsv=w=iw/2:h=ih/2:mode=compute:async_depth=1[tile3];\
    [tile0][tile1][tile2][tile3]xstack_qsv=inputs=4:layout=0_0|0_h0|w0_0|w0_h0[multiview]" \
    -map "[multiview]" -c:v hevc_qsv /config/$OUTPUT_FILE_NAME