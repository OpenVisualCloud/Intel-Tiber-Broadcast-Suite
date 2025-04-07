/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

package containercontroller

import (
	"context"
	"reflect"
	"strconv"

	"bcs.pod.launcher.intel/resources_library/resources/general"
	"bcs.pod.launcher.intel/resources_library/utils"

	// "bcs.pod.launcher.intel/resources_library/utils"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/go-logr/logr"
)
const (
	MediaProxyAgentContainerName = "mesh-agent"
	MediaProxyContainerName      = "media-proxy"
	BCSPipelineContainerName     = "bcs-ffmpeg-pipeline"
)

type ContainerController interface {
	CreateAndRunContainers(ctx context.Context, launcherConfigName string, log logr.Logger) error
	IsContainerRunning(containerID string) (bool, error)
	// ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
	ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error)
	// ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error)
    // ImagePull(ctx context.Context, ref string, options image.PullOptions) (io.ReadCloser, error)
    // ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error)
	// // ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) (container.CreateResponse, error)
    // ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error
    // ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error
}

type DockerContainerController struct {
	cli *client.Client
}

func (d *DockerContainerController) ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error) {
	return d.cli.ImageList(ctx, options)
}

func NewDockerContainerController() (*DockerContainerController, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerContainerController{cli: cli}, nil
}

// func (d *DockerContainerController) ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error) {
//     return d.cli.ImageList(ctx, options)
// }

// func (d *DockerContainerController) ImagePull(ctx context.Context, ref string, options image.PullOptions) (io.ReadCloser, error) {
//     return d.cli.ImagePull(ctx, ref, options)
// }

// func (d *DockerContainerController) ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error){
//     return d.cli.ContainerList(ctx, options)
// }

// // func (d *DockerContainerController) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) (container.CreateResponse, error){
// // 	// Convert networkingConfig to the correct type if necessary
	
// // 	return d.cli.ContainerCreate(ctx, config, hostConfig, networkingConfig, nil, containerName)
// // }

// func (d *DockerContainerController) ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error{
//     return d.cli.ContainerStart(ctx, containerID, options)
// }

// func (d *DockerContainerController) ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error{
//     return d.cli.ContainerRemove(ctx, containerID, options)
// }

