# syntax=docker/dockerfile:1

# build stage
FROM ubuntu:22.04 as buildstage

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
  X265=3.5 \
  VSR=v23.11 \
  CARTWHEEL_COMMIT_ID=6.1 \
  FFMPEG_COMMIT_ID=n6.1.1 \
  MTL_VER=95673e279cf37e22e664b8b921b7da950976008b \
  DPDK_VER=23.11

# Install dependencies
RUN \
  echo "**** INSTALL BUILD PACKAGES ****" && \
  apt-get update && \
  apt-get install -y \
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
  git \
  cmake \
  meson \
  curl \
  g++ \
  nasm \
  autoconf \
  automake \
  autoconf \
  automake \
  pkg-config \
  bzip2 \
  cmake \
  curl \
  diffutils \
  g++ \
  gcc \
  git \
  xxd \
  wget \
  zip \
  python3-pyelftools \
  systemtap-sdt-dev \
  sudo

RUN \
  echo "**** DOWNLOAD LIBVA ****" && \
  mkdir -p /tmp/libva && \
  curl -Lf \
    https://github.com/intel/libva/archive/${LIBVA}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/libva
RUN \
  echo "**** BUILD LIBVA ****" && \
  cd /tmp/libva && \
  ./autogen.sh && \
  ./configure \
    --disable-static \
    --enable-shared && \
  make && \
  make install && \
  strip -d \
    /usr/local/lib/libva.so \
    /usr/local/lib/libva-drm.so \
    /usr/local/lib/libva-glx.so \
    /usr/local/lib/libva-wayland.so \
    /usr/local/lib/libva-x11.so

RUN \
  echo "**** DOWNLOAD GMMLIB ****" && \
  mkdir -p /tmp/gmmlib && \
  curl -Lf \
    https://github.com/intel/gmmlib/archive/refs/tags/intel-gmmlib-${GMMLIB}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/gmmlib
RUN \
  echo "**** BUILD GMMLIB ****" && \
  mkdir -p /tmp/gmmlib/build && \
  cd /tmp/gmmlib/build && \
  cmake \
    -DCMAKE_BUILD_TYPE=Release \
    .. && \
  make && \
  make install && \
  strip -d /usr/local/lib/libigdgmm.so

RUN \
  echo "**** DOWNLOAD IHD ****" && \
  mkdir -p /tmp/ihd && \
  curl -Lf \
    https://github.com/intel/media-driver/archive/refs/tags/intel-media-${IHD}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/ihd
RUN \
  echo "**** BUILD IHD ****" && \
  mkdir -p /tmp/ihd/build && \
  cd /tmp/ihd/build && \
  cmake \
    -DLIBVA_DRIVERS_PATH=/usr/lib/x86_64-linux-gnu/dri/ \
    .. && \
  make && \
  make install && \
  strip -d /usr/lib/x86_64-linux-gnu/dri/iHD_drv_video.so

RUN \
  echo "**** DOWNLOAD LIBVPL ****" && \
  mkdir -p /tmp/libvpl && \
  curl -Lf \
    https://github.com/oneapi-src/oneVPL/archive/refs/tags/v${LIBVPL}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/libvpl
RUN \
  echo "**** BUILD LIBVPL ****" && \
  mkdir -p /tmp/libvpl/build && \
  cd /tmp/libvpl/build && \
  cmake .. && \
  cmake --build . --config Release && \
  cmake --build . --config Release --target install && \
  strip -d /usr/local/lib/libvpl.so

RUN \
  echo "**** DOWNLOAD ONEVPL ****" && \
  mkdir -p /tmp/onevpl && \
  curl -Lf \
    https://github.com/oneapi-src/oneVPL-intel-gpu/archive/refs/tags/intel-onevpl-${ONEVPL}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/onevpl

