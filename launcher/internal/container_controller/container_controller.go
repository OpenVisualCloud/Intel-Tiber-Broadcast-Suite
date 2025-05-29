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

	"bcs.pod.launcher.intel/resources_library/parser"
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
	config, err := parser.ParseLauncherConfiguration(launcherConfigName)
	if err != nil {
		log.Error(err, "Failed to parse launcher configuration file")
		return err
	}
	if d.isEmptyStruct(config) {
		log.Error(err, "Failed to parse launcher configuration file. Configuration is empty")
		return err
	}
	//pass the yaml configuration to the Container struct
	if !d.isEmptyStruct(config.RunOnce.MediaProxyAgent) {

		mcmAgentContainer := general.Containers{}
		mcmAgentContainer.Type = general.MediaProxyAgent
		mcmAgentContainer.ContainerName = MediaProxyAgentContainerName
		mcmAgentContainer.Image = config.RunOnce.MediaProxyAgent.ImageAndTag
		if !config.RunOnce.MediaProxyAgent.Network.Enable {
			config.RunOnce.MediaProxyAgent.Network.Name = "host"
			// Note - When you use the host network mode in Docker, the container shares
			// the host machine's network stack.
		}

		err := utils.CreateAndRunContainer(ctx, d.cli, log, &mcmAgentContainer, &config)
		if err != nil {
			log.Error(err, "Failed to create container MCM MediaProxy Agent!")
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
		if !config.RunOnce.MediaProxyMcm.Network.Enable {
			config.RunOnce.MediaProxyMcm.Network.Name = "host"
			// Note - When you use the host network mode in Docker, the container shares
			// the host machine's network stack.
		}
		err := utils.CreateAndRunContainer(ctx, d.cli, log, &mediaProxyContainer, &config)
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

	for n, instance := range config.WorkloadToBeRun {
		if d.isEmptyStruct(instance.FfmpegPipeline) || d.isEmptyStruct(instance.NmosClient) {
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
		err = utils.CreateAndRunContainer(ctx, d.cli, log, &bcsPipelinesContainer, &config)
		if err != nil {
			log.Error(err, "Failed to create container for FFMPEG pipeline instance %d!")
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

		err = utils.CreateAndRunContainer(ctx, d.cli, log, &bcsNmosContainer, &config)
		if err != nil {
			log.Error(err, "Failed to create container!")
			return err
		}

	}
	return nil
}

func (d *DockerContainerController) IsContainerRunning(containerName string) (bool, error) {
	containerStatus, err := d.cli.ContainerInspect(context.Background(), containerName)
	if err != nil {
		return false, err
	}
	return containerStatus.State.Running, nil
}
