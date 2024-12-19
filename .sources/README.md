# Intel® Tiber™ Broadcast Suite Source Browser Docker

## Overview

This Docker image is created solely for the purpose of browsing the source code of the Intel® Tiber™ Broadcast Suite. It contains all the source files but does not provide any functionality or executable binaries.

## ⚠️ Warning

> **This Docker image has no functionality and only contains the source files.**  
> It is not intended for running or testing the Intel® Tiber™ Broadcast Suite.  
> If you are looking to build and run the suite, please refer to the [README](../docs/README.md).

## Prerequisites
- Docker installed on your machine.

## Build Instructions

To build the Docker image, use the following command:

```bash
# access the root of the folder
cd Intel-Tiber-Broadcast-Suite

# build the sources docker
ENV_PROXY_ARGS=()
while IFS='' read -r line; do
    ENV_PROXY_ARGS+=("--build-arg")
    ENV_PROXY_ARGS+=("${line}=${!line}")
done < <(compgen -e | grep -E "_(proxy|PROXY)")
docker build "${ENV_PROXY_ARGS[@]}" "$@" -f .sources/Dockerfile.sources -t 2024.1.0-sources .
```

## Usage

To use this Docker image:

```bash
docker run -it 2024.1.0-sources
```
