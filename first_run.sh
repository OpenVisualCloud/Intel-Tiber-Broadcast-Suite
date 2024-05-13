#!/bin/bash
touch /tmp/kahawai_lcore.lock
sysctl -w vm.nr_hugepages=4096
mkdir /tmp/hugepages
mkdir /hugepages
mount -t hugetlbfs hugetlbfs /tmp/hugepages -o pagesize=2M
mount -t hugetlbfs hugetlbfs /hugepages -o pagesize=1G
for f in /sys/devices/system/node/node*; do echo 1 > "$f/hugepages/hugepages-1048576kB/nr_hugepages"; done
docker network create --subnet 192.168.2.0/24 --gateway 192.168.2.100 -o parent=ens801f0 my_net_801f0
output=$(lspci | grep "Ethernet controller: Intel Corporation Ethernet Controller" | awk '{print "0000:"$1}')
IFS=$'\n'
while IFS= read -r line; do find ${HOME} -name "nicctl.sh" -exec {} create_vf $line \;; done <<< $output
chmod 777 -R /dev/vfio