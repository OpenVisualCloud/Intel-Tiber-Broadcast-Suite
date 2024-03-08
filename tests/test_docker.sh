#!/bin/bash 

INPUT_FILE_NAME="input_video.yuv"
INPUT_FILE_NAME_MP4="input_video.mp4"
REF_FILE_NAME="reference_video.yuv"
OUTPUT_FILE_NAME="output_video.yuv"
OUTPUT_FILE_NAME_MP4="output_video.mp4"
REF_FILE_NAME_MP4="reference_video.mp4"

rm -rf *.yuv
rm -rf *.mp4

########### PSNR MINIMAL VALUE TO PASS TESTS ###########
PSNR_LIM=38

########### HELP MESSAGE ###########
help_message_arg_1 () {
   echo "Please chose one option as first argument:"
   echo "   imtl           -  IMTL transmition only test"
   echo "   qsv            -  QSV filters only test"
   echo "   jxs            -  JPEG-XS codec only test"
   echo "   vsr            -  VSR codec only test"
   echo "   imtl_qsv       -  IMTL transmition and QSV filters test"
   echo "   imtl_jxs       -  IMTL transmition and JPEG-XS coding test"
   echo "   imtl_qsv_vsr   -  IMTL transmition and JPEG-XS coding test"
   echo "   imtl_qsv_jxs   -  IMTL transmition, QSV filters and JPEG-XS coding test"
}

########### EQUALITY TEST ARGS: 1: resolution | 2: framerate | 3: pixel format | 4: output video ###########
equality_test() {
   if [ "$6" = "rawvideo" ]; then
      docker run \
      --user root\
      --privileged \
      -v $(pwd):/videos \
      -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
      video_production_image -hide_banner -loglevel quiet -an -y -video_size $2 -f rawvideo -pix_fmt $3 -i /videos/$4 \
      -video_size $2 -f rawvideo -pix_fmt $3 -i /videos/$5 \
      -filter_complex libvmaf=log_path=/videos/vmaf_output.xml -f null -

      docker run \
      --user root\
      --privileged \
      -v $(pwd):/videos \
      -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
      video_production_image -hide_banner -loglevel quiet -video_size $2 -pixel_format $3 -i /videos/$4 \
      -video_size $2 -pixel_format $3 -i /videos/$5 -filter_complex "psnr=f=/videos/psnr_out.log" -f null /dev/null
   else
      docker run \
      --user root\
      --privileged \
      -v $(pwd):/videos \
      -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
      video_production_image -hide_banner -loglevel quiet -an -y -i /videos/$4 -i /videos/$5 \
      -filter_complex libvmaf=log_path=/videos/vmaf_output.xml -f null -

      docker run \
      --user root\
      --privileged \
      -v $(pwd):/videos \
      -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
      video_production_image -hide_banner -loglevel quiet -i /videos/$4 -i /videos/$5 \
      -filter_complex "psnr=f=/videos/psnr_out.log" -f null /dev/null 
   fi

   psnr_array=$(cat psnr_out.log |grep -o -P 'psnr_avg:.{0,5}' | awk -F":" '{print $2}' | awk -F"." '{print $1}')

   for element in $psnr_array
   do
      if [ $((element)) -ge $PSNR_LIM ]; then
         psnr_err_val=$(($psnr_err_val+1))
      fi
   done

   echo "****** TEST $1 | $2 | $3 ******"
   if [ ! -e $5 ] || [ ! -e $4 ]; then
      echo -e "****** \033[0;31mTEST FAILED\033[0m ******      |  file:$5 or $4 not exist ******"
   else
      hash_ref=$(md5sum $5  | awk '{print $1}')
      hash_test=$(md5sum $4  | awk '{print $1}')
      if [ $hash_ref == $hash_test ] || [ $psnr_err_val > 0 ]; then
         echo -e "****** \033[0;32mTEST PASSED\033[0m ******"
      else

         echo -e "****** \033[0;31mTEST FAILED\033[0m ******   |  file:$5 and $4 are different ******"
      fi
   fi
}

########### TEST FILE PREPARATION ARGS: 1: resolution | 2: framerate | 3: pixel format | 4: output video ###########
prepare_reference_file() {
   rm -rf $5
   if [ ! -e $5 ]; then
      if [ $4 == "rawvideo" ]; then
         docker run -it \
            --user root\
            --privileged \
            -v $(pwd):/videos \
            -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
            video_production_image -hide_banner -loglevel quiet -threads $(nproc) -an -y -f lavfi -i testsrc=d=2:s=$1:r=$2,format=$3 -f $4 /videos/$5
      else
         docker run -it \
            --user root\
            --privileged \
            --device=/dev/dri:/dev/dri \
            -v $(pwd):/videos \
            -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
            video_production_image -hide_banner -loglevel quiet -threads $(nproc) -qsv_device /dev/dri/renderD128 -an -y -hwaccel qsv -f lavfi -i testsrc=d=2:s=$1:r=$2,format=$3 -c:v $4 /videos/$5 > logs.txt
      fi
   fi
}

