package workloads

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestFfmpegPipelineConfig_UnmarshalYAML(t *testing.T) {
	yamlData := `
name: test-pipeline
imageAndTag: ffmpeg:latest
gRPCPort: 50051
environmentVariables:
  - VAR1=value1
  - VAR2=value2
volumes:
  videos: /path/to/videos
  dri: /path/to/dri
  kahawai: /path/to/kahawai
  devnull: /dev/null
  tmpHugepages: /path/to/tmpHugepages
  hugepages: /path/to/hugepages
  imtl: /path/to/imtl
  shm: /dev/shm
devices:
  vfio: /dev/vfio
  dri: /dev/dri
custom_network:
  enable: true
  name: custom-net
  ip: 192.168.1.100
`
	var config FfmpegPipelineConfig
	err := yaml.Unmarshal([]byte(yamlData), &config)
	assert.NoError(t, err)
	assert.Equal(t, "test-pipeline", config.Name)
	assert.Equal(t, "ffmpeg:latest", config.ImageAndTag)
	assert.Equal(t, 50051, config.GRPCPort)
	assert.Equal(t, []string{"VAR1=value1", "VAR2=value2"}, config.EnvironmentVariables)
	assert.Equal(t, "/path/to/videos", config.Volumes.Videos)
	assert.Equal(t, "/dev/vfio", config.Devices.Vfio)
	assert.True(t, config.Network.Enable)
	assert.Equal(t, "custom-net", config.Network.Name)
	assert.Equal(t, "192.168.1.100", config.Network.IP)
}

func TestNmosClientConfig_UnmarshalYAML(t *testing.T) {
	yamlData := `
name: test-nmos-client
imageAndTag: nmos-client:latest
environmentVariables:
  - NMOS_VAR1=value1
  - NMOS_VAR2=value2
nmosConfigPath: /path/to/config
nmosConfigFileName: config.json
custom_network:
  enable: true
  name: nmos-net
  ip: 192.168.1.101
nmosPort: 8080
ffmpegConnectionAddress: 192.168.1.102
ffmpegConnectionPort: 9090
`
	var config NmosClientConfig
	err := yaml.Unmarshal([]byte(yamlData), &config)
	assert.NoError(t, err)
	assert.Equal(t, "test-nmos-client", config.Name)
	assert.Equal(t, "nmos-client:latest", config.ImageAndTag)
	assert.Equal(t, []string{"NMOS_VAR1=value1", "NMOS_VAR2=value2"}, config.EnvironmentVariables)
	assert.Equal(t, "/path/to/config", config.NmosConfigPath)
	assert.Equal(t, "config.json", config.NmosConfigFileName)
	assert.True(t, config.Network.Enable)
	assert.Equal(t, "nmos-net", config.Network.Name)
	assert.Equal(t, "192.168.1.101", config.Network.IP)
	assert.Equal(t, 8080, config.NmosPort)
	assert.Equal(t, "192.168.1.102", config.FfmpegConnectionAddress)
	assert.Equal(t, "9090", config.FfmpegConnectionPort)
}

func TestNetworkConfig_UnmarshalYAML(t *testing.T) {
	yamlData := `
enable: true
name: test-network
ip: 192.168.1.1
`
	var config NetworkConfig
	err := yaml.Unmarshal([]byte(yamlData), &config)
	assert.NoError(t, err)
	assert.True(t, config.Enable)
	assert.Equal(t, "test-network", config.Name)
	assert.Equal(t, "192.168.1.1", config.IP)
}
