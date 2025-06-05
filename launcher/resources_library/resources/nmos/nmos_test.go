package nmos

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_UnmarshalJSON(t *testing.T) {
	jsonData := `
{
    "logging_level": 10,
    "http_port": 5004,
    "label": "nmos-test",
    "device_tags": {
        "pipeline": ["tx", "rx"]
    },
    "function": "tx",
    "activate_senders": true,
    "stream_loop": -1,
    "gpu_hw_acceleration": "none",
    "domain": "local",
    "sender_payload_type": 112,
    "sender": [
        {
            "stream_payload": {
                "video": {
                    "frame_width": 1920,
                    "frame_height": 1080,
                    "frame_rate": {
                        "numerator": 60,
                        "denominator": 1
                    },
                    "pixel_format": "yuv422p10le",
                    "video_type": "rawvideo"
                }
            },
            "stream_type": {
                "st2110": {
                    "transport": "st2110-20",
                    "payloadType": 112,
                    "queues_cnt": 0
                }
            }
        }
    ],
    "receiver": [
        {
            "stream_payload": {
                "audio": {
                    "channels": 2,
                    "sampleRate": 48000,
                    "format": "pcm_s24be",
                    "packetTime": "1ms"
                }
            },
            "stream_type": {
                "file": {
                    "path": "/videos",
                    "filename": "testfile.yuv"
                }
            }
        }
    ]
}
`
	var config Config
	err := json.Unmarshal([]byte(jsonData), &config)
	assert.NoError(t, err)

	assert.Equal(t, 10, config.LoggingLevel)
	assert.Equal(t, 5004, config.HttpPort)
	assert.Equal(t, "nmos-test", config.Label)
	assert.Equal(t, []string{"tx", "rx"}, config.DeviceTags.Pipeline)
	assert.Equal(t, "tx", config.Function)
	assert.True(t, config.ActivateSenders)
	assert.Equal(t, -1, config.StreamLoop)
	assert.Equal(t, "none", config.GpuHwAcceleration)
	assert.Equal(t, "local", config.Domain)
	assert.Equal(t, 112, config.SenderPayloadType)

	assert.Len(t, config.Sender, 1)
	assert.Equal(t, 1920, config.Sender[0].StreamPayload.Video.FrameWidth)
	assert.Equal(t, 1080, config.Sender[0].StreamPayload.Video.FrameHeight)
	assert.Equal(t, 60, config.Sender[0].StreamPayload.Video.FrameRate.Numerator)
	assert.Equal(t, 1, config.Sender[0].StreamPayload.Video.FrameRate.Denominator)
	assert.Equal(t, "yuv422p10le", config.Sender[0].StreamPayload.Video.PixelFormat)
	assert.Equal(t, "rawvideo", config.Sender[0].StreamPayload.Video.VideoType)
	assert.Equal(t, "st2110-20", config.Sender[0].StreamType.St2110.Transport)
	assert.Equal(t, 112, config.Sender[0].StreamType.St2110.Payload_type)
	assert.Equal(t, 0, config.Sender[0].StreamType.St2110.QueuesCount)

	assert.Len(t, config.Receiver, 1)
	assert.Equal(t, 2, config.Receiver[0].StreamPayload.Audio.Channels)
	assert.Equal(t, 48000, config.Receiver[0].StreamPayload.Audio.SampleRate)
	assert.Equal(t, "pcm_s24be", config.Receiver[0].StreamPayload.Audio.Format)
	assert.Equal(t, "1ms", config.Receiver[0].StreamPayload.Audio.PacketTime)
	assert.Equal(t, "/videos", config.Receiver[0].StreamType.File.Path)
	assert.Equal(t, "testfile.yuv", config.Receiver[0].StreamType.File.Filename)
}
