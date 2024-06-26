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
   --network=my_net_801f0 \
   --ip=192.168.2.1 \
   --expose=20000-20170 \
   --ipc=host -v /dev/shm:/dev/shm \
   --cpuset-cpus=11-20 \
   -e MTL_PARAM_LCORES=15-20 \
   -e MTL_PARAM_DATA_QUOTA=10356 \
      video_production_image \
      -video_size 3840x2160 -f rawvideo -pix_fmt yuv422p10le -i /videos/src/2160p_yuv422_10b.yuv -filter:v fps=25 \
      -p_port 0000:4b:01.2 -p_sip 192.168.2.1 -p_tx_ip 192.168.2.2 -udp_port 20000 -payload_type 112 -f mtl_st20p -
