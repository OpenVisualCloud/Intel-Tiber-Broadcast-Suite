# Alternative manual build method

## 1. Alternative manual build method
> **Note:** This is not the recommended method of bare metal installation only use as a reference if you encounter erros during the build.sh execution

### 1.1. Installing dependencies
First, add the CUDA APT repository to your system.

1. Update Package Lists
    ```bash
    sudo apt-get update
    ```

1. Install `wget` if it is not already installed:
    ```bash
    sudo apt-get install -y wget
    ```

1. Download the CUDA Keyring Package
    ```bash
    . versions.env && wget ${LINK_CUDA_REPO}
    ```

1. Install the CUDA Keyring Package
    ```bash
    sudo dpkg -i cuda-keyring_*.deb
    ```

1. Clean up
    ```bash
    rm cuda-keyring_*.deb
    ```

1. Now install all necesery packages
    ```bash
    sudo apt-get install --no-install-recommends -y \
            libigdgmm-dev \
            libva-dev \
            intel-media-va-driver \
            libvpl-dev \
            ca-certificates \
            build-essential \
            libarchive-tools \
            libnuma-dev \
            libjson-c-dev \
            libpcap-dev \
            libgtest-dev \
            libsdl2-dev \
            libsdl2-ttf-dev \
            libssl-dev \
            libtool \
            libx11-dev \
            libx11-xcb-dev \
            libwayland-dev \
            libxcb-dri3-dev \
            libxcb-present-dev \
            libxext-dev \
            libxfixes-dev \
            libxml2-dev \
            libdrm-dev \
            librdmacm-dev \
            zlib1g-dev \
            libelf-dev \
            git \
            cmake \
            meson \
            m4 \
            clang \
            llvm \
            curl \
            g++ \
            nasm \
            autoconf \
            automake \
            pkg-config \
            diffutils \
            gcc \
            gcc-multilib \
            xxd \
            zip \
            python3-pyelftools \
            systemtap-sdt-dev \
            sudo \
            libbsd-dev \
            ocl-icd-opencl-dev \
            libcap2-bin \
            ubuntu-drivers-common \
            libc6-dev \
            cuda-toolkit-12-6 \
            libnvidia-compute-550 \
            libfdt-dev
    ```

### 1.2. Step-by-Step Instructions to Download and install VMAF

1. Create a Directory for VMAF:
    ```bash
    mkdir vmaf
    ```

1. Download and Extract VMAF:
    ```bash
    . versions.env && curl -Lf https://github.com/Netflix/vmaf/archive/refs/tags/v${LIBVMAF}.tar.gz | tar -zx --strip-components=1 -C vmaf
    ```

1. Build and Install VMAF:
    ```bash
    # Navigate to the VMAF Directory:
    cd vmaf

    # Prepare then build VMAF
    meson setup --prefix=/usr --libdir=/usr/lib/x86_64-linux-gnu --buildtype release build
    ninja -j$(nproc) -vC build
    sudo ninja -j$(nproc) -vC build install
    ```

1. Clean Up:
    ```bash
    cd -
    rm -drf vmaf
    ```

### 1.3. Step-by-Step Instructions to Download, Build, and Clean Up SVT-AV1

1. Create a Directory for SVT-AV1:
    ```bash
    mkdir svt-av1
    ```

1. Download and Extract SVT-AV1:
    ```bash
    . versions.env && curl -Lf https://gitlab.com/AOMediaCodec/SVT-AV1/-/archive/v${SVTAV1}/SVT-AV1-v${SVTAV1}.tar.gz | tar -zx --strip-components=1 -C svt-av1
    ```

1. Build and Install SVT-AV1:
    ```bash
    # Navigate to the SVT-AV1 Directory:
    cd svt-av1

    # Create a Build Directory:
    mkdir Build && cd Build

    # Prepare then build SVT-AV1
    cmake .. -G"Unix Makefiles" -DCMAKE_BUILD_TYPE=Release
    make
    sudo make install
    ```

1. Clean Up:
    ```bash
    cd ../..
    rm -drf svt-av1
    ```

### 1.4. Download and Extract Vulkan Headers
1. Download Vulkan Headers:
    ```bash
    . versions.env && curl -Lf https://github.com/KhronosGroup/Vulkan-Headers/archive/refs/tags/${VULKANSDK}.tar.gz | tar -zx --strip-components=1 -C vulkan-headers
    ```

1. Build and Install Vulkan Headers:
    ```bash
    # Navigate to the Vulkan Headers Directory:
    cd vulkan-headers

    # Prepare the build directory:
    cmake -S . -B build/

    # Install Vulkan Headers:
    cmake --install build --prefix /usr/local
    ```

