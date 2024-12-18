#!/bin/bash


# script acording to readme
# init submodule
echo "**** BUILD gRPC ****"
cd gRPC

# add rebuild flag to make command
sed -i '$s/make/make -B/' compile.sh
./compile.sh