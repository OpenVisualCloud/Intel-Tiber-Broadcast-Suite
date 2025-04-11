#!/bin/bash
# Prepare intel-node-tx.json configurations intel-node-tx-2.json, intel-node-tx-3.json, ...
echo "Usage $0 <number_of_additional_configuration_files_to_create> <localhost/ip>"

file_to_delete="intel-node-tx.json"
if [ -e "$file_to_delete" ]; then
  # If the file exists, delete it
  rm "$file_to_delete"
fi

number_of_configs=$1
host=$2

if [ -z "$number_of_configs" ]; then
  echo "Usage: $0 <number of additional tx configurations files todo>"
  exit 1
fi

if [ -z "$host" ]; then
  echo "Usage: $0 $1 <localhost/ip>"
  exit 1
fi

base_port=50090
base_http_port=100

if [[ "$host" == "localhost" ]]; then
  jq '.ffmpeg_grpc_server_port = "50055" | .http_port=90 | .ffmpeg_grpc_server_address = "localhost"' \
    intel-node-tx_example.json > intel-node-tx.json
else
  jq '.ffmpeg_grpc_server_port = "50055" | .http_port=90' \
    intel-node-tx_example.json > intel-node-tx.json
fi

for ((i=1; i<=number_of_configs; i++)); do 
   new_file="1920x1080p10le_multi_${i}.yuv"
   port=\"$((base_port + i))\"
   http=$((base_http_port + i * 5))
   ip=$(awk "NR==${i}" IP_launcher.txt)

  if [[ "$host" == "localhost" ]]; then
    jq --arg new_filename "$new_file" --argjson port_number "$port" --argjson http_port "$http" \
      '.ffmpeg_grpc_server_port = $port_number | .http_port = $http_port | .receiver[0].stream_type.file.filename = $new_filename' \
      intel-node-tx.json > intel-node-tx-${i}.json
  else
    jq --arg new_filename "$new_file" --argjson port_number "$port" --argjson http_port "$http" --arg ip_add "$ip" \
     '.ffmpeg_grpc_server_port = $port_number | .http_port = $http_port | .receiver[0].stream_type.file.filename = $new_filename | .ffmpeg_grpc_server_address = $ip_add' \
     intel-node-tx.json > intel-node-tx-${i}.json
  fi
done