1. Clean Up:
    ```bash
    cd ..
    rm -drf vulkan-headers
    ```

### 1.5. Download and Extract XDP-Tools and libbpf

1. Download XDP-Tools and libbpf:
    ```bash
    . versions.env && \
    curl -Lf https://github.com/xdp-project/xdp-tools/archive/${XDP_VER}.tar.gz | tar -zx --strip-components=1 -C xdp-tools && \
    curl -Lf https://github.com/libbpf/libbpf/archive/${BPF_VER}.tar.gz | tar -zx --strip-components=1 -C xdp-tools/lib/libbpf
    ```

1. Build and Install XDP-Tools and libbpf:
    ```bash
    # Navigate to the XDP-Tools Directory:
    cd xdp-tools

    # Configure the build:
    ./configure

    # Build and install XDP-Tools:
    make
    sudo make -j$(nproc) install

    # Install libbpf:
    sudo make -C lib/libbpf/src install
    ```

1. Clean Up:
    ```bash
    cd ..
    rm -drf xdp-tools
    ```

### 1.6. Media Transport Library
#### Option #1: Install Media Transport Library using a Debian Package

1. Depending on your distribution, assign the appropriate value:
    ```bash
    # for Ubuntu LTS 22.04
    . versions.env && export LINK_MTL_DEBIAN_ZIP=$LINK_MTL_DEBIAN_v2204_ZIP
    ```
    ```bash
    # for Ubuntu LTS 24.04
    . versions.env && export LINK_MTL_DEBIAN_ZIP=$LINK_MTL_DEBIAN_v2404_ZIP
    ```

1. Download and Install the Debian Package:
    ```bash
    wget -O mtl.zip $LINK_MTL_DEBIAN_ZIP
    mkdir mtl && unzip mtl.zip -d mtl
    sudo dpkg -i mtl/*.deb
    ```

1. Clean Up:
    ```bash
    rm -rf mtl mtl.zip
    ```

#### Option #2: Alternatively Download and Build Media Transport Library (MTL) from source

1. Download Media Transport Library:
    ```bash
    . versions.env && curl -Lf https://github.com/OpenVisualCloud/Media-Transport-Library/archive/refs/tags/${MTL_VER}.tar.gz | tar -zx --strip-components=1 -C Media-Transport-Library
    ```

1. Build patch Media Transport Library:
    ```bash
    # Navigate to the Media Transport Library Directory:
    cd Media-Transport-Library

    git apply ../patches/imtl/0001-cartwheel-imtl-y210le-support.patch

    # Run the build script:
    ./build.sh
    ```
### 1.7. JPEG-XS
#### Option #1: Installing JPEG-XS Using a Debian Package

1. Depending on your distribution, assign the appropriate value:
    ```bash
    # for Ubuntu LTS 22.04
    . versions.env && export LINK_JPEGXS_DEBIAN_ZIP=$LINK_JPEGXS_DEBIAN_v2204_ZIP
    ```
    ```bash
    # for Ubuntu LTS 24.04
    . versions.env && export LINK_JPEGXS_DEBIAN_ZIP=$LINK_JPEGXS_DEBIAN_v2404_ZIP
    ```

1. Download and Install the Debian Package:
    ```bash
    wget -O jpegxs.zip $LINK_JPEGXS_DEBIAN_ZIP
    mkdir jpegxs && unzip jpegxs.zip -d jpegxs
    sudo dpkg -i jpegxs/*.deb
    ```

1. Clean Up:
    ```bash
    rm -rf jpegxs jpegxs.zip
    ```

#### Option #2: Alternatively Download, Build, and Install JPEG-XS from source code

1. Download and Extract JPEG-XS:
    ```bash
    . versions.env && curl -Lf https://github.com/OpenVisualCloud/SVT-JPEG-XS/archive/${JPEG_XS_COMMIT_ID}.tar.gz | tar -zx --strip-components=1 -C jpegxs
    ```

1. Build and Install JPEG-XS:
    ```bash
    # Navigate to the JPEG-XS Build Directory:
    cd jpegxs/Build/linux

    # Run the build script:
    ./build.sh install

    # Copy the kahawai.json file:
    cp ../imtl-plugin/kahawai.json ./
    ```

1. Build the IMTL Plugin:
    ```bash
    # Navigate to the IMTL Plugin Directory:
    cd ../imtl-plugin

    # Run the build script:
    ./build.sh
    ```
### 1.8. Download, Install, and Set Up Intel IPP

