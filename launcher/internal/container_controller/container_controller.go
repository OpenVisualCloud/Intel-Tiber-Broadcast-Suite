/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

// This package will be exploited in task SDBQ-1261
package containercontroller

import (
	"context"
	"reflect"

	"bcs.pod.launcher.intel/resources_library/resources/general"
	"bcs.pod.launcher.intel/resources_library/utils"
	"github.com/docker/docker/client"
	"github.com/go-logr/logr"
)
const (
	MediaProxyAgentContainerName = "mesh-agent"
	MediaProxyContainerName      = "media-proxy"
	BCSPipelineContainerName     = "bcs-ffmpeg-pipeline"
)
type ContainerController interface {
	CreateAndRunContainers(ctx context.Context, log logr.Logger) error
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
	if !d.isEmptyStruct(config.RunOnce.MediaProxyAgent) {
		mcmAgentContainer := general.Containers{
			ContainerName:   MediaProxyAgentContainerName, //the name is predefined because only one instance is required on the machine
			Ip:              config.RunOnce.MediaProxyAgent.IP,
			ExposedPort:     config.RunOnce.MediaProxyAgent.ExposedPort, //"80/tcp",
			BindingHostPort: config.RunOnce.MediaProxyAgent.BindingHostPort,
			Image:           config.RunOnce.MediaProxyAgent.ImageAndTag,
			Privileged:      config.RunOnce.MediaProxyAgent.Privileged,
		}
		err := utils.CreateAndRunContainer(ctx, d.cli, log, mcmAgentContainer)
		if err != nil {
			log.Error(err, "Failed to create contianer MCM MediaProxy Agent!")
			return err
		}
	} else {
		log.Info("No information about MCM MediaProxy Agent provided. Omitting creation of MCM MediaProxy Agent container")
	}

	if !d.isEmptyStruct(config.RunOnce.MediaProxyMcm) {
		mediaProxyContainer := general.Containers{
			ContainerName:   MediaProxyContainerName, //the name is predefined because only one instance is required on the machine
			Ip:              config.RunOnce.MediaProxyMcm.IP,
			ExposedPort:     config.RunOnce.MediaProxyMcm.ExposedPort, //"80/tcp",
			BindingHostPort: config.RunOnce.MediaProxyMcm.BindingHostPort,
			Image:           config.RunOnce.MediaProxyMcm.ImageAndTag,
			VolumeMount:     config.RunOnce.MediaProxyMcm.Volumes,
			Privileged:      config.RunOnce.MediaProxyMcm.Privileged,
		}
		err := utils.CreateAndRunContainer(ctx, d.cli, log, mediaProxyContainer)
		if err != nil {
			log.Error(err, "Failed to create contianer MCM MediaProxy!")
			return err
		}
	} else {
		log.Info("No information about MCM MediaProxy provided. Omitting creation of MCM MediaProxy container")
	}

	if !d.isEmptyStruct(config.WorkloadToBeRun.NmosClient) {
		bcsNmosContainer := general.Containers{
			ContainerName:   config.WorkloadToBeRun.NmosClient.Name,
			Ip:              config.WorkloadToBeRun.NmosClient.IP,
			ExposedPort:     config.WorkloadToBeRun.NmosClient.ExposedPort,
			BindingHostPort: config.WorkloadToBeRun.NmosClient.BindingHostPort,
			Image:           config.WorkloadToBeRun.NmosClient.ImageAndTag,
			VolumeMount:     config.WorkloadToBeRun.NmosClient.Volumes,
			EnviromentVariables: config.WorkloadToBeRun.NmosClient.EnvironmentVariables,
			NetworkMode:    config.WorkloadToBeRun.NmosClient.Network,
		}
		err = utils.CreateAndRunContainer(ctx, d.cli, log, bcsNmosContainer)
		if err != nil {
			log.Error(err, "Failed to create contianer!")
			return err
		}
	} else {
		log.Info("No information about BCS NMOS client container provided. Omitting creation of BCS NMOS client container")

	}

	bcsPipelinesContainer := general.Containers{
		ContainerName:   config.WorkloadToBeRun.FFmpegPipeline.Name,
		Ip:              config.WorkloadToBeRun.FFmpegPipeline.IP,
		ExposedPort:     config.WorkloadToBeRun.FFmpegPipeline.ExposedPort,
		BindingHostPort: config.WorkloadToBeRun.FFmpegPipeline.BindingHostPort,
		Image:           config.WorkloadToBeRun.FFmpegPipeline.ImageAndTag,
	}

	err = utils.CreateAndRunContainer(ctx, d.cli, log, bcsPipelinesContainer)
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
