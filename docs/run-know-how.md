# Docker command breakdown

Full Intel® Tiber™ Broadcast Suite pipeline command consists of `docker run [docker_parameters] video_production_image [broadcast_suite_parameters]`.
This document describes a docker-related part of command - all `[docker_parameters]` switches until the name of the image, called `video_production_image`.

Information about the `[broadcast_suite_parameters]` part of the command may be found, depending on the plugin, under:
- [FFmpeg Intel® Media Communications Mesh Muxer Parameters Table](plugins/media-communications-mesh.md)
- [Media Transport Library](plugins/media-transport-library.md)
- [FFmpeg Intel® JPEG XS Parameters Table](plugins/svt-jpeg-xs.md)
- [Raisr FFmpeg Filter Plugin Parameters Table](plugins/video-super-resolution.md)


## Command example

> **Note:** The example below is based on [../pipelines/jpeg_xs_rx.sh](pipelines/jpeg_xs_rx.sh). Some of the pipelines may require a different number of parameters in order to run.

```bash
docker run -it \
   --user root\
   --privileged \
   --device=/dev/vfio:/dev/vfio \
   --device=/dev/dri:/dev/dri \
   --cap-add ALL \
   -v "$(pwd)":/videos \
   -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
   -v /tmp/kahawai_lcore.lock:/tmp/kahawai_lcore.lock \
   -v /dev/null:/dev/null \
   -v /tmp/hugepages:/tmp/hugepages \
   -v /hugepages:/hugepages \
   --network=my_net_801f0 \
   --ip=192.168.2.2 \
   --expose=20000-20170 \
   --ipc=host -v /dev/shm:/dev/shm \
   --cpuset-cpus=20-40 \
   -e MTL_PARAM_LCORES=30-40 \
   -e MTL_PARAM_DATA_QUOTA=10356 \
   video_production_image  [broadcast_suite_parameters]
```

## Docker parameters breakdown

- `-it`: Runs the container in interactive mode with a TTY for interaction.
- `--user root`: Sets the user inside the container to `root`.
- `--privileged`: Grants the container full access to the host system.
- `--device=/dev/vfio:/dev/vfio`: Mounts the host's `/dev/vfio` directory inside the container.
- `--device=/dev/dri:/dev/dri`: Mounts the host's `/dev/dri` directory inside the container.
- `--cap-add ALL`: Gives the container all capabilities, similar to root access.
- `-v "$(pwd)":/videos`: Binds the current working directory on the host to `/videos` inside the container.
- `-v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/`: Mounts the host's DRI drivers into the container.
- `-v /tmp/kahawai_lcore.lock:/tmp/kahawai_lcore.lock`: Shares a lock file between the host and the container.
- `-v /dev/null:/dev/null`: Makes `/dev/null` available inside the container.
- `-v /tmp/hugepages:/tmp/hugepages`: Shares the hugepages directory for memory management optimizations.
- `-v /hugepages:/hugepages`: Shares another hugepages directory.
- `--network=my_net_801f0`: Connects the container to the `my_net_801f0` network.
- `--ip=192.168.2.2`: Assigns the IP address `192.168.2.2` to the container.
- `--expose=20000-20170`: Exposes a range of ports for the container.
- `--ipc=host -v /dev/shm:/dev/shm`: Shares the host's IPC namespace and mounts the shared memory directory.
- `--cpuset-cpus=20-40`: Limits the container to specific CPUs on the host machine.
- `-e MTL_PARAM_LCORES=30-40`: Sets the `MTL_PARAM_LCORES` environment variable inside the container.
- `-e MTL_PARAM_DATA_QUOTA=10356`: Sets the `MTL_PARAM_DATA_QUOTA` environment variable inside the container.
- `video_production_image`: Specifies the Docker image to be used for the container.
