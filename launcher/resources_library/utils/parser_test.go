package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLauncherConfiguration(t *testing.T) {
    mockConfig := `
k8s: true
configuration:
  runOnce:
    mediaProxyAgent:
      imageAndTag: "mesh-agent:25.04"
      gRPCPort: "50051"
      restPort: "80"
      network:
        enable: true
        name: "host"
        ip: "localhost"
    mediaProxyMcm:
      imageAndTag: "mcm:25.04"
      interfaceName: "eth1"
      volumes:
        - "/dev/vfio:/dev/vfio"
      network:
        enable: true
        name: "host"
        ip: "localhost"
  workloadToBeRun:
    ffmpegPipeline:
      name: "ffmpeg"
      imageAndTag: "ffmpeg:latest"
      gRPCPort: 58000
      sourcePort: 5004
      environmentVariables:
        - "ENV_VAR_1=value1"
        - "ENV_VAR_2=value2"
      volumes:
        videos: "/videos"
        dri: "/dev/dri"
        kahawai: "/kahawai"
        devnull: "/dev/null"
        tmpHugepages: "/tmp/hugepages"
        hugepages: "/hugepages"
        imtl: "/imtl"
        shm: "/dev/shm"
      devices:
        vfio: "/dev/vfio"
        dri: "/dev/dri"
      network:
        enable: true
        name: "host"
        ip: "localhost"
    nmosClient:
      name: "nmos-client"
      imageAndTag: "nmos-client:latest"
      environmentVariables:
        - "ENV_VAR_1=value1"
        - "ENV_VAR_2=value2"
      nmosConfigPath: "/path/to/nmos/config"
      nmosConfigFileName: "nmos_config.json"
      network:
        enable: true
        name: "host"
        ip: "localhost"
      ffmpegConectionAddress: "localhost"
      ffmpegConnectionPort: "58000"
`
    tempFile, err := os.CreateTemp("", "config-*.yaml")
    assert.NoError(t, err)
    defer os.Remove(tempFile.Name())

    _, err = tempFile.Write([]byte(mockConfig))
    assert.NoError(t, err)
    tempFile.Close()

    // Act: Parse the configuration using the function
    config, err := ParseLauncherConfiguration(tempFile.Name())

    // Assert: Verify the parsed configuration matches the expected values
    assert.NoError(t, err)
    assert.True(t, config.RunOnce.MediaProxyAgent.Network.Enable)
    assert.Equal(t, "localhost",config.RunOnce.MediaProxyAgent.Network.IP)
    assert.Equal(t, "host",config.RunOnce.MediaProxyAgent.Network.Name)
    assert.Equal(t, "mesh-agent:25.04", config.RunOnce.MediaProxyAgent.ImageAndTag)
    assert.Equal(t, "50051", config.RunOnce.MediaProxyAgent.GRPCPort)
    assert.Equal(t, "80", config.RunOnce.MediaProxyAgent.RestPort)

    assert.Equal(t, "mcm:25.04", config.RunOnce.MediaProxyMcm.ImageAndTag)
    assert.Equal(t, "eth1", config.RunOnce.MediaProxyMcm.InterfaceName)
    assert.Equal(t, []string{"/dev/vfio:/dev/vfio"}, config.RunOnce.MediaProxyMcm.Volumes)
    assert.Contains(t, []string{"/dev/vfio:/dev/vfio"}[0], ":")


    assert.Equal(t, "ffmpeg", config.WorkloadToBeRun.FfmpegPipeline.Name)
    assert.Equal(t, "ffmpeg:latest", config.WorkloadToBeRun.FfmpegPipeline.ImageAndTag)
    assert.Equal(t, 58000, config.WorkloadToBeRun.FfmpegPipeline.GRPCPort)
    assert.Equal(t, 5004, config.WorkloadToBeRun.FfmpegPipeline.SourcePort)
    assert.Contains(t, config.WorkloadToBeRun.FfmpegPipeline.EnvironmentVariables, "ENV_VAR_1=value1")
    assert.Contains(t, config.WorkloadToBeRun.FfmpegPipeline.EnvironmentVariables, "ENV_VAR_2=value2")

    assert.Equal(t, "nmos-client", config.WorkloadToBeRun.NmosClient.Name)
    assert.Equal(t, "nmos-client:latest", config.WorkloadToBeRun.NmosClient.ImageAndTag)
    assert.Equal(t, "/path/to/nmos/config", config.WorkloadToBeRun.NmosClient.NmosConfigPath)
    assert.Equal(t, "nmos_config.json", config.WorkloadToBeRun.NmosClient.NmosConfigFileName)
    assert.Equal(t, "localhost", config.WorkloadToBeRun.NmosClient.FfmpegConectionAddress)
    assert.Equal(t, "58000", config.WorkloadToBeRun.NmosClient.FfmpegConnectionPort)
}
func TestParseLauncherMode(t *testing.T) {
    mockConfig := `
k8s: true
configuration:
  runOnce:
    mediaProxyAgent:
      imageAndTag: "mesh-agent:25.04"
      gRPCPort: "50051"
      restPort: "80"
      network:
        enable: true
        name: "host"
        ip: "localhost"
    mediaProxyMcm:
      imageAndTag: "mcm:25.04"
      interfaceName: "eth1"
      volumes:
        - "/dev/vfio:/dev/vfio"
      network:
        enable: true
        name: "host"
        ip: "localhost"
  workloadToBeRun:
    ffmpegPipeline:
      name: "ffmpeg"
      imageAndTag: "ffmpeg:latest"
      gRPCPort: 58000
      sourcePort: 5004
      environmentVariables:
        - "ENV_VAR_1=value1"
        - "ENV_VAR_2=value2"
      volumes:
        videos: "/videos"
        dri: "/dev/dri"
        kahawai: "/kahawai"
        devnull: "/dev/null"
        tmpHugepages: "/tmp/hugepages"
        hugepages: "/hugepages"
        imtl: "/imtl"
        shm: "/dev/shm"
      devices:
        vfio: "/dev/vfio"
        dri: "/dev/dri"
      network:
        enable: true
        name: "host"
        ip: "localhost"
    nmosClient:
      name: "nmos-client"
      imageAndTag: "nmos-client:latest"
      environmentVariables:
        - "ENV_VAR_1=value1"
        - "ENV_VAR_2=value2"
      nmosConfigPath: "/path/to/nmos/config"
      nmosConfigFileName: "nmos_config.json"
      network:
        enable: true
        name: "host"
        ip: "localhost"
      ffmpegConectionAddress: "localhost"
      ffmpegConnectionPort: "58000"
`
    tempFile, err := os.CreateTemp("", "mode-*.yaml")
    assert.NoError(t, err)
    defer os.Remove(tempFile.Name())

    _, err = tempFile.Write([]byte(mockConfig))
    assert.NoError(t, err)
    tempFile.Close()

    // Act: Parse the mode using the function
    modeK8s, err := ParseLauncherMode(tempFile.Name())

    // Assert: Verify the parsed mode matches the expected value
    assert.NoError(t, err)
    assert.True(t, modeK8s)
}

