/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

// This package is exploited SDBQ-1261
package utils

import (
	"os"

	"bcs.pod.launcher.intel/resources_library/workloads"
	"gopkg.in/yaml.v2"
)

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
	NmosPort             int           `yaml:"nmosPort"`
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
	Name   string `yaml:"name"`
	IP     string `yaml:"ip"`
}
type Config struct {
	ModeK8s       bool          `yaml:"k8s"`
	Configuration Configuration `yaml:"configuration"`
}

type Configuration struct {
	RunOnce         RunOnce          `yaml:"runOnce"`
	WorkloadToBeRun []WorkloadConfig `yaml:"workloadToBeRun"`
}

type RunOnce struct {
	MediaProxyAgent workloads.MediaProxyAgentConfig `yaml:"mediaProxyAgent"`
	MediaProxyMcm   workloads.MediaProxyMcmConfig   `yaml:"mediaProxyMcm"`
}

func ParseLauncherMode(filename string) (bool, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return false, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return false, err
	}
	return config.ModeK8s, nil
}

func ParseLauncherConfiguration(filename string) (Configuration, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Configuration{}, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Configuration{}, err
	}
	return config.Configuration, nil
}
