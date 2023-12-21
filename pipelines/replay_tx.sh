#!/bin/bash
NIC_PORT="0000:b1:01.1"
LOCAL_IP_ADDRESS="192.168.2.1"
DEST_IP_ADDRESS="192.168.2.2"
INPUT_FILE_NAME="input_test_4k.mkv"

if [ ! -e $INPUT_FILE_NAME ]; then
  docker run -it \
    --user root\
    --privileged \
    --device=/dev/dri:/dev/dri \
    -v $(pwd):/videos \
    -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
    my_ffmpeg \
    -hwaccel qsv -hwaccel_device /dev/dri/renderD128 \
    -y -f lavfi -i testsrc=d=10:s=3840x2160:r=50,format=y210le \
    -c:v libx265 -an -x265-params crf=25 /videos/$INPUT_FILE_NAME
fi

docker run -it \
  --user root\
  --privileged \
  --device=/dev/vfio:/dev/vfio \
  --device=/dev/dri:/dev/dri \
  --cap-add ALL \
  -v $(pwd):/videos \
  -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
  -v /tmp/kahawai_lcore.lock:/tmp/kahawai_lcore.lock \
  -v /dev/null:/dev/null \
  -v /tmp/hugepages:/tmp/hugepages \
  -v /hugepages:/hugepages \
  --network=my_net_801f0 \
  --ip=192.168.2.1 \
  --expose=20000-20170 \
  --ipc=host -v /dev/shm:/dev/shm \
  --cpuset-cpus="32-63" \
  -e MTL_PARAM_LCORES="59-63" \
  -e MTL_PARAM_DATA_QUOTA=10356 \
  my_ffmpeg \
  -y \
  -an \
  -i /videos/$INPUT_FILE_NAME \
  -vframes 500 -port $NIC_PORT -local_addr $LOCAL_IP_ADDRESS -dst_addr $DEST_IP_ADDRESS -udp_port 20000 -total_sessions 1 -f kahawai_mux -\
