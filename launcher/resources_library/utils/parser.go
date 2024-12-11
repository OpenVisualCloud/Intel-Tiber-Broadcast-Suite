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
	ModeK8s bool `yaml:"k8s"`
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
