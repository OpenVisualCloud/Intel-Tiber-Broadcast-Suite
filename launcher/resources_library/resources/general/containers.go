/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

package general

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

type Containers struct {
	Type            Workload
	Image           string // image + tag
	ContainerName   string
	Ip              string
	ExposedPort     string // "format should be: 80/tcp"
	BindingHostPort string
	NetworkMode     NetworkMode
	Overridden      string
	Privileged      bool
	VolumeMount     []string
}
