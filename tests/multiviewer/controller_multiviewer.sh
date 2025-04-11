#!/bin/bash
echo "Usage $0 <number_of_receivers> <localhost/ip>"

number_of_receivers=$1
host=$2

if [ -z "$number_of_receivers" ]; then
  echo "Usage: $0 <number of receivers>"
  exit 1
fi

if [ -z "$host" ]; then
  echo "Usage: $0 $1 <localhost/ip>"
  exit 1
fi

base_http_port=100
cd ../../src/nmos/nmos-is05-controller

for ((i=1; i<=number_of_receivers; i++)); do 
   http=$((base_http_port + (i-1) * 5))
  if [[ "$i" -eq 1 ]]; then
    http=90
  fi

  if [[ "$host" == "localhost" ]]; then
    receiver_ip="localhost"
    sender_ip="localhost"
  else
    receiver_ip="192.168.2.7"
    sender_ip=$(awk "NR==${i}" ../../../tests/multisender/IP_node.txt)
  fi

  python3 threaded-nmos-controller05.py --receiver_ip $receiver_ip --sender_ip $sender_ip --receiver_port 95 --sender_port $http
done