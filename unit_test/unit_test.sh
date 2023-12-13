#!/bin/bash
pkill -f -9 ffmpeg
rm -f gradients.* received.*
./generator.sh
./docker_camera.sh & timeout 40 ./docker_receiver.sh
./compressor.sh
./compressor_received.sh
sha256sum -b *.mp4
if [[ `sha256sum -b gradients.mp4 | cut -d ' ' -f 1` == `sha256sum -b received.mp4 | cut -d ' ' -f 1` ]]; then
 echo "TEST SUCCEEDED"
else
 echo "TEST FAILED"
fi


