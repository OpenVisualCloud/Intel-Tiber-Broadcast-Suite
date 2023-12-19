#!/bin/bash
NIC_PORT="0000:b1:01.1"
LOCAL_IP_ADDRESS="192.168.2.1"
DEST_IP_ADDRESS="192.168.2.2"

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
  --ip=192.168.2.1 \
  --expose=20000-20170 \
  --ipc=host -v /dev/shm:/dev/shm \
  --cpuset-cpus="32-63" \
  -e MTL_PARAM_LCORES="59-63" \
  -e MTL_PARAM_DATA_QUOTA=10356 \
  my_ffmpeg \
  -y \
  -an \
  -i /config/test0_4k.mkv \
  -filter_complex "[0:v]format=yuv422p10le,fps=50,split=4[in1][in2][in3][in4]" \
  -map "[in1]" -vframes 4000 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -dst_addr $DEST_IP_ADDRESS -udp_port 20000 -total_sessions 4 -f kahawai_mux -\
  -map "[in2]" -vframes 4000 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -dst_addr $DEST_IP_ADDRESS -udp_port 20001 -total_sessions 4 -f kahawai_mux -\
  -map "[in3]" -vframes 4000 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -dst_addr $DEST_IP_ADDRESS -udp_port 20002 -total_sessions 4 -f kahawai_mux -\
  -map "[in4]" -vframes 4000 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -dst_addr $DEST_IP_ADDRESS -udp_port 20003 -total_sessions 4 -f kahawai_mux -
