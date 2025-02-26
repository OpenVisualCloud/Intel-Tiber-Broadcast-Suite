/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

package utils

import (
	"context"
	"encoding/json"
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
	"gopkg.in/yaml.v2"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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
		fmt.Printf(">> MediaProxyAgentConfig: %+v\n", containerInfo.Configuration.MediaProxyAgentConfig)
		containerConfig = &container.Config{
			User: "root",
			Image: containerInfo.Configuration.MediaProxyAgentConfig.ImageAndTag,
			Cmd:   []string{"-c", containerInfo.Configuration.MediaProxyAgentConfig.RestPort, "-p", containerInfo.Configuration.MediaProxyAgentConfig.GRPCPort},
		}
	
		hostConfig = &container.HostConfig{
			Privileged: true,
			PortBindings: nat.PortMap{
				nat.Port(fmt.Sprintf("%s/tcp", containerInfo.Configuration.MediaProxyAgentConfig.RestPort)): []nat.PortBinding{{HostPort: containerInfo.Configuration.MediaProxyAgentConfig.RestPort}},
				nat.Port(fmt.Sprintf("%s/tcp", containerInfo.Configuration.MediaProxyAgentConfig.GRPCPort)): []nat.PortBinding{{HostPort: containerInfo.Configuration.MediaProxyAgentConfig.GRPCPort}},
			},
		}
	    if containerInfo.Configuration.MediaProxyAgentConfig.Network.Enable {
			networkConfig = &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					containerInfo.Configuration.MediaProxyAgentConfig.Network.Name: {
						IPAMConfig: &network.EndpointIPAMConfig{
							IPv4Address: containerInfo.Configuration.MediaProxyAgentConfig.Network.IP,
						},
					},
				},
			}
		}else{
			networkConfig = &network.NetworkingConfig{}
		}
	case general.MediaProxyMCM:
		fmt.Printf(">> MediaProxyMcmConfig: %+v\n", containerInfo.Configuration.MediaProxyMcmConfig)
		containerConfig = &container.Config{
			Image: containerInfo.Configuration.MediaProxyMcmConfig.ImageAndTag,
			Cmd:   []string{"-d", fmt.Sprintf("kernel:%s", containerInfo.Configuration.MediaProxyMcmConfig.InterfaceName),"-i", "localhost"},
		}
	
		hostConfig = &container.HostConfig{
			Privileged: true,
			Binds:      containerInfo.Configuration.MediaProxyMcmConfig.Volumes,
		}
	
	    if containerInfo.Configuration.MediaProxyMcmConfig.Network.Enable {
			networkConfig = &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					containerInfo.Configuration.MediaProxyMcmConfig.Network.Name: {
						IPAMConfig: &network.EndpointIPAMConfig{
							IPv4Address: containerInfo.Configuration.MediaProxyMcmConfig.Network.IP,
						},
					},
				},
			}
		}else{
			networkConfig = &network.NetworkingConfig{}
		}
    case general.BcsPipelineFfmpeg:
		fmt.Printf(">> BcsPipelineFfmpeg: %+v\n", containerInfo.Configuration.WorkloadConfig.FfmpegPipeline)

		containerConfig = &container.Config{
			User:       "root",
			Image: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.ImageAndTag,
			Cmd:   []string{containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Network.IP, fmt.Sprintf("%d", containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.GRPCPort)},
			Env: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.EnvironmentVariables,
			ExposedPorts: nat.PortSet{
				"20000/tcp": struct{}{},
				"20170/tcp": struct{}{},
			},
		}
	
		hostConfig = &container.HostConfig{
			Privileged: true,
			CapAdd:     []string{"ALL"},
			
			Mounts: []mount.Mount{
				{Type: mount.TypeBind, Source: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Videos, Target: "/videos"},
				{Type: mount.TypeBind, Source: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Dri, Target: "/usr/local/lib/x86_64-linux-gnu/dri/"},
				{Type: mount.TypeBind, Source: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Kahawai, Target: "/tmp/kahawai_lcore.lock"},
				{Type: mount.TypeBind, Source: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Devnull, Target: "/dev/null"},
				{Type: mount.TypeBind, Source: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.TmpHugepages, Target: "/tmp/hugepages"},
				{Type: mount.TypeBind, Source: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Hugepages, Target: "/hugepages"},
				{Type: mount.TypeBind, Source: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Imtl, Target: "/var/run/imtl"},
				{Type: mount.TypeBind, Source: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Shm, Target: "/dev/shm"},
			},
			IpcMode: "host",
		}
		hostConfig.Devices= []container.DeviceMapping{
			{PathOnHost: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Devices.Vfio, PathInContainer: "/dev/vfio"},
			{PathOnHost: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Dri, PathInContainer: "/dev/dri"},
		}
	
		networkConfig = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Network.Name: {
					IPAMConfig: &network.EndpointIPAMConfig{
						IPv4Address: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Network.IP,
					},
				},
			},
		}
	case general.BcsPipelineNmosClient:
		fmt.Printf(">> NmosClient: %+v\n", containerInfo.Configuration.WorkloadConfig.NmosClient)
		containerConfig = &container.Config{
			Image: containerInfo.Configuration.WorkloadConfig.NmosClient.ImageAndTag,
			Cmd: []string{"config/node.json"},
			Env: containerInfo.Configuration.WorkloadConfig.NmosClient.EnvironmentVariables,
			User:       "root",
		}
	
		hostConfig = &container.HostConfig{
			Privileged: true,
			Binds:      []string{fmt.Sprintf("%s:/home/config/", containerInfo.Configuration.WorkloadConfig.NmosClient.NmosConfigPath)},
		}
	
		networkConfig = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				containerInfo.Configuration.WorkloadConfig.NmosClient.Network.Name: {
					IPAMConfig: &network.EndpointIPAMConfig{
						IPv4Address: containerInfo.Configuration.WorkloadConfig.NmosClient.Network.IP,
					},
					Aliases: []string{containerInfo.Configuration.WorkloadConfig.NmosClient.Network.Name},
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

func boolPtr(b bool) *bool    { return &b }
func intstrPtr(i int) intstr.IntOrString {
    return intstr.IntOrString{IntVal: int32(i)}
}
		
func CreateMeshAgentDeployment(name string) *appsv1.Deployment {
	 return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mesh-agent-deployment",
			Namespace: "mcm",
			Labels: map[string]string{
				"app": "mesh-agent",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "mesh-agent",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "mesh-agent",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "mesh-agent",
							Image: "mcm/mesh-agent:bcs",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command: []string{
								"mesh-agent","-c", "8100", "-p", "50051",
							},
							Ports: []corev1.ContainerPort{
								{ContainerPort: 8100},
								{ContainerPort: 50051},
							},
							SecurityContext: &corev1.SecurityContext{
								Privileged: boolPtr(true),
							},
						},
					},
				},
			},
		},
	}
}

