/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

// This package will be exploited in task SDBQ-1261
package containercontroller

import (
	"context"

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

// Hardcoded poc (it will be generalized in next PR)
// use case covers running coaiers on single host
func (d *DockerContainerController) CreateAndRunContainers(ctx context.Context, log logr.Logger) error {
	mcmAgentContainer := general.Containers{
		ContainerName:   "mcm-agent",
		Ip:              "0.0.0.0",
		ExposedPort:     "80/tcp",
		BindingHostPort: "8089",
		Image:           "nginx:latest",
	}
	mediaProxyContainer := general.Containers{
		ContainerName:   "media-proxy",
		Ip:              "0.0.0.0",
		ExposedPort:     "80/tcp",
		BindingHostPort: "8085",
		Image:           "nginx:latest",
	}
	bcsNmosContainer := general.Containers{
		ContainerName:   "bcsNmos",
		Ip:              "0.0.0.0",
		ExposedPort:     "80/tcp",
		BindingHostPort: "8082",
		Image:           "nginx:latest",
	}
	bcsPipelinesContainer := general.Containers{
		ContainerName:   "bcsPipeline",
		Ip:              "0.0.0.0",
		ExposedPort:     "80/tcp",
		BindingHostPort: "8083",
		Image:           "nginx:latest",
	}

	err := utils.CreateAndRunContainer(ctx, d.cli, log, mcmAgentContainer)
	if err != nil {
		log.Error(err, "Failed to create contianer!")
		return err
	}
	err = utils.CreateAndRunContainer(ctx, d.cli, log, mediaProxyContainer)
	if err != nil {
		log.Error(err, "Failed to create contianer!")
		return err
	}
	err = utils.CreateAndRunContainer(ctx, d.cli, log, bcsNmosContainer)
	if err != nil {
		log.Error(err, "Failed to create contianer!")
		return err
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