func (d *DockerContainerController) isEmptyStruct(s interface{}) bool {
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

func (d *DockerContainerController) s(ctx context.Context, launcherConfigName string, log logr.Logger) error {
	config, err := utils.ParseLauncherConfiguration(launcherConfigName)
	if err != nil {
		log.Error(err, "Failed to parse launcher configuration file")
		return err
	}
	if d.isEmptyStruct(config) {
		log.Error(err, "Failed to parse launcher configuration file. Configuration is empty")
		return err
	}
	//pass the yaml configuration to the Contaier struct
	if !d.isEmptyStruct(config.RunOnce.MediaProxyAgent) {

		mcmAgentContainer := general.Containers{}
		mcmAgentContainer.Type = general.MediaProxyAgent
		mcmAgentContainer.ContainerName = MediaProxyAgentContainerName
		mcmAgentContainer.Image = config.RunOnce.MediaProxyAgent.ImageAndTag
		mcmAgentContainer.Configuration.MediaProxyAgentConfig.ImageAndTag = config.RunOnce.MediaProxyAgent.ImageAndTag
		mcmAgentContainer.Configuration.MediaProxyAgentConfig.GRPCPort = config.RunOnce.MediaProxyAgent.GRPCPort
		mcmAgentContainer.Configuration.MediaProxyAgentConfig.RestPort = config.RunOnce.MediaProxyAgent.RestPort
		mcmAgentContainer.Configuration.MediaProxyAgentConfig.Network.Enable = config.RunOnce.MediaProxyAgent.Network.Enable
		mcmAgentContainer.Configuration.MediaProxyAgentConfig.Network.Name = config.RunOnce.MediaProxyAgent.Network.Name
		mcmAgentContainer.Configuration.MediaProxyAgentConfig.Network.IP = config.RunOnce.MediaProxyAgent.Network.IP

		err := utils.CreateAndRunContainer(ctx, d.cli, log, &mcmAgentContainer)
		if err != nil {
			log.Error(err, "Failed to create contianer MCM MediaProxy Agent!")
			return err
		}
	} else {
		log.Info("No information about MCM MediaProxy Agent provided. Omitting creation of MCM MediaProxy Agent container")
	}

	if !d.isEmptyStruct(config.RunOnce.MediaProxyMcm) {
		mediaProxyContainer := general.Containers{}
		mediaProxyContainer.Type = general.MediaProxyMCM
		mediaProxyContainer.ContainerName = MediaProxyContainerName
		mediaProxyContainer.Image = config.RunOnce.MediaProxyMcm.ImageAndTag
		mediaProxyContainer.Configuration.MediaProxyMcmConfig.ImageAndTag = config.RunOnce.MediaProxyMcm.ImageAndTag
		mediaProxyContainer.Configuration.MediaProxyMcmConfig.InterfaceName = config.RunOnce.MediaProxyMcm.InterfaceName
		mediaProxyContainer.Configuration.MediaProxyMcmConfig.Volumes = config.RunOnce.MediaProxyMcm.Volumes
		mediaProxyContainer.Configuration.MediaProxyMcmConfig.Network.Enable = config.RunOnce.MediaProxyMcm.Network.Enable
		mediaProxyContainer.Configuration.MediaProxyMcmConfig.Network.Name = config.RunOnce.MediaProxyMcm.Network.Name
		mediaProxyContainer.Configuration.MediaProxyMcmConfig.Network.IP = config.RunOnce.MediaProxyMcm.Network.IP
		
		err := utils.CreateAndRunContainer(ctx, d.cli, log, &mediaProxyContainer)
		if err != nil {
			log.Error(err, "Failed to create contianer MCM MediaProxy!")
			return err
		}
	} else {
		log.Info("No information about MCM MediaProxy provided. Omitting creation of MCM MediaProxy container")
	}

	if !d.isEmptyStruct(config.WorkloadToBeRun.NmosClient) {
		bcsNmosContainer := general.Containers{}
		bcsNmosContainer.Type = general.BcsPipelineNmosClient
		bcsNmosContainer.ContainerName = config.WorkloadToBeRun.NmosClient.Name
		bcsNmosContainer.Image = config.WorkloadToBeRun.NmosClient.ImageAndTag
		bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.ImageAndTag = config.WorkloadToBeRun.NmosClient.ImageAndTag
		bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.EnvironmentVariables = config.WorkloadToBeRun.NmosClient.EnvironmentVariables
		bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.NmosConfigPath = config.WorkloadToBeRun.NmosClient.NmosConfigPath
		bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.NmosConfigFileName = config.WorkloadToBeRun.NmosClient.NmosConfigFileName
		bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.Network.Enable = config.WorkloadToBeRun.NmosClient.Network.Enable
		bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.Network.Name = config.WorkloadToBeRun.NmosClient.Network.Name
		bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.Network.IP = config.WorkloadToBeRun.NmosClient.Network.IP
		bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.FfmpegConectionAddress = config.WorkloadToBeRun.FfmpegPipeline.Network.IP
		bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.FfmpegConnectionPort = strconv.Itoa(config.WorkloadToBeRun.FfmpegPipeline.GRPCPort)

		err = utils.CreateAndRunContainer(ctx, d.cli, log, &bcsNmosContainer)
		if err != nil {
			log.Error(err, "Failed to create contianer!")
			return err
		}
	} else {
		log.Info("No information about BCS NMOS client container provided. Omitting creation of BCS NMOS client container")
	}

	bcsPipelinesContainer := general.Containers{}
	bcsPipelinesContainer.Type = general.BcsPipelineFfmpeg
	bcsPipelinesContainer.ContainerName = config.WorkloadToBeRun.FfmpegPipeline.Name
	bcsPipelinesContainer.Image = config.WorkloadToBeRun.FfmpegPipeline.ImageAndTag
	bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Name = config.WorkloadToBeRun.FfmpegPipeline.Name
	bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.ImageAndTag = config.WorkloadToBeRun.FfmpegPipeline.ImageAndTag
	bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.GRPCPort = config.WorkloadToBeRun.FfmpegPipeline.GRPCPort
	bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.SourcePort = config.WorkloadToBeRun.FfmpegPipeline.SourcePort
	bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.EnvironmentVariables = config.WorkloadToBeRun.FfmpegPipeline.EnvironmentVariables
	bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Devnull = config.WorkloadToBeRun.FfmpegPipeline.Volumes.Devnull
	bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Dri = config.WorkloadToBeRun.FfmpegPipeline.Volumes.Dri
	bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Hugepages = config.WorkloadToBeRun.FfmpegPipeline.Volumes.Hugepages
	bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Imtl = config.WorkloadToBeRun.FfmpegPipeline.Volumes.Imtl
	bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Kahawai = config.WorkloadToBeRun.FfmpegPipeline.Volumes.Kahawai
	bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Shm = config.WorkloadToBeRun.FfmpegPipeline.Volumes.Shm
	bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.TmpHugepages = config.WorkloadToBeRun.FfmpegPipeline.Volumes.TmpHugepages
	bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Videos = config.WorkloadToBeRun.FfmpegPipeline.Volumes.Videos
	bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Devices.Dri = config.WorkloadToBeRun.FfmpegPipeline.Devices.Dri
	bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Devices.Vfio = config.WorkloadToBeRun.FfmpegPipeline.Devices.Vfio
	bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Network.Enable = config.WorkloadToBeRun.FfmpegPipeline.Network.Enable
	bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Network.Name = config.WorkloadToBeRun.FfmpegPipeline.Network.Name
	bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Network.IP = config.WorkloadToBeRun.FfmpegPipeline.Network.IP

	err = utils.CreateAndRunContainer(ctx, d.cli, log, &bcsPipelinesContainer)
	if err != nil {
		log.Error(err, "Failed to create contianer!")
		return err
	}
	return nil
}

func (d *DockerContainerController) IsContainerRunning(contaierName string) (bool, error) {
	cotainerStatus, err := d.cli.ContainerInspect(context.Background(), contaierName)
	if err != nil {
		return false, err
	}
	return cotainerStatus.State.Running, nil
}
