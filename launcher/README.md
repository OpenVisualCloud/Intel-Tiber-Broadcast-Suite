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

```bash
# currently only in dev phase
cd <repo>/launcher/cmd/
# Adjust to your needs: ./configuration_files/bcslauncher-static-config.yaml
go run main.go
```

### To Deploy on the cluster

**Build image:**  
Modify `./launcher/configuration_files/bcslauncher-static-config.yaml`. `k8s: false` defines docker (containers on single host) mode.  
`docker build -t controller:bcs_pod_launcher .`

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
