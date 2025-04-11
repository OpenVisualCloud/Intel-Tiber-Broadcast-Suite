/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"bcs.pod.launcher.intel/resources_library/resources/general"
	"bcs.pod.launcher.intel/resources_library/resources/nmos"
	"bcs.pod.launcher.intel/resources_library/workloads"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"

	bcsv1 "bcs.pod.launcher.intel/api/v1"
)

func TestConstructContainerConfig_MediaProxyAgent(t *testing.T) {
    log := testr.New(t)
    containerInfo := &general.Containers{
        Type: general.MediaProxyAgent,
        Configuration: general.ContainersConfig{
            MediaProxyAgentConfig: workloads.MediaProxyAgentConfig{
                ImageAndTag: "test-image:latest",
                RestPort:    "8080",
                GRPCPort:    "9090",
                Network: workloads.NetworkConfig{
                    Enable: true,
                    Name:   "test-network",
                    IP:     "192.168.1.100",
                },
            },
        },
    }

    containerConfig, hostConfig, networkConfig := ConstructContainerConfig(containerInfo, log)

    assert.NotNil(t, containerConfig)
    assert.NotNil(t, hostConfig)
    assert.NotNil(t, networkConfig)

    assert.Equal(t, "test-image:latest", containerConfig.Image)
    assert.ElementsMatch(t, strslice.StrSlice(strslice.StrSlice{"-c", "8080", "-p", "9090"}), containerConfig.Cmd)

    assert.True(t, hostConfig.Privileged)
    assert.Equal(t, nat.PortMap{
        "8080/tcp": []nat.PortBinding{{HostPort: "8080"}},
        "9090/tcp": []nat.PortBinding{{HostPort: "9090"}},
    }, hostConfig.PortBindings)

    assert.Equal(t, "192.168.1.100", networkConfig.EndpointsConfig["test-network"].IPAMConfig.IPv4Address)
}

func TestConstructContainerConfig_MediaProxyMCM(t *testing.T) {
    log := testr.New(t)
    containerInfo := &general.Containers{
        Type: general.MediaProxyMCM,
        Configuration: general.ContainersConfig{
            MediaProxyMcmConfig: workloads.MediaProxyMcmConfig{
                ImageAndTag:   "mcm-image:latest",
                InterfaceName: "eth0",
                Volumes:       []string{"/host/path:/container/path"},
                Network: workloads.NetworkConfig{
                    Enable: true,
                    Name:   "mcm-network",
                    IP:     "192.168.1.101",
                },
            },
        },
    }

    containerConfig, hostConfig, networkConfig := ConstructContainerConfig(containerInfo, log)

    assert.NotNil(t, containerConfig)
    assert.NotNil(t, hostConfig)
    assert.NotNil(t, networkConfig)

    assert.Equal(t, "mcm-image:latest", containerConfig.Image)
    assert.ElementsMatch(t, []string{"-d", "kernel:eth0", "-i", "localhost"}, containerConfig.Cmd)

    assert.True(t, hostConfig.Privileged)
    assert.Equal(t, []string{"/host/path:/container/path"}, hostConfig.Binds)

    assert.Equal(t, "192.168.1.101", networkConfig.EndpointsConfig["mcm-network"].IPAMConfig.IPv4Address)
}


func TestConstructContainerConfig_BcsPipelineFfmpeg(t *testing.T) {
    log := testr.New(t)
    containerInfo := &general.Containers{
        Type: general.BcsPipelineFfmpeg,
        Configuration: general.ContainersConfig{
            WorkloadConfig: workloads.WorkloadConfig{
                FfmpegPipeline: workloads.FfmpegPipelineConfig{
                    ImageAndTag: "ffmpeg-image:latest",
                    GRPCPort:    50051,
                    EnvironmentVariables: []string{"ENV_VAR=VALUE"},
                    Network: workloads.NetworkConfig{
                        Name: "ffmpeg-network",
                        IP:   "192.168.1.102",
                    },
                    Volumes: workloads.Volumes{
                        Videos: "/host/videos",
                        Dri:    "/host/dri",
                    },
                    Devices: workloads.Devices{
                        Vfio: "/dev/vfio",
                        Dri:  "/dev/dri",
                    },
                },
            },
        },
    }

    containerConfig, hostConfig, networkConfig := ConstructContainerConfig(containerInfo, log)

    assert.NotNil(t, containerConfig)
    assert.NotNil(t, hostConfig)
    assert.NotNil(t, networkConfig)

    assert.Equal(t, "ffmpeg-image:latest", containerConfig.Image)
    assert.ElementsMatch(t, []string{"192.168.1.102", "50051"}, containerConfig.Cmd)
    assert.Equal(t, nat.PortSet{
        "20000/tcp": {},
        "20170/tcp": {},
    }, containerConfig.ExposedPorts)

    assert.True(t, hostConfig.Privileged)
    assert.Equal(t, "192.168.1.102", networkConfig.EndpointsConfig["ffmpeg-network"].IPAMConfig.IPv4Address)
}

func TestUnmarshalK8sConfig_InvalidYAML(t *testing.T) {
    yamlData := `
k8s: true
definition:
  meshAgent:
    image: "mesh-agent-image:latest"
    restPort: "invalid-port"
`
    config, err := UnmarshalK8sConfig([]byte(yamlData))
    assert.Error(t, err)
    assert.Nil(t, config)
}

