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
  -v /var/run/imtl:/var/run/imtl \
  -v /dev/null:/dev/null \
  -v /tmp/hugepages:/tmp/hugepages \
  -v /hugepages:/hugepages \
  --network=my_net_801f0 \
  --ip=192.168.2.1 \
  --expose=20000-20170 \
  --ipc=host -v /dev/shm:/dev/shm \
  --cpuset-cpus="28-55" \
  video_production_image \
  -y \
  -an \
  -video_size 3840x2160 -framerate 50 -pix_fmt y210le \
  -i /videos/digit1.yuv \
  -map 0 \
  -vframes 100 \
  -pix_fmt yuv422p10le \
  -port $CAM_NIC_PORT \
  -local_addr $CAM_LOCAL_IP_ADDRESS \
  -dst_addr $CAM_DEST_IP_ADDRESS \
  -udp_port 20000 \
  -total_sessions 1 \
  -f kahawai_mux -
