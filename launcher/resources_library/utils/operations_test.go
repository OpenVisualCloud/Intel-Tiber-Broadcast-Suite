/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

package utils

import (
	"testing"

	"bcs.pod.launcher.intel/resources_library/resources/general"
	"bcs.pod.launcher.intel/resources_library/workloads"

	"github.com/docker/engine-api/types/strslice"
	"github.com/docker/go-connections/nat"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
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

	containerConfig, hostConfig, networkConfig := constructContainerConfig(containerInfo, log)

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

	containerConfig, hostConfig, networkConfig := constructContainerConfig(containerInfo, log)

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

	containerConfig, hostConfig, networkConfig := constructContainerConfig(containerInfo, log)

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