########### SCALE DOWN TEST FILE ARGS: 1: codec | 2: pixel format | 3: pixel format | 4: file name ###########
prepare_input_for_vsr() {
   rm -rf $4
   docker run -it \
      --user root\
      --privileged \
      --device=/dev/dri:/dev/dri \
      -v $(pwd):/videos \
      -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
      video_production_image -hide_banner -loglevel quiet -threads $(nproc) -y -an \
      -qsv_device /dev/dri/renderD128 \
      -hwaccel qsv -c:v $1 -i /videos/$3 -vf "hwupload=extra_hw_frames=64,scale_qsv=w=iw/2:h=ih/2,hwdownload" -c:v $1 -pixel_format $2 /videos/$4  > logs.txt
}

prepare_input_for_vsr_rawvideo() {
   rm -rf $5
   docker run -it \
      --user root\
      --privileged \
      --device=/dev/dri:/dev/dri \
      -v $(pwd):/videos \
      -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
      video_production_image -hide_banner -loglevel quiet -threads $(nproc) -y -an -qsv_device /dev/dri/renderD128 \
      -hwaccel qsv -f $1 -video_size $3 -pix_fmt $2 -i /videos/$4 \
      -vf "hwupload=extra_hw_frames=64,scale_qsv=w=iw/2:h=ih/2,hwdownload,format=$2" -f $1 -pix_fmt $2 /videos/$5 > logs.txt
}
########### END OF TEST FILE PREPARATION ###########

run_qsv () {
   docker run -it \
      --user root\
      --privileged \
      -v $(pwd):/videos \
      -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
      video_production_image -hide_banner -loglevel quiet -an -y -threads $(nproc) -video_size $1 -f rawvideo -pix_fmt $3 -i /videos/$INPUT_FILE_NAME \
      -filter_complex "[0:v]format=$3,fps=$2,split=4[in1][in2][in3][in4]; \
      [in1]scale=iw/2:ih/2:flags=area[f_out_1]; \
      [in2]scale=iw/2:ih/2:flags=area[f_out_2]; \
      [in3]scale=iw/2:ih/2:flags=area[f_out_3]; \
      [in4]scale=iw/2:ih/2:flags=area[f_out_4]; \
      [f_out_1][f_out_2][f_out_3][f_out_4]xstack=inputs=4:layout=0_0|0_h0|w0_0|w0_h0[out_mutiview]" -sws_flags area -map "[out_mutiview]" -f rawvideo -pix_fmt $3 /videos/$REF_FILE_NAME > logs.txt 

   docker run -it \
      --user root\
      --privileged \
      --device=/dev/dri:/dev/dri \
      -v $(pwd):/videos \
      -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
      video_production_image -hide_banner -loglevel quiet -threads $(nproc) -y -an \
      -qsv_device /dev/dri/renderD128 \
      -hwaccel qsv -hwaccel_output_format qsv \
      -f rawvideo -pix_fmt $3 -s:v $1 \
      -i /videos/$INPUT_FILE_NAME \
      -noauto_conversion_filters \
      -filter_complex "[0:v]format=$3,fps=$2,split=4[in1][in2][in3][in4]; \
         [in1]hwupload=extra_hw_frames=1,scale_qsv=w=iw/2:h=ih/2:mode=compute:async_depth=1[tile0];\
         [in2]hwupload=extra_hw_frames=1,scale_qsv=w=iw/2:h=ih/2:mode=compute:async_depth=1[tile1];\
         [in3]hwupload=extra_hw_frames=1,scale_qsv=w=iw/2:h=ih/2:mode=compute:async_depth=1[tile2];\
         [in4]hwupload=extra_hw_frames=1,scale_qsv=w=iw/2:h=ih/2:mode=compute:async_depth=1[tile3];\
         [tile0][tile1][tile2][tile3]xstack_qsv=inputs=4:layout=0_0|0_h0|w0_0|w0_h0[multiview];[multiview]hwdownload,format=y210le[out_multiview]" \
      -map "[out_multiview]" -f rawvideo -pix_fmt $3 /videos/$OUTPUT_FILE_NAME > logs.txt
      # -filter_complex_frames 2 \
      # -filter_complex_policy 1 \
}