func TestUnmarshalK8sConfig_MissingFields(t *testing.T) {
    yamlData := `
k8s: true
definition:
  meshAgent:
    image: "mesh-agent-image:latest"
`
    config, err := UnmarshalK8sConfig([]byte(yamlData))
    assert.NoError(t, err)
    assert.NotNil(t, config)
    assert.True(t, config.K8s)
    assert.Equal(t, "mesh-agent-image:latest", config.Definition.MeshAgent.Image)
    assert.Equal(t, 0, config.Definition.MeshAgent.RestPort) // Default value for int
    assert.Equal(t, 0, config.Definition.MeshAgent.GrpcPort) // Default value for int
}
func TestUpdateNmosJsonFile_Success(t *testing.T) {
    tempFile, err := os.CreateTemp("", "nmos_config_*.json")
    assert.NoError(t, err)
    defer os.Remove(tempFile.Name())

    initialConfig := `{
  "logging_level": 10,
  "http_port": 95,
  "label": "intel-broadcast-suite",
  "device_tags": {
    "pipeline": [
      "rx"
    ]
  },
  "function": "rx",
  "gpu_hw_acceleration": "none",
  "domain": "local",
  "ffmpeg_grpc_server_address": "192.168.2.6",
  "ffmpeg_grpc_server_port": "50053",
  "sender_payload_type": 0,
  "sender": [
    {
      "stream_payload": {
        "video": {
          "frame_width": 1920,
          "frame_height": 1080,
          "frame_rate": {
            "numerator": 60,
            "denominator": 1
          },
          "pixel_format": "yuv422p10le",
          "video_type": "rawvideo"
        },
        "audio": {
          "channels": 2,
          "sampleRate": 48000,
          "format": "pcm_s24be",
          "packetTime": "1ms"
        }
      },
      "stream_type": {
        "file": {
          "path": "/videos/recv",
          "filename": "1920x1080p10le_0.yuv"
        }
      }
    }
  ],
  "receiver": [
    {
      "stream_payload": {
        "video": {
          "frame_width": 1920,
          "frame_height": 1080,
          "frame_rate": {
            "numerator": 60,
            "denominator": 1
          },
          "pixel_format": "yuv422p10le",
          "video_type": "rawvideo"
        },
        "audio": {
          "channels": 2,
          "sampleRate": 48000,
          "format": "pcm_s24be",
          "packetTime": "1ms"
        }
      },
      "stream_type": {
        "st2110": {
          "transport": "st2110-20",
          "payloadType": 112
        }
      }
    }
  ]
}`
    _, err = tempFile.WriteString(initialConfig)
    assert.NoError(t, err)
    tempFile.Close()

    err = updateNmosJsonFile(tempFile.Name(), "192.168.1.100", "8080")
    assert.NoError(t, err)

    updatedContent, err := os.ReadFile(tempFile.Name())
    assert.NoError(t, err)

    var updatedConfig nmos.Config
    err = json.Unmarshal(updatedContent, &updatedConfig)
    assert.NoError(t, err)

    assert.Equal(t, "192.168.1.100", updatedConfig.FfmpegGrpcServerAddress)
    assert.Equal(t, "8080", updatedConfig.FfmpegGrpcServerPort)
}

func TestUpdateNmosJsonFile_FileNotFound(t *testing.T) {
    err := updateNmosJsonFile("nonexistent_file.json", "192.168.1.100", "8080")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "no such file")
}

func TestUpdateNmosJsonFile_InvalidJson(t *testing.T) {
    tempFile, err := os.CreateTemp("", "invalid_nmos_config_*.json")
    assert.NoError(t, err)
    defer os.Remove(tempFile.Name())

    _, err = tempFile.WriteString("invalid_json")
    assert.NoError(t, err)
    tempFile.Close()

    err = updateNmosJsonFile(tempFile.Name(), "192.168.1.100", "8080")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "invalid character")
}
func TestFileExists_FileExists(t *testing.T) {
    tempFile, err := os.CreateTemp("", "test_file_exists_*.txt")
    assert.NoError(t, err)
    defer os.Remove(tempFile.Name())

    exists := FileExists(tempFile.Name())
    assert.True(t, exists)
}

func TestFileExists_FileDoesNotExist(t *testing.T) {
    nonExistentFile := "nonexistent_file.txt"
    exists := FileExists(nonExistentFile)
    assert.False(t, exists)
}
func TestConstructContainerConfig_BcsPipelineNmosClient(t *testing.T) {
    log := testr.New(t)
    containerInfo := &general.Containers{
        Type: general.BcsPipelineNmosClient,
        Configuration: general.ContainersConfig{
            WorkloadConfig: workloads.WorkloadConfig{
                NmosClient: workloads.NmosClientConfig{
                    Name:                 "nmos-client",
                    ImageAndTag:              "nmos-client-image:latest",
                    NmosConfigFileName:       "nmos_config.json",
                    NmosConfigPath:           "/host/config",
                    FfmpegConectionAddress:   "192.168.1.103",
                    FfmpegConnectionPort:     "8081",
                    EnvironmentVariables:     []string{"ENV_VAR=VALUE"},
                    Network: workloads.NetworkConfig{
                        Enable: true,
                        Name: "nmos-network",
                        IP:   "192.168.1.103",
                    },
                },
            },
        },
    }

    // Create a temporary NMOS config file
    tempDir := t.TempDir()
    nmosFilePath := tempDir + "/nmos_config.json"
    err := os.WriteFile(nmosFilePath, []byte(`{"ffmpeg_grpc_server_address": "", "ffmpeg_grpc_server_port": ""}`), 0644)
    assert.NoError(t, err)

    // Update the container info to use the temporary file path
    containerInfo.Configuration.WorkloadConfig.NmosClient.NmosConfigPath = tempDir

    containerConfig, hostConfig, networkConfig := ConstructContainerConfig(containerInfo, log)

    assert.NotNil(t, containerConfig)
    assert.NotNil(t, hostConfig)
    assert.NotNil(t, networkConfig)

    assert.Equal(t, "nmos-client-image:latest", containerConfig.Image)
    assert.ElementsMatch(t, []string{"config/nmos_config.json"}, containerConfig.Cmd)
    assert.ElementsMatch(t, []string{"ENV_VAR=VALUE"}, containerConfig.Env)

    assert.True(t, hostConfig.Privileged)
    assert.ElementsMatch(t, []string{fmt.Sprintf("%s:/home/config/", tempDir)}, hostConfig.Binds)

    assert.Equal(t, "192.168.1.103", networkConfig.EndpointsConfig["nmos-network"].IPAMConfig.IPv4Address)
    assert.ElementsMatch(t, []string{"nmos-network"}, networkConfig.EndpointsConfig["nmos-network"].Aliases)
}

func TestConstructContainerConfig_UnknownType(t *testing.T) {
    log := testr.New(t)
    containerInfo := &general.Containers{
    Type: 7, // based on enum iota
    }

    containerConfig, hostConfig, networkConfig := ConstructContainerConfig(containerInfo, log)

    assert.Nil(t, containerConfig)
    assert.Nil(t, hostConfig)
    assert.Nil(t, networkConfig)
}
func TestConstructContainerConfig_InvalidType(t *testing.T) {
    log := testr.New(t)
    containerInfo := &general.Containers{
        Type: -1, // Invalid type
    }

    containerConfig, hostConfig, networkConfig := ConstructContainerConfig(containerInfo, log)

    assert.Nil(t, containerConfig)
    assert.Nil(t, hostConfig)
    assert.Nil(t, networkConfig)
}

