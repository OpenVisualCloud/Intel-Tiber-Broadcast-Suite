/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

package utils

import (
	"os"
	"testing"

	bcsv1 "bcs.pod.launcher.intel/api/v1"

	"bcs.pod.launcher.intel/resources_library/resources/general"
	"bcs.pod.launcher.intel/resources_library/workloads"
	corev1 "k8s.io/api/core/v1"

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
					ImageAndTag:          "ffmpeg-image:latest",
					GRPCPort:             50051,
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

func TestUpdateNmosJsonFile_Success(t *testing.T) {
	tempFile, err := os.CreateTemp("", "nmos_config_*.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())
	initialJson := `{
		"domain":"",
		"ffmpeg_grpc_server_address":"192.9.1.200",
		"ffmpeg_grpc_server_port":"6001",
		"function":"",
		"gpu_hw_acceleration":"",
		"http_port":0,
		"receiver": [],
		"receivers": [],
		"receivers_count": [],
		"sender": [],
		"sender_payload_type":0, 
		"senders":[],
		"senders_count":[]
		}`
	_, err = tempFile.WriteString(initialJson)
	assert.NoError(t, err)
	tempFile.Close()

	err = updateNmosJsonFile(tempFile.Name(), "192.168.1.200", "60051")
	assert.NoError(t, err)

	updatedContent, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err)

	expectedJson := `{
		"domain":"",
		"ffmpeg_grpc_server_address":"192.168.1.200",
		"ffmpeg_grpc_server_port":"60051",
		"function":"",
		"gpu_hw_acceleration":"",
		"http_port":0,
		"receiver": [],
		"receivers": [],
		"receivers_count": [],
		"sender": [],
		"sender_payload_type":0, 
		"senders":[],
		"senders_count":[]
		}`
	assert.JSONEq(t, expectedJson, string(updatedContent))
}

func TestUpdateNmosJsonFile_FileNotFound(t *testing.T) {
	err := updateNmosJsonFile("non_existent_file.json", "192.168.1.200", "60051")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "open non_existent_file.json: no such file or directory")
}
func TestCreateNamespace(t *testing.T) {
	namespaceName := "test-namespace"
	namespace := CreateNamespace(namespaceName)

	assert.NotNil(t, namespace)
	assert.Equal(t, namespaceName, namespace.ObjectMeta.Name)
}

func TestConvertEnvVars(t *testing.T) {
	input := []bcsv1.EnvVar{
		{Name: "ENV_VAR1", Value: "VALUE1"},
		{Name: "ENV_VAR2", Value: "VALUE2"},
	}

	expected := []corev1.EnvVar{
		{Name: "ENV_VAR1", Value: "VALUE1"},
		{Name: "ENV_VAR2", Value: "VALUE2"},
	}

	result := convertEnvVars(input)

	assert.Equal(t, expected, result)
}

func TestConvertEnvVars_EmptyInput(t *testing.T) {
	input := []bcsv1.EnvVar{}

	expected := []corev1.EnvVar(nil)

	result := convertEnvVars(input)

	assert.Equal(t, expected, result)
}

func TestConvertEnvVars_NilInput(t *testing.T) {
	var input []bcsv1.EnvVar

	expected := []corev1.EnvVar(nil)

	result := convertEnvVars(input)

	assert.Equal(t, expected, result)
}

// func TestCreateBcsDeployment(t *testing.T) {
// 	bcsConfig := &bcsv1.BcsConfig{
// 		Spec: bcsv1.BcsConfigSpec{
// 			Name:      "test-bcs",
// 			Namespace: "test-namespace",
// 			Nmos: bcsv1.Nmos{
// 				Image: "nmos-image:latest",
// 				Args:  []string{"--arg1", "--arg2"},
// 				EnvironmentVariables: []bcsv1.EnvVar{
// 					{Name: "ENV_VAR1", Value: "VALUE1"},
// 					{Name: "ENV_VAR2", Value: "VALUE2"},
// 				},
// 			},
// 			App: bcsv1.App{
// 				Image:    "app-image:latest",
// 				GrpcPort: 50051,
// 				Volumes: map[string]string{
// 					"videos":       "/host/videos",
// 					"dri":          "/host/dri",
// 					"kahawaiLock":  "/host/kahawai.lock",
// 					"devNull":      "/dev/null",
// 					"hugepagesTmp": "/host/hugepages/tmp",
// 					"hugepages":    "/host/hugepages",
// 					"imtl":         "/host/imtl",
// 					"shm":          "/dev/shm",
// 					"vfio":         "/dev/vfio",
// 					"dri-dev":       "/dev/dri",
// 				},
// 				EnvironmentVariables: []bcsv1.EnvVar{
// 					{Name: "APP_ENV_VAR1", Value: "APP_VALUE1"},
// 					{Name: "APP_ENV_VAR2", Value: "APP_VALUE2"},
// 				},
// 			},
// 		},
// 	}

// 	deployment := CreateBcsDeployment(bcsConfig)

// 	assert.NotNil(t, deployment)
// 	assert.Equal(t, "test-bcs", deployment.ObjectMeta.Name)
// 	assert.Equal(t, "test-namespace", deployment.ObjectMeta.Namespace)

// 	assert.Equal(t, 2, len(deployment.Spec.Template.Spec.Containers))

// 	nmosContainer := deployment.Spec.Template.Spec.Containers[0]
// 	assert.Equal(t, "tiber-broadcast-suite-nmos-node", nmosContainer.Name)
// 	assert.Equal(t, "nmos-image:latest", nmosContainer.Image)
// 	assert.ElementsMatch(t, []string{"--arg1", "--arg2"}, nmosContainer.Args)
// 	assert.ElementsMatch(t, []corev1.EnvVar{
// 		{Name: "ENV_VAR1", Value: "VALUE1"},
// 		{Name: "ENV_VAR2", Value: "VALUE2"},
// 	}, nmosContainer.Env)

// 	appContainer := deployment.Spec.Template.Spec.Containers[1]
// 	assert.Equal(t, "tiber-broadcast-suite", appContainer.Name)
// 	assert.Equal(t, "app-image:latest", appContainer.Image)
// 	assert.ElementsMatch(t, []string{"localhost", "50051"}, appContainer.Args)
// 	assert.ElementsMatch(t, []corev1.EnvVar{
// 		{Name: "APP_ENV_VAR1", Value: "APP_VALUE1"},
// 		{Name: "APP_ENV_VAR2", Value: "APP_VALUE2"},
// 	}, appContainer.Env)

//		assert.Equal(t, 11, len(deployment.Spec.Template.Spec.Volumes))
//		assert.Equal(t, "/host/videos", deployment.Spec.Template.Spec.Volumes[0].VolumeSource.HostPath.Path)
//		assert.Equal(t, "/host/dri", deployment.Spec.Template.Spec.Volumes[1].VolumeSource.HostPath.Path)
//		assert.Equal(t, "/host/kahawai.lock", deployment.Spec.Template.Spec.Volumes[2].VolumeSource.HostPath.Path)
//		assert.Equal(t, "/dev/null", deployment.Spec.Template.Spec.Volumes[3].VolumeSource.HostPath.Path)
//		assert.Equal(t, "/host/hugepages/tmp", deployment.Spec.Template.Spec.Volumes[4].VolumeSource.HostPath.Path)
//		assert.Equal(t, "/host/hugepages", deployment.Spec.Template.Spec.Volumes[5].VolumeSource.HostPath.Path)
//		assert.Equal(t, "/host/imtl", deployment.Spec.Template.Spec.Volumes[6].VolumeSource.HostPath.Path)
//		assert.Equal(t, "/dev/shm", deployment.Spec.Template.Spec.Volumes[7].VolumeSource.HostPath.Path)
//		assert.Equal(t, "/dev/vfio", deployment.Spec.Template.Spec.Volumes[8].VolumeSource.HostPath.Path)
//		assert.Equal(t, "/dev/dri", deployment.Spec.Template.Spec.Volumes[9].VolumeSource.HostPath.Path)
//	}
func TestCreateBcsService(t *testing.T) {
	bcsConfig := &bcsv1.BcsConfig{
		Spec: bcsv1.BcsConfigSpec{
			Name:      "test-bcs-service",
			Namespace: "test-namespace",
			Nmos: bcsv1.Nmos{
				NmosApiNodePort: 30080,
			},
		},
	}

	service := CreateBcsService(bcsConfig)

	assert.NotNil(t, service)
	assert.Equal(t, "test-bcs-service", service.ObjectMeta.Name)
	assert.Equal(t, "test-namespace", service.ObjectMeta.Namespace)
	assert.Equal(t, corev1.ServiceTypeNodePort, service.Spec.Type)

	assert.Equal(t, 2, len(service.Spec.Ports))

	assert.Equal(t, map[string]string{"app": "test-bcs-service"}, service.Spec.Selector)
}
func TestCreateService(t *testing.T) {
	serviceName := "test-service"
	service := CreateService(serviceName)

	assert.NotNil(t, service)
	assert.Equal(t, serviceName, service.ObjectMeta.Name)
	assert.Equal(t, "default", service.ObjectMeta.Namespace)
	assert.Equal(t, map[string]string{"app": serviceName}, service.Spec.Selector)

	assert.Equal(t, 1, len(service.Spec.Ports))
	port := service.Spec.Ports[0]
	assert.Equal(t, corev1.ProtocolTCP, port.Protocol)
	assert.Equal(t, int32(80), port.Port)
}
