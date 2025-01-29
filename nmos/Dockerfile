# syntax=docker/dockerfile:1

#
# SPDX-License-Identifier: BSD-3-Clause
# Copyright(©) 2024 Intel Corporation
# Intel® Tiber™ Broadcast Suite
#
# build stage

FROM ubuntu:jammy AS build-stage

ARG nproc

SHELL ["/bin/bash", "-ex", "-o", "pipefail", "-c"]

ENV \
  TZ="Europe/Warsaw" \
  DEBIAN_FRONTEND="noninteractive" \
  APT_KEY_DONT_WARN_ON_DANGEROUS_USAGE=DontWarn

RUN \
    apt-get update --fix-missing && \
    apt-get install -y --no-install-recommends \
    g++ \
    build-essential \
    autoconf \
    libtool \
    pkg-config \
    clang \
    libc++-dev \
    openssl \
    libssl-dev \
    git \
    wget \
    gnupg \
    curl \
    ca-certificates \
    nano \
    python3 \
    python3-pip \
    python3-setuptools \
    rdma-core \
    unzip

RUN \
    curl -sL https://deb.nodesource.com/setup_16.x | bash - && \
    apt-get update --fix-missing && \
    apt-get install -y --no-install-recommends \
    nodejs && \
    corepack enable && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean -y --no-install-recommends && \
    apt-get autoclean -y --no-install-recommends


# Avahi:    dbus avahi-daemon libavahi-compat-libdnssd-dev libnss-mdns AND NOT make \

## Install latest versions of CMake and Conan using pip3 package installer
RUN \
    python3 -m pip install --upgrade pip setuptools wheel cmake conan~=1.47

