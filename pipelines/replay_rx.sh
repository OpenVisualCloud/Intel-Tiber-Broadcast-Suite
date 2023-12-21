#!/bin/bash
RECV_NIC_PORT="0000:b1:01.2"
RECV_LOCAL_IP_ADDRESS="192.168.2.2"
RECV_SOURCE_IP_ADDRESS="192.168.2.1"
OUTPUT_FILE_NAME="output_test"


docker run \
  --user root\
  --privileged \
  --device=/dev/vfio:/dev/vfio \
  --device=/dev/dri:/dev/dri \
  --cap-add ALL \
  -v $(pwd):/videos \
  -v /tmp/kahawai_lcore.lock:/tmp/kahawai_lcore.lock \
  -v /dev/null:/dev/null \
  -v /tmp/hugepages:/tmp/hugepages \
  -v /hugepages:/hugepages \
  --network=my_net_801f0 \
  --ip=192.168.2.2 \
  --expose=20000-20170 \
  --ipc=host -v /dev/shm:/dev/shm \
  --cpuset-cpus="84-111" \
  -e MTL_PARAM_LCORES="84-111" \
  -e MTL_PARAM_DATA_QUOTA=2589 \
  my_ffmpeg \
  -y -an \
  -pixel_format y210le -width 3840 -height 2160 -port $RECV_NIC_PORT -local_addr $RECV_LOCAL_IP_ADDRESS -src_addr $RECV_SOURCE_IP_ADDRESS \
    -udp_port 20000 -total_sessions 1 -f kahawai -i "0" \
  -map 0:v -vframes 500 -r 50 -vf scale=960:540 -c:v libx265 -an -x265-params crf=25 /videos/$OUTPUT_FILE_NAME"_540p.mkv" \
  -map 0:v -vframes 500 -r 50 -vf scale=480:270 -c:v libx265 -an -x265-params crf=25 /videos/$OUTPUT_FILE_NAME"_270p.mkv" 

