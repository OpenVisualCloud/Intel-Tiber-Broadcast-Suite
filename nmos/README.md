# Environment

- Ubuntu 22.04.4 LTS
- kernel: 5.15.0-88-generic
- minikube: v1.34.0 (1 node)
- docker: 27.3.1
- kubectl: v1.31.2
- NMOS supports standards IS-04 and IS-05: <https://specs.amwa.tv/is-04/>, <https://specs.amwa.tv/is-05/>

## Configuration note for NMOS simplified client in the contect of Intel Broadcast Suite

`nmos-cpp` repository has been simplified to **IS-04** & **IS-05** implementation only.
The key change is in configuration of senders and receivers for BCS pipeline.

BCS Pipeline is a NMOS client that is treated as one node that has 1 device and has x senders and y receivers that are provided from the level of json config `node.json`.
Here is sample config `node.json`:

```json
{
    "logging_level": 0,
    "http_port": 90,
    "activate_senders": true,
    "label": "intel-broadcast-suite",
    "senders": ["v"],
    "senders_count": [1],
    "receivers": ["v"],
    "receivers_count": [0],
    "device_tags": {
        "pipeline": ["tx-sender"]
    },
    "frame_rate": { "numerator": 60, "denominator": 1 },
    "frame_width": 1920,
    "frame_height": 1080,
    "video_type": "video/jxsv",
    "domain": "local",
    "function" : "tx",
    "gpu_hw_acceleration": "none",
    "sender_ffmpeg_video_type": "rawvideo",
    "sender_payload_type": 96,
    "sender_pixel_format": "yuv422p10le",
    "sender_transportFormat": "mcm",
    "sender_conn_type": "st2110",
    "sender_transport": "st2110-20",
    "sender_input_path": "/root",
    "sender_input_path_name": "1920x1080p10le_1.yuv",
    "receiver_transportFormat": "mcm",
    "receiver_conn_type": "st2110",
    "receiver_transport": "st2110-20",
    "ffmpeg_grpc_server_address": "localhost",
    "ffmpeg_grpc_server_port": "50051"
}
```

The crucial params are:  

```

"senders": ["v","d"],
"senders_count": [2, 1],
"receivers": ["v"],
"receivers_count": [4],

```

> `senders` and `receivers` are arrays that specifies the kind of ports. The possible options are video "v", audio "a", data "d", and mux "m"

> `senders_count` and `receivers_count` are corresponding arrays to senders and receivers arrays that provide count by kind of port. For example, for `senders`: `["v", "a", "d"]`, the `senders_count`: `[3, 1, 1]` should be defined. It means that there are 3 senders of type video, 1 sender of type audio and 1 sender of type data.

From point of view POC, only `http_port` has the relevant role. It must be provided in further configurations for example for building process of nmos client node image.

NMOS in BCS is provided in the form of NMOS client node single container that is about to 'stick' to the appropriate BCS pipeline container (within 1 pod).
For testing purposes there is also NMOS registry pod and NMOS testing tool for validation of features.

## Installation

### Docker option

- using docker compose and customized network (bridge)
  
``` bash
cd <repo>/nmos/docker
./run.sh --source-dir <source_dir> --build-dir <build_dir> --patch-dir <patch_dir> --run-dir <RUN_DIR> --update-submodules --apply-patches --build-images --run-docker-compose 
```

- ...or using docker command and running using host network

``` bash
cd <repo>/nmos/docker
./run.sh --source-dir <source_dir> --build-dir <build_dir> --patch-dir <patch_dir> --run-dir <RUN_DIR> --update-submodules --apply-patches --build-images
```

#### Usage and description of options

```

Pattern: ./run.sh --source-dir <source_dir> --build-dir <build_dir> --patch-dir <patch_dir> --run-dir <RUN_DIR> [--prepare-only] [--apply-patches] [--build-images] [--run-docker-compose] [--update-submodules]
  --source-dir <source_dir>  : Absolute path to directory with source code of repository nmos-cpp 3rd party submodule
  --build-dir <build_dir>    : Absolute path to directory with dockerfile and other build files in build-nmos-cpp 3rd party submodule
  --patch-dir <patch_dir>    : Absolute path to directory with patches for both 3rd party submodules
  --run-dir <run_dir>        : Absolute path to directory with run.sh and docker-compose.yaml
  --prepare-only             : Run steps in script that prepares images for nmos but option with runninng docker containers is not applicable
  --build-images             : Build docker images for nmos-client and nmos-registry
  --update-submodules        : Update git submodules
  --apply-patches            : Apply patches for 3rd party submodules
  --run-docker-compose       : Run Docker Compose (nmos-client + nmos-registry + nmos-testing)
                               in customized network of bridge type.
                               Else, by default the <docker run> command will run:
                               (nmos-client + nmos-registry without nmos-testing tool container)
                               in host network

```

