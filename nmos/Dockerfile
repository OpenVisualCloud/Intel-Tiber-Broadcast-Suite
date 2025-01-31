# syntax=docker/dockerfile:1

#
# SPDX-License-Identifier: BSD-3-Clause
# Copyright(©) 2024 Intel Corporation
# Intel® Tiber™ Broadcast Suite
#
# build stage

FROM ubuntu:jammy AS build-stage

USER root

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
    unzip && \
    pip install --upgrade \
    cmake \
    conan~=2.4

RUN \
    echo "**** DOWNLOAD AND INSTALL gRPC v1.58 ****" && \
    git clone --branch "v1.58.0" --recurse-submodules --depth 1 --shallow-submodules https://github.com/grpc/grpc /tmp/grpc-source && \
    mkdir -p "/tmp/grpc-source/cmake/build" && \
    cmake -S "/tmp/grpc-source" -B "/tmp/grpc-source/cmake/build" -DgRPC_BUILD_TESTS=OFF -DgRPC_INSTALL=ON && \
    make -C "/tmp/grpc-source/cmake/build" -j${nproc} && \
    make -C "/tmp/grpc-source/cmake/build" install && \
    rm -rf /tmp/grpc-source

## Get source for Sony nmos-cpp
## Commit f5497129 corresponds to Conan package nmos-cpp/cci.20240223
ENV NMOS_CPP_VERSION=f54971298c47a633969e9e9adac824b56fc08da7
WORKDIR /home
RUN \
    curl --output - -s -k https://codeload.github.com/sony/nmos-cpp/tar.gz/$NMOS_CPP_VERSION | tar zxvf - -C . && \
    mv ./nmos-cpp-${NMOS_CPP_VERSION} ./nmos-cpp

COPY \
    "patches/" "/home/patches"

RUN \
    patch -d /home/nmos-cpp -p1 -i /home/patches/nmos-cpp.patch

## mDNSResponder 878.260.1 Build and install
RUN \
    curl --output - -s -k https://codeload.github.com/apple-oss-distributions/mDNSResponder/tar.gz/mDNSResponder-878.260.1 | tar zxvf - -C . && \
    mv ./mDNSResponder-mDNSResponder-878.260.1 ./mDNSResponder && \
    patch -d mDNSResponder/ -p1 <nmos-cpp/Development/third_party/mDNSResponder/unicast.patch && \
    patch -d mDNSResponder/ -p1 <nmos-cpp/Development/third_party/mDNSResponder/permit-over-long-service-types.patch && \
    patch -d mDNSResponder/ -p1 <nmos-cpp/Development/third_party/mDNSResponder/poll-rather-than-select.patch

WORKDIR /home/mDNSResponder/mDNSPosix
RUN \
    make os=linux && make os=linux install

## Build Sony nmos-cpp from sources - with workaround to keep latest packages for nmos
## Remember to use the same veriosn of nmos and registry
## First download all packages automatically using conan package manager, than build project
WORKDIR /home/nmos-cpp/Development
RUN \
    conan profile detect && \
    conan install --requires=nmos-cpp/cci.20240223 --deployer=direct_deploy --build=missing

COPY gRPC/ /home/gRPC

WORKDIR /home/gRPC
RUN \
    cmake -B build && \
    make -j${nproc} -C build

RUN \
    mkdir /home/nmos-cpp/Development/build && \
    mkdir /home/nmos-cpp/Development/nmos-cpp-node/build && \
    cp /home/gRPC/build/libFFmpeg_wrapper_client.a /home/nmos-cpp/Development/build && \
    cp /home/gRPC/build/libhw_grpc_proto.a /home/nmos-cpp/Development/build && \
    cp /home/gRPC/config_params.hpp /home/nmos-cpp/Development/nmos-cpp-node && \
    cp /home/gRPC/FFmpeg_wrapper_client.h /home/nmos-cpp/Development/nmos-cpp-node && \
    cp /home/gRPC/build/ffmpeg_cmd_wrap.pb.h /home/nmos-cpp/Development/nmos-cpp-node/build && \
    cp /home/gRPC/build/ffmpeg_cmd_wrap.grpc.pb.h /home/nmos-cpp/Development/nmos-cpp-node/build

ENV LIBRARY_PATH=/home/nmos-cpp/Development/build

WORKDIR /home/nmos-cpp/Development/build
RUN \
    cmake \
    -DCMAKE_PROJECT_TOP_LEVEL_INCLUDES:STRING="third_party/cmake/conan_provider.cmake" \
    -G "Unix Makefiles" \
    -DCMAKE_BUILD_TYPE:STRING="MinSizeRel" \
    -DCMAKE_CONFIGURATION_TYPES:STRING="MinSizeRel" \
    -DCXXFLAGS:STRING="-Os" \
    -DNMOS_CPP_USE_AVAHI:BOOL="0" \
    /home/nmos-cpp/Development/build .. && \
    make -j${nproc} && \
    cp ./nmos-cpp-node /home

## Move executables, libraries and clean up container as much as possible
WORKDIR /home
RUN \
    rm -rf .git nmos-cpp

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

WORKDIR /home/mDNSResponder/mDNSPosix
RUN \
    make os=linux install

WORKDIR /home
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
    rm -rf /usr/local/share/man/* && rm -rf /usr/local/share/.cache/* \
    rm -rf /home/mDNSResponder /etc/nsswitch.conf.pre-mdns

##Copy entrypoint.sh script and master config to image
COPY entrypoint.sh /home/

##Set script to executable
RUN chmod +x /home/entrypoint.sh

##Expose correct default ports to allow quick publishing
EXPOSE 8010 8011 11000 11001 1883 5353/udp

HEALTHCHECK --interval=30s --timeout=5s CMD ps aux | grep "nmos-cpp-node" || exit 1

USER "nmos"

WORKDIR /home/
ENTRYPOINT ["/home/entrypoint.sh"]