run_vsr () {
   if [ $1 == "qsv" ]; then
      docker run -it \
         --user root\
         --privileged \
         --device=/dev/dri:/dev/dri \
         -v $(pwd):/videos \
         -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
         video_production_image -hide_banner -loglevel quiet -threads $(nproc) -y -an -init_hw_device vaapi=va -init_hw_device qsv=qs@va -init_hw_device opencl=ocl@va -hwaccel qsv \
         -i /videos/$INPUT_FILE_NAME_MP4 \
         -vf "hwmap=derive_device=opencl,format=opencl,raisr_opencl,hwmap=derive_device=qsv:reverse=1:extra_hw_frames=16" \
         -c:v h264_qsv /videos/$OUTPUT_FILE_NAME_MP4 > logs.txt
   elif [ $1 == "cpu" ]; then
      docker run -it \
         --user root\
         --privileged \
         --device=/dev/dri:/dev/dri \
         -v $(pwd):/videos \
         -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
         video_production_image -hide_banner -loglevel quiet -threads $(nproc) -y -an -i /videos/$INPUT_FILE_NAME_MP4 \
         -vf "raisr=asm=opencl:threadcount=26:passes=2:filterfolder=filters_2x/filters_highres" /videos/$OUTPUT_FILE_NAME_MP4 > logs.txt
   fi
}

run_vsr_raw () {
   if [ $1 == "qsv" ]; then
      docker run -it \
         --user root\
         --privileged \
         --device=/dev/dri:/dev/dri \
         -v $(pwd):/videos \
         -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
         video_production_image -hide_banner -loglevel quiet -threads $(nproc) -y -an -init_hw_device vaapi=va -init_hw_device qsv=qs@va -init_hw_device opencl=ocl@va -hwaccel qsv \
         -f rawvideo -pixel_format $2 -video_size $3 -i /videos/$INPUT_FILE_NAME \
         -vf "hwmap=derive_device=opencl,format=opencl,raisr_opencl,hwmap=derive_device=qsv:reverse=1:extra_hw_frames=16" \
         -f rawvideo -pixel_format $2 -video_size $4 /videos/$OUTPUT_FILE_NAME
   elif [ $1 == "cpu" ]; then
      docker run -it \
         --user root\
         --privileged \
         --device=/dev/dri:/dev/dri \
         -v $(pwd):/videos \
         -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
         video_production_image -hide_banner -loglevel quiet -threads $(nproc) -y -an -f rawvideo -pixel_format $2 -video_size $3 -i /videos/$INPUT_FILE_NAME \
         -vf "raisr=asm=opencl:threadcount=26:passes=2:filterfolder=filters_2x/filters_highres" -f rawvideo -pixel_format $2 -video_size $4 /videos/$OUTPUT_FILE_NAME
   fi
}

run_jxs() {
   docker run \
      --user root\
      --privileged \
      --device=/dev/dri:/dev/dri \
      -v $(pwd):/videos \
      -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
      video_production_image -hide_banner -loglevel quiet -threads $(nproc) -y -i /videos/$INPUT_FILE_NAME_MP4 -c:v jpegxs -bpp 2 /videos/encoded_tmp.mov > logs.txt 

   docker run \
      --user root\
      --privileged \
      --device=/dev/dri:/dev/dri \
      -v $(pwd):/videos \
      -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
      video_production_image -hide_banner -loglevel quiet -y -threads $(nproc) -i /videos/encoded_tmp.mov -c:v $1 /videos/$OUTPUT_FILE_NAME_MP4 > logs.txt 

   rm -rf /videos/encoded_tmp.mov
}

########### TEST FILE PREPARATION ARGS: 1: width 2: heigh 3: framerate 4: pixel format 5: src IP 6: dst IP 7: CPUs 8: MTL CPUs ###########
run_imtl_rx_qsv_vsr () {
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
   --ip=$5 \
   --expose=20000-20170 \
   --ipc=host -v /dev/shm:/dev/shm \
   --cpuset-cpus=$7 \
   -e MTL_PARAM_LCORES=$8 \
   -e MTL_PARAM_DATA_QUOTA=10356 \
      video_production_image -an -y -init_hw_device vaapi=va -init_hw_device qsv=qs@va -init_hw_device opencl=ocl@va -hwaccel qsv \
      -port 0000:b1:01.2 -local_addr $5 -rx_addr $6 -fps $3 -pix_fmt yuv422p10le \
      -video_size "$1"x"$2" -udp_port 20000 -payload_type 112 -f mtl_st20p -i "0" \
      -vf "raisr=asm=opencl:threadcount=26:passes=2:bits=10:filterfolder=filters_2x/filters_highres" \
      -map 0:0 -vframes 100 -c:v hevc_qsv -pixel_format yuv422p10le /videos/$OUTPUT_FILE_NAME_MP4
}