COPY \
  /patches/onevpl/*.patch /

RUN \
  echo "**** FFMPEG ONEVPL PATCH ****" && \
  cd /tmp/onevpl && \
  git apply /*.patch
RUN \
  echo "**** BUILD ONEVPL ****" && \
  mkdir -p /tmp/onevpl/build && \
  cd /tmp/onevpl/build && \
  cmake \
    -DCMAKE_INSTALL_PREFIX=/usr \
    -DCMAKE_INSTALL_LIBDIR=/usr/local/lib \
    .. && \
  make && \
  make install && \
  strip -d /usr/local/lib/libmfx-gen.so

RUN \
  echo "**** DOWNLOAD LIBMFX ****" && \
  mkdir -p /tmp/libmfx && \
  curl -Lf \
    https://github.com/Intel-Media-SDK/MediaSDK/archive/refs/tags/intel-mediasdk-${LIBMFX}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/libmfx
RUN \
  echo "**** BUILD LIBMFX ****" && \
  mkdir -p /tmp/libmfx/build && \
  cd /tmp/libmfx/build && \
  cmake \
    -DCMAKE_INSTALL_PREFIX=/usr \
    -DCMAKE_INSTALL_LIBDIR=/usr/local/lib \
    -DBUILD_SAMPLES=OFF \
    -DENABLE_X11_DRI3=ON \
    -DBUILD_DISPATCHER=OFF \
    -DBUILD_TUTORIALS=OFF \
    .. && \
  make && \
  make install && \
  strip -d \
    /usr/local/lib/libmfxhw64.so \
    /usr/local/lib/mfx/libmfx_*.so

RUN \
  echo "**** DOWNLOAD VMAF ****" && \
  mkdir -p /tmp/vmaf && \
  curl -Lf \
    https://github.com/Netflix/vmaf/archive/refs/tags/v${LIBVMAF}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/vmaf
RUN \
  echo "**** BUILD VMAF ****" && \
  cd /tmp/vmaf/libvmaf && \
  meson setup \
    --prefix=/usr --libdir=/usr/local/lib \
    --buildtype release \
    build && \
  ninja -vC build && \
  ninja -vC build install

RUN \
  echo "**** DOWNLOAD SVT-AV1 ****" && \
  mkdir -p /tmp/svt-av1 && \
  curl -Lf \
    https://gitlab.com/AOMediaCodec/SVT-AV1/-/archive/v${SVTAV1}/SVT-AV1-v${SVTAV1}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/svt-av1
RUN \
  echo "**** BUILD SVT-AV1 ****" && \
  cd /tmp/svt-av1/Build && \
  cmake .. -G"Unix Makefiles" -DCMAKE_BUILD_TYPE=Release && \
  make && \
  make install

RUN \
  echo "**** DOWNLOAD VULKAN HEADERS ****" && \
  mkdir -p /tmp/vulkan-headers && \
  curl -Lf \
  https://github.com/KhronosGroup/Vulkan-Headers/archive/refs/tags/${VULKANSDK}.tar.gz | \
  tar -zx --strip-components=1 -C /tmp/vulkan-headers
RUN \
  echo "**** BUILD VULKAN HEADERS ****" && \
  cd /tmp/vulkan-headers && \
  cmake -S . -B build/ && \
  cmake --install build --prefix /usr/local

RUN \
  echo "**** DOWNLOAD x264 ****" && \
  mkdir -p /tmp/x264 && \
  curl -Lf \
    https://code.videolan.org/videolan/x264/-/archive/master/x264-stable.tar.bz2 | \
    tar -jx --strip-components=1 -C /tmp/x264
RUN \
  echo "**** BUILD x264 ****" && \
  cd /tmp/x264 && \
  ./configure \
    --disable-cli \
    --disable-static \
    --enable-pic \
    --enable-shared && \
  make && \
  make install

RUN \
  echo "**** DOWNLOAD x265 ****" && \
  mkdir -p /tmp/x265 && \
  curl -Lf \
    https://bitbucket.org/multicoreware/x265_git/downloads/x265_${X265}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/x265
RUN \
  echo "**** BUILD x265 ****" && \
  cd /tmp/x265/build/linux && \
  ./multilib.sh && \
  make -C 8bit install

RUN \
  echo "**** DOWNLOAD MTL ****" && \
  mkdir -p /tmp/Media-Transport-Library && \
  curl -LO \
  https://github.com/OpenVisualCloud/Media-Transport-Library/archive/${MTL_VER}.zip && \
  unzip ${MTL_VER}.zip -d /tmp/Media-Transport-Library && \
  mv /tmp/Media-Transport-Library/Media-Transport-Library-${MTL_VER}/* /tmp/Media-Transport-Library && \
  rm -rf /tmp/Media-Transport-Library/Media-Transport-Library-${MTL_VER} && \
  mkdir /tmp/Media-Transport-Library/patches/video_production_image

COPY \
  /patches/imtl/* /tmp/Media-Transport-Library/patches/video_production_image/

RUN \
  echo "**** DOWNLOAD DPDK ****" && \
  mkdir -p /tmp/dpdk && \
  curl -Lf \
    https://github.com/DPDK/dpdk/archive/refs/tags/v${DPDK_VER}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/dpdk

RUN \
  echo "**** BUILD DPDK ****"  && \
  cd /tmp/dpdk && \
  git apply /tmp/Media-Transport-Library/patches/dpdk/$DPDK_VER/*.patch && \
  meson build && \
  ninja -C build && \
  sudo ninja -C build install && \
  cd ..

RUN \
  echo "**** BUILD MTL ****"  && \
  cd /tmp/Media-Transport-Library && \
  git apply /tmp/Media-Transport-Library/patches/video_production_image/*.patch && \
  ./build.sh && \
  cd ..

RUN \
  echo "**** DOWNLOAD JPEG-XS ****" && \
  mkdir -p /tmp/jpegxs && \
  curl -LO \
  https://github.com/OpenVisualCloud/SVT-JPEG-XS/archive/refs/heads/main.zip && \
  unzip main.zip -d /tmp/jpegxs && \
  mv /tmp/jpegxs/SVT-JPEG-XS-main/* /tmp/jpegxs && \
  rm -rf /tmp/jpegxs/SVT-JPEG-XS-main

RUN \
  echo "**** BUILD JPEG-XS ****" && \
  mkdir /tmp/jpegxs/Build/linux/install && \
  cd /tmp/jpegxs/Build/linux && \
  ./build.sh install --prefix=/tmp/jpegxs/Build/linux/install

RUN \
  echo "**** BUILD JPEG-XS MTL PLUGIN ****" && \
  cd /tmp/jpegxs/imtl-plugin && \
  ./build.sh --prefix=/tmp/jpegxs/Build/linux/install

ENV \
  LD_LIBRARY_PATH="/tmp/jpegxs/Build/linux/install/lib:${LD_LIBRARY_PATH}"
ENV \
  PKG_CONFIG_PATH="/tmp/jpegxs/Build/linux/install/lib/pkgconfig:${PKG_CONFIG_PATH}"

RUN \
  echo "**** DOWNLOAD FFMPEG ****" && \
  mkdir -p /tmp/ffmpeg && \
  curl -Lf \
    https://github.com/ffmpeg/ffmpeg/archive/${FFMPEG_COMMIT_ID}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/ffmpeg
RUN \
  echo "**** APPLY CARTWHEEL PATCHES ****" && \
  mkdir -p /tmp/cartwheel && \
  curl -Lf \
    https://github.com/intel/cartwheel-ffmpeg/archive/${CARTWHEEL_COMMIT_ID}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/cartwheel && \
  cd /tmp/ffmpeg && \
  git apply /tmp/cartwheel/patches/*.patch

COPY \
  patches/ffmpeg/*.diff /ffmpeg_patches/

RUN \
  echo "**** APPLY MTL PATCHES ****" && \
  cp /tmp/Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_*.c -rf /tmp/ffmpeg/libavdevice/ && \
  cp /tmp/Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_*.h -rf /tmp/ffmpeg/libavdevice/ && \
  cd /tmp/ffmpeg/ && \
  git apply /tmp/Media-Transport-Library/ecosystem/ffmpeg_plugin/6.1/*.patch

RUN \
  echo "**** APPLY JPEG-XS PATCHES ****" && \
  cd /tmp/ffmpeg/ && \
  git apply --whitespace=fix /tmp/jpegxs/ffmpeg-plugin/*.patch

RUN \
  echo "**** APPLY FFMPEG patches ****" && \
  cd /tmp/ffmpeg && \
  git apply /ffmpeg_patches/hwupload_async.diff && \
  git apply /ffmpeg_patches/qsv_aligned_malloc.diff && \
  git apply /ffmpeg_patches/qsvvpp_async.diff && \
  git apply /ffmpeg_patches/filtergraph_async.diff

RUN \
  echo "**** DOWNLOAD AND INSTALL IPP ****" && \
  wget https://registrationcenter-download.intel.com/akdlm/IRC_NAS/046b1402-c5b8-4753-9500-33ffb665123f/l_ipp_oneapi_p_2021.10.1.16_offline.sh && \
  chmod +x l_ipp_oneapi_p_2021.10.1.16_offline.sh && \
  ./l_ipp_oneapi_p_2021.10.1.16_offline.sh -a -s --eula accept && \
  echo "source /opt/intel/oneapi/ipp/latest/env/vars.sh" | tee -a ~/.bash_profile

ENV \
  PKG_CONFIG_PATH="/usr/local/lib/pkgconfig:$PKG_CONFIG_PATH"

RUN \
  echo "**** DOWNLOAD VIDEO SUPER RESOLUTION (VSR) ****" && \
  mkdir -p /tmp/vsr && \
  curl -Lf \
    https://github.com/OpenVisualCloud/Video-Super-Resolution-Library/archive/refs/tags/${VSR}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/vsr

RUN apt-get install -y  ocl-icd-opencl-dev

COPY \
  patches/vsr/*.patch /vsr_patches/

RUN \
  echo "**** BUILD VIDEO SUPER RESOLUTION (VSR) ****" && \
  cd /tmp/vsr/ && . /opt/intel/oneapi/ipp/latest/env/vars.sh && \
  git apply /vsr_patches/0003-missing-header-fix.patch && \
  ./build.sh -DCMAKE_INSTALL_PREFIX="$PWD/install" -DENABLE_RAISR_OPENCL=ON

RUN \
  echo "**** APPLY VSR PATCHES ****" && \
  cd tmp/ffmpeg/ && \
  git apply /vsr_patches/0001-ffmpeg-raisr-filter.patch && \
  git apply /vsr_patches/0002-libavfilter-raisr_opencl-Add-raisr_opencl-filter.patch && \
  cp /tmp/vsr/ffmpeg/vf_raisr*.c /tmp/ffmpeg/libavfilter

RUN \
  echo "**** BUILD FFMPEG ****" && \
  . /opt/intel/oneapi/ipp/latest/env/vars.sh && \
  cd /tmp/ffmpeg && \
    ./configure \
    --disable-debug \
    --disable-doc \
    --enable-static \
    --enable-ffprobe \
    --enable-gpl \
    --enable-libsvtav1 \
    --enable-libvpl \
    --enable-libvmaf \
    --enable-version3 \
    --enable-libx264 \
    --enable-libx265 \
    --enable-libxml2 \
    --enable-mtl \
    --enable-opencl \
    --enable-shared \
    --enable-stripping \
    --enable-vaapi \
    --enable-vulkan \
    --enable-libsvtjpegxs \
    --enable-libipp \
    --extra-cflags="-fopenmp -I/tmp/vsr/install/include/ -I/opt/intel/oneapi/ipp/latest/include/ipp/" \
    --extra-ldflags="-fopenmp -L/tmp/vsr/install/lib" \
    --extra-libs='-lraisr -lstdc++ -lippcore -lippvm -lipps -lippi -lm' \
    --enable-cross-compile && \
  make

RUN \
  echo "**** ARRANGE FILES ****" && \
  sudo ldconfig && \
  sudo mkdir -p \
    /buildout/usr/local/bin \
    /buildout/usr/local/lib/libmfx-gen \
    /buildout/usr/local/lib/mfx \
    /buildout/usr/local/lib/vpl \
    /buildout/usr/local/lib/x86_64-linux-gnu/dri \
    /buildout/usr/local/lib/x86_64-linux-gnu/dpdk/pmds-24.0/ \
    /buildout/etc/OpenCL/vendors \
    /buildout/dpdk && \
  sudo cp \
    /tmp/ffmpeg/ffmpeg \
    /buildout/usr/local/bin && \
  sudo cp \
    /tmp/ffmpeg/ffprobe \
    /buildout/usr/local/bin && \
  sudo cp \
    /tmp/ffmpeg/ffplay \
    /buildout/usr/local/bin && \
  sudo cp -a \
    /usr/local/lib/lib*so* \
    /buildout/usr/local/lib/ && \
  sudo cp -a \
    /usr/local/lib/libmfx-gen/*.so \
    /buildout/usr/local/lib/libmfx-gen/ && \
  sudo cp -a \
    /usr/local/lib/mfx/*.so \
    /buildout/usr/local/lib/mfx/ && \
  sudo cp -a \
    /usr/local/lib/x86_64-linux-gnu/lib*so* \
    /buildout/usr/local/lib/x86_64-linux-gnu/ && \
  sudo cp -a \
    /usr/local/lib/x86_64-linux-gnu/dpdk/pmds-24.0/* \
    /buildout/usr/local/lib/x86_64-linux-gnu/dpdk/pmds-24.0/ && \
  sudo cp -a \
    /usr/lib/x86_64-linux-gnu/dri/*.so \
    /buildout/usr/local/lib/x86_64-linux-gnu/dri/ && \
  sudo cp -a \
    /tmp/ffmpeg/libavdevice/* \
    /buildout/usr/local/lib/ && \
  sudo cp -a \
    /tmp/ffmpeg/libavfilter/* \
    /buildout/usr/local/lib/ && \
  sudo cp -a \
    /tmp/ffmpeg/libavformat/* \
    /buildout/usr/local/lib/ && \
  sudo cp -a \
    /tmp/ffmpeg/libavcodec/* \
    /buildout/usr/local/lib/ && \
  sudo cp -a \
    /tmp/ffmpeg/libpostproc/* \
    /buildout/usr/local/lib/ && \
  sudo cp -a \
    /tmp/ffmpeg/libavutil/* \
    /buildout/usr/local/lib/ && \
  sudo cp -a \
    /tmp/ffmpeg/libswscale/* \
    /buildout/usr/local/lib/ && \
  sudo cp -a \
    /tmp/ffmpeg/libswresample/* \
    /buildout/usr/local/lib/ && \
  sudo cp -a \
    /tmp/jpegxs/Build/linux/install/lib/* \
    /buildout/usr/local/lib/ && \
  sudo cp -a \
    /tmp/vsr/install/lib/* \
    /buildout/usr/local/lib/ && \
  sudo cp -a \
    /tmp/vsr/filters* \
    /buildout/ && \
  sudo cp -a \
    /opt/intel/oneapi/ipp/2021.10/lib/libipp* \
    /buildout/usr/local/lib/ && \
  sudo cp -a \
    /tmp/Media-Transport-Library/build/lib/libmtl.so* \
    /buildout/usr/local/lib/ && \
  sudo cp -a \
    /tmp/dpdk/buildtools/* \
    /buildout/dpdk/ && \
  sudo echo \
    'libnvidia-opencl.so.1' | sudo tee \
    /buildout/etc/OpenCL/vendors/nvidia.icd

# runtime stage
FROM ubuntu:22.04 as finalstage

# Add files from binstage
COPY \
  --from=buildstage /buildout/ /

# set version label
ARG BUILD_DATE
ARG VERSION

ARG DEBIAN_FRONTEND="noninteractive"

# hardware env
ENV \
  LIBVA_DRIVERS_PATH="/usr/local/lib/x86_64-linux-gnu/dri" \
  LD_LIBRARY_PATH="/usr/local/lib" \
  NVIDIA_DRIVER_CAPABILITIES="compute,video,utility" \
  NVIDIA_VISIBLE_DEVICES="all"

RUN \
  echo "**** INSTALL RUNTIME PACKAGES ****" && \
  apt-get update -y && \
  apt-get install -y \
  meson \
  python3-pyelftools \
  libnuma-dev \
  sudo \
  autoconf \
  libtool \
  libsdl2-dev \
  libpcap-dev \
  libssl-dev \
  libxcb-shape0 \
  intel-opencl-icd \
  opencl-headers \
  ocl-icd-opencl-dev \
  libjson-c-dev

RUN \
  echo "**** ENABLE DPDK ****" && \
  chmod +x dpdk/symlink-drivers-solibs.sh && \
  ./dpdk/symlink-drivers-solibs.sh lib/x86_64-linux-gnu dpdk/pmds-24.0

RUN \
  mkdir /hugetlbfs

ENTRYPOINT ["./usr/local/bin/ffmpeg"]
