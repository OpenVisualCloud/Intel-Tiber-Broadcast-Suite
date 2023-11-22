# Video Production Pipeline Project

To run ffmpeg in Docker please run command which creates docker image:

```
docker build -t my_ffmpeg .
```

Step 1. If IMTL plugin support is needed then please run commands on host as a root (echo for each numa node):

```
mount -t hugetlbfs hugetlbfs /hugepages -o pagesize=1G
echo 1 >  /sys/devices/system/node/node0/hugepages/hugepages-1048576kB/nr_hugepages
echo 1 >  /sys/devices/system/node/node1/hugepages/hugepages-1048576kB/nr_hugepages
(...)
echo 1 >  /sys/devices/system/node/node{n}/hugepages/hugepages-1048576kB/nr_hugepages

docker run -it \
  --privileged \
  --device=/dev/dri:/dev/dri \
  --cap-add=ALL \
  -v /dev:/dev \
  -v /tmp/kahawai_lcore.lock:/tmp/kahawai_lcore.lock \
  -v /hugepages:/hugepages \
  -v /lib/modules/$(uname -r):/lib/modules/$(uname -r) \
  --net=host \
  --entrypoint /bin/bash \
  my_ffmpeg
```

Step 2. In my_ffmpeg container please run (NIC port is in format 0000:xx:00.0):

```
mkdir /tmp/hugepages
mount -t hugetlbfs hugetlbfs /tmp/hugepages -o pagesize=2M
mount -t hugetlbfs hugetlbfs /hugepages -o pagesize=1G
./tmp/MTL/script/nicctl.sh create_vf <paste NIC port here>
exit
```

Steps 1 and 2 are needed each time host is restarted and IMTL is needed.

Step 3. Run .sh script with ffmpeg parameters. Examples are in [test_scripts](./test_scripts) directory.


