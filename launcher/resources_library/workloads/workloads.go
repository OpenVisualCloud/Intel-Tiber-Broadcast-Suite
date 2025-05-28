//
//  SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
//
//  SPDX-License-Identifier: BSD-3-Clause
//

package workloads

type MediaProxyAgentConfig struct {
	ImageAndTag string        `yaml:"imageAndTag"`
	GRPCPort    string        `yaml:"gRPCPort"`
	RestPort    string        `yaml:"restPort"`
	Network     NetworkConfig `yaml:"custom_network"`
}

type MediaProxyMcmConfig struct {
	ImageAndTag   string        `yaml:"imageAndTag"`
	InterfaceName string        `yaml:"interfaceName"`
	Volumes       []string      `yaml:"volumes"`
	Network       NetworkConfig `yaml:"custom_network"`
}

type WorkloadConfig struct {
	FfmpegPipeline FfmpegPipelineConfig `yaml:"ffmpegPipeline"`
	NmosClient     NmosClientConfig     `yaml:"nmosClient"`
}

type Volumes struct {
	Videos       string `yaml:"videos"`
	Dri          string `yaml:"dri"`
	Kahawai      string `yaml:"kahawai"`
	Devnull      string `yaml:"devnull"`
	TmpHugepages string `yaml:"tmpHugepages,omitempty"`
	Hugepages    string `yaml:"hugepages,omitempty"`
	Imtl         string `yaml:"imtl"`
	Shm          string `yaml:"shm"`
}

type Devices struct {
	Vfio string `yaml:"vfio"`
	Dri  string `yaml:"dri"`
}

type FfmpegPipelineConfig struct {
	Name                 string        `yaml:"name"`
	ImageAndTag          string        `yaml:"imageAndTag"`
	GRPCPort             int           `yaml:"gRPCPort"`
	EnvironmentVariables []string      `yaml:"environmentVariables"`
	Volumes              Volumes       `yaml:"volumes"`
	Devices              Devices       `yaml:"devices"`
	Network              NetworkConfig `yaml:"custom_network"`
}

type NmosClientConfig struct {
	Name                    string        `yaml:"name"`
	ImageAndTag             string        `yaml:"imageAndTag"`
	EnvironmentVariables    []string      `yaml:"environmentVariables"`
	NmosConfigPath          string        `yaml:"nmosConfigPath"`
	NmosConfigFileName      string        `yaml:"nmosConfigFileName"`
	Network                 NetworkConfig `yaml:"custom_network"`
	NmosPort                int           `yaml:"nmosPort"`
	FfmpegConnectionAddress string        `yaml:"ffmpegConnectionAddress"`
	FfmpegConnectionPort    string        `yaml:"ffmpegConnectionPort"`
}

type NetworkConfig struct {
	Enable bool   `yaml:"enable"`
	Name   string `yaml:"name,omitempty"`
	IP     string `yaml:"ip,omitempty"`
}
