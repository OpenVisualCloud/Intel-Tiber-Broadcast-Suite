#!/bin/bash
. /opt/intel/oneapi/ipp/latest/env/vars.sh
echo "**************** BUILD FFMPEG ****************"
cd /tmp/ffmpeg
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
        --extra-cflags="-march=native -fopenmp -I/tmp/vsr/install/include/ -I/opt/intel/oneapi/ipp/latest/include/ipp/" \
        --extra-ldflags="-fopenmp -L/tmp/vsr/install/lib" \
        --extra-libs='-lraisr -lstdc++ -lippcore -lippvm -lipps -lippi -lpthread -lm -lz -lbsd -lrdmacm' \
        --enable-cross-compile
make -B
