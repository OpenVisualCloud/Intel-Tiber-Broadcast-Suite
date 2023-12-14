#!/bin/bash
pkill -f -9 ffmpeg
rm -f gradients*.* outfile*.*
./generator.sh
cp gradients.yuv gradients1.yuv 
cp gradients.yuv gradients2.yuv 
cp gradients.yuv gradients3.yuv 
./4sources_CPU.sh
./4sources_GPU.sh
./compressor.sh outfile_CPU
./compressor.sh outfile_GPU
sha256sum -b *.mp4
if [[ `sha256sum -b outfile_CPU.mp4 | cut -d ' ' -f 1` == `sha256sum -b outfile_GPU.mp4 | cut -d ' ' -f 1` ]]; then
 rm -f gradients*.* outfile*.*
 echo "TEST SUCCEEDED"
else
 rm -f gradients*.* outfile*.*
 echo "TEST FAILED"
fi


