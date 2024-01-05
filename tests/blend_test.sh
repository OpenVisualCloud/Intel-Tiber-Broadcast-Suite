#!/bin/bash
pkill -f -9 ffmpeg
rm -f digit1.yuv digit2.yuv blend_output.yuv blend_output.mp4
./blend_digit_generator.sh 1
./blend_digit_generator.sh 2
./blend_camera.sh & timeout 40 ./blend_blender.sh
./blend_compressor.sh blend_output
rm -f digit1.yuv digit2.yuv blend_output.yuv

