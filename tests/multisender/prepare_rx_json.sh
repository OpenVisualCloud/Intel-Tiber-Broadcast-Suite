#!/bin/bash
echo "Usage $0 <number_of_additional_configuration_files_to_create>"

# Prepare intel-node-rx.json configurations intel-node-rx-2.json, intel-node-rx-3.json, ...
number_of_configs=$1

if [ -z "$number_of_configs" ]; then
  echo "Usage: $0 <number of rx configurations files todo>"
  exit 1
fi

base_port=50056
base_http_port=100

cp -n ../intel-node-rx.json intel-node-rx-example.json

jq '.ffmpeg_grpc_server_port = "50056" | .ffmpeg_grpc_server_address = "localhost"' \
  intel-node-rx-example.json > intel-node-rx.json

for ((i=1; i<=number_of_configs; i++)); do 
   new_file="multisenders_${i}.yuv"
   port=\"$((base_port + i))\"
   http=$((base_http_port + i * 5))


   jq --arg new_filename "$new_file" --argjson port_number "$port" --argjson http_port "$http" \
    '.ffmpeg_grpc_server_port = $port_number | .http_port = $http_port | .sender[0].stream_type.file.filename = $new_filename' \
    intel-node-rx.json > intel-node-rx-${i}.json

done

file_to_delete="intel-node-rx-example.json"
if [ -e "$file_to_delete" ]; then
  # If the file exists, delete it
  rm "$file_to_delete"
fi