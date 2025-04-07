package general

import (
	"testing"

	"bcs.pod.launcher.intel/resources_library/workloads"
	"github.com/stretchr/testify/assert"
)

func TestContainersConfigInitialization(t *testing.T) {
	// Arrange: Create mock data for workloads
    net1 := workloads.NetworkConfig {
		Enable: true,
		Name  : "host",
		IP    : "localhost",
	  }
	net2 := workloads.NetworkConfig {
		Enable: true,
		Name  : "host",
		IP    : "localhost",
	  }
	net3 := workloads.NetworkConfig {
		Enable: true,
		Name  : "host",
		IP    : "localhost",
	  }

	  net4 := workloads.NetworkConfig {
		Enable: true,
		Name  : "host",
		IP    : "localhost",
	  }

	mockMediaProxyAgentConfig := workloads.MediaProxyAgentConfig{
		ImageAndTag: "mesh-agent:25.04",
        GRPCPort: "50051",
        RestPort: "80",  
        Network: net1,
	}
	mockMediaProxyMcmConfig := workloads.MediaProxyMcmConfig{
		ImageAndTag:     "mcm:25.04",
		InterfaceName:   "eth1",
		Volumes:         []string{"/tmp"},
		Network:         net2,
	}

	app:= workloads.FfmpegPipelineConfig{
		Name: "ffmpeg",
		ImageAndTag: "ffmpeg:latest",
		GRPCPort: 58000,
		SourcePort: 5004,
		EnvironmentVariables: []string{
			"ENV_VAR_1=value1",
			"ENV_VAR_2=value2",
		},
		Volumes: workloads.Volumes{
			Videos:      "/videos",
			Dri:         "/dev/dri",
			Kahawai:     "/kahawai",
			Devnull:     "/dev/null",
			TmpHugepages: "/tmp/hugepages",
			Hugepages:   "/hugepages",
			Imtl:        "/imtl",
			Shm:         "/dev/shm",
		},
		Devices: workloads.Devices{
			Vfio: "/dev/vfio",
			Dri:  "/dev/dri",
		},
		Network: net3,
	  }
	
	  nmos:= workloads.NmosClientConfig{
		Name: "nmos-client",
		ImageAndTag: "nmos-client:latest",
		EnvironmentVariables: []string{
			"ENV_VAR_1=value1",
			"ENV_VAR_2=value2",
		},
		NmosConfigPath: "/path/to/nmos/config",
		NmosConfigFileName: "nmos_config.json",
		Network: net4,
		FfmpegConectionAddress: "localhost",
		FfmpegConnectionPort: "58000",
	  }

	mockWorkloadConfig := workloads.WorkloadConfig{
		FfmpegPipeline: app,
		NmosClient:     nmos,
	}

	// Act: Initialize a ContainersConfig struct
	config := ContainersConfig{
		MediaProxyAgentConfig: mockMediaProxyAgentConfig,
		MediaProxyMcmConfig:   mockMediaProxyMcmConfig,
		WorkloadConfig:        mockWorkloadConfig,
	}

	// Assert: Verify the struct fields
	assert.Equal(t, mockMediaProxyAgentConfig, config.MediaProxyAgentConfig, "MediaProxyAgentConfig should match")
	assert.Equal(t, mockMediaProxyMcmConfig, config.MediaProxyMcmConfig, "MediaProxyMcmConfig should match")
	assert.Equal(t, mockWorkloadConfig, config.WorkloadConfig, "WorkloadConfig should match")
}

func TestWorkloadStringMethod(t *testing.T) {
    // Arrange: Define expected string representations
    expectedStrings := []string{"MediaProxyAgent", "MediaProxyMCM", "BcsPipelineFfmpeg", "BcsPipelineNmosClient"}

    // Act & Assert: Verify the String() method for each Workload
    for i, expected := range expectedStrings {
        assert.Equal(t, expected, Workload(i).String(), "Workload string representation should match")
    }
}