#!/bin/bash
echo "Usage $0 <video_type> <localhost/ip>"

if [ -z "$1" ]; then
  echo "Usage: $0 <video_type>"
  exit 1
fi

if [ -z "$2" ]; then
  echo "Usage: $0 $1 <localhost/ip>"
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

host=$2
config_file="intel-node-rx${format}.json"
port_rx=0000:27:01.1

if [[ "$host" == "localhost" ]]; then
  ip=""
  network=host
else
  ip=192.168.2.7
  network=my_net_801f0
fi

echo "Starting container"
echo "Configuration file: $config_file"
echo "VFIO address: $port_rx"
echo "IP: $ip"

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
   -e VFIO_PORT_RX=$port_rx \
   --network=$network \
   --ip=$ip \
   --ipc=host \
   -v /dev/shm:/dev/shm \
      tiber-broadcast-suite-nmos-node config/$config_file
