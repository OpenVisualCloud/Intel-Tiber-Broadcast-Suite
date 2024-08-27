#!/bin/bash

#
# SPDX-License-Identifier: BSD-3-Clause
# Copyright(©) 2024 Intel Corporation
# Intel® Tiber™ Broadcast Suite
#

. VARIABLES.rc

# Check if VFIO_PORT_PROC is set
if [ -z "$VFIO_PORT_PROC" ]; then
    echo -e "\e[31mError: VFIO_PORT_PROC is not set.\e[0m"
    echo "Use dpdk-devbind.py -s to check pci address of vfio device"
    exit 1
fi

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
   -v /var/run/imtl:/var/run/imtl \
   --network=my_net_801f0 \
   --ip=192.168.2.2 \
   --expose=20000-20170 \
   --ipc=host -v /dev/shm:/dev/shm \
      video_production_image -y \
      -qsv_device /dev/dri/renderD128 -hwaccel qsv \
      -p_port "${VFIO_PORT_PROC}" -p_sip 192.168.2.2 -p_rx_ip 192.168.2.1 -udp_port 20000 -payload_type 112 -fps 25 -pix_fmt yuv422p10le -video_size 1920x1080 -f mtl_st20p -i "0" \
      -p_port "${VFIO_PORT_PROC}" -p_sip 192.168.2.2 -p_rx_ip 192.168.2.1 -udp_port 20001 -payload_type 112 -fps 25 -pix_fmt yuv422p10le -video_size 1920x1080 -f mtl_st20p -i "1" \
      -p_port "${VFIO_PORT_PROC}" -p_sip 192.168.2.2 -p_rx_ip 192.168.2.1 -udp_port 20002 -payload_type 112 -fps 25 -pix_fmt yuv422p10le -video_size 1920x1080 -f mtl_st20p -i "2" \
      -p_port "${VFIO_PORT_PROC}" -p_sip 192.168.2.2 -p_rx_ip 192.168.2.1 -udp_port 20003 -payload_type 112 -fps 25 -pix_fmt yuv422p10le -video_size 1920x1080 -f mtl_st20p -i "3" \
      -p_port "${VFIO_PORT_PROC}" -p_sip 192.168.2.2 -p_rx_ip 192.168.2.1 -udp_port 20004 -payload_type 112 -fps 25 -pix_fmt yuv422p10le -video_size 1920x1080 -f mtl_st20p -i "4" \
      -p_port "${VFIO_PORT_PROC}" -p_sip 192.168.2.2 -p_rx_ip 192.168.2.1 -udp_port 20005 -payload_type 112 -fps 25 -pix_fmt yuv422p10le -video_size 1920x1080 -f mtl_st20p -i "5" \
      -p_port "${VFIO_PORT_PROC}" -p_sip 192.168.2.2 -p_rx_ip 192.168.2.1 -udp_port 20006 -payload_type 112 -fps 25 -pix_fmt yuv422p10le -video_size 1920x1080 -f mtl_st20p -i "6" \
      -p_port "${VFIO_PORT_PROC}" -p_sip 192.168.2.2 -p_rx_ip 192.168.2.1 -udp_port 20007 -payload_type 112 -fps 25 -pix_fmt yuv422p10le -video_size 1920x1080 -f mtl_st20p -i "7" \
      -filter_complex "[0:v]hwupload,scale_qsv=iw/4:ih/2[out0]; \
                       [1:v]hwupload,scale_qsv=iw/4:ih/2[out1]; \
                       [2:v]hwupload,scale_qsv=iw/4:ih/2[out2]; \
                       [3:v]hwupload,scale_qsv=iw/4:ih/2[out3]; \
                       [4:v]hwupload,scale_qsv=iw/4:ih/2[out4]; \
                       [5:v]hwupload,scale_qsv=iw/4:ih/2[out5]; \
                       [6:v]hwupload,scale_qsv=iw/4:ih/2[out6]; \
                       [7:v]hwupload,scale_qsv=iw/4:ih/2[out7]; \
                       [out0][out1][out2][out3] \
                       [out4][out5][out6][out7] \
                       xstack_qsv=inputs=8:\
                       layout=0_0|w0_0|0_h0|w0_h0|w0+w1_0|w0+w1+w2_0|w0+w1_h0|w0+w1+w2_h0, \
                       format=y210le,format=yuv422p10le" \
      -p_port "${VFIO_PORT_PROC}" -p_sip 192.168.2.2 -p_tx_ip 192.168.2.3 -udp_port 20000 -payload_type 112 -fps 25 -pix_fmt yuv422p10le -f mtl_st20p -
