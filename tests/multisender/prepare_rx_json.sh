#!/bin/bash
echo "Usage $0 <number_of_additional_configuration_files_to_create> <localhost/ip>"

# Prepare intel-node-rx.json configurations intel-node-rx-2.json, intel-node-rx-3.json, ...
number_of_configs=$1
host=$2

if [ -z "$number_of_configs" ]; then
  echo "Usage: $0 <number of rx configurations files todo>"
  exit 1
fi

if [ -z "$host" ]; then
  echo "Usage: $0 $1 <localhost/ip>"
  exit 1
fi

base_port=50056
base_http_port=100

cp -n ../intel-node-rx.json intel-node-rx-example.json

if [[ "$host" == "localhost" ]]; then
  jq '.ffmpeg_grpc_server_port = "50056" | .ffmpeg_grpc_server_address = "localhost"' \
    intel-node-rx-example.json > intel-node-rx.json
else
    jq '.ffmpeg_grpc_server_port = "50056"' \
    intel-node-rx-example.json > intel-node-rx.json
fi

for ((i=1; i<=number_of_configs; i++)); do 
   new_file="multisenders_${i}.yuv"
   port=\"$((base_port + i))\"
   http=$((base_http_port + i * 5))
   ip=$(awk "NR==${i}" IP_launcher.txt)

  if [[ "$host" == "localhost" ]]; then
    jq --arg new_filename "$new_file" --argjson port_number "$port" --argjson http_port "$http" \
      '.ffmpeg_grpc_server_port = $port_number | .http_port = $http_port | .sender[0].stream_type.file.filename = $new_filename' \
      intel-node-rx.json > intel-node-rx-${i}.json
  else
    jq --arg new_filename "$new_file" --argjson port_number "$port" --argjson http_port "$http" --arg ip_add "$ip" \
      '.ffmpeg_grpc_server_port = $port_number | .http_port = $http_port | .sender[0].stream_type.file.filename = $new_filename | .ffmpeg_grpc_server_address = $ip_add' \
      intel-node-rx.json > intel-node-rx-${i}.json
  fi
done

file_to_delete="intel-node-rx-example.json"
if [ -e "$file_to_delete" ]; then
  # If the file exists, delete it
  rm "$file_to_delete"
fi