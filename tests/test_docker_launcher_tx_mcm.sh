#!/bin/bash

docker run -it \
   --user root\
   --privileged \
   --device=/dev/vfio:/dev/vfio \
   --device=/dev/dri:/dev/dri \
   --cap-add ALL \
   -v "$(pwd)":/videos \
   -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
   -v /tmp/kahawai_lcore.lock:/tmp/kahawai_lcore.lock \
   -v /dev/null:/dev/null \
   -v /tmp/hugepages:/tmp/hugepages \
   -v /hugepages:/hugepages \
   -v /run/mcm:/var/run/mcm \
   -v /var/run/imtl:/var/run/imtl \
   -e http_proxy="" \
   -e https_proxy="" \
   -e HTTP_PROXY="" \
   -e HTTPS_PROXY="" \
   --expose=5000-20170 \
   --ipc=host -v /dev/shm:/dev/shm \
   --net=host \
   -e MCM_MEDIA_PROXY_PORT=8002 \
      tiber-broadcast-suite localhost 50057