## Get Certificates and scripts from AMWA-TV/nmos-testing
RUN \
    cd /home && mkdir certs && git config --global http.sslVerify false && \
    git clone https://github.com/AMWA-TV/nmos-testing.git && \
    mv /home/nmos-testing/test_data/BCP00301/ca/* /home/certs && \
    rm -rf /home/nmos-testing

## Get source for Sony nmos-cpp
## Commit 0fb6b51 corresponds to Conan package nmos-cpp/cci.20221203
ENV NMOS_CPP_VERSION=f54971298c47a633969e9e9adac824b56fc08da7
RUN \
    cd /home/ && curl --output - -s -k https://codeload.github.com/sony/nmos-cpp/tar.gz/$NMOS_CPP_VERSION | tar zxvf - -C . && \
    mv ./nmos-cpp-${NMOS_CPP_VERSION} ./nmos-cpp


    #to remove and copy nmos-cpp ptach
# FOR NMOS-CPP-CLIENT-NODE in BCS

COPY \
    "patches/" "/home/patches"

RUN \
    patch -d /home/nmos-cpp -p1 -i /home/patches/nmos-cpp.patch


#COPY ./node_implementation.cpp /home/nmos-cpp/Development/nmos-cpp-node/node_implementation.cpp
#COPY ./main.cpp /home/nmos-cpp/Development/nmos-cpp-node/main.cpp

RUN \
    echo "**** DOWNLOAD AND INSTALL gRPC v1.58 ****" && \
    git clone --branch "v1.58.0" --recurse-submodules --depth 1 --shallow-submodules https://github.com/grpc/grpc /tmp/grpc-source && \
    mkdir -p "/tmp/grpc-source/cmake/build" && \
    cmake -S "/tmp/grpc-source" -B "/tmp/grpc-source/cmake/build" -DgRPC_BUILD_TESTS=OFF -DgRPC_INSTALL=ON && \
    make -C "/tmp/grpc-source/cmake/build" -j${nproc} && \
    make -C "/tmp/grpc-source/cmake/build" install && \
    rm -rf /tmp/grpc-source

# RUN curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v25.6/protoc-25.6-linux-x86_64.zip && \
#     unzip protoc-25.6-linux-x86_64.zip -d /home/.local

# ENV PATH="$PATH:/home/.local/bin"
# Install BCS gRPC client
#COPY ./gRPC/build/ /home/gRPC/build/
# COPY gRPC/ /home/gRPC
# RUN ls -l /home/gRPC 
# RUN cd /home/gRPC && \
#     ./compile.sh --unit-testing

## You should use either Avahi or Apple mDNS - DO NOT use both
##
## mDNSResponder 878.260.1 Build and install
RUN cd /home/ && curl --output - -s -k https://codeload.github.com/apple-oss-distributions/mDNSResponder/tar.gz/mDNSResponder-878.260.1 | tar zxvf - -C . && \
    mv ./mDNSResponder-mDNSResponder-878.260.1 ./mDNSResponder && \
    patch -d mDNSResponder/ -p1 <nmos-cpp/Development/third_party/mDNSResponder/unicast.patch && \
    patch -d mDNSResponder/ -p1 <nmos-cpp/Development/third_party/mDNSResponder/permit-over-long-service-types.patch && \
    patch -d mDNSResponder/ -p1 <nmos-cpp/Development/third_party/mDNSResponder/poll-rather-than-select.patch && \
    cd /home/mDNSResponder/mDNSPosix && make os=linux && make os=linux install

## Build Sony nmos-cpp from sources - with workaround to keep latest packages for nmos
## Remember to use the same veriosn of nmos and registry
## First download all packages automatically using conan package manager, than build project
RUN pip install --upgrade conan~=2.4 && \
    cd /home/nmos-cpp/Development/ && \
    conan profile detect && \
    conan install --requires=nmos-cpp/cci.20240223 --deployer=direct_deploy --build=missing

COPY gRPC/ /home/gRPC
RUN ls -la /home/gRPC 
RUN cd /home/gRPC && \
    cmake -B build && \
    make -j100 -C build

RUN \
    mkdir /home/nmos-cpp/Development/build && \
    mkdir /home/nmos-cpp/Development/nmos-cpp-node/build

RUN ls -la /home/gRPC/build/
RUN ls -la /home/gRPC/

    # Copy files with error checking
RUN \
    cp /home/gRPC/build/libFFmpeg_wrapper_client.a /home/nmos-cpp/Development/build && \
    cp /home/gRPC/build/libhw_grpc_proto.a /home/nmos-cpp/Development/build && \
    cp /home/gRPC/config_params.hpp /home/nmos-cpp/Development/nmos-cpp-node && \
    cp /home/gRPC/FFmpeg_wrapper_client.h /home/nmos-cpp/Development/nmos-cpp-node && \
    cp /home/gRPC/build/ffmpeg_cmd_wrap.pb.h /home/nmos-cpp/Development/nmos-cpp-node/build && \
    cp /home/gRPC/build/ffmpeg_cmd_wrap.grpc.pb.h /home/nmos-cpp/Development/nmos-cpp-node/build

ENV LIBRARY_PATH=/home/nmos-cpp/Development/build

RUN \
    cd /home/nmos-cpp/Development/build && \
    cmake \
    -DCMAKE_PROJECT_TOP_LEVEL_INCLUDES:STRING="third_party/cmake/conan_provider.cmake" \
    -G "Unix Makefiles" \
    -DCMAKE_BUILD_TYPE:STRING="MinSizeRel" \
    -DCMAKE_CONFIGURATION_TYPES:STRING="MinSizeRel" \
    -DCXXFLAGS:STRING="-Os" \
    -DNMOS_CPP_USE_AVAHI:BOOL="0" \
    /home/nmos-cpp/Development/build .. && \
    make -j${nproc}

    #if [ -n "$nproc" ]; then echo "Making multi-threaded with $nproc jobs"; make -j100; else echo "Making single-threaded"; make -j100; fi

## Generate Example Certificates and position into correct locations
RUN cd /home/certs && mkdir run-certs && ./generateCerts registration1 nmos.tv query1.nmos.tv && \
    cd /home/certs/certs && \
    cp ca.cert.pem /home/certs/run-certs/ca.cert.pem && \
    cd /home/certs/intermediate/certs && \
    mv ecdsa.registration1.nmos.tv.cert.chain.pem /home/certs/run-certs/ecdsa.cert.chain.pem && \
    mv rsa.registration1.nmos.tv.cert.chain.pem /home/certs/run-certs/rsa.cert.chain.pem && \
    cd /home/certs/intermediate/private && \
    mv ecdsa.registration1.nmos.tv.key.pem /home/certs/run-certs/ecdsa.key.pem && \
    mv rsa.registration1.nmos.tv.key.pem /home/certs/run-certs/rsa.key.pem && \
    cp dhparam.pem /home/certs/run-certs/dhparam.pem

## Create relevant configuration files for Sony Registry and Node
# RUN cd /home/ && mkdir example-conf && mkdir admin
# ADD example-conf /home/example-conf

# ## Get and build source for Sony nmos-js
# RUN cd /home/ && git config --global http.sslVerify false && git clone https://github.com/sony/nmos-js.git

## Custom branding
#COPY NVIDIA_Logo_H_ForScreen_ForLightBG.png nmos-js.patch /home/nmos-js/Development/src/assets/
# RUN cd /home && \
#     mv /home/nmos-js/Development/src/assets/nmos-js.patch /home && \
#     patch -p0 <nmos-js.patch && \
#     rm /home/nmos-js/Development/src/assets/sea-lion.png && \
#     rm nmos-js.patch

## Build and install Sony nmos-js
# RUN cd /home/nmos-js/Development && \
#     yarn install --network-timeout 1000000 && \
#     yarn build && \
#     cp -rf /home/nmos-js/Development/build/* /home/admin

## Move executables, libraries and clean up container as much as possible
RUN cd /home/nmos-cpp/Development/build/ && \
    cp  ./nmos-cpp-node /home && \
    cd /home && rm -rf .git nmos-cpp

## Re-build container for optimised runtime environment using clean Ubuntu Bionic release
FROM ubuntu:jammy

ENV \
  DEBIAN_FRONTEND="noninteractive"

##Copy required files from build container
COPY --from=build-stage /home /home

##Update container with latest patches and needed packages
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    openssl \
    make \
    nano \
    curl \
    jq \
    gnupg \
    rdma-core 

RUN \
    cd /home/mDNSResponder/mDNSPosix && make os=linux install && \
    cd /home && rm -rf /home/mDNSResponder /etc/nsswitch.conf.pre-mdns

RUN \
    curl -sS -k "https://keyserver.ubuntu.com/pks/lookup?op=get&search=0x77b7346a59027b33c10cafe35e64e954262c4500" | apt-key add - && \
    echo "deb http://ppa.launchpad.net/mosquitto-dev/mosquitto-ppa/ubuntu jammy main" | tee /etc/apt/sources.list.d/mosquitto.list && \
    apt-get update && apt-get install -y --no-install-recommends mosquitto libc6 && \
    apt-get remove --purge -y make gnupg && \
    apt-get autoremove -y && \
    apt-get clean -y --no-install-recommends && \
    apt-get autoclean -y --no-install-recommends && \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf /usr/share/doc/ && rm -rf /usr/share/man/ && rm -rf /usr/share/locale/ && \
    rm -rf /usr/local/share/man/* && rm -rf /usr/local/share/.cache/*

# Avahi:    dbus avahi-daemon libavahi-compat-libdnssd-dev libnss-mdns AND NOT make \

##Copy entrypoint.sh script and master config to image
COPY build-nmos-cpp/entrypoint.sh build-nmos-cpp/container-config build-nmos-cpp/node.json /home/

##Set script to executable
RUN chmod +x /home/entrypoint.sh

##Set default config variable to run registry (FALSE) or node (TRUE)
ARG runnode=FALSE
ENV RUN_NODE=$runnode

##Expose correct default ports to allow quick publishing
EXPOSE 8010 8011 11000 11001 1883 5353/udp

WORKDIR /home/
ENTRYPOINT ["/home/entrypoint.sh"]
#CMD []
