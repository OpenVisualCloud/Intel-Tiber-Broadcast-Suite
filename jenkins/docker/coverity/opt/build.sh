#! /bin/bash -e 

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
    -DCMAKE_INSTALL_LIBDIR=/usr/local/lib \
    .. 
make
make install



cd /tmp/vsr/ && . /opt/intel/oneapi/ipp/latest/env/vars.sh 
./build.sh -DCMAKE_INSTALL_PREFIX="$(pwd)/install" -DENABLE_RAISR_OPENCL=ON
