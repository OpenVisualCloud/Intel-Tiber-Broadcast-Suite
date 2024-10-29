#!/bin/bash

rm -rf /tmp/jpegxs/Build/linux/install /tmp/jpegxs/Build/linux/Release /tmp/jpegxs/imtl-plugin/build

echo "**************** BUILD JPEG XS ****************"
cd /tmp/jpegxs/Build/linux
./build.sh install --prefix=/tmp/jpegxs/Build/linux/install

echo "**************** BUILD JPEG XS PLUGIN ****************"
cd /tmp/jpegxs/imtl-plugin
./build.sh
