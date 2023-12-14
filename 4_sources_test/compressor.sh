#!/bin/bash
NIC_PORT="0000:b1:01.1"
LOCAL_IP_ADDRESS="192.168.2.1"
DST_IP_ADDRESS="192.168.2.2"

docker run \
  --user root\
  --privileged \
  --device=/dev/vfio:/dev/vfio \
  --device=/dev/dri:/dev/dri \
  --cap-add ALL \
  -v $(pwd):/videos \
  -v /tmp/kahawai_lcore.lock:/tmp/kahawai_lcore.lock \
  -v /dev/null:/dev/null \
  --device=/dev/urandom \
  -v /tmp/hugepages:/tmp/hugepages \
  -v /hugepages:/hugepages \
  -e LD_PRELOAD="/usr/lib/gcc/x86_64-linux-gnu/11/libasan.so" \
  -e LD_LIBRARY_PATH=/usr/local/lib/x86_64-linux-gnu/ \
  --network=my_net_801f0 \
  --ip=192.168.2.1 \
  --expose=20000-20170 \
  my_ffmpeg \
   -video_size 3840x2160 -framerate 25 -pix_fmt y210le -i /videos/outfile.yuv /videos/outfile.mp4