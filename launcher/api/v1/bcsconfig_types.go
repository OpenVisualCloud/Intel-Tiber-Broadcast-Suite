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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BcsConfigSpec defines the desired state of BcsConfig
type BcsConfigSpec struct {
	AppParams   AppParams   `json:"appParams"`
	Connections Connections `json:"connections"`
}

type AppParams struct {
	UniqueName  string `json:"uniqueName"`
	Codec       string `json:"codec"`
	PixelFormat string `json:"pixelFormat"`
	Height      int    `json:"height"`
	Width       int    `json:"width"`
}

type Connections struct {
	DataConnection    DataConnection    `json:"dataConnection"`
	ControlConnection ControlConnection `json:"controlConnection"`
}

type DataConnection struct {
	ConnType            string `json:"connType"`
	MediaProxyIpAddress string `json:"mediaProxyIpAddress"`
	Port                int    `json:"port"`
}

type ControlConnection struct {
	IpAddress string `json:"ipAddress"`
	Port      int    `json:"port"`
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
