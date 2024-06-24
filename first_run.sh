#!/bin/bash -e

#
# SPDX-License-Identifier: BSD-3-Clause
# Copyright(©) 2024 Intel Corporation
# Intel® Tiber™ Broadcast Suite
#

STACK_DEBUG="${STACK_DEBUG:-1}"

function add_fstab_line()
{
    ADD_LINE_STRING="${1:-nodev /hugepages hugetlbfs pagesize=1GB 0 0}"
    grep "${ADD_LINE_STRING}" || echo -e "\n${ADD_LINE_STRING}" >> /etc/fstab
}

if [[ "${STACK_DEBUG}" == "0" ]]
then
    getent group 2110 || sudo groupadd -g 2110 vfio
    sudo usermod -aG vfio "$USER"
    echo 'SUBSYSTEM=="vfio", GROUP="vfio", MODE="0660"' | sudo tee /etc/udev/rules.d/10-vfio.rules
    udevadm control --reload-rules
    udevadm trigger
else
    chmod 777 -R /dev/vfio
fi

mkdir -p /tmp/hugepages /hugepages
for pt in /sys/devices/system/node/node*
do
    # sysctl -w vm.nr_hugepages=4096
    echo 2048 > "$pt/hugepages/hugepages-2048kB/nr_hugepages";
    echo 1 > "$pt/hugepages/hugepages-1048576kB/nr_hugepages";
done
mount -t hugetlbfs hugetlbfs /tmp/hugepages -o pagesize=2M
mount -t hugetlbfs hugetlbfs /hugepages -o pagesize=1G
# add_fstab_line "nodev /tmp/hugepages hugetlbfs pagesize=2M 0 0"
# add_fstab_line "nodev /hugepages hugetlbfs pagesize=1GB 0 0"

touch /tmp/kahawai_lcore.lock
docker network create --subnet 192.168.2.0/24 --gateway 192.168.2.100 -o parent=ens801f0 my_net_801f0
output=$(lspci | grep "Ethernet controller: Intel Corporation Ethernet Controller" | awk '{print "0000:"$1}')
IFS=$'\n'
while IFS= read -r line; do find "${HOME}" -name "nicctl.sh" -exec {} create_vf "$line" \;; done <<< "$output"
