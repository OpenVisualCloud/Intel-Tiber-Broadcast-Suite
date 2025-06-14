#!/bin/bash

echo "Usage $0 <number_of_receivers/senders> <video_type> <IP>"

if [ -z "$1" ] || [ -z "$2" ] || [ -z "$3" ]; then
  echo "Usage: $0 <number_of_receivers/senders> <video_type> <IP>"
  exit 1
fi

case $2 in
  "raw")
    format=""
    video_type="rawvideo"
    pixel_format="yuv422p10le"
    video_format="yuv"
    gpu_hw_acceleration="none"
    gpu_hw_acceleration_device=""
    ;;
  "h264")
    format="-h264"
    video_type="h264"
    pixel_format="yuv420p"
    video_format="h264"
    gpu_hw_acceleration="none"
    gpu_hw_acceleration_device=""
    preset="ultrafast"
    ;;
  "h265")
    format="-h265"
    video_type="hevc_qsv"
    pixel_format="yuv422p10le"
    video_format="h265"
    gpu_hw_acceleration="intel"
    gpu_hw_acceleration_device="/dev/dri/renderD128"
    preset="veryfast"
    ;;
  *)
    echo "Invalid option: $2"
    echo "Valid options are: raw, h264, h265"
    exit 1
    ;;
esac

file_to_delete="intel-node-rx${format}.json"
if [ -e "$file_to_delete" ]; then
  # If the file exists, delete it
  rm "$file_to_delete"
fi

# Prepare intel-node-rx.json with many senders and receivers
number=$1
queues=16

if [ "$number" -gt 16 ]; then
  queues=$number
fi

# Input parameters
port=50056
http=95
ip=$3
numerator=50


#Json block for receiver
json_receiver_block="{
  \"stream_payload\": {
    \"video\": {
      \"frame_width\": 1920,
      \"frame_height\": 1080,
      \"frame_rate\": { \"numerator\": ${numerator}, \"denominator\": 1 },
      \"pixel_format\": \"yuv422p10le\",
      \"video_type\": \"rawvideo\"
    },
    \"audio\": {
      \"channels\": 2,
      \"sampleRate\": 48000,
      \"format\": \"pcm_s24be\",
      \"packetTime\": \"1ms\"
    }
  },
  \"stream_type\": {
    \"st2110\" : {
      \"transport\" : \"st2110-20\",
      \"payloadType\" :  112,
      \"queues_cnt\": ${queues}
    }
  }
}"

json_sender_blocks=()

for ((i=1; i<=number; i++)); do 
  new_file="output_${i}.${video_format}"

#Json block for sender
case $3 in
    "raw")
      json_sender_block="{
        \"stream_payload\": {
          \"video\": {
            \"frame_width\": 1920,
            \"frame_height\": 1080,
            \"frame_rate\": {
              \"numerator\": ${numerator},
              \"denominator\": 1
            },
            \"pixel_format\": \"${pixel_format}\",
            \"video_type\": \"${video_type}\"
          },
          \"audio\": {
            \"channels\": 2,
            \"sampleRate\": 48000,
            \"format\": \"pcm_s24be\",
            \"packetTime\": \"1ms\"
          }
        },
        \"stream_type\": {
          \"file\": {
            \"path\": \"/videos/recv\",
            \"filename\": \"${new_file}\"
          }
        }
      }"
    ;;
    *)
      json_sender_block="{
        \"stream_payload\": {
          \"video\": {
            \"frame_width\": 1920,
            \"frame_height\": 1080,
            \"frame_rate\": {
              \"numerator\": ${numerator},
              \"denominator\": 1
            },
            \"pixel_format\": \"${pixel_format}\",
            \"video_type\": \"${video_type}\",
            \"preset\": \"${preset}\",
            \"profile\": \"main\"
          },
          \"audio\": {
            \"channels\": 2,
            \"sampleRate\": 48000,
            \"format\": \"pcm_s24be\",
            \"packetTime\": \"1ms\"
          }
        },
        \"stream_type\": {
          \"file\": {
            \"path\": \"/videos/recv\",
            \"filename\": \"${new_file}\"
          }
        }
      }"
    ;;
  esac
  json_sender_blocks+=("$json_sender_block")
done


# Construct blocks
receiver_blocks=$(for _ in $(seq 1 $number); do echo "$json_receiver_block"; done | jq -s '.')
sender_blocks=$(printf "%s\n" "${json_sender_blocks[@]}" | jq -s '.')


jq --argjson blocksR "$receiver_blocks" --argjson blocksS "$sender_blocks" --argjson count "$number" --arg port_number "$port" --argjson http_port $http --arg ip_add "$ip" --arg gpu_hw_acceleration "$gpu_hw_acceleration" --arg gpu_hw_acceleration_device "$gpu_hw_acceleration_device" \
    '.receiver += $blocksR | .sender += $blocksS | .ffmpeg_grpc_server_port = $port_number | .http_port = $http_port | .ffmpeg_grpc_server_address = $ip_add | .gpu_hw_acceleration = $gpu_hw_acceleration | .gpu_hw_acceleration_device = $gpu_hw_acceleration_device' \
    intel-node-rx_example.json > intel-node-rx${format}.json

