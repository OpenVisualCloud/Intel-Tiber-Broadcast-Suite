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
  -v $(pwd):/config \
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
  video_production_image \
  -y \
  -an \
  -video_size 640x480 -framerate 25 -pix_fmt y210le \
  -i /config/gradients.yuv \
  -filter_complex "[0:v]format=yuv422p10le,fps=25,split=4[in1][in2][in3][in4]" \
  -map "[in1]" -vframes 250 -pix_fmt yuv422p10le -port $CAM_NIC_PORT -local_addr $CAM_LOCAL_IP_ADDRESS -dst_addr $CAM_DEST_IP_ADDRESS -udp_port 20000 -total_sessions 4 -f kahawai_mux -\
  -map "[in2]" -vframes 250 -pix_fmt yuv422p10le -port $CAM_NIC_PORT -local_addr $CAM_LOCAL_IP_ADDRESS -dst_addr $CAM_DEST_IP_ADDRESS -udp_port 20001 -total_sessions 4 -f kahawai_mux -\
  -map "[in3]" -vframes 250 -pix_fmt yuv422p10le -port $CAM_NIC_PORT -local_addr $CAM_LOCAL_IP_ADDRESS -dst_addr $CAM_DEST_IP_ADDRESS -udp_port 20002 -total_sessions 4 -f kahawai_mux -\
  -map "[in4]" -vframes 250 -pix_fmt yuv422p10le -port $CAM_NIC_PORT -local_addr $CAM_LOCAL_IP_ADDRESS -dst_addr $CAM_DEST_IP_ADDRESS -udp_port 20003 -total_sessions 4 -f kahawai_mux -