func TestConstructContainerConfig_MediaProxyAgent_NoNetwork(t *testing.T) {
    log := testr.New(t)
    containerInfo := &general.Containers{
        Type: general.MediaProxyAgent,
        Configuration: general.ContainersConfig{
            MediaProxyAgentConfig: workloads.MediaProxyAgentConfig{
                ImageAndTag: "test-image:latest",
                RestPort:    "8080",
                GRPCPort:    "9090",
                Network: workloads.NetworkConfig{
                    Enable: false, // Network disabled
                },
            },
        },
    }

    containerConfig, hostConfig, networkConfig := ConstructContainerConfig(containerInfo, log)

    assert.NotNil(t, containerConfig)
    assert.NotNil(t, hostConfig)
    assert.NotNil(t, networkConfig)

    assert.Equal(t, "test-image:latest", containerConfig.Image)
    assert.ElementsMatch(t, strslice.StrSlice{"-c", "8080", "-p", "9090"}, containerConfig.Cmd)

    assert.True(t, hostConfig.Privileged)
    assert.Equal(t, nat.PortMap{
        "8080/tcp": []nat.PortBinding{{HostPort: "8080"}},
        "9090/tcp": []nat.PortBinding{{HostPort: "9090"}},
    }, hostConfig.PortBindings)

    assert.Empty(t, networkConfig.EndpointsConfig)
}


func TestConstructContainerConfig_BcsPipelineNmosClient_MissingConfigFile(t *testing.T) {
    log := testr.New(t)
    containerInfo := &general.Containers{
        Type: general.BcsPipelineNmosClient,
        Configuration: general.ContainersConfig{
            WorkloadConfig: workloads.WorkloadConfig{
                NmosClient: workloads.NmosClientConfig{
                    ImageAndTag:        "nmos-client-image:latest",
                    NmosConfigFileName: "missing_config.json",
                    NmosConfigPath:     "/nonexistent/path",
                    Network: workloads.NetworkConfig{
                        Name: "nmos-network",
                        IP:   "192.168.1.103",
                    },
                },
            },
        },
    }

    containerConfig, hostConfig, networkConfig := ConstructContainerConfig(containerInfo, log)

    assert.Nil(t, containerConfig)
    assert.Nil(t, hostConfig)
    assert.Nil(t, networkConfig)
}
// type ContainerControllerMock interface {
func TestConstructContainerConfig_InvalidConfiguration(t *testing.T) {
    log := testr.New(t)
    containerInfo := &general.Containers{
        Type: general.MediaProxyAgent,
        Configuration: general.ContainersConfig{
            MediaProxyAgentConfig: workloads.MediaProxyAgentConfig{
                ImageAndTag: "",
                RestPort:    "",
                GRPCPort:    "",
                Network: workloads.NetworkConfig{
                    Enable: true,
                    Name:   "",
                    IP:     "",
                },
            },
        },
    }

    containerConfig, hostConfig, networkConfig := ConstructContainerConfig(containerInfo, log)

    assert.NotNil(t, containerConfig)
    assert.NotNil(t, hostConfig)
    assert.NotNil(t, networkConfig)

    assert.Equal(t, "", containerConfig.Image)
    assert.ElementsMatch(t, []string{"-c", "", "-p", ""}, containerConfig.Cmd)

    assert.True(t, hostConfig.Privileged)
    assert.Equal(t, nat.PortMap{
        "/tcp": []nat.PortBinding{{HostPort: ""}},
    }, hostConfig.PortBindings)

    assert.Equal(t, "", networkConfig.EndpointsConfig[""].IPAMConfig.IPv4Address)
}

func TestConstructContainerConfig_BcsPipelineNmosClient_InvalidJsonUpdate(t *testing.T) {
    log := testr.New(t)
    containerInfo := &general.Containers{
        Type: general.BcsPipelineNmosClient,
        Configuration: general.ContainersConfig{
            WorkloadConfig: workloads.WorkloadConfig{
                NmosClient: workloads.NmosClientConfig{
                    ImageAndTag:        "nmos-client-image:latest",
                    NmosConfigFileName: "invalid_config.json",
                    NmosConfigPath:     "/tmp",
                    FfmpegConectionAddress: "192.168.1.103",
                    FfmpegConnectionPort:   "8081",
                    EnvironmentVariables:   []string{"ENV_VAR=VALUE"},
                    Network: workloads.NetworkConfig{
                        Name: "nmos-network",
                        IP:   "192.168.1.103",
                    },
                },
            },
        },
    }

    // Create an invalid NMOS config file
    tempDir := t.TempDir()
    nmosFilePath := tempDir + "/invalid_config.json"
    err := os.WriteFile(nmosFilePath, []byte("invalid_json"), 0644)
    assert.NoError(t, err)

    // Update the container info to use the temporary file path
    containerInfo.Configuration.WorkloadConfig.NmosClient.NmosConfigPath = tempDir

    containerConfig, hostConfig, networkConfig := ConstructContainerConfig(containerInfo, log)

    assert.Nil(t, containerConfig)
    assert.Nil(t, hostConfig)
    assert.Nil(t, networkConfig)
}
func TestBoolPtr_True(t *testing.T) {
    result := boolPtr(true)
    assert.NotNil(t, result)
    assert.True(t, *result)
}

func TestBoolPtr_False(t *testing.T) {
    result := boolPtr(false)
    assert.NotNil(t, result)
    assert.False(t, *result)
}
func TestCreateMeshAgentDeployment_ValidConfig(t *testing.T) {
    cm := &corev1.ConfigMap{
        Data: map[string]string{
            "config.yaml": `
k8s: true
definition:
  meshAgent:
    image: "mesh-agent-image:latest"
    restPort: 8080
    grpcPort: 9090
`,
        },
    }

    deployment := CreateMeshAgentDeployment(cm)
    assert.NotNil(t, deployment)
    assert.Equal(t, "mesh-agent-deployment", deployment.ObjectMeta.Name)
    assert.Equal(t, "mcm", deployment.ObjectMeta.Namespace)
    assert.Equal(t, "mesh-agent", deployment.Spec.Template.Labels["app"])
    assert.Equal(t, "mesh-agent-image:latest", deployment.Spec.Template.Spec.Containers[0].Image)
    assert.ElementsMatch(t, []string{"mesh-agent", "-c", "8080", "-p", "9090"}, deployment.Spec.Template.Spec.Containers[0].Command)
    assert.Equal(t, int32(8080), deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort)
    assert.Equal(t, int32(9090), deployment.Spec.Template.Spec.Containers[0].Ports[1].ContainerPort)
    assert.NotNil(t, deployment.Spec.Template.Spec.Containers[0].SecurityContext)
    assert.True(t, *deployment.Spec.Template.Spec.Containers[0].SecurityContext.Privileged)
}

func TestCreateMeshAgentDeployment_InvalidConfig(t *testing.T) {
    cm := &corev1.ConfigMap{
        Data: map[string]string{
            "config.yaml": `
invalid_yaml
`,
        },
    }

    deployment := CreateMeshAgentDeployment(cm)
    assert.Nil(t, deployment)
}