func CreateMeshAgentService(name string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mesh-agent-service",
			Namespace: "mcm",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "mesh-agent",
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "rest",
					Protocol:   corev1.ProtocolTCP,
					Port:       8100,
					TargetPort: intstrPtr(8100),
				},
				{
					Name: 	    "grpc",
					Protocol:   corev1.ProtocolTCP,
					Port:       50051,
					TargetPort: intstrPtr(50051),
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

func CreateBcsService(name string) *corev1.Service {
	return &corev1.Service{
        ObjectMeta: metav1.ObjectMeta{
          Name: "tiber-broadcast-suite",
		  Namespace: "default",
        },
        Spec: corev1.ServiceSpec{
          Type: corev1.ServiceTypeNodePort,
          Selector: map[string]string{
            "app": "tiber-broadcast-suite",
          },
          Ports: []corev1.ServicePort{
            {
              Protocol:   corev1.ProtocolTCP,
			  Name: "nmos-node-api",
              Port:       84,
              TargetPort: intstr.FromInt(84),
              NodePort:   30084,
            },
            {
              Protocol:   corev1.ProtocolTCP,
			  Name: "nmos-app-communication",
              Port:       5004,
              TargetPort: intstr.FromInt(5004),
              NodePort:   32054,
            },
          },
        },
      }
}

func CreatePersistentVolume(name string) *corev1.PersistentVolume {
	return &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mtl-pv",
		},
		Spec: corev1.PersistentVolumeSpec{
			Capacity: corev1.ResourceList{
				corev1.ResourceStorage: resource.MustParse("1Gi"),
			},
			VolumeMode: func() *corev1.PersistentVolumeMode { mode := corev1.PersistentVolumeFilesystem; return &mode }(),
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimRetain,
            StorageClassName:              "manual",
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/var/run/imtl",
				},
			},
		},
	}
}

func CreatePersistentVolumeClaim(name string) *corev1.PersistentVolumeClaim {
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mtl-pvc",
			Namespace: "mcm",
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
			StorageClassName: func() *string { s := "manual"; return &s }(),
			VolumeName:       "mtl-pv",
		},
	}
}

