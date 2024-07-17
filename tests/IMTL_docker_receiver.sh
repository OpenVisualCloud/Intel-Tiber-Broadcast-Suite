#!/bin/bash
RECV_NIC_PORT="0000:b1:01.2"
RECV_LOCAL_IP_ADDRESS="192.168.2.2"
RECV_SOURCE_IP_ADDRESS="192.168.2.1"

docker run \
  --user root\
  --privileged \
  --device=/dev/vfio:/dev/vfio \
  --device=/dev/dri:/dev/dri \
  --cap-add ALL \
  -v $(pwd):/config \
  -v /tmp/kahawai_lcore.lock:/tmp/kahawai_lcore.lock \
  -v /var/run/imtl:/var/run/imtl \
  -v /dev/null:/dev/null \
  -v /tmp/hugepages:/tmp/hugepages \
  -v /hugepages:/hugepages \
  --network=my_net_801f0 \
  --ip=192.168.2.2 \
  --expose=20000-20170 \
  --ipc=host -v /dev/shm:/dev/shm \
  --cpuset-cpus="84-111" \
  video_production_image \
  -y \
  -an \
  -thread_queue_size 512 \
  -framerate 25 -pixel_format y210le -pix_fmt y210le -width 640 -height 480 -port $RECV_NIC_PORT -local_addr $RECV_LOCAL_IP_ADDRESS -src_addr $RECV_SOURCE_IP_ADDRESS -udp_port 20000 -total_sessions 4 -f kahawai -i "0" \
  -framerate 25 -pixel_format y210le -pix_fmt y210le -width 640 -height 480 -port $RECV_NIC_PORT -local_addr $RECV_LOCAL_IP_ADDRESS -src_addr $RECV_SOURCE_IP_ADDRESS -udp_port 20001 -total_sessions 4 -f kahawai -i "1" \
  -framerate 25 -pixel_format y210le -pix_fmt y210le -width 640 -height 480 -port $RECV_NIC_PORT -local_addr $RECV_LOCAL_IP_ADDRESS -src_addr $RECV_SOURCE_IP_ADDRESS -udp_port 20002 -total_sessions 4 -f kahawai -i "2" \
  -framerate 25 -pixel_format y210le -pix_fmt y210le -width 640 -height 480 -port $RECV_NIC_PORT -local_addr $RECV_LOCAL_IP_ADDRESS -src_addr $RECV_SOURCE_IP_ADDRESS -udp_port 20003 -total_sessions 4 -f kahawai -i "3" \
  -map 0:0 -f rawvideo -pix_fmt y210le /config/received.yuv \
  -map 1:0 -f rawvideo /dev/null \
  -map 2:0 -f rawvideo /dev/null \
  -map 3:0 -f rawvideo /dev/null 

