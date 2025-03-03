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

USER root

# common env
ENV \
  TZ="Europe/Warsaw" \
  DEBIAN_FRONTEND="noninteractive" \
  PKG_CONFIG_PATH=/usr/lib/pkgconfig:/usr/local/lib/pkgconfig:/usr/lib64/pkgconfig:/usr/local/lib/x86_64-linux-gnu/pkgconfig

# versions variables are contained in the versions.env
ARG nproc
ARG VERSIONS_ENVIRONMENT_FILE="versions.env"
ARG REQUIRED_ENVIRONMENT_VARIABLES="LIBVMAF ONEVPL SVTAV1 VULKANSDK VSR CARTWHEEL_COMMIT_ID FFMPEG_COMMIT_ID XDP_VER BPF_VER MTL_VER MCM_VER JPEG_XS_COMMIT_ID DPDK_VER FFNVCODED_VER LINK_CUDA_REPO FFMPEG_PLUGIN_VER"

SHELL ["/bin/bash", "-ex", "-o", "pipefail", "-c"]

COPY "/patches" "/tmp/patches"
COPY "${VERSIONS_ENVIRONMENT_FILE}" "/tmp/versions.env"
RUN echo -e "nproc=${nproc:-$(nproc)}" >> "/tmp/versions.env"
ENV BASH_ENV=/tmp/versions.env

# Check dependencies
RUN for pt in $REQUIRED_ENVIRONMENT_VARIABLES; do \
      if [ -z "${!pt}" ]; then \
        echo "Error: Wrong arguments. See docs/build.md. [${pt} = \"${!pt}\"]"; exit 1;\
      fi; \
    done

