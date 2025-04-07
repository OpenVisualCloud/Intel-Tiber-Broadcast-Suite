#!/bin/bash
echo "Usage $0 <number_of_receivers>"

# Prepare intel-node-rx.json configurations intel-node-rx-2.json, intel-node-rx-3.json, ...
number_of_receivers=$1

if [ -z "$number_of_receivers" ]; then
  echo "Usage: $0 <number of receivers>"
  exit 1
fi

base_http_port=100
cd ../../src/nmos/nmos-is05-controller
python3 threaded-nmos-controller05.py --receiver_ip localhost --sender_ip localhost --receiver_port 95 --sender_port 90

for ((i=1; i<number_of_receivers; i++)); do 
   http=$((base_http_port + i * 5))
   python3 threaded-nmos-controller05.py --receiver_ip localhost --sender_ip localhost --receiver_port 95 --sender_port $http
done