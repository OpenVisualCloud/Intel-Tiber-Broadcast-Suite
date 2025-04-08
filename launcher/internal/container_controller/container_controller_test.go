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
	"strings"
	"testing"

	"bcs.pod.launcher.intel/resources_library/resources/general"
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

type errorReader struct {
	err error
}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}

func (e *errorReader) Close() error {
	return nil
}