//
//  SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
//
//  SPDX-License-Identifier: BSD-3-Clause
//

package workloads

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
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
ffmpegConnectionAddress: "192.168.1.101"
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
		assert.Equal(t, "192.168.1.101", config.FfmpegConnectionAddress)
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
		assert.Empty(t, config.FfmpegConnectionAddress)
		assert.Empty(t, config.FfmpegConnectionPort)
	})
}
