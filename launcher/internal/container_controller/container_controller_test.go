/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

package containercontroller

import (
	"context"
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"bcs.pod.launcher.intel/resources_library/resources/general"
	"bcs.pod.launcher.intel/resources_library/utils"
	"bcs.pod.launcher.intel/resources_library/workloads"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/go-logr/logr"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/mock"
)

// MockContainerController is a mock implementation of the ContainerController interface
type MockContainerController struct {
    mock.Mock
}

func (m *MockContainerController) CreateAndRunContainers(ctx context.Context, launcherConfigName string, log logr.Logger) error {
    args := m.Called(ctx, launcherConfigName, log)
    return args.Error(0)
}

func (m *MockContainerController) ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error) {
    args := m.Called(ctx, options)
    return args.Get(0).([]image.Summary), args.Error(1)
}

func (m *MockContainerController) ImagePull(ctx context.Context, ref string, options image.PullOptions) (io.ReadCloser, error) {
    args := m.Called(ctx, ref, options)
    return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockContainerController) ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error {
    args := m.Called(ctx, containerID, options)
    return args.Error(0)
}

func (m *MockContainerController) ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
    args := m.Called(ctx, options)
    return args.Get(0).([]types.Container), args.Error(1)
}

func (m *MockContainerController) ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error {
    args := m.Called(ctx, containerID, options)
    return args.Error(0)
}

func (m *MockContainerController) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) (container.CreateResponse, error) {
    args := m.Called(ctx, config, hostConfig, networkingConfig, nil, containerName)
    return args.Get(0).(container.CreateResponse), args.Error(1)
}