func TestCreateMeshAgentDeployment_MissingConfig(t *testing.T) {
    cm := &corev1.ConfigMap{
        Data: map[string]string{},
    }

    deployment := CreateMeshAgentDeployment(cm)
    assert.Equal(t,  types.UID(""), deployment.ObjectMeta.UID)
}
func TestCreateMeshAgentService_ValidConfig(t *testing.T) {
    cm := &corev1.ConfigMap{
        Data: map[string]string{
            "config.yaml": `
k8s: true
definition:
  meshAgent:
    image: "mesh-agent-image:latest"
    restPort: 8080
    grpcPort: 9090
`,
        },
    }

    service := CreateMeshAgentService(cm)
    assert.NotNil(t, service)
    assert.Equal(t, "mesh-agent-service", service.ObjectMeta.Name)
    assert.Equal(t, "mcm", service.ObjectMeta.Namespace)
    assert.Equal(t, map[string]string{"app": "mesh-agent"}, service.Spec.Selector)

    assert.Len(t, service.Spec.Ports, 2)
    assert.Equal(t, "rest", service.Spec.Ports[0].Name)
    assert.Equal(t, corev1.ProtocolTCP, service.Spec.Ports[0].Protocol)
    assert.Equal(t, int32(8080), service.Spec.Ports[0].Port)
    assert.Equal(t, intstr.FromInt(8080), service.Spec.Ports[0].TargetPort)

    assert.Equal(t, "grpc", service.Spec.Ports[1].Name)
    assert.Equal(t, corev1.ProtocolTCP, service.Spec.Ports[1].Protocol)
    assert.Equal(t, int32(9090), service.Spec.Ports[1].Port)
    assert.Equal(t, intstr.FromInt(9090), service.Spec.Ports[1].TargetPort)
}

func TestCreateMeshAgentService_InvalidConfig(t *testing.T) {
    cm := &corev1.ConfigMap{
        Data: map[string]string{
            "config.yaml": `
invalid_yaml
`,
        },
    }

    service := CreateMeshAgentService(cm)
    assert.Nil(t, service)
}

func TestCreateMeshAgentService_MissingConfig(t *testing.T) {
    cm := &corev1.ConfigMap{
        Data: map[string]string{},
    }
    service := CreateMeshAgentService(cm)
    assert.Equal(t, types.UID(""), service.ObjectMeta.UID)
}
func TestCreateService_ValidName(t *testing.T) {
    serviceName := "test-service"
    service := CreateService(serviceName)

    assert.NotNil(t, service)
    assert.Equal(t, serviceName, service.ObjectMeta.Name)
    assert.Equal(t, "default", service.ObjectMeta.Namespace)
    assert.Equal(t, map[string]string{"app": serviceName}, service.Spec.Selector)

    assert.Len(t, service.Spec.Ports, 1)
    assert.Equal(t, corev1.ProtocolTCP, service.Spec.Ports[0].Protocol)
    assert.Equal(t, int32(80), service.Spec.Ports[0].Port)
}

func TestCreateService_EmptyName(t *testing.T) {
    serviceName := ""
    service := CreateService(serviceName)

    assert.NotNil(t, service)
    assert.Equal(t, serviceName, service.ObjectMeta.Name)
    assert.Equal(t, "default", service.ObjectMeta.Namespace)
    assert.Equal(t, map[string]string{"app": serviceName}, service.Spec.Selector)

    assert.Len(t, service.Spec.Ports, 1)
    assert.Equal(t, corev1.ProtocolTCP, service.Spec.Ports[0].Protocol)
    assert.Equal(t, int32(80), service.Spec.Ports[0].Port)
}

func TestCreateService_NamespaceDefault(t *testing.T) {
    serviceName := "test-service"
    service := CreateService(serviceName)

    assert.NotNil(t, service)
    assert.Equal(t, "default", service.ObjectMeta.Namespace)
}
func TestCreateBcsService_ValidConfig(t *testing.T) {
    bcs := &bcsv1.BcsConfig{
        Spec: bcsv1.BcsConfigSpec{
            Name:      "test-bcs",
            Namespace: "test-namespace",
            Nmos: bcsv1.Nmos{
                NmosApiPort:                3000,
                NmosApiNodePort:            32000,
                NmosAppCommunicationPort:   4000,
                NmosAppCommunicationNodePort: 42000,
            },
        },
    }

    service := CreateBcsService(bcs)

    assert.NotNil(t, service)
    assert.Equal(t, "test-bcs", service.ObjectMeta.Name)
    assert.Equal(t, "test-namespace", service.ObjectMeta.Namespace)
    assert.Equal(t, corev1.ServiceTypeNodePort, service.Spec.Type)
    assert.Equal(t, map[string]string{"app": "test-bcs"}, service.Spec.Selector)

    assert.Len(t, service.Spec.Ports, 2)

    assert.Equal(t, "nmos-node-api", service.Spec.Ports[0].Name)
    assert.Equal(t, corev1.ProtocolTCP, service.Spec.Ports[0].Protocol)
    assert.Equal(t, int32(3000), service.Spec.Ports[0].Port)
    assert.Equal(t, intstr.FromInt(3000), service.Spec.Ports[0].TargetPort)
    assert.Equal(t, int32(32000), service.Spec.Ports[0].NodePort)

    assert.Equal(t, "nmos-app-communication", service.Spec.Ports[1].Name)
    assert.Equal(t, corev1.ProtocolTCP, service.Spec.Ports[1].Protocol)
    assert.Equal(t, int32(4000), service.Spec.Ports[1].Port)
    assert.Equal(t, intstr.FromInt(4000), service.Spec.Ports[1].TargetPort)
    assert.Equal(t, int32(42000), service.Spec.Ports[1].NodePort)
}

func TestCreateBcsService_MissingConfig(t *testing.T) {
    bcs := &bcsv1.BcsConfig{
        Spec: bcsv1.BcsConfigSpec{
            Name:      "",
            Namespace: "",
        },
    }

    service := CreateBcsService(bcs)

    assert.NotNil(t, service)
    assert.Equal(t, "", service.ObjectMeta.Name)
    assert.Equal(t, "", service.ObjectMeta.Namespace)
    assert.Equal(t, corev1.ServiceTypeNodePort, service.Spec.Type)
    assert.Equal(t, map[string]string{"app": ""}, service.Spec.Selector)
}
func TestCreatePersistentVolume_ValidConfig(t *testing.T) {
    cm := &corev1.ConfigMap{
        Data: map[string]string{
            "config.yaml": `
k8s: true
definition:
  mediaProxy:
    pvHostPath: "/mnt/data"
    pvStorageClass: "manual"
    pvStorage: "10Gi"
`,
        },
    }

    pv := CreatePersistentVolume(cm)
    assert.NotNil(t, pv)
    assert.Equal(t, "mtl-pv", pv.ObjectMeta.Name)
    assert.Equal(t, v1.ResourceList{"storage": resource.MustParse("10Gi")}, pv.Spec.Capacity)
    assert.Equal(t, corev1.PersistentVolumeFilesystem, *pv.Spec.VolumeMode)
    assert.Equal(t, corev1.ReadWriteOnce, pv.Spec.AccessModes[0])
    assert.Equal(t, corev1.PersistentVolumeReclaimRetain, pv.Spec.PersistentVolumeReclaimPolicy)
    assert.Equal(t, "manual", pv.Spec.StorageClassName)
    assert.Equal(t, "/mnt/data", pv.Spec.PersistentVolumeSource.HostPath.Path)
}

