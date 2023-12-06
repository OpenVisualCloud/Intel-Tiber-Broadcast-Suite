#! /bin/bash
FULL_ARGS=( "$@" )
mkdir /tmp/hugepages
mount -t hugetlbfs hugetlbfs /tmp/hugepages -o pagesize=2M
mount -t hugetlbfs hugetlbfs /hugepages -o pagesize=1G
#./tmp/MTL/script/nicctl.sh create_vf ${NIC_PORT}
modprobe ice
modprobe vfio-pci
#gdb --args /usr/local/bin/ffmpeg_g "${FULL_ARGS[@]}"
/usr/local/bin/ffmpeg_g "${FULL_ARGS[@]}"