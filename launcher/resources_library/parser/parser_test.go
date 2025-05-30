package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLauncherMode(t *testing.T) {
	// Create a temporary YAML file for testing
	yamlData := `
k8s: false
configuration:
  runOnce:
    mediaProxyAgent:
      imageAndTag: mcm/mesh-agent:latest
      gRPCPort: 50051
      restPort: 8100
      custom_network:
        enable: false
    mediaProxyMcm:
      imageAndTag: mcm/media-proxy:latest
      interfaceName: eth0
      volumes:
        - /dev/vfio:/dev/vfio
      custom_network:
        enable: false
  workloadToBeRun:
    - ffmpegPipeline:
        name: bcs-ffmpeg-pipeline-tx
        imageAndTag: tiber-broadcast-suite:latest
        gRPCPort: 50088
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
        custom_network:
          enable: false
          ip: 10.123.x.x
      nmosClient:
        name: bcs-ffmpeg-pipeline-nmos-client-tx
        imageAndTag: tiber-broadcast-suite-nmos-node:latest
        environmentVariables:
          - "http_proxy="
          - "https_proxy=" 
          - "VFIO_PORT_TX=0000:ca:11.0"
        nmosConfigPath: /root/path/to/intel-node-tx/json/file
        nmosPort: 5045
        nmosConfigFileName: intel-node-tx.json
        custom_network:
          enable: false
`
	tempFile, err := os.CreateTemp("", "launcher_mode_test_*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte(yamlData))
	assert.NoError(t, err)
	tempFile.Close()

	// Test ParseLauncherMode
	modeK8s, err := ParseLauncherMode(tempFile.Name())
	assert.NoError(t, err)
	assert.False(t, modeK8s)
}

func TestParseLauncherConfiguration(t *testing.T) {
	// Create a temporary YAML file for testing
	yamlData := `
k8s: false
configuration:
  runOnce:
    mediaProxyAgent:
      imageAndTag: mcm/mesh-agent:latest
      gRPCPort: 50051
      restPort: 8100
      custom_network:
        enable: false
    mediaProxyMcm:
      imageAndTag: mcm/media-proxy:latest
      interfaceName: eth0
      volumes:
        - /dev/vfio:/dev/vfio
      custom_network:
        enable: false
  workloadToBeRun:
    - ffmpegPipeline:
        name: bcs-ffmpeg-pipeline-tx
        imageAndTag: tiber-broadcast-suite:latest
        gRPCPort: 50088
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
        custom_network:
          enable: false
          ip: 10.123.1.1
      nmosClient:
        name: bcs-ffmpeg-pipeline-nmos-client-tx
        imageAndTag: tiber-broadcast-suite-nmos-node:latest
        environmentVariables:
          - "http_proxy="
          - "https_proxy=" 
          - "VFIO_PORT_TX=0000:ca:11.0"
        nmosConfigPath: /root/path/to/intel-node-tx/json/file
        nmosPort: 5045
        nmosConfigFileName: intel-node-tx.json
        custom_network:
          enable: false
`
	tempFile, err := os.CreateTemp("", "launcher_config_test_*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte(yamlData))
	assert.NoError(t, err)
	tempFile.Close()

	// Test ParseLauncherConfiguration
	config, err := ParseLauncherConfiguration(tempFile.Name())
	assert.NoError(t, err)

	// Validate RunOnce
	assert.Equal(t, "mcm/mesh-agent:latest", config.RunOnce.MediaProxyAgent.ImageAndTag)
	assert.Equal(t, "50051", config.RunOnce.MediaProxyAgent.GRPCPort)
	assert.Equal(t, "8100", config.RunOnce.MediaProxyAgent.RestPort)
	assert.False(t, config.RunOnce.MediaProxyAgent.Network.Enable)
	assert.Equal(t, "mcm/media-proxy:latest", config.RunOnce.MediaProxyMcm.ImageAndTag)
	assert.Equal(t, "eth0", config.RunOnce.MediaProxyMcm.InterfaceName)
	assert.Len(t, config.RunOnce.MediaProxyMcm.Volumes, 1)
	assert.Equal(t, "/dev/vfio:/dev/vfio", config.RunOnce.MediaProxyMcm.Volumes[0])
	assert.False(t, config.RunOnce.MediaProxyMcm.Network.Enable)

	// Validate WorkloadToBeRun
	assert.Len(t, config.WorkloadToBeRun, 1)
	assert.Equal(t, "bcs-ffmpeg-pipeline-tx", config.WorkloadToBeRun[0].FfmpegPipeline.Name)
	assert.Equal(t, "tiber-broadcast-suite:latest", config.WorkloadToBeRun[0].FfmpegPipeline.ImageAndTag)
	assert.Equal(t, 50088, config.WorkloadToBeRun[0].FfmpegPipeline.GRPCPort)
	assert.Equal(t, []string{"http_proxy=", "https_proxy="}, config.WorkloadToBeRun[0].FfmpegPipeline.EnvironmentVariables)
	assert.Equal(t, "/root", config.WorkloadToBeRun[0].FfmpegPipeline.Volumes.Videos)
	assert.Equal(t, "/usr/lib/x86_64-linux-gnu/dri", config.WorkloadToBeRun[0].FfmpegPipeline.Volumes.Dri)
	assert.Equal(t, "/tmp/kahawai_lcore.lock", config.WorkloadToBeRun[0].FfmpegPipeline.Volumes.Kahawai)
	assert.Equal(t, "/dev/null", config.WorkloadToBeRun[0].FfmpegPipeline.Volumes.Devnull)
	assert.Equal(t, "/tmp/hugepages", config.WorkloadToBeRun[0].FfmpegPipeline.Volumes.TmpHugepages)
	assert.Equal(t, "/hugepages", config.WorkloadToBeRun[0].FfmpegPipeline.Volumes.Hugepages)
	assert.Equal(t, "/var/run/imtl", config.WorkloadToBeRun[0].FfmpegPipeline.Volumes.Imtl)
	assert.Equal(t, "/dev/shm", config.WorkloadToBeRun[0].FfmpegPipeline.Volumes.Shm)
	assert.Equal(t, "/dev/vfio", config.WorkloadToBeRun[0].FfmpegPipeline.Devices.Vfio)
	assert.Equal(t, "/dev/dri", config.WorkloadToBeRun[0].FfmpegPipeline.Devices.Dri)
	assert.False(t, config.WorkloadToBeRun[0].FfmpegPipeline.Network.Enable)
	assert.Equal(t, "10.123.1.1", config.WorkloadToBeRun[0].FfmpegPipeline.Network.IP)

	assert.Equal(t, "bcs-ffmpeg-pipeline-nmos-client-tx", config.WorkloadToBeRun[0].NmosClient.Name)
	assert.Equal(t, "tiber-broadcast-suite-nmos-node:latest", config.WorkloadToBeRun[0].NmosClient.ImageAndTag)
	assert.Equal(t, []string{"http_proxy=", "https_proxy=", "VFIO_PORT_TX=0000:ca:11.0"}, config.WorkloadToBeRun[0].NmosClient.EnvironmentVariables)
	assert.Equal(t, "/root/path/to/intel-node-tx/json/file", config.WorkloadToBeRun[0].NmosClient.NmosConfigPath)
	assert.Equal(t, 5045, config.WorkloadToBeRun[0].NmosClient.NmosPort)
	assert.Equal(t, "intel-node-tx.json", config.WorkloadToBeRun[0].NmosClient.NmosConfigFileName)
	assert.False(t, config.WorkloadToBeRun[0].NmosClient.Network.Enable)

}

func TestParseLauncherConfigurationWithMultipleWorkloads(t *testing.T) {
	// Create a temporary YAML file for testing
	yamlData := `
k8s: true
configuration:
  runOnce:
    mediaProxyAgent:
      imageAndTag: mcm/mesh-agent:latest
      gRPCPort: 50051
      restPort: 8100
      custom_network:
        enable: true
    mediaProxyMcm:
      imageAndTag: mcm/media-proxy:latest
      interfaceName: eth0
      volumes:
        - /dev/vfio:/dev/vfio
      custom_network:
        enable: true
  workloadToBeRun:
    - ffmpegPipeline:
        name: bcs-ffmpeg-pipeline-tx-1
        imageAndTag: tiber-broadcast-suite:latest
        gRPCPort: 50088
        environmentVariables:
          - "http_proxy="
          - "https_proxy="
        volumes:
          videos: /root/videos1
          dri: /usr/lib/x86_64-linux-gnu/dri1
        devices:
          vfio: /dev/vfio1
          dri: /dev/dri1
        custom_network:
          enable: true
          ip: 10.123.1.1
      nmosClient:
        name: bcs-ffmpeg-pipeline-nmos-client-tx-1
        imageAndTag: tiber-broadcast-suite-nmos-node:latest
        environmentVariables:
          - "http_proxy="
          - "https_proxy="
        nmosConfigPath: /root/path/to/intel-node-tx/json/file1
        nmosPort: 5045
        nmosConfigFileName: intel-node-tx-1.json
        custom_network:
          enable: true
    - ffmpegPipeline:
        name: bcs-ffmpeg-pipeline-tx-2
        imageAndTag: tiber-broadcast-suite:latest
        gRPCPort: 50089
        environmentVariables:
          - "http_proxy="
          - "https_proxy="
        volumes:
          videos: /root/videos2
          dri: /usr/lib/x86_64-linux-gnu/dri2
        devices:
          vfio: /dev/vfio2
          dri: /dev/dri2
        custom_network:
          enable: true
          ip: 10.123.1.2
      nmosClient:
        name: bcs-ffmpeg-pipeline-nmos-client-tx-2
        imageAndTag: tiber-broadcast-suite-nmos-node:latest
        environmentVariables:
          - "http_proxy="
          - "https_proxy="
        nmosConfigPath: /root/path/to/intel-node-tx/json/file2
        nmosPort: 5046
        nmosConfigFileName: intel-node-tx-2.json
        custom_network:
          enable: true
`
	tempFile, err := os.CreateTemp("", "launcher_config_test_multiple_*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte(yamlData))
	assert.NoError(t, err)
	tempFile.Close()

	// Test ParseLauncherConfiguration
	config, err := ParseLauncherConfiguration(tempFile.Name())
	assert.NoError(t, err)

	// Validate WorkloadToBeRun
	assert.Len(t, config.WorkloadToBeRun, 2)

	// Validate first workload
	assert.Equal(t, "bcs-ffmpeg-pipeline-tx-1", config.WorkloadToBeRun[0].FfmpegPipeline.Name)
	assert.Equal(t, "tiber-broadcast-suite:latest", config.WorkloadToBeRun[0].FfmpegPipeline.ImageAndTag)
	assert.Equal(t, 50088, config.WorkloadToBeRun[0].FfmpegPipeline.GRPCPort)
	assert.Equal(t, "/root/videos1", config.WorkloadToBeRun[0].FfmpegPipeline.Volumes.Videos)
	assert.Equal(t, "/usr/lib/x86_64-linux-gnu/dri1", config.WorkloadToBeRun[0].FfmpegPipeline.Volumes.Dri)
	assert.Equal(t, "/dev/vfio1", config.WorkloadToBeRun[0].FfmpegPipeline.Devices.Vfio)
	assert.Equal(t, "/dev/dri1", config.WorkloadToBeRun[0].FfmpegPipeline.Devices.Dri)
	assert.True(t, config.WorkloadToBeRun[0].FfmpegPipeline.Network.Enable)
	assert.Equal(t, "10.123.1.1", config.WorkloadToBeRun[0].FfmpegPipeline.Network.IP)

	assert.Equal(t, "bcs-ffmpeg-pipeline-nmos-client-tx-1", config.WorkloadToBeRun[0].NmosClient.Name)
	assert.Equal(t, "tiber-broadcast-suite-nmos-node:latest", config.WorkloadToBeRun[0].NmosClient.ImageAndTag)
	assert.Equal(t, "/root/path/to/intel-node-tx/json/file1", config.WorkloadToBeRun[0].NmosClient.NmosConfigPath)
	assert.Equal(t, 5045, config.WorkloadToBeRun[0].NmosClient.NmosPort)
	assert.Equal(t, "intel-node-tx-1.json", config.WorkloadToBeRun[0].NmosClient.NmosConfigFileName)
	assert.True(t, config.WorkloadToBeRun[0].NmosClient.Network.Enable)

	// Validate second workload
	assert.Equal(t, "bcs-ffmpeg-pipeline-tx-2", config.WorkloadToBeRun[1].FfmpegPipeline.Name)
	assert.Equal(t, "tiber-broadcast-suite:latest", config.WorkloadToBeRun[1].FfmpegPipeline.ImageAndTag)
	assert.Equal(t, 50089, config.WorkloadToBeRun[1].FfmpegPipeline.GRPCPort)
	assert.Equal(t, "/root/videos2", config.WorkloadToBeRun[1].FfmpegPipeline.Volumes.Videos)
	assert.Equal(t, "/usr/lib/x86_64-linux-gnu/dri2", config.WorkloadToBeRun[1].FfmpegPipeline.Volumes.Dri)
	assert.Equal(t, "/dev/vfio2", config.WorkloadToBeRun[1].FfmpegPipeline.Devices.Vfio)
	assert.Equal(t, "/dev/dri2", config.WorkloadToBeRun[1].FfmpegPipeline.Devices.Dri)
	assert.True(t, config.WorkloadToBeRun[1].FfmpegPipeline.Network.Enable)
	assert.Equal(t, "10.123.1.2", config.WorkloadToBeRun[1].FfmpegPipeline.Network.IP)

	assert.Equal(t, "bcs-ffmpeg-pipeline-nmos-client-tx-2", config.WorkloadToBeRun[1].NmosClient.Name)
	assert.Equal(t, "tiber-broadcast-suite-nmos-node:latest", config.WorkloadToBeRun[1].NmosClient.ImageAndTag)
	assert.Equal(t, "/root/path/to/intel-node-tx/json/file2", config.WorkloadToBeRun[1].NmosClient.NmosConfigPath)
	assert.Equal(t, 5046, config.WorkloadToBeRun[1].NmosClient.NmosPort)
	assert.Equal(t, "intel-node-tx-2.json", config.WorkloadToBeRun[1].NmosClient.NmosConfigFileName)
	assert.True(t, config.WorkloadToBeRun[1].NmosClient.Network.Enable)
}
