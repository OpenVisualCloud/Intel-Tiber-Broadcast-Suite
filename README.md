# Video Production Pipeline Project

Before using IMTL in Docker please install required packages on host machine:
[Build_instruction](https://github.com/OpenVisualCloud/Media-Transport-Library/blob/main/doc/build.md)

To run ffmpeg in Docker please run command which creates docker image:

```
docker build -t my_ffmpeg .
```

Step 1. If IMTL plugin support is needed then please run commands on host as a root:

```
./first_run.sh
```

Steps 1 is required each time host is restarted and IMTL is needed.

Step 2. Run .sh script with ffmpeg parameters. Examples are in [test_scripts](./test_scripts) directory.


