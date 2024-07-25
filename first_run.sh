#!/bin/bash -e

#
# SPDX-License-Identifier: BSD-3-Clause
# Copyright(©) 2024 Intel Corporation
# Intel® Tiber™ Broadcast Suite
#

STACK_DEBUG="${STACK_DEBUG:-1}"
if [ -z "$mtl_source_code" ]; then
    mtl_source_code="$HOME";
    echo -e '\e[33mmtl_source_code variable was empty and has been set to '"$HOME"' directory.\e[0m'
fi

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
if ! docker network create --subnet 192.168.2.0/24 --gateway 192.168.2.100 -o parent=ens801f0 my_net_801f0 2>/dev/null; then
    echo -e '\e[33mNetwork with the name my_net_801f0 already exists\e[0m'
fi

output=$(lspci | grep "Ethernet controller: Intel Corporation Ethernet Controller" | awk '{print "0000:"$1}')
IFS=$'\n'

NICCTL=$(find "${mtl_source_code}" -name "nicctl.sh" -print -quit 2>/dev/null);
if [ -z "$NICCTL" ]; then
    echo -e '\e[31mnicctl.sh script not found inside '"${mtl_source_code}"'\e[0m'
    exit 1
fi

while IFS= read -r line; do
    if ! "$NICCTL" create_vf "$line" ; then
        echo -e '\e[31mError occurred while creating VF for device: '"$line"'\e[0m'
        exit 2
    fi
done <<< "$output"

container_id=$(docker ps -aq -f name=^mtl-manager$)

if [ -n "$container_id" ]; then
    if [ "$(docker inspect -f '{{.State.Running}}' "$container_id")" = "true" ]; then
        echo -e '\e[32mContainer mtl-manager is already running.\e[0m'
    else
        echo -e '\e[33mContainer mtl-manager exists but is not running. Removing it...\e[0m'
        docker rm "$container_id"

        docker run -d \
          --name mtl-manager \
          --privileged --net=host \
          -v /var/run/imtl:/var/run/imtl \
          -v /sys/fs/bpf:/sys/fs/bpf \
          mtl-manager:latest
    fi
else
    docker run -d \
      --name mtl-manager \
      --privileged --net=host \
      -v /var/run/imtl:/var/run/imtl \
      -v /sys/fs/bpf:/sys/fs/bpf \
      mtl-manager:latest
fi