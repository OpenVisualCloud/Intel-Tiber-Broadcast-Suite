#! /bin/bash

# SPDX-License-Identifier: BSD-3-Clause
# Copyright 2024 Intel Corporation
# Intel® Tiber™ Broadcast Suite

mkdir -p /tmp/hugepages
mount -t hugetlbfs hugetlbfs /tmp/hugepages -o pagesize=2M
mount -t hugetlbfs hugetlbfs /hugepages -o pagesize=1G
#./tmp/MTL/script/nicctl.sh create_vf ${NIC_PORT}
modprobe ice irdma vfio-pci
#gdb --args /usr/local/bin/ffmpeg_g $@
/usr/local/bin/ffmpeg_g $@
