#! /usr/bin/bash -e

echo "/************************************************************************************/"
echo "e810 driver firmware update steps."
echo "for more details see:"
echo "https://github.com/OpenVisualCloud/Media-Transport-Library/blob/main/doc/e810.md#2-update-firmware-version-to-latest"
echo "/************************************************************************************/"




RELEASE_ZIP=${1:-Release_29.0.zip}
ETH_NAME=${2:-eth0}

if [ ! -f  "${RELEASE_ZIP}" ]
then
    echo "/********************************* no ${ICE_DRIVER} file found, exiting. *********************************/"
    exit 1
fi

echo "/********************************* extracting ${RELEASE_ZIP} ... *********************************/"
mkdir -p release-zip
unzip ${RELEASE_ZIP} -d release-zip
cd release-zip/NVMUpdatePackage/E810
tar xvf E810_NVMUpdatePackage_v4_40_Linux.tar.gz
cd E810/Linux_x64/

echo "/********************************* updating driver ... *********************************/"
sudo ./nvmupdate64e



echo "/********************************* print updated version *********************************/"
ethtool -i ${ETH_NAME}

