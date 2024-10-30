# syntax=docker/dockerfile:1

#
# SPDX-License-Identifier: BSD-3-Clause
# Copyright(©) 2024 Intel Corporation
# Intel® Tiber™ Broadcast Suite
#
# build stage

ARG IMAGE_CACHE_REGISTRY=docker.io
ARG IMAGE_NAME=library/ubuntu:24.04@sha256:8a37d68f4f73ebf3d4efafbcf66379bf3728902a8038616808f04e34a9ab63ee

FROM ${IMAGE_CACHE_REGISTRY}/${IMAGE_NAME} AS build-stage

ARG nproc=100
USER root

# common env
ENV \
  nproc=$nproc \
  TZ="Europe/Warsaw" \
  DEBIAN_FRONTEND="noninteractive" \
  PKG_CONFIG_PATH=/usr/lib/pkgconfig:/usr/local/lib/pkgconfig:/usr/lib64/pkgconfig:/usr/local/lib/x86_64-linux-gnu/pkgconfig:/tmp/jpegxs/Build/linux/install/lib/pkgconfig \
  MAKEFLAGS=-j${nproc}

# versions variables are contained in the versions.env
ARG REQUIRED_ENVIRONMENT_VARIABLES="LIBVMAF ONEVPL SVTAV1 VULKANSDK VSR CARTWHEEL_COMMIT_ID FFMPEG_COMMIT_ID XDP_VER BPF_VER MTL_VER MCM_VER JPEG_XS_VER DPDK_VER FFNVCODED_VER LINK_CUDA_REPO FFMPEG_PLUGIN_VER"

ARG \
  LIBVMAF \
  ONEVPL \
  SVTAV1 \
  VULKANSDK \
  VSR \
  CARTWHEEL_COMMIT_ID \
  FFMPEG_COMMIT_ID \
  XDP_VER \
  BPF_VER \
  MTL_VER \
  MCM_VER \
  JPEG_XS_VER \
  DPDK_VER \
  FFNVCODED_VER \
  LINK_CUDA_REPO \
  FFMPEG_PLUGIN_VER

SHELL ["/bin/bash", "-ex", "-o", "pipefail", "-c"]

RUN for var in $REQUIRED_ENVIRONMENT_VARIABLES; do \
      if [ -z "\${!var}" ]; then \
        echo \$var = \${!var}; \
        echo "Error: WRONG BUILD ARGUMENTS SEE docs/build.md "; \
        exit 1; \
      fi; \
    done