func TestCreatePersistentVolume_InvalidConfig(t *testing.T) {
    cm := &corev1.ConfigMap{
        Data: map[string]string{
            "config.yaml": `
invalid_yaml
`,
        },
    }

    pv := CreatePersistentVolume(cm)
    assert.Nil(t, pv)
}
func TestCreateConfigMap_ValidConfig(t *testing.T) {
    bcs := &bcsv1.BcsConfig{
        Spec: bcsv1.BcsConfigSpec{
            Name:      "test-bcs",
            Namespace: "test-namespace",
            Nmos: bcsv1.Nmos{
                NmosInputFile: bcsv1.NmosInputFile{
                    LoggingLevel:            10,
                    HttpPort:                8080,
                    Label:                   "test-label",
                    DeviceTags:              bcsv1.DeviceTags{Pipeline: []string{"rx"}},
                    Function:                "rx",
                    GpuHwAcceleration:       "none",
                    Domain:                  "local",
                    FfmpegGrpcServerAddress: "192.168.1.100",
                    FfmpegGrpcServerPort:    "50051",
                    SenderPayloadType:       96,
                    Sender: []bcsv1.Sender{
                        {
                            StreamPayload: bcsv1.StreamPayload{
                                Video: bcsv1.Video{
                                    FrameWidth:  1920,
                                    FrameHeight: 1080,
                                    FrameRate: bcsv1.FrameRate{
                                        Numerator:   30,
                                        Denominator: 1,
                                    },
                                    PixelFormat: "yuv420p",
                                    VideoType:   "rawvideo",
                                },
                                Audio: bcsv1.Audio{
                                    Channels:   2,
                                    SampleRate: 48000,
                                    Format:     "pcm_s16le",
                                    PacketTime: "1ms",
                                },
                            },
                            StreamType: bcsv1.StreamType{
                                File: &bcsv1.File{
                                    Path:     "/videos",
                                    Filename: "test_video.yuv",
                                },
                            },
                        },
                    },
                    Receiver: []bcsv1.Receiver{
                        {
                            StreamPayload: bcsv1.StreamPayload{
                                Video: bcsv1.Video{
                                    FrameWidth:  1920,
                                    FrameHeight: 1080,
                                    FrameRate: bcsv1.FrameRate{
                                        Numerator:   30,
                                        Denominator: 1,
                                    },
                                    PixelFormat: "yuv420p",
                                    VideoType:   "rawvideo",
                                },
                                Audio: bcsv1.Audio{
                                    Channels:   2,
                                    SampleRate: 48000,
                                    Format:     "pcm_s16le",
                                    PacketTime: "1ms",
                                },
                            },
                            StreamType: bcsv1.StreamType{
                                St2110: &bcsv1.St2110{
                                    Transport:   "st2110-20",
                                    Payload_type : 112,
                                },
                            },
                        },
                    },
                },
            },
        },
    }

	configMap := CreateConfigMap(bcs)
	assert.NotNil(t, configMap)
	assert.Equal(t, "test-bcs-config", configMap.ObjectMeta.Name)
	assert.Equal(t, "test-namespace", configMap.ObjectMeta.Namespace)

	var data map[string]interface{}
	err := json.Unmarshal([]byte(configMap.Data["config"]), &data)
	assert.NoError(t, err)
	assert.Equal(t, float64(8080), data["http_port"])
	assert.Equal(t, "rx", data["function"])
    assert.Equal(t, "test-label", data["label"])
    assert.Equal(t, "none", data["gpu_hw_acceleration"])
    assert.Equal(t, "local", data["domain"])
    assert.Equal(t, "192.168.1.100", data["ffmpeg_grpc_server_address"])
    assert.Equal(t, "50051", data["ffmpeg_grpc_server_port"])
    assert.Equal(t, float64(96), data["sender_payload_type"])

    senders := data["sender"].([]interface{})
    assert.Len(t, senders, 1)
    sender := senders[0].(map[string]interface{})
    streamPayload := sender["stream_payload"].(map[string]interface{})
    video := streamPayload["video"].(map[string]interface{})
    assert.Equal(t, float64(1920), video["frame_width"])
    assert.Equal(t, float64(1080), video["frame_height"])
    assert.Equal(t, "yuv420p", video["pixel_format"])
    assert.Equal(t, "rawvideo", video["video_type"])

    audio := streamPayload["audio"].(map[string]interface{})
    assert.Equal(t, float64(2), audio["channels"])
    assert.Equal(t, float64(48000), audio["sampleRate"])
    assert.Equal(t, "pcm_s16le", audio["format"])

    receivers := data["receiver"].([]interface{})
    assert.Len(t, receivers, 1)
    receiver := receivers[0].(map[string]interface{})
    receiverStreamPayload := receiver["stream_payload"].(map[string]interface{})
    receiverVideo := receiverStreamPayload["video"].(map[string]interface{})
    assert.Equal(t, float64(1920), receiverVideo["frame_width"])
    assert.Equal(t, float64(1080), receiverVideo["frame_height"])
    assert.Equal(t, "yuv420p", receiverVideo["pixel_format"])
    assert.Equal(t, "rawvideo", receiverVideo["video_type"])

    receiverAudio := receiverStreamPayload["audio"].(map[string]interface{})
    assert.Equal(t, float64(2), receiverAudio["channels"])
    assert.Equal(t, float64(48000), receiverAudio["sampleRate"])
    assert.Equal(t, "pcm_s16le", receiverAudio["format"])
}

