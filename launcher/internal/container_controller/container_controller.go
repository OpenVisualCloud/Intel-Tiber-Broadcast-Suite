/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

package containercontroller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"bcs.pod.launcher.intel/resources_library/resources/general"
	"bcs.pod.launcher.intel/resources_library/utils"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/go-logr/logr"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)
const (
	MediaProxyAgentContainerName = "mesh-agent"
	MediaProxyContainerName      = "media-proxy"
	BCSPipelineContainerName     = "bcs-ffmpeg-pipeline"
)

type ContainerController interface {
	// CreateAndRunContainers(ctx context.Context, launcherConfigName string, log logr.Logger) error
	ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error)
    ImagePull(ctx context.Context, ref string, options image.PullOptions) (io.ReadCloser, error)
	ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error
	ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error)
	ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) (container.CreateResponse, error)
}

type DockerContainerController struct {
	cli *client.Client
}

func (d *DockerContainerController) ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error) {
	return d.cli.ImageList(ctx, options)
}

func (d *DockerContainerController) ImagePull(ctx context.Context, ref string, options image.PullOptions) (io.ReadCloser, error) {
	return d.cli.ImagePull(ctx, ref, options)
}
func (d *DockerContainerController) ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
	return d.cli.ContainerList(ctx, options)
}
func (d *DockerContainerController) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) (container.CreateResponse, error) {
	return d.cli.ContainerCreate(ctx, config, hostConfig, networkingConfig, platform, containerName)
}
func (d *DockerContainerController) ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error {
	return d.cli.ContainerStart(ctx, containerID, options)
}
func (d *DockerContainerController) ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error {
	return d.cli.ContainerRemove(ctx, containerID, options)
}

func NewDockerContainerController() (*DockerContainerController, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerContainerController{cli: cli}, nil
}

func IsEmptyStruct(s interface{}) bool {
	return reflect.DeepEqual(s, reflect.Zero(reflect.TypeOf(s)).Interface())
}

// Use case covers running containers on single host
// CreateAndRunContainers creates and runs Docker containers based on the provided launcher configuration.
// It parses the launcher configuration file, checks for the presence of specific container configurations,
// and creates and runs the containers accordingly.
//
// Parameters:
//   - ctx: The context for managing the lifecycle of the container creation process.
//   - launcherConfigName: The name of the launcher configuration file to be parsed.
//   - log: The logger for logging errors and information.
//
// Returns:
//   - error: An error if any step in the container creation process fails, otherwise nil.
//
// The function performs the following steps:
//   1. Parses the launcher configuration file.
//   2. Checks if the configuration is empty and logs an error if it is.
//   3. Creates and runs the MCM MediaProxy Agent container if its configuration is provided.
//   4. Creates and runs the MCM MediaProxy container if its configuration is provided.
//   5. Creates and runs the BCS NMOS client container if its configuration is provided.
//   6. Creates and runs the BCS FFmpeg pipeline container with predefined settings.

func prepareContainerMediaProxyAgent(config *utils.Configuration)  general.Containers{
	containerApp := general.Containers{}
	
		containerApp.Type = general.MediaProxyAgent
		containerApp.ContainerName = MediaProxyAgentContainerName
		containerApp.Image = config.RunOnce.MediaProxyAgent.ImageAndTag
		containerApp.Configuration.MediaProxyAgentConfig.ImageAndTag = config.RunOnce.MediaProxyAgent.ImageAndTag
		containerApp.Configuration.MediaProxyAgentConfig.GRPCPort = config.RunOnce.MediaProxyAgent.GRPCPort
		containerApp.Configuration.MediaProxyAgentConfig.RestPort = config.RunOnce.MediaProxyAgent.RestPort
		containerApp.Configuration.MediaProxyAgentConfig.Network.Enable = config.RunOnce.MediaProxyAgent.Network.Enable
		containerApp.Configuration.MediaProxyAgentConfig.Network.Name = config.RunOnce.MediaProxyAgent.Network.Name
		containerApp.Configuration.MediaProxyAgentConfig.Network.IP = config.RunOnce.MediaProxyAgent.Network.IP
	
	return containerApp
}

