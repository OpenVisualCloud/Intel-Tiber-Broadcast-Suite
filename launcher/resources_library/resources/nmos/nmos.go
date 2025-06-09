//
//  SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
//
//  SPDX-License-Identifier: BSD-3-Clause
//

package nmos

type Config struct {
	LoggingLevel            int        `json:"logging_level"`
	HttpPort                int        `json:"http_port"`
	Label                   string     `json:"label"`
	DeviceTags              DeviceTags `json:"device_tags"`
	Function                string     `json:"function"`
	ActivateSenders         bool       `json:"activate_senders"`
	MultiviewerColumns      int        `json:"multiviewer_columns,omitempty"`
	StreamLoop              int        `json:"stream_loop"`
	GpuHwAcceleration       string     `json:"gpu_hw_acceleration"`
	GpuHwAccelerationDevice string     `json:"gpu_hw_acceleration_device,omitempty"`
	Domain                  string     `json:"domain"`
	DefaultSourceIP         string     `json:"default_source_ip,omitempty"`
	DefaultDestinationIP    string     `json:"default_destination_ip,omitempty"`
	DefaultSourcePort       int        `json:"default_source_port,omitempty"`
	DefaultDestinationPort  int        `json:"default_destination_port,omitempty"`
	DefaultInterfaceIP      string     `json:"default_interface_ip,omitempty"`
	FfmpegGrpcServerAddress string     `json:"ffmpeg_grpc_server_address,omitempty"`
	FfmpegGrpcServerPort    string     `json:"ffmpeg_grpc_server_port,omitempty"`
	SenderPayloadType       int        `json:"sender_payload_type"`
	Sender                  []Sender   `json:"sender"`
	Receiver                []Receiver `json:"receiver"`
}

type DeviceTags struct {
	Pipeline []string `json:"pipeline"`
}

type Sender struct {
	StreamPayload StreamPayload `json:"stream_payload"`
	StreamType    StreamType    `json:"stream_type"`
}

type Receiver struct {
	StreamPayload StreamPayload `json:"stream_payload"`
	StreamType    StreamType    `json:"stream_type"`
}

type StreamPayload struct {
	Video Video `json:"video,omitempty"`
	Audio Audio `json:"audio,omitempty"`
}

type Video struct {
	FrameWidth  int       `json:"frame_width"`
	FrameHeight int       `json:"frame_height"`
	FrameRate   FrameRate `json:"frame_rate"`
	PixelFormat string    `json:"pixel_format"`
	VideoType   string    `json:"video_type"`
	Preset      string    `json:"preset,omitempty"`
	Profile     string    `json:"profile,omitempty"`
}

type FrameRate struct {
	Numerator   int `json:"numerator"`
	Denominator int `json:"denominator"`
}

type Audio struct {
	Channels   int    `json:"channels"`
	SampleRate int    `json:"sampleRate"`
	Format     string `json:"format"`
	PacketTime string `json:"packetTime"`
}

type StreamType struct {
	St2110 *St2110 `json:"st2110,omitempty"`
	Mcm    *Mcm    `json:"mcm,omitempty"`
	File   *File   `json:"file,omitempty"`
}

type St2110 struct {
	Transport    string `json:"transport"`
	Payload_type int    `json:"payloadType"`
	QueuesCount  int    `json:"queues_cnt"`
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
