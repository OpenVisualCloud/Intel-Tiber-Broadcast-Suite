#!/bin/bash -e

#
# SPDX-License-Identifier: BSD-3-Clause
# Copyright(©) 2024 Intel Corporation
# Intel® Tiber™ Broadcast Suite
#

set -eo pipefail

SCRIPT_DIR="$(readlink -f "$(dirname -- "${BASH_SOURCE[0]}")")"
. "${SCRIPT_DIR}/scripts/common.sh"

TIBER_STACK_DEBUG="${TIBER_STACK_DEBUG:-1}" # (future) Force where possible instead of try to configure
rm -f /tmp/kahawai_lcore.lock               # Remove MtlManager legacy indicator/switch if exists
print_logo_anim                             # Print intel animated terminal logo

function setup_jq_package()
{
    local PM=""
    prompt 'Starting setup jq package sequence.'
    if [[ -x "$(command -v jq)" ]]; then
        prompt 'Found jq package installed and PATH available.'
        return 0
    fi
    PM="$(setup_package_manager)" || exit 1
    $PM update && $PM install -y jq && return 0
    error "Got non zero return code from '$PM update && $PM install -y jq && return 0'"
    return 1
}

function add_fstab_line()
{
    ADD_LINE_STRING="${1:-nodev /hugepages hugetlbfs pagesize=1GB 0 0}"
    grep "${ADD_LINE_STRING}" || echo -e "\n${ADD_LINE_STRING}" >> /etc/fstab
}

function copy_nicctl_script()
{
    local script_result=""
    script_result="0"
    prompt 'Starting copy_nicctl_script sequence.'
    if [ ! -f "/usr/local/bin/nicctl.sh" ]; then
        docker create --name mtl-tmp mtl-manager:latest 2>&1 && \
        docker cp mtl-tmp:/home/mtl/nicctl.sh /usr/local/bin 2>&1 && \
        docker rm mtl-tmp 2>&1
        script_result="$?"
    fi
    if [ "$script_result" == "0" ]; then
        prompt 'Finished copy_nicctl_script sequence. Success.'
        return 0
    fi
    error 'Finished copy_nicctl_script sequence.'
    return 1
}

function setup_vfio_subsytem()
{
    prompt 'Starting setup_vfio_subsytem sequence.'
    if [[ "${TIBER_STACK_DEBUG}" != "0" ]]; then
        getent group 2110 > /dev/null || groupadd -g 2110 vfio
        usermod -aG vfio "$USER"
        touch /etc/udev/rules.d/10-vfio.rules
        if ! grep -q '^SUBSYSTEM=="vfio", GROUP="vfio"' /etc/udev/rules.d/10-vfio.rules; then
            echo 'SUBSYSTEM=="vfio", GROUP="vfio", MODE="0660"' >> /etc/udev/rules.d/10-vfio.rules
            udevadm control --reload-rules
            udevadm trigger
        fi
    else
        chmod 666 -R /dev/vfio
    fi
    prompt 'Finished setup_vfio_subsytem sequence. Success'
}

function setup_hugepages()
{
  prompt 'Starting setup_hugepages sequence.'
  mkdir -p /tmp/hugepages /hugepages
  # lsmem --json | jq '.memory[].size'
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
  prompt 'Finished setup_hugepages sequence. Success'
}

function setup_docker_network()
{
    prompt 'Starting setup_docker_network sequence.'
    local parent_nic=""
    parent_nic="$(get_default_route_nic)"
    if ! docker network create --subnet 192.168.2.0/24 --gateway 192.168.2.100 -o parent="${parent_nic}" my_net_801f0 2>/dev/null; then
        warning 'Network with name my_net_801f0 already exists'
    fi
    prompt 'Finished setup_docker_network sequence. Success'
}

function setup_nic_virtual_functions()
{
    prompt 'Starting create virtual functions sequence.'
    output=$(get_intel_nic_device | cut -f1 -d' ' | awk '{print "0000:"$1}')
    IFS=$'\n'

    [ -f "/usr/local/bin/nicctl.sh" ] || copy_nicctl_script
    if [ "$?" -ne "0" ]; then
        error 'Container mtl-manager:latest or nicctl.sh script failed. Exiting.'
        exit 1
    fi

    while IFS= read -r line; do
        /usr/local/bin/nicctl.sh disable_vf "$line" 1>/dev/null
        if ! /usr/local/bin/nicctl.sh create_vf "$line" ; then
            error "Error occurred while creating VF for device: '$line'"
            exit 2
        fi
    done <<< "$output"
    prompt 'Finished create virtual functions sequence. Success.'
}

function setup_mtl_manager_container
{
    prompt 'Starting run sequence for mtl-manager:latest image.'
    container_id="$(docker ps -aq -f name=^mtl-manager$)"

    if [ -n "$container_id" ]; then
        if [ "$(docker inspect -f '{{.State.Running}}' "$container_id")" = "true" ]; then
            prompt 'Container mtl-manager is already running.'
        else
            warning 'Container mtl-manager exists but is not running. Removing it.'
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
    prompt 'Finished run sequence for mtl-manager:latest image. Success.'
}

setup_jq_package
setup_vfio_subsytem
setup_hugepages
setup_docker_network
setup_nic_virtual_functions
setup_mtl_manager_container

