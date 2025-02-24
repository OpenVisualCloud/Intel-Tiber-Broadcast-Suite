/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

package general

import "bcs.pod.launcher.intel/resources_library/workloads"

type Workload int
type NetworkMode string

const (
	NetworkModeHost NetworkMode = "host"
)

const (
	MediaProxyAgent Workload = iota
	MediaProxyMCM
	BcsPipelineFfmpeg
	BcsPipelineNmosClient
)

func (w Workload) String() string {
	return [...]string{"MediaProxyAgent", "MediaProxyMCM", "BcsPipelineFfmpeg", "BcsPipelineNmosClient"}[w]
}

type ContainersConfig struct{
MediaProxyAgentConfig  workloads.MediaProxyAgentConfig
MediaProxyMcmConfig  workloads.MediaProxyMcmConfig
WorkloadConfig  workloads.WorkloadConfig
}

type Containers struct {
	Type            Workload
	Configuration   ContainersConfig
	ContainerName   string
	Image           string // image + tag
	// Command    		string
	// ContainerName   string
	// Ip              string
	// ExposedPort     []string // "format should be: 80/tcp"
	// BindingHostPort []string
	// NetworkMode     string
	// Overridden      string
	// Privileged      bool
	// VolumeMount     []string
	// EnviromentVariables []string
	// Network         string
    // DeviceDri       string
    // DeviceVfio      string
}
