#!/bin/bash

# SPDX-License-Identifier: BSD-3-Clause
# Copyright 2024 Intel Corporation

# One of: yes|no|accept-new
SSH_STRICT_HOST_KEY_CHECKING="accept-new"
SSH_CMD="ssh -oStrictHostKeyChecking=${SSH_STRICT_HOST_KEY_CHECKING} -t -o RemoteCommand="

function exec_command()
{
    local values_returned=""
    local user_at_address=""
    [[ "$#" -eq "2" ]] && user_at_address="${2}"
    [[ "$#" -eq "3" ]] && user_at_address="${3}@${2}"
    if [ "$#" -eq "1" ]; then
        if values_returned="$(eval "${1}")"; then
            echo "$values_returned"
        fi
    elif [[ "$#" -eq "2" ]] || [[ "$#" -eq "3" ]]; then
        if values_returned="$($SSH_CMD"eval \"${1}\"" "${user_at_address}")"; then
            echo "$values_returned"
        fi
    else
        echo "Error: Wrong arguments for exec_command(). Valid number is 1 or 3, got $#" 1>&2
    fi
    if [ -z "$values_returned" ]; then
        echo "Error: Unable to collect results or results are empty." 1>&2
        return 1;
    fi
}

function get_hostname() {
    exec_command 'hostname'
}

function get_intel_nic_device() {
    exec_command "lspci | grep 'Intel Corporation.*\(810\|X722\)'"
}

function get_default_route_nic() {
    exec_command "ip -json r show default | jq '.[0].dev' -r"
}

function get_cpu_arch() {
    local arch=""
    if arch="$(exec_command 'cat /sys/devices/cpu/caps/pmu_name')"; then
        case $arch in
            icelake)
                echo "Xeon IceLake CPU (icx)" 1>&2
                echo "icx"
                ;;
            sapphire_rapids)
                echo "Xeon Sapphire Rapids CPU (spr)" 1>&2
                echo "spr"
                ;;
            skylake)
                echo "Xeon SkyLake" 1>&2
                ;;
            *)
                echo "Error: Unsupported architecture: $arch. Please edit the script or setup the architecture manually." 1>&2
                return 1
                ;;
        esac
    else
        echo "Error: Unable to connect to $1"
        return 1
    fi
    return 0
}
