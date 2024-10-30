#!/bin/bash

#
# SPDX-License-Identifier: BSD-3-Clause
# Copyright(©) 2024 Intel Corporation
# Intel® Tiber™ Broadcast Suite
#

function help() {
    echo "Usage: $0 [-l]"
    echo
    echo "Options:"
    echo "  -l    Run the pipeline on bare metal locally."
    echo
    echo "For more information, please refer to docs/run.md."
    exit 0
}

while getopts "lh" opt; do
    case ${opt} in
        l )
            echo "Running pipeline on bare metal locally..."
            MCM_MEDIA_PROXY_PORT=8002
            ./ffmpeg  -y \
            -video_size 1920x1080 -f rawvideo -pix_fmt yuv422p10le -i src/1080p_yuv422_10b_1.yuv \
            -video_size 1920x1080 -f rawvideo -pix_fmt yuv422p10le -i src/1080p_yuv422_10b_2.yuv \
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
          exit 0
            ;;
        h )
            help
            ;;
        \? )
            echo "Invalid option: -$OPTARG" >&2
            help
            ;;
    esac
done

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
   -v /var/run/imtl:/var/run/imtl \
   --expose=20000-20170 \
   --ipc=host -v /dev/shm:/dev/shm \
   --net=host \
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
