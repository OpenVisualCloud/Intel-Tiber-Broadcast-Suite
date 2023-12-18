#!/bin/bash
CAM_NIC_PORT="0000:b1:01.1"
CAM_LOCAL_IP_ADDRESS="192.168.2.1"
CAM_DEST_IP_ADDRESS="192.168.2.2"

docker run \
  --user root\
  --privileged \
  --device=/dev/vfio:/dev/vfio \
  --device=/dev/dri:/dev/dri \
  --cap-add ALL \
  -v $(pwd):/videos \
  -v /usr/lib/x86_64-linux-gnu/dri \
  -v /tmp/kahawai_lcore.lock:/tmp/kahawai_lcore.lock \
  -v /dev/null:/dev/null \
  -v /tmp/hugepages:/tmp/hugepages \
  -v /hugepages:/hugepages \
  --network=my_net_801f0 \
  --ip=192.168.2.1 \
  --expose=20000-20170 \
  --ipc=host -v /dev/shm:/dev/shm \
  --cpuset-cpus="28-55" \
  -e MTL_PARAM_LCORES="28-55" \
  -e MTL_PARAM_DATA_QUOTA=10356 \
  my_ffmpeg \
  -qsv_device /dev/dri/renderD128 \
  -hwaccel qsv -f rawvideo -pix_fmt y210le -s:v 3840x2160 -r 25 -i /videos/gradients.yuv \
  -hwaccel qsv -f rawvideo -pix_fmt y210le -s:v 3840x2160 -r 25 -i /videos/gradients1.yuv \
  -hwaccel qsv -f rawvideo -pix_fmt y210le -s:v 3840x2160 -r 25 -i /videos/gradients2.yuv \
  -hwaccel qsv -f rawvideo -pix_fmt y210le -s:v 3840x2160 -r 25 -i /videos/gradients3.yuv \
  -noauto_conversion_filters \
  -filter_complex "
    [0:v]scale_qsv=w=iw/2:h=ih/2[tile0];
    [1:v]scale_qsv=w=iw/2:h=ih/2[tile1];
    [2:v]scale_qsv=w=iw/2:h=ih/2[tile2];
    [3:v]scale_qsv=w=iw/2:h=ih/2[tile3];
    [tile0][tile1][tile2][tile3]xstack_qsv=inputs=4:layout=0_0|0_h0|w0_0|w0_h0[out];" \
  -map [out] -f rawvideo -pix_fmt y210le /videos/outfile_GPU.yuv
