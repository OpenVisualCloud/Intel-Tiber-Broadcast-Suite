/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

package v1

import (
	"bcs.pod.launcher.intel/resources_library/resources/bcs"
	"bcs.pod.launcher.intel/resources_library/resources/nmos"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BcsConfigSpec struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	App       App `json:"app"`
	Nmos      Nmos `json:"nmos"`
	ScheduleOnNode []string `yaml:"scheduleOnNode,omitempty"`
	DoNotScheduleOnNode []string `yaml:"doNotScheduleOnNode,omitempty"`
  }

  type App struct {
	Image               string `json:"image"`
	GrpcPort            int    `json:"grpcPort"`
	EnvironmentVariables []EnvVar `json:"environmentVariables"`
	Volumes             map[string]string `json:"volumes"`
	Resources bcs.HwResources `json:"resources,omitempty"`
  }

  type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
  }

  type Nmos struct {
	Image                     string `json:"image"`
	Args                      []string `json:"args"`
	EnvironmentVariables      []EnvVar `json:"environmentVariables"`
	NmosApiNodePort           int `json:"nmosApiNodePort"`
	NmosInputFile             nmos.Config `json:"nmosInputFile"`
	Resources bcs.HwResources `json:"resources,omitempty"`
  }


// BcsConfigStatus defines the observed state of BcsConfig
type BcsConfigStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// BcsConfig is the Schema for the bcsconfigs API
type BcsConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BcsConfigSpec   `json:"spec,omitempty"`
	Status BcsConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BcsConfigList contains a list of BcsConfig
type BcsConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BcsConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BcsConfig{}, &BcsConfigList{})
}
