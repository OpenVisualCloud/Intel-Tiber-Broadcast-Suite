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

type BcsConfigSpec struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	App       App `json:"app"`
	Nmos      Nmos `json:"nmos"`
  }

  type App struct {
	Image               string `json:"image"`
	GrpcPort            int    `json:"grpcPort"`
	EnvironmentVariables []EnvVar `json:"environmentVariables"`
	Volumes             map[string]string `json:"volumes"`
  }

  type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
  }

  type Nmos struct {
	Image                     string `json:"image"`
	Args                      []string `json:"args"`
	EnvironmentVariables      []EnvVar `json:"environmentVariables"`
	NmosApiPort               int `json:"nmosApiPort"`
	NmosApiNodePort           int `json:"nmosApiNodePort"`
	NmosAppCommunicationPort  int `json:"nmosAppCommunicationPort"`
	NmosAppCommunicationNodePort int `json:"nmosAppCommunicationNodePort"`
	NmosInputFile             NmosInputFile `json:"nmosInputFile"`
  }

  type NmosInputFile struct {
	LoggingLevel            int `json:"logging_level"`
	HttpPort                int `json:"http_port"`
	Label                   string `json:"label"`
	Senders                 []string `json:"senders"`
	SendersCount            []int `json:"senders_count"`
	Receivers               []string `json:"receivers"`
	ReceiversCount          []int `json:"receivers_count"`
	DeviceTags              DeviceTags `json:"device_tags"`
	Function                string `json:"function"`
	GpuHwAcceleration       string `json:"gpu_hw_acceleration"`
	Domain                  string `json:"domain"`
	FfmpegGrpcServerAddress string `json:"ffmpeg_grpc_server_address"`
	FfmpegGrpcServerPort    string `json:"ffmpeg_grpc_server_port"`
	SenderPayloadType       int `json:"sender_payload_type"`
	Sender                  []Sender `json:"sender"`
	Receiver                []Receiver `json:"receiver"`
  }

  type DeviceTags struct {
	Pipeline []string `json:"pipeline"`
  }

  type Sender struct {
	StreamPayload StreamPayload `json:"stream_payload"`
	StreamType    StreamType `json:"stream_type"`
  }

  type Receiver struct {
	StreamPayload StreamPayload `json:"stream_payload"`
	StreamType    StreamType `json:"stream_type"`
  }

  type StreamPayload struct {
	Video Video `json:"video,omitempty"`
	Audio Audio `json:"audio,omitempty"`
  }

  type Video struct {
	FrameWidth  int `json:"frame_width"`
	FrameHeight int `json:"frame_height"`
	FrameRate   FrameRate `json:"frame_rate"`
	PixelFormat string `json:"pixel_format"`
	VideoType   string `json:"video_type"`
  }

  type FrameRate struct {
	Numerator   int `json:"numerator"`
	Denominator int `json:"denominator"`
  }

  type Audio struct {
	Channels   int `json:"channels"`
	SampleRate int `json:"sampleRate"`
	Format     string `json:"format"`
	PacketTime string `json:"packetTime"`
  }

  type StreamType struct {
	St2100 *St2100 `json:"st2100,omitempty"`
	Mcm  *Mcm  `json:"mcm,omitempty"`
	File *File `json:"file,omitempty"`
  }

  type St2100 struct {
	Transport string `json:"transport"`
	Payload_type string `json:"payload_type"`
  }

  type Mcm struct {
	ConnType             string `json:"conn_type"`
	Transport            string `json:"transport"`
	Urn                  string `json:"urn"`
	TransportPixelFormat string `json:"transportPixelFormat"`
  }

  type File struct {
	Path     string `json:"path"`
	Filename string `json:"filename"`
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
