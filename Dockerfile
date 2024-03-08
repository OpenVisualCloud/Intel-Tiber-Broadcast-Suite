# syntax=docker/dockerfile:1

# build stage
FROM ubuntu as buildstage

# set version label
ARG FFMPEG_VERSION

# common env
ENV \
  DEBIAN_FRONTEND="noninteractive" \
  MAKEFLAGS="-j100"

# versions
ENV \
  AOM=v3.7.1 \
  FDKAAC=2.0.2 \
  FFMPEG_HARD=6.1 \
  FONTCONFIG=2.14.2 \
  FREETYPE=2.13.2 \
  FRIBIDI=1.0.13 \
  GMMLIB=22.3.12 \
  IHD=23.3.5 \
  KVAZAAR=2.2.0 \
  LAME=3.100 \
  LIBASS=0.17.1 \
  LIBDRM=2.4.118 \
  LIBMFX=22.5.4 \
  LIBVA=2.20.0 \
  LIBVDPAU=1.5 \
  LIBVIDSTAB=1.1.1 \
  LIBVMAF=2.3.1 \
  LIBVPL=2023.3.1 \
  NVCODEC=n12.1.14.0 \
  OGG=1.3.5 \
  ONEVPL=23.3.4 \
  OPENCOREAMR=0.1.6 \
  OPENJPEG=2.5.0 \
  OPUS=1.4 \
  SHADERC=v2023.7 \
  SVTAV1=1.7.0 \
  THEORA=1.1.1 \
  VORBIS=1.3.7 \
  VPX=1.13.1 \
  VULKANSDK=vulkan-sdk-1.3.268.0 \
  WEBP=1.3.2 \
  X265=3.5 \
  XVID=1.3.7

RUN \
  echo "**** install build packages ****" && \
  apt-get update && \ 
  apt-get install -y \
    autoconf \
    automake \
    bzip2 \
    cmake \
    curl \
    diffutils \
    doxygen \
    g++ \
    gcc \
    git \
    gperf \
    i965-va-driver-shaders \
    libexpat1-dev \
    libgcc-10-dev \
    libgomp1 \
    libharfbuzz-dev \
    libpciaccess-dev \
    libssl-dev \
    libtool \
    libv4l-dev \
    libwayland-dev \
    libx11-dev \
    libx11-xcb-dev \
    libxcb-dri3-dev \
    libxcb-present-dev \
    libxext-dev \
    libxfixes-dev \
    libxml2-dev \
    make \
    nasm \
    ninja-build \
    ocl-icd-opencl-dev \
    perl \
    pkg-config \
    python3-pip \
    python3-venv \
    wayland-protocols \
    x11proto-xext-dev \
    xserver-xorg-dev \
    xxd \
    yasm \
    zlib1g-dev \
    sudo && \
  python3 -m venv /lsiopy && \
  pip install -U --no-cache-dir \
    pip \
    setuptools \
    wheel && \
  pip install --no-cache-dir meson cmake

# compile 3rd party libs

RUN \
  echo "**** grabbing aom ****" && \
  mkdir -p /tmp/aom && \
  git clone \
    --branch ${AOM} \
    --depth 1 https://aomedia.googlesource.com/aom \
    /tmp/aom
RUN \
  echo "**** compiling aom ****" && \
  cd /tmp/aom && \
  rm -rf \
    CMakeCache.txt \
    CMakeFiles && \
  mkdir -p \
    aom_build && \
  cd aom_build && \
  cmake \
    -DBUILD_STATIC_LIBS=0 .. && \
  make && \
  make install
RUN \
  echo "**** grabbing fdk-aac ****" && \
  mkdir -p /tmp/fdk-aac && \
  curl -Lf \
    https://github.com/mstorsjo/fdk-aac/archive/v${FDKAAC}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/fdk-aac
RUN \
  echo "**** compiling fdk-aac ****" && \
  cd /tmp/fdk-aac && \
  autoreconf -fiv && \
  ./configure \
    --disable-static \
    --enable-shared && \
  make && \
  make install
RUN \
  echo "**** grabbing ffnvcodec ****" && \
  mkdir -p /tmp/ffnvcodec && \
  git clone \
    --branch ${NVCODEC} \
    --depth 1 https://git.videolan.org/git/ffmpeg/nv-codec-headers.git \
    /tmp/ffnvcodec
RUN \
  echo "**** compiling ffnvcodec ****" && \
  cd /tmp/ffnvcodec && \
  make install
RUN \
  echo "**** grabbing freetype ****" && \
  mkdir -p /tmp/freetype && \
  curl -Lf \
    https://downloads.sourceforge.net/project/freetype/freetype2/${FREETYPE}/freetype-${FREETYPE}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/freetype
