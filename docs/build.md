# Build guide

Building of the Intel® Tiber™ Broadcast Suite is based on docker container building.

## 1. Prerequisites

Steps to perform before run Intel® Tiber™ Broadcast Suite on host with Ubuntu operating system installed.

### 1.1 BIOS settings
> **Note:** It is recommended to properly setup BIOS settings before proceeding. Depending on manufacturer, labels may vary. Please consult an instruction manual or ask a platform vendor for detailed steps.

Following technologies must be enabled for Media Transport Library (MTL) to function properly:
- [Intel® Virtualization for Directed I/O (VT-d)](https://en.wikipedia.org/wiki/X86_virtualization#Intel_virtualization_(VT-x))
- [Single-root input/output virtualization (SR-IOV)](https://en.wikipedia.org/wiki/Single-root_input/output_virtualization)
- Bifurcation on PCI-E lanes of Intel® E810 Series Ethernet Adapter card may be required in some cases <!--TODO: Document which cases require bifurcation-->


### 1.2 Install Docker build environment

To install Docker environment please refer to the official Docker Engine on Ubuntu installation manual's [Install using the apt repository](https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository) section.

> **Note:** Do not skip `docker-buildx-plugin` installation, otherwise the `build.sh` script may not run properly.

### 1.3 Setup proxy

Depending on the network environment it could be required to set up the proxy. In that case please refer to [Configure the Docker client](https://docs.docker.com/network/proxy/#configure-the-docker-client) section of _Configure Docker to use a proxy server_ guide.

### 1.4 Install Flex GPU driver

To install Flex GPU driver follow the [1.4.3. Ubuntu Install Steps](https://dgpu-docs.intel.com/driver/installation.html#ubuntu-install-steps) part of the Installation guide for Intel® Data Center GPUs.

> **Note:** If prompted with `Unable to locate package`, please ensure repository key `intel-graphics.key` is properly dearmored and installed as `/usr/share/keyrings/intel-graphics.gpg`.

### 1.5 Install and configure host's NIC drivers and related software

1. Gather information about currently used Media Transport Library tag with:
    ```shell
    grep "MTL_VER=" Dockerfile | awk -F "=" '{print gensub(/ \\/,"","g",$NF)}'
    ```
2. Clone Media Transport Library repository and checkout to the tag detected in a previous step with
    ```shell
    git clone https://github.com/OpenVisualCloud/Media-Transport-Library.git
    cd Media-Transport-Library
    git checkout <tag>
    ```
3. While in `Media-Transport-Library` folder, set `mtl_source_code` variable with:
    ```shell
    export mtl_source_code=${PWD}
    ```
4. Install patched ice driver for Intel® E810 Series Ethernet Adapter NICs based on the [Intel® E810 Series Ethernet Adapter driver install steps](https://github.com/OpenVisualCloud/Media-Transport-Library/blob/2f1c2a3be417065a4dc9276e2d7344d768e95118/doc/e810.md) instruction.

    > **Note:** Please ensure Intel® Ethernet Adapter Complete Driver Pack is downloaded in a version specified in the instruction from a link containing the `MTL_VER` commit hash.

5.  Install Data Plane with Media Transport Library patches included using Ubuntu-related steps from [Build Guide](https://github.com/OpenVisualCloud/Media-Transport-Library/blob/2f1c2a3be417065a4dc9276e2d7344d768e95118/doc/build.md).
    > **Note:** PIP package manager for Python reads proxy settings from environment variables, thus it might be required to re-setup proxy before proceeding, if sudo is used.
6. Configure VFIO (IOMMU) required by PMD-based DPDK using [Run Guide](https://github.com/OpenVisualCloud/Media-Transport-Library/blob/2f1c2a3be417065a4dc9276e2d7344d768e95118/doc/run.md), chapters 1-4, and (optionally) 7 for PTP configuration.

7. From the root of the Intel® Tiber™ Broadcast Suite repository, execute `first_run.sh` script that sets up the hugepages, locks for MTL, and E810 NIC's virtual controllers:
    ```shell
    sudo -E ./first_run.sh | tee virtual_functions.txt
    ```
    > **Note:** Please ensure the command is executed with `-E` switch, to copy all the necessary environment variables. Lack of the switch may cause the script to fail silently.

    > **Note:** In order to avoid unnecessary reruns, preserve the command's output as a file to note which interface was bound to which Virtual Functions.

## 2. Build

### 2.1. Using build.sh script
> **Note:** This method recommended instead of 2.2 - layers are built in parallel, cross-compability is possible.

Download the project from GitHub repo

```shell
git clone https://github.com/OpenVisualCloud/Intel-Tiber-Broadcast-Suite
cd Intel-Tiber-Broadcast-Suite
```

Run build.sh script

> **Note:** For `build.sh` script to run without errors, `docker-buildx-plugin` must be installed. The error thrown without the plugin does not inform about that fact, rather that the flags are not correct. See chapter [1.2. Install Docker build environment](#12-install-docker-build-environment) for installation details.

```shell
./build.sh
```

## 2.2. Alternative manual build method

> **Note:** Below method does not require buildx, but lacks cross-compability and may prolongate the build process.

Build image using Dockerfile
```shell
docker build -t video_production_image -f Dockerfile .
```

Change number of cores used to build by make can be changed  by _--build-arg nproc={number of proc}_

```shell
docker build --build-arg nproc=1 -t video_production_image -f Dockerfile .
```

Build the mtl manager docker

```shell
cd "$mtl_source_code"/manager
docker build -t mtl-manager:latest .
cd -
```


## 3. Test run the image

```shell
docker run --rm -it --user=root --privileged video_production_image --help
```
