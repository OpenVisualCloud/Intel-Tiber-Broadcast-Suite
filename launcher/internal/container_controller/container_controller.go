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
	"sync"

	"bcs.pod.launcher.intel/resources_library/parser"
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

var fileMutex sync.Mutex

const (
	MediaProxyAgentContainerName = "mesh-agent"
	MediaProxyContainerName      = "media-proxy"
	BCSPipelineContainerName     = "bcs-ffmpeg-pipeline"
)

type ContainerController interface {
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

func CreateAndRunContainers(ctx context.Context, cli ContainerController, log logr.Logger, config *parser.Configuration) error {

	//pass the yaml configuration to the Container struct
	if !IsEmptyStruct(config.RunOnce.MediaProxyAgent) {

		mcmAgentContainer := general.Containers{}
		mcmAgentContainer.Type = general.MediaProxyAgent
		mcmAgentContainer.ContainerName = MediaProxyAgentContainerName
		mcmAgentContainer.Image = config.RunOnce.MediaProxyAgent.ImageAndTag
		if !config.RunOnce.MediaProxyAgent.Network.Enable {
			config.RunOnce.MediaProxyAgent.Network.Name = "host"
			// Note - When you use the host network mode in Docker, the container shares
			// the host machine's network stack.
		}
		err := createAndRunContainer(ctx, cli, log, &mcmAgentContainer, config)
		if err != nil {
			log.Error(err, "Failed to create container MCM MediaProxy Agent!")
			return err
		}
	} else {
		log.Info("No information about MCM MediaProxy Agent provided. Omitting creation of MCM MediaProxy Agent container")
	}

	if !IsEmptyStruct(config.RunOnce.MediaProxyMcm) {
		mediaProxyContainer := general.Containers{}
		mediaProxyContainer.Type = general.MediaProxyMCM
		mediaProxyContainer.ContainerName = MediaProxyContainerName
		mediaProxyContainer.Image = config.RunOnce.MediaProxyMcm.ImageAndTag
		if !config.RunOnce.MediaProxyMcm.Network.Enable {
			config.RunOnce.MediaProxyMcm.Network.Name = "host"
			// Note - When you use the host network mode in Docker, the container shares
			// the host machine's network stack.
		}
		err := createAndRunContainer(ctx, cli, log, &mediaProxyContainer, config)
		if err != nil {
			log.Error(err, "Failed to create container MCM MediaProxy!")
			return err
		}
	} else {
		log.Info("No information about MCM MediaProxy provided. Omitting creation of MCM MediaProxy container")
	}

	if IsEmptyStruct(config.WorkloadToBeRun) || len(config.WorkloadToBeRun) == 0 {
		log.Info("No workloads provided under workloadToBeRun. Omitting creation of BCS pipeline and NMOS node containers")
	}

	for n, instance := range config.WorkloadToBeRun {
		if IsEmptyStruct(instance.FfmpegPipeline) || IsEmptyStruct(instance.NmosClient) {
			return fmt.Errorf("no information about BCS pipeline provided. Either FfmpegPipeline or NmosClient is empty for instance Ffmpeg: %s; Nmos: %s", instance.FfmpegPipeline.Name, instance.NmosClient.Name)
		}
		bcsPipelinesContainer := general.Containers{}
		bcsPipelinesContainer.Type = general.BcsPipelineFfmpeg
		bcsPipelinesContainer.ContainerName = instance.FfmpegPipeline.Name
		bcsPipelinesContainer.Image = instance.FfmpegPipeline.ImageAndTag
		bcsPipelinesContainer.Id = n // use the index of the instance as the ID for the container

		if !instance.FfmpegPipeline.Network.Enable {
			instance.FfmpegPipeline.Network.Name = "host"
			// Note - When you use the host network mode in Docker, the container shares
			// the host machine's network stack.
		}
		err := createAndRunContainer(ctx, cli, log, &bcsPipelinesContainer, config)
		if err != nil {
			log.Error(err, "Failed to create container for FFMPEG pipeline instance %d!", n)
			return err
		}
		bcsNmosContainer := general.Containers{}
		bcsNmosContainer.Type = general.BcsPipelineNmosClient
		bcsNmosContainer.ContainerName = instance.NmosClient.Name
		bcsNmosContainer.Image = instance.NmosClient.ImageAndTag
		bcsNmosContainer.Id = n // use the index of the instance as the ID for the container

		if !instance.NmosClient.Network.Enable {
			// do not forget to set the network ip address despite disabling the custom network!
			instance.NmosClient.Network.IP = "host"
			// Note - When you use the host network mode in Docker, the container shares
			// the host machine's network stack.
		}
		instance.NmosClient.FfmpegConnectionAddress = instance.FfmpegPipeline.Network.IP
		instance.NmosClient.FfmpegConnectionPort = strconv.Itoa(instance.FfmpegPipeline.GRPCPort)
		err = createAndRunContainer(ctx, cli, log, &bcsNmosContainer, config)
		if err != nil {
			log.Error(err, "Failed to create container!")
			return err
		}

	}
	return nil
}

func createAndRunContainer(ctx context.Context, cli ContainerController, log logr.Logger, containerInfo *general.Containers, config *parser.Configuration) error {
	err, isRunning := isContainerRunning(ctx, cli, containerInfo.ContainerName)
	if err != nil {
		log.Error(err, "Failed to parse launcher configuration file")
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
	containerConfig, hostConfig, networkConfig := utils.ConstructContainerConfig(containerInfo, config, log)

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
