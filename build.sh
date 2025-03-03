#!/bin/bash

#
# SPDX-License-Identifier: BSD-3-Clause
# Copyright(©) 2024 Intel Corporation
# Intel® Tiber™ Broadcast Suite
#

set -eo pipefail
SCRIPT_DIR="$(readlink -f "$(dirname -- "${BASH_SOURCE[0]}")")"
LOCAL_INSTALL=false
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m'
CLEAR_LINE='\r\033[K'
VERSIONS_ENVIRONMENT_FILE=versions.env
LOCAL_INSTALL_LOG_DIRECTORY=/tmp/tiberbuild/log
LOCAL_INSTALL_LOG_FILE="tiber_local_build_$(date '+%Y-%m-%d-%H%M%S')"
LOCAL_INSTALL_DEBIAN_DIRECTORY="${SCRIPT_DIR}"/local_debians
LOCAL_INSTALL_DEPENDENCIES_DIRECTORY=/tmp/tiberbuild/build
LOCAL_INSTALL_CLEANUP=true
log_file="$LOCAL_INSTALL_LOG_DIRECTORY/$LOCAL_INSTALL_LOG_FILE"

# Function to display text in the center of the terminal
# Arguments:
#   $1 -- text (the text you want to display)
# Arguments:
function center {
    local TERM_WIDTH
    TERM_WIDTH=$(tput cols)
    local TEXT_LENGTH=${#1}
    local PADDING_NUMBER=$(( (TERM_WIDTH - TEXT_LENGTH) / 2 ))
    printf "%*s%s\n" $PADDING_NUMBER "" "$1"
}

# Function to display a progress bar leaves the terminal cursor at the end
#   of the line (\r) when the status is different from completed
# Arguments:
#   $1 -- progress (current progress)
#   $2 -- out of (total progress)
function progress_bar {
    local TERM_WIDTH
    TERM_WIDTH=$(tput cols)
    local BAR_WIDTH=$((TERM_WIDTH * 80 / 100 - 4))
    local PROGRESS=$(( ($1 * 100 / $2 * 100) / 100 ))
    local PROGRESS_DONE=$(( (PROGRESS * BAR_WIDTH) / 100 ))
    local PROGRESS_LEFT=$(( BAR_WIDTH - PROGRESS_DONE ))
    local PROGRESS_DONE_STR
    PROGRESS_DONE_STR=$(printf "%${PROGRESS_DONE}s")
    local PROGRESS_LEFT_STR
    PROGRESS_LEFT_STR=$(printf "%${PROGRESS_LEFT}s")

    printf "\rProgress : [${PROGRESS_DONE_STR// /#}${PROGRESS_LEFT_STR// /-}] ${PROGRESS}%%"
    if [[ $PROGRESS_LEFT == 0 ]]; then
        echo
    fi
}

# Function to clean up a directory
# Arguments:
#   $1 -- directory name to clean up
function cleanup_directory {
    local dir_name=$1

    if [ "$LOCAL_INSTALL_CLEANUP" != true ]; then
        return 0
    fi

    if [ -z "${log_file}" ] || [ ! -f "${log_file}" ]; then
        log_file=/dev/null
    fi

    if ! rm -drf "$dir_name" >>$log_file 2>&1; then
        echo -e $CLEAR_LINE
        echo -e "${YELLOW}[WARNING] $dir_name cleanup failed ${NC}"
        echo
    fi
}

# Function to call another function and display a progress bar
# Arguments:
#   $1 -- progress (current progress)
#   $2 -- total (total progress)
#   $3 -- function_to_call (the function to be called)
progress_function() {
    local progress=$1
    local total=$2
    local function_to_call=$3

    $function_to_call
    function_return=$?

    if [ $function_return -eq 0 ]; then
        progress_bar $progress $total
    else
        return $function_return
    fi
}

# Function to download a ZIP file and unzip it to a specified folder
# Arguments:
#   $1 -- URL of the ZIP file
#   $2 -- Name of the folder to unzip the contents into
function download_install_debian {
    local url=$1
    local folder_name="$2"
    local folder_path="$LOCAL_INSTALL_DEBIAN_DIRECTORY/$folder_name"

    if [ -z "$url" ]; then
        echo -e "${RED}[ERROR] Error: download_install_debian function first argument missing.${NC}"
        return 1; exit
    fi

    if [ -z "$folder_name" ]; then
        echo -e "${RED}[ERROR] Error: download_install_debian function $1 second argument missing.${NC}"
        echo "Please ensure that the versions.env file is properly loaded."
        return 1
    fi

    if ( [ ! -d "${folder_path}" ] ||
         [ ! -f "${folder_path}/*deb" ] ||
         sudo dpkg -i $folder_path/*.deb ) >>$log_file 2>&1; then
        if ! (mkdir -p "${folder_path}" &&
              wget -O "${folder_path}/download.tar.gz" "$url" &&
              echo "[INFO] Downloaded successfully from $url") >>$log_file 2>&1; then
            echo -e $CLEAR_LINE
            echo -e "${YELLOW}[WARNING] Failed to download $folder_name debian from $url ${NC}"
            echo "Attempting to install from source code as a fallback mechanism."
            return 2
        fi

        if ! (tar -xzvf "${folder_path}/download.tar.gz" -C "$folder_path" &&
            echo "[INFO] Unzipped successfully to $folder_path" &&
            rm -f "${folder_path}/download.tar.gz") >>$log_file 2>&1; then
            echo -e $CLEAR_LINE
            echo -e "${YELLOW}[WARNING] Failed to unzip $folder_name debian to $folder_path ${NC}"
            echo
            return 2
        fi
    fi

    if (sudo dpkg -i $folder_path/*.deb) >>$log_file 2>&1; then
        return 0
    else
        if [ -f "${folder_path}/download.tar.gz" ]; then
            rm -f "${folder_path}/download.tar.gz"
        fi
        echo -e $CLEAR_LINE
        echo -e "${YELLOW}[WARNING] installation of $folder_name debians failed -- trying to install then from source${NC}"
        echo
        return 2
    fi
}

### local installation section

function install_dependencies {
    if ! command -v apt-get >/dev/null; then #TODO add yum support
        echo
        echo -e "${RED}[ERROR] For now only debian distribution are supported by this script ${NC}"
        return 1
    fi

    if ! (sudo apt-get -y --fix-broken install &&
          sudo apt-get update &&
          sudo apt-get install -y wget &&
          wget ${LINK_CUDA_REPO} &&
          sudo dpkg -i cuda-keyring_*.deb &&
          rm cuda-keyring_*.deb) >>$log_file 2>&1; then
        echo
        echo -e "${RED}[ERROR] Nvidia cuda-keyring installation failed ${NC}"
        return 2
    fi

    if ! sudo apt-get install --no-install-recommends -y \
        ca-certificates \
        build-essential \
        libarchive-tools \
        libnuma-dev \
        libjson-c-dev \
        libpcap-dev \
        libgtest-dev \
        libsdl2-dev \
        libsdl2-ttf-dev \
        libssl-dev \
        libtool \
        libx11-dev \
        libx11-xcb-dev \
        libwayland-dev \
        libxcb-dri3-dev \
        libxcb-present-dev \
        libxext-dev \
        libxfixes-dev \
        libxml2-dev \
        libdrm-dev \
        librdmacm-dev \
        zlib1g-dev \
        libelf-dev \
        git \
        cmake \
        meson \
        m4 \
        clang \
        llvm \
        curl \
        g++ \
        nasm \
        autoconf \
        automake \
        pkg-config \
        diffutils \
        gcc \
        gcc-multilib \
        xxd \
        zip \
        python3-pyelftools \
        systemtap-sdt-dev \
        sudo \
        libbsd-dev \
        ocl-icd-opencl-dev \
        libcap2-bin \
        ubuntu-drivers-common \
        libc6-dev \
        cuda-toolkit-12-6 \
        libnvidia-compute-550 \
        python3-pip \
        libfdt-dev 1>>$log_file 2>&1; then
        echo
        echo -e "${RED}[ERROR] apt-get installing dependencies failed ${NC}"
        return 2
    fi
}

function libva_download_build_cleanup {
    if ! (mkdir -p "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/libva" &&
          curl -Lf https://github.com/intel/libva/archive/${LIBVA}.tar.gz | \
          tar -zx --strip-components=1 -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/libva") >>$log_file 2>&1; then
        echo
        echo -e "${RED}[ERROR] libva download failed ${NC}"
        return 2
    fi

    if ! (cd "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/libva" &&
          meson setup build --strip -Dprefix=/usr -Dlibdir=/usr/lib/x86_64-linux-gnu -Ddefault_library=shared &&
          ninja -j"$(nproc)" -C build &&
          sudo meson install -C build &&
          cd -) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] libva build failed ${NC}
        return 2
    fi

    if ! (mkdir -p "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/gmmlib/build" &&
          curl -Lf https://github.com/intel/gmmlib/archive/refs/tags/intel-gmmlib-${GMMLIB}.tar.gz | \
          tar -zx --strip-components=1 -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/gmmlib" &&
          cd "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/gmmlib/build" &&
          cmake -DCMAKE_BUILD_TYPE=Release .. &&
          make -j"$(nproc)" &&
          sudo make install &&
          cd -) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] GMMLIB build failed ${NC}
        return 2
    fi

    if ! (mkdir -p "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/ihd/build" &&
          curl -Lf https://github.com/intel/media-driver/archive/refs/tags/intel-media-${IHD}.tar.gz | \
          tar -zx --strip-components=1 -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/ihd" &&
          cd "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/ihd/build" &&
          cmake -DLIBVA_INSTALL_PATH=/usr/lib/x86_64-linux-gnu -DLIBVA_DRIVERS_PATH=/usr/lib/x86_64-linux-gnu/dri/ .. &&
          make -j"$(nproc)" &&
          sudo make install &&
          cd -) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] IHD build failed ${NC}
        return 2
    fi

    if ! (mkdir -p "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/libvpl/build" &&
          curl -Lf https://github.com/intel/libvpl/archive/refs/tags/v${LIBVPL}.tar.gz | \
          tar -zx --strip-components=1 -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/libvpl" &&
          cd "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/libvpl/build" &&
          cmake -DCMAKE_INSTALL_PREFIX=/usr -DCMAKE_INSTALL_LIBDIR=/usr/lib/x86_64-linux-gnu .. &&
          cmake --build . --config Release &&
          sudo cmake --build . --config Release --target install &&
          cd -) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] LIBVPL build failed ${NC}
        return 2
    fi

    if ! (mkdir -p "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/onevpl/build" &&
          curl -Lf https://github.com/intel/vpl-gpu-rt/archive/refs/tags/intel-onevpl-${ONEVPL}.tar.gz | \
            tar -zx --strip-components=1 -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/onevpl" &&
          cd "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/onevpl/build" &&
          git -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/onevpl" apply "${SCRIPT_DIR}/patches"/onevpl/*.patch &&
          cmake -DCMAKE_INSTALL_PREFIX=/usr -DCMAKE_INSTALL_LIBDIR=/usr/lib/x86_64-linux-gnu .. &&
          make -j${nproc} &&
          sudo make install -j${nproc} &&
          cd -) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] ONEVPL build failed ${NC}
        return 2
    fi

    # This package will block libvapi
    (sudo apt-get purge -y intel-media-va-driver) >>$log_file 2>&1

    if ! (sudo apt-get install -y --reinstall intel-opencl-icd intel-level-zero-gpu level-zero \
          intel-media-va-driver-non-free libmfx1 libmfxgen1 \
          libegl-mesa0 libegl1-mesa libegl1-mesa-dev libgbm1 libgl1-mesa-dev libgl1-mesa-dri \
          libglapi-mesa libgles2-mesa-dev libglx-mesa0 libxatracker2 mesa-va-drivers \
          mesa-vdpau-drivers mesa-vulkan-drivers va-driver-all vainfo hwinfo clinfo) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] Reinstalation of gpu packages failed ${NC}
        return 2
    fi

    cleanup_directory "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/libva"
    cleanup_directory "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/gmmlib"
    cleanup_directory "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/ihd"
    cleanup_directory "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/libvpl"
    cleanup_directory "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/onevpl"
}

function vmaf_download_build_cleanup {
    if ! (mkdir -p "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/vmaf/libvmaf" &&
          curl -Lf https://github.com/Netflix/vmaf/archive/refs/tags/v${LIBVMAF}.tar.gz | \
          tar -zx --strip-components=1 -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/vmaf") >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] VMAF download failed ${NC}
        return 2
    fi

    if ! (cd "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/vmaf/libvmaf" &&
          meson setup \
          --prefix=/usr --libdir=/usr/lib/x86_64-linux-gnu \
          --buildtype release build &&
          ninja -j"$(nproc)" -vC build &&
          sudo ninja -j"$(nproc)" -vC build install &&
          cd -) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] VMAF build failed ${NC}
        return 2
    fi

    cleanup_directory "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/vmaf"
}

function svt_av1_download_build_cleanup {
    if ! (mkdir -p "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/svt-av1/Build" &&
          curl -Lf https://gitlab.com/AOMediaCodec/SVT-AV1/-/archive/v${SVTAV1}/SVT-AV1-v${SVTAV1}.tar.gz | \
          tar -zx --strip-components=1 -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/svt-av1" ) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] SVT-AV1 download failed ${NC}
        return 2
    fi

    if ! (cd "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/svt-av1/Build" &&
          cmake .. -G"Unix Makefiles" -DCMAKE_BUILD_TYPE=Release &&
          make -j "$(nproc)"&&
          sudo make install &&
          cd -) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] SVT-AV1 build failed ${NC}
        return 2
    fi

    cleanup_directory "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/svt-av1"
}

function vulkan_headers_download_build_cleanup {
    if ! (mkdir -p "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/vulkan-headers" &&
      curl -Lf https://github.com/KhronosGroup/Vulkan-Headers/archive/refs/tags/${VULKANSDK}.tar.gz | \
      tar -zx --strip-components=1 -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/vulkan-headers" ) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] Vulkan Headers download failed ${NC}
        return 2
    fi

    if ! (cd ${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/vulkan-headers &&
          cmake -S . -B build/ &&
          sudo cmake --install build --prefix /usr/local &&
          cd - ) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] Vulkan Headers build failed ${NC}
        return 2
    fi

    cleanup_directory "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/vulkan-headers"
}

function xdp_tools_download_build_cleanup {
    if ! (mkdir -p "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/xdp-tools" &&
          curl -Lf https://github.com/xdp-project/xdp-tools/archive/${XDP_VER}.tar.gz | \
          tar -zx --strip-components=1 -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/xdp-tools" &&
          curl -Lf https://github.com/libbpf/libbpf/archive/${BPF_VER}.tar.gz | \
          tar -zx --strip-components=1 -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/xdp-tools/lib/libbpf") >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] XDP-Tools or libbpf download failed ${NC}
        return 2
    fi

    if ! (cd "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/xdp-tools" &&
          ./configure &&
          make &&
          sudo make -j"$(nproc)" install &&
          sudo make -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/xdp-tools/lib/libbpf/src" install)>>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] XDP-Tools or libbpf installation failed ${NC}
        return 2
    fi

    cleanup_directory "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/xdp-tools"
}

function mtl_download {
    if ! (mkdir -p ${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/Media-Transport-Library &&
          curl -Lf https://github.com/OpenVisualCloud/Media-Transport-Library/archive/refs/tags/${MTL_VER}.tar.gz | \
          tar -zx --strip-components=1 -C ${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/Media-Transport-Library ) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] MTL download failed ${NC}
        return 2
    fi
}

function mtl_build {
    if ! (cd ${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/Media-Transport-Library &&
          ./build.sh &&
          cd -) >>$log_file 2>&1; then
          echo
        echo -e ${RED}[ERROR] MTL build script failed ${NC}
        return 2
    fi

    if ! (sudo apt-get -y install make m4 clang llvm zlib1g-dev libelf-dev libpcap-dev libcap-ng-dev gcc-multilib && \
          cd ${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/Media-Transport-Library/manager &&
          meson setup build &&
          meson compile -C build &&
          sudo meson install -C build &&
          cd -) >>$log_file 2>&1; then
          echo
        echo -e ${RED}[ERROR] MTL manager build failed ${NC}
        return 2
    fi
}

function dpdk_download_patch_build {
    if ! (mkdir -p "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/dpdk" &&
          curl -Lf https://github.com/DPDK/dpdk/archive/refs/tags/v${DPDK_VER}.tar.gz | \
          tar -zx --strip-components=1 -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/dpdk") >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] DPDK download failed ${NC}
        return 2
    fi

    if ! (cd "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/dpdk" &&
          git apply ${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/Media-Transport-Library/patches/dpdk/$DPDK_VER/*.patch &&
          cd -) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] Patching DPDK with Media-Transport-Library patches failed ${NC}
        return 2
    fi

    if ! (cd "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/dpdk" &&
          meson build &&
          ninja -C build &&
          sudo ninja -C build install &&
          cd -) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] Dpdk build failed ${NC}
        return 2
    fi
}

function jpegxs_download_build {
    if ! (mkdir -p "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/jpegxs" &&
          curl -Lf https://github.com/OpenVisualCloud/SVT-JPEG-XS/archive/${JPEG_XS_COMMIT_ID}.tar.gz | \
          tar -zx --strip-components=1 -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/jpegxs") >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] JPEG download failed ${NC}
        return 2
    fi

    if ! (cd "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/jpegxs/Build/linux" &&
          ./build.sh install &&
          cp "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/jpegxs/imtl-plugin/kahawai.json" ${SCRIPT_DIR} &&
          cd - ) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] JPEG-XS build failed ${NC}
        return 2
    fi

    if ! (cd "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/jpegxs/imtl-plugin" &&
          ./build.sh &&
          cd - ) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] JPEG-XS imtl plugin build failed ${NC}
        return 2
    fi
}

function ipp_download_build {
    if ! (wget --progress=dot:giga \
          https://registrationcenter-download.intel.com/akdlm/IRC_NAS/046b1402-c5b8-4753-9500-33ffb665123f/l_ipp_oneapi_p_2021.10.1.16_offline.sh) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] IPP download failed ${NC}
        return 2
    fi

    if [ -z "$LOCAL_INSTALL_SKIP_IPP_BUILD" ]; then

        if grep -q "source /opt/intel/oneapi/ipp/latest/env/vars.sh" ~/.bash_profile > /dev/null 2>&1; then
            echo -e $CLEAR_LINE
            echo -e ${YELLOW}"[WARNING] IPP source variables command already present in ~/.bash_profile"${NC}
            echo -e "Skipping IPP installation"
            return 0
        fi

        if ! (chmod +x l_ipp_oneapi_p_2021.10.1.16_offline.sh &&
            sudo ./l_ipp_oneapi_p_2021.10.1.16_offline.sh -a -s --eula accept) >>$log_file 2>&1; then
            echo
            echo -e ${RED}[ERROR] IPP installation failed ${NC}
            echo
            echo "If IPP was previously installed, the build log will show:"
            echo "'Cannot install intel.oneapi.lin.ipp.product. It is already installed.'"
            echo "Please uninstall IPP before retrying or use -i to skip ipp installation"
            return 2
        fi
    fi

    if ! (echo "source /opt/intel/oneapi/ipp/latest/env/vars.sh" | tee -a ~/.bash_profile) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] IPP environment setup failed ${NC}
        return 2
    fi
}


# Depends on ipp_download_build
function vsr_download_build {
    if ! (mkdir -p "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/vsr" &&
          curl -Lf https://github.com/OpenVisualCloud/Video-Super-Resolution-Library/archive/refs/tags/${VSR}.tar.gz | \
          tar -zx --strip-components=1 -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/vsr") >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] VSR download failed ${NC}
        return 2
    fi

    if ! (cd "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/vsr" &&
          . /opt/intel/oneapi/ipp/latest/env/vars.sh &&
          sed -i 's/clan//g' "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/vsr/build.sh" &&
          "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/vsr/build.sh" \
          -DCMAKE_INSTALL_PREFIX="${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/vsr/install" \
          -DENABLE_RAISR_OPENCL=ON) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] VSR build failed ${NC}
        return 2
    fi
}

function mcm_download_build {
    if ! (mkdir -p "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/mcm" &&
          curl -Lf https://github.com/OpenVisualCloud/Media-Communications-Mesh/archive/refs/tags/${MCM_VER}.tar.gz |
          tar -zx --strip-components=1 -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/mcm") >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] MCM download failed ${NC}
        return 2
    fi

    if ! (cd "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/mcm" &&
          cmake -S "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/mcm/sdk" -B "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/mcm/sdk/out" &&
          cmake --build "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/mcm/sdk/out" &&
          sudo cmake --install "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/mcm/sdk/out") >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] MCM build failed ${NC}
        return 2
    fi
}

function ffnvcodec_download_build {
    if ! (mkdir -p "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/nv-codec-headers" &&
          curl -Lf https://github.com/FFmpeg/nv-codec-headers/archive/${FFNVCODED_VER}.tar.gz | \
          tar -zx --strip-components=1 -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/nv-codec-headers") >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] FFmpeg NV-Codec-Headers download failed ${NC}
        return 2
    fi

    if ! (cd "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/nv-codec-headers" &&
          make &&
          sudo make install PREFIX=/usr &&
          cd -) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] FFmpeg NV-Codec-Headers build failed ${NC}
        return 2
    fi
}

function ffmpeg_download_patch_build {
    if ! (mkdir -p "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/ffmpeg" &&
          curl -Lf https://github.com/ffmpeg/ffmpeg/archive/${FFMPEG_COMMIT_ID}.tar.gz | \
          tar -zx --strip-components=1 -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/ffmpeg") >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] FFmpeg download failed ${NC}
        return 2
    fi

    if ! (mkdir -p "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/cartwheel" &&
          curl -Lf https://github.com/intel/cartwheel-ffmpeg/archive/${CARTWHEEL_COMMIT_ID}.tar.gz | \
          tar -zx --strip-components=1 -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/cartwheel) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] Cartwheel patches application failed ${NC}
        return 2
    fi

    if ! (cd "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/ffmpeg &&
          git -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/ffmpeg apply "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/cartwheel/patches/*.patch &&
          cp "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_*.c -rf "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/ffmpeg/libavdevice/ &&
          cp "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_*.h -rf "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/ffmpeg/libavdevice/ &&
          git -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/ffmpeg apply "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/Media-Transport-Library/ecosystem/ffmpeg_plugin/7.0/*.patch &&
          cp "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/jpegxs/ffmpeg-plugin/libsvtjpegxs* "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/ffmpeg/libavcodec/ &&
          git -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/ffmpeg apply --whitespace=fix "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/jpegxs/ffmpeg-plugin/7.0/*.patch &&
          git -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/ffmpeg apply "${SCRIPT_DIR}/patches"/ffmpeg/0001-hwupload_async.diff &&
          git -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/ffmpeg apply "${SCRIPT_DIR}/patches"/ffmpeg/0002-qsv_aligned_malloc.diff &&
          git -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/ffmpeg apply "${SCRIPT_DIR}/patches"/ffmpeg/0003-qsvvpp_async.diff &&
          git -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/ffmpeg apply "${SCRIPT_DIR}/patches"/ffmpeg/0004-filtergraph_async.diff &&
          git -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/ffmpeg apply "${SCRIPT_DIR}/patches"/ffmpeg/0005-ffmpeg_scheduler.diff &&
          git -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/ffmpeg apply -v --whitespace=fix --ignore-space-change "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/mcm/ffmpeg-plugin/7.0/*.patch &&
          cp -f "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/mcm/ffmpeg-plugin/mcm_* "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/ffmpeg/libavdevice/ &&
          cd -) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] FFmpeg patches application failed ${NC}
        return 2
    fi

    if ! (. /opt/intel/oneapi/ipp/latest/env/vars.sh &&
          cd "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"/ffmpeg &&
          ./configure \
            --disable-debug \
            --disable-doc \
            --disable-shared \
            --enable-static \
            --enable-ffprobe \
            --enable-libsvtav1 \
            --enable-libvpl \
            --enable-libvmaf \
            --enable-version3 \
            --enable-libxml2 \
            --enable-mtl \
            --enable-opencl \
            --enable-stripping \
            --enable-vaapi \
            --enable-vulkan \
            --enable-libsvtjpegxs \
            --enable-libipp \
            --enable-mcm \
            --enable-pthreads \
            --enable-ffnvcodec \
            --enable-cuda \
            --enable-cuvid \
            --enable-nvenc \
            --enable-nvdec \
            --enable-cuda-llvm \
            --extra-cflags="-march=native -fopenmp -I${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/vsr/install/include/ -I/opt/intel/oneapi/ipp/latest/include/ipp/ -I/usr/local/cuda/include" \
            --extra-ldflags="-fopenmp -L${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}/vsr/install/lib -L/usr/local/cuda/lib64 -L/usr/lib64 -L/usr/local/lib" \
            --extra-libs='-lraisr -lstdc++ -lippcore -lippvm -lipps -lippi -lpthread -lm -lz -lbsd -lrdmacm -lbpf -lxdp' \
            --enable-cross-compile &&
          make -j"$(nproc)" &&
          sudo make install) >>$log_file 2>&1; then
        echo
        echo -e ${RED}[ERROR] FFmpeg build failed ${NC}
        return 2
    fi
}

# When option -l is tagged the enviroment will be installed from this function
function install_locally {
    UBUNTU_DISTRIBUTION_VERSION=$(grep 'VERSION_ID' /etc/os-release | cut -d '"' -f 2)

    if [ ! -d "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}" ]; then
        mkdir -p "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}" || return 2
    fi


    if git -C "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}" rev-parse --is-inside-work-tree > /dev/null 2>&1; then
        echo
        echo -e ${RED}[ERROR] Git is initialized in "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}" ${NC}
        echo folder for local dependencies cannot have git initialized in it
        return 2
    fi

    center "Installing Intel® Tiber™ Broadcast Suite locally"
    progress_bar 0 100

    if ! [[ "$UBUNTU_DISTRIBUTION_VERSION" == "22.04" || "$UBUNTU_DISTRIBUTION_VERSION" == "24.04" ]]; then
        echo
        echo -e ${RED}[ERROR] currently only supported distribution is 22.04 and 24.04 lts ubuntu ${NC}
        return 2
    fi

    if ! [ -d "${LOCAL_INSTALL_LOG_DIRECTORY}" ]; then
        mkdir -p "${LOCAL_INSTALL_LOG_DIRECTORY}"
    fi

    progress_bar 2 100
    progress_function 15 100 install_dependencies                   || return 3
    progress_function 30 100 vmaf_download_build_cleanup            || return 5
    progress_function 35 100 svt_av1_download_build_cleanup         || return 6
    progress_function 37 100 vulkan_headers_download_build_cleanup  || return 7
    progress_function 40 100 xdp_tools_download_build_cleanup       || return 8

    # THESE REPOSITORIES CAN'T BE DELETED UNTIL FFMPEG IS INSTALLED


    if [[ "$UBUNTU_DISTRIBUTION_VERSION" == "24.04" ]]; then
        if ! (download_install_debian $LINK_DPDK_DEBIAN_v2404_ZIP dpdk &&
              progress_bar 50 100 &&
              download_install_debian $LINK_MTL_DEBIAN_v2404_ZIP mtl &&
              progress_bar 55 100 &&
              download_install_debian $LINK_JPEG_XS_DEBIAN_v2404_ZIP jpegxs &&
              progress_bar 60 100 &&
              download_install_debian $LINK_MCM_DEBIAN_v2404_ZIP mcm &&
              progress_bar 65 100 &&
              download_install_debian $LINK_FFMPEG_DEBIAN_v2404_ZIP ffmpeg); then
            # if installation of debians failed then fallback to installation from source
            progress_bar 49 100
            progress_function 50 100 mtl_download                   || return 9
            progress_function 52 100 mtl_build                      || return 9
            progress_function 55 100 dpdk_download_patch_build      || return 10
            progress_function 60 100 jpegxs_download_build          || return 11
            progress_function 65 100 mcm_download_build             || return 12
            progress_function 70 100 ipp_download_build             || return 13
            progress_function 75 100 vsr_download_build             || return 14
            progress_function 85 100 ffnvcodec_download_build       || return 15
            progress_function 95 100 ffmpeg_download_patch_build    || return 16
        else
            # if installation from debians succeded progress
            progress_function 75 100 ipp_download_build             || return 13
            progress_function 80 100 vsr_download_build             || return 14
            progress_function 95 100 ffnvcodec_download_build       || return 15
        fi

    # on ubuntu 22 we don't have support for dpdk debian
    elif [[ "$UBUNTU_DISTRIBUTION_VERSION" == "22.04" ]]; then
        progress_function 45 100 libva_download_build_cleanup       || return 4
        progress_function 47 100 mtl_download                       || return 9
        progress_function 50 100 dpdk_download_patch_build          || return 10

        if ! (download_install_debian $LINK_MTL_DEBIAN_v2204_ZIP mtl &&
              progress_bar 55 100 &&
              download_install_debian $LINK_JPEG_XS_DEBIAN_v2204_ZIP jpegxs &&
              progress_bar 60 100 &&
              download_install_debian $LINK_MCM_DEBIAN_v2204_ZIP mcm &&
              progress_bar 65 100 &&
              download_install_debian $LINK_FFMPEG_DEBIAN_v2204_ZIP ffmpeg); then
            progress_bar 50 100
            progress_function 55 100 mtl_build                      || return 9
            progress_function 60 100 jpegxs_download_build          || return 11
            progress_function 65 100 mcm_download_build             || return 12
            progress_function 70 100 ipp_download_build             || return 13
            progress_function 75 100 vsr_download_build             || return 14
            progress_function 85 100 ffnvcodec_download_build       || return 15
            progress_function 95 100 ffmpeg_download_patch_build    || return 16

        else
            # if installation from debians succeded progress
            progress_function 75 100 ipp_download_build             || return 13
            progress_function 80 100 vsr_download_build             || return 14
            progress_function 95 100 ffnvcodec_download_build       || return 15
        fi

    fi

    if ! sudo ldconfig; then
        echo -e ${RED}[ERROR] ldconfig failed ${NC}
        return 2
    fi

    cleanup_directory "${LOCAL_INSTALL_DEPENDENCIES_DIRECTORY}"
    cleanup_directory "${LOCAL_INSTALL_DEBIAN_DIRECTORY}"
    progress_bar 100 100
}
### local installation section end



### docker installation

# this should support every distribution
function docker_host_prerequisites {
    if ! (mkdir -p "${HOME}/dpdk" &&
          curl -Lf https://github.com/DPDK/dpdk/archive/refs/tags/v${DPDK_VER}.tar.gz | \
          tar -zx --strip-components=1 -C "${HOME}/dpdk" ); then
        echo
        echo -e ${RED}[ERROR] DPDK download and extraction failed ${NC}
        return 2
    fi

    if [ ! -d "${HOME}/Media-Transport-Library" ] && ! (mkdir -p ${HOME}/Media-Transport-Library &&
          curl -Lf https://github.com/OpenVisualCloud/Media-Transport-Library/archive/refs/tags/${MTL_VER}.tar.gz | \
          tar -zx --strip-components=1 -C ${HOME}/Media-Transport-Library ); then
        echo
        echo -e ${RED}[ERROR] MTL download failed ${NC}
        return 2
    fi

    if ! (cd "${HOME}/dpdk" &&
          git apply "${HOME}/Media-Transport-Library/patches/dpdk/$DPDK_VER"/*.patch &&
          cd .. &&
          rm -rf "${HOME}/Media-Transport-Library"); then
        echo
        echo -e ${RED}[ERROR] Patching DPDK with Media-Transport-Library patches failed ${NC}
        return 2
    fi

    if ! (cd "${HOME}/dpdk" &&
          meson build &&
          ninja -C build &&
          sudo ninja -C build install); then
        echo
        echo -e ${RED}[ERROR] DPDK build and installation failed ${NC}
        return 2
    fi

    cleanup_directory "${HOME}/dpdk"
    cleanup_directory "${HOME}/Media-Transport-Library"
}

function install_in_docker_enviroment {
    ENV_PROXY_ARGS=()
    while IFS='' read -r line; do
        ENV_PROXY_ARGS+=("--build-arg")
        ENV_PROXY_ARGS+=("${line}=${!line}")
    done < <(compgen -e | grep -E "_(proxy|PROXY)")

    IMAGE_CACHE_REGISTRY="${IMAGE_CACHE_REGISTRY:-docker.io}"
    IMAGE_REGISTRY="${IMAGE_REGISTRY:-docker.io}"
    IMAGE_TAG="${IMAGE_TAG:-latest}"
    cat "${VERSIONS_ENVIRONMENT_FILE:-${SCRIPT_DIR}/versions.env}" > "${SCRIPT_DIR}/.temp.env"

    docker buildx build "${ENV_PROXY_ARGS[@]}" \
        --build-arg VERSIONS_ENVIRONMENT_FILE=".temp.env" \
        --build-arg IMAGE_CACHE_REGISTRY="${IMAGE_CACHE_REGISTRY}" \
        -t "${IMAGE_REGISTRY}/tiber-broadcast-suite:${IMAGE_TAG}" \
        -f "${SCRIPT_DIR}/Dockerfile" \
        --target final-stage \
        "${SCRIPT_DIR}"

    docker buildx build "${ENV_PROXY_ARGS[@]}" \
        --build-arg VERSIONS_ENVIRONMENT_FILE=".temp.env" \
        --build-arg IMAGE_CACHE_REGISTRY="${IMAGE_CACHE_REGISTRY}" \
        -t "${IMAGE_REGISTRY}/mtl-manager:${IMAGE_TAG}" \
        -f "${SCRIPT_DIR}/Dockerfile" \
        --target manager-stage \
        "${SCRIPT_DIR}"

    cp -r "${SCRIPT_DIR}/gRPC" "${SCRIPT_DIR}/nmos"

    docker buildx build \
        -t "${IMAGE_REGISTRY}/tiber-broadcast-suite-nmos-node:${IMAGE_TAG}" \
        -f "${SCRIPT_DIR}/nmos/Dockerfile" \
        --target final-stage \
        "${SCRIPT_DIR}/nmos"

    docker tag "${IMAGE_REGISTRY}/tiber-broadcast-suite:${IMAGE_TAG}" video_production_image:latest
    docker tag "${IMAGE_REGISTRY}/mtl-manager:${IMAGE_TAG}" mtl-manager:latest
}
### docker installation end

function display_help {
cat <<- HELPTEXT
Usage: $0 [-l] [-h] [-d DEBIAN_DIRECTORY] [-p DEPENDENCIES_DIRECTORY] [-c CLEANUP] [-i SKIP_IPP]

Options:
  -l    Install locally (bare metal installation) instead of using Docker
  -h    Display this help message
  -d    Specify the directory from where Debian packages are used for local installation
        (default is the $LOCAL_INSTALL_DEBIAN_DIRECTORY directory)
  -p    Specify the directory where locally downloaded and compiled dependencies are placed
        (default is the $LOCAL_INSTALL_DEPENDENCIES_DIRECTORY directory)
  -c    Skip the cleanup process of local dependencies
  -i    Skip the IPP (Intel Performance Primitives) build

Description:
  This script will build / install the Intel® Tiber™ Broadcast Suite by default in Docker.
  It will also build the MtlManager Docker container, which is needed for the Intel® Tiber™ Broadcast Suite.
  If you prefer, you can also install it locally from Debian packages using the -l option.
  Local installation will install all of the components on your machine (bare metal installation).
  It will first search the $LOCAL_INSTALL_DEPENDENCIES_DIRECTORY
  (or other location you can specify with the -d option).
  Then it will try to install all of the Debian packages for Tiber.
  If there are none, it will try to download the Debian packages from the Tiber repo.
  If that also fails, the script will download the source code and install it from source.

Note: The -d, -p, -c, and -i options only work when you install Tiber locally using the -l option.

Return Codes:
  1/2 - General error
  3  - (bare metal installation only) Failed to install dependencies
  4  - (bare metal installation only) Failed to download, patch, build, and clean up libva-dev
  5  - (bare metal installation only) Failed to download, build, and clean up VMAF
  6  - (bare metal installation only) Failed to download, build, and clean up SVT-AV1
  7  - (bare metal installation only) Failed to download, build, and clean up Vulkan headers
  8  - (bare metal installation only) Failed to download, build, and clean up XDP tools
  9  - (bare metal installation only) Failed to download, build Media Transport Library (MTL)
  10 - (bare metal installation only) Failed to download, patch, and build Data Plane Development Kit (DPDK)
  11 - (bare metal installation only) Failed to download, build JPEG XS
  12 - (bare metal installation only) Failed to download, build Media Comunication Mesh (MCM)
  13 - (bare metal installation only) Failed to download, build Intel Performance Primitives (IPP)
  14 - (bare metal installation only) Failed to download, build Video Super Resolution (VSR)
  15 - (bare metal installation only) Failed to download, build FFNVCodec
  16 - (bare metal installation only) Failed to download, patch and build FFmpeg
HELPTEXT

}

#TODO check prerequisites

while getopts "lhid:p:c:" opt; do
    case ${opt} in
        l )
            LOCAL_INSTALL=true
            ;;
        h )
            display_help
            exit 0
            ;;
        d )
            LOCAL_INSTALL_DEBIAN_DIRECTORY=${OPTARG}
            ;;
        p )
            LOCAL_INSTALL_DEPENDENCIES_DIRECTORY=${OPTARG}
            ;;
        c )
            LOCAL_INSTALL_CLEANUP=1
            ;;
        i )
            LOCAL_INSTALL_SKIP_IPP_BUILD=1
            ;;
        \? )
            display_help
            exit 1
            ;;
    esac
done

if ! cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null; then
    echo
    echo -e ${RED}[ERROR] failed change the current directory to the script directory ${NC}
    return 2
fi
# shellcheck source=versions.env
if ! . "${VERSIONS_ENVIRONMENT_FILE}" 2> /dev/null; then
    echo
    echo -e ${RED}[ERROR] failed to source "${VERSIONS_ENVIRONMENT_FILE}"  ${NC}
    return 2
fi

ret=0
if [ "$LOCAL_INSTALL" = true ]; then
    install_locally || ret=$?
else
    {
        install_in_docker_enviroment && \
        docker_host_prerequisites
    } || ret=$?
fi

if [ "$ret" -ne 0 ]; then
    echo
    echo -e ${YELLOW}Please check the "$LOCAL_INSTALL_LOG_DIRECTORY/$LOCAL_INSTALL_LOG_FILE" ${NC}
    echo "For detailed installation steps you can refer to docs/manual_bare_metal_installation_helper.md."
    exit $ret
else
    echo
    echo -e ${GREEN}Intel® Tiber™ Broadcast Suite installed sucessfuly ${NC}
    echo -e ${YELLOW}Please restart your computer ${NC}
    echo
    exit 0
fi
