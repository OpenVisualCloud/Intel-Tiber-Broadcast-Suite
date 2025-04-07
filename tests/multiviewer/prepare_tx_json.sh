#!/bin/bash

# Prepare intel-node-tx.json configurations intel-node-tx-2.json, intel-node-tx-3.json, ...

echo "Usage $0 <number_of_configs>"
number_of_configs=$1

if [ -z "$number_of_configs" ]; then
  echo "Usage: $0 <number of additional tx configurations files todo>"
  exit 1
fi

base_port=50090
base_http_port=100

for ((i=1; i<=number_of_configs; i++)); do 
   new_file="1920x1080p10le_multi_${i}.yuv"
   port=\"$((base_port + i))\"
   http=$((base_http_port + i * 5))
   
   jq --arg new_filename "$new_file" --argjson port_number "$port" --argjson http_port "$http" \
    '.ffmpeg_grpc_server_port = $port_number | .http_port = $http_port | .receiver[0].stream_type.file.filename = $new_filename' \
    intel-node-tx.json > intel-node-tx-${i}.json
done