func TestParseLauncherMode_InvalidFile(t *testing.T) {
    // Act: Attempt to parse a non-existent file
    _, err := ParseLauncherMode("non-existent-file.yaml")

    // Assert: Verify an error is returned
    assert.Error(t, err)
}

func TestParseLauncherMode_InvalidYAML(t *testing.T) {
    invalidYAML := `
k8s: true
configuration:
  runOnce:
    mediaProxyAgent:
      imageAndTag: "mesh-agent:25.04"
      gRPCPort: "50051"
      restPort: "80"
      network:
        enable: true
        name: "host"
        ip: "localhost"
    mediaProxyMcm:
      imageAndTag: "mcm:25.04"
      interfaceName: "eth1"
      volumes:
        - "/dev/vfio:/dev/vfio"
      network:
        enable: true
        name: "host"
        ip: "localhost"
  workloadToBeRun:
    ffmpegPipeline:
      name: "ffmpeg"
      imageAndTag: "ffmpeg:latest"
      gRPCPort: 58000
      sourcePort: 5004
      environmentVariables:
        - "ENV_VAR_1=value1"
        - "ENV_VAR_2=value2"
      volumes:
        videos: "/videos"
        dri: "/dev/dri"
        kahawai: "/kahawai"
        devnull: "/dev/null"
        tmpHugepages: "/tmp/hugepages"
        hugepages: "/hugepages"
        imtl: "/imtl"
        shm: "/dev/shm"
      devices:
        vfio: "/dev/vfio"
        dri: "/dev/dri"
      network:
        enable: true
        name: "host"
        ip: "localhost"
    nmosClient:
      name: "nmos-client"
      imageAndTag: "nmos-client:latest"
      environmentVariables:
        - "ENV_VAR_1=value1"
        - "ENV_VAR_2=value2"
      nmosConfigPath: "/path/to/nmos/config"
      nmosConfigFileName: "nmos_config.json"
      network:
        enable: true
        name: "host"
        ip: "localhost"
      ffmpegConectionAddress: "localhost"
      ffmpegConnectionPort: "58000"
  invalidField: [unclosed-bracket
`
    tempFile, err := os.CreateTemp("", "invalid-*.yaml")
    assert.NoError(t, err)
    defer os.Remove(tempFile.Name())

    _, err = tempFile.Write([]byte(invalidYAML))
    assert.NoError(t, err)
    tempFile.Close()

    // Act: Attempt to parse the invalid YAML
    _, err = ParseLauncherMode(tempFile.Name())

    // Assert: Verify an error is returned
    assert.Error(t, err)
}