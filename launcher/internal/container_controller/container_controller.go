/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

package containercontroller

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"bcs.pod.launcher.intel/resources_library/resources/general"
	"bcs.pod.launcher.intel/resources_library/utils"

	// "bcs.pod.launcher.intel/resources_library/utils"
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
}

type DockerContainerController struct {
	cli *client.Client
}

func NewDockerContainerController() (*DockerContainerController, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerContainerController{cli: cli}, nil
}

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

func (d *DockerContainerController) CreateAndRunContainers(ctx context.Context, launcherConfigName string, log logr.Logger) error {
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
		if config.RunOnce.MediaProxyAgent.Network.Enable {
			mcmAgentContainer.Configuration.MediaProxyAgentConfig.Network.Name = config.RunOnce.MediaProxyAgent.Network.Name
			mcmAgentContainer.Configuration.MediaProxyAgentConfig.Network.IP = config.RunOnce.MediaProxyAgent.Network.IP
		} else {
			mcmAgentContainer.Configuration.MediaProxyAgentConfig.Network.Name = "host"
			// Note - When you use the host network mode in Docker, the container shares
			// the host machine's network stack.
		}

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
		if config.RunOnce.MediaProxyMcm.Network.Enable {
			mediaProxyContainer.Configuration.MediaProxyMcmConfig.Network.Name = config.RunOnce.MediaProxyMcm.Network.Name
			mediaProxyContainer.Configuration.MediaProxyMcmConfig.Network.IP = config.RunOnce.MediaProxyMcm.Network.IP
		} else {
			mediaProxyContainer.Configuration.MediaProxyMcmConfig.Network.Name = "host"
			// Note - When you use the host network mode in Docker, the container shares
			// the host machine's network stack.
		}
		err := utils.CreateAndRunContainer(ctx, d.cli, log, &mediaProxyContainer)
		if err != nil {
			log.Error(err, "Failed to create contianer MCM MediaProxy!")
			return err
		}
	} else {
		log.Info("No information about MCM MediaProxy provided. Omitting creation of MCM MediaProxy container")
	}

	if d.isEmptyStruct(config.WorkloadToBeRun) || len(config.WorkloadToBeRun) == 0 {
		log.Info("No workloads provided under workloadToBeRun. Omitting creation of BCS pipeline and NMOS node containers")
	}

	for _, instance := range config.WorkloadToBeRun {
		if d.isEmptyStruct(instance.FfmpegPipeline) || d.isEmptyStruct(instance.NmosClient) {
			return fmt.Errorf("no information about BCS pipeline provided. Either FfmpegPipeline or NmosClient is empty for instance Ffmpeg: %s; Nmos: %s", instance.FfmpegPipeline.Name, instance.NmosClient.Name)
		}
		bcsPipelinesContainer := general.Containers{}
		bcsPipelinesContainer.Type = general.BcsPipelineFfmpeg
		bcsPipelinesContainer.ContainerName = instance.FfmpegPipeline.Name
		bcsPipelinesContainer.Image = instance.FfmpegPipeline.ImageAndTag
		bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Name = instance.FfmpegPipeline.Name
		bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.ImageAndTag = instance.FfmpegPipeline.ImageAndTag
		bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.GRPCPort = instance.FfmpegPipeline.GRPCPort
		bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.EnvironmentVariables = instance.FfmpegPipeline.EnvironmentVariables
		bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Devnull = instance.FfmpegPipeline.Volumes.Devnull
		bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Dri = instance.FfmpegPipeline.Volumes.Dri
		bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Hugepages = instance.FfmpegPipeline.Volumes.Hugepages
		bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Imtl = instance.FfmpegPipeline.Volumes.Imtl
		bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Kahawai = instance.FfmpegPipeline.Volumes.Kahawai
		bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Shm = instance.FfmpegPipeline.Volumes.Shm
		bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.TmpHugepages = instance.FfmpegPipeline.Volumes.TmpHugepages
		bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Videos = instance.FfmpegPipeline.Volumes.Videos
		bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Devices.Dri = instance.FfmpegPipeline.Devices.Dri
		bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Devices.Vfio = instance.FfmpegPipeline.Devices.Vfio
		bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Network.Enable = instance.FfmpegPipeline.Network.Enable
		if instance.FfmpegPipeline.Network.Enable {
			bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Network.Name = instance.FfmpegPipeline.Network.Name
			bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Network.IP = instance.FfmpegPipeline.Network.IP
		} else {
			bcsPipelinesContainer.Configuration.WorkloadConfig.FfmpegPipeline.Network.Name = "host"
			// Note - When you use the host network mode in Docker, the container shares
			// the host machine's network stack.
		}
		err = utils.CreateAndRunContainer(ctx, d.cli, log, &bcsPipelinesContainer)
		if err != nil {
			log.Error(err, "Failed to create contianer!")
			return err
		}
		bcsNmosContainer := general.Containers{}
		bcsNmosContainer.Type = general.BcsPipelineNmosClient
		bcsNmosContainer.ContainerName = instance.NmosClient.Name
		bcsNmosContainer.Image = instance.NmosClient.ImageAndTag
		bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.ImageAndTag = instance.NmosClient.ImageAndTag
		bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.EnvironmentVariables = instance.NmosClient.EnvironmentVariables
		bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.NmosConfigPath = instance.NmosClient.NmosConfigPath
		bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.NmosConfigFileName = instance.NmosClient.NmosConfigFileName
		bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.Network.Enable = instance.NmosClient.Network.Enable
		bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.NmosPort = instance.NmosClient.NmosPort

		if instance.NmosClient.Network.Enable {
			bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.Network.Name = instance.NmosClient.Network.Name
			bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.Network.IP = instance.NmosClient.Network.IP
		} else {
			bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.Network.Name = "host"
			bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.Network.IP = instance.NmosClient.Network.IP
			// Note - When you use the host network mode in Docker, the container shares
			// the host machine's network stack.
		}
		bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.FfmpegConnectionAddress = instance.FfmpegPipeline.Network.IP
		bcsNmosContainer.Configuration.WorkloadConfig.NmosClient.FfmpegConnectionPort = strconv.Itoa(instance.FfmpegPipeline.GRPCPort)

		err = utils.CreateAndRunContainer(ctx, d.cli, log, &bcsNmosContainer)
		if err != nil {
			log.Error(err, "Failed to create contianer!")
			return err
		}

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