func prepareContainerMediaProxyMcm(config *utils.Configuration) general.Containers {
	containerApp := general.Containers{}
		containerApp.Type = general.MediaProxyMCM
		containerApp.ContainerName = MediaProxyContainerName
		containerApp.Image = config.RunOnce.MediaProxyMcm.ImageAndTag
		containerApp.Configuration.MediaProxyMcmConfig.ImageAndTag = config.RunOnce.MediaProxyMcm.ImageAndTag
		containerApp.Configuration.MediaProxyMcmConfig.InterfaceName = config.RunOnce.MediaProxyMcm.InterfaceName
		containerApp.Configuration.MediaProxyMcmConfig.Volumes = config.RunOnce.MediaProxyMcm.Volumes
		containerApp.Configuration.MediaProxyMcmConfig.Network.Enable = config.RunOnce.MediaProxyMcm.Network.Enable
		containerApp.Configuration.MediaProxyMcmConfig.Network.Name = config.RunOnce.MediaProxyMcm.Network.Name
		containerApp.Configuration.MediaProxyMcmConfig.Network.IP = config.RunOnce.MediaProxyMcm.Network.IP
	return containerApp
}

func prepareContainerNmosClient(config *utils.Configuration) general.Containers {
	containerApp := general.Containers{}
		containerApp.Type = general.BcsPipelineNmosClient
		containerApp.ContainerName = config.WorkloadToBeRun.NmosClient.Name
		containerApp.Image = config.WorkloadToBeRun.NmosClient.ImageAndTag
		containerApp.Configuration.WorkloadConfig.NmosClient.ImageAndTag = config.WorkloadToBeRun.NmosClient.ImageAndTag
		containerApp.Configuration.WorkloadConfig.NmosClient.EnvironmentVariables = config.WorkloadToBeRun.NmosClient.EnvironmentVariables
		containerApp.Configuration.WorkloadConfig.NmosClient.NmosConfigPath = config.WorkloadToBeRun.NmosClient.NmosConfigPath
		containerApp.Configuration.WorkloadConfig.NmosClient.NmosConfigFileName = config.WorkloadToBeRun.NmosClient.NmosConfigFileName
		containerApp.Configuration.WorkloadConfig.NmosClient.Network.Enable = config.WorkloadToBeRun.NmosClient.Network.Enable
		containerApp.Configuration.WorkloadConfig.NmosClient.Network.Name = config.WorkloadToBeRun.NmosClient.Network.Name
		containerApp.Configuration.WorkloadConfig.NmosClient.Network.IP = config.WorkloadToBeRun.NmosClient.Network.IP
		containerApp.Configuration.WorkloadConfig.NmosClient.FfmpegConectionAddress = config.WorkloadToBeRun.FfmpegPipeline.Network.IP
		containerApp.Configuration.WorkloadConfig.NmosClient.FfmpegConnectionPort = strconv.Itoa(config.WorkloadToBeRun.FfmpegPipeline.GRPCPort)
	return containerApp
}

