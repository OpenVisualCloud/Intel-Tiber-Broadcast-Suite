#!/bin/bash

#
# SPDX-License-Identifier: BSD-3-Clause
# Copyright(©) 2024 Intel Corporation
# Intel® Tiber™ Broadcast Suite
#

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
   -v /run/mcm:/run/mcm \
   --expose=20000-20170 \
   --ipc=host -v /dev/shm:/dev/shm \
   --cpuset-cpus=20-40 \
   --net=host \
   -e MTL_PARAM_LCORES=30-40 \
   -e MTL_PARAM_DATA_QUOTA=10356 \
   -e MCM_MEDIA_PROXY_PORT=8003 \
      video_production_image -y \
        -f mcm \
           -frame_rate 25 \
           -video_size 1920x1080 \
           -pixel_format yuv422p10le \
           -protocol_type auto \
           -payload_type st22 \
           -ip_addr 192.168.96.1 \
           -port 9001 \
           -i "0" \
        -f mcm \
           -frame_rate 25 \
           -video_size 1920x1080 \
           -pixel_format yuv422p10le \
           -protocol_type auto \
           -payload_type st22 \
           -ip_addr 192.168.96.1 \
           -port 9002 \
           -i "1" \
         -map 0:v /videos/recv-mcm_1.yuv \
         -map 1:v /videos/recv-mcm_2.yuv
