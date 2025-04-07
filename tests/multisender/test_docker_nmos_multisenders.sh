#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <video_type>"
  exit 1
fi

case $1 in
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

config_file="intel-node-multisenders${format}.json"

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
   -e VFIO_PORT_TX=0000:27:01.0 \
   --network=host \
   --ipc=host \
   -v /dev/shm:/dev/shm \
      tiber-broadcast-suite-nmos-node config/$config_file
