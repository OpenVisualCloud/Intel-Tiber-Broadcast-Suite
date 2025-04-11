#!/bin/bash

echo "Usage $0 <number_of_containers> <video_type> <localhost/ip>"

number_of_containers=$1

if [ -z "$1" ]; then
  echo "Usage: $0 <number_of_containers>"
  exit 1
fi

if [ -z "$2" ]; then
  echo "Usage: $0 $1 <video_type>"
  exit 1
fi

if [ -z "$3" ]; then
  echo "Usage: $0 $1 $2 <localhost/ip>"
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

host=$3
config_file="intel-node-tx${format}.json"

for ((i=1; i<=number_of_containers; i++)); do 
  port_tx=$(awk "NR==$((i+2))" vfio_addresses.txt)

  if [[ "$host" == "localhost" ]]; then
    ip=""
    network=host
  else
    ip=$(awk "NR==${i}" IP_node.txt)
    network=my_net_801f0
  fi

  echo "Starting container $i"
  echo "Configuration file: $config_file"
  echo "VFIO address: $port_tx"
  echo "IP: $ip"

 container_id=$(docker run -d \
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
    -e VFIO_PORT_TX=$port_tx \
    --network=$network \
    --ip=$ip \
    --ipc=host \
    -v /dev/shm:/dev/shm \
    tiber-broadcast-suite-nmos-node config/$config_file)

  if [ $? -eq 0 ]; then
    echo "Container $i started successfully with ID $container_id."
    # Save logs to a file
    #docker logs "$container_id" > "container_${i}_logs.txt"
    #echo "Logs for container $i saved to container_${i}_logs.txt."
  else
    echo "Failed to start container $i."
  fi
  config_file="intel-node-tx-${i}${format}.json"

done 