run_imtl_multiple_tx () {
   docker run \
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
   --ip=$5 \
   --expose=20000-20170 \
   --ipc=host -v /dev/shm:/dev/shm \
   --cpuset-cpus=$7 \
   -e MTL_PARAM_LCORES=$8 \
   -e MTL_PARAM_DATA_QUOTA=10356 \
      video_production_image -y -an -hide_banner -loglevel error -video_size $1"x"$2 -f rawvideo -pix_fmt $4 -i /videos/$INPUT_FILE_NAME \
      -filter_complex "[0:v]format=yuv422p10le,fps=25,split=4[in1][in2][in3][in4]" -total_sessions 4 \
      -map "[in1]" -pix_fmt yuv422p10le -udp_port 20000 -port 0000:b1:01.1 -local_addr $5 -dst_addr $6 -f kahawai_mux - \
      -map "[in2]" -pix_fmt yuv422p10le -udp_port 20001 -port 0000:b1:01.1 -local_addr $5 -dst_addr $6 -f kahawai_mux - \
      -map "[in3]" -pix_fmt yuv422p10le -udp_port 20002 -port 0000:b1:01.1 -local_addr $5 -dst_addr $6 -f kahawai_mux - \
      -map "[in4]" -pix_fmt yuv422p10le -udp_port 20003 -port 0000:b1:01.1 -local_addr $5 -dst_addr $6 -f kahawai_mux - 
}

run_imtl_multiple_rx () {
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
   --ip=$5 \
   --expose=20000-20170 \
   --ipc=host -v /dev/shm:/dev/shm \
   --cpuset-cpus=$7 \
   -e MTL_PARAM_LCORES=$8 \
   -e MTL_PARAM_DATA_QUOTA=10356 \
      video_production_image -y -an -hide_banner -loglevel error -thread_queue_size 32768 \
      -framerate $3 -pixel_format $4 -width $1 -height $2 -udp_port 20000 -port 0000:b1:01.2 -local_addr $5 -src_addr $6 -ext_frames_mode 1 -total_sessions 4 -f kahawai -i "0" \
      -framerate $3 -pixel_format $4 -width $1 -height $2 -udp_port 20001 -port 0000:b1:01.2 -local_addr $5 -src_addr $6 -ext_frames_mode 1 -total_sessions 4 -f kahawai -i "1" \
      -framerate $3 -pixel_format $4 -width $1 -height $2 -udp_port 20002 -port 0000:b1:01.2 -local_addr $5 -src_addr $6 -ext_frames_mode 1 -total_sessions 4 -f kahawai -i "2" \
      -framerate $3 -pixel_format $4 -width $1 -height $2 -udp_port 20003 -port 0000:b1:01.2 -local_addr $5 -src_addr $6 -ext_frames_mode 1 -total_sessions 4 -f kahawai -i "3" \
      -map 0:0 -f rawvideo -pixel_format $4 -width $1 -height $2 -vframes 25 /videos/output_1080p_y210le_1_frames_1.yuv -y \
      -map 1:0 -f rawvideo -pixel_format $4 -width $1 -height $2 -vframes 25 /videos/output_1080p_y210le_1_frames_2.yuv -y \
      -map 2:0 -f rawvideo -pixel_format $4 -width $1 -height $2 -vframes 25 /videos/output_1080p_y210le_1_frames_3.yuv -y \
      -map 3:0 -f rawvideo -pixel_format $4 -width $1 -height $2 -vframes 25 /videos/output_1080p_y210le_1_frames_4.yuv -y
}

imtl_latest_tx () {
  docker run \
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
   --ip=$5 \
   --expose=20000-20170 \
   --ipc=host -v /dev/shm:/dev/shm \
   --cpuset-cpus=$7 \
   -e MTL_PARAM_LCORES=$8 \
   -e MTL_PARAM_DATA_QUOTA=10356 \
      video_production_image -hide_banner -loglevel quiet -video_size "$1"x"$2" -f rawvideo -pix_fmt yuv422p10le -i /videos/$INPUT_FILE_NAME -filter:v fps=$3 \
      -total_sessions 1 -port 0000:b1:01.1 -local_addr $5 -tx_addr $6 -udp_port 20000 -payload_type 112 -f mtl_st20p -
}

imtl_latest_rx () {
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
   --ip=$5 \
   --expose=20000-20170 \
   --ipc=host -v /dev/shm:/dev/shm \
   --cpuset-cpus=$7 \
   -e MTL_PARAM_LCORES=$8 \
   -e MTL_PARAM_DATA_QUOTA=10356 \
      video_production_image -y \
      -qsv_device /dev/dri/renderD128 -hwaccel qsv \
      -port 0000:b1:01.2 -local_addr $5 -rx_addr $6 -fps $3 -pix_fmt yuv422p10le \
      -video_size "$1"x"$2" -udp_port 20000 -payload_type 112 -f mtl_st20p -i "0" \
      -map 0:0 -c:v h264_qsv /videos/$OUTPUT_FILE_NAME_MP4 #-f rawvideo /videos/$OUTPUT_FILE_NAME -y
}

