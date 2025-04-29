/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

package bcs

import "bcs.pod.launcher.intel/resources_library/resources/general"

type BcsApp struct {
	Name       string
	Namespace  string
	Containers general.Containers
}

type HwResources struct {
	Requests struct {
		CPU    string `yaml:"cpu"`
		Memory string `yaml:"memory"`
		Hugepages1Gi string `yaml:"hugepages-1Gi,omitempty"`
	    Hugepages2Mi string `yaml:"hugepages-2Mi,omitempty"`
	} `yaml:"requests"`
	Limits struct {
		CPU    string `yaml:"cpu"`
		Memory string `yaml:"memory"`
		Hugepages1Gi string `yaml:"hugepages-1Gi,omitempty"`
	    Hugepages2Mi string `yaml:"hugepages-2Mi,omitempty"`
	} `yaml:"limits"`
}