func CreateConfigMap(name string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tiber-broadcast-suite-config",
			Namespace: "default",
		},
		Data: map[string]string{
			"data": "bcsdata",
		},
	}
}

func ConvertYAMLToJSON(yamlData []byte) ([]byte, error) {
	var jsonData map[string]interface{}
	err := yaml.Unmarshal(yamlData, &jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return jsonBytes, nil
}


func CreateConfigMapFromFile(name, namespace, filePath string) (*corev1.ConfigMap, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tiber-broadcast-suite-config",
			Namespace: "default",
		},
		Data: map[string]string{
			"config": string(data),
		},
	}, nil
}

func CreateBcsDeployment(name string) *appsv1.Deployment {
	return &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
          Name:      "tiber-broadcast-suite",
          Namespace: "default",
        },
        Spec: appsv1.DeploymentSpec{
          Replicas: int32Ptr(1),
          Selector: &metav1.LabelSelector{
            MatchLabels: map[string]string{
              "app": "tiber-broadcast-suite",
            },
          },
          Template: corev1.PodTemplateSpec{
            ObjectMeta: metav1.ObjectMeta{
              Labels: map[string]string{
                "app": "tiber-broadcast-suite",
              },
            },
            Spec: corev1.PodSpec{
              Containers: []corev1.Container{
                {
                  Name:  "tiber-broadcast-suite-nmos-node",
                  Image: "tiber-broadcast-suite-nmos-node:latest",
				  ImagePullPolicy: corev1.PullIfNotPresent,
                  Args:  []string{"config/intel-node-tx.json"},
                  SecurityContext: &corev1.SecurityContext{
                    RunAsUser:  int64Ptr(0),
                    Privileged: boolPtr(true),
                    Capabilities: &corev1.Capabilities{
                      Add: []corev1.Capability{"ALL"},
                    },
                  },
                  VolumeMounts: []corev1.VolumeMount{
                    {
                      Name:      "config",
                      MountPath: "/home/config",
                    },
                  },
                  Env: []corev1.EnvVar{
                    {Name: "http_proxy", Value: ""},
                    {Name: "https_proxy", Value: ""},
                    {Name: "VFIO_PORT_TX", Value: "0000:ca:11.0"},
                  },
                  Ports: []corev1.ContainerPort{
                    {ContainerPort: 20000},
                    {ContainerPort: 20170},
                  },
                  Resources: corev1.ResourceRequirements{
                    Requests: corev1.ResourceList{
                      corev1.ResourceMemory: resource.MustParse("64Mi"),
                      corev1.ResourceCPU:    resource.MustParse("250m"),
                    },
                    Limits: corev1.ResourceList{
                      corev1.ResourceMemory: resource.MustParse("128Mi"),
                      corev1.ResourceCPU:    resource.MustParse("500m"),
                    },
                  },
                },
                {
                  Name:  "tiber-broadcast-suite",
                  Image: "tiber-broadcast-suite:latest",
				  ImagePullPolicy: corev1.PullIfNotPresent,
                  Args:  []string{"192.168.2.4", "50051"},
                  SecurityContext: &corev1.SecurityContext{
                    RunAsUser:  int64Ptr(0),
                    Privileged: boolPtr(true),
                    Capabilities: &corev1.Capabilities{
                      Add: []corev1.Capability{"ALL"},
                    },
                  },
                  VolumeMounts: []corev1.VolumeMount{
                    {Name: "videos", MountPath: "/videos"},
                    {Name: "dri", MountPath: "/usr/local/lib/x86_64-linux-gnu/dri"},
                    {Name: "kahawai-lock", MountPath: "/tmp/kahawai_lcore.lock"},
                    {Name: "dev-null", MountPath: "/dev/null"},
                    {Name: "hugepages-tmp", MountPath: "/tmp/hugepages"},
                    {Name: "hugepages", MountPath: "/hugepages"},
                    {Name: "imtl", MountPath: "/var/run/imtl"},
                    {Name: "shm", MountPath: "/dev/shm"},
                    {Name: "driDev", MountPath: "/dev/dri"},
                    {Name: "vfio", MountPath: "/dev/vfio"},
                  },
                  Env: []corev1.EnvVar{
                    {Name: "http_proxy", Value: ""},
                    {Name: "https_proxy", Value: ""},
                  },
                  Ports: []corev1.ContainerPort{
                    {ContainerPort: 20000},
                    {ContainerPort: 20170},
                  },
                  Resources: corev1.ResourceRequirements{
                    Requests: corev1.ResourceList{
                      corev1.ResourceMemory: resource.MustParse("64Mi"),
                      corev1.ResourceCPU:    resource.MustParse("250m"),
                    },
                    Limits: corev1.ResourceList{
                      corev1.ResourceMemory: resource.MustParse("128Mi"),
                      corev1.ResourceCPU:    resource.MustParse("500m"),
                    },
                  },
                },
              },
              Volumes: []corev1.Volume{
                {Name: "videos", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/root/DEMO_NMOS/move/nmos/nmos-cpp/Development/nmos-cpp-node/"}}},
                {Name: "dri", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/usr/lib/x86_64-linux-gnu/dri"}}},
                {Name: "kahawai-lock", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/tmp/kahawai_lcore.lock"}}},
                {Name: "dev-null", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/dev/null"}}},
                {Name: "hugepages-tmp", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/tmp/hugepages"}}},
                {Name: "hugepages", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/hugepages"}}},
                {Name: "imtl", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/var/run/imtl"}}},
                {Name: "shm", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/dev/shm"}}},
                {Name: "vfio", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/dev/vfio"}}},
                {Name: "driDev", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/dev/dri"}}},
				{Name: "config", VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: "tiber-broadcast-suite-config"}}}},
              },
            },
          },
        },
      }
}

func CreateDaemonSet(name string) *appsv1.DaemonSet {
	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "media-proxy",
			Namespace: "mcm",
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "media-proxy",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "media-proxy",
					},
				},
				Spec: corev1.PodSpec{
					// NodeSelector: map[string]string{
					// 	"node-role.kubernetes.io/worker": "true",
					// },
					Containers: []corev1.Container{
						{
							Name:    "media-proxy",
							Image:   "mcm/media-proxy:latest",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command: []string{"media_proxy"},
							Args:    []string{"-d", "kernel:eth0", "-i", "$(POD_IP)"},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:              resource.MustParse("2"),
									corev1.ResourceMemory:           resource.MustParse("8Gi"),
									corev1.ResourceHugePagesPrefix + "2Mi": resource.MustParse("1Gi"),
									corev1.ResourceHugePagesPrefix + "1Gi": resource.MustParse("2Gi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:              resource.MustParse("2"),
									corev1.ResourceMemory:           resource.MustParse("8Gi"),
									corev1.ResourceHugePagesPrefix + "2Mi": resource.MustParse("1Gi"),
									corev1.ResourceHugePagesPrefix + "1Gi": resource.MustParse("2Gi"),
								},
							},
							SecurityContext: &corev1.SecurityContext{
								Privileged: boolPtr(true),
								RunAsUser:  int64Ptr(0),
								RunAsGroup: int64Ptr(0),
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8001,
									HostPort:      8001,
									Protocol:      corev1.ProtocolTCP,
									Name:          "grpc-port",
								},
								{
									ContainerPort: 8002,
									HostPort:      8002,
									Protocol:      corev1.ProtocolTCP,
									Name:          "sdk-port",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "memif-dir",
									MountPath: "/run/mcm",
								},
								{
									Name:      "dev-vfio",
									MountPath: "/dev/vfio",
								},
								// {
								// 	Name:      "hugepage-2mi",
								// 	MountPath: "/hugepages-2Mi",
								// },
								// {
								// 	Name:      "hugepage-1gi",
								// 	MountPath: "/hugepages-1Gi",
								// },
								{
									Name:      "cache-volume",
									MountPath: "/dev/shm",
								},
								{
									Name:      "mtl-mgr",
									MountPath: "/var/run/imtl",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name: "NODE_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "spec.nodeName",
										},
									},
								},
								{
									Name: "NODE_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.hostIP",
										},
									},
								},
								{
									Name: "POD_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "memif-dir",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/tmp/mcm/memif",
								},
							},
						},
						{
							Name: "dev-vfio",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/dev/vfio",
								},
							},
						},
						{
							Name: "hugepage-2mi",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{
									Medium: corev1.StorageMediumHugePages,
								},
							},
						},
						{
							Name: "hugepage-1gi",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{
									Medium: corev1.StorageMediumHugePages,
								},
							},
						},
						{
							Name: "cache-volume",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{
									Medium:    corev1.StorageMediumMemory,
									SizeLimit: resource.NewQuantity(4*1024*1024*1024, resource.BinarySI),
								},
							},
						},
						{
							Name: "mtl-mgr",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "mtl-pvc",
								},
							},
						},
					},
				},
			},
		},
	}
}

func int64Ptr(i int64) *int64 {
	return &i
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
