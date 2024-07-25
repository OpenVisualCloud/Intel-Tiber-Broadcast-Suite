#!/bin/bash

#
# SPDX-License-Identifier: BSD-3-Clause
# Copyright(©) 2024 Intel Corporation
# Intel® Tiber™ Broadcast Suite
#

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
    ger-is-registry.caas.intel.com/nex-vs-cicd-automation/mcm/media-proxy:latest \
       /usr/local/bin/media_proxy \
          -d 0000:4b:01.2 \
          -i 192.168.96.2 \
          -t 8003
