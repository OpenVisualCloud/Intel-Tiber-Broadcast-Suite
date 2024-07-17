#!/bin/bash
NIC_PORT="0000:b1:01.0"
LOCAL_IP_ADDRESS="192.168.2.1"

docker run \
  --user root \
  --privileged \
  -it \
  --device /dev/vfio \
  --cap-add ALL \
  -v /tmp/kahawai_lcore.lock:/tmp/kahawai_lcore.lock \
  -v /var/run/imtl:/var/run/imtl \
  -v /dev/null:/dev/null \
  -v /tmp/hugepages:/tmp/hugepages \
  -v /hugepages:/hugepages \
  --network=my_net_801f0 \
  --ip=192.168.2.1 \
  --expose=20010-20170 \
  video_production_image \
  -loglevel debug \
  -framerate 50 \
  -pixel_format yuv422p10le \
  -width 3840 \
  -height 2160 \
  -udp_port 20010 \
  -port $NIC_PORT \
  -local_addr $LOCAL_IP_ADDRESS \
  -src_addr $LOCAL_IP_ADDRESS \
  -total_sessions 2 \
  -ext_frames_mode 1 \
  -f kahawai \
  -i "1" \
  -framerate 50 \
  -pixel_format yuv422p10le \
  -width 3840 \
  -height 2160 \
  -udp_port 20020 \
  -port $NIC_PORT \
  -local_addr $LOCAL_IP_ADDRESS \
  -src_addr $LOCAL_IP_ADDRESS \
  -total_sessions 2 \
  -ext_frames_mode 1 \
  -f kahawai \
  -i "2" \
  -map 0:0 \
  -vframes 9000 \
  -f rawvideo /dev/null \
  -y \
  -map 1:0 \
  -vframes 9000 \
  -f rawvideo /dev/null \
  -y
