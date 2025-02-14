/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

// This package is exploited SDBQ-1261
package utils

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ModeK8s       bool          `yaml:"k8s"`
	Configuration Configuration `yaml:"configuration"`
}

type Configuration struct {
	RunOnce         RunOnce         `yaml:"runOnce"`
	WorkloadToBeRun WorkloadToBeRun `yaml:"workloadToBeRun"`
}

type RunOnce struct {
	MediaProxyAgent MediaProxyAgent `yaml:"mediaProxyAgent"`
	MediaProxyMcm   MediaProxyMcm   `yaml:"mediaProxyMcm"`
}

type MediaProxyAgent struct {
	ImageAndTag     string   `yaml:"imageAndTag"`
	Command         string   `yaml:"command"`
	ExposedPort     []string   `yaml:"exposedPort"`
	BindingHostPort []string   `yaml:"bindingHostPort"`
	IP              string   `yaml:"ip"`
	Volumes         []string `yaml:"volumes"`
	Privileged      bool     `yaml:"privileged"`
}

type MediaProxyMcm struct {
	ImageAndTag     string   `yaml:"imageAndTag"`
	Command         string   `yaml:"command"`
	ExposedPort     []string   `yaml:"exposedPort"`
	BindingHostPort []string   `yaml:"bindingHostPort"`
	IP              string   `yaml:"ip"`
	Volumes         []string `yaml:"volumes"`
	Privileged      bool     `yaml:"privileged"`
}

type WorkloadToBeRun struct {
	FFmpegPipeline FFmpegPipeline `yaml:"ffmpegPipeline"`
	NmosClient     NmosClient     `yaml:"nmosClient"`
}

type FFmpegPipeline struct {
	Name            string `yaml:"name"`
	ImageAndTag     string `yaml:"imageAndTag"`
	Command         string `yaml:"command"`
	ExposedPort     []string `yaml:"exposedPort"`
	BindingHostPort []string `yaml:"bindingHostPort"`
	IP              string `yaml:"ip"`
    Network         string `yaml:"network"`
    DeviceDri       string `yaml:"deviceDri"`
    DeviceVfio      string `yaml:"deviceVfio"`
    Volumes         []string `yaml:"volumes"`
}

type NmosClient struct {
	Name            string `yaml:"name"`
	ImageAndTag     string `yaml:"imageAndTag"`
	Command         string `yaml:"command"`
	ExposedPort     []string `yaml:"exposedPort"`
	BindingHostPort []string `yaml:"bindingHostPort"`
	IP              string `yaml:"ip"`
	Network         string `yaml:"network"`
	Volumes         []string `yaml:"volumes"`
	EnvironmentVariables []string `yaml:"environmentVariables"`
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
