#!/bin/bash

#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <number>"
  exit 1
fi

case $1 in
  0)
    docker run -it \
      --user root \
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
      -e VFIO_PORT_TX=0000:31:01.0 \
      --network=host \
      --ipc=host \
      -v /dev/shm:/dev/shm \
      tiber-broadcast-suite-nmos-node config/intel-node-tx.json
    ;;
  1)
    docker run -it \
      --user root \
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
      -e VFIO_PORT_TX=0000:31:01.2 \
      --network=host \
      --ipc=host \
      -v /dev/shm:/dev/shm \
      tiber-broadcast-suite-nmos-node config/intel-node-tx-1.json
    ;;
  2)
    docker run -it \
      --user root \
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
      -e VFIO_PORT_TX=0000:31:01.3 \
      --network=host \
      --ipc=host \
      -v /dev/shm:/dev/shm \
      tiber-broadcast-suite-nmos-node config/intel-node-tx-2.json
    ;;
  3)
    docker run -it \
      --user root \
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
      -e VFIO_PORT_TX=0000:31:01.3 \
      --network=host \
      --ipc=host \
      -v /dev/shm:/dev/shm \
      tiber-broadcast-suite-nmos-node config/intel-node-tx-h265.json
    ;;
  *)
    echo "Invalid number: $1. Please provide 0, 1 or 2."
    exit 1
    ;;
esac