1. Download Intel IPP:
    ```bash
    wget --progress=dot:giga https://registrationcenter-download.intel.com/akdlm/IRC_NAS/046b1402-c5b8-4753-9500-33ffb665123f/l_ipp_oneapi_p_2021.10.1.16_offline.sh
    ```

1. Install Intel IPP (if not skipped):
    ```bash
    # Make the installer executable:
    chmod +x l_ipp_oneapi_p_2021.10.1.16_offline.sh

    # Run the installer:
    ./l_ipp_oneapi_p_2021.10.1.16_offline.sh -a -s --eula accept
    ```

1. Set Up the IPP Environment:
    ```bash
    echo "source /opt/intel/oneapi/ipp/latest/env/vars.sh" | tee -a ~/.bash_profile
    ```

1. Clean Up:
    ```bash
    rm -f l_ipp_oneapi_p_2021.10.1.16_offline.sh
    ```

### 1.9. Download, Patch, Build, and Install Video Super Resolution (VSR)

1. Download and Extract VSR:
    ```bash
    . versions.env && curl -Lf https://github.com/OpenVisualCloud/Video-Super-Resolution-Library/archive/refs/tags/${VSR}.tar.gz | tar -zx --strip-components=1 -C vsr
    ```

1. Build and Install VSR:
    ```bash
    # Navigate to the VSR Directory:
    cd vsr

    # Remove 'clan' from the build script:
    sed -i 's/clan//g' build.sh

    # Source the IPP environment:
    . /opt/intel/oneapi/ipp/latest/env/vars.sh

    # Run the build script with the specified options:
    ./build.sh -DCMAKE_INSTALL_PREFIX=$(pwd)/install -DENABLE_RAISR_OPENCL=ON
    ```

1. Clean Up:
    ```bash
    cd ..
    rm -drf vsr
    ```

### 1.10. Media Communications Mesh

#### Option #1: Install Media Communications Mesh Using a Debian Package

1. Depending on your distribution, assign the appropriate value:
    ```bash
    # for Ubuntu LTS 22.04
    . versions.env && export LINK_MCM_DEBIAN_ZIP=$LINK_MCM_DEBIAN_v2204_ZIP
    ```
    ```bash
    # for Ubuntu LTS 24.04
    . versions.env && export LINK_MCM_DEBIAN_ZIP=$LINK_MCM_DEBIAN_v2404_ZIP
    ```

1. Download and Install the Debian Package:
    ```bash
    wget -O mcm.zip $LINK_MCM_DEBIAN_ZIP
    mkdir mcm && unzip mcm.zip -d mcm
    sudo dpkg -i mcm/*.deb
    ```

1. Clean Up:
    ```bash
    rm -rf mcm mcm.zip
    ```

#### Option #2: Alternatively Download, Build, and Install JPEG-XS froum source code
1. Download and Extract Media Communications Mesh:
    ```bash
    . versions.env && curl -Lf https://github.com/OpenVisualCloud/Media-Communications-Mesh/archive/refs/tags/${MCM_VER}.tar.gz | tar -zx --strip-components=1 -C mcm
    ```

2. Build and Install Media Communications Mesh:
    ```bash
    # Navigate to the Media Communications Mesh Directory:
    cd mcm

    # Prepare the build directory:
    cmake -S sdk -B sdk/out

    # Build Media Communications Mesh:
    cmake --build sdk/out

    # Install Media Communications Mesh:
    sudo cmake --install sdk/out
    ```


### 1.11. Download, Build, and Install FFmpeg NV-Codec-Headers

1. Download and Extract FFmpeg NV-Codec-Headers:
    ```bash
    . versions.env && curl -Lf https://github.com/FFmpeg/nv-codec-headers/archive/${FFNVCODED_VER}.tar.gz | tar -zx --strip-components=1 -C nv-codec-headers
    ```

1. Build and Install FFmpeg NV-Codec-Headers:
    ```bash
    # Navigate to the NV-Codec-Headers Directory:
    cd nv-codec-headers

    # Build NV-Codec-Headers:
    make

    # Install NV-Codec-Headers:
    sudo make install PREFIX=/usr
    ```

1. Clean Up:
    ```bash
    cd ..
    rm -drf nv-codec-headers
    ```

### 1.12. FFmpeg
#### Option #1: Install FFmpeg Using a Debian Package