if [[ -z $1 ]]; then
   echo -e "\033[0;31mERROR: No argument provided\033[0m"
   help_message_arg_1
elif [ $1 == "qsv" ]; then
   ####### TEST QSV | y210le | 1920x1080 #######
   prepare_reference_file "1920x1080" "25" "y210le" "rawvideo" $INPUT_FILE_NAME
   run_qsv "1920x1080" "25" "y210le"
   equality_test "QSV" "1920x1080" "y210le" $OUTPUT_FILE_NAME $REF_FILE_NAME "rawvideo"

   ####### TEST QSV | y210le | 3840x2160 #######
   prepare_reference_file "3840x2160" "25" "y210le" "rawvideo" $INPUT_FILE_NAME
   run_qsv "3840x2160" "25" "y210le"
   equality_test "QSV" "3840x2160" "y210le" $OUTPUT_FILE_NAME $REF_FILE_NAME "rawvideo"

elif [ $1 == "imtl" ]; then
   ###### TEST IMTL | y210le | 1920x1080 #######
   prepare_reference_file "1920x1080" "25" "y210le" "rawvideo" $INPUT_FILE_NAME
   imtl_latest_tx "1920" "1080" "25" "y210le" "192.168.2.1" "192.168.2.2" "32-63" "59-63" &
   sleep 15
   imtl_latest_rx "1920" "1080" "25" "y210le" "192.168.2.2" "192.168.2.1" "96-127" "123-127"
   equality_test "IMTL" "1920x1080" "y210le" $OUTPUT_FILE_NAME $INPUT_FILE_NAME "rawvideo"

   ###### TEST IMTL | y210le | 3840x2160 #######
   prepare_reference_file "3840x2160" "25" "y210le" "rawvideo" $INPUT_FILE_NAME
   imtl_latest_tx "3840" "2160" "25" "y210le" "192.168.2.1" "192.168.2.2" "32-63" "59-63" &
   sleep 15
   imtl_latest_rx "3840" "2160" "25" "y210le" "192.168.2.2" "192.168.2.1" "96-127" "123-127"
   equality_test "IMTL" "3840x2160" "y210le" $OUTPUT_FILE_NAME $INPUT_FILE_NAME "rawvideo"

   #TODO: Add tests for multiple streams 

   ####### TEST IMTL MULTIPLE STREAMS | y210le | 1920x1080 #######
      # prepare_reference_file "1920x1080" "25" "y210le" "rawvideo" $INPUT_FILE_NAME
      # run_imtl_multiple_tx "1920" "1080" "25" "y210le" "192.168.2.1" "192.168.2.2" "32-63" "59-63" &
      # run_imtl_multiple_rx "1920" "1080" "25" "y210le" "192.168.2.2" "192.168.2.1" "96-127" "123-127"
      # equality_test "IMTL" "1920x1080" "y210le"

   ####### TEST LATEST IMTL | y210le | 1920x1080 #######
      # prepare_reference_file "1920x1080" "25" "y210le" "rawvideo" $INPUT_FILE_NAME
      # imtl_latest_tx "1920" "1080" "25" "y210le" "192.168.2.1" "192.168.2.2" "32-63" "59-63" &
      # imtl_latest_rx "1920" "1080" "25" "y210le" "192.168.2.2" "192.168.2.1" "96-127" "123-127"
      # equality_test "IMTL" "1920x1080" "y210le"

   # prepare_reference_file "1920x1080" "25" "y210le" $INPUT_FILE_NAME

elif [ $1 == "vsr" ]; then
   prepare_reference_file "3840x2160" "25" "yuv420p" "h264_qsv" $REF_FILE_NAME_MP4
   prepare_input_for_vsr "h264_qsv" "yuv420p" $REF_FILE_NAME_MP4 $INPUT_FILE_NAME_MP4
   run_vsr "cpu"
   equality_test "VSR CPU" "3840x2160" "yuv420p" $OUTPUT_FILE_NAME_MP4 $REF_FILE_NAME_MP4

   prepare_reference_file "3840x2160" "25" "yuv420p" "h264_qsv" $REF_FILE_NAME_MP4
   prepare_input_for_vsr "h264_qsv" "yuv420p" $REF_FILE_NAME_MP4 $INPUT_FILE_NAME_MP4
   run_vsr "qsv"
   equality_test "VSR GPU (QSV)" "3840x2160" "yuv420p" $OUTPUT_FILE_NAME_MP4 $REF_FILE_NAME_MP4

   prepare_reference_file "3840x2160" "25" "yuv420p" "rawvideo" $REF_FILE_NAME
   prepare_input_for_vsr_rawvideo "rawvideo" "yuv420p" "3840x2160" $REF_FILE_NAME $INPUT_FILE_NAME
   run_vsr_raw "cpu" "yuv420p" "1920x1080" "3840x2160"
   equality_test "VSR rawvideo CPU" "3840x2160" "yuv420p" $OUTPUT_FILE_NAME $REF_FILE_NAME "rawvideo"