# Install dependencies
RUN \
  echo "**** ADD CUDA APT REPO ****" && \
  apt-get update --fix-missing && \
  apt-get install -y wget && \
  wget ${LINK_CUDA_REPO} && \
  dpkg -i cuda-keyring_1.1-1_all.deb && \
  echo "**** INSTALL BUILD PACKAGES ****" && \
  apt-get update --fix-missing && \
  apt-get full-upgrade -y && \
  apt-get install --no-install-recommends -y \
    libigdgmm-dev \
    libva-dev \
    intel-media-va-driver \
    libvpl-dev \
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
    libnvidia-compute-550-server \
    libfdt-dev && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/*

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
    -DCMAKE_INSTALL_LIBDIR=/usr/lib/x86_64-linux-gnu .. && \
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
    https://github.com/OpenVisualCloud/Media-Transport-Library/archive/refs/tags/${MTL_VER}.tar.gz | \
  tar -zx --strip-components=1 -C /tmp/Media-Transport-Library

WORKDIR /tmp/dpdk
RUN \
  echo "**** DOWNLOAD DPDK ****" && \
  curl -Lf \
    https://github.com/DPDK/dpdk/archive/refs/tags/v${DPDK_VER}.tar.gz | \
  tar -zx --strip-components=1 -C /tmp/dpdk

# Download and build the xdp-tools project
WORKDIR /tmp/xdp-tools
RUN curl -Lf https://github.com/xdp-project/xdp-tools/archive/${XDP_VER}.tar.gz | \
      tar -zx --strip-components=1 -C /tmp/xdp-tools && \
    curl -Lf https://github.com/libbpf/libbpf/archive/${BPF_VER}.tar.gz | \
      tar -zx --strip-components=1 -C /tmp/xdp-tools/lib/libbpf && \
    ./configure && \
    make && \
    make install && \
    DESTDIR=/buildout make install && \
    make -C /tmp/xdp-tools/lib/libbpf/src install && \
    DESTDIR=/buildout make -C /tmp/xdp-tools/lib/libbpf/src install

WORKDIR /tmp/dpdk
RUN \
  echo "**** BUILD DPDK ****"  && \
  git apply /tmp/Media-Transport-Library/patches/dpdk/$DPDK_VER/*.patch && \
  meson build && \
  ninja -C build && \
  ninja -C build install

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
  git -C /tmp/ffmpeg/ apply /tmp/Media-Transport-Library/ecosystem/ffmpeg_plugin/${FFMPEG_PLUGIN_VER}/*.patch

RUN \
  echo "**** APPLY JPEG-XS PATCHES ****" && \
  git -C /tmp/ffmpeg apply --whitespace=fix /tmp/patches/jpegxs/*.patch

RUN \
  echo "**** APPLY FFMPEG patches ****" && \
  git -C /tmp/ffmpeg apply /tmp/patches/ffmpeg/hwupload_async.diff && \
  git -C /tmp/ffmpeg apply /tmp/patches/ffmpeg/qsv_aligned_malloc.diff && \
  git -C /tmp/ffmpeg apply /tmp/patches/ffmpeg/qsvvpp_async.diff && \
  git -C /tmp/ffmpeg apply /tmp/patches/ffmpeg/filtergraph_async.diff && \
  git -C /tmp/ffmpeg apply /tmp/patches/ffmpeg/ffmpeg_scheduler.diff

WORKDIR /tmp
RUN \
  echo "**** DOWNLOAD AND INSTALL IPP ****" && \
  no_proxy="" wget --progress=dot:giga \
    https://registrationcenter-download.intel.com/akdlm/IRC_NAS/046b1402-c5b8-4753-9500-33ffb665123f/l_ipp_oneapi_p_2021.10.1.16_offline.sh && \
  chmod +x l_ipp_oneapi_p_2021.10.1.16_offline.sh && \
  ./l_ipp_oneapi_p_2021.10.1.16_offline.sh -a -s --eula accept && \
  echo "source /opt/intel/oneapi/ipp/latest/env/vars.sh" | tee -a ~/.bash_profile

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
  # Clang 18 have bug that break compilation, force compiler to GCC
  sed -i 's/clan//g' build.sh && \
  ./build.sh -DCMAKE_INSTALL_PREFIX="/tmp/vsr/install" -DENABLE_RAISR_OPENCL=ON

# RUN \
#   echo "**** APPLY VIDEO SUPER RESOLUTION PATCHES ****" && \
#   git -C /tmp/ffmpeg apply /tmp/patches/vsr/0001-ffmpeg-raisr-filter.patch && \
#   git -C /tmp/ffmpeg apply /tmp/patches/vsr/0002-libavfilter-raisr_opencl-Add-raisr_opencl-filter.patch && \
#   cp /tmp/vsr/ffmpeg/vf_raisr*.c /tmp/ffmpeg/libavfilter

WORKDIR /tmp/mcm
RUN \
  echo "**** DOWNLOAD MEDIA COMMUNICATIONS MESH ****" && \
  curl -Lf \
    https://github.com/OpenVisualCloud/Media-Communications-Mesh/archive/refs/tags/${MCM_VER}.tar.gz | \
  tar -zx --strip-components=1 -C /tmp/mcm

RUN \
  echo "**** BUILD MEDIA COMMUNICATIONS MESH ****" && \
  cmake -S /tmp/mcm/sdk -B /tmp/mcm/sdk/out && \
  cmake --build /tmp/mcm/sdk/out && \
  cmake --install /tmp/mcm/sdk/out

WORKDIR /tmp/nv-codec-headers
RUN \
  echo "**** DOWNLOAD AND INSTALL FFNVCODED HEADERS ****" && \
  curl -Lf https://github.com/FFmpeg/nv-codec-headers/archive/${FFNVCODED_VER}.tar.gz  | \
    tar -zx --strip-components=1 -C /tmp/nv-codec-headers && \
  make && \
  make install PREFIX=/usr

WORKDIR /tmp/ffmpeg/
RUN \
  echo "**** APPLY MEDIA COMMUNICATIONS MESH PATCHES ****" && \
  git -C /tmp/ffmpeg apply -v --whitespace=fix --ignore-space-change /tmp/mcm/ffmpeg-plugin/${FFMPEG_PLUGIN_VER}/*.patch && \
  cp -f /tmp/mcm/ffmpeg-plugin/mcm_* /tmp/ffmpeg/libavdevice/

RUN \
  echo "**** BUILD FFMPEG ****" && \
  . /opt/intel/oneapi/ipp/latest/env/vars.sh && \
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
    --extra-cflags="-march=native -fopenmp -I/tmp/vsr/install/include/ -I/opt/intel/oneapi/ipp/latest/include/ipp/ -I/usr/local/cuda/include" \
    --extra-ldflags="-fopenmp -L/tmp/vsr/install/lib -L/usr/local/cuda/lib64 -L/usr/lib64 -L/usr/local/lib" \
    --extra-libs='-lraisr -lstdc++ -lippcore -lippvm -lipps -lippi -lpthread -lm -lz -lbsd -lrdmacm -lbpf -lxdp' \
    --enable-cross-compile && \
  make

RUN \
  echo "**** ARRANGE FILES ****" && \
  ldconfig && \
  mkdir -p \
    /buildout/usr/bin \
    /buildout/usr/lib/x86_64-linux-gnu/libmfx-gen \
    /buildout/usr/lib/x86_64-linux-gnu/mfx \
    /buildout/usr/local/lib/vpl \
    /buildout/usr/local/lib/x86_64-linux-gnu/dri \
    /buildout/usr/local/lib/x86_64-linux-gnu/dpdk/pmds-24.0/ \
    /buildout/etc/OpenCL/vendors \
    /buildout/usr/local/etc/ && \
  cp \
    /tmp/ffmpeg/ffmpeg \
    /buildout/usr/bin && \
  cp \
    /tmp/ffmpeg/ffprobe \
    /buildout/usr/bin && \
  cp \
    /tmp/ffmpeg/ffplay \
    /buildout/usr/bin && \
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
    /buildout/usr/local/lib/x86_64-linux-gnu/dri && \
  cp -a \
    /tmp/jpegxs/Build/linux/install/lib/* \
    /buildout/usr/lib/x86_64-linux-gnu/ && \
  cp -a \
    /tmp/jpegxs/imtl-plugin/kahawai.json \
    /buildout/usr/local/etc/jpegxs.json && \
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
    /tmp/mcm/sdk/out/lib/libmcm_dp.so* \
    /buildout/usr/lib/x86_64-linux-gnu/

# ===============================================//
#         Tiber Suite final-stage
# ===============================================//
ARG IMAGE_NAME
ARG IMAGE_CACHE_REGISTRY
FROM ${IMAGE_CACHE_REGISTRY}/${IMAGE_NAME} AS final-stage

LABEL org.opencontainers.image.authors="andrzej.wilczynski@intel.com,milosz.linkiewicz@intel.com"
LABEL org.opencontainers.image.url="https://github.com/OpenVisualCloud/Intel-Tiber-Broadcast-Suite"
LABEL org.opencontainers.image.title="Intel® Tiber™ Broadcast Suite"
LABEL org.opencontainers.image.description="Intel® Tiber™ Broadcast Suite. Open Visual Cloud from Intel® Corporation, collaboration on FFmpeg with plugins on Ubuntu. Release image"
LABEL org.opencontainers.image.documentation="https://github.com/OpenVisualCloud/Intel-Tiber-Broadcast-Suite/tree/main/docs"
LABEL org.opencontainers.image.version="0.9.0"
LABEL org.opencontainers.image.vendor="Intel® Corporation"
LABEL org.opencontainers.image.licenses="BSD 3-Clause License"

ENV \
  DEBIAN_FRONTEND="noninteractive" \
  LIBVA_DRIVERS_PATH="/usr/local/lib/x86_64-linux-gnu/dri" \
  LD_LIBRARY_PATH="/usr/lib64:/usr/local/lib:/usr/local/lib/x86_64-linux-gnu:/usr/lib/x86_64-linux-gnu" \
  NVIDIA_DRIVER_CAPABILITIES="compute,video,utility" \
  NVIDIA_VISIBLE_DEVICES="all"

ENV TZ=Europe/Warsaw
ENV KAHAWAI_CFG_PATH="/usr/local/etc/jpegxs.json"

# Install dependencies
SHELL ["/bin/bash", "-e", "-o", "pipefail", "-c"]
RUN \
  echo "**** INSTALL RUNTIME PACKAGES ****" && \
  apt-get update --fix-missing && \
  apt-get full-upgrade -y && \
  apt-get install --no-install-recommends -y \
    libigdgmm12 \
    libva2 \
    intel-media-va-driver \
    libvpl2 \
    sudo \
    ca-certificates \
    libtool \
    libnuma1 \
    libsdl2-2.0-0 \
    libpcap0.8 \
    libssl3 \
    libxcb-shape0 \
    librdmacm1 \
    libsdl2-ttf-2.0-0 \
    libcap-ng0 \
    libatomic1 \
    intel-opencl-icd \
    opencl-headers \
    ocl-icd-libopencl1 \
    libjson-c5 \
    zlib1g \
    libelf1 \
    libcap2-bin \
    libfdt1 && \
  apt-get remove linux-libc-dev -y && \
  apt-get autoremove -y && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/* &&\
  groupadd -g 2110 vfio && \
  groupadd -g 1001 mtl && \
  groupadd -g 1002 mcm && \
  useradd -m -s /bin/bash -G vfio,mtl,mcm -u 1003 tiber && \
  usermod -aG sudo tiber && \
  mkdir -p /var/run/imtl /var/run/mcm /workspace && \
  chown -R tiber:tiber /var/run/imtl /var/run/mcm /workspace && \
  chmod 775 /var/run/imtl /var/run/mcm /workspace

VOLUME ["/var/run/imtl", "/var/run/mcm", "/workspace"]
COPY --chown=tiber --from=build-stage /buildout/ /

RUN ldconfig

EXPOSE 8001/tcp 8002/tcp
HEALTHCHECK --interval=30s --timeout=5s CMD ps aux | grep "ffmpeg" || exit 1

USER "tiber"

CMD ["--help"]
SHELL ["/bin/bash", "-c"]
ENTRYPOINT ["/usr/bin/ffmpeg"]

# ===============================================//
#          MtlManager stage
# ===============================================//
ARG IMAGE_NAME
ARG IMAGE_CACHE_REGISTRY
FROM ${IMAGE_CACHE_REGISTRY}/${IMAGE_NAME} AS manager-stage

LABEL org.opencontainers.image.authors="andrzej.wilczynski@intel.com,milosz.linkiewicz@intel.com"
LABEL org.opencontainers.image.url="https://github.com/OpenVisualCloud/Intel-Tiber-Broadcast-Suite"
LABEL org.opencontainers.image.title="Intel® MTL Manager"
LABEL org.opencontainers.image.description="Intel® MTL Manager. Open Visual Cloud Media Transport Library Manager required for live software defined broadcast stack optimizations. Ubuntu release image"
LABEL org.opencontainers.image.documentation="https://github.com/OpenVisualCloud/Intel-Tiber-Broadcast-Suite/tree/main/docs"
LABEL org.opencontainers.image.version="1.0.0"
LABEL org.opencontainers.image.vendor="Intel® Corporation"
LABEL org.opencontainers.image.licenses="BSD 3-Clause License"

ENV DEBIAN_FRONTEND="noninteractive"
ENV LD_LIBRARY_PATH="/usr/local/lib"
ENV TZ=Europe/Warsaw

SHELL ["/bin/bash", "-e", "-o", "pipefail", "-c"]
RUN \
  apt-get update --fix-missing && \
  apt-get full-upgrade -y && \
  apt-get install --no-install-recommends -y \
    sudo \
    ca-certificates \
    ethtool \
    libelf1 \
    libfdt1 && \
  apt-get autoremove -y && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/* && \
  groupadd -g 2110 vfio && \
  useradd -m -s /bin/bash -G vfio,root mtl && \
  usermod -aG sudo mtl && \
  mkdir -p /var/run/imtl /usr/local/lib/bpf && \
  chown mtl:root /var/run/imtl

VOLUME ["/var/run/imtl"]
WORKDIR "/home/mtl/"
# USER "mtl"
USER "root"
COPY --chown=mtl --chmod=755 --from=build-stage /usr/local/bin/MtlManager /usr/local/bin/MtlManager
COPY --chown=mtl --chmod=755 --from=build-stage /usr/local/lib/libxdp.so.1 /usr/local/lib
COPY --chown=mtl --chmod=755 --from=build-stage /usr/lib64/libbpf.so.1 /usr/local/lib
COPY --chown=mtl --chmod=755 --from=build-stage /usr/local/lib/bpf/ /usr/local/lib/bpf
COPY --chown=mtl --chmod=755 --from=build-stage /tmp/Media-Transport-Library/script/nicctl.sh /home/mtl/

HEALTHCHECK --interval=30s --timeout=5s CMD ps aux | grep "MtlManager" || exit 1
SHELL ["/bin/bash", "-c"]
ENTRYPOINT ["/usr/local/bin/MtlManager"]
