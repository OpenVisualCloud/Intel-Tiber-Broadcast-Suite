#! /bin/bash -e 

MTL_VER=2f1c2a3be417065a4dc9276e2d7344d768e95118
ONEVPL=23.3.4
VSR=v23.11
JPEG_XS_VER=0.9.0

cp /patches /tmp/patches
echo  " downloading repositories ..." 

echo "**** DOWNLOAD MTL ****"
curl -Lf \
  https://github.com/OpenVisualCloud/Media-Transport-Library/archive/${MTL_VER}.tar.gz | \
  tar -zx --strip-components=1 -C /tmp/Media-Transport-Library

echo "**** DOWNLOAD and PATCH ONEVPL ****" && \
curl -Lf \
  https://github.com/oneapi-src/oneVPL-intel-gpu/archive/refs/tags/intel-onevpl-${ONEVPL}.tar.gz | \
  tar -zx --strip-components=1 -C /tmp/onevpl
git -C /tmp/onevpl apply /tmp/patches/onevpl/*.patch

echo "**** DOWNLOAD JPEG-XS ****" 
curl -Lf \
  https://github.com/OpenVisualCloud/SVT-JPEG-XS/archive/refs/tags/v${JPEG_XS_VER}.tar.gz | \
  tar -zx --strip-components=1 -C /tmp/jpegxs 
chmod +x ./build.sh /tmp/jpegxs/imtl-plugin/build.sh
git -C /tmp/ffmpeg apply --whitespace=fix /tmp/jpegxs/ffmpeg-plugin/*.patch

echo "**** DOWNLOAD VIDEO SUPER RESOLUTION ****"
curl -Lf \
  https://github.com/OpenVisualCloud/Video-Super-Resolution-Library/archive/refs/tags/${VSR}.tar.gz | \
  tar -zx --strip-components=1 -C /tmp/vsr

echo "**** APPLY VIDEO SUPER RESOLUTION PATCHES ****" 
git -C /tmp/ffmpeg apply /tmp/patches/vsr/0001-ffmpeg-raisr-filter.patch 
git -C /tmp/ffmpeg apply /tmp/patches/vsr/0002-libavfilter-raisr_opencl-Add-raisr_opencl-filter.patch 
cp /tmp/vsr/ffmpeg/vf_raisr*.c /tmp/ffmpeg/libavfilter

echo "**** PATCH MTL ****" 
cp /tmp/Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_*.c -rf /tmp/ffmpeg/libavdevice/
cp /tmp/Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_*.h -rf /tmp/ffmpeg/libavdevice/
git -C /tmp/ffmpeg/ apply /tmp/Media-Transport-Library/ecosystem/ffmpeg_plugin/6.1/*.patch

echo "**** APPLY JPEG-XS PATCHES ****"
git -C /tmp/ffmpeg apply --whitespace=fix /tmp/jpegxs/ffmpeg-plugin/*.patch

echo "**** APPLY FFMPEG patches ****" 
git -C /tmp/ffmpeg apply /tmp/patches/ffmpeg/hwupload_async.diff
git -C /tmp/ffmpeg apply /tmp/patches/ffmpeg/qsv_aligned_malloc.diff 
git -C /tmp/ffmpeg apply /tmp/patches/ffmpeg/qsvvpp_async.diff 
git -C /tmp/ffmpeg apply /tmp/patches/ffmpeg/filtergraph_async.diff

export LD_LIBRARY_PATH="/tmp/jpegxs/Build/linux/install/lib"
export PKG_CONFIG_PATH="/tmp/jpegxs/Build/linux/install/lib/pkgconfig"

echo "performing coverity scan ..."
echo "**** BUILD FFMPEG ****" 
. /opt/intel/oneapi/ipp/latest/env/vars.sh 
cd /tmp/ffmpeg 
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
    --extra-cflags="-fopenmp -I/tmp/vsr/install/include/ -I/opt/intel/oneapi/ipp/latest/include/ipp/" \
    --extra-ldflags="-fopenmp -L/tmp/vsr/install/lib" \
    --extra-libs='-lraisr -lstdc++ -lippcore -lippvm -lipps -lippi -lm' \
    --enable-cross-compile 
make

cd /tmp/Media-Transport-Library 
./build.sh

cd /tmp/onevpl/build 
cmake \
    -DCMAKE_INSTALL_PREFIX=/usr \
    -DCMAKE_INSTALL_LIBDIR=/usr/lib/x86_64-linux-gnu \
    ..
make
make install
strip -d /usr/lib/x86_64-linux-gnu/libmfx-gen.so

cd /tmp/vsr/ && . /opt/intel/oneapi/ipp/latest/env/vars.sh 
./build.sh -DCMAKE_INSTALL_PREFIX="$(pwd)/install" -DENABLE_RAISR_OPENCL=ON
