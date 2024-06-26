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
   -v $(pwd):/videos \
   -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
   -v /tmp/kahawai_lcore.lock:/tmp/kahawai_lcore.lock \
   -v /dev/null:/dev/null \
   -v /tmp/hugepages:/tmp/hugepages \
   -v /hugepages:/hugepages \
   -v /run/mcm:/run/mcm \
   --expose=20000-20170 \
   --ipc=host -v /dev/shm:/dev/shm \
   --cpuset-cpus=0-20 \
   --net=host \
   -e MTL_PARAM_LCORES=10-20 \
   -e MTL_PARAM_DATA_QUOTA=10356 \
   -e MCM_MEDIA_PROXY_PORT=8002 \
      video_production_image \
         -video_size 1920x1080 -f rawvideo -pix_fmt yuv422p10le -i /videos/src/1080p_yuv422_10b_1.yuv \
         -video_size 1920x1080 -f rawvideo -pix_fmt yuv422p10le -i /videos/src/1080p_yuv422_10b_2.yuv \
         -map 0:v -f mcm \
            -frame_rate 25 \
            -video_size 1920x1080 \
            -pixel_format yuv422p10le \
            -protocol_type auto \
            -payload_type st22 \
            -ip_addr 192.168.96.2 \
            -port 9001 \
            - \
         -map 1:v -f mcm \
            -frame_rate 25 \
            -video_size 1920x1080 \
            -pixel_format yuv422p10le \
            -protocol_type auto \
            -payload_type st22 \
            -ip_addr 192.168.96.2 \
            -port 9002 \
            -
