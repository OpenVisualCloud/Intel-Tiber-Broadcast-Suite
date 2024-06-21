#!/bin/bash

docker run -it \
   --user root\
   --privileged \
   --device=/dev/vfio:/dev/vfio \
   --device=/dev/dri:/dev/dri \
   --cap-add ALL \
   -v "$(pwd)":/videos \
   -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
   -v /tmp/kahawai_lcore.lock:/tmp/kahawai_lcore.lock \
   -v /dev/null:/dev/null \
   -v /tmp/hugepages:/tmp/hugepages \
   -v /hugepages:/hugepages \
   --network=my_net_801f0 \
   --ip=192.168.2.2 \
   --expose=20000-20170 \
   --ipc=host -v /dev/shm:/dev/shm \
   --cpuset-cpus=20-40 \
   -e MTL_PARAM_LCORES=30-40 \
   -e MTL_PARAM_DATA_QUOTA=10356 \
      video_production_image -y \
      -p_port 0000:4b:01.2 -p_sip 192.168.2.2 -p_rx_ip 192.168.2.1 -udp_port 20000 -payload_type 112 -fps 25 -f mtl_st22p -i "0" \
      -p_port 0000:4b:01.2 -p_sip 192.168.2.2 -p_rx_ip 192.168.2.1 -udp_port 20001 -payload_type 112 -fps 25 -f mtl_st22p -i "1" \
      -map 0:v /videos/recv-jpeg_1.yuv \
      -map 1:v /videos/recv-jpeg_2.yuv
