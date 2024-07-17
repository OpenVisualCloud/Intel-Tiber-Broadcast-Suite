#!/bin/bash

# USAGE: ./blend_compressor.sh [reference_video].yuv

NIC_PORT="0000:b1:01.1"
LOCAL_IP_ADDRESS="192.168.2.1"
DST_IP_ADDRESS="192.168.2.2"
if [ $1 == "" ]; then
  SOURCE_VIDEO_NAME="reference_video.yuv" # output changes .yuv to .mp4
else
  SOURCE_VIDEO_NAME=$1 # output changes .yuv to .mp4
fi

docker run \
  --user root\
  --privileged \
  --device=/dev/vfio:/dev/vfio \
  --device=/dev/dri:/dev/dri \
  --cap-add ALL \
  -v $(pwd)/videos:/videos \
  -v /tmp/kahawai_lcore.lock:/tmp/kahawai_lcore.lock \
  -v /var/run/imtl:/var/run/imtl \
  -v /dev/null:/dev/null \
  --device=/dev/urandom \
  -v /tmp/hugepages:/tmp/hugepages \
  -v /hugepages:/hugepages \
  -e LD_PRELOAD="/usr/lib/gcc/x86_64-linux-gnu/11/libasan.so" \
  -e LD_LIBRARY_PATH=/usr/local/lib/x86_64-linux-gnu/ \
  --network=my_net_801f0 \
  --ip=192.168.2.1 \
  --expose=20000-20170 \
  video_production_image \
   -video_size 3840x2160 -framerate 50 -pix_fmt y210le -i /videos/${SOURCE_VIDEO_NAME} /videos/${SOURCE_VIDEO_NAME/.yuv/.mp4}