elif [ $1 == "imtl_rx" ]; then
   # docker run -it \
   #    --user root\
   #    --privileged \
   #    --device=/dev/vfio:/dev/vfio \
   #    --device=/dev/dri:/dev/dri \
   #    --cap-add ALL \
   #    -v $(pwd):/videos \
   #    -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
   #    -v /tmp/kahawai_lcore.lock:/tmp/kahawai_lcore.lock \
   #    -v /dev/null:/dev/null \
   #    -v /tmp/hugepages:/tmp/hugepages \
   #    -v /hugepages:/hugepages \
   #    --network=my_net_801f0 \
   #    --ip=192.168.2.2 \
   #    --expose=20000-20170 \
   #    --ipc=host -v /dev/shm:/dev/shm \
   #    --cpuset-cpus="96-127" \
   #    -e MTL_PARAM_LCORES="123-127" \
   #    -e MTL_PARAM_DATA_QUOTA=10356 \
   #       video_production_image -loglevel verbose -qsv_device /dev/dri/renderD128 -an -y -hwaccel qsv \
   #       -port 0000:b1:01.2 -local_addr 192.168.2.2 -rx_addr 192.168.2.1 -pix_fmt yuv422p10le \
   #       -video_size 1920x1080 -udp_port 20000 -payload_type 112 -f mtl_st20p -i "k" -vf "hwmap=derive_device=opencl,format=opencl,raisr_opencl,hwmap=derive_device=qsv:reverse=1:extra_hw_frames=16" \
   #       -f rawvideo -pix_fmt y210le -video_size 3840x2160 /videos/$OUTPUT_FILE_NAME

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
      --cpuset-cpus="96-127" \
      -e MTL_PARAM_LCORES="123-127" \
      -e MTL_PARAM_DATA_QUOTA=10356 \
         video_production_image -loglevel verbose -an -y -init_hw_device vaapi=va -init_hw_device qsv=qs@va -init_hw_device opencl=ocl@va -hwaccel qsv \
         -port 0000:b1:01.2 -local_addr 192.168.2.2 -rx_addr 192.168.2.1 -pix_fmt yuv422p10le \
         -video_size 1920x1080 -udp_port 20000 -payload_type 112 -f mtl_st20p -i "k" -vf "raisr=asm=opencl:threadcount=26:passes=2:bits=10:filterfolder=filters_2x/filters_highres" \
         -f rawvideo -pix_fmt y210le -video_size 3840x2160 /videos/$OUTPUT_FILE_NAME
elif [ $1 == "jxs" ]; then
   prepare_reference_file "1920x1080" "25" "y210le" "hevc_qsv" $INPUT_FILE_NAME_MP4
   run_jxs "hevc_qsv"
   equality_test "JPEG-XS" "1920x1080" "y210le" $OUTPUT_FILE_NAME_MP4 $INPUT_FILE_NAME_MP4

   prepare_reference_file "3840x2160" "25" "y210le" "hevc_qsv" $INPUT_FILE_NAME_MP4
   run_jxs "hevc_qsv"
   equality_test "JPEG-XS" "3840x2160" "y210le" $OUTPUT_FILE_NAME_MP4 $INPUT_FILE_NAME_MP4

elif [ $1 == "imtl_qsv" ]; then
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
      video_production_image 
      # -video_size 1920x1080 -f rawvideo -pix_fmt yuv422p10le -i /videos/input_1080p_y210le_500_frames.yuv \
      # -filter:v fps=59.94 -total_sessions 1 -port 0000:b1:01.1 -local_addr "192.168.2.1" -tx_addr "192.168.2.2" -udp_port 20000 -payload_type 112 -f mtl_st20p -
      #video_production_image -loglevel verbose -i /videos/input_1080p_y210le_500_frames.yuv -filter:v fps=25 -udp_port 20000 -port 0000:b1:01.1 -local_addr 192.168.2.1 -dst_addr 192.168.2.2 -f kahawai_mux -
   
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
   --cpuset-cpus="96-127" \
   -e MTL_PARAM_LCORES="123-127" \
   -e MTL_PARAM_DATA_QUOTA=10356 \
      video_production_image #\
      # -loglevel verbose \
      # -framerate 25 \
      # -pixel_format y210le \
      # -width 1920 -height 1080 \
      # -udp_port 20000 -port 0000:b1:01.2 -local_addr 192.168.2.2 -src_addr 192.168.2.1 \
      # -dma_dev 0000:00:01.0 \
      # -ext_frames_mode 1 \
      # -f kahawai -i "0" -vframes 150 -map 0:0 -c:v libx265 -an -x265-params crf=25 /videos/output_1080p_y210le_500_frames.mkv -y