func TestIsContainerRunning(t *testing.T) {
    ctx := context.Background()

    t.Run("Container is running", func(t *testing.T) {
        mockController := new(MockContainerController)
        containerName := "test-container"
        mockController.On("ContainerList", ctx, container.ListOptions{All: true}).Return([]types.Container{
            {
                Names: []string{"/" + containerName},
                State: "running",
            },
        }, nil)

        err, isRunning := isContainerRunning(ctx, mockController, containerName)
        mockController.AssertExpectations(t)

        if err != nil {
            t.Fatalf("Expected no error, got %v", err)
        }
        if !isRunning {
            t.Fatalf("Expected container to be running, got %v", isRunning)
        }
    })

    t.Run("Container is not running", func(t *testing.T) {
        mockController := new(MockContainerController)
        containerName := "test-container"
        mockController.On("ContainerList", ctx, container.ListOptions{All: true}).Return([]types.Container{
            {
                Names: []string{"/" + containerName},
                State: "exited",
            },
        }, nil)

        err, isRunning := isContainerRunning(ctx, mockController, containerName)
        mockController.AssertExpectations(t)

        if err != nil {
            t.Fatalf("Expected no error, got %v", err)
        }
        if isRunning {
            t.Fatalf("Expected container to not be running, got %v", isRunning)
        }
    })

    t.Run("Container does not exist", func(t *testing.T) {
        mockController := new(MockContainerController)
        containerName := "non-existent-container"
        mockController.On("ContainerList", ctx, container.ListOptions{All: true}).Return([]types.Container{}, nil)

        err, isRunning := isContainerRunning(ctx, mockController, containerName)
        mockController.AssertExpectations(t)

        if err != nil {
            t.Fatalf("Expected no error, got %v", err)
        }
        if isRunning {
            t.Fatalf("Expected container to not be running, got %v", isRunning)
        }
    })
}
func TestRemoveContainer(t *testing.T) {
    ctx := context.Background()

    t.Run("Successfully removes container", func(t *testing.T) {
        mockController := new(MockContainerController)
        containerID := "test-container-id"
        mockController.On("ContainerRemove", ctx, containerID, container.RemoveOptions{Force: true}).Return(nil)

        err := removeContainer(ctx, mockController, containerID)
        mockController.AssertExpectations(t)

        if err != nil {
            t.Fatalf("Expected no error, got %v", err)
        }
    })

    t.Run("Fails to remove container", func(t *testing.T) {
        mockController := new(MockContainerController)
        containerID := "test-container-id"
        expectedError := errors.New("failed to remove container")
        mockController.On("ContainerRemove", ctx, containerID, container.RemoveOptions{Force: true}).Return(expectedError)

        err := removeContainer(ctx, mockController, containerID)
        mockController.AssertExpectations(t)

        if err == nil {
            t.Fatalf("Expected error, got nil")
        }
        if err != expectedError {
            t.Fatalf("Expected error %v, got %v", expectedError, err)
        }
    })
}
func TestCreateAndRunContainer(t *testing.T) {
    ctx := context.Background()
    log := logr.Discard()

    t.Run("Container is already running", func(t *testing.T) {
        mockController := new(MockContainerController)
        containerInfo := &general.Containers{
            ContainerName: "test-container",
            Image:         "test-image",
        }

        mockController.On("ContainerList", ctx, container.ListOptions{All: true}).Return([]types.Container{
            {
                Names: []string{"/" + containerInfo.ContainerName},
                State: "running",
            },
        }, nil)

        err := createAndRunContainer(ctx, mockController, log, containerInfo)
        mockController.AssertExpectations(t)

        if err != nil {
            t.Fatalf("Expected no error, got %v", err)
        }
    })

    t.Run("Container exists but is not running", func(t *testing.T) {
        mockController := new(MockContainerController)
        containerInfo := &general.Containers{
            ContainerName: "test-container",
            Image:         "test-image",
        }

        mockController.On("ContainerList", ctx, container.ListOptions{All: true}).Return([]types.Container{
            {
                Names: []string{"/" + containerInfo.ContainerName},
                State: "exited",
            },
        }, nil)
        mockController.On("ImageList", ctx, image.ListOptions{}).Return([]image.Summary{}, nil)
        mockController.On("ImagePull", ctx, containerInfo.Image, image.PullOptions{}).Return(io.NopCloser(strings.NewReader("")), nil)
        mockController.On("ContainerRemove", ctx, containerInfo.ContainerName, container.RemoveOptions{Force: true}).Return(nil)
        mockController.On("ContainerCreate", ctx, mock.Anything, mock.Anything, mock.Anything, nil, containerInfo.ContainerName).Return(container.CreateResponse{ID: "test-id"}, nil)
        mockController.On("ContainerStart", ctx, "test-id", container.StartOptions{}).Return(nil)

        err := createAndRunContainer(ctx, mockController, log, containerInfo)
        mockController.AssertExpectations(t)

        if err != nil {
            t.Fatalf("Expected no error, got %v", err)
        }
    })


    t.Run("Fails to create container", func(t *testing.T) {
        mockController := new(MockContainerController)
        containerInfo := &general.Containers{
            ContainerName: "test-container",
            Image:         "test-image",
            Type:          general.BcsPipelineFfmpeg,
            Configuration: general.ContainersConfig{
                WorkloadConfig: workloads.WorkloadConfig{
                    FfmpegPipeline: workloads.FfmpegPipelineConfig{
                        Name:        "pipeline-name",
                        ImageAndTag: "test-image",
                        GRPCPort:    50051,
                        SourcePort:  8080,
                        EnvironmentVariables: []string{
                            "ENV_VAR1=value1",
                            "ENV_VAR2=value2",
                        },
                        Volumes: workloads.Volumes{
                            Videos:      "/path/to/videos",
                            Dri:         "/path/to/dri",
                            Kahawai:     "/path/to/kahawai",
                            Devnull:     "/path/to/devnull",
                            TmpHugepages: "/path/to/tmpHugepages",
                            Hugepages:   "/path/to/hugepages",
                            Imtl:        "/path/to/imtl",
                            Shm:         "/path/to/shm",
                        },
                        Devices: workloads.Devices{
                            Vfio: "/dev/vfio",
                            Dri:  "/dev/dri",
                        },
                        Network: workloads.NetworkConfig{
                            Enable: true,
                            Name: "host",
                            IP: "localhost",
                        },
                    },
                },
            },
        }

        mockController.On("ContainerList", ctx, container.ListOptions{All: true}).Return([]types.Container{}, nil)
        mockController.On("ImageList", ctx, image.ListOptions{}).Return([]image.Summary{}, nil)
        mockController.On("ImagePull", ctx, containerInfo.Image, image.PullOptions{}).Return(io.NopCloser(strings.NewReader("")), nil)
        mockController.On("ContainerCreate", ctx, mock.Anything, mock.Anything, mock.Anything, nil, containerInfo.ContainerName).Return(container.CreateResponse{}, errors.New("failed to create container"))

        err:= createAndRunContainer(ctx, mockController, log, containerInfo)
        mockController.AssertExpectations(t)

        if err == nil {
            t.Fatalf("Expected error, got nil")
        }
        if err.Error() != "failed to create container" {
            t.Fatalf("Expected error 'failed to create container', got %v", err)
        }
    })

    t.Run("Fails to start container", func(t *testing.T) {
        mockController := new(MockContainerController)
        containerInfo := &general.Containers{
            ContainerName: "test-container",
            Image:         "test-image",
        }

        mockController.On("ContainerList", ctx, container.ListOptions{All: true}).Return([]types.Container{}, nil)
        mockController.On("ImageList", ctx, image.ListOptions{}).Return([]image.Summary{}, nil)
        mockController.On("ImagePull", ctx, containerInfo.Image, image.PullOptions{}).Return(io.NopCloser(strings.NewReader("")), nil)
        mockController.On("ContainerCreate", ctx, mock.Anything, mock.Anything, mock.Anything, nil, containerInfo.ContainerName).Return(container.CreateResponse{ID: "test-id"}, nil)
        mockController.On("ContainerStart", ctx, "test-id", container.StartOptions{}).Return(errors.New("failed to start container"))

        err := createAndRunContainer(ctx, mockController, log, containerInfo)
        mockController.AssertExpectations(t)

        if err == nil {
            t.Fatalf("Expected error, got nil")
        }
        if err.Error() != "failed to start container" {
            t.Fatalf("Expected error 'failed to start container', got %v", err)
        }
    })
}
func TestDoesContainerExist(t *testing.T) {
    ctx := context.Background()

    t.Run("Container exists and is exited", func(t *testing.T) {
        mockController := new(MockContainerController)
        containerName := "test-container"
        mockController.On("ContainerList", ctx, container.ListOptions{All: true}).Return([]types.Container{
            {
                Names: []string{"/" + containerName},
                State: "exited",
            },
        }, nil)

        err, exists := doesContainerExist(ctx, mockController, containerName)
        mockController.AssertExpectations(t)

        if err != nil {
            t.Fatalf("Expected no error, got %v", err)
        }
        if !exists {
            t.Fatalf("Expected container to exist, got %v", exists)
        }
    })

    t.Run("Container exists and is running", func(t *testing.T) {
        mockController := new(MockContainerController)
        containerName := "test-container"
        mockController.On("ContainerList", ctx, container.ListOptions{All: true}).Return([]types.Container{
            {
                Names: []string{"/" + containerName},
                State: "running",
            },
        }, nil)

        err, exists := doesContainerExist(ctx, mockController, containerName)
        mockController.AssertExpectations(t)

        if err != nil {
            t.Fatalf("Expected no error, got %v", err)
        }
        if exists {
            t.Fatalf("Expected container to not be considered 'exited', got %v", exists)
        }
    })

    t.Run("Container does not exist", func(t *testing.T) {
        mockController := new(MockContainerController)
        containerName := "non-existent-container"
        mockController.On("ContainerList", ctx, container.ListOptions{All: true}).Return([]types.Container{}, nil)

        err, exists := doesContainerExist(ctx, mockController, containerName)
        mockController.AssertExpectations(t)

        if err != nil {
            t.Fatalf("Expected no error, got %v", err)
        }
        if exists {
            t.Fatalf("Expected container to not exist, got %v", exists)
        }
    })
}
func TestPullImageIfNotExists(t *testing.T) {
    ctx := context.Background()
    log := logr.Discard()

    t.Run("Image is already pulled", func(t *testing.T) {
        mockController := new(MockContainerController)
        imageName := "test-image"

        mockController.On("ImageList", ctx, image.ListOptions{}).Return([]image.Summary{
            {RepoTags: []string{imageName}},
        }, nil)

        err := pullImageIfNotExists(ctx, mockController, imageName, log)
        mockController.AssertExpectations(t)

        if err != nil {
            t.Fatalf("Expected no error, got %v", err)
        }
    })

    t.Run("Image is not pulled and successfully pulled", func(t *testing.T) {
        mockController := new(MockContainerController)
        imageName := "test-image"

        mockController.On("ImageList", ctx, image.ListOptions{}).Return([]image.Summary{}, nil)
        mockController.On("ImagePull", ctx, imageName, image.PullOptions{}).Return(io.NopCloser(strings.NewReader("pulled")), nil)

        err := pullImageIfNotExists(ctx, mockController, imageName, log)
        mockController.AssertExpectations(t)

        if err != nil {
            t.Fatalf("Expected no error, got %v", err)
        }
    })

    t.Run("Error while reading pull output", func(t *testing.T) {
        mockController := new(MockContainerController)
        imageName := "test-image"
        expectedError := errors.New("failed to read output")

        mockController.On("ImageList", ctx, image.ListOptions{}).Return([]image.Summary{}, nil)
        mockController.On("ImagePull", ctx, imageName, image.PullOptions{}).Return(io.NopCloser(&errorReader{err: expectedError}), nil)

        err := pullImageIfNotExists(ctx, mockController, imageName, log)
        mockController.AssertExpectations(t)

        if err == nil {
            t.Fatalf("Expected error, got nil")
        }
        if err != expectedError {
            t.Fatalf("Expected error %v, got %v", expectedError, err)
        }
    })
}
func TestIsEmptyStruct(t *testing.T) {

    t.Run("Empty struct", func(t *testing.T) {
        type EmptyStruct struct{}
        empty := EmptyStruct{}

        if !IsEmptyStruct(empty) {
            t.Fatalf("Expected struct to be empty, but it was not")
        }
    })

    t.Run("Non-empty struct", func(t *testing.T) {
        type NonEmptyStruct struct {
            Field string
        }
        nonEmpty := NonEmptyStruct{Field: "value"}

        if IsEmptyStruct(nonEmpty) {
            t.Fatalf("Expected struct to be non-empty, but it was empty")
        }
    })
}
// errorReader is a helper struct to simulate an error during io.Reader operations
func TestPrepareContainerMediaProxyAgent(t *testing.T) {
    t.Run("Valid configuration", func(t *testing.T) {
        config := &utils.Configuration{
            RunOnce: utils.RunOnce{
                MediaProxyAgent: workloads.MediaProxyAgentConfig{
                    ImageAndTag: "test-image:latest",
                    GRPCPort:    "50051",
                    RestPort:    "8080",
                    Network: workloads.NetworkConfig{
                        Enable: true,
                        Name:   "test-network",
                        IP:     "192.168.1.1",
                    },
                },
            },
        }

        expected := general.Containers{
            Type:          general.MediaProxyAgent,
            ContainerName: MediaProxyAgentContainerName,
            Image:         "test-image:latest",
            Configuration: general.ContainersConfig{
                MediaProxyAgentConfig: workloads.MediaProxyAgentConfig{
                    ImageAndTag: "test-image:latest",
                    GRPCPort:    "50051",
                    RestPort:    "8080",
                    Network: workloads.NetworkConfig{
                        Enable: true,
                        Name:   "test-network",
                        IP:     "192.168.1.1",
                    },
                },
            },
        }

        result := prepareContainerMediaProxyAgent(config)

        if !reflect.DeepEqual(result, expected) {
            t.Fatalf("Expected %+v, got %+v", expected, result)
        }
    })

    t.Run("Empty configuration", func(t *testing.T) {
        config := &utils.Configuration{}

        expected := general.Containers{
            Type:          general.MediaProxyAgent,
            ContainerName: MediaProxyAgentContainerName,
        }

        result := prepareContainerMediaProxyAgent(config)

        if !reflect.DeepEqual(result, expected) {
            t.Fatalf("Expected %+v, got %+v", expected, result)
        }
    })
}
func TestPrepareContainerMediaProxyMcm(t *testing.T) {
    t.Run("Valid configuration", func(t *testing.T) {
        config := &utils.Configuration{
            RunOnce: utils.RunOnce{
                MediaProxyMcm: workloads.MediaProxyMcmConfig{
                    ImageAndTag:   "test-image:latest",
                    InterfaceName: "eth0",
                    Volumes: []string{"/dev/vfio:/dev/vfio"},
                    Network: workloads.NetworkConfig{
                        Enable: true,
                        Name:   "test-network",
                        IP:     "192.168.1.1",
                    },
                },
            },
        }

        expected := general.Containers{
            Type:          general.MediaProxyMCM,
            ContainerName: MediaProxyContainerName,
            Image:         "test-image:latest",
            Configuration: general.ContainersConfig{
                MediaProxyMcmConfig: workloads.MediaProxyMcmConfig{
                    ImageAndTag:   "test-image:latest",
                    InterfaceName: "eth0",
                    Volumes: []string{"/dev/vfio:/dev/vfio"},
                    Network: workloads.NetworkConfig{
                        Enable: true,
                        Name:   "test-network",
                        IP:     "192.168.1.1",
                    },
                },
            },
        }

        result := prepareContainerMediaProxyMcm(config)

        if !reflect.DeepEqual(result, expected) {
            t.Fatalf("Expected %+v, got %+v", expected, result)
        }
    })

    t.Run("Empty configuration", func(t *testing.T) {
        config := &utils.Configuration{}
        expected := general.Containers{
            Type:          general.MediaProxyMCM,
            ContainerName: MediaProxyContainerName,
        }
        result := prepareContainerMediaProxyMcm(config)
        if !reflect.DeepEqual(result, expected) {
            t.Fatalf("Expected %+v, got %+v", expected, result)
        }
    })
}
func TestPrepareContainerNmosClient(t *testing.T) {
    t.Run("Valid configuration", func(t *testing.T) {
        config := &utils.Configuration{
            WorkloadToBeRun: utils.WorkloadConfig{
                NmosClient: workloads.NmosClientConfig{
                    Name:                  "nmos-client",
                    ImageAndTag:           "test-image:latest",
                    EnvironmentVariables:  []string{"ENV_VAR1=value1", "ENV_VAR2=value2"},
                    NmosConfigPath:        "/path/to/config",
                    NmosConfigFileName:    "config.json",
                    Network: workloads.NetworkConfig{
                        Enable: true,
                        Name:   "test-network",
                        IP:     "192.168.1.1",
                    },
                },
                FfmpegPipeline: workloads.FfmpegPipelineConfig{
                    Network: workloads.NetworkConfig{
                        IP: "192.168.1.2",
                    },
                    GRPCPort: 50051,
                },
            },
        }

        expected := general.Containers{
            Type:          general.BcsPipelineNmosClient,
            ContainerName: "nmos-client",
            Image:         "test-image:latest",
            Configuration: general.ContainersConfig{
                WorkloadConfig: workloads.WorkloadConfig{
                    NmosClient: workloads.NmosClientConfig{
                        ImageAndTag:           "test-image:latest",
                        EnvironmentVariables:  []string{"ENV_VAR1=value1", "ENV_VAR2=value2"},
                        NmosConfigPath:        "/path/to/config",
                        NmosConfigFileName:    "config.json",
                        Network: workloads.NetworkConfig{
                            Enable: true,
                            Name:   "test-network",
                            IP:     "192.168.1.1",
                        },
                        FfmpegConectionAddress: "192.168.1.2",
                        FfmpegConnectionPort:   "50051",
                    },
                },
            },
        }

        result := prepareContainerNmosClient(config)

        if !reflect.DeepEqual(result, expected) {
            t.Fatalf("Expected %+v, got %+v", expected, result)
        }
    })
}
func TestPrepareContainerFfmpegPipeline(t *testing.T) {
    t.Run("Valid configuration", func(t *testing.T) {
        config := &utils.Configuration{
            WorkloadToBeRun: utils.WorkloadConfig{
                FfmpegPipeline: workloads.FfmpegPipelineConfig{
                    Name:        "ffmpeg-pipeline",
                    ImageAndTag: "test-image:latest",
                    GRPCPort:    50051,
                    SourcePort:  8080,
                    EnvironmentVariables: []string{"ENV_VAR1=value1", "ENV_VAR2=value2"},
                    Volumes: workloads.Volumes{
                        Videos:      "/path/to/videos",
                        Dri:         "/path/to/dri",
                        Kahawai:     "/path/to/kahawai",
                        Devnull:     "/path/to/devnull",
                        TmpHugepages: "/path/to/tmpHugepages",
                        Hugepages:   "/path/to/hugepages",
                        Imtl:        "/path/to/imtl",
                        Shm:         "/path/to/shm",
                    },
                    Devices: workloads.Devices{
                        Vfio: "/dev/vfio",
                        Dri:  "/dev/dri",
                    },
                    Network: workloads.NetworkConfig{
                        Enable: true,
                        Name:   "test-network",
                        IP:     "192.168.1.1",
                    },
                },
            },
        }

        expected := general.Containers{
            Type:          general.BcsPipelineFfmpeg,
            ContainerName: "ffmpeg-pipeline",
            Image:         "test-image:latest",
            Configuration: general.ContainersConfig{
                WorkloadConfig: workloads.WorkloadConfig{
                    FfmpegPipeline: workloads.FfmpegPipelineConfig{
                        Name:        "ffmpeg-pipeline",
                        ImageAndTag: "test-image:latest",
                        GRPCPort:    50051,
                        SourcePort:  8080,
                        EnvironmentVariables: []string{"ENV_VAR1=value1", "ENV_VAR2=value2"},
                        Volumes: workloads.Volumes{
                            Videos:      "/path/to/videos",
                            Dri:         "/path/to/dri",
                            Kahawai:     "/path/to/kahawai",
                            Devnull:     "/path/to/devnull",
                            TmpHugepages: "/path/to/tmpHugepages",
                            Hugepages:   "/path/to/hugepages",
                            Imtl:        "/path/to/imtl",
                            Shm:         "/path/to/shm",
                        },
                        Devices: workloads.Devices{
                            Vfio: "/dev/vfio",
                            Dri:  "/dev/dri",
                        },
                        Network: workloads.NetworkConfig{
                            Enable: true,
                            Name:   "test-network",
                            IP:     "192.168.1.1",
                        },
                    },
                },
            },
        }

        result := prepareContainerFfmpegPipeline(config)

        if !reflect.DeepEqual(result, expected) {
            t.Fatalf("Expected %+v, got %+v", expected, result)
        }
    })

    t.Run("Empty configuration", func(t *testing.T) {
        config := &utils.Configuration{}
        expected := general.Containers{
            Type:          general.BcsPipelineFfmpeg,
            ContainerName: "",
        }
        result := prepareContainerFfmpegPipeline(config)
        if !reflect.DeepEqual(result, expected) {
            t.Fatalf("Expected %+v, got %+v", expected, result)
        }
    })
}
func TestDockerContainerController(t *testing.T) {
	ctx := context.Background()

	t.Run("ImageList returns images successfully", func(t *testing.T) {
		mockClient := new(MockContainerController)
		expectedImages := []image.Summary{
			{RepoTags: []string{"test-image:latest"}},
		}
		mockClient.On("ImageList", ctx, image.ListOptions{}).Return(expectedImages, nil)

		images, err := mockClient.ImageList(ctx, image.ListOptions{})
		mockClient.AssertExpectations(t)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !reflect.DeepEqual(images, expectedImages) {
			t.Fatalf("Expected images %+v, got %+v", expectedImages, images)
		}
	})

	t.Run("ImagePull pulls image successfully", func(t *testing.T) {
		mockClient := new(MockContainerController)
		imageName := "test-image:latest"
		mockClient.On("ImagePull", ctx, imageName, image.PullOptions{}).Return(io.NopCloser(strings.NewReader("pulled")), nil)

		reader, err := mockClient.ImagePull(ctx, imageName, image.PullOptions{})
		mockClient.AssertExpectations(t)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		defer reader.Close()
		output, _ := io.ReadAll(reader)
		if string(output) != "pulled" {
			t.Fatalf("Expected output 'pulled', got %s", string(output))
		}
	})

	t.Run("ContainerList returns containers successfully", func(t *testing.T) {
		mockClient := new(MockContainerController)
		expectedContainers := []types.Container{
			{Names: []string{"/test-container"}, State: "running"},
		}
		mockClient.On("ContainerList", ctx, container.ListOptions{All: true}).Return(expectedContainers, nil)

		containers, err := mockClient.ContainerList(ctx, container.ListOptions{All: true})
		mockClient.AssertExpectations(t)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !reflect.DeepEqual(containers, expectedContainers) {
			t.Fatalf("Expected containers %+v, got %+v", expectedContainers, containers)
		}
	})

	t.Run("ContainerCreate creates container successfully", func(t *testing.T) {
		mockClient := new(MockContainerController)
		containerName := "test-container"
		expectedResponse := container.CreateResponse{ID: "test-id"}
		mockClient.On("ContainerCreate", ctx, mock.Anything, mock.Anything, mock.Anything, nil, containerName).Return(expectedResponse, nil)

		resp, err := mockClient.ContainerCreate(ctx, &container.Config{}, &container.HostConfig{}, &network.NetworkingConfig{}, nil, containerName)
		mockClient.AssertExpectations(t)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if resp.ID != expectedResponse.ID {
			t.Fatalf("Expected container ID %s, got %s", expectedResponse.ID, resp.ID)
		}
	})

	t.Run("ContainerStart starts container successfully", func(t *testing.T) {
		mockClient := new(MockContainerController)
		containerID := "test-id"
		mockClient.On("ContainerStart", ctx, containerID, container.StartOptions{}).Return(nil)

		err := mockClient.ContainerStart(ctx, containerID, container.StartOptions{})
		mockClient.AssertExpectations(t)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("ContainerRemove removes container successfully", func(t *testing.T) {
		mockClient := new(MockContainerController)
		containerID := "test-id"
		mockClient.On("ContainerRemove", ctx, containerID, container.RemoveOptions{Force: true}).Return(nil)

		err := mockClient.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
		mockClient.AssertExpectations(t)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}
func TestNewDockerContainerController(t *testing.T) {
	t.Run("Successfully creates DockerContainerController", func(t *testing.T) {
		controller, err := NewDockerContainerController()

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if controller == nil {
			t.Fatalf("Expected a valid DockerContainerController, got nil")
		}
		if controller.cli == nil {
			t.Fatalf("Expected a valid Docker client, got nil")
		}
	})

	t.Run("Fails to create DockerContainerController due to invalid client options", func(t *testing.T) {
		// Temporarily override the environment variable to simulate an error
		originalEnv := os.Getenv("DOCKER_HOST")
		os.Setenv("DOCKER_HOST", "invalid-host")

		defer os.Setenv("DOCKER_HOST", originalEnv)

		controller, err := NewDockerContainerController()

		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
		if controller != nil {
			t.Fatalf("Expected nil DockerContainerController, got %v", controller)
		}
	})
}

type errorReader struct {
    err error
}

func (e *errorReader) Read(p []byte) (n int, err error) {
    return 0, e.err
}

func (e *errorReader) Close() error {
    return nil
}

