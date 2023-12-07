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
  --cpuset-cpus="28-55" \
  my_ffmpeg \
  -y \
  -an \
  -threads 16 \
  -i /config/test0_4k.mkv \
  -i /config/test1_4k.mkv \
  -i /config/test2_4k.mkv \
  -i /config/test3_4k.mkv \
  -map 0:v -vframes 2000 -filter:v format=y210le,fps=60 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -dst_addr $DEST_IP_ADDRESS -udp_port 20000 -total_sessions 4 -f kahawai_mux -\
  -map 1:v -vframes 2000 -filter:v format=y210le,fps=60 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -dst_addr $DEST_IP_ADDRESS -udp_port 20001 -total_sessions 4 -f kahawai_mux -\
  -map 2:v -vframes 2000 -filter:v format=y210le,fps=60 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -dst_addr $DEST_IP_ADDRESS -udp_port 20002 -total_sessions 4 -f kahawai_mux -\
  -map 3:v -vframes 2000 -filter:v format=y210le,fps=60 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -dst_addr $DEST_IP_ADDRESS -udp_port 20003 -total_sessions 4 -f kahawai_mux -
