#! /usr/bin/bash -e

echo "/************************************************************************************/"
echo "e810 driver installation steps."
echo "for more details see:"
echo "    https://github.com/OpenVisualCloud/Media-Transport-Library/blob/main/doc/e810.md"
echo "/************************************************************************************/"
SCRIPT_DIR=$(pwd)
function cleanup(){
    cd ${SCRIPT_DIR}
    rm -rf mtl
    rm -rf ${ICE_DRIVER}
}


ICE_DRIVER="${1:-ice-1.13.7}"

if [ ! -f  "${ICE_DRIVER}.tar.gz" ]
then
    echo "/********************************* no ${ICE_DRIVER}.tar.gz file found, exiting *********************************/"
    exit 1
fi

git clone https://github.com/OpenVisualCloud/Media-Transport-Library.git mtl

echo "/********************************* 1.2 Unzip ${ICE_DRIVER/ice-/} driver and enter into the source code directory *********************************/"
tar xvzf ${ICE_DRIVER}.tar.gz

echo "/********************************* 1.3 Patch 1.13.7 driver with rate limit patches *********************************/"
cd ${ICE_DRIVER}
git init
git add .
git config --local user.email "gta@intel.com"
git config --global user.name "gta"
git commit -m "init version ${ICE_DRIVER/ice-/}"
git am ../mtl/patches/ice_drv/${ICE_DRIVER/ice-/}/*.patch
COMMIT_MESSAGE=$(git log | grep 'update to Kahawai_1.13.7_20240220')

if [ -z "${COMMIT_MESSAGE}" ]
then
    echo "update failed exiting"
    cleanup
    exit 1
fi

echo "/********************************* 1.4 Build and install the driver *********************************/"
cd src
make
sudo make install
sudo rmmod ice
sudo modprobe ice

echo "/********************************* All done, cleaning up... *********************************/"
cleanup