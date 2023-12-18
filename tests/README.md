# Video Production Pipeline Project Tests

These scripts are created for validation purposes.\
16ins_4k.sh - assuming that there exist 16 source 4K video files it generates multiviewer output 4K hevc video file using qsv.\
16ins_4k_y210le.sh - the same as previous one but uses y210le video format.\
2sources_IMTL.sh - it tests IMTL with 2 sources setup.\
docker_receiver.sh - it tests IMTL with 4 sources in y210le format.\
IMTL_test.sh - script uses IMTL\*.sh scripts. Script generates gradient movie, send it and receives on the same network interface using IMTL, compress and compare input and output video using checksums. It will print message if everything succeded or not.\
multiviewer_test.sh - script uses multiviewer\*.sh scripts. Script generates 4 identical gradient 4K videos in y210le format. Then it scales them to half of size and put them in 2x2 grid 4K video in the same format. It compares CPU video checksum with GPU video checksum. It will print message if everything succeded or not.\

NOTES: In case of failures please check if first_run.sh was executed after system reboot and if --cpuset-cpus="..." and -e MTL_PARAM_LCORES="..." contain valid values.\
It is recommended to run scripts as root or using sudo.