# BCS pod launcher

BCS pod launcher starts once Media Proxy Agent instance (on one machine) and MCM Media Proxy instances on each machine. It enables to starts BCS ffmpeg pipeline with bound NMOS client node application.

## Description

The tool can operate in two modes:

- Kubernetes Mode: For multi-node cluster deployment.
- Docker Mode: For single-node using Docker containers.

**Flow (Common to Both Modes)**

1. Run MediaProxy Agent
2. Run MCM Media Proxy
3. Run BcsFfmpeg pipeline with NMOS

In case of docker, MediaProxy/MCM things should only start/run once and on every run of launcher, one should start the app according to input file. It does not store the state of apps, just check appropriate conditions.

In case of kuberenetes, MediaProxy/MCM things should only be run once and BCS pod launcher works as operator in the understanding of Kuberenetes operators within pod. That is way, input file in this way is CustomReaource called BcsConfig.

## Getting Started

### Prerequisites

- go version v1.22.0+
- docker version 17.03+.
- kubectl version v1.27+
- Access to a Kubernetes v1.11.3+ cluster.

### To Run containers on single node

Note that you have to adjust **NMOS** node configuration file. Examples with use cases you can find under the path `<repo>/tests` or `<repo>/launcher/configuration_files` (Files in json format).

Remember the path of above mentioned configuration NMOS file because it must be provided in the next config below: `<repo>/launcher/configuration_files/bcslauncher-static-config.yaml`

```json
      nmosConfigPath: /root/demo
      nmosConfigFileName: intel-node-example.json
```

Edit this configuration file under path `<repo>/launcher/configuration_files/bcslauncher-static-config.yaml`:

```yaml
# 
# SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
# 
# SPDX-License-Identifier: BSD-3-Clause
# 

k8s: false
configuration: # Configuration should be used only for docker mode
  runOnce:
    mediaProxyAgent:
      imageAndTag: mcm/mesh-agent:latest
      gRPCPort: 50051
      restPort: 8100
      network: 
        enable: false
        name: my_net_801f0
        ip: 192.168.2.1
    mediaProxyMcm:
      imageAndTag:  mcm/media-proxy:latest
      interfaceName: eth0
      volumes:
        - /dev/vfio:/dev/vfio
      network: 
        enable: false
        name: my_net_801f0
        ip: 192.168.2.2
  workloadToBeRun:
    ffmpegPipeline:
      name: bcs-ffmpeg-pipeline
      imageAndTag: tiber-broadcast-suite:latest
      gRPCPort: 50051
      sourcePort: 5004
      environmentVariables:
        - "http_proxy="
        - "https_proxy=" 
      volumes:
        videos: /root #for videos
        dri: /usr/lib/x86_64-linux-gnu/dri
        kahawai: /tmp/kahawai_lcore.lock
        devnull: /dev/null
        tmpHugepages: /tmp/hugepages
        hugepages: /hugepages
        imtl: /var/run/imtl
        shm: /dev/shm
      devices:
        vfio: /dev/vfio
        dri: /dev/dri
      network: 
        enable: true
        name: my_net_801f0
        ip: 192.168.2.4
    nmosClient:
      name: bcs-ffmpeg-pipeline-nmos-client
      imageAndTag: tiber-broadcast-suite-nmos-node:latest
      environmentVariables:
        - "http_proxy="
        - "https_proxy=" 
        - "VFIO_PORT_TX=0000:ca:11.0"
      nmosConfigPath: /root/demo
      nmosConfigFileName: intel-node-example.json
      network: 
        enable: true
        name: my_net_801f0
        ip: 192.168.2.5
```

> It is worth noting that workloads under key `runOnce` are configurable globally and only once, whereas `workloadToBeRun` can be defined many times for diffrent workloads (so under the same path, change the content of the file). The flag should be set as `k8s: false`.

Run `<repo>/first_run.sh` script then run `<repo>/build.sh`.
Next, follow all guidelines [here](https://github.com/OpenVisualCloud/Media-Communications-Mesh/blob/main/media-proxy/README.md)

Remember to export RX and TX vfio ports (consequently TX port for sender and RX port por receiver):

``` bash
 # Use dpdk-devbind.py -s to check pci address of vfio device
 export VFIO_PORT_TX="0000:ca:11.0"
 export VFIO_PORT_RX="0000:ca:11.1"
```

```bash

cd <repo>/launcher/cmd/
go build main.go
./main <pass path to file ./launcher/configuration_files/bcslauncher-k8s-config.yaml>
# Alternatively instead of go build main.go && ./main, you can type: go run main.go <pass path to file ./launcher/configuration_files/bcslauncher-k8s-config.yaml>
```

### To Deploy on the cluster

Follow instructions to build minikube cluster: [here](https://github.com/OpenVisualCloud/Media-Communications-Mesh/blob/main/media-proxy/README.md)
**Build image:**

Modify `./launcher/configuration_files/bcslauncher-k8s-config.yaml`. `k8s: true` should be set. Resources for media proxy and msh agent are only configured once.

`docker build -t controller:bcs_pod_launcher .`

Modify `./launcher/configuration_files/bcsconfig-example.yaml` to prepare information for bcs pipeline and nmos node. There may be many custom resources that specifies diffrent `workloads` with nmos node.

**BCS pod launcher installer in k8s cluster:**  

Users can just run kubectl apply -f <file> to install the project:

```bash
cd <repo>/launcher/
kubectl apply -f ./configuration_files/bcslauncher-k8s-config.yaml
kubectl apply -f ./configuration_files/bcsconfig-crd.yaml
kubectl apply -f ./configuration_files/bcs-launcher.yaml
# Adjust to your needs: ./configuration_files/bcsconfig-example.yaml
kubectl apply -f ./configuration_files/bcsconfig-example.yaml
```

**BCS pod launcher roles of files in k8s cluster:**  

- `configuration_files/bcslauncher-k8s-config.yaml` -> configmap for setting up the mode of launcher. `k8s: true` defines kuberenets mode. Currently, you should not modify this in that file.
- `configuration_files/bcs-launcher.yaml` -> install set of kuberenetes resources that are needed to run bcs pod luancher, no additional configuration required
- `configuration_files/bcsconfig-crd.yaml` -> Custom Resource Definition for `BcsConfig`  
- `configuration_files/bcsconfig-example.yaml` -> example `BcsConfig` file that it is an input to provide information about **bcs ffmpeg piepeline and NMOS client**, you can adjust file to your needs,
- `configuration_files/bcslauncher-static-config.yaml` -> static config for docker mode. `k8s: false` defines docker mode. Currently, you should not modify this in that file.

**BCS pod launcher deletion of implementationn of BCS pod launcher in k8s cluster:**  

```bash
cd <repo>/launcher/
kubectl delete -f ./configuration_files/bcslauncher-k8s-config.yaml
kubectl delete -f ./configuration_files/bcsconfig-crd.yaml
kubectl delete -f ./configuration_files/bcs-launcher.yaml
kubectl delete -f ./configuration_files/bcsconfig-example.yaml
```

## License

SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation

SPDX-License-Identifier: BSD-3-Clause

===============================================================

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
