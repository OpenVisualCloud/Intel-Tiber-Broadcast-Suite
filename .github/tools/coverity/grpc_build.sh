#!/bin/bash


# script acording to readme
# init submodule
echo "**** BUILD gRPC ****"
cd ${1}/gRPC

# add rebuild flag to make command
sed -i '$s/make/make -B/' compile.sh
./compile.sh
cd ${1}