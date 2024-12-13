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

// use case covers running containers on single host
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
			ContainerName:   "mcm-agent", //the name is predefined because only one instance is required on the machine
			Ip:              config.RunOnce.MediaProxyAgent.IP,
			ExposedPort:     config.RunOnce.MediaProxyAgent.ExposedPort, //"80/tcp",
			BindingHostPort: config.RunOnce.MediaProxyAgent.BindingHostPort,
			Image:           config.RunOnce.MediaProxyAgent.ImageAndTag,
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
			ContainerName:   "media-proxy", //the name is predefined because only one instance is required on the machine
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
			BindingHostPort: "8082",
			Image:           "nginx:latest",
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
		ContainerName:   "bcs-ffmpeg-pipeline",
		Ip:              "0.0.0.0",
		ExposedPort:     "80/tcp",
		BindingHostPort: "8083",
		Image:           "nginx:latest",
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
