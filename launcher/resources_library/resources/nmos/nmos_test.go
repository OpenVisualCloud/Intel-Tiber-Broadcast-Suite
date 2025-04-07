//
//  SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
//
//  SPDX-License-Identifier: BSD-3-Clause
//

package nmos

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigSerialization(t *testing.T) {
	// Arrange: Create a sample Config object
	config := Config{
		MultiviewerColumns:     2,
		GpuHwAccelerationDevice: "/dev/dri/renderD128",
		LoggingLevel:           3,
		HttpPort:               8080,
		Label:                  "TestLabel",
		DeviceTags:             map[string][]string{"device1": {"tag1", "tag2"}},
		Function:               "TestFunction",
		StreamLoop:             1,
		GpuHwAcceleration:      "enabled",
		Domain:                 "test.domain",
		FfmpegGrpcServerAddress: "127.0.0.1",
		FfmpegGrpcServerPort:   "50051",
		SenderPayloadType:      96,
		Sender: []Stream{
			{
				StreamPayload: StreamPayload{
					Video: Video{
						FrameWidth:  1920,
						FrameHeight: 1080,
						FrameRate: FrameRate{
							Numerator:   30,
							Denominator: 1,
						},
						PixelFormat: "yuv420p",
						VideoType:   "H264",
					},
				},
				StreamType: StreamType{
					St2110: &St2110{
						Transport:   "udp",
						PayloadType: 96,
					},
				},
			},
		},
		Receiver: []Stream{
			{
				StreamPayload: StreamPayload{
					Audio: Audio{
						Channels:   2,
						SampleRate: 48000,
						Format:     "pcm_s16le",
						PacketTime: "1ms",
					},
				},
				StreamType: StreamType{
					File: &File{
						Path:     "/media",
						Filename: "test.mp4",
					},
				},
			},
		},
	}

	// Act: Serialize the Config object to JSON
	jsonData, err := json.Marshal(config)
	assert.NoError(t, err, "JSON serialization should not produce an error")

	// Act: Deserialize the JSON back to a Config object
	var deserializedConfig Config
	err = json.Unmarshal(jsonData, &deserializedConfig)
	assert.NoError(t, err, "JSON deserialization should not produce an error")

	// Assert: Verify that the original and deserialized Config objects are equal
	assert.Equal(t, config, deserializedConfig, "Original and deserialized Config objects should be equal")
}

func TestConfigDefaultValues(t *testing.T) {
	// Arrange: Create a Config object with minimal fields
	config := Config{
		HttpPort: 8080,
		Function: "TestFunction",
		Domain:   "test.domain",
	}

	// Act: Serialize and deserialize the Config object
	jsonData, err := json.Marshal(config)
	assert.NoError(t, err, "JSON serialization should not produce an error")

	var deserializedConfig Config
	err = json.Unmarshal(jsonData, &deserializedConfig)
	assert.NoError(t, err, "JSON deserialization should not produce an error")

	// Assert: Verify that default values are preserved
	assert.Equal(t, config.HttpPort, deserializedConfig.HttpPort, "HttpPort should match")
	assert.Equal(t, config.Function, deserializedConfig.Function, "Function should match")
	assert.Equal(t, config.Domain, deserializedConfig.Domain, "Domain should match")
	assert.Empty(t, deserializedConfig.Sender, "Sender should be empty by default")
	assert.Empty(t, deserializedConfig.Receiver, "Receiver should be empty by default")
}