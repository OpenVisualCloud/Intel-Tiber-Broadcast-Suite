#!/bin/bash -e

#
# SPDX-License-Identifier: BSD-3-Clause
# Copyright(©) 2024 Intel Corporation
# Intel® Tiber™ Broadcast Suite
#

set -eo pipefail

SCRIPT_DIR="$(readlink -f "$(dirname -- "${BASH_SOURCE[0]}")")"
. "${SCRIPT_DIR}/common.sh"


cat "${VERSIONS_ENVIRONMENT_FILE:-${SCRIPT_DIR}/../versions.env}" > "${SCRIPT_DIR}/../.temp.env"
source "${SCRIPT_DIR}/../.temp.env"

INSTALL_DPDK="${INSTALL_DPDK:-false}"
TIBER_STACK_DEBUG="${TIBER_STACK_DEBUG:-1}" # (future) Force where possible instead of try to configure
rm -f /tmp/kahawai_lcore.lock               # Remove MtlManager legacy indicator/switch if exists
print_logo_anim                             # Print intel animated terminal logo

function cleanup_directory {
    local dir_name=$1

    if [ "$LOCAL_INSTALL_CLEANUP" != true ]; then
        return 0
    fi

    if [ -z "${log_file}" ] || [ ! -f "${log_file}" ]; then
        log_file=/dev/null
    fi

    if ! rm -drf "$dir_name" >>$log_file 2>&1; then
        log_warning "$dir_name cleanup failed"
    fi
}

function setup_jq_package()
{
    local PM=""
    log_info 'Starting setup jq package sequence.'
    if [[ -x "$(command -v jq)" ]]; then
        log_info 'Found jq package installed and PATH available.'
        return 0
    fi
    PM="$(setup_package_manager)" || exit 1
    $PM update && $PM install -y jq && return 0
    log_error "Got non zero return code from '$PM update && $PM install -y jq && return 0'"
    return 1
}

function add_fstab_line()
{
    ADD_LINE_STRING="${1:-nodev /hugepages hugetlbfs pagesize=1GB 0 0}"
    grep "${ADD_LINE_STRING}" || echo -e "\n${ADD_LINE_STRING}" >> /etc/fstab
}

