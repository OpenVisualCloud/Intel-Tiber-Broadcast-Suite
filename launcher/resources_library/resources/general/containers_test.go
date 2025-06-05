package general

import (
	"testing"

	"bcs.pod.launcher.intel/resources_library/workloads"

	"github.com/stretchr/testify/assert"
)

func TestWorkload_String(t *testing.T) {
	assert.Equal(t, "MediaProxyAgent", MediaProxyAgent.String())
	assert.Equal(t, "MediaProxyMCM", MediaProxyMCM.String())
	assert.Equal(t, "BcsPipelineFfmpeg", BcsPipelineFfmpeg.String())
	assert.Equal(t, "BcsPipelineNmosClient", BcsPipelineNmosClient.String())
}

func TestNetworkModeConstants(t *testing.T) {
	assert.Equal(t, "host", string(NetworkModeHost))
}

func TestContainersStruct(t *testing.T) {
	container := Containers{
		Type:          MediaProxyAgent,
		ContainerName: "test-container",
		Image:         "test-image:latest",
		Id:            1,
	}

	assert.Equal(t, MediaProxyAgent, container.Type)
	assert.Equal(t, "test-container", container.ContainerName)
	assert.Equal(t, "test-image:latest", container.Image)
	assert.Equal(t, 1, container.Id)
}

func TestContainersConfigStruct(t *testing.T) {
	config := ContainersConfig{
		MediaProxyAgentConfig: workloads.MediaProxyAgentConfig{ImageAndTag: "agent-config:latest"},
		MediaProxyMcmConfig:   workloads.MediaProxyMcmConfig{ImageAndTag: "mcm-config:latest"},
		WorkloadConfig: workloads.WorkloadConfig{FfmpegPipeline: workloads.FfmpegPipelineConfig{ImageAndTag: "workload-config:latest"},
			NmosClient: workloads.NmosClientConfig{ImageAndTag: "nmos-config:latest"}},
	}

	assert.Equal(t, "agent-config:latest", config.MediaProxyAgentConfig.ImageAndTag)
	assert.Equal(t, "mcm-config:latest", config.MediaProxyMcmConfig.ImageAndTag)
	assert.Equal(t, "workload-config:latest", config.WorkloadConfig.FfmpegPipeline.ImageAndTag)
	assert.Equal(t, "nmos-config:latest", config.WorkloadConfig.NmosClient.ImageAndTag)
}
