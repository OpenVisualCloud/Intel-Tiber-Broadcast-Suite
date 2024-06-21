# syntax=docker/dockerfile:1

# SPDX-License-Identifier: BSD-3-Clause
# Copyright 2024 Intel Corporation
# Intel® Tiber™ Broadcast Suite

# build stage
ARG IMAGE_CACHE_REGISTRY=docker.io
FROM ${IMAGE_CACHE_REGISTRY}/library/ubuntu:22.04 AS buildstage

ARG nproc=100

# common env
ENV \
  nproc=$nproc \
  DEBIAN_FRONTEND="noninteractive" \
  MAKEFLAGS=-j${nproc}

# versions
ENV \
  GMMLIB=22.3.12 \
  IHD=23.3.5 \
  LIBMFX=22.5.4 \
  LIBVA=2.20.0 \
  LIBVMAF=2.3.1 \
  LIBVPL=2023.3.1 \
  ONEVPL=23.3.4 \
  SVTAV1=1.7.0 \
  VULKANSDK=vulkan-sdk-1.3.268.0 \
  VSR=v23.11 \
  CARTWHEEL_COMMIT_ID=6.1 \
  FFMPEG_COMMIT_ID=n6.1.1 \
  MTL_VER=2f1c2a3be417065a4dc9276e2d7344d768e95118 \
  MCM_VER=9e921f714a3559e78df28c3b4b0160ab7c855582 \
  JPEG_XS_VER=0.9.0 \
  DPDK_VER=23.11

