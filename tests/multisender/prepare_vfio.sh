#!/bin/bash

echo "Usage $0 <address> <number_of_vfio_to_create>"

if [ -z "$1" ]; then
  echo "Usage: $0 <0000:27:00.0>"
  exit 1
fi

if [ -z "$2" ]; then
  echo "Usage: $0 $1 <number_of_vfio>"
  exit 1
fi

addres=$1
number_of_vfio=$2

nicctl.sh create_vf $addres $number_of_vfio
nicctl.sh list all | awk '/[0-9a-f]{4}:[0-9a-f]{2}:[0-9a-f]{2}\.[0-9a-f]/ {print $2}' | head -n $number_of_vfio > vfio_addresses.txt