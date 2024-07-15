#!/bin/bash

echo "**************** BUILD VSR ****************"
cd /tmp/vsr/
rm -rf  /tmp/vsr/build
sed -i s'/make\ install\ -j/make\ install\ -B -j/'g  ./build.sh
./build.sh -DCMAKE_INSTALL_PREFIX="/tmp/vsr/install" -DENABLE_RAISR_OPENCL=ON