### Kubernetes option

Run script that prepares images dor NMOS client node and NMOS registry:

``` bash
cd <repo>/nmos/docker
./run.sh --source-dir <source_dir> --build-dir <build_dir> --patch-dir <patch_dir> --run-dir <RUN_DIR> --update-submodules --apply-patches --build-images --prepare-only
```

```bash
cd <repo>/nmos/k8s
# Install minikube https://minikube.sigs.k8s.io/docs/start/?arch=%2Fwindows%2Fx86-64%2Fstable%2F.exe+download
minikube start
# Build iamges. Refer to 4. Build images
# Adjust ConfigMaps in <repo>/nmos/k8s/nmos-client.yaml, <repo>/nmos/k8s/nmos-registry.yaml and <repo>/nmos/k8s/nmos-testing.yaml
kubectl apply -f <repo>/nmos/k8s/nmos-client.yaml
kubectl apply -f <repo>/nmos/k8s/nmos-registry.yaml
kubectl apply -f <repo>/nmos/k8s/nmos-testing.yaml
# Useful for accessing testing tool browser: https://minikube.sigs.k8s.io/docs/handbook/accessing/
```

### From terminal

#### 1. Git

``` bash
git submodule update --init --recursive
```

#### 2. Patch

``` bash
cd <repo>/nmos
cd build-nmos-cpp
git apply ../patches/build-nmos-cpp.patch
cd ../nmos-cpp
git apply ../patches/nmos-cpp.patch
```

#### 3. Build NMOS binaries (client & registry/controller) (optional for user, useful for dev)

``` bash
cd <repo>/nmos/nmos-cpp/Development/
pip install --upgrade conan~=2.4 
conan profile detect
conan install --requires=nmos-cpp/cci.20240223 --deployer=direct_deploy --build=missing
cd <repo>/nmos/
./run-before-nmos-cpp-build.sh
```

#### 4. Build images

``` bash
cd <repo>/
cp <repo>/nmos/nmos-cpp/Development/nmos-cpp-node/node_implementation.h <repo>/nmos/build-nmos-cpp/
cp <repo>/nmos/nmos-cpp/Development/nmos-cpp-node/node_implementation.cpp <repo>/nmos/build-nmos-cpp/
cp <repo>/nmos/nmos-cpp/Development/nmos-cpp-node/main.cpp <repo>/nmos/build-nmos-cpp/


cd <repo>/nmos/build-nmos-cpp/
make build # build NMOS registry and controller
make buildnode # build NMOS client node

# NMOS testing: https://github.com/AMWA-TV/nmos-testing/blob/master/docs/1.2.%20Installation%20-%20Docker.md

```

#### 5. Docker-compose for running NMOS registry and controller, client and testing tool [OPTION #1]

``` bash
cd <repo>/nmos/docker
# Adjust configs: <repo>/nmos/docker/node.json and registry <repo>/nmos/docker/registry.json and <repo>/nmos/docker/docker-compose.yaml
docker compose up
```

#### 6. Kubernetes for running NMOS registry and controller, client and testing tool [OPTION #2]

``` bash
cd <repo>/nmos/k8s
# Install minikube https://minikube.sigs.k8s.io/docs/start/?arch=%2Fwindows%2Fx86-64%2Fstable%2F.exe+download
minikube start
# Build iamges. Refer to 4. Build images
# Adjust ConfigMaps in <repo>/nmos/k8s/nmos-client.yaml, <repo>/nmos/k8s/nmos-registry.yaml and <repo>/nmos/k8s/nmos-testing.yaml
kubectl apply -f <repo>/nmos/k8s/nmos-client.yaml
kubectl apply -f <repo>/nmos/k8s/nmos-registry.yaml
kubectl apply -f <repo>/nmos/k8s/nmos-testing.yaml
# Useful for accessing testing tool browser: https://minikube.sigs.k8s.io/docs/handbook/accessing/
```

### License

```text
SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation

SPDX-License-Identifier: BSD-3-Clause
```
