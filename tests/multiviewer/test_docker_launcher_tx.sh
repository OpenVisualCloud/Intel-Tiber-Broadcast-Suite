#!/bin/bash

echo "Usage $0 <host_name> <number_of_containers>"

if [ -z "$1" ]; then
  echo "Usage: $0 <hostname/ip> e.g. $0 192.168.2.1"
  exit 1
fi

number_of_containers=$2

if [ -z "$2" ]; then
  echo "Usage: $0 $1  <number_of_containers>"
  exit 1
fi

host=$1
base_port=50089

for ((i=1; i<=number_of_containers; i++)); do 
  if [[ "$host" == "localhost" ]]; then
    ip=""
    network=host
    HOSTNAME=localhost
  else
    ip=$(awk "NR==$((i-1))" IP_launcher.txt)
    HOSTNAME=$(awk "NR==$((i-1))" IP_launcher.txt)
    if [[ "$i" -eq 1 ]]; then
        ip=192.168.2.4
        HOSTNAME=192.168.2.4
    fi
    network=my_net_801f0
  fi

  if [ "$i" -eq 1 ]; then
    PORT=50055
  else
    PORT=$((base_port + i))
  fi

  echo "Starting container $i"
  echo "Host name: $HOSTNAME"
  echo "Port: $PORT"
  echo "IP: $ip"

  container_id=$(docker run -d \
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
     tiber-broadcast-suite "$HOSTNAME" "$PORT")

  if [ $? -eq 0 ]; then
    echo "Container $i started successfully with ID $container_id."
    # Save logs to a file
    #docker logs "$container_id" > "container_launcher_tx_${i}_logs.txt"
    #echo "Logs for container $i saved to container_launcher_tx_${i}_logs.txt."
  else
    echo "Failed to start container $i."
  fi
done