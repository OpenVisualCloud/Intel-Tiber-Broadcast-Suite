//
//  SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
//
//  SPDX-License-Identifier: BSD-3-Clause
//

package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestNmosClientConfig_UnmarshalYAML(t *testing.T) {
	t.Run("Valid NmosClientConfig YAML", func(t *testing.T) {
		yamlContent := `
name: "nmos-client"
imageAndTag: "nmos-client:latest"
environmentVariables:
  - "ENV_VAR1=value1"
  - "ENV_VAR2=value2"
nmosConfigPath: "/etc/nmos"
nmosConfigFileName: "config.yaml"
network:
  enable: true
  name: "nmos-network"
  ip: "192.168.1.100"
ffmpegConectionAddress: "192.168.1.101"
ffmpegConnectionPort: "8080"
`
		var config NmosClientConfig
		err := yaml.Unmarshal([]byte(yamlContent), &config)
		assert.NoError(t, err)

		assert.Equal(t, "nmos-client", config.Name)
		assert.Equal(t, "nmos-client:latest", config.ImageAndTag)
		assert.Equal(t, []string{"ENV_VAR1=value1", "ENV_VAR2=value2"}, config.EnvironmentVariables)
		assert.Equal(t, "/etc/nmos", config.NmosConfigPath)
		assert.Equal(t, "config.yaml", config.NmosConfigFileName)
		assert.True(t, config.Network.Enable)
		assert.Equal(t, "nmos-network", config.Network.Name)
		assert.Equal(t, "192.168.1.100", config.Network.IP)
		assert.Equal(t, "192.168.1.101", config.FfmpegConectionAddress)
		assert.Equal(t, "8080", config.FfmpegConnectionPort)
	})

	t.Run("Invalid NmosClientConfig YAML", func(t *testing.T) {
		yamlContent := `
name: "nmos-client"
imageAndTag: "nmos-client:latest"
environmentVariables: "not-an-array"
`
		var config NmosClientConfig
		err := yaml.Unmarshal([]byte(yamlContent), &config)
		assert.Error(t, err)
	})
}

func TestNmosClientConfig_DefaultValues(t *testing.T) {
	t.Run("Default values for missing fields", func(t *testing.T) {
		yamlContent := `
name: "nmos-client"
imageAndTag: "nmos-client:latest"
`
		var config NmosClientConfig
		err := yaml.Unmarshal([]byte(yamlContent), &config)
		assert.NoError(t, err)

		assert.Equal(t, "nmos-client", config.Name)
		assert.Equal(t, "nmos-client:latest", config.ImageAndTag)
		assert.Nil(t, config.EnvironmentVariables)
		assert.Empty(t, config.NmosConfigPath)
		assert.Empty(t, config.NmosConfigFileName)
		assert.False(t, config.Network.Enable)
		assert.Empty(t, config.Network.Name)
		assert.Empty(t, config.Network.IP)
		assert.Empty(t, config.FfmpegConectionAddress)
		assert.Empty(t, config.FfmpegConnectionPort)
	})
}
func TestParseLauncherMode(t *testing.T) {
	t.Run("Valid YAML with ModeK8s true", func(t *testing.T) {
		yamlContent := `
    k8s: true
    definition:
      meshAgent:
        image: "mesh-agent:latest"
        restPort: 8100
        grpcPort: 50051
      mediaProxy:
        image: mcm/media-proxy:latest
        command: ["media-proxy"]
        args: ["-d", "kernel:eth0", "-i", "$(POD_IP)"]
        grpcPort: 8001
        sdkPort: 8002
        volumes:
          memif: /tmp/mcm/memif
          vfio: /dev/vfio
        pvHostPath: /var/run/imtl
        pvStorageClass: manual
        pvStorage: 1Gi
        pvcStorage: 1Gi
`
		tempFile, err := os.CreateTemp("", "valid_modek8s_true_*.yaml")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		_, err = tempFile.Write([]byte(yamlContent))
		assert.NoError(t, err)
		tempFile.Close()

		modeK8s, err := ParseLauncherMode(tempFile.Name())
		assert.NoError(t, err)
		assert.True(t, modeK8s)
	})

	t.Run("Valid YAML with ModeK8s false", func(t *testing.T) {
		yamlContent := `
    k8s: false
    definition:
      meshAgent:
        image: "mesh-agent:latest"
        restPort: 8100
        grpcPort: 50051
      mediaProxy:
        image: mcm/media-proxy:latest
        command: ["media-proxy"]
        args: ["-d", "kernel:eth0", "-i", "$(POD_IP)"]
        grpcPort: 8001
        sdkPort: 8002
        volumes:
          memif: /tmp/mcm/memif
          vfio: /dev/vfio
        pvHostPath: /var/run/imtl
        pvStorageClass: manual
        pvStorage: 1Gi
        pvcStorage: 1Gi
`
		tempFile, err := os.CreateTemp("", "valid_modek8s_false_*.yaml")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		_, err = tempFile.Write([]byte(yamlContent))
		assert.NoError(t, err)
		tempFile.Close()

		modeK8s, err := ParseLauncherMode(tempFile.Name())
		assert.NoError(t, err)
		assert.False(t, modeK8s)
	})

	t.Run("Invalid YAML", func(t *testing.T) {
		yamlContent := `
k8s: true
configuration:
  runOnce:
	mediaProxyAgent:
	  imageAndTag: "agent:latest"
	mediaProxyMcm:
	  imageAndTag: "mcm:latest"
  workloadToBeRun:
	ffmpegPipeline:
	  name: "pipeline"
	nmosClient:
	  name: "client
`
		tempFile, err := os.CreateTemp("", "invalid_yaml_*.yaml")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		_, err = tempFile.Write([]byte(yamlContent))
		assert.NoError(t, err)
		tempFile.Close()

		_, err = ParseLauncherMode(tempFile.Name())
		assert.Error(t, err)
	})

	t.Run("File does not exist", func(t *testing.T) {
		_, err := ParseLauncherMode("non_existent_file.yaml")
		assert.Error(t, err)
	})
}

// func TestParseLauncherConfiguration(t *testing.T) {
// 	t.Run("Valid YAML with full configuration", func(t *testing.T) {
// 		yamlContent := `
// k8s: true
// configuration:
//   runOnce:
//   mediaProxyAgent:
//     imageAndTag: agent:latest
//     gRPCPort: "50051"
//     restPort: "8100"
//     network:
//     enable: true
//     name: "agent-network"
//     ip: "192.168.1.10"
//   mediaProxyMcm:
//     imageAndTag: "mcm:latest"
//     interfaceName: "eth0"
//     volumes:
//     - "/data"
//     network:
//     enable: true
//     name: "mcm-network"
//     ip: "192.168.1.20"
//   workloadToBeRun:
//   ffmpegPipeline:
//     name: "pipeline"
//     imageAndTag: "pipeline:latest"
//     gRPCPort: 9000
//     nmosPort: 8000
//     environmentVariables:
//     - "ENV_VAR1=value1"
//     - "ENV_VAR2=value2"
//     volumes:
//     videos: "/videos"
//     dri: "/dri"
//     kahawai: "/kahawai"
//     devnull: "/dev/null"
//     tmpHugepages: "/hugepages/tmp"
//     hugepages: "/hugepages"
//     imtl: "/imtl"
//     shm: "/shm"
//     devices:
//     vfio: "/dev/vfio"
//     dri: "/dev/dri"
//     network:
//     enable: true
//     name: "pipeline-network"
//     ip: "192.168.1.30"
//   nmosClient:
//     name: "nmos-client"
//     imageAndTag: "nmos-client:latest"
//     environmentVariables:
//     - "ENV_VAR3=value3"
//     - "ENV_VAR4=value4"
//     nmosConfigPath: "/etc/nmos"
//     nmosConfigFileName: "config.yaml"
//     network:
//     enable: true
//     name: "nmos-network"
//     ip: "192.168.1.40"
//     ffmpegConectionAddress: "192.168.1.50"
//     ffmpegConnectionPort: "8080"
// `
// 		tempFile, err := os.CreateTemp("", "valid_configuration_*.yaml")
// 		assert.NoError(t, err)
// 		defer os.Remove(tempFile.Name())

// 		_, err = tempFile.Write([]byte(yamlContent))
// 		assert.NoError(t, err)
// 		tempFile.Close()

// 		config, err := ParseLauncherConfiguration(tempFile.Name())
// 		fmt.Printf("Parsed config: %+v\n", config)
// 		assert.NoError(t, err)

// 		assert.True(t, config.RunOnce.MediaProxyAgent.Network.Enable)
// 		assert.Equal(t, "agent:latest", config.RunOnce.MediaProxyAgent.ImageAndTag)
// 		assert.Equal(t, "192.168.1.10", config.RunOnce.MediaProxyAgent.Network.IP)

// 		assert.True(t, config.RunOnce.MediaProxyMcm.Network.Enable)
// 		assert.Equal(t, "mcm:latest", config.RunOnce.MediaProxyMcm.ImageAndTag)
// 		assert.Equal(t, "192.168.1.20", config.RunOnce.MediaProxyMcm.Network.IP)

// 		assert.Equal(t, "pipeline", config.WorkloadToBeRun.FfmpegPipeline.Name)
// 		assert.Equal(t, "pipeline:latest", config.WorkloadToBeRun.FfmpegPipeline.ImageAndTag)
// 		assert.Equal(t, 9000, config.WorkloadToBeRun.FfmpegPipeline.GRPCPort)
// 		assert.Equal(t, "/videos", config.WorkloadToBeRun.FfmpegPipeline.Volumes.Videos)

// 		assert.Equal(t, "nmos-client", config.WorkloadToBeRun.NmosClient.Name)
// 		assert.Equal(t, "nmos-client:latest", config.WorkloadToBeRun.NmosClient.ImageAndTag)
// 		assert.Equal(t, "192.168.1.50", config.WorkloadToBeRun.NmosClient.FfmpegConectionAddress)
// 	})

// 	t.Run("Invalid YAML", func(t *testing.T) {
// 		yamlContent := `
// k8s: true
// configuration:
//   runOnce:
// 	mediaProxyAgent:
// 	  imageAndTag: "agent:latest"
// 	  gRPCPort: "50051"
// 	  restPort: "8100"
// 	  network:
// 		enable: true
// 		name: "agent-network"
// 		ip: "192.168.1.10"
//   workloadToBeRun:
// 	ffmpegPipeline:
// 	  name: "pipeline
// `
// 		tempFile, err := os.CreateTemp("", "invalid_configuration_*.yaml")
// 		assert.NoError(t, err)
// 		defer os.Remove(tempFile.Name())

// 		_, err = tempFile.Write([]byte(yamlContent))
// 		assert.NoError(t, err)
// 		tempFile.Close()

// 		_, err = ParseLauncherConfiguration(tempFile.Name())
// 		assert.Error(t, err)
// 	})

// 	t.Run("File does not exist", func(t *testing.T) {
// 		_, err := ParseLauncherConfiguration("non_existent_file.yaml")
// 		assert.Error(t, err)
// 	})
// }