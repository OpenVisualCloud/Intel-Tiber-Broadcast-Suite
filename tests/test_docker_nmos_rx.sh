#!/bin/bash

docker run -it \
   --user root\
   --privileged \
   --device=/dev/vfio:/dev/vfio \
   --device=/dev/dri:/dev/dri \
   --cap-add ALL \
   -v "$(pwd)":/home/config/ \
   -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
   -v /tmp/kahawai_lcore.lock:/tmp/kahawai_lcore.lock \
   -v /dev/null:/dev/null \
   -v /tmp/hugepages:/tmp/hugepages \
   -v /hugepages:/hugepages \
   -v /var/run/imtl:/var/run/imtl \
   -e http_proxy="" \
   -e https_proxy="" \
   -e VFIO_PORT_RX=0000:31:01.2 \
   --network=my_net_801f0 \
   --ip=192.168.2.3 \
   --expose=20000-20170 \
   --ipc=host -v /dev/shm:/dev/shm \
      tiber-broadcast-suite-nmos-node config/intel-node-rx.json
