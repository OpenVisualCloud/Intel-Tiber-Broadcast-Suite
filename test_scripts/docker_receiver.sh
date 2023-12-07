#!/bin/bash
NIC_PORT="0000:b1:01.2"
LOCAL_IP_ADDRESS="192.168.2.2"
SOURCE_IP_ADDRESS="192.168.2.1"

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
  -thread_queue_size 512 \
  -framerate 60 -pixel_format y210le -width 3840 -height 2160 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -src_addr $SOURCE_IP_ADDRESS -udp_port 20000 -total_sessions 4 -ext_frames_mode 1 -f kahawai -i "0" \
  -framerate 60 -pixel_format y210le -width 3840 -height 2160 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -src_addr $SOURCE_IP_ADDRESS -udp_port 20001 -total_sessions 4 -ext_frames_mode 1 -f kahawai -i "1" \
  -framerate 60 -pixel_format y210le -width 3840 -height 2160 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -src_addr $SOURCE_IP_ADDRESS -udp_port 20002 -total_sessions 4 -ext_frames_mode 1 -f kahawai -i "2" \
  -framerate 60 -pixel_format y210le -width 3840 -height 2160 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -src_addr $SOURCE_IP_ADDRESS -udp_port 20003 -total_sessions 4 -ext_frames_mode 1 -f kahawai -i "3" \
  -map 0:0 -f rawvideo /dev/null \
  -map 1:0 -f rawvideo /dev/null \
  -map 2:0 -f rawvideo /dev/null \
  -map 3:0 -f rawvideo /dev/null 