1. Depending on your distribution, assign the appropriate value:
    ```bash
    # for Ubuntu LTS 22.04
    . versions.env && export LINK_FFMPEG_DEBIAN_ZIP=$LINK_FFMPEG_DEBIAN_v2204_ZIP
    ```
    ```bash
    # for Ubuntu LTS 24.04
    . versions.env && export LINK_FFMPEG_DEBIAN_ZIP=$LINK_FFMPEG_DEBIAN_v2404_ZIP
    ```

1. Download and Install the Debian Package:
    ```bash
    wget -O ffmpeg.zip $LINK_FFMPEG_DEBIAN_ZIP
    mkdir ffmpeg && unzip ffmpeg.zip -d ffmpeg
    sudo dpkg -i ffmpeg/*.deb
    ```

#### Option #2: Alternatively Download, Patch, Build, and Install FFmpeg

1. Download and Extract FFmpeg:
    ```bash
    . versions.env && curl -Lf https://github.com/ffmpeg/ffmpeg/archive/${FFMPEG_COMMIT_ID}.tar.gz | tar -zx --strip-components=1 -C ffmpeg
    ```

1. Download and Extract Cartwheel Patches:
    ```bash
    . versions.env && curl -Lf https://github.com/intel/cartwheel-ffmpeg/archive/${CARTWHEEL_COMMIT_ID}.tar.gz | tar -zx --strip-components=1 -C cartwheel
    ```

1. Apply Patches to FFmpeg:
    ```bash
    # Navigate to the FFmpeg Directory:
    cd ffmpeg

    # Apply Cartwheel patches:
    git apply ../cartwheel/patches/*.patch

    # Copy MTL plugin files:
    cp ../Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_*.c libavdevice/
    cp ../Media-Transport-Library/ecosystem/ffmpeg_plugin/mtl_*.h libavdevice/

    # Apply MTL patches:
    git apply ../Media-Transport-Library/ecosystem/ffmpeg_plugin/7.0/*.patch

    # Apply JPEG-XS patches:
    cp ../jpegxs/ffmpeg-plugin/libsvtjpegxs* libavcodec/
    git apply --whitespace=fix ../jpegxs/ffmpeg-plugin/7.0/*.patch

    # Apply additional FFmpeg patches:
    git apply ../patches/ffmpeg/001-hwupload_async.diff
    git apply ../patches/ffmpeg/002-qsv_aligned_malloc.diff
    git apply ../patches/ffmpeg/003-qsvvpp_async.diff
    git apply ../patches/ffmpeg/004-filtergraph_async.diff
    git apply ../patches/ffmpeg/005-ffmpeg_scheduler.diff

    # Apply Media Communications Mesh patches:
    git apply -v --whitespace=fix --ignore-space-change ../mcm/ffmpeg-plugin/7.0/*.patch
    cp -f ../mcm/ffmpeg-plugin/mcm_* libavdevice/
    ```

1. Build and Install FFmpeg:
    ```bash
    # Source the IPP environment:
    . /opt/intel/oneapi/ipp/latest/env/vars.sh

    # Configure FFmpeg:
    ./configure \
    --disable-debug \
    --disable-doc \
    --disable-shared \
    --enable-static \
    --enable-ffprobe \
    --enable-libsvtav1 \
    --enable-libvpl \
    --enable-libvmaf \
    --enable-version3 \
    --enable-libxml2 \
    --enable-mtl \
    --enable-opencl \
    --enable-stripping \
    --enable-vaapi \
    --enable-vulkan \
    --enable-libsvtjpegxs \
    --enable-libipp \
    --enable-mcm \
    --enable-pthreads \
    --enable-ffnvcodec \
    --enable-cuda \
    --enable-cuvid \
    --enable-nvenc \
    --enable-nvdec \
    --enable-cuda-llvm \
    --extra-cflags="-march=native -fopenmp -I../vsr/install/include/ -I/opt/intel/oneapi/ipp/latest/include/ipp/ -I/usr/local/cuda/include" \
    --extra-ldflags="-fopenmp -L../vsr/install/lib -L/usr/local/cuda/lib64 -L/usr/lib64 -L/usr/local/lib" \
    --extra-libs='-lraisr -lstdc++ -lippcore -lippvm -lipps -lippi -lpthread -lm -lz -lbsd -lrdmacm -lbpf -lxdp' \
    --enable-cross-compile

    # Build FFmpeg:
    make -j$(nproc)

    # Install FFmpeg:
    sudo make install
    sudo ldconfig
    ```

1. Clean Up:
    ```bash
    cd ..
    rm -drf ffmpeg cartwheel
    ```

## 2. Go to the run.md instruction for more details on how to run the image
**[Running Intel® Tiber™ Broadcast Suite Pipelines](./run.md)**