#!/bin/bash

#
# SPDX-License-Identifier: BSD-3-Clause
# Copyright(©) 2024 Intel Corporation
# Intel® Tiber™ Broadcast Suite
#

. VARIABLES.rc 2>/dev/null

# Check if VFIO_PORT_R is set
if [ -z "$VFIO_PORT_R" ]; then
    echo -e "\e[31mError: VFIO_PORT_R is not set.\e[0m"
    echo "Use dpdk-devbind.py -s to check pci address of vfio device"
    exit 1
fi

while getopts "l" opt; do
    case ${opt} in
        l )
            media_proxy -d "${VFIO_PORT_R}" -i 192.168.96.2 -t 8003
            exit 0
            ;;
        \? )
            echo "Invalid option: -$OPTARG" >&2
            ;;
    esac
done


 docker run -it \
    --privileged \
    -u 0:0 \
    --net=host \
    --device=/dev/vfio:/dev/vfio \
    --device=/dev/dri:/dev/dri \
    -v /var/run/imtl:/var/run/imtl \
    -v /run/mcm:/run/mcm \
    -v /tmp/hugepages:/dev/hugepages \
    -v /hugepages:/dev/hugepages1G \
    -v /var/run/imtl:/var/run/imtl \
    -v /sys/fs/bpf:/sys/fs/bpf \
    -v /dev/shm:/dev/shm \
    -v /dev/vfio:/dev/vfio \
    --ipc=host \
    --expose 8000-9100 \
    media-proxy:latest \
       /usr/local/bin/media_proxy \
          -d "${VFIO_PORT_R}" \
          -i 192.168.96.2 \
          -t 8003
