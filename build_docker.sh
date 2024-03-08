#!/bin/bash

rm -rf libraries.media.encoding.svt-jpeg-xs

echo "**** DOWNLOADING JPEG-XS ****"
git clone https://github.com/intel-innersource/libraries.media.encoding.svt-jpeg-xs.git
cd libraries.media.encoding.svt-jpeg-xs
git checkout release_v0.8_ww04_24_beta3
git switch -c release_v0.8_ww04_24_beta3
cd ..


docker build -t video_production_image -f Dockerfile .