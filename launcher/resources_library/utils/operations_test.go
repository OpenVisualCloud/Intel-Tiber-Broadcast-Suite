package utils

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// // MockDockerClient is a mock implementation of the Docker client
// type MockDockerClient struct {
// 	mock.Mock
// }

// func (m *MockDockerClient) ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
// 	args := m.Called(ctx, options)
// 	return args.Get(0).([]types.ImageSummary), args.Error(1)
// }

// func (m *MockDockerClient) ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error) {
// 	args := m.Called(ctx, ref, options)
// 	return args.Get(0).(io.ReadCloser), args.Error(1)
// }

// func (m *MockDockerClient) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
// 	args := m.Called(ctx, options)
// 	return args.Get(0).([]types.Container), args.Error(1)
// }

// func (m *MockDockerClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *types.Platform, containerName string) (container.CreateResponse, error) {
// 	args := m.Called(ctx, config, hostConfig, networkingConfig, platform, containerName)
// 	return args.Get(0).(container.CreateResponse), args.Error(1)
// }

// func (m *MockDockerClient) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error {
// 	args := m.Called(ctx, containerID, options)
// 	return args.Error(0)
// }

// func (m *MockDockerClient) ContainerRemove(ctx context.Context, containerID string, options types.ContainerRemoveOptions) error {
// 	args := m.Called(ctx, containerID, options)
// 	return args.Error(0)
// }

// func TestIsImagePulled(t *testing.T) {
// 	mockClient := new(MockDockerClient)
// 	ctx := context.Background()

// 	mockClient.On("ImageList", ctx, types.ImageListOptions{}).Return([]types.ImageSummary{
// 		{RepoTags: []string{"test-image:latest"}},
// 	}, nil)

// 	err, isPulled := isImagePulled(ctx, mockClient, "test-image:latest")
// 	assert.NoError(t, err)
// 	assert.True(t, isPulled)

// 	err, isPulled = isImagePulled(ctx, mockClient, "non-existent-image:latest")
// 	assert.NoError(t, err)
// 	assert.False(t, isPulled)
// }

// func TestPullImageIfNotExists(t *testing.T) {
// 	mockClient := new(MockDockerClient)
// 	ctx := context.Background()
// 	log := logr.Discard()

// 	mockClient.On("ImageList", ctx, types.ImageListOptions{}).Return([]types.ImageSummary{
// 		{RepoTags: []string{"existing-image:latest"}},
// 	}, nil)

// 	mockClient.On("ImagePull", ctx, "new-image:latest", types.ImagePullOptions{}).Return(io.NopCloser(nil), nil)

// 	// Test when image already exists
// 	err := pullImageIfNotExists(ctx, mockClient, "existing-image:latest", log)
// 	assert.NoError(t, err)

// 	// Test when image does not exist
// 	err = pullImageIfNotExists(ctx, mockClient, "new-image:latest", log)
// 	assert.NoError(t, err)
// }

// func TestDoesContainerExist(t *testing.T) {
// 	mockClient := new(MockDockerClient)
// 	ctx := context.Background()

// 	mockClient.On("ContainerList", ctx, types.ContainerListOptions{All: true}).Return([]types.Container{
// 		{Names: []string{"/test-container"}, State: "exited"},
// 	}, nil)

// 	err, exists := doesContainerExist(ctx, mockClient, "test-container")
// 	assert.NoError(t, err)
// 	assert.True(t, exists)

// 	err, exists = doesContainerExist(ctx, mockClient, "non-existent-container")
// 	assert.NoError(t, err)
// 	assert.False(t, exists)
// }

// func TestIsContainerRunning(t *testing.T) {
// 	mockClient := new(MockDockerClient)
// 	ctx := context.Background()

// 	mockClient.On("ContainerList", ctx, types.ContainerListOptions{All: true}).Return([]types.Container{
// 		{Names: []string{"/running-container"}, State: "running"},
// 	}, nil)

// 	err, isRunning := isContainerRunning(ctx, mockClient, "running-container")
// 	assert.NoError(t, err)
// 	assert.True(t, isRunning)

// 	err, isRunning = isContainerRunning(ctx, mockClient, "stopped-container")
// 	assert.NoError(t, err)
// 	assert.False(t, isRunning)
// }

func TestUpdateNmosJsonFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "nmos_test_*.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	jsonData := `{
        "ffmpegGrpcServerAddress": "old-address",
        "ffmpegGrpcServerPort": "old-port"
    }`
	_, err = tempFile.Write([]byte(jsonData))
	assert.NoError(t, err)
	tempFile.Close()

	err = updateNmosJsonFile(tempFile.Name(), "new-address", "new-port")
	assert.NoError(t, err)

	updatedData, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err)

	var config map[string]string
	err = json.Unmarshal(updatedData, &config)
	assert.NoError(t, err)
	assert.Equal(t, "new-address", config["ffmpegGrpcServerAddress"])
	assert.Equal(t, "new-port", config["ffmpegGrpcServerPort"])
}

func TestFileExists(t *testing.T) {
	tempFile, err := os.CreateTemp("", "file_exists_test_*.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	assert.True(t, FileExists(tempFile.Name()))
	assert.False(t, FileExists("non-existent-file.txt"))
}
