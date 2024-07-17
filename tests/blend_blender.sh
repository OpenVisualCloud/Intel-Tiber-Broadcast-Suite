#!/bin/bash
RECV_NIC_PORT="0000:b1:01.2"
RECV_LOCAL_IP_ADDRESS="192.168.2.2"
RECV_SOURCE_IP_ADDRESS="192.168.2.1"

docker run \
  --user root\
  --privileged \
  --device=/dev/vfio:/dev/vfio \
  --device=/dev/dri:/dev/dri \
  --cap-add ALL \
  -v $(pwd):/videos \
  -v /tmp/kahawai_lcore.lock:/tmp/kahawai_lcore.lock \
  -v /var/run/imtl:/var/run/imtl \
  -v /dev/null:/dev/null \
  -v /tmp/hugepages:/tmp/hugepages \
  -v /hugepages:/hugepages \
  --network=my_net_801f0 \
  --ip=192.168.2.2 \
  --expose=20000-20170 \
  --ipc=host -v /dev/shm:/dev/shm \
  --cpuset-cpus="84-111" \
  video_production_image \
  -y \
  -an \
  -thread_queue_size 4096 \
  -framerate 50 -t 2 -pixel_format y210le -pix_fmt y210le -width 3840 -height 2160 -port $RECV_NIC_PORT -local_addr $RECV_LOCAL_IP_ADDRESS -src_addr $RECV_SOURCE_IP_ADDRESS -udp_port 20000 -total_sessions 1 -f kahawai -i "0" \
  -thread_queue_size 4096 \
  -video_size 3840x2160 -framerate 50 -t 2 -pix_fmt y210le -i /videos/digit2.yuv \
  -f rawvideo -filter_complex "[0]scale=3840:2160,framerate=50,format=y210le,setpts=PTS-STARTPTS,fps=50[in1];[1]scale=3840:2160,framerate=50,format=y210le,setpts=PTS-STARTPTS,fps=50[in2];[in1]fade=t=out:d=2:alpha=1, setpts=PTS-STARTPTS, fps=50[fade1];[in2]fade=t=in:d=2:alpha=1, setpts=PTS-STARTPTS, fps=50[fade2];[fade1][fade2]overlay, setpts=PTS-STARTPTS, fps=50[out];" \
  -map [out] -f rawvideo -vframes 100 -t 2 -framerate 50 -pixel_format y210le -pix_fmt y210le -width 3840 -height 2160 /videos/blend_output.yuv
