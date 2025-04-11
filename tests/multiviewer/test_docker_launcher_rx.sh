#!/bin/bash

if [ $# -ne 2 ]; then
  echo "Usage: $0 <hostname/ip> <port> e.g. $0 192.168.2.1 50056"
  exit 1
fi

host=$1
PORT=$2
  
if [[ "$host" == "localhost" ]]; then
  ip=""
  network=host
  HOSTNAME=localhost
else
  ip=192.168.2.6
  network=my_net_801f0
  HOSTNAME=192.168.2.6
fi

echo "Starting container $i"
echo "Host name: $HOSTNAME"
echo "Port: $PORT"
echo "IP: $ip"

docker run -it \
   --user root \
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
   -v /var/run/imtl:/var/run/imtl \
   -e http_proxy="" \
   -e https_proxy="" \
   --network=$network \
   --ip=$ip \
   --ipc=host \
   -v /dev/shm:/dev/shm \
   tiber-broadcast-suite "$HOSTNAME" "$PORT"