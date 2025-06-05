/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */
package bcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestHwResources_UnmarshalYAML(t *testing.T) {
	yamlData := `
requests:
  cpu: "500m"
  memory: "256Mi"
  hugepages-1Gi: "1Gi"
  hugepages-2Mi: "2Mi"
limits:
  cpu: "1000m"
  memory: "512Mi"
  hugepages-1Gi: "2Gi"
  hugepages-2Mi: "4Mi"
`
	var hwResources HwResources
	err := yaml.Unmarshal([]byte(yamlData), &hwResources)
	assert.NoError(t, err)
	assert.Equal(t, "500m", hwResources.Requests.CPU)
	assert.Equal(t, "256Mi", hwResources.Requests.Memory)
	assert.Equal(t, "1Gi", hwResources.Requests.Hugepages1Gi)
	assert.Equal(t, "2Mi", hwResources.Requests.Hugepages2Mi)
	assert.Equal(t, "1000m", hwResources.Limits.CPU)
	assert.Equal(t, "512Mi", hwResources.Limits.Memory)
	assert.Equal(t, "2Gi", hwResources.Limits.Hugepages1Gi)
	assert.Equal(t, "4Mi", hwResources.Limits.Hugepages2Mi)
}