func TestCreateConfigMap_InvalidJson(t *testing.T) {
	bcs := &bcsv1.BcsConfig{
		Spec: bcsv1.BcsConfigSpec{
			Name:      "test-bcs",
			Namespace: "test-namespace",
            Nmos: bcsv1.Nmos{},
		},
	}

	configMap := CreateConfigMap(bcs)
	assert.Equal(t, types.UID(""), configMap.ObjectMeta.UID)
}
func TestCreateBcsDeployment_ValidConfig(t *testing.T) {
	bcs := &bcsv1.BcsConfig{
		Spec: bcsv1.BcsConfigSpec{
			Name:      "test-bcs",
			Namespace: "test-namespace",
			Nmos: bcsv1.Nmos{
				Image: "nmos-image:latest",
				Args:  []string{"--arg1", "--arg2"},
				EnvironmentVariables: []bcsv1.EnvVar{
					{Name: "ENV_VAR_1", Value: "VALUE_1"},
					{Name: "ENV_VAR_2", Value: "VALUE_2"},
				},
			},
			App: bcsv1.App{
				Image:    "app-image:latest",
				GrpcPort: 50051,
				Volumes: map[string]string{
					"videos":       "/host/videos",
					"dri":          "/host/dri",
					"kahawaiLock":  "/host/kahawai.lock",
					"devNull":      "/dev/null",
					"hugepagesTmp": "/host/tmp/hugepages",
					"hugepages":    "/host/hugepages",
					"imtl":         "/host/imtl",
					"shm":          "/dev/shm",
					"vfio":         "/dev/vfio",
					"driDev":       "/dev/dri",
				},
				EnvironmentVariables: []bcsv1.EnvVar{
					{Name: "APP_ENV_VAR_1", Value: "APP_VALUE_1"},
					{Name: "APP_ENV_VAR_2", Value: "APP_VALUE_2"},
				},
			},
		},
	}

	deployment := CreateBcsDeployment(bcs)

	assert.NotNil(t, deployment)
	assert.Equal(t, "test-bcs", deployment.ObjectMeta.Name)
	assert.Equal(t, "test-namespace", deployment.ObjectMeta.Namespace)

	assert.Len(t, deployment.Spec.Template.Spec.Containers, 2)

	nmosContainer := deployment.Spec.Template.Spec.Containers[0]
	assert.Equal(t, "tiber-broadcast-suite-nmos-node", nmosContainer.Name)
	assert.Equal(t, "nmos-image:latest", nmosContainer.Image)
	assert.ElementsMatch(t, []string{"--arg1", "--arg2"}, nmosContainer.Args)
	assert.ElementsMatch(t, []corev1.EnvVar{
		{Name: "ENV_VAR_1", Value: "VALUE_1"},
		{Name: "ENV_VAR_2", Value: "VALUE_2"},
	}, nmosContainer.Env)

	appContainer := deployment.Spec.Template.Spec.Containers[1]
	assert.Equal(t, "tiber-broadcast-suite", appContainer.Name)
	assert.Equal(t, "app-image:latest", appContainer.Image)
	assert.ElementsMatch(t, []string{"localhost", "50051"}, appContainer.Args)
	assert.ElementsMatch(t, []corev1.EnvVar{
		{Name: "APP_ENV_VAR_1", Value: "APP_VALUE_1"},
		{Name: "APP_ENV_VAR_2", Value: "APP_VALUE_2"},
	}, appContainer.Env)

	assert.Len(t, deployment.Spec.Template.Spec.Volumes, 11)
	assert.Equal(t, "videos", deployment.Spec.Template.Spec.Volumes[0].Name)
	assert.Equal(t, "/host/videos", deployment.Spec.Template.Spec.Volumes[0].VolumeSource.HostPath.Path)
}

func TestCreateBcsDeployment_MissingConfig(t *testing.T) {
	bcs := &bcsv1.BcsConfig{
		Spec: bcsv1.BcsConfigSpec{
			Name:      "",
			Namespace: "",
		},
	}

	deployment := CreateBcsDeployment(bcs)

	assert.NotNil(t, deployment)
	assert.Equal(t, types.UID(""), deployment.ObjectMeta.UID)
}
func TestCreateDaemonSet_ValidConfig(t *testing.T) {
    cm := &corev1.ConfigMap{
        Data: map[string]string{
            "config.yaml": `
k8s: true
definition:
  mediaProxy:
    image: "media-proxy-image:latest"
    command: ["media-proxy"]
    args: ["--grpc-port", "50051", "--sdk-port", "50052"]
    grpcPort: 50051
    sdkPort: 50052
    volumes:
      memif: "/run/mcm"
      vfio: "/dev/vfio"
    pvHostPath: "/mnt/data"
    pvStorageClass: "manual"
    pvStorage: "10Gi"
    pvcStorage: "5Gi"
`,
        },
    }

    daemonSet := CreateDaemonSet(cm)
    assert.NotNil(t, daemonSet)
    assert.Equal(t, "media-proxy", daemonSet.ObjectMeta.Name)
    assert.Equal(t, "mcm", daemonSet.ObjectMeta.Namespace)
    assert.Equal(t, map[string]string{"app": "media-proxy"}, daemonSet.Spec.Selector.MatchLabels)
    assert.Equal(t, "media-proxy-image:latest", daemonSet.Spec.Template.Spec.Containers[0].Image)
    assert.ElementsMatch(t, []string{"media-proxy"}, daemonSet.Spec.Template.Spec.Containers[0].Command)
    assert.ElementsMatch(t, []string{"--grpc-port", "50051", "--sdk-port", "50052"}, daemonSet.Spec.Template.Spec.Containers[0].Args)
    assert.Equal(t, int32(50051), daemonSet.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort)
    assert.Equal(t, int32(50052), daemonSet.Spec.Template.Spec.Containers[0].Ports[1].ContainerPort)
    assert.Equal(t, "/run/mcm", daemonSet.Spec.Template.Spec.Volumes[0].HostPath.Path)
    assert.Equal(t, "/dev/vfio", daemonSet.Spec.Template.Spec.Volumes[1].HostPath.Path)
}

func TestCreateDaemonSet_InvalidConfig(t *testing.T) {
    cm := &corev1.ConfigMap{
        Data: map[string]string{
            "config.yaml": `
invalid_yaml
`,
        },
    }

    daemonSet := CreateDaemonSet(cm)
    assert.Nil(t, daemonSet)
}

func TestCreateDaemonSet_MissingConfig(t *testing.T) {
    cm := &corev1.ConfigMap{
        Data: map[string]string{},
    }

    daemonSet := CreateDaemonSet(cm)
    assert.Equal(t,types.UID(""), daemonSet.ObjectMeta.UID)
}
func TestConvertEnvVars_ValidInput(t *testing.T) {
	input := []bcsv1.EnvVar{
		{Name: "ENV_VAR_1", Value: "VALUE_1"},
		{Name: "ENV_VAR_2", Value: "VALUE_2"},
	}

	expected := []corev1.EnvVar{
		{Name: "ENV_VAR_1", Value: "VALUE_1"},
		{Name: "ENV_VAR_2", Value: "VALUE_2"},
	}

	result := convertEnvVars(input)
	assert.ElementsMatch(t, expected, result)
}

func TestConvertEnvVars_EmptyInput(t *testing.T) {
	input := []bcsv1.EnvVar{}

	expected := []corev1.EnvVar{}

	result := convertEnvVars(input)
	assert.ElementsMatch(t, expected, result)
}

func TestConvertEnvVars_NilInput(t *testing.T) {
	var input []bcsv1.EnvVar

	expected := []corev1.EnvVar{}

	result := convertEnvVars(input)
	assert.ElementsMatch(t, expected, result)
}
func TestCreateNamespace_ValidName(t *testing.T) {
	namespaceName := "test-namespace"
	namespace := CreateNamespace(namespaceName)

	assert.NotNil(t, namespace)
	assert.Equal(t, namespaceName, namespace.ObjectMeta.Name)
}

func TestCreateNamespace_EmptyName(t *testing.T) {
	namespaceName := ""
	namespace := CreateNamespace(namespaceName)

	assert.NotNil(t, namespace)
	assert.Equal(t, namespaceName, namespace.ObjectMeta.Name)
}
// }

// type DockerContainerControllerMock struct {
// 	cli *client.Client
//     mock.Mock
// }

// func (m *DockerContainerControllerMock) ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error) {
// 	args := m.Called(ctx, options)
// 	return args.Get(0).([]image.Summary), args.Error(1)
// }

// // ListImages lists images using the Docker client
// func ListImages(client* DockerContainerControllerMock) ([]image.Summary, error) {
// 	return client.ImageList(context.Background(), image.ListOptions{})
// }

// func mockPullImage(ctx context.Context, cli *client.Client, imageName string) error {
// 	out, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
// 	if err != nil {
// 		return err
// 	}
// 	defer out.Close()

// 	// Copy the output to stdout
// 	_, err = io.Copy(os.Stdout, out)
// 	return err
// }

// func mockDeleteImage(ctx context.Context, cli *client.Client, imageName string) error {
// 	removedImages, err := cli.ImageRemove(ctx, imageName, image.RemoveOptions{Force: true})
// 	if err != nil {
// 		return err
// 	}
// 	for _, removedImage := range removedImages {
// 		fmt.Printf("Deleted image: %s\n", removedImage.Untagged)
// 		if removedImage.Deleted != "" {
// 			fmt.Printf("Deleted image ID: %s\n", removedImage.Deleted)
// 		}
// 	}

// 	return nil
// }


// func TestIsImagePulled_ImageExists(t *testing.T) {
// 	ctx := context.Background()
//     imageMock := "busybox:latest"
// 	mockClient := new(DockerContainerControllerMock)

// 	mockClient.On("ImageList", ctx, image.ListOptions{}).Return([]image.Summary{
// 		{RepoTags: []string{imageMock}},
// 	}, nil)

//     cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// 	client := DockerContainerControllerMock{cli: cli}

//     if err := mockPullImage(ctx, client.cli, imageMock); err != nil {
// 		fmt.Printf("Error pulling image: %v\n", err)
// 	} else {
// 		fmt.Println("Image pulled successfully.")
// 	}

// 	err, isPulled := isImagePulled(ctx, client.cli, imageMock)
// 	assert.NoError(t, err)
// 	assert.True(t, isPulled)
// 	client.AssertExpectations(t)
// }

// func TestIsImagePulled_ImageDoesNotExist(t *testing.T) {
// 	ctx := context.Background()
// 	imageMock := "busybox:latest"
// 	mockClient := new(DockerContainerControllerMock)

//     cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// 	client := DockerContainerControllerMock{cli: cli}

// 	mockClient.On("ImageList", ctx, image.ListOptions{}).Return([]image.Summary{
// 		{RepoTags: []string{imageMock}},
// 	}, nil)

//     mockClient.On("ImageRemove", ctx, imageMock, image.RemoveOptions{Force: true}).Return([]image.DeleteResponse{}, nil)

// 	err, isPulled := isImagePulled(ctx, client.cli, imageMock)
// 	assert.NoError(t, err)
// 	assert.False(t, isPulled)
// 	client.AssertExpectations(t)
// }

// func TestIsImagePulled_ErrorFetchingImages(t *testing.T) {
// 	ctx := context.Background()
// 	imageMock := "abcd:latest"
// 	mockClient := new(DockerContainerControllerMock)

//     cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// 	client := DockerContainerControllerMock{cli: cli}

// 	mockClient.On("ImageList", ctx, image.ListOptions{}).Return(nil, errors.New("failed to list images"))

// 	_, isPulled := isImagePulled(ctx, client.cli, imageMock)
// 	assert.False(t, isPulled)
// 	client.AssertExpectations(t)
// }

// func TestPullImageIfNotExists_ImageNotPulled(t *testing.T) {
// 	ctx := context.Background()
// 	log := testr.New(t)
// 	imageName := "busybox:latest"
// 	mockClient := new(DockerContainerControllerMock)
// 	mockClient.On("ImageList", ctx, image.ListOptions{}).Return([]image.Summary{}, nil)
// 	mockClient.On("ImagePull", ctx, imageName, image.PullOptions{}).Return(io.NopCloser(strings.NewReader("pulling image")), nil)
//     cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// 	client := DockerContainerControllerMock{cli: cli}
//     mockClient.On("ImageRemove", ctx, imageName, image.RemoveOptions{Force: true}).Return([]image.DeleteResponse{}, nil)
//     mockDeleteImage(ctx, client.cli, imageName)
// 	err := pullImageIfNotExists(ctx, client.cli, imageName, log)
// 	assert.NoError(t, err)
// 	client.AssertExpectations(t)
// }

// func TestPullImageIfNotExists_ImageAlreadyPulled(t *testing.T) {
// 	ctx := context.Background()
// 	log := testr.New(t)
// 	imageName := "busybox:latest"

// 	mockClient := new(DockerContainerControllerMock)
// 	mockClient.On("ImageList", ctx, image.ListOptions{}).Return([]image.Summary{
// 		{RepoTags: []string{imageName}},
// 	}, nil)
    

//     cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// 	client := DockerContainerControllerMock{cli: cli}

//     if err := mockPullImage(ctx, client.cli, imageName); err != nil {
// 		fmt.Printf("Error pulling image: %v\n", err)
// 	} else {
// 		fmt.Println("Image pulled successfully.")
// 	}


// 	err := pullImageIfNotExists(ctx, client.cli, imageName, log)
// 	assert.NoError(t, err)
// 	client.AssertExpectations(t)
// }

// func TestPullImageIfNotExists_ImagePullError(t *testing.T) {
// 	ctx := context.Background()
// 	log := testr.New(t)
// 	imageName := "test-image:latest"

