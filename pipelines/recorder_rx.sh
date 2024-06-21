#!/bin/bash

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
   --ip=192.168.2.2 \
   --expose=20000-20170 \
   --ipc=host -v /dev/shm:/dev/shm \
   --cpuset-cpus=0-10 \
   -e MTL_PARAM_LCORES="5-10" \
   -e MTL_PARAM_DATA_QUOTA=10356 \
      video_production_image -y \
      -qsv_device /dev/dri/renderD128 -hwaccel qsv \
      -p_port 0000:4b:01.1 -p_sip 192.168.2.2 -p_rx_ip 192.168.2.1 -udp_port 20000 -payload_type 112 -fps 25 -pix_fmt yuv422p10le \
      -video_size 3840x2160 -f mtl_st20p -i "0" \
      -filter_complex "[0:v]split=2[in1][in2];\
          [in1]hwupload,scale_qsv=iw/2:ih/2[out1]; \
          [in2]hwupload,scale_qsv=iw/4:ih/4[out2]" \
      -map "[out1]" -c:v hevc_qsv /videos/recv-quarter.mp4 \
      -map "[out2]" -c:v hevc_qsv /videos/recv-sixteenth.mp4
