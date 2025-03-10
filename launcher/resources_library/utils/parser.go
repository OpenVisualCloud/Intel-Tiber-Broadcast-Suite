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


type Config struct {
	ModeK8s       bool          `yaml:"k8s"`
	Configuration Configuration `yaml:"configuration"`
}

type Configuration struct {
	RunOnce         RunOnce         `yaml:"runOnce"`
	WorkloadToBeRun WorkloadConfig `yaml:"workloadToBeRun"`
}

type RunOnce struct {
	MediaProxyAgent workloads.MediaProxyAgentConfig `yaml:"mediaProxyAgent"`
	MediaProxyMcm   workloads.MediaProxyMcmConfig   `yaml:"mediaProxyMcm"`
}

  type WorkloadConfig struct {
	FfmpegPipeline  workloads.FfmpegPipelineConfig `yaml:"ffmpegPipeline"`
	NmosClient      workloads.NmosClientConfig     `yaml:"nmosClient"`
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