RUN \
  echo "**** compiling freetype ****" && \
  cd /tmp/freetype && \
  ./configure \
    --disable-static \
    --enable-shared && \
  make && \
  make install
RUN \
  echo "**** grabbing fribidi ****" && \
  mkdir -p /tmp/fribidi && \
  curl -Lf \
    https://github.com/fribidi/fribidi/archive/v${FRIBIDI}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/fribidi
RUN \
  echo "**** compiling fribidi ****" && \
  cd /tmp/fribidi && \
  ./autogen.sh && \
  ./configure \
    --disable-static \
    --enable-shared && \
  make -j 1 && \
  make install
RUN \
  echo "**** grabbing kvazaar ****" && \
  mkdir -p /tmp/kvazaar && \
  curl -Lf \
    https://github.com/ultravideo/kvazaar/archive/v${KVAZAAR}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/kvazaar
RUN \
  echo "**** compiling kvazaar ****" && \
  cd /tmp/kvazaar && \
  ./autogen.sh && \
  ./configure \
    --disable-static \
    --enable-shared && \
  make && \
  make install
RUN \
  echo "**** grabbing lame ****" && \
  mkdir -p /tmp/lame && \
  curl -Lf \
    http://downloads.sourceforge.net/project/lame/lame/${LAME}/lame-${LAME}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/lame
RUN \
  echo "**** compiling lame ****" && \
  cd /tmp/lame && \
  cp \
    /usr/share/automake-1.16/config.guess \
    config.guess && \
  cp \
    /usr/share/automake-1.16/config.sub \
    config.sub && \
  ./configure \
    --disable-frontend \
    --disable-static \
    --enable-nasm \
    --enable-shared && \
  make && \
  make install
RUN \
  echo "**** grabbing libass ****" && \
  mkdir -p /tmp/libass && \
  curl -Lf \
    https://github.com/libass/libass/archive/${LIBASS}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/libass
RUN \
  echo "**** compiling libass ****" && \
  cd /tmp/libass && \
  ./autogen.sh && \
  ./configure \
    --disable-static \
    --disable-require-system-font-provider \
    --enable-shared && \
  make && \
  make install
RUN \
  echo "**** grabbing libva ****" && \
  mkdir -p /tmp/libva && \
  curl -Lf \
    https://github.com/intel/libva/archive/${LIBVA}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/libva
RUN \
  echo "**** compiling libva ****" && \
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
  echo "**** grabbing gmmlib ****" && \
  mkdir -p /tmp/gmmlib && \
  curl -Lf \
    https://github.com/intel/gmmlib/archive/refs/tags/intel-gmmlib-${GMMLIB}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/gmmlib
RUN \
  echo "**** compiling gmmlib ****" && \
  mkdir -p /tmp/gmmlib/build && \
  cd /tmp/gmmlib/build && \
  cmake \
    -DCMAKE_BUILD_TYPE=Release \
    .. && \
  make && \
  make install && \
  strip -d /usr/local/lib/libigdgmm.so
RUN \
  echo "**** grabbing IHD ****" && \
  mkdir -p /tmp/ihd && \
  curl -Lf \
    https://github.com/intel/media-driver/archive/refs/tags/intel-media-${IHD}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/ihd
RUN \
  echo "**** compiling IHD ****" && \
  mkdir -p /tmp/ihd/build && \
  cd /tmp/ihd/build && \
  cmake \
    -DLIBVA_DRIVERS_PATH=/usr/lib/x86_64-linux-gnu/dri/ \
    .. && \
  make && \
  make install && \
  strip -d /usr/lib/x86_64-linux-gnu/dri/iHD_drv_video.so
RUN \
  echo "**** grabbing libvpl ****" && \
  mkdir -p /tmp/libvpl && \
  curl -Lf \
    https://github.com/oneapi-src/oneVPL/archive/refs/tags/v${LIBVPL}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/libvpl
RUN \
  echo "**** compiling libvpl ****" && \
  mkdir -p /tmp/libvpl/build && \
  cd /tmp/libvpl/build && \
  cmake .. && \ 
  cmake --build . --config Release && \
  cmake --build . --config Release --target install && \
  strip -d /usr/local/lib/libvpl.so
RUN \
  echo "**** grabbing onevpl ****" && \
  mkdir -p /tmp/onevpl && \
  curl -Lf \
    https://github.com/oneapi-src/oneVPL-intel-gpu/archive/refs/tags/intel-onevpl-${ONEVPL}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/onevpl
