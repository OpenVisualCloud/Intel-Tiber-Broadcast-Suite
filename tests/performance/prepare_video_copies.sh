#!/bin/bash

echo "Usage $0 <number_of_videos> <video_format>, output_1.format need to be already prepared"



# Prepare intel-node-rx.json with many senders and receivers
number=$1
format=$2


if [ -z "$number" ]; then
  echo "Usage: $0 <number of videos>"
  exit 1
fi

if [ -z "$format" ]; then
  echo "Usage: $0 $1 <video_format>"
  exit 1
fi


for ((i=1; i<=number; i++)); do 
  cp output_1.${format} output_$((i+1)).${format}
done

