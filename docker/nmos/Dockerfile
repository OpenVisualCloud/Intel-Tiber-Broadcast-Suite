# syntax=docker/dockerfile:1

#
# SPDX-License-Identifier: BSD-3-Clause
# Copyright(©) 2024 Intel Corporation
# Intel® Tiber™ Broadcast Suite
#
# build stage

FROM ubuntu:24.04 AS build-stage

USER root

SHELL ["/bin/bash", "-ex", "-o", "pipefail", "-c"]

ARG nproc

# Commit f5497129 corresponds to Conan package nmos-cpp/cci.20240223
ENV \
  TZ="Europe/Warsaw" \
  DEBIAN_FRONTEND="noninteractive" \
  nproc=${nproc:-$(nproc)} \
  NMOS_CPP_VERSION=f54971298c47a633969e9e9adac824b56fc08da7

RUN \
    apt-get update --fix-missing && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    build-essential \
    g++ \
    git \
    curl \
    cmake \
    libboost-all-dev \
    libcpprest-dev \
    libavahi-client-dev \
    libavahi-compat-libdnssd-dev && \
    apt-get autoremove -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN \
    echo "**** DOWNLOAD AND INSTALL gRPC v1.58 ****" && \
    git clone --branch "v1.58.0" --recurse-submodules --depth 1 --shallow-submodules https://github.com/grpc/grpc /tmp/grpc-source && \
    mkdir -p "/tmp/grpc-source/cmake/build" && \
    cmake -S "/tmp/grpc-source" -B "/tmp/grpc-source/cmake/build" -DgRPC_BUILD_TESTS=OFF -DgRPC_INSTALL=ON && \
    make -C "/tmp/grpc-source/cmake/build" -j$(nproc) && \
    make -C "/tmp/grpc-source/cmake/build" install && \
    rm -rf /tmp/grpc-source

COPY src /tmp

## Get source for Sony nmos-cpp library
WORKDIR /tmp
RUN \
    curl --output - -s -k https://codeload.github.com/sony/nmos-cpp/tar.gz/$NMOS_CPP_VERSION | tar zxvf - -C . && \
    mv ./nmos-cpp-${NMOS_CPP_VERSION} ./nmos-cpp

WORKDIR /tmp/nmos-cpp/Development/build
RUN \
    cmake .. -DNMOS_CPP_USE_SUPPLIED_JSON_SCHEMA_VALIDATOR=ON \
    -DNMOS_CPP_USE_SUPPLIED_JWT_CPP=ON \
    -DNMOS_CPP_BUILD_EXAMPLES=OFF \
    -DNMOS_CPP_BUILD_TESTS=OFF && \
    make -j$(nproc) && \
    make install

## mDNSResponder 878.260.1 Build and install
WORKDIR /home
RUN \
    curl --output - -s -k https://codeload.github.com/apple-oss-distributions/mDNSResponder/tar.gz/mDNSResponder-878.260.1 | tar zxvf - -C . && \
    mv ./mDNSResponder-mDNSResponder-878.260.1 /home/mDNSResponder && \
    patch -d mDNSResponder/ -p1 </tmp/nmos-cpp/Development/third_party/mDNSResponder/unicast.patch && \
    patch -d mDNSResponder/ -p1 </tmp/nmos-cpp/Development/third_party/mDNSResponder/permit-over-long-service-types.patch && \
    patch -d mDNSResponder/ -p1 </tmp/nmos-cpp/Development/third_party/mDNSResponder/poll-rather-than-select.patch

WORKDIR /home/mDNSResponder/mDNSPosix
RUN \
    make os=linux && \
    mkdir /tmp/buildout && \
    cp -a /home/mDNSResponder /tmp/buildout/mDNSResponder

## Build gRPC
WORKDIR /tmp/gRPC
RUN \
    cmake -B build && \
    make -j$(nproc) -C build

## Build NMOS Node
WORKDIR /tmp/nmos/nmos-node/build
RUN \
    cmake .. && \
    make -j$(nproc) && \
    cp ./bcs-nmos-node /tmp/buildout

## Re-build container for optimised runtime environment using clean Ubuntu Bionic release
FROM ubuntu:24.04 AS final-stage

ENV \
  DEBIAN_FRONTEND="noninteractive"

##Copy required files from build container
COPY --from=build-stage /tmp/buildout /home
COPY scripts/entrypoint.sh /home/entrypoint.sh

##Update container with latest patches and needed packages
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    make \
    curl \
    gnupg \
    libcpprest-dev && \
    apt-get autoremove -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /home/mDNSResponder/mDNSPosix
RUN \
    make os=linux install

WORKDIR /home
RUN \
    curl -sS -k "https://keyserver.ubuntu.com/pks/lookup?op=get&search=0x77b7346a59027b33c10cafe35e64e954262c4500" | apt-key add - && \
    echo "deb http://ppa.launchpad.net/mosquitto-dev/mosquitto-ppa/ubuntu noble main" | tee /etc/apt/sources.list.d/mosquitto.list && \
    apt-get update && apt-get install -y --no-install-recommends mosquitto libc6 && \
    apt-get remove --purge -y make gnupg && \
    apt-get autoremove -y && \
    apt-get clean -y --no-install-recommends && \
    apt-get autoclean -y --no-install-recommends && \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf /usr/share/doc/ && rm -rf /usr/share/man/ && rm -rf /usr/share/locale/ && \
    rm -rf /usr/local/share/man/* && rm -rf /usr/local/share/.cache/* \
    rm -rf /home/mDNSResponder /etc/nsswitch.conf.pre-mdns

HEALTHCHECK --interval=30s --timeout=5s CMD ps aux | grep "bcs-nmos-node" || exit 1

USER "nmos"

RUN chmod +x /home/entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]
