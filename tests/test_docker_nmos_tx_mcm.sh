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
   -v /run/mcm:/run/mcm \
   -e http_proxy="" \
   -e https_proxy="" \
   -e HTTP_PROXY="" \
   -e HTTPS_PROXY="" \
   -e MCM_MEDIA_PROXY_PORT=8002 \
   -e VFIO_PORT_TX=0000:31:01.0 \
   --network=host \
   --ipc=host -v /dev/shm:/dev/shm \
      tiber-broadcast-suite-nmos-node config/intel-node-tx-mcm.json