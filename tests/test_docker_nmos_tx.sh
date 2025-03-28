#!/bin/bash

#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <number>"
  exit 1
fi

if [ -z "$2" ]; then
  echo "Usage: $0 $1 <video_type>"
  exit 1
fi

case $2 in
  "raw")
    format=""
    ;;
  "h264")
    format="-h264"
    ;;
  "h265")
    format="-h265"
    ;;
  *)
    echo "Invalid option: $2"
    echo "Valid options are: raw, h264, h265"
    exit 1
    ;;
esac

config_file="intel-node-tx${format}.json"
config_file1="intel-node-tx-1${format}.json"
config_file2="intel-node-tx-2${format}.json"

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
      tiber-broadcast-suite-nmos-node config/$config_file
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
      tiber-broadcast-suite-nmos-node config/$config_file1
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
      tiber-broadcast-suite-nmos-node config/$config_file2
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