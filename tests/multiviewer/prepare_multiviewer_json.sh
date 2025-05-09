#!/bin/bash

echo "Usage $0 <number_of_receivers> <localhost/ip>"

file_to_delete="intel-node-multiviewer.json"
if [ -e "$file_to_delete" ]; then
  # If the file exists, delete it
  rm "$file_to_delete"
fi

# Prepare intel-node-multiviewer.json with many receivers
number_of_receivers=$1
host=$2

if [ -z "$number_of_receivers" ]; then
  echo "Usage: $0 <number of receivers blocks to configure>"
  exit 1
fi

if [ -z "$host" ]; then
  echo "Usage: $0 $1 <localhost/ip>"
  exit 1
fi

#Json block for receiver
json_block='{
    "stream_payload": {
      "video": {
        "frame_width": 1920,
        "frame_height": 1080,
        "frame_rate": { "numerator": 50, "denominator": 1 },
        "pixel_format": "yuv422p10le",
        "video_type": "rawvideo"
      },
      "audio": {
        "channels": 2,
        "sampleRate": 48000,
        "format": "pcm_s24be",
        "packetTime": "1ms"
      }
    },
    "stream_type": {
        "st2110" : {
          "transport" : "st2110-20",
          "payloadType" :  112
        }
      }
  }'

# Construct blocks
blocks=$(for _ in $(seq 1 $number_of_receivers); do echo "$json_block"; done | jq -s '.')

if [[ "$host" == "localhost" ]]; then
  jq --argjson blocks "$blocks" --argjson count "$number_of_receivers" \
     '.receiver += $blocks | .multiviewer_columns = $count' \
     intel-node-multiviewer_example.json > intel-node-multiviewer.json
else
  jq --argjson blocks "$blocks" --argjson count "$number_of_receivers" \
     '.receiver += $blocks | .multiviewer_columns = $count | .ffmpeg_grpc_server_address = "192.168.2.6"' \
     intel-node-multiviewer_example.json > intel-node-multiviewer.json
fi
