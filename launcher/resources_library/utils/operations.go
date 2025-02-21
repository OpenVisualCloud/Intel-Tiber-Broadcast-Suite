/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"bcs.pod.launcher.intel/resources_library/resources/general"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/go-logr/logr"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func isImagePulled(ctx context.Context, cli *client.Client, imageName string) (error, bool) {
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

func pullImageIfNotExists(ctx context.Context, cli *client.Client, imageName string, log logr.Logger) error {

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

func doesContainerExist(ctx context.Context, cli *client.Client, containerName string) (error, bool) {
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

func isContainerRunning(ctx context.Context, cli *client.Client, containerName string) (error, bool) {
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

func removeContainer(ctx context.Context, cli *client.Client, containerID string) error {
	return cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
}


func constructContainerConfig(containerInfo general.Containers) (*container.Config, *container.HostConfig, *network.NetworkingConfig) {
	var containerConfig *container.Config
	var hostConfig *container.HostConfig
	var networkConfig *network.NetworkingConfig

	switch containerInfo.Type {
	case general.MediaProxyAgent:
		containerConfig = &container.Config{
			User: "root",
			Image: "mcm/mesh-agent:latest",
			Cmd:   []string{"-c", "8100", "-p", "50051"},
		}
	
		hostConfig = &container.HostConfig{
			Privileged: true,
			PortBindings: nat.PortMap{
				"8100/tcp": []nat.PortBinding{{HostPort: "8100"}},
				"50051/tcp": []nat.PortBinding{{HostPort: "50051"}},
			},
		}
		isNetwork := true
	    if isNetwork {
			networkConfig = &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					"my_net_801f0": {
						IPAMConfig: &network.EndpointIPAMConfig{
							IPv4Address: "192.168.2.11",
						},
					},
				},
			}
		}else{
			networkConfig = &network.NetworkingConfig{}
		}
	case general.MediaProxyMCM:
		containerConfig = &container.Config{
			Image: "mcm/media-proxy:latest",
			Cmd:   []string{"-d", "kernel:eth0", "-i", "localhost"},
		}
	
		hostConfig = &container.HostConfig{
			Privileged: true,
			Binds:      []string{"/dev/vfio:/dev/vfio"},
		}
	
		isNetwork := true
	    if isNetwork {
			networkConfig = &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					"my_net_801f0": {
						IPAMConfig: &network.EndpointIPAMConfig{
							IPv4Address: "192.168.2.14",
						},
					},
				},
			}
		}else{
			networkConfig = &network.NetworkingConfig{}
		}
    case general.BcsPipelineFfmpeg:
		containerConfig = &container.Config{
			User:       "root",
			Image: "tiber-broadcast-suite:latest",
			Cmd:   []string{"192.168.2.4", "50051"},
			Env: []string{
				"http_proxy=",
				"https_proxy=",
			},
			ExposedPorts: nat.PortSet{
				"20000/tcp": struct{}{},
				"20170/tcp": struct{}{},
			},
		}
	
		hostConfig = &container.HostConfig{
			Privileged: true,
			CapAdd:     []string{"ALL"},
			
			Mounts: []mount.Mount{
				{Type: mount.TypeBind, Source: "/root", Target: "/videos"}, //customizable
				{Type: mount.TypeBind, Source: "/usr/lib/x86_64-linux-gnu/dri", Target: "/usr/local/lib/x86_64-linux-gnu/dri/"},
				{Type: mount.TypeBind, Source: "/tmp/kahawai_lcore.lock", Target: "/tmp/kahawai_lcore.lock"},
				{Type: mount.TypeBind, Source: "/dev/null", Target: "/dev/null"},
				{Type: mount.TypeBind, Source: "/tmp/hugepages", Target: "/tmp/hugepages"},
				{Type: mount.TypeBind, Source: "/hugepages", Target: "/hugepages"},
				{Type: mount.TypeBind, Source: "/var/run/imtl", Target: "/var/run/imtl"},
				{Type: mount.TypeBind, Source: "/dev/shm", Target: "/dev/shm"},
			},
			IpcMode: "host",
		}
		hostConfig.Devices= []container.DeviceMapping{
			{PathOnHost: "/dev/vfio", PathInContainer: "/dev/vfio"},
			{PathOnHost: "/dev/dri", PathInContainer: "/dev/dri"},
		}
	
		networkConfig = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"my_net_801f0": {
					IPAMConfig: &network.EndpointIPAMConfig{
						IPv4Address: "192.168.2.4",
					},
				},
			},
		}
	case general.BcsPipelineNmosClient:
		containerConfig = &container.Config{
			Image: "tiber-broadcast-suite-nmos-node:latest",
			Cmd:   []string{"config/node.json"},
			Env: []string{
				"http_proxy=",
				"https_proxy=",
				"VFIO_PORT_TX=0000:ca:11.0",
			},
			User:       "root",
			// WorkingDir: "/home/config/",
		}
	
		hostConfig = &container.HostConfig{
			Privileged: true,
			// CapAdd:     []string{"ALL"},
			Binds:      []string{fmt.Sprintf("%s:/home/config/", "/root/DEMO_NMOS/move/nmos/build-nmos-cpp/")},
			// IpcMode:    "host",
		}
	
		networkConfig = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"my_net_801f0": {
					IPAMConfig: &network.EndpointIPAMConfig{
						IPv4Address: "192.168.2.2",
					},
					Aliases: []string{"my_net_801f0"},
				},
			},
		}
	default:
		containerConfig, hostConfig, networkConfig = nil, nil, nil
	}

	return containerConfig, hostConfig, networkConfig
}

func CreateAndRunContainer(ctx context.Context, cli *client.Client, log logr.Logger, containerInfo general.Containers) error {
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

	containerConfig, hostConfig, networkConfig := constructContainerConfig(containerInfo)
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

// func CreateDeployment(myApp *bcsv1.BcsConfig, name string) appsv1.Deployment {
func CreateDeployment(name string) *appsv1.Deployment {

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
						},
					},
				},
			},
		},
	}
}

func CreateService(name string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": name},
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					Port:     80,
				},
			},
		},
	}
}

func CreatePersistentVolume(name string) *corev1.PersistentVolume {
	return &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: corev1.PersistentVolumeSpec{
			Capacity: corev1.ResourceList{
				corev1.ResourceStorage: resource.MustParse("1Gi"),
			},
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/mnt/data",
				},
			},
		},
	}
}

func CreatePersistentVolumeClaim(name string) *corev1.PersistentVolumeClaim {
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
		},
	}
}

func CreateConfigMap(name string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Data: map[string]string{
			"data": "bcsdata",
		},
	}
}

func CreateDaemonSet(name string) *appsv1.DaemonSet {
	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "example-container",
							Image: "nginx:latest",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "example-volume",
									MountPath: "/usr/share/nginx/html",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "example-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: name,
								},
							},
						},
					},
				},
			},
		},
	}
}

func int32Ptr(i int32) *int32 { return &i }

func CreateNamespace(namespaceName string) *corev1.Namespace {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceName,
		},
	}
	return namespace
}
