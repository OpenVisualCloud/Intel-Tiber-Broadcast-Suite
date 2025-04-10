#!/bin/bash
echo "Usage $0 <number_of_senders> <localhost/ip>"

# Prepare intel-node-rx.json configurations intel-node-rx-2.json, intel-node-rx-3.json, ...
number_of_senders=$1
host=$2

if [ -z "$number_of_senders" ]; then
  echo "Usage: $0 <number of senders>"
  exit 1
fi

if [ -z "$host" ]; then
  echo "Usage: $0 $1 <localhost/ip>"
  exit 1
fi

base_http_port=100
cd ../../src/nmos/nmos-is05-controller
 # python3 threaded-nmos-controller05.py --receiver_ip localhost --sender_ip localhost --receiver_port 95 --sender_port 90

for ((i=1; i<=number_of_senders; i++)); do 
  http=$((base_http_port + (i-1) * 5))
  if [[ "$i" -eq 1 ]]; then
    http=95
  fi

  if [[ "$host" == "localhost" ]]; then
    receiver_ip="localhost"
    sender_ip="localhost"
  else
    receiver_ip=$(awk "NR==${i}" ../../../tests/multisender/IP_node.txt)
    sender_ip="192.168.2.5"
  fi

  python3 threaded-nmos-controller05.py --receiver_ip $receiver_ip --sender_ip $sender_ip --receiver_port $http --sender_port 90
done
