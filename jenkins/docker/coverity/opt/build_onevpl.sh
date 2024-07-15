#!/bin/bash

. /opt/intel/oneapi/ipp/latest/env/vars.sh
echo "**************** BUILD ONEVPL ****************"
cd /tmp/onevpl/build
cmake \
    -DCMAKE_INSTALL_PREFIX=/usr \
    -DCMAKE_INSTALL_LIBDIR=/usr/lib/x86_64-linux-gnu .. 
make -B
make install 
strip -d /usr/lib/x86_64-linux-gnu/libmfx-gen.so
