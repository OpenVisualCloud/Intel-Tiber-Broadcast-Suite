# Build guide

Building of the Intel® Tiber™ Broadcast Suite is based on docker container building.

## 1. Prerequisites

Steps to perform before run Intel® Tiber™ Broadcast Suite on host with Ubuntu operating system installed.

### 1.1. BIOS settings
> **Note:** It is recommended to properly setup BIOS settings before proceeding. Depending on manufacturer, labels may vary. Please consult an instruction manual or ask a platform vendor for detailed steps.

Following technologies must be enabled for Media Transport Library (MTL) to function properly:
- [Intel® Virtualization for Directed I/O (VT-d)](https://en.wikipedia.org/wiki/X86_virtualization#Intel_virtualization_(VT-x))
- [Single-root input/output virtualization (SR-IOV)](https://en.wikipedia.org/wiki/Single-root_input/output_virtualization)
- For 200 GbE throughput on [Intel® Ethernet Network Adapter E810-2CQDA2 card](https://ark.intel.com/content/www/us/en/ark/products/210969/intel-ethernet-network-adapter-e810-2cqda2.html) a PCI-E lane bifurcation is required.

### 1.2. Install Docker

> **Note:**  This step is optional if you want to install Intel® Tiber™ Broadcast Suite locally.

#### 1.2.1. Install Docker build environment

To install Docker environment please refer to the official Docker Engine on Ubuntu installation manual's [Install using the apt repository](https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository) section.

> **Note:** Do not skip `docker-buildx-plugin` installation, otherwise the `build.sh` script may not run properly.

#### 1.2.2. Setup docker proxy

Depending on the network environment it could be required to set up the proxy. In that case please refer to [Configure the Docker client](https://docs.docker.com/network/proxy/#configure-the-docker-client) section of _Configure Docker to use a proxy server_ guide.

### 1.3 Install GPU driver
#### 1.3.1 Intel Flex GPU driver

To install Flex GPU driver follow the [1.4.3. Ubuntu Install Steps](https://dgpu-docs.intel.com/driver/installation.html#ubuntu-install-steps) part of the Installation guide for Intel® Data Center GPUs.

> **Note:** If prompted with `Unable to locate package`, please ensure repository key `intel-graphics.key` is properly dearmored and installed as `/usr/share/keyrings/intel-graphics.gpg`.

Use vainfo command to check the gpu installation
```shell
vainfo
```


#### 1.3.2 Nvidia GPU driver

In case of using an Nvidia GPU, please follow the steps below:
```
sudo apt install --install-suggests nvidia-driver-550-server
sudo apt install nvidia-utils-550-server
```

In case of any issues please follow [Nvidia GPU driver install steps](https://ubuntu.com/server/docs/nvidia-drivers-installation#heading--manual-driver-installation-using-apt)

> **Note:** Supported version of Nvidia driver compatible with packages inside Docker container is
>* **Driver Version: 550.90.07**
>* **CUDA Version: 12.4**

### 1.4. Install and configure host's NIC drivers and related software

1. If you didn't do it already, then download the project from GitHub repo.
    ```
    git clone https://github.com/OpenVisualCloud/Intel-Tiber-Broadcast-Suite
    cd Intel-Tiber-Broadcast-Suite
    ```

1. Install patched ice driver for Intel® E810 Series Ethernet Adapter NICs.

   1. Download ice driver.
       ```shell
       mkdir -p ${HOME}/ice_patched
       . versions.env && wget -qO- $LINK_ICE_DRIVER | tar -xz -C ${HOME}/ice_patched
       ```

   1. Patch the ice driver.
       ```shell
        # Ensure the target directory exists
        mkdir -p ${HOME}/Media-Transport-Library

       # Download Media Transport Library:
       . versions.env && curl -Lf https://github.com/OpenVisualCloud/Media-Transport-Library/archive/refs/tags/${MTL_VER}.tar.gz | tar -zx --strip-components=1 -C ${HOME}/Media-Transport-Library

       . versions.env && git -C ${HOME}/ice_patched/ice* apply ~/Media-Transport-Library/patches/ice_drv/${ICE_VER}/*.patch
       ```

   1. Install the ice driver.
       ```shell
       cd ${HOME}/ice_patched/ice-*/src
       make
       sudo make install
       sudo rmmod irdma 2>/dev/null
       sudo rmmod ice
       sudo modprobe ice
       cd -
       ```

   1. Check if the driver is installed properly, and if so clean up.
        ```shell
        # should give you output
        sudo dmesg | grep "Intel(R) Ethernet Connection E800 Series Linux Driver - version Kahawai"
        rm -rf ${HOME}/ice_patched ${HOME}/Media-Transport-Library
        ```

   1. Update firmware
        ```shell
        . versions.env && wget ${LINK_ICE_FIRMWARE}
        unzip Release_*.zip
        cd NVMUpdatePackage/E810
        tar xvf E810_NVMUpdatePackage_v*_Linux.tar.gz
        cd E810/Linux_x64/
        sudo ./nvmupdate64e
        ```

    1. Verify installation
        ```shell
        # replace with your device
        ethtool -i ens801f0
        ```
        Result should look like:
        ```
        driver: ice
        version: Kahawai_1.14.9_20240613
        firmware-version: 4.60 0x8001e8dc 1.3682.0
        ```

    > **Note:** if you encountered any problems please go to <https://github.com/OpenVisualCloud/Media-Transport-Library/blob/maint-24.09/doc/e810.md>.


1. Configure VFIO (IOMMU) required by PMD-based DPDK using [Run Guide](https://github.com/OpenVisualCloud/Media-Transport-Library/blob/maint-24.09/doc/run.md), chapter 1, and (optionally) 7 for PTP configuration.

### 1.5 Optional: Install MCM Proxy

> **Note:** This step is required for the [MCM Proxy Pipelines](../pipelines/mcm_media_proxy_tx.sh).

Please install the MCM Proxy
[MCM Dockerized](https://github.com/OpenVisualCloud/Media-Communications-Mesh/tree/main?tab=readme-ov-file#dockerfiles-build)

If you want to avoid using docker and want to run the mcm-proxy on bare metal
[MCM instalation](https://github.com/OpenVisualCloud/Media-Communications-Mesh/tree/main?tab=readme-ov-file#getting-started)

## 2. Install Intel Tiber™ Broadcast Suite

### Option #1: Build Docker image from Dockerfile using build.sh script
> **Note:** This method is recommended instead of Option 2 - layers are built in parallel, cross-compability is possible.

Access the project directory

```shell
cd Intel-Tiber-Broadcast-Suite
```

Run build.sh script

> **Note:** For `build.sh` script to run without errors, `docker-buildx-plugin` must be installed. The error thrown without the plugin does not inform about that fact, rather that the flags are not correct. See chapter [1.2. Install Docker build environment](#12-install-docker-build-environment) for installation details.

```shell
./build.sh
```

### Option #2: Local installation from Debian packages

You can install the Intel® Tiber™ Broadcast Suite localy on bare metal. This
installation allows you to skip installing docker altogether.

```shell
git clone https://github.com/OpenVisualCloud/Intel-Tiber-Broadcast-Suite
cd Intel-Tiber-Broadcast-Suite
./build.sh -l
```

### Option #3: Install Docker image from Docker Hub
Visit <https://hub.docker.com/r/intel/intel-tiber-broadcast-suite/> Intel® Tiber™ Broadcast Suite image docker hub to select the most appropriate version.

Pull the Intel® Tiber™ Broadcast Suite image from Docker Hub
```shell
docker pull intel/intel-tiber-broadcast-suite:latest
```

### Option #4: Build Docker image from Dockerfile manually

> **Note:** Below method does not require buildx, but lacks cross-compability and may prolongate the build process.

Download, Patch, Build, and Install DPDK from source code

   1. Download and Extract DPDK and MTL:
        ```bash
       . versions.env && curl -Lf https://github.com/OpenVisualCloud/Media-Transport-Library/archive/refs/tags/${MTL_VER}.tar.gz | tar -zx --strip-components=1 -C ${HOME}/Media-Transport-Library

        . versions.env && curl -Lf https://github.com/DPDK/dpdk/archive/refs/tags/v${DPDK_VER}.tar.gz | tar -zx --strip-components=1 -C dpdk
        ```

   1. Apply Patches from Media Transport Library:
        ```bash
        # Apply patches:
        . versions.env && cd dpdk && git apply ${HOME}/Media-Transport-Library/patches/dpdk/$DPDK_VER/*.patch
        ```


   1. Build and Install DPDK:
        ```bash
        # Prepare the build directory:
        meson build

        # Build DPDK:
        ninja -C build

        # Install DPDK:
        sudo ninja -C build install
        ```

   1. Clean Up:
        ```bash
        cd ..
        rm -drf dpdk
        ```

Build image using Dockerfile
```shell
docker build $(cat versions.env | xargs -I {} echo --build-arg {}) -t video_production_image -f Dockerfile .
```

Change number of cores used to build by make can be changed  by _--build-arg nproc={number of proc}_

```shell
docker build $(cat versions.env | xargs -I {} echo --build-arg {}) --build-arg nproc=1 -t video_production_image -f Dockerfile .
```

Build the mtl manager docker

```shell
cd ${HOME}/Media-Transport-Library/manager
docker build --build-arg VERSION=1.0.0.TIBER -t mtl-manager:latest .
cd -
```


## 3. Running Intel Tiber™ Broadcast Suite

### 3.1. First run script

> **Note:** first_run.sh needs to be run after every reset of the machine

From the root of the Intel® Tiber™ Broadcast Suite repository, execute `first_run.sh` script that sets up the hugepages, locks for MTL, E810 NIC's virtual controllers and runs MtlManager docker container:

```shell
sudo -E ./first_run.sh | tee virtual_functions.txt
```
> **Note:** Please ensure the command is executed with `-E` switch, to copy all the necessary environment variables. Lack of the switch may cause the script to fail silently.

When running the Intel Tiber™ Broadcast Suite locally, please execute first_run with the -l argument.
```shell
sudo -E ./first_run.sh -l | tee virtual_functions.txt
```
This script will start the Mtl Manager locally. To avoid issues with core assignment in Docker, ensure that the Mtl Manager is running. The Mtl Manager is typically run within a Docker container, but the `-l` argument allows it to be executed directly from the terminal.

> **Note:** Ensure that `MtlManager` is running when using the Intel Tiber™ Broadcast Suite locally. You can check this by running `pgrep -l "MtlManager"`. If it is not running, start it with the command `sudo MtlManager`.

> **Note:** In order to avoid unnecessary reruns, preserve the command's output as a file to note which interface was bound to which Virtual Functions.

### 3.2. Test docker installation

```shell
docker run --rm -it --user=root --privileged video_production_image --help
```

### 3.3. Test local installation
```shell
ffmpeg --help
```


## Go to the run.md instruction for more details on how to run the image
#### [Running Intel® Tiber™ Broadcast Suite Pipelines](./run.md)