SHELL ["/bin/bash", "-ex", "-o", "pipefail", "-c"]
# Install dependencies
RUN \
  echo "**** INSTALL BUILD PACKAGES ****" && \
  apt-get update --fix-missing && \
  apt-get full-upgrade -y && \
  apt-get install --no-install-recommends -y \
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
    curl \
    g++ \
    nasm \
    autoconf \
    automake \
    pkg-config \
    diffutils \
    gcc \
    xxd \
    wget \
    zip \
    python3-pyelftools \
    systemtap-sdt-dev \
    sudo \
    libbsd-dev \
    ocl-icd-opencl-dev \
    libcap2-bin && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/*

WORKDIR /tmp/libva
RUN \
  echo "**** DOWNLOAD LIBVA ****" && \
  curl -Lf \
    https://github.com/intel/libva/archive/${LIBVA}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/libva
RUN \
  echo "**** BUILD LIBVA ****" && \
  meson setup build --strip -Dprefix=/usr -Dlibdir=/usr/lib/x86_64-linux-gnu -Ddefault_library=shared && \
  ninja -j${nproc} -C build && \
  meson install -C build && \
  strip -d "/usr/lib/x86_64-linux-gnu/libva"*.so

WORKDIR /tmp/gmmlib/build
RUN \
  echo "**** DOWNLOAD and BUILD GMMLIB ****" && \
  curl -Lf \
    https://github.com/intel/gmmlib/archive/refs/tags/intel-gmmlib-${GMMLIB}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/gmmlib && \
  cmake \
    -DCMAKE_BUILD_TYPE=Release .. && \
  make && \
  make install && \
  strip -d /usr/local/lib/libigdgmm.so

WORKDIR /tmp/ihd/build
RUN \
  echo "**** DOWNLOAD and BUILD IHD ****" && \
  curl -Lf \
    https://github.com/intel/media-driver/archive/refs/tags/intel-media-${IHD}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/ihd && \
  cmake \
    -DLIBVA_INSTALL_PATH=/usr/lib/x86_64-linux-gnu \
    -DLIBVA_DRIVERS_PATH=/usr/lib/x86_64-linux-gnu/dri/  .. && \
  make && \
  make install && \
  strip -d /usr/lib/x86_64-linux-gnu/dri/iHD_drv_video.so

WORKDIR /tmp/libvpl/build
RUN \
  echo "**** DOWNLOAD and BUILD LIBVPL ****" && \
  curl -Lf \
    https://github.com/oneapi-src/oneVPL/archive/refs/tags/v${LIBVPL}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/libvpl && \
  cmake \
    -DCMAKE_INSTALL_PREFIX=/usr \
    -DCMAKE_INSTALL_LIBDIR=/usr/lib/x86_64-linux-gnu \
    .. && \
  cmake --build . --config Release && \
  cmake --build . --config Release --target install && \
  strip -d /usr/lib/x86_64-linux-gnu/libvpl.so

COPY /patches /tmp/patches
WORKDIR /tmp/onevpl/build
RUN \
  echo "**** DOWNLOAD and PATCH ONEVPL ****" && \
  curl -Lf \
    https://github.com/oneapi-src/oneVPL-intel-gpu/archive/refs/tags/intel-onevpl-${ONEVPL}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/onevpl && \
  git -C /tmp/onevpl apply /tmp/patches/onevpl/*.patch

RUN \
  echo "**** BUILD ONEVPL ****" && \
  cmake \
    -DCMAKE_INSTALL_PREFIX=/usr \
    -DCMAKE_INSTALL_LIBDIR=/usr/lib/x86_64-linux-gnu \
    .. && \
  make && \
  make install && \
  strip -d /usr/lib/x86_64-linux-gnu/libmfx-gen.so

WORKDIR /tmp/vmaf/libvmaf
RUN \
  echo "**** DOWNLOAD VMAF ****" && \
  curl -Lf \
    https://github.com/Netflix/vmaf/archive/refs/tags/v${LIBVMAF}.tar.gz | \
  tar -zx --strip-components=1 -C /tmp/vmaf && \
  meson setup \
    --prefix=/usr --libdir=/usr/lib/x86_64-linux-gnu \
    --buildtype release build && \
  ninja -j${nproc} -vC build && \
  ninja -j${nproc} -vC build install

WORKDIR /tmp/svt-av1/Build
RUN \
  echo "**** DOWNLOAD SVT-AV1 ****" && \
  curl -Lf \
    https://gitlab.com/AOMediaCodec/SVT-AV1/-/archive/v${SVTAV1}/SVT-AV1-v${SVTAV1}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/svt-av1
RUN \
  echo "**** BUILD SVT-AV1 ****" && \
  cmake .. -G"Unix Makefiles" -DCMAKE_BUILD_TYPE=Release && \
  make && \
  make install

WORKDIR /tmp/vulkan-headers
RUN \
  echo "**** DOWNLOAD VULKAN HEADERS ****" && \
  curl -Lf \
    https://github.com/KhronosGroup/Vulkan-Headers/archive/refs/tags/${VULKANSDK}.tar.gz | \
  tar -zx --strip-components=1 -C /tmp/vulkan-headers
RUN \
  echo "**** BUILD VULKAN HEADERS ****" && \
  cmake -S . -B build/ && \
  cmake --install build --prefix /usr/local

WORKDIR /tmp/Media-Transport-Library
RUN \
  echo "**** DOWNLOAD MTL ****" && \
  curl -Lf \
    https://github.com/OpenVisualCloud/Media-Transport-Library/archive/${MTL_VER}.tar.gz | \
  tar -zx --strip-components=1 -C /tmp/Media-Transport-Library

WORKDIR /tmp/dpdk
RUN \
  echo "**** DOWNLOAD DPDK ****" && \
  curl -Lf \
    https://github.com/DPDK/dpdk/archive/refs/tags/v${DPDK_VER}.tar.gz | \
  tar -zx --strip-components=1 -C /tmp/dpdk
RUN \
  echo "**** BUILD DPDK ****"  && \
  git apply /tmp/Media-Transport-Library/patches/dpdk/$DPDK_VER/*.patch && \
  meson build && \
  ninja -C build && \
  ninja -C build install

# git -C /tmp/onevpl apply /tmp/patches/onevpl/*.patch
WORKDIR /tmp/Media-Transport-Library
RUN \
  echo "**** BUILD MTL ****"  && \
  git -C /tmp/Media-Transport-Library apply /tmp/patches/imtl/*.patch && \
  ./build.sh

WORKDIR /tmp/jpegxs/Build/linux
RUN \
  echo "**** DOWNLOAD JPEG-XS ****" && \
  curl -Lf \
    https://github.com/OpenVisualCloud/SVT-JPEG-XS/archive/refs/tags/v${JPEG_XS_VER}.tar.gz | \
  tar -zx --strip-components=1 -C /tmp/jpegxs && \
  chmod +x ./build.sh /tmp/jpegxs/imtl-plugin/build.sh

RUN \
  echo "**** BUILD JPEG-XS ****" && \
  ./build.sh install --prefix=/tmp/jpegxs/Build/linux/install

ENV LD_LIBRARY_PATH="/tmp/jpegxs/Build/linux/install/lib"
ENV PKG_CONFIG_PATH="/tmp/jpegxs/Build/linux/install/lib/pkgconfig"

WORKDIR /tmp/jpegxs/imtl-plugin
RUN \
  echo "**** BUILD JPEG-XS MTL PLUGIN ****" && \
  ./build.sh

WORKDIR /tmp/ffmpeg
RUN \
  echo "**** DOWNLOAD FFMPEG ****" && \
  curl -Lf \
    https://github.com/ffmpeg/ffmpeg/archive/${FFMPEG_COMMIT_ID}.tar.gz | \
  tar -zx --strip-components=1 -C /tmp/ffmpeg

WORKDIR /tmp/cartwheel
RUN \
  echo "**** APPLY CARTWHEEL PATCHES ****" && \
  curl -Lf \
    https://github.com/intel/cartwheel-ffmpeg/archive/${CARTWHEEL_COMMIT_ID}.tar.gz | \
  tar -zx --strip-components=1 -C /tmp/cartwheel && \
  git -C /tmp/ffmpeg apply /tmp/cartwheel/patches/*.patch

RUN \
  echo "**** APPLY MTL PATCHES ****" && \
  cp /tmp/Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_*.c -rf /tmp/ffmpeg/libavdevice/ && \
  cp /tmp/Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_*.h -rf /tmp/ffmpeg/libavdevice/ && \
  git -C /tmp/ffmpeg/ apply /tmp/Media-Transport-Library/ecosystem/ffmpeg_plugin/6.1/*.patch

RUN \
  echo "**** APPLY JPEG-XS PATCHES ****" && \
  git -C /tmp/ffmpeg apply --whitespace=fix /tmp/jpegxs/ffmpeg-plugin/*.patch

RUN \
  echo "**** APPLY FFMPEG patches ****" && \
  git -C /tmp/ffmpeg apply /tmp/patches/ffmpeg/hwupload_async.diff && \
  git -C /tmp/ffmpeg apply /tmp/patches/ffmpeg/qsv_aligned_malloc.diff && \
  git -C /tmp/ffmpeg apply /tmp/patches/ffmpeg/qsvvpp_async.diff && \
  git -C /tmp/ffmpeg apply /tmp/patches/ffmpeg/filtergraph_async.diff

WORKDIR /tmp
RUN \
  echo "**** DOWNLOAD AND INSTALL IPP ****" && \
  no_proxy="" wget --progress=dot:giga \
    https://registrationcenter-download.intel.com/akdlm/IRC_NAS/046b1402-c5b8-4753-9500-33ffb665123f/l_ipp_oneapi_p_2021.10.1.16_offline.sh && \
  chmod +x l_ipp_oneapi_p_2021.10.1.16_offline.sh && \
  ./l_ipp_oneapi_p_2021.10.1.16_offline.sh -a -s --eula accept && \
  echo "source /opt/intel/oneapi/ipp/latest/env/vars.sh" | tee -a ~/.bash_profile

ENV PKG_CONFIG_PATH="/usr/local/lib/pkgconfig:$PKG_CONFIG_PATH"

WORKDIR /tmp/vsr
RUN \
  echo "**** DOWNLOAD VIDEO SUPER RESOLUTION ****" && \
  curl -Lf \
    https://github.com/OpenVisualCloud/Video-Super-Resolution-Library/archive/refs/tags/${VSR}.tar.gz | \
  tar -zx --strip-components=1 -C /tmp/vsr

RUN \
  echo "**** BUILD VIDEO SUPER RESOLUTION ****" && \
  . /opt/intel/oneapi/ipp/latest/env/vars.sh && \
  git -C /tmp/vsr/ apply /tmp/patches/vsr/0003-missing-header-fix.patch && \
  ./build.sh -DCMAKE_INSTALL_PREFIX="/tmp/vsr/install" -DENABLE_RAISR_OPENCL=ON

RUN \
  echo "**** APPLY VIDEO SUPER RESOLUTION PATCHES ****" && \
  git -C /tmp/ffmpeg apply /tmp/patches/vsr/0001-ffmpeg-raisr-filter.patch && \
  git -C /tmp/ffmpeg apply /tmp/patches/vsr/0002-libavfilter-raisr_opencl-Add-raisr_opencl-filter.patch && \
  cp /tmp/vsr/ffmpeg/vf_raisr*.c /tmp/ffmpeg/libavfilter

WORKDIR /tmp/mcm
RUN \
  echo "**** DOWNLOAD MEDIA COMMUNICATIONS MESH ****" && \
  curl -Lf \
    https://github.com/OpenVisualCloud/Media-Communications-Mesh/archive/${MCM_VER}.tar.gz | \
  tar -zx --strip-components=1 -C /tmp/mcm

RUN \
  echo "**** BUILD MEDIA COMMUNICATIONS MESH ****" && \
  cmake -S /tmp/mcm/sdk -B /tmp/mcm/sdk/out && \
  cmake --build /tmp/mcm/sdk/out && \
  cmake --install /tmp/mcm/sdk/out

WORKDIR /tmp/ffmpeg/
RUN \
  echo "**** APPLY MEDIA COMMUNICATIONS MESH PATCHES ****" && \
  git -C /tmp/ffmpeg apply -v --whitespace=fix --ignore-space-change /tmp/mcm/ffmpeg-plugin/6.1/*.patch && \
  cp -f /tmp/mcm/ffmpeg-plugin/mcm_* /tmp/ffmpeg/libavdevice/

RUN \
  echo "**** BUILD FFMPEG ****" && \
  . /opt/intel/oneapi/ipp/latest/env/vars.sh && \
  ./configure \
    --disable-debug \
    --disable-doc \
    --enable-static \
    --enable-ffprobe \
    --enable-libsvtav1 \
    --enable-libvpl \
    --enable-libvmaf \
    --enable-version3 \
    --enable-libxml2 \
    --enable-mtl \
    --enable-opencl \
    --enable-shared \
    --enable-stripping \
    --enable-vaapi \
    --enable-vulkan \
    --enable-libsvtjpegxs \
    --enable-libipp \
    --enable-mcm \
    --extra-cflags="-march=native -fopenmp -I/tmp/vsr/install/include/ -I/opt/intel/oneapi/ipp/latest/include/ipp/" \
    --extra-ldflags="-fopenmp -L/tmp/vsr/install/lib" \
    --extra-libs='-lraisr -lstdc++ -lippcore -lippvm -lipps -lippi -lm -lz -lbsd -lrdmacm' \
    --enable-cross-compile && \
  make

RUN \
  echo "**** ARRANGE FILES ****" && \
  ldconfig && \
  mkdir -p \
    /buildout/usr/local/bin \
    /buildout/usr/lib/x86_64-linux-gnu/libmfx-gen \
    /buildout/usr/lib/x86_64-linux-gnu/mfx \
    /buildout/usr/local/lib/vpl \
    /buildout/usr/lib/x86_64-linux-gnu/dri \
    /buildout/usr/local/lib/x86_64-linux-gnu/dpdk/pmds-24.0/ \
    /buildout/etc/OpenCL/vendors \
    /buildout/dpdk && \
  cp \
    /tmp/ffmpeg/ffmpeg \
    /buildout/usr/local/bin && \
  cp \
    /tmp/ffmpeg/ffprobe \
    /buildout/usr/local/bin && \
  cp \
    /tmp/ffmpeg/ffplay \
    /buildout/usr/local/bin && \
  cp -a \
    /usr/local/lib/lib*so* \
    /buildout/usr/local/lib/ && \
  cp -a \
    /usr/lib/x86_64-linux-gnu/lib*so* \
    /buildout/usr/lib/x86_64-linux-gnu/ && \
  cp -a \
    /usr/local/lib/x86_64-linux-gnu/lib*so* \
    /buildout/usr/local/lib/x86_64-linux-gnu/ && \
  cp -a \
    /usr/local/lib/x86_64-linux-gnu/dpdk/pmds-24.0/* \
    /buildout/usr/local/lib/x86_64-linux-gnu/dpdk/pmds-24.0/ && \
  cp -a \
    /usr/lib/x86_64-linux-gnu/dri/*.so \
    /buildout/usr/lib/x86_64-linux-gnu/dri/ && \
  cp -a \
    /tmp/ffmpeg/libavdevice/libavdevice* \
    /buildout/usr/lib/x86_64-linux-gnu/ && \
  cp -a \
    /tmp/ffmpeg/libavfilter/libavfilter* \
    /buildout/usr/lib/x86_64-linux-gnu/ && \
  cp -a \
    /tmp/ffmpeg/libavformat/libavformat* \
    /buildout/usr/lib/x86_64-linux-gnu/ && \
  cp -a \
    /tmp/ffmpeg/libavcodec/libavcodec* \
    /buildout/usr/lib/x86_64-linux-gnu/ && \
  cp -a \
    /tmp/ffmpeg/libpostproc/libpostproc* \
    /buildout/usr/lib/x86_64-linux-gnu/ && \
  cp -a \
    /tmp/ffmpeg/libavutil/libavutil* \
    /buildout/usr/lib/x86_64-linux-gnu/ && \
  cp -a \
    /tmp/ffmpeg/libswscale/libswscale* \
    /buildout/usr/lib/x86_64-linux-gnu/ && \
  cp -a \
    /tmp/ffmpeg/libswresample/libswresample* \
    /buildout/usr/lib/x86_64-linux-gnu/ && \
  cp -a \
    /tmp/jpegxs/Build/linux/install/lib/* \
    /buildout/usr/lib/x86_64-linux-gnu/ && \
  cp -a \
    /tmp/jpegxs/imtl-plugin/kahawai.json \
    /buildout/kahawai.json && \
  cp -a \
    /tmp/vsr/install/lib/* \
    /buildout/usr/lib/x86_64-linux-gnu/ && \
  cp -a \
    /tmp/vsr/filters* \
    /buildout/ && \
  cp -a \
    /opt/intel/oneapi/ipp/2021.10/lib/libipp*.so.* \
    /buildout/usr/lib/x86_64-linux-gnu/ && \
  cp -a \
    /tmp/Media-Transport-Library/build/lib/libmtl.so* \
    /buildout/usr/lib/x86_64-linux-gnu/ && \
  cp -a \
    /tmp/dpdk/buildtools/* \
    /buildout/dpdk/ && \
  cp -a \
    /tmp/mcm/sdk/out/lib/libmcm_dp.so* \
    /buildout/usr/lib/x86_64-linux-gnu/

# runtime stage
ARG IMAGE_CACHE_REGISTRY
FROM ${IMAGE_CACHE_REGISTRY}/library/ubuntu:22.04 AS finalstage

LABEL maintainer="andrzej.wilczynski@intel.com,milosz.linkiewicz@intel.com"
LABEL org.opencontainers.image.title="Intel® Tiber™ Broadcast Suite"
LABEL org.opencontainers.image.description="Intel® Tiber™ Broadcast Suite. Open Visual Cloud from Intel® Corporation, collaboration on FFmpeg with plugins on Ubuntu. Release image"
LABEL org.opencontainers.image.version="0.9.0"
LABEL org.opencontainers.image.vendor="Intel® Corporation"
LABEL org.opencontainers.image.licenses="BSD 3-Clause License"

ENV \
  DEBIAN_FRONTEND="noninteractive" \
  LIBVA_DRIVERS_PATH="/usr/local/lib/x86_64-linux-gnu/dri" \
  LD_LIBRARY_PATH="/usr/local/lib:/usr/local/lib/x86_64-linux-gnu/" \
  NVIDIA_DRIVER_CAPABILITIES="compute,video,utility" \
  NVIDIA_VISIBLE_DEVICES="all"

# Install dependencies
SHELL ["/bin/bash", "-e", "-o", "pipefail", "-c"]
RUN \
  echo "**** INSTALL RUNTIME PACKAGES ****" && \
  apt-get update --fix-missing && \
  apt-get full-upgrade -y && \
  apt-get install --no-install-recommends -y \
    sudo \
    ca-certificates \
    libtool \
    libnuma1 \
    libsdl2-2.0-0 \
    libpcap0.8 \
    libssl3 \
    libxcb-shape0 \
    librdmacm1 \
    ocl-icd-libopencl1 \
    libjson-c5 \
    zlib1g \
    libelf1 \
    libcap2-bin && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/* &&\
  groupadd -g 2110 vfio && \
  groupadd -g 1001 imtl && \
  groupadd -g 1002 mcm && \
  useradd -m -s /bin/bash -G vfio,imtl,mcm -u 1003 ffmpeg-vpp && \
  usermod -aG sudo ffmpeg-vpp

COPY --chown=ffmpeg-vpp --from=buildstage /buildout/ /

RUN \
  echo "**** ENABLE DPDK ****" && \
  chmod +x dpdk/symlink-drivers-solibs.sh && \
  ./dpdk/symlink-drivers-solibs.sh lib/x86_64-linux-gnu dpdk/pmds-24.0 && \
  ldconfig

SHELL ["/bin/bash", "-c"]
ENTRYPOINT ["/usr/local/bin/ffmpeg"]

HEALTHCHECK --interval=30s --timeout=5s \
  CMD ps aux | grep "ffmpeg" || exit 1
