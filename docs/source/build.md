# Build guide

Building of the Intel® Tiber™ Broadcast Suite is based on docker container building.

## 1. Prerequisites

Steps to perform before runIntel® Tiber™ Broadcast Suite on host

### 1.1 Install Docker build environment

To build Docker build environment please refer to the official manual: [Docker installation](https://docs.docker.com/engine/install/ubuntu/)

### 1.2 Setup proxy

Depending on the network environment it could be required to set up the proxy. In that case please refer to the official Docker proxy setup manual: [Docker proxy](https://docs.docker.com/network/proxy/)

### 1.3 Install Flex GPU driver

To install Flex GPU dirver follow the instruction: [Flex GPU driver install steps](https://dgpu-docs.intel.com/driver/installation.html#ubuntu-install-steps)

### 1.4 Configure network

1.  Install patched ice driver for Intel® E810 Series Ethernet Adapter NICs:
    [Intel® E810 Series Ethernet Adapter driver install steps](https://github.com/OpenVisualCloud/Media-Transport-Library/blob/main/docs/e810.md)
2.  Install Data Plain with Intel® Media Transport Library patches included:
    [Patched DPDK install steps](https://github.com/OpenVisualCloud/Media-Transport-Library/blob/main/docs/build.md)
3. Configure VFIO (IOMMU) required by PMD based DPDK:
    [Configuration of the VFIO (IOMMU)](https://github.com/OpenVisualCloud/Media-Transport-Library/blob/main/docs/run.md)

## 2. Build


Download the project from GitHub repo
```
git clone https://github.com/OpenVisualCloud/Intel-Tiber-Broadcast-Suite
```

Build image using dockerfile
```
docker build -t video_production_image -f Dockerfile .
```

Change number of cores used to build by make can be changed  by _--build-arg nproc={number of proc}_

```
docker build --build-arg nproc=1 -t video_production_image -f Dockerfile .
```
