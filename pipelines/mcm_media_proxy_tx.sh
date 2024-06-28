#!/bin/bash

#
# SPDX-License-Identifier: BSD-3-Clause
# Copyright(©) 2024 Intel Corporation
# Intel® Tiber™ Broadcast Suite
#

# docker run -it \
#    --privileged \
#    -u 0:0 \
#    --net=host \
#    --device=/dev/vfio:/dev/vfio \
#    --device=/dev/dri:/dev/dri \
#    -v /var/run/imtl:/var/run/imtl \
#    -v /run/mcm:/run/mcm \
#    -v /dev/hugepages:/dev/hugepages \
#    -v /dev/hugepages1G:/dev/hugepages1G \
#    -v /sys/fs/bpf:/sys/fs/bpf \
#    -v /dev/shm:/dev/shm \
#    -v /dev/vfio:/dev/vfio \
#    --ipc=host \
#    ger-is-registry.caas.intel.com/nex-vs-cicd-automation/mcm/media-proxy:latest \
#       /usr/local/bin/media_proxy \
#          -d 0000:32:11.0 \
#          -i 192.168.96.1 \
#          -t 8002

# The container environment for MCM media_proxy is Work in Progress.
# Temporary workaround:
#  1. Install Media Communications Mesh
#  2. Run media_proxy in the host OS.
media_proxy -d 0000:32:11.0 -i 192.168.96.1 -t 8002