# Install dependencies
RUN \
  echo "**** ADD CUDA APT REPO ****" && \
  apt-get update --fix-missing && \
  apt-get install --no-install-recommends -y ca-certificates curl && \
  curl -Lf ${LINK_CUDA_REPO} -o /tmp/cuda-keyring_1.1-1_all.deb && \
  dpkg -i /tmp/cuda-keyring_1.1-1_all.deb && \
  echo "**** INSTALL BUILD PACKAGES ****" && \
  apt-get update --fix-missing && \
  apt-get full-upgrade -y && \
  apt-get install --no-install-recommends -y \
    autoconf \
    automake \
    build-essential \
    clang \
    cmake \
    cuda-toolkit-12-6 \
    diffutils \
    g++ \
    gcc \
    gcc-multilib \
    git \
    libarchive-tools \
    libbsd-dev \
    libc6-dev \
    libcap2-bin \
    libdrm-dev \
    libelf-dev \
    libfdt-dev \
    libgtest-dev \
    libjson-c-dev \
    libnuma-dev \
    libnvidia-compute-550-server \
    libpcap-dev \
    librdmacm-dev \
    libsdl2-dev \
    libsdl2-ttf-dev \
    libssl-dev \
    libtool \
    libwayland-dev \
    libx11-dev \
    libx11-xcb-dev \
    libxcb-dri3-dev \
    libxcb-present-dev \
    libxext-dev \
    libxfixes-dev \
    libxml2-dev \
    llvm \
    m4 \
    meson \
    nasm \
    ocl-icd-opencl-dev \
    pkg-config \
    python3-pyelftools \
    sudo \
    systemtap-sdt-dev \
    ubuntu-drivers-common \
    xxd \
    zip \
    zlib1g-dev && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/* /tmp/cuda-keyring_1.1-1_all.deb

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
    make -j${nproc} && \
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
    make -j${nproc} && \
    make install && \
    strip -d /usr/lib/x86_64-linux-gnu/dri/iHD_drv_video.so

  WORKDIR /tmp/libvpl/build
  RUN \
    echo "**** DOWNLOAD and BUILD LIBVPL ****" && \
    curl -Lf \
      https://github.com/intel/libvpl/archive/refs/tags/v${LIBVPL}.tar.gz | \
      tar -zx --strip-components=1 -C /tmp/libvpl && \
    cmake \
      -DCMAKE_INSTALL_PREFIX=/usr \
      -DCMAKE_INSTALL_LIBDIR=/usr/lib/x86_64-linux-gnu \
      .. && \
    cmake --build . --config Release && \
    cmake --build . --config Release --target install && \
    strip -d /usr/lib/x86_64-linux-gnu/libvpl.so

WORKDIR /tmp
RUN \
  echo "**** DOWNLOAD AND INSTALL IPP ****" && \
  curl -Lf https://registrationcenter-download.intel.com/akdlm/IRC_NAS/046b1402-c5b8-4753-9500-33ffb665123f/l_ipp_oneapi_p_2021.10.1.16_offline.sh -o /tmp/l_ipp_oneapi_p_2021.10.1.16_offline.sh && \
  chmod a+x /tmp/l_ipp_oneapi_p_2021.10.1.16_offline.sh && \
  /tmp/l_ipp_oneapi_p_2021.10.1.16_offline.sh -a -s --eula accept && \
  echo "**** DOWNLOAD AND INSTALL gRPC v1.58 ****" && \
  git clone --branch "v1.58.0" --recurse-submodules --depth 1 --shallow-submodules https://github.com/grpc/grpc /tmp/grpc-source && \
  mkdir -p "/tmp/grpc-source/cmake/build" && \
  cmake -S "/tmp/grpc-source" -B "/tmp/grpc-source/cmake/build" -DgRPC_BUILD_TESTS=OFF -DgRPC_INSTALL=ON && \
  make -C "/tmp/grpc-source/cmake/build" "-j${nproc}" && \
  make -C "/tmp/grpc-source/cmake/build" install && \
  rm -rf /tmp/grpc-source /tmp/l_ipp_oneapi_p_2021.10.1.16_offline.sh

WORKDIR /tmp/onevpl/build
RUN \
  echo "**** DOWNLOAD and PATCH ONEVPL ****" && \
  curl -Lf \
    https://github.com/intel/vpl-gpu-rt/archive/refs/tags/intel-onevpl-${ONEVPL}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/onevpl && \
  git -C /tmp/onevpl apply /tmp/patches/onevpl/*.patch

RUN \
  echo "**** BUILD ONEVPL ****" && \
  cmake \
    -DCMAKE_INSTALL_PREFIX=/usr \
    -DCMAKE_INSTALL_LIBDIR=/usr/lib/x86_64-linux-gnu .. && \
  make -j${nproc} && \
  make install -j${nproc} && \
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
  make -j${nproc} && \
  make install -j${nproc}

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
    make -j${nproc} && \
    make install -j${nproc} && \
    make -C /tmp/xdp-tools/lib/libbpf/src -j${nproc} && \
    make -C /tmp/xdp-tools/lib/libbpf/src install -j${nproc}

WORKDIR /tmp/dpdk
RUN \
  echo "**** BUILD DPDK ****"  && \
  git apply /tmp/Media-Transport-Library/patches/dpdk/$DPDK_VER/*.patch && \
  meson build && \
  ninja -j${nproc} -C build && \
  ninja -j${nproc} -C build install

WORKDIR /tmp/Media-Transport-Library
RUN \
  echo "**** BUILD MTL ****"  && \
  ./build.sh

WORKDIR /tmp/jpegxs/Build/linux
RUN \
  echo "**** DOWNLOAD JPEG-XS ****" && \
  curl -Lf \
    https://github.com/OpenVisualCloud/SVT-JPEG-XS/archive/${JPEG_XS_COMMIT_ID}.tar.gz | \
  tar -zx --strip-components=1 -C /tmp/jpegxs && \
  chmod +x ./build.sh /tmp/jpegxs/imtl-plugin/build.sh

RUN \
  echo "**** BUILD JPEG-XS ****" && \
  ./build.sh release --prefix="/usr/local" install

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
  cp /tmp/Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_* -rf /tmp/ffmpeg/libavdevice/ && \
  patch -d "/tmp/ffmpeg" -p1 -i <(cat "/tmp/Media-Transport-Library/ecosystem/ffmpeg_plugin/${FFMPEG_PLUGIN_VER}/"*.patch) && \
  echo "**** APPLY JPEG-XS PATCHES ****" && \
  cp /tmp/jpegxs/ffmpeg-plugin/libsvtjpegxs* -rf /tmp/ffmpeg/libavcodec/ && \
  patch -d "/tmp/ffmpeg" -p1 -i <(cat "/tmp/jpegxs/ffmpeg-plugin/7.0/"*.patch) && \
  echo "**** APPLY FFMPEG patches ****" && \
  patch -d "/tmp/ffmpeg" -p1 -i <(cat "/tmp/patches/ffmpeg/"*.diff)


WORKDIR /tmp/vsr
# hadolint ignore=SC1091
RUN \
  echo "**** DOWNLOAD AND BUILD VIDEO SUPER RESOLUTION ****" && \
  curl -Lf \
    https://github.com/OpenVisualCloud/Video-Super-Resolution-Library/archive/refs/tags/${VSR}.tar.gz | \
  tar -zx --strip-components=1 -C "/tmp/vsr" && \
  echo "Fix for clang 18 bug that breaks compilation. Force compiler to GCC." && \
  sed -i 's/clan//g' build.sh && \
  . "/opt/intel/oneapi/ipp/latest/env/vars.sh" && \
  ./build.sh -DCMAKE_INSTALL_PREFIX="/usr/local" -DENABLE_RAISR_OPENCL=ON

WORKDIR /tmp/mcm
RUN \
  echo "**** DOWNLOAD MEDIA COMMUNICATIONS MESH ****" && \
  curl -Lf \
    https://github.com/OpenVisualCloud/Media-Communications-Mesh/archive/refs/tags/${MCM_VER}.tar.gz | \
  tar -zx --strip-components=1 -C /tmp/mcm

RUN \
  echo "**** BUILD MEDIA COMMUNICATIONS MESH ****" && \
  cmake -S "/tmp/mcm/sdk" -B "/tmp/mcm/sdk/out" \
    -DCMAKE_BUILD_TYPE="Release" \
    -DCMAKE_INSTALL_PREFIX="/usr/local" && \
  cmake --build "/tmp/mcm/sdk/out" && \
  cmake --install "/tmp/mcm/sdk/out"

WORKDIR /tmp/nv-codec-headers
RUN \
  echo "**** DOWNLOAD AND INSTALL FFNVCODED HEADERS ****" && \
  curl -Lf https://github.com/FFmpeg/nv-codec-headers/archive/${FFNVCODED_VER}.tar.gz  | \
    tar -zx --strip-components=1 -C /tmp/nv-codec-headers && \
  make -j${nproc} && \
  make install -j${nproc} PREFIX=/usr

WORKDIR /tmp/ffmpeg/
RUN \
  echo "**** APPLY MEDIA COMMUNICATIONS MESH PATCHES ****" && \
  patch -d "/tmp/ffmpeg" -p1 -i <(cat "/tmp/mcm/ffmpeg-plugin/${FFMPEG_PLUGIN_VER}/"*.patch) && \
  cp -f "/tmp/mcm/ffmpeg-plugin/mcm_"* "/tmp/ffmpeg/libavdevice/"

# hadolint ignore=SC1091
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
    --extra-cflags="-march=native -fopenmp -I/usr/local/include/ -I/opt/intel/oneapi/ipp/latest/include/ipp/ -I/usr/local/cuda/include" \
    --extra-ldflags="-fopenmp -L/usr/local/cuda/lib64 -L/usr/lib64 -L/usr/local/lib" \
    --extra-libs='-lraisr -lstdc++ -lippcore -lippvm -lipps -lippi -lpthread -lm -lz -lbsd -lrdmacm -lbpf -lxdp' \
    --enable-cross-compile && \
  make -j${nproc}

COPY /gRPC /tmp/gRPC
RUN /tmp/gRPC/compile.sh

WORKDIR /tmp/ffmpeg/

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
  mv \
    /tmp/ffmpeg/ffmpeg \
    /buildout/usr/bin && \
  mv \
    /tmp/gRPC/build/FFmpeg_wrapper_service \
    /buildout/usr/bin && \
  mv \
    /tmp/ffmpeg/ffprobe \
    /buildout/usr/bin && \
  mv \
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
    /tmp/jpegxs/imtl-plugin/kahawai.json \
    /buildout/usr/local/etc/jpegxs.json && \
  cp -a \
    /tmp/vsr/filters* \
    /buildout/ && \
  cp -a \
    /opt/intel/oneapi/ipp/2021.10/lib/libipp*.so.* \
    /buildout/usr/lib/x86_64-linux-gnu/ && \
  cp -a \
    /tmp/Media-Transport-Library/build/lib/libmtl.so* \
    /buildout/usr/lib/x86_64-linux-gnu/

# ===============================================//
#         Tiber Suite final-stage
# ===============================================//
ARG IMAGE_NAME
ARG IMAGE_CACHE_REGISTRY
FROM ${IMAGE_CACHE_REGISTRY}/${IMAGE_NAME} AS final-stage

LABEL org.opencontainers.image.authors="andrzej.wilczynski@intel.com,milosz.linkiewicz@intel.com,dawid.wesierski@intel.com"
LABEL org.opencontainers.image.url="https://github.com/OpenVisualCloud/Intel-Tiber-Broadcast-Suite"
LABEL org.opencontainers.image.title="Intel® Tiber™ Broadcast Suite"
LABEL org.opencontainers.image.description="Intel® Tiber™ Broadcast Suite. Open Visual Cloud from Intel® Corporation, collaboration on FFmpeg with plugins on Ubuntu. Release image"
LABEL org.opencontainers.image.documentation="https://github.com/OpenVisualCloud/Intel-Tiber-Broadcast-Suite/tree/main/docs"
LABEL org.opencontainers.image.version="24.11.0"
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
COPY --chown=tiber --chmod=755 --from=build-stage /usr/lib64/libbpf.so.1 /usr/local/lib

RUN ldconfig

EXPOSE 8001/tcp 8002/tcp
HEALTHCHECK --interval=30s --timeout=5s CMD ps aux | grep "ffmpeg" || exit 1

USER "tiber"

CMD ["--help"]
SHELL ["/bin/bash", "-c"]
ENTRYPOINT ["/usr/bin/FFmpeg_wrapper_service"]

# ===============================================//
#          MtlManager stage
# ===============================================//
ARG IMAGE_NAME
ARG IMAGE_CACHE_REGISTRY
FROM ${IMAGE_CACHE_REGISTRY}/${IMAGE_NAME} AS manager-stage

LABEL org.opencontainers.image.authors="andrzej.wilczynski@intel.com,milosz.linkiewicz@intel.com,dawid.wesierski@intel.com"
LABEL org.opencontainers.image.url="https://github.com/OpenVisualCloud/Intel-Tiber-Broadcast-Suite"
LABEL org.opencontainers.image.title="Intel® MTL Manager"
LABEL org.opencontainers.image.description="Intel® MTL Manager. Open Visual Cloud Media Transport Library Manager required for live software defined broadcast stack optimizations. Ubuntu release image"
LABEL org.opencontainers.image.documentation="https://github.com/OpenVisualCloud/Intel-Tiber-Broadcast-Suite/tree/main/docs"
LABEL org.opencontainers.image.version="24.11.0"
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