function install_dpdk() {
    if ! [ "$INSTALL_DPDK" = true ]; then
        return 0
    fi
    if ! (mkdir -p "${HOME}/dpdk" &&
          curl -Lf https://github.com/DPDK/dpdk/archive/refs/tags/v${DPDK_VER}.tar.gz | \
          tar -zx --strip-components=1 -C "${HOME}/dpdk" ); then
        log_error "DPDK download and extraction failed"
        return 2
    fi

    if [ ! -d "${HOME}/Media-Transport-Library" ] && ! (mkdir -p ${HOME}/Media-Transport-Library &&
          curl -Lf https://github.com/OpenVisualCloud/Media-Transport-Library/archive/refs/heads/${MTL_VER}.tar.gz | \
          tar -zx --strip-components=1 -C ${HOME}/Media-Transport-Library ); then
        log_error "MTL download failed"
        return 2
    fi

    if ! (cd "${HOME}/dpdk" &&
          git apply "${HOME}/Media-Transport-Library/patches/dpdk/$DPDK_VER"/*.patch &&
          cd .. &&
          rm -rf "${HOME}/Media-Transport-Library"); then
        log_error "Patching DPDK with Media-Transport-Library patches failed"
        return 2
    fi

    if ! (cd "${HOME}/dpdk" &&
          meson build &&
          ninja -C build &&
          sudo ninja -C build install); then
        log_error "DPDK build and installation failed"
        return 2
    fi

    cleanup_directory "${HOME}/dpdk"
    cleanup_directory "${HOME}/Media-Transport-Library"
}


function copy_nicctl_script()
{
    local script_result=""
    script_result="0"

    log_info 'Starting copy_nicctl_script sequence.'
    if [ ! -f "/usr/local/bin/nicctl.sh" ]; then
        docker create --name mtl-tmp mtl-manager:latest 2>&1 && \
        docker cp mtl-tmp:/home/mtl/nicctl.sh /usr/local/bin 2>&1 && \
        docker rm mtl-tmp 2>&1
        script_result="$?"
    fi

    if [ "$script_result" != "0" ]; then
        . versions.env &&
        STRIPPED_VER=${MTL_VER#v} &&
        if ! sudo wget -O /usr/local/bin/nicctl.sh https://raw.githubusercontent.com/OpenVisualCloud/Media-Transport-Library/refs/heads/"${STRIPPED_VER}"/script/nicctl.sh; then
            log_error "Failed to download nicctl.sh script"
            return 1
        fi
        if ! sudo chmod +x /usr/local/bin/nicctl.sh; then
            log_error "Failed to set executable permissions for nicctl.sh script"
            return 1
        fi
        script_result="$?"
    fi

    if [ "$script_result" == "0" ]; then
        log_info 'Finished copy_nicctl_script sequence. Success.'
        return 0
    fi
    log_error 'Sequence copy_nicctl_script failed'
    return 1
}

function setup_vfio_subsytem()
{
    log_info 'Starting setup_vfio_subsytem sequence.'
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
        chmod 777 -R /dev/vfio
    fi
    log_success 'Finished setup_vfio_subsytem sequence. Success'
}

function setup_hugepages()
{
  log_info 'Starting setup_hugepages sequence.'
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
  log_success 'Finished setup_hugepages sequence. Success'
}

function setup_docker_network()
{
    log_info 'Starting setup_docker_network sequence.'
    local parent_nic=""
    parent_nic="$(get_default_route_nic)"
    if ! docker network create --subnet 192.168.2.0/24 --gateway 192.168.2.100 -o parent="${parent_nic}" my_net_801f0 2>/dev/null; then
        log_warning 'Network with name my_net_801f0 already exists'
    fi
    log_success 'Finished setup_docker_network sequence. Success'
}

function setup_nic_virtual_functions()
{
    log_info 'Starting create virtual functions sequence.'

    if [ -n "$E810_PCIE_SPECIFIED" ]; then
        output=$(echo "$E810_PCIE_SPECIFIED" | tr ' ' '\n')
        log_info "Selected NICs $E810_PCIE_SPECIFIED"
    else
        output=$(get_intel_nic_device | cut -f1 -d' ' | awk '{print "0000:"$1}')
    fi
    IFS=$'\n'


    if [ ! -f "/usr/local/bin/nicctl.sh" ]; then
        if ! copy_nicctl_script; then
            log_error "Failed to copy nicctl.sh script. Exiting."
            exit 1
        fi
    fi
    if [ "$?" -ne "0" ]; then
        log_error 'Container mtl-manager:latest or nicctl.sh script failed. Exiting.'
        exit 1
    fi

    while IFS= read -r line; do
        sudo /usr/local/bin/nicctl.sh disable_vf "$line" 1>/dev/null
        if ! sudo /usr/local/bin/nicctl.sh create_vf "$line" ; then
            log_error "Error occurred while creating VF for device: '$line'"
            exit 2
        fi
    done <<< "$output"
    log_success 'Finished create virtual functions sequence. Success.'
}

function setup_mtl_manager_container
{
    log_info 'Starting run sequence for mtl-manager:latest image.'
    container_id="$(docker ps -aq -f name=^mtl-manager$)"

    if [ -n "$container_id" ]; then
        if [ "$(docker inspect -f '{{.State.Running}}' "$container_id")" = "true" ]; then
            log_info 'Container mtl-manager is already running.'
        else
            log_warning 'Container mtl-manager exists but is not running. Removing it.'
            docker rm "$container_id"
            docker run -d \
              --name mtl-manager \
              --privileged --net=host \
              -v /var/run/imtl:/var/run/imtl \
              -v /sys/fs/bpf:/sys/fs/bpf \
              mtl-manager:latest || return 2
        fi
    else
        if ! docker run -d \
          --name mtl-manager \
          --privileged --net=host \
          -v /var/run/imtl:/var/run/imtl \
          -v /sys/fs/bpf:/sys/fs/bpf \
          mtl-manager:latest; then
            log_warning "Failed to start mtl-manager container"
        else
            log_success 'Finished run sequence for mtl-manager:latest image. Success.'
        fi
    fi
}

print_help() {
    echo "Usage: $0 [-l] [-h] [-e PCIe_ADDRESSES]"
    echo "Options:"
    echo "  -l    For users running Intel® Tiber™ Broadcast Suite on bare metal."
    echo "  -h    Display this help message."
    echo "  -d    Install DPDK with MTL patches."
    echo "  -e    Specify the PCIe addresses for the E810 NIC."
    echo "        By default, all PCIe E810 addresses are selected."
    echo "        e.g. ./first_run.sh -e \"0000:4b:00.0 0000:4b:00.1 \""
    exit 0
}

while getopts "lhe:d" opt; do
    case ${opt} in
    l )
        setup_vfio_subsytem
        setup_hugepages
        setup_nic_virtual_functions

        if ! setup_mtl_manager_container 1>/dev/null 2>&1 && ! pgrep -x "MtlManager" > /dev/null; then 
            log_info 'Now starting Mtl Manager'
            nohup sudo MtlManager > /dev/null 2>&1 &
            log_info 'Mtl Manager running in background'
        fi

        exit 0
    ;;
    h )
        print_help
    ;;
    e )
        E810_PCIE_SPECIFIED=${OPTARG}
    ;;
    d )
        INSTALL_DPDK=true
    ;;
    \? )
        exit 1
        ;;
    esac
done

setup_jq_package
install_dpdk
setup_vfio_subsytem
setup_hugepages
setup_docker_network
setup_nic_virtual_functions
setup_mtl_manager_container

