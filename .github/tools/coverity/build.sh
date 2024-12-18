#!/bin/bash

BASE_DIR=$(pwd)
# script acording to readme
# init submodule
echo "**** BUILD gRPC ****"
cd gRPC

# add rebuild flag to make command
sed -i '$s/make/make -B/' compile.sh
./compile.sh


cd $BASE_DIR

echo "**** BUILD pod Launcher ****"
cd launcher/cmd/
go build main.go

