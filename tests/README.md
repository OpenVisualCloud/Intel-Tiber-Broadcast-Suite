# Video Production Pipeline Project Tests

Status from check on an SPR machine with Flex 170 GPU on 2024-05-10.

## Separate test scripts

These scripts are created for validation purposes.

| Works? | Script filename | Description | Output | Comment |
|:-------|:----------------|:------------|:-------|:--------|
| ❌NO | [16ins_4k.sh](./16ins_4k.sh) | Assuming 16x source 4K video files exist, it generates multiviewer output 4K hevc video file using qsv. | out_4k.mkv | `Decode error rate 1 exceeds maximum 0.666667` is thrown regardless of input. See `16ins_4k.log`. |
| ✔️YES | [16ins_4k_y210le.sh](./16ins_4k_y210le.sh) | Same as previous one but uses y210le video format. | out_y210le.yuv | In version adjusted for 4096x1744 resolution Sintel videos, see `16ins_4k_y210le_4096x1744.log` and details below.* |
| ❌NO | [2sources_IMTL.sh](./2sources_IMTL.sh) | The script tests IMTL with 2 sources setup. | only ffmpeg output | `Unrecognized option 'src_addr'` is thrown, more errors when removed |
| ❌NO | [blend_test.sh](./blend_test.sh)<br />>blend_digit_generator.sh<br />>blend_digit_generator.sh<br />>blend_camera.sh<br />>blend_blender.sh<br />>blend_compressor.sh | The script tests blending IMTL source with existing file in y210le format. After that result file is compressed. | blend_output.mp4 | Fails on `./blend_camera.sh` with `Unrecognized option 'dst_addr'.` |
| ❌NO | [docker_receiver.sh](./docker_receiver.sh) | The script tests IMTL with 4 sources in y210le format. | only ffmpeg output | Fails on `Unrecognized option 'src_addr'.` |
| ❌NO | ~~[multiviewer_local.sh](./multiviewer_local.sh)~~<br />multiviewer_4x_local.sh<br />or multiviewer_9x_local.sh | The script tests multiviever pipeline in performace mode with 4 dummy IMTL sources (frames always available) and dummy output. | only ffmpeg output | Fails on `Unrecognized option 'src_addr'.` |
| ❌NO | [IMTL_test.sh](./IMTL_test.sh) | The script uses IMTL\*.sh scripts. Script generates gradient movie, send it and receives on the same network interface using IMTL, compress and compare input and output video using checksums. It will print message if everything succeeded or not. | TEST SUCCEEDED or FAILED | `Unrecognized option 'src_addr'.` and other errors appear |
| ❌NO | [multiviewer_test.sh](./multiviewer_test.sh) | The script uses multiviewer\*.sh scripts. It generates 4 identical gradient 4K videos in y210le format. Then it scales them to half of size and put them in 2x2 grid 4K video in the same format. It compares CPU video checksum with GPU video checksum. It will print message if everything succeeded or not. | TEST SUCCEEDED or FAILED | `Unrecognized option 'filter_complex_frames'.` when running `multiviewer_4sources_GPU.sh` |

NOTES: In case of failures please check if first_run.sh was executed after system reboot and if `--cpuset-cpus="..."` and `-e MTL_PARAM_LCORES="..."` contain valid values.\
It is recommended to run scripts as root or using sudo.

*25-frame sequence clips generated out of downloaded 4K TIF files from https://media.xiph.org/sintel/sintel-4k-tiff16/ by `wget https://media.xiph.org/sintel/sintel-4k-tiff16/00000${num}.tif` and `docker run --rm -v .../input_clips:/clips video_production_image -framerate 25 -pattern_type glob -i /clips/00000${num}*.tif -c:v rawvideo -pix_fmt y210le /clips/out${num}00.yuv`.

## Combined test_docker.sh tests
All outputs are `****** TEST PASSED ******` or `****** TEST FAILED ******`.

| Works? | Script switch/es | Description | Comment |
|:-------|:----------------|:------------|:--------|
| ❌NO | imtl | IMTL transmission only test | `Error: mt_dev_get_socket_id, failed to get port for 0000:b1:01.2` |
| ✔️YES | qsv | QSV filters only test | - |
| ✔️YES | jxs | JPEG-XS codec only test | - |
| ✔️YES | vsr | VSR codec only test | - |
| ❌NO | imtl_qsv | IMTL transmission and QSV filters test | Not useful `Use -h to get full help or, even better, run 'man ffmpeg'` error is thrown. |
| ❌NO | imtl_jxs | IMTL transmission and JPEG-XS coding test | Not useful `Use -h to get full help or, even better, run 'man ffmpeg'` error is thrown. |
| ❌NO | imtl_qsv_vsr | IMTL transmission and JPEG-XS coding test | `Error: mt_dev_get_socket_id, failed to get port for 0000:b1:01.2` |
| ❌NO | imtl_qsv_jxs | IMTL transmission, QSV filters and JPEG-XS coding test | Not useful `Use -h to get full help or, even better, run 'man ffmpeg'` error is thrown. |