elif [ $1 == "imtl_jxs" ]; then
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
      video_production_image 
      # -video_size 1920x1080 -f rawvideo -pix_fmt yuv422p10le -i /videos/input_1080p_y210le_500_frames.yuv \
      # -filter:v fps=59.94 -total_sessions 1 -port 0000:b1:01.1 -local_addr "192.168.2.1" -tx_addr "192.168.2.2" -udp_port 20000 -payload_type 112 -f mtl_st20p -
      #video_production_image -loglevel verbose -i /videos/input_1080p_y210le_500_frames.yuv -filter:v fps=25 -udp_port 20000 -port 0000:b1:01.1 -local_addr 192.168.2.1 -dst_addr 192.168.2.2 -f kahawai_mux -

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
   --cpuset-cpus="96-127" \
   -e MTL_PARAM_LCORES="123-127" \
   -e MTL_PARAM_DATA_QUOTA=10356 \
      video_production_image #\
      # -loglevel verbose \
      # -framerate 25 \
      # -pixel_format y210le \
      # -width 1920 -height 1080 \
      # -udp_port 20000 -port 0000:b1:01.2 -local_addr 192.168.2.2 -src_addr 192.168.2.1 \
      # -dma_dev 0000:00:01.0 \
      # -ext_frames_mode 1 \
      # -f kahawai -i "0" -vframes 150 -map 0:0 -c:v libx265 -an -x265-params crf=25 /videos/output_1080p_y210le_500_frames.mkv -y
elif [ $1 == "imtl_vsr" ]; then
      echo "STUB"
elif [ $1 == "imtl_qsv_vsr" ]; then
   # #prepare_reference_file "3840x2160" "25" "y210le" "hevc_qsv" $REF_FILE_NAME_MP4
   # prepare_reference_file "1920x1080" "25" "y210le"
   # #prepare_input_for_vsr
   # imtl_latest_tx "1920" "1080" "25" "y210le" "192.168.2.1" "192.168.2.2" "32-63" "59-63" &
   # run_imtl_rx_qsv_vsr "1920" "1080" "25" "y210le" "192.168.2.2" "192.168.2.1" "96-127" "123-127"
   # #equality_test "IMTL SR GPU (QSV)" "3840x2160" "y210le" $OUTPUT_FILE_NAME_MP4 $REF_FILE_NAME_MP4
   # #equality_test "IMTL QSV VSR" "1920x1080" "y210le" $OUTPUT_FILE_NAME $REF_FILE_NAME "rawvideo"

   # rm -rf *mp4

   # echo "### STEP 1 ###"

   # docker run -it \
   #    --user root\
   #    --privileged \
   #    --device=/dev/dri:/dev/dri \
   #    -v $(pwd):/videos \
   #    -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
   #    video_production_image -loglevel verbose -qsv_device /dev/dri/renderD128 -an -y -hwaccel qsv -f lavfi -i testsrc=d=10:s=3840x2160:r=25,format=y210le -c:v h264_qsv /videos/$REF_FILE_NAME_MP4

   # echo "### STEP 2 ###"

   # docker run -it \
   #    --user root\
   #    --privileged \
   #    --device=/dev/dri:/dev/dri \
   #    -v $(pwd):/videos \
   #    -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
   #    video_production_image -y -an \
   #    -qsv_device /dev/dri/renderD128 \
   #    -hwaccel qsv -c:v h264_qsv -i /videos/$REF_FILE_NAME_MP4 -vf "scale_qsv=w=iw/2:h=ih/2,hwdownload" -c:v h264_qsv /videos/$INPUT_FILE_NAME_MP4

   # echo "### STEP 3 ###"

   # docker run \
   #    --user root\
   #    --privileged \
   #    --device=/dev/vfio:/dev/vfio \
   #    --device=/dev/dri:/dev/dri \
   #    --cap-add ALL \
   #    -v $(pwd):/videos \
   #    -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
   #    -v /tmp/kahawai_lcore.lock:/tmp/kahawai_lcore.lock \
   #    -v /dev/null:/dev/null \
   #    -v /tmp/hugepages:/tmp/hugepages \
   #    -v /hugepages:/hugepages \
   #    --network=my_net_801f0 \
   #    --ip=192.168.2.1 \
   #    --expose=20000-20170 \
   #    --ipc=host -v /dev/shm:/dev/shm \
   #    --cpuset-cpus="32-63" \
   #    -e MTL_PARAM_LCORES="59-63" \
   #    -e MTL_PARAM_DATA_QUOTA=10356 \
   #       video_production_image -loglevel verbose -qsv_device /dev/dri/renderD128 -an -y -hwaccel qsv -stream_loop 1 -video_size "1920x1080" -f rawvideo -pix_fmt yuv422p10le -i /videos/$INPUT_FILE_NAME -filter:v fps=25 \
   #       -port 0000:b1:01.1 -local_addr 192.168.2.1 -tx_addr 192.168.2.2 -udp_port 20000 -payload_type 112 -f mtl_st20p - &

   echo "### STEP 4 ###"

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
      --cpuset-cpus="96-127" \
      -e MTL_PARAM_LCORES="123-127" \
      -e MTL_PARAM_DATA_QUOTA=10356 \
         video_production_image -loglevel verbose -an -y -total_sessions 1 -port 0000:b1:01.2 -local_addr 192.168.2.2 -rx_addr 192.168.2.1 -fps 59.94 -pix_fmt yuv422p10le \
         -video_size 1920x1080 -udp_port 20000 -payload_type 112 -f mtl_st20p -i "k" -vframes 2000 -f rawvideo /dev/null -y
         
         # -hwaccel_output_format qsv \
         # -port 0000:b1:01.2 -local_addr 192.168.2.2 -rx_addr 192.168.2.1 -fps 25 -pix_fmt yuv422p10le -video_size 1920x1080 -udp_port 20000 -payload_type 112 -f mtl_st20p -i "0" \
         # -map 0:0 -vframes 100 -pix_fmt y210le -video_size 3840x2160 /videos/$OUTPUT_FILE_NAME

   #  -filter_complex "[0:v]format=$3,fps=$2,split=4[in1][in2][in3][in4]; \
   #          [in1]hwupload=extra_hw_frames=1,scale_qsv=w=iw/2:h=ih/2:mode=compute:async_depth=1[tile0];\
   #          [in2]hwupload=extra_hw_frames=1,scale_qsv=w=iw/2:h=ih/2:mode=compute:async_depth=1[tile1];\
   #          [in3]hwupload=extra_hw_frames=1,scale_qsv=w=iw/2:h=ih/2:mode=compute:async_depth=1[tile2];\
   #          [in4]hwupload=extra_hw_frames=1,scale_qsv=w=iw/2:h=ih/2:mode=compute:async_depth=1[tile3];\
   #          [tile0][tile1][tile2][tile3]xstack_qsv=inputs=4:layout=0_0|0_h0|w0_0|w0_h0[multiview];[multiview]hwdownload,format=y210le[out_multiview]" \
   #       -map "[out_multiview]" -f rawvideo -pix_fmt $3 /videos/$OUTPUT_FILE_NAME

elif [ $1 == "imtl_qsv_jxs" ]; then
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
      video_production_image 
      # -video_size 1920x1080 -f rawvideo -pix_fmt yuv422p10le -i /videos/input_1080p_y210le_500_frames.yuv \
      # -filter:v fps=59.94 -total_sessions 1 -port 0000:b1:01.1 -local_addr "192.168.2.1" -tx_addr "192.168.2.2" -udp_port 20000 -payload_type 112 -f mtl_st20p -
      #video_production_image -loglevel verbose -i /videos/input_1080p_y210le_500_frames.yuv -filter:v fps=25 -udp_port 20000 -port 0000:b1:01.1 -local_addr 192.168.2.1 -dst_addr 192.168.2.2 -f kahawai_mux -

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
   --cpuset-cpus="96-127" \
   -e MTL_PARAM_LCORES="123-127" \
   -e MTL_PARAM_DATA_QUOTA=10356 \
      video_production_image #\
      # -loglevel verbose \
      # -framerate 25 \
      # -pixel_format y210le \
      # -width 1920 -height 1080 \
      # -udp_port 20000 -port 0000:b1:01.2 -local_addr 192.168.2.2 -src_addr 192.168.2.1 \
      # -dma_dev 0000:00:01.0 \
      # -ext_frames_mode 1 \
      # -f kahawai -i "0" -vframes 150 -map 0:0 -c:v libx265 -an -x265-params crf=25 /videos/output_1080p_y210le_500_frames.mkv -y
else
   help_message_arg_1
fi