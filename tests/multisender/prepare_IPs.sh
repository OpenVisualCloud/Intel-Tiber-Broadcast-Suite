#!/bin/bash

echo "Usage $0 <number_of_IPs>"

number_of_IPs=$1

if [ -z "$1" ]; then
  echo "Usage: $0 <number_of_IPs>"
  exit 1
fi

# Check and delete the second file

file_to_delete="IP_launcher.txt"
file_to_delete1="IP_receiver.txt"

if [ -e "$file_to_delete" ]; then
  rm "$file_to_delete"
fi

if [ -e "$file_to_delete1" ]; then
  rm "$file_to_delete1"
fi

for ((i=1; i<=number_of_IPs*2; i++)); do
  ip=$((i + 9))
  if ((ip % 2 == 0)); then
    echo "192.168.2.${ip}" >> IP_launcher.txt
  else
    echo "192.168.2.${ip}" >> IP_node.txt
  fi
done
