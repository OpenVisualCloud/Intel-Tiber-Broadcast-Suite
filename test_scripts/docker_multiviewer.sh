#!/bin/bash
NIC_PORT="0000:b1:01.2"
LOCAL_IP_ADDRESS="192.168.2.2"
SOURCE_IP_ADDRESS="192.168.2.1"
DEST_IP_ADDRESS="192.168.2.3"

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
  --cpuset-cpus="56-84" \
  my_ffmpeg \
  -y \
  -an \
  -qsv_device /dev/dri/renderD128 \
  -hwaccel qsv -hwaccel_output_format qsv \
  -thread_queue_size 512 -framerate 60 -pixel_format y210le -width 3840 -height 2160 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -src_addr $SOURCE_IP_ADDRESS -udp_port 20000 -total_sessions 4 -ext_frames_mode 1 -f kahawai -i "0" \
  -thread_queue_size 512 -framerate 60 -pixel_format y210le -width 3840 -height 2160 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -src_addr $SOURCE_IP_ADDRESS -udp_port 20001 -total_sessions 4 -ext_frames_mode 1 -f kahawai -i "1" \
  -thread_queue_size 512 -framerate 60 -pixel_format y210le -width 3840 -height 2160 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -src_addr $SOURCE_IP_ADDRESS -udp_port 20002 -total_sessions 4 -ext_frames_mode 1 -f kahawai -i "2" \
  -thread_queue_size 512 -framerate 60 -pixel_format y210le -width 3840 -height 2160 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -src_addr $SOURCE_IP_ADDRESS -udp_port 20003 -total_sessions 4 -ext_frames_mode 1 -f kahawai -i "3" \
  -noauto_conversion_filters \
  -filter_complex "\
    [0:v]hwupload=extra_hw_frames=4,scale_qsv=w=iw/2:h=ih/2:mode=compute[tile0];\
    [1:v]hwupload,scale_qsv=w=iw/2:h=ih/2:mode=compute[tile1];\
    [2:v]hwupload,scale_qsv=w=iw/2:h=ih/2:mode=compute[tile2];\
    [3:v]hwupload,scale_qsv=w=iw/2:h=ih/2:mode=compute[tile3];\
    [tile0][tile1][tile2][tile3]xstack_qsv=inputs=4:layout=0_0|0_h0|w0_0|w0_h0[out];\
    [out]hwdownload,format=y210le[multiview]" \
-map [multiview] -c:v hevc_qsv /config/out_hevc.mkv
#-map [multiview] -f rawvideo -pix_fmt y210le /dev/null
#TODO: enable transmit path IMTL supports RX and TX in same plugin
#-map "[multiview]" -filter:v fps=60 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -dst_addr $DEST_IP_ADDRESS -udp_port 20000 -total_sessions 1 -f kahawai_mux -