func prepareContainerFfmpegPipeline(config *utils.Configuration) general.Containers {
	containerApp := general.Containers{}
	containerApp.Type = general.BcsPipelineFfmpeg
	containerApp.ContainerName = config.WorkloadToBeRun.FfmpegPipeline.Name
	containerApp.Image = config.WorkloadToBeRun.FfmpegPipeline.ImageAndTag
	containerApp.Configuration.WorkloadConfig.FfmpegPipeline.Name = config.WorkloadToBeRun.FfmpegPipeline.Name
	containerApp.Configuration.WorkloadConfig.FfmpegPipeline.ImageAndTag = config.WorkloadToBeRun.FfmpegPipeline.ImageAndTag
	containerApp.Configuration.WorkloadConfig.FfmpegPipeline.GRPCPort = config.WorkloadToBeRun.FfmpegPipeline.GRPCPort
	containerApp.Configuration.WorkloadConfig.FfmpegPipeline.SourcePort = config.WorkloadToBeRun.FfmpegPipeline.SourcePort
	containerApp.Configuration.WorkloadConfig.FfmpegPipeline.EnvironmentVariables = config.WorkloadToBeRun.FfmpegPipeline.EnvironmentVariables
	containerApp.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Devnull = config.WorkloadToBeRun.FfmpegPipeline.Volumes.Devnull
	containerApp.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Dri = config.WorkloadToBeRun.FfmpegPipeline.Volumes.Dri
	containerApp.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Hugepages = config.WorkloadToBeRun.FfmpegPipeline.Volumes.Hugepages
	containerApp.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Imtl = config.WorkloadToBeRun.FfmpegPipeline.Volumes.Imtl
	containerApp.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Kahawai = config.WorkloadToBeRun.FfmpegPipeline.Volumes.Kahawai
	containerApp.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Shm = config.WorkloadToBeRun.FfmpegPipeline.Volumes.Shm
	containerApp.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.TmpHugepages = config.WorkloadToBeRun.FfmpegPipeline.Volumes.TmpHugepages
	containerApp.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Videos = config.WorkloadToBeRun.FfmpegPipeline.Volumes.Videos
	containerApp.Configuration.WorkloadConfig.FfmpegPipeline.Devices.Dri = config.WorkloadToBeRun.FfmpegPipeline.Devices.Dri
	containerApp.Configuration.WorkloadConfig.FfmpegPipeline.Devices.Vfio = config.WorkloadToBeRun.FfmpegPipeline.Devices.Vfio
	containerApp.Configuration.WorkloadConfig.FfmpegPipeline.Network.Enable = config.WorkloadToBeRun.FfmpegPipeline.Network.Enable
	containerApp.Configuration.WorkloadConfig.FfmpegPipeline.Network.Name = config.WorkloadToBeRun.FfmpegPipeline.Network.Name
	containerApp.Configuration.WorkloadConfig.FfmpegPipeline.Network.IP = config.WorkloadToBeRun.FfmpegPipeline.Network.IP
	return containerApp
}

func CreateAndRunContainers(ctx context.Context, cli ContainerController, launcherConfigName string, log logr.Logger) error {
	config, err := utils.ParseLauncherConfiguration(launcherConfigName)
	if err != nil {
		log.Error(err, "Failed to parse launcher configuration file")
		return err
	}
	if IsEmptyStruct(config) {
		log.Error(err, "Failed to parse launcher configuration file. Configuration is empty")
		return err
	}
	//pass the yaml configuration to the Contaier struct
	if IsEmptyStruct(config.RunOnce.MediaProxyAgent) {
		mcmAgentContainer := prepareContainerMediaProxyAgent(&config)
		err := createAndRunContainer(ctx, cli, log, &mcmAgentContainer)
		if err != nil {
			log.Error(err, "Failed to create contianer MCM MediaProxy Agent!")
			return err
		}
	} else {
		log.Info("No information about MCM MediaProxy Agent provided. Omitting creation of MCM MediaProxy Agent container")
	}

	if IsEmptyStruct(config.RunOnce.MediaProxyMcm) {
		mediaProxyContainer := prepareContainerMediaProxyMcm(&config)
		err := createAndRunContainer(ctx, cli, log, &mediaProxyContainer)
		if err != nil {
			log.Error(err, "Failed to create contianer MCM MediaProxy!")
			return err
		}
	} else {
		log.Info("No information about MCM MediaProxy provided. Omitting creation of MCM MediaProxy container")
	}

	if IsEmptyStruct(config.WorkloadToBeRun.NmosClient) {
		bcsNmosContainer := prepareContainerNmosClient(&config)
		err = createAndRunContainer(ctx, cli, log, &bcsNmosContainer)
		if err != nil {
			log.Error(err, "Failed to create contianer!")
			return err
		}
	} else {
		log.Info("No information about BCS NMOS client container provided. Omitting creation of BCS NMOS client container")
	}

	if IsEmptyStruct(config.WorkloadToBeRun.FfmpegPipeline) {
		bcsPipelinesContainer := prepareContainerFfmpegPipeline(&config)
		err = createAndRunContainer(ctx, cli, log, &bcsPipelinesContainer)
	if err != nil {
		log.Error(err, "Failed to create contianer!")
		return err
	}
	} else {
		log.Info("No information about BCS ffmpeg container provided. Omitting creation of BCS NMOS client container")
	}
	return nil
}

func isImagePulled(ctx context.Context, cli ContainerController, imageName string) (error, bool) {
	images, err := cli.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return err, false
	}

	imageMap := make(map[string]bool)
	for _, image := range images {
		for _, tag := range image.RepoTags {
			imageMap[tag] = true
		}
	}

	_, isPulled := imageMap[imageName]
	return nil, isPulled
}

func pullImageIfNotExists(ctx context.Context, cli ContainerController, imageName string, log logr.Logger) error {

	// Check if the Docker client is nil
	if cli == nil {
		err := errors.New("docker client is nil")
		log.Error(err, "Docker client is not initialized")
		return err
	}

	// Check if the context is nil
	if ctx == nil {
		err := errors.New("context is nil")
		log.Error(err, "Context is not initialized")
		return err
	}

	// Check if the image is already pulled
	err, pulled := isImagePulled(ctx, cli, imageName)
	if err != nil {
		log.Error(err, "Error checking if image is pulled")
		return err
	}

	// Pull the image if it is not already pulled
	if !pulled {
		reader, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
		if err != nil {
			log.Error(err, "Error pulling image")
			return err
		}
		defer reader.Close()

		_, err = io.Copy(os.Stdout, reader)
		if err != nil {
			log.Error(err, "Error reading output")
			return err
		}
		log.Info("Image pulled successfully", "image", imageName)
	}

	return nil
}

func doesContainerExist(ctx context.Context, cli ContainerController, containerName string) (error, bool) {
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return err, false
	}

	containerMap := make(map[string]string)
	for _, container := range containers {
		for _, name := range container.Names {
			containerMap[name] = strings.ToLower(container.State)
		}
	}

	state, exists := containerMap["/"+containerName]
	if !exists {
		return nil, false
	}

	return nil, state == "exited"
}

func isContainerRunning(ctx context.Context, cli ContainerController, containerName string) (error, bool) {
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return err, false
	}

	containerMap := make(map[string]string)
	for _, container := range containers {
		for _, name := range container.Names {
			containerMap[name] = strings.ToLower(container.State)
		}
	}

	state, exists := containerMap["/"+containerName]
	if !exists {
		return nil, false
	}

	return nil, state == "running"
}

func removeContainer(ctx context.Context, cli ContainerController, containerID string) error {
	return cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
}

func createAndRunContainer(ctx context.Context, cli ContainerController, log logr.Logger, containerInfo *general.Containers) error {
	err, isRunning := isContainerRunning(ctx, cli, containerInfo.ContainerName)
	if err != nil {
		log.Error(err, "Failed to read container status (if it is in running state)")
		return err
	}

	if isRunning {
		log.Info("Container ", containerInfo.ContainerName, " is running. Omitting this container creation.")
		return nil
	}

	err, exists := doesContainerExist(ctx, cli, containerInfo.ContainerName)
	if err != nil {
		log.Error(err, "Failed to read container status (if it exists)")
		return err
	}

	if exists {
		log.Info("Removing container to re-create and re-run because container with a such name exists but with status exited:", "container", containerInfo.ContainerName)
		err = removeContainer(ctx, cli, containerInfo.ContainerName)
		if err != nil {
			log.Error(err, "Failed to remove container")
			return err
		}

	}

	err = pullImageIfNotExists(ctx, cli, containerInfo.Image, log)
	if err != nil {
		log.Error(err, "Error pulling image for container")
		return err
	}
	// Define the container configuration
	fmt.Println("--------------------------------------------------------------", containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.GRPCPort)

	containerConfig, hostConfig, networkConfig := utils.ConstructContainerConfig(containerInfo, log)

	if containerConfig == nil || hostConfig == nil || networkConfig == nil {
		// log.Error(errors.New("container configuration is nil"), "Failed to construct container configuration")
		return errors.New("container configuration is nil")
	}
	// Create the container
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, nil, containerInfo.ContainerName)

	if err != nil {
		log.Error(err, "Error creating container")
		return err
	}

	// Start the container
	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		log.Error(err, "Error starting container")
		return err
	}

	log.Info("Container is created and started successfully", "name", containerInfo.ContainerName, "container id: ", resp.ID)
	return nil
}