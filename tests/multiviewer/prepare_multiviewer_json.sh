#!/bin/bash

# Prepare intel-node-multiviewer.json

echo "Usage $0 <number_of_receivers>"

file_to_delete="intel-node-multiviewer.json"
if [ -e "$file_to_delete" ]; then
  # If the file exists, delete it
  rm "$file_to_delete"
fi

number_of_receivers=$1

if [ -z "$number_of_receivers" ]; then
  echo "Usage: $0 <number of receivers blocks to configure>"
  exit 1
fi

#Json block for sender
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
blocks=$(for i in $(seq 1 $number_of_receivers); do echo "$json_block"; done | jq -s '.')


jq --argjson blocks "$blocks" --argjson count "$number_of_receivers" \
   '.receiver += $blocks | .receivers_count = [$count] | .multiviewer_columns = $count' \
   intel-node-multiviewer_example.json > intel-node-multiviewer.json