COPY /patches/onevpl/*.patch /
RUN \
  echo "**** FFMPEG oneVPL patch ****" && \
  cd /tmp/onevpl && \
  git apply /*.patch
RUN \
  echo "**** compiling onevpl ****" && \
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
  echo "**** grabbing libmfx ****" && \
  mkdir -p /tmp/libmfx && \
  curl -Lf \
    https://github.com/Intel-Media-SDK/MediaSDK/archive/refs/tags/intel-mediasdk-${LIBMFX}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/libmfx
RUN \
  echo "**** compiling libmfx ****" && \
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
  echo "**** grabbing libvdpau ****" && \
  mkdir -p /tmp/libvdpau && \
  git clone \
    --branch ${LIBVDPAU} \
    --depth 1 https://gitlab.freedesktop.org/vdpau/libvdpau.git \
    /tmp/libvdpau
RUN \
  echo "**** compiling libvdpau ****" && \
  cd /tmp/libvdpau && \
  meson setup \
    --prefix=/usr --libdir=/usr/local/lib \
    -Ddocumentation=false \
    build && \
  ninja -C build install && \
  strip -d /usr/local/lib/libvdpau.so
RUN \
  echo "**** grabbing vmaf ****" && \
  mkdir -p /tmp/vmaf && \
  curl -Lf \
    https://github.com/Netflix/vmaf/archive/refs/tags/v${LIBVMAF}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/vmaf
RUN \
  echo "**** compiling libvmaf ****" && \
  cd /tmp/vmaf/libvmaf && \
  meson setup \
    --prefix=/usr --libdir=/usr/local/lib \
    --buildtype release \
    build && \
  ninja -vC build && \
  ninja -vC build install
RUN \
  echo "**** grabbing ogg ****" && \
  mkdir -p /tmp/ogg && \
  curl -Lf \
    http://downloads.xiph.org/releases/ogg/libogg-${OGG}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/ogg
RUN \
  echo "**** compiling ogg ****" && \
  cd /tmp/ogg && \
  ./configure \
    --disable-static \
    --enable-shared && \
  make && \
  make install
RUN \
  echo "**** grabbing opencore-amr ****" && \
  mkdir -p /tmp/opencore-amr && \
  curl -Lf \
    http://downloads.sourceforge.net/project/opencore-amr/opencore-amr/opencore-amr-${OPENCOREAMR}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/opencore-amr
RUN \
  echo "**** compiling opencore-amr ****" && \
  cd /tmp/opencore-amr && \
  ./configure \
    --disable-static \
    --enable-shared  && \
  make && \
  make install
RUN \
  echo "**** grabbing openjpeg ****" && \
  mkdir -p /tmp/openjpeg && \
  curl -Lf \
    https://github.com/uclouvain/openjpeg/archive/v${OPENJPEG}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/openjpeg
RUN \
  echo "**** compiling openjpeg ****" && \
  cd /tmp/openjpeg && \
  rm -Rf \
    thirdparty/libpng/* && \
  curl -Lf \
    https://download.sourceforge.net/libpng/libpng-1.6.37.tar.gz | \
    tar -zx --strip-components=1 -C thirdparty/libpng/ && \
  cmake \
    -DBUILD_STATIC_LIBS=0 \
    -DBUILD_THIRDPARTY:BOOL=ON . && \
  make && \
  make install
RUN \
  echo "**** grabbing opus ****" && \
  mkdir -p /tmp/opus && \
  curl -Lf \
    https://downloads.xiph.org/releases/opus/opus-${OPUS}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/opus
RUN \
  echo "**** compiling opus ****" && \
  cd /tmp/opus && \
  autoreconf -fiv && \
  ./configure \
    --disable-static \
    --enable-shared && \
  make && \
  make install
RUN \
  echo "**** grabbing shaderc ****" && \
  mkdir -p /tmp/shaderc && \
  git clone \
    --branch ${SHADERC} \
    --depth 1 https://github.com/google/shaderc.git \
    /tmp/shaderc
RUN \
  echo "**** compiling shaderc ****" && \
  cd /tmp/shaderc && \
  ./utils/git-sync-deps && \
  mkdir -p build && \
  cd build && \
  cmake -GNinja \
    -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_INSTALL_PREFIX=/usr/local \
    .. && \
  ninja install
RUN \
  echo "**** grabbing SVT-AV1 ****" && \
  mkdir -p /tmp/svt-av1 && \
  curl -Lf \
    https://gitlab.com/AOMediaCodec/SVT-AV1/-/archive/v${SVTAV1}/SVT-AV1-v${SVTAV1}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/svt-av1
RUN \
  echo "**** compiling SVT-AV1 ****" && \
  cd /tmp/svt-av1/Build && \
  cmake .. -G"Unix Makefiles" -DCMAKE_BUILD_TYPE=Release && \
  make && \
  make install
RUN \
  echo "**** grabbing theora ****" && \
  mkdir -p /tmp/theora && \
  curl -Lf \
    http://downloads.xiph.org/releases/theora/libtheora-${THEORA}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/theora
RUN \
  echo "**** compiling theora ****" && \
  cd /tmp/theora && \
  cp \
    /usr/share/automake-1.16/config.guess \
    config.guess && \
  cp \
    /usr/share/automake-1.16/config.sub \
    config.sub && \
  curl -fL \
    'https://gitlab.xiph.org/xiph/theora/-/commit/7288b539c52e99168488dc3a343845c9365617c8.diff' \
    > png.patch && \
  patch ./examples/png2theora.c < png.patch && \
  ./configure \
    --disable-static \
    --enable-shared && \
  make && \
  make install
RUN \
  echo "**** grabbing vid.stab ****" && \
  mkdir -p /tmp/vid.stab && \
  curl -Lf \
    https://github.com/georgmartius/vid.stab/archive/v${LIBVIDSTAB}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/vid.stab
RUN \
  echo "**** compiling vid.stab ****" && \
  cd /tmp/vid.stab && \
  cmake \
    -DBUILD_STATIC_LIBS=0 . && \
  make && \
  make install
RUN \
  echo "**** grabbing vorbis ****" && \
  mkdir -p /tmp/vorbis && \
  curl -Lf \
    http://downloads.xiph.org/releases/vorbis/libvorbis-${VORBIS}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/vorbis
RUN \
  echo "**** compiling vorbis ****" && \
  cd /tmp/vorbis && \
  ./configure \
    --disable-static \
    --enable-shared && \
  make && \
  make install
RUN \
  echo "**** grabbing vpx ****" && \
  mkdir -p /tmp/vpx && \
  curl -Lf \
  https://github.com/webmproject/libvpx/archive/v${VPX}.tar.gz | \
  tar -zx --strip-components=1 -C /tmp/vpx
RUN \
  echo "**** compiling vpx ****" && \
  cd /tmp/vpx && \
  ./configure \
    --disable-debug \
    --disable-docs \
    --disable-examples \
    --disable-install-bins \
    --disable-static \
    --disable-unit-tests \
    --enable-pic \
    --enable-shared \
    --enable-vp8 \
    --enable-vp9 \
    --enable-vp9-highbitdepth && \
  make && \
  make install
RUN \
  echo "**** grabbing vulkan headers ****" && \
  mkdir -p /tmp/vulkan-headers && \
  git clone \
    --branch ${VULKANSDK} \
    --depth 1 https://github.com/KhronosGroup/Vulkan-Headers.git \
    /tmp/vulkan-headers
RUN \
  echo "**** compiling vulkan headers ****" && \
  cd /tmp/vulkan-headers && \
  cmake -S . -B build/ && \
  cmake --install build --prefix /usr/local
RUN \
  echo "**** grabbing webp ****" && \
  mkdir -p /tmp/webp && \
  curl -Lf \
    https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-${WEBP}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/webp
RUN \
  echo "**** compiling webp ****" && \
  cd /tmp/webp && \
  ./configure && \
  make && \
  make install
RUN \
  echo "**** grabbing x264 ****" && \
  mkdir -p /tmp/x264 && \
  curl -Lf \
    https://code.videolan.org/videolan/x264/-/archive/master/x264-stable.tar.bz2 | \
    tar -jx --strip-components=1 -C /tmp/x264
RUN \
  echo "**** compiling x264 ****" && \
  cd /tmp/x264 && \
  ./configure \
    --disable-cli \
    --disable-static \
    --enable-pic \
    --enable-shared && \
  make && \
  make install
RUN \
  echo "**** grabbing x265 ****" && \
  mkdir -p /tmp/x265 && \
  curl -Lf \
    https://bitbucket.org/multicoreware/x265_git/downloads/x265_${X265}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/x265
RUN \
  echo "**** compiling x265 ****" && \
  cd /tmp/x265/build/linux && \
  ./multilib.sh && \
  make -C 8bit install
RUN \
  echo "**** grabbing xvid ****" && \
  mkdir -p /tmp/xvid && \
  curl -Lf \
    https://downloads.xvid.com/downloads/xvidcore-${XVID}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/xvid
RUN \
  echo "**** compiling xvid ****" && \
  cd /tmp/xvid/build/generic && \
  ./configure && \ 
  make && \
  make install

######################################################### IMTL

ENV MTL_REPO=Media-Transport-Library
ENV DPDK_REPO=dpdk
ENV DPDK_VER=23.11
ENV IMTL_USER=imtl

RUN apt-get update -y

# Install dependencies
RUN apt-get install -y git gcc meson python3 python3-pip pkg-config libnuma-dev libjson-c-dev libpcap-dev libgtest-dev libsdl2-dev libsdl2-ttf-dev libssl-dev

RUN pip install pyelftools ninja

RUN apt-get install -y sudo

# some misc tools
RUN apt-get install -y vim htop

RUN apt clean all

# user: imtl
RUN adduser $IMTL_USER
RUN usermod -G sudo $IMTL_USER
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
USER $IMTL_USER

WORKDIR /home/$IMTL_USER/

RUN git config --global user.email "you@example.com" && \
    git config --global user.name "Your Name"

RUN git clone https://github.com/OpenVisualCloud/$MTL_REPO.git && \
    cd $MTL_REPO && \
    git checkout b210f1a85f571507f317d156b105dbe5690a234d && \
    git switch -c b210f1a85f571507f317d156b105dbe5690a234d

RUN \
    sed -i '531 d' ./Media-Transport-Library/include/mtl_api.h && \
    sed -i '531 i char lcores[256];' ./Media-Transport-Library/include/mtl_api.h && \
    sed -i '261 d' ./Media-Transport-Library/app/src/args.c && \
    sed -i '261 i strcpy(p->lcores, list);' ./Media-Transport-Library/app/src/args.c && \
    sed -i '1856 d' ./Media-Transport-Library/app/v4l2_to_ip/v4l2_to_ip.c && \
    sed -i '1856 i strcpy(st_v4l2_tx->param.lcores, tx_lcore);' ./Media-Transport-Library/app/v4l2_to_ip/v4l2_to_ip.c && \
    sed -i '132 d' ./Media-Transport-Library/tests/src/tests.cpp && \
    sed -i '132 i strcpy(p->lcores, optarg);' ./Media-Transport-Library/tests/src/tests.cpp && \
    sed -i '399 d' ./Media-Transport-Library/tests/src/tests.cpp && \
    sed -i '399 i strcpy(p->lcores, ctx->lcores_list);' ./Media-Transport-Library/tests/src/tests.cpp && \
    sed -i '75 i if (getenv("MTL_PARAM_LCORES")) {' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '76 i sprintf(p.lcores, "%s", getenv("MTL_PARAM_LCORES"));' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '77 i } else {' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '78 i   sprintf(p.lcores, "");' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '79 i }' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '80 i if (getenv("MTL_PARAM_DATA_QUOTA")) {' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '81 i   uint32_t data_quota_mbs_per_sch = (uint32_t)atol(getenv("MTL_PARAM_DATA_QUOTA"));' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '82 i   av_log(NULL, AV_LOG_DEBUG, "p.data_quota_mbs_per_sch = %u", data_quota_mbs_per_sch);' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '83 i       p.data_quota_mbs_per_sch = data_quota_mbs_per_sch;' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '84 i } else {' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '85 i   av_log(NULL, AV_LOG_DEBUG, "p.data_quota_mbs_per_sch set do default: %lu", 4 * st20_1080p59_yuv422_10bit_bandwidth_mps());' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '86 i   p.data_quota_mbs_per_sch = (uint32_t)(4 * st20_1080p59_yuv422_10bit_bandwidth_mps());' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '87 i } ' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c

RUN git clone https://github.com/DPDK/$DPDK_REPO.git && \
    cd $DPDK_REPO && \
    git checkout v$DPDK_VER && \
    git switch -c v$DPDK_VER

# build dpdk
RUN cd $DPDK_REPO && \
    git am ../Media-Transport-Library/patches/dpdk/$DPDK_VER/*.patch && \
    meson build && \
    ninja -C build && \
    sudo ninja -C build install && \
    cd ..

# build mtl
RUN cd $MTL_REPO && \
    ./build.sh && \
    cd ..

USER root
WORKDIR /

COPY libraries.media.encoding.svt-jpeg-xs/ /tmp/jpegxs

RUN echo "**** BUILD JPEG-XS ****" && \
    mkdir /tmp/jpegxs/Build/linux/install && \
    cd /tmp/jpegxs/Build/linux && \
    ./build.sh install --prefix=/tmp/jpegxs/Build/linux/install

RUN ls -lR /tmp/jpegxs/Build/linux/install

RUN echo "**** BUILD JPEG-XS MTL PLUGIN ****" && \
    cd /tmp/jpegxs/imtl-plugin && \
    ./build.sh --prefix=/tmp/jpegxs/Build/linux/install

ENV LD_LIBRARY_PATH="/tmp/jpegxs/Build/linux/install/lib:${LD_LIBRARY_PATH}"
ENV PKG_CONFIG_PATH="/tmp/jpegxs/Build/linux/install/lib/pkgconfig:${PKG_CONFIG_PATH}"

# add ffmpeg patches
ENV \
  CARTWHEEL_COMMIT_ID=6.1 \
  FFMPEG_COMMIT_ID=n6.1.1

RUN git config --global user.email "you@example.com" && \
    git config --global user.name "Your Name"

RUN \
  echo "**** Get FFMPEG from git branch ****" && \
  mkdir -p /tmp/ffmpeg && \
  curl -Lf \
    https://github.com/ffmpeg/ffmpeg/archive/${FFMPEG_COMMIT_ID}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/ffmpeg
RUN \
  echo "**** Cartwheel patch ****" && \
  mkdir -p /tmp/cartwheel && \
  curl -Lf \
    https://github.com/intel/cartwheel-ffmpeg/archive/${CARTWHEEL_COMMIT_ID}.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/cartwheel && \
  cd /tmp/ffmpeg && \
  git init && \
  git add . && \
  git commit -m "initial" && \
  git am /tmp/cartwheel/patches/*.patch
RUN \
  echo "**** openh264 library ****" && \
  mkdir -p /tmp/openh264 && \
  curl -Lf \
    https://github.com/cisco/openh264/archive/master.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/openh264 && \
  cd /tmp/openh264 && \
  make && \
  sudo make install
COPY \
  patches/ffmpeg/*.diff /
# RUN \
#   echo "**** FFMPEG patches ****" && \
#   cd /tmp/ffmpeg && \
#   git apply /*.diff

RUN rm -rf /tmp/jpegxs/ffmpeg-plugin/*.patch
COPY patches/jpegxs /tmp/jpegxs/ffmpeg-plugin/


RUN cp /home/$IMTL_USER/$MTL_REPO/ecosystem/ffmpeg_plugin/mtl_*.c -rf /tmp/ffmpeg/libavdevice/ && \
    cp /home/$IMTL_USER/$MTL_REPO/ecosystem/ffmpeg_plugin/mtl_*.h -rf /tmp/ffmpeg/libavdevice/

RUN echo "**** APPLY MTL PATCHES ****" && \
    cd /tmp/ffmpeg/ && \
    git am /home/$IMTL_USER/$MTL_REPO/ecosystem/ffmpeg_plugin/6.1/*.patch

RUN echo "**** APPLY JPEG-XS PATCHES ****" && \
    cd /tmp/ffmpeg/ && \
    git apply -v --whitespace=fix /tmp/jpegxs/ffmpeg-plugin/*.patch

RUN \
  echo "**** FFMPEG 1 patches ****" && \
  cd /tmp/ffmpeg && \
  git apply -v ../../hwupload_async.diff

RUN \
  echo "**** FFMPEG 2 patches ****" && \
  cd /tmp/ffmpeg && \
  git apply -v ../../qsv_aligned_malloc.diff
RUN \
  echo "**** FFMPEG 3 patches ****" && \
  cd /tmp/ffmpeg && \
  git apply -v ../../qsvvpp_async.diff
RUN \
  echo "**** FFMPEG 4 patches ****" && \
  cd /tmp/ffmpeg && \
  git apply -v ../../filtergraph_async.diff

RUN sudo apt-get install wget

RUN wget https://registrationcenter-download.intel.com/akdlm/IRC_NAS/046b1402-c5b8-4753-9500-33ffb665123f/l_ipp_oneapi_p_2021.10.1.16_offline.sh
RUN chmod +x l_ipp_oneapi_p_2021.10.1.16_offline.sh
RUN ./l_ipp_oneapi_p_2021.10.1.16_offline.sh -a -s --eula accept
RUN echo "source /opt/intel/oneapi/ipp/latest/env/vars.sh" | tee -a ~/.bash_profile
ENV PKG_CONFIG_PATH="/usr/local/lib/pkgconfig:$PKG_CONFIG_PATH"

RUN echo "**** DOWNLOADING VIDEO SUPER RESOLUTION (VSR) ****" && \
    git clone https://github.com/OpenVisualCloud/Video-Super-Resolution-Library.git /tmp/vsr && \
    cd /tmp/vsr/ && \
    git checkout v23.11

RUN cd /tmp/vsr/ && \
    sed -i '24 d' ./ffmpeg/vf_raisr_opencl.c && \
    sed -i '24 i #include "video.h"' ./ffmpeg/vf_raisr_opencl.c && \
    sed -i '25 d' ./ffmpeg/vf_raisr_opencl.c

RUN echo "**** BUILD VIDEO SUPER RESOLUTION (VSR) ****" && \
    cd /tmp/vsr/ && . /opt/intel/oneapi/ipp/latest/env/vars.sh && \
    ./build.sh -DCMAKE_INSTALL_PREFIX="$PWD/install" -DENABLE_RAISR_OPENCL=ON

COPY patches/vsr/*.patch /vsr_patches/

RUN echo "**** APPLY VSR PATCHES ****" && \
    cd tmp/ffmpeg/ && \
    git apply -v /vsr_patches/*.patch

RUN cp /tmp/vsr/ffmpeg/vf_raisr*.c /tmp/ffmpeg/libavfilter

RUN \
  echo "**** compiling ffmpeg ****" && . /opt/intel/oneapi/ipp/latest/env/vars.sh && \
  cd /tmp/ffmpeg && \
    ./configure \
    --disable-debug \
    --disable-doc \
    --enable-static \
    --enable-cuvid \
    --enable-encoder=libopenh264 \
    --enable-ffprobe \
    --enable-gpl \
    --enable-libaom \
    --enable-libass \
    --enable-libfdk_aac \
    --enable-libfreetype \
    --enable-libkvazaar \
    --enable-libmp3lame \
    --enable-libopencore-amrnb \
    --enable-libopencore-amrwb \
    --enable-libopenjpeg \
    --enable-libopus \
    --enable-libshaderc \
    --enable-libsvtav1 \
    --enable-libtheora \
    --enable-libv4l2 \
    --enable-libvidstab \
    --enable-libvmaf \
    --enable-libvorbis \
    --enable-libvpl \
    --enable-libvpx \
    --enable-libwebp \
    --enable-libx264 \
    --enable-libx265 \
    --enable-libxml2 \
    --enable-libxvid \
    --enable-mtl \
    --enable-nonfree \
    --enable-nvdec \
    --enable-nvenc \
    --enable-opencl \
    --enable-openssl \
    --enable-pic \
    --enable-shared \
    --enable-stripping \
    --enable-vaapi \
    --enable-vdpau \
    --enable-version3 \
    --enable-vulkan \
    --enable-libsvtjpegxs \
    --enable-libipp \
    --extra-cflags="-fopenmp -I/tmp/vsr/install/include/ -I/opt/intel/oneapi/ipp/latest/include/ipp/" \
    --extra-ldflags="-fopenmp -L/tmp/vsr/install/lib" \
    --extra-libs='-lraisr -lstdc++ -lippcore -lippvm -lipps -lippi -lm' \
    --enable-cross-compile && \
  make

RUN \
  echo "**** arrange files ****" && \
  sudo ldconfig && \
  sudo mkdir -p \
    /buildout/usr/local/bin \
    /buildout/usr/local/lib/libmfx-gen \
    /buildout/usr/local/lib/mfx \
    /buildout/usr/local/lib/vpl \
    /buildout/usr/local/lib/x86_64-linux-gnu/dri \
    /buildout/etc/OpenCL/vendors && \
  sudo cp \
    /tmp/ffmpeg/ffmpeg \
    /buildout/usr/local/bin && \
  sudo cp \
    /tmp/ffmpeg/ffprobe \
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
    /usr/local/lib/vpl/*.so \
    /buildout/usr/local/lib/vpl/ && \
  sudo cp -a \
    /usr/local/lib/x86_64-linux-gnu/lib*so* \
    /buildout/usr/local/lib/x86_64-linux-gnu/ && \
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
  sudo echo \
    'libnvidia-opencl.so.1' | sudo tee \
    /buildout/etc/OpenCL/vendors/nvidia.icd

# runtime stage
FROM ubuntu as finalstage

# Add files from binstage
COPY --from=buildstage /buildout/ /

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
  echo "**** install runtime ****" && \
    apt-get update && \
    apt-get install -y \
    libexpat1 \
    libglib2.0-0 \
    libgomp1 \
    libharfbuzz0b \
    libpciaccess0 \
    libv4l-0 \
    libwayland-client0 \
    libx11-6 \
    libx11-xcb1 \
    libxcb-dri3-0 \
    libxcb-shape0 \
    libxcb-xfixes0 \
    libxcb1 \
    libxext6 \
    libxfixes3 \
    libxml2 \
    ocl-icd-libopencl1 && \
  echo "**** clean up ****" && \
  rm -rf \
    /var/lib/apt/lists/* \
    /var/tmp/*
    
############################# IMTL
ENV MTL_REPO=Media-Transport-Library
ENV DPDK_REPO=dpdk
ENV DPDK_VER=23.11
ENV IMTL_USER=imtl

RUN apt-get update -y

# Install dependencies
RUN apt-get install -y git gcc meson python3 python3-pip pkg-config libnuma-dev libjson-c-dev libpcap-dev libgtest-dev libsdl2-dev libsdl2-ttf-dev libssl-dev sudo autoconf curl tar libtool

RUN pip install pyelftools ninja

RUN \
  echo "**** grabbing libva ****" && \
  mkdir -p /tmp/libva && \
  ls -ltr /tmp/libva/ && \
  curl -Lf \
    https://github.com/intel/libva/archive/2.20.0.tar.gz | \
    tar -zx --strip-components=1 -C /tmp/libva
RUN \
  echo "**** compiling libva ****" && \
  cd /tmp/libva && \
  ./autogen.sh && \
  ./configure \
    --disable-static \
    --enable-shared && \
  make && \
  make install

# user: imtl
RUN adduser $IMTL_USER
RUN usermod -G sudo $IMTL_USER
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
USER $IMTL_USER

WORKDIR /home/$IMTL_USER/

RUN git config --global user.email "you@example.com" && \
    git config --global user.name "Your Name"
    
RUN git clone https://github.com/OpenVisualCloud/$MTL_REPO.git && \
    cd $MTL_REPO && \
    git checkout b210f1a85f571507f317d156b105dbe5690a234d && \
    git switch -c b210f1a85f571507f317d156b105dbe5690a234d

RUN \
    sed -i '531 d' ./Media-Transport-Library/include/mtl_api.h && \
    sed -i '531 i char lcores[256];' ./Media-Transport-Library/include/mtl_api.h && \
    sed -i '261 d' ./Media-Transport-Library/app/src/args.c && \
    sed -i '261 i strcpy(p->lcores, list);' ./Media-Transport-Library/app/src/args.c && \
    sed -i '1856 d' ./Media-Transport-Library/app/v4l2_to_ip/v4l2_to_ip.c && \
    sed -i '1856 i strcpy(st_v4l2_tx->param.lcores, tx_lcore);' ./Media-Transport-Library/app/v4l2_to_ip/v4l2_to_ip.c && \
    sed -i '132 d' ./Media-Transport-Library/tests/src/tests.cpp && \
    sed -i '132 i strcpy(p->lcores, optarg);' ./Media-Transport-Library/tests/src/tests.cpp && \
    sed -i '399 d' ./Media-Transport-Library/tests/src/tests.cpp && \
    sed -i '399 i strcpy(p->lcores, ctx->lcores_list);' ./Media-Transport-Library/tests/src/tests.cpp && \
    sed -i '75 i if (getenv("MTL_PARAM_LCORES")) {' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '76 i sprintf(p.lcores, "%s", getenv("MTL_PARAM_LCORES"));' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '77 i } else {' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '78 i   sprintf(p.lcores, "");' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '79 i }' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '80 i if (getenv("MTL_PARAM_DATA_QUOTA")) {' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '81 i   uint32_t data_quota_mbs_per_sch = (uint32_t)atol(getenv("MTL_PARAM_DATA_QUOTA"));' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '82 i   av_log(NULL, AV_LOG_DEBUG, "p.data_quota_mbs_per_sch = %u", data_quota_mbs_per_sch);' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '83 i       p.data_quota_mbs_per_sch = data_quota_mbs_per_sch;' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '84 i } else {' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '85 i   av_log(NULL, AV_LOG_DEBUG, "p.data_quota_mbs_per_sch set do default: %lu", 4 * st20_1080p59_yuv422_10bit_bandwidth_mps());' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '86 i   p.data_quota_mbs_per_sch = (uint32_t)(4 * st20_1080p59_yuv422_10bit_bandwidth_mps());' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c && \ 
    sed -i '87 i } ' ./Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_common.c

RUN git clone https://github.com/DPDK/$DPDK_REPO.git && \
    cd $DPDK_REPO && \
    git checkout v$DPDK_VER && \
    git switch -c v$DPDK_VER

# build dpdk
RUN cd $DPDK_REPO && \
    git am ../Media-Transport-Library/patches/dpdk/$DPDK_VER/*.patch && \
    meson build && \
    ninja -C build && \
    sudo ninja -C build install && \
    cd ..

# build mtl
RUN cd $MTL_REPO && \
    ./build.sh && \
    cd ..

USER root
WORKDIR /

RUN apt-get update -y
RUN apt-get install -y wget cmake make pkg-config intel-opencl-icd opencl-headers ocl-icd-opencl-dev nasm

RUN mkdir /hugetlbfs

RUN \
  apt-get install -y kmod libkmod-dev pciutils gdb

ENTRYPOINT ["./usr/local/bin/ffmpeg"]