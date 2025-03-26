//
//  SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
//
//  SPDX-License-Identifier: BSD-3-Clause
//

package nmos

type Config struct {
	MultiviewerColumns     int				 `json:"multiviewer_columns,omitempty"`
	GpuHwAccelerationDevice string			 `json:"gpu_hw_acceleration_device,omitempty"`
	LoggingLevel           int               `json:"logging_level,omitempty"`
	HttpPort               int               `json:"http_port"`
	Label                  string            `json:"label,omitempty"`
	Senders                []string          `json:"senders"`
	SendersCount           []int             `json:"senders_count"`
	Receivers              []string          `json:"receivers"`
	ReceiversCount         []int             `json:"receivers_count"`
	DeviceTags             map[string][]string `json:"device_tags,omitempty"`
	Function               string            `json:"function"`
	StreamLoop             int               `json:"stream_loop,omitempty"`
	GpuHwAcceleration      string            `json:"gpu_hw_acceleration"`
	Domain                 string            `json:"domain"`
	FfmpegGrpcServerAddress string           `json:"ffmpeg_grpc_server_address"`
	FfmpegGrpcServerPort   string            `json:"ffmpeg_grpc_server_port"`
	SenderPayloadType      int               `json:"sender_payload_type"`
	Sender                 []Stream          `json:"sender"`
	Receiver               []Stream          `json:"receiver"`
}

type Stream struct {
	StreamPayload StreamPayload `json:"stream_payload"`
	StreamType    StreamType    `json:"stream_type"`
}

type StreamPayload struct {
	Video Video `json:"video"`
	Audio Audio `json:"audio,omitempty"`
}

type Video struct {
	FrameWidth   int    `json:"frame_width"`
	FrameHeight  int    `json:"frame_height"`
	FrameRate    FrameRate `json:"frame_rate"`
	PixelFormat  string `json:"pixel_format,omitempty"`
	VideoType    string `json:"video_type"`
}

type FrameRate struct {
	Numerator   int `json:"numerator,omitempty"`
	Denominator int `json:"denominator,omitempty"`
}

type Audio struct {
	Channels   int    `json:"channels,omitempty"`
	SampleRate int    `json:"sampleRate,omitempty"`
	Format     string `json:"format,omitempty"`
	PacketTime string `json:"packetTime,omitempty"`
}

type StreamType struct {
	St2110 *St2110 `json:"st2110,omitempty"`
	File   *File   `json:"file,omitempty"`
	MCM	*MCM    `json:"mcm,omitempty"`
}

type St2110 struct {
	Transport   string `json:"transport,omitempty"`
	PayloadType int    `json:"payloadType,omitempty"`
}

type File struct {
	Path     string `json:"path,omitempty"`
	Filename string `json:"filename,omitempty"`
}

type MCM struct {
	ConnType             string `json:"conn_type,omitempty"`
	Transport            string `json:"transport,omitempty"`
	URN                  string `json:"urn,omitempty"`
	TransportPixelFormat string `json:"transportPixelFormat,omitempty"`
}