// 	mockClient := new(DockerContainerControllerMock)
// 	mockClient.On("ImageList", ctx, image.ListOptions{}).Return([]image.Summary{}, nil)
// 	mockClient.On("ImagePull", ctx, imageName, image.PullOptions{}).Return(nil, errors.New("failed to pull image"))

//     cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// 	client := DockerContainerControllerMock{cli: cli}

// 	err := pullImageIfNotExists(ctx, client.cli, imageName, log)
// 	assert.Error(t, err)
// 	client.AssertExpectations(t)
// }

// func TestPullImageIfNotExists_ImageListError(t *testing.T) {
// 	ctx := context.Background()
// 	log := testr.New(t)
// 	imageName := "test-image:latest"

//     cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// 	client := DockerContainerControllerMock{cli: cli}

// 	mockClient := new(DockerContainerControllerMock)
// 	mockClient.On("ImageList", ctx, image.ListOptions{}).Return(nil, errors.New("failed to list images"))

// 	err := pullImageIfNotExists(ctx, client.cli, imageName, log)
// 	assert.Error(t, err)
// 	client.AssertExpectations(t)
// }

// func TestPullImageIfNotExists_NilDockerClient(t *testing.T) {
// 	ctx := context.Background()
// 	log := testr.New(t)
// 	imageName := "test-image:latest"

// 	err := pullImageIfNotExists(ctx, nil, imageName, log)
// 	assert.Error(t, err)
// 	assert.Equal(t, "docker client is nil", err.Error())
// }

// func TestPullImageIfNotExists_NilContext(t *testing.T) {
// 	log := testr.New(t)
// 	imageName := "busybox:latest"

//     cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// 	client := DockerContainerControllerMock{cli: cli}

// 	err := pullImageIfNotExists(nil, client.cli, imageName, log)
// 	assert.Error(t, err)
// 	assert.Contains(t, "context is nil", err.Error())
// }
// func TestIsContainerRunning_ContainerRunning(t *testing.T) {
// 	ctx := context.Background()
// 	containerName := "running-container"
// 	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

// 	mockClient := new(DockerContainerControllerMock)

// 	mockClient.On("ContainerList", ctx, container.ListOptions{All: true}).Return([]types.Container{
// 		{
// 			Names: []string{"/" + containerName},
// 			State: "running",
// 		},
// 	}, nil)

// 	// client := DockerContainerControllerMock{cli: cli}

// 	err, isRunning := isContainerRunning(ctx, mockClient.cli, containerName)
// 	assert.NoError(t, err)
// 	assert.True(t, isRunning)
// 	mockClient.AssertExpectations(t)
// }

// func TestIsContainerRunning_ContainerNotRunning(t *testing.T) {
// 	ctx := context.Background()
// 	containerName := "stopped-container"
// 	mockClient := new(DockerContainerControllerMock)

// 	mockClient.On("ContainerList", ctx, container.ListOptions{All: true}).Return([]types.Container{
// 		{
// 			Names: []string{"/" + containerName},
// 			State: "exited",
// 		},
// 	}, nil)

// 	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// 	client := DockerContainerControllerMock{cli: cli}

// 	err, isRunning := isContainerRunning(ctx, client.cli, containerName)
// 	assert.NoError(t, err)
// 	assert.False(t, isRunning)
// 	client.AssertExpectations(t)
// }

// func TestIsContainerRunning_ContainerDoesNotExist(t *testing.T) {
// 	ctx := context.Background()
// 	containerName := "nonexistent-container"
// 	mockClient := new(DockerContainerControllerMock)

// 	mockClient.On("ContainerList", ctx, container.ListOptions{All: true}).Return([]types.Container{}, nil)

// 	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// 	client := DockerContainerControllerMock{cli: cli}

// 	err, isRunning := isContainerRunning(ctx, client.cli, containerName)
// 	assert.NoError(t, err)
// 	assert.False(t, isRunning)
// 	client.AssertExpectations(t)
// }

// // func TestIsContainerRunning_ErrorFetchingContainers(t *testing.T) {
// // 	ctx := context.Background()
// // 	containerName := "test-container"
// // 	mockClient := new(DockerContainerControllerMock)

// // 	mockClient.On("ContainerList", ctx, container.ListOptions{All: true}).Return(nil, errors.New("failed to list containers"))

// // 	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// // 	mockClient.cli = DockerContainerControllerMock{cli: cli}.cli

// // 	err, isRunning := isContainerRunning(ctx, mockClient.cli, containerName)
// // 	assert.Error(t, err)
// // 	assert.False(t, isRunning)
// // 	mockClient.AssertExpectations(t)
// // }
// //todo fix!!
// // func TestDoesContainerExist_ContainerExists(t *testing.T) {
// // 	ctx := context.Background()
// // 	containerName := "nmos-container"
// // 	mockClient := new(DockerContainerControllerMock)
// // 	mockClient.On("ContainerList", ctx, container.ListOptions{All: true}).Return([]containers.Container{
// // 		{
// // 			Names: []string{"/" + containerName},
// // 			State: "exited",
// // 		},
// // 	}, nil,)

// // 	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// // 	client := DockerContainerControllerMock{cli: cli}

// // 	err, exists := doesContainerExist(ctx, client.cli, containerName)
// // 	assert.NoError(t, err)
// // 	assert.True(t, exists)
// // 	client.AssertExpectations(t)
// // }

// // func TestDoesContainerExist_ContainerDoesNotExist(t *testing.T) {
// // 	ctx := context.Background()
// // 	containerName := "nonexistent-container"
// // 	mockClient := new(DockerContainerControllerMock)

// // 	mockClient.On("ContainerList", ctx, container.ListOptions{All: true}).Return([]container.Container{}, nil)

// // 	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// // 	client := DockerContainerControllerMock{cli: cli}

// // 	err, exists := doesContainerExist(ctx, client.cli, containerName)
// // 	assert.NoError(t, err)
// // 	assert.False(t, exists)
// // 	client.AssertExpectations(t)
// // }

// // func TestDoesContainerExist_ErrorFetchingContainers(t *testing.T) {
// // 	ctx := context.Background()
// // 	containerName := "test-container"
// // 	mockClient := new(DockerContainerControllerMock)

// // 	mockClient.On("ContainerList", ctx, container.ListOptions{All: true}).Return(nil, errors.New("failed to list containers"))

// // 	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// // 	client := DockerContainerControllerMock{cli: cli}

// // 	err, exists := doesContainerExist(ctx, client.cli, containerName)
// // 	assert.Error(t, err)
// // 	assert.False(t, exists)
// // 	client.AssertExpectations(t)
// // }







