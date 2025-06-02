/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"bcs.pod.launcher.intel/resources_library/parser"
	"bcs.pod.launcher.intel/resources_library/resources/bcs"
	"bcs.pod.launcher.intel/resources_library/resources/general"
	"bcs.pod.launcher.intel/resources_library/resources/nmos"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/go-logr/logr"
	"gopkg.in/yaml.v2"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	bcsv1 "bcs.pod.launcher.intel/api/v1"
)

func updateNmosJsonFile(filePath string, ip string, port string) error {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	var config nmos.Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return err
	}

	config.FfmpegGrpcServerAddress = ip
	config.FfmpegGrpcServerPort = port

	updatedJson, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}

	err = os.WriteFile(filePath, updatedJson, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	fmt.Println("Updated configuration saved to ", filePath)
	return nil
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func ConstructContainerConfig(containerInfo *general.Containers, config *parser.Configuration, log logr.Logger) (*container.Config, *container.HostConfig, *network.NetworkingConfig) {
	var containerConfig *container.Config
	var hostConfig *container.HostConfig
	var networkConfig *network.NetworkingConfig

	switch containerInfo.Type {
	case general.MediaProxyAgent:
		fmt.Printf(">> MediaProxyAgentConfig: %+v\n", config.RunOnce.MediaProxyAgent)
		containerConfig = &container.Config{
			User:  "root",
			Image: config.RunOnce.MediaProxyAgent.ImageAndTag,
			Cmd:   []string{"-c", config.RunOnce.MediaProxyAgent.RestPort, "-p", config.RunOnce.MediaProxyAgent.GRPCPort},
			ExposedPorts: nat.PortSet{
				nat.Port(fmt.Sprintf("%s/tcp", config.RunOnce.MediaProxyAgent.RestPort)): struct{}{},
				nat.Port(fmt.Sprintf("%s/tcp", config.RunOnce.MediaProxyAgent.GRPCPort)): struct{}{},
			},
		}

		hostConfig = &container.HostConfig{
			Privileged: true,
			PortBindings: nat.PortMap{
				nat.Port(fmt.Sprintf("%s/tcp", config.RunOnce.MediaProxyAgent.RestPort)): []nat.PortBinding{{HostPort: config.RunOnce.MediaProxyAgent.RestPort}},
				nat.Port(fmt.Sprintf("%s/tcp", config.RunOnce.MediaProxyAgent.GRPCPort)): []nat.PortBinding{{HostPort: config.RunOnce.MediaProxyAgent.GRPCPort}},
			},
		}
		if config.RunOnce.MediaProxyAgent.Network.Enable {
			hostConfig.NetworkMode = container.NetworkMode(config.RunOnce.MediaProxyAgent.Network.Name)
			networkConfig = &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					config.RunOnce.MediaProxyAgent.Network.Name: {
						IPAMConfig: &network.EndpointIPAMConfig{
							IPv4Address: config.RunOnce.MediaProxyAgent.Network.IP,
						},
					},
				},
			}
		} else {
			networkConfig = &network.NetworkingConfig{}
			hostConfig.NetworkMode = "host"
		}
	case general.MediaProxyMCM:
		fmt.Printf(">> MediaProxyMcmConfig: %+v\n", config.RunOnce.MediaProxyMcm)
		containerConfig = &container.Config{
			Image: config.RunOnce.MediaProxyMcm.ImageAndTag,
			Cmd:   []string{"-d", fmt.Sprintf("kernel:%s", config.RunOnce.MediaProxyMcm.InterfaceName), "-i", "localhost"},
		}

		hostConfig = &container.HostConfig{
			Privileged: true,
			Binds:      config.RunOnce.MediaProxyMcm.Volumes,
		}

		if config.RunOnce.MediaProxyMcm.Network.Enable {
			hostConfig.NetworkMode = container.NetworkMode(config.RunOnce.MediaProxyMcm.Network.Name)

			networkConfig = &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					config.RunOnce.MediaProxyMcm.Network.Name: {
						IPAMConfig: &network.EndpointIPAMConfig{
							IPv4Address: config.RunOnce.MediaProxyMcm.Network.IP,
						},
					},
				},
			}
		} else {
			networkConfig = &network.NetworkingConfig{}
			hostConfig.NetworkMode = "host"
		}
	case general.BcsPipelineFfmpeg:
		fmt.Printf(">> BcsPipelineFfmpeg: %+v\n", config.WorkloadToBeRun)

		containerConfig = &container.Config{
			User:  "root",
			Image: config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.ImageAndTag,
			Env:   config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.EnvironmentVariables,
			ExposedPorts: nat.PortSet{
				nat.Port(fmt.Sprintf("%d/tcp", config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.GRPCPort)): struct{}{},
			},
		}

		if config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.Network.Enable {
			containerConfig.Cmd = []string{config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.Network.IP, fmt.Sprintf("%d", config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.GRPCPort)}
		} else {
			containerConfig.Cmd = []string{"localhost", fmt.Sprintf("%d", config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.GRPCPort)}
		}

		hostConfig = &container.HostConfig{
			Privileged: true,
			CapAdd:     []string{"ALL"},
			PortBindings: nat.PortMap{
				nat.Port(fmt.Sprintf("%d/tcp", config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.GRPCPort)): []nat.PortBinding{
					{
						HostPort: fmt.Sprintf("%d", config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.GRPCPort),
					},
				},
			},
			Mounts: []mount.Mount{
				{Type: mount.TypeBind, Source: config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.Volumes.Videos, Target: "/videos"},
				{Type: mount.TypeBind, Source: config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.Volumes.Dri, Target: "/usr/local/lib/x86_64-linux-gnu/dri"},
				{Type: mount.TypeBind, Source: config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.Volumes.Kahawai, Target: "/tmp/kahawai_lcore.lock"},
				{Type: mount.TypeBind, Source: config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.Volumes.Devnull, Target: "/dev/null"},
				{Type: mount.TypeBind, Source: config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.Volumes.TmpHugepages, Target: "/tmp/hugepages"},
				{Type: mount.TypeBind, Source: config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.Volumes.Hugepages, Target: "/hugepages"},
				{Type: mount.TypeBind, Source: config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.Volumes.Imtl, Target: "/var/run/imtl"},
				{Type: mount.TypeBind, Source: config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.Volumes.Shm, Target: "/dev/shm"},
			},
			IpcMode: "host",
		}
		hostConfig.Devices = []container.DeviceMapping{
			{PathOnHost: config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.Devices.Vfio, PathInContainer: "/dev/vfio"},
			{PathOnHost: config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.Devices.Dri, PathInContainer: "/dev/dri"},
		}
		if config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.Network.Enable {
			hostConfig.NetworkMode = container.NetworkMode(config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.Network.Name)
			networkConfig = &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.Network.Name: {
						IPAMConfig: &network.EndpointIPAMConfig{
							IPv4Address: config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.Network.IP,
						},
						Aliases: []string{config.WorkloadToBeRun[containerInfo.Id].FfmpegPipeline.Network.Name},
					},
				},
			}
		} else {
			networkConfig = &network.NetworkingConfig{}
			hostConfig.NetworkMode = "host"
		}
	case general.BcsPipelineNmosClient:
		nmosFileNameJson := config.WorkloadToBeRun[containerInfo.Id].NmosClient.NmosConfigFileName
		nmosFilePathJson := config.WorkloadToBeRun[containerInfo.Id].NmosClient.NmosConfigPath + "/" + nmosFileNameJson
		if !FileExists(nmosFilePathJson) {
			log.Error(errors.New("NMOS json file does not exist"), "NMOS json file does not exist")
			return nil, nil, nil
		}
		errUpdateJson := updateNmosJsonFile(nmosFilePathJson,
			config.WorkloadToBeRun[containerInfo.Id].NmosClient.FfmpegConnectionAddress,
			config.WorkloadToBeRun[containerInfo.Id].NmosClient.FfmpegConnectionPort)
		if errUpdateJson != nil {
			log.Error(errUpdateJson, "Error updating NMOS json file")
			return nil, nil, nil
		}
		configPathContainer := "config/" + nmosFileNameJson
		containerConfig = &container.Config{
			Image: config.WorkloadToBeRun[containerInfo.Id].NmosClient.ImageAndTag,
			Cmd:   []string{configPathContainer},
			Env:   config.WorkloadToBeRun[containerInfo.Id].NmosClient.EnvironmentVariables,
			User:  "root",
			ExposedPorts: nat.PortSet{
				nat.Port(fmt.Sprintf("%d/tcp", config.WorkloadToBeRun[containerInfo.Id].NmosClient.NmosPort)): struct{}{},
			},
		}

		hostConfig = &container.HostConfig{
			Privileged: true,
			PortBindings: nat.PortMap{
				nat.Port(fmt.Sprintf("%d/tcp", config.WorkloadToBeRun[containerInfo.Id].NmosClient.NmosPort)): []nat.PortBinding{
					{
						HostPort: fmt.Sprintf("%d", config.WorkloadToBeRun[containerInfo.Id].NmosClient.NmosPort),
					},
				},
			},
			Binds: []string{fmt.Sprintf("%s:/home/config/", config.WorkloadToBeRun[containerInfo.Id].NmosClient.NmosConfigPath)},
		}

		if config.WorkloadToBeRun[containerInfo.Id].NmosClient.Network.Enable {
			hostConfig.NetworkMode = container.NetworkMode(config.WorkloadToBeRun[containerInfo.Id].NmosClient.Network.Name)
			networkConfig = &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					config.WorkloadToBeRun[containerInfo.Id].NmosClient.Network.Name: {
						IPAMConfig: &network.EndpointIPAMConfig{
							IPv4Address: config.WorkloadToBeRun[containerInfo.Id].NmosClient.Network.IP,
						},
						Aliases: []string{config.WorkloadToBeRun[containerInfo.Id].NmosClient.Network.Name},
					},
				},
			}
		} else {
			networkConfig = &network.NetworkingConfig{}
			hostConfig.NetworkMode = "host"
		}
	default:
		containerConfig, hostConfig, networkConfig = nil, nil, nil
	}

	return containerConfig, hostConfig, networkConfig
}

func boolPtr(b bool) *bool { return &b }

type K8sConfig struct {
	K8s        bool `yaml:"k8s"`
	Definition struct {
		MeshAgent struct {
			Image               string          `yaml:"image"`
			RestPort            int             `yaml:"restPort"`
			GrpcPort            int             `yaml:"grpcPort"`
			Resources           bcs.HwResources `yaml:"resources"`
			ScheduleOnNode      []string        `yaml:"scheduleOnNode,omitempty"`
			DoNotScheduleOnNode []string        `yaml:"doNotScheduleOnNode,omitempty"`
		} `yaml:"meshAgent"`
		MediaProxy struct {
			Image     string          `yaml:"image"`
			Command   []string        `yaml:"command"`
			Args      []string        `yaml:"args"`
			GrpcPort  int             `yaml:"grpcPort"`
			SdkPort   int             `yaml:"sdkPort"`
			Resources bcs.HwResources `yaml:"resources"`
			Volumes   struct {
				Memif     string `yaml:"memif"`
				Vfio      string `yaml:"vfio"`
				CacheSize string `yaml:"cache-size"`
			} `yaml:"volumes"`
			PvHostPath          string   `yaml:"pvHostPath"`
			PvStorageClass      string   `yaml:"pvStorageClass"`
			PvStorage           string   `yaml:"pvStorage"`
			PvcStorage          string   `yaml:"pvcStorage"`
			PvcAssignedName     string   `yaml:"pvcAssignedName"`
			ScheduleOnNode      []string `yaml:"scheduleOnNode,omitempty"`
			DoNotScheduleOnNode []string `yaml:"doNotScheduleOnNode,omitempty"`
		} `yaml:"mediaProxy"`
		MtlManager struct {
			Image      string          `yaml:"image"`
			Resources  bcs.HwResources `yaml:"resources"`
			VolumesMtl struct {
				ImtlHostPath string `yaml:"imtlHostPath"`
				BpfPath      string `yaml:"bpfPath"`
			} `yaml:"volumes"`
			ScheduleOnNode      []string `yaml:"scheduleOnNode,omitempty"`
			DoNotScheduleOnNode []string `yaml:"doNotScheduleOnNode,omitempty"`
		} `yaml:"mtlManager"`
	} `yaml:"definition"`
}

func UnmarshalK8sConfig(yamlData []byte) (*K8sConfig, error) {
	var config K8sConfig
	err := yaml.Unmarshal(yamlData, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}
	return &config, nil
}

func CreateMtlManagerDeployment(cm *corev1.ConfigMap) *appsv1.Deployment {
	data, err := UnmarshalK8sConfig([]byte(cm.Data["config.yaml"]))
	if err != nil {
		fmt.Println("Error unmarshalling K8s config:", err)
		return nil
	}
	// Assign default values if CPU or Memory requests/limits are empty
	if data.Definition.MtlManager.Resources.Requests.CPU == "" {
		data.Definition.MtlManager.Resources.Requests.CPU = "500m"
	}
	if data.Definition.MtlManager.Resources.Requests.Memory == "" {
		data.Definition.MtlManager.Resources.Requests.Memory = "256Mi"
	}
	if data.Definition.MtlManager.Resources.Limits.CPU == "" {
		data.Definition.MtlManager.Resources.Limits.CPU = "1000m"
	}
	if data.Definition.MtlManager.Resources.Limits.Memory == "" {
		data.Definition.MtlManager.Resources.Limits.Memory = "512Mi"
	}

	depl := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mtl-manager",
			Namespace: "mcm",
			Labels: map[string]string{
				"app": "mtl-manager",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "mtl-manager",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "mtl-manager",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "mtl-manager",
							Image: data.Definition.MtlManager.Image,
							SecurityContext: &corev1.SecurityContext{
								Privileged: boolPtr(true),
								Capabilities: &corev1.Capabilities{
									Add: []corev1.Capability{"ALL"},
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(data.Definition.MtlManager.Resources.Requests.CPU),
									corev1.ResourceMemory: resource.MustParse(data.Definition.MtlManager.Resources.Requests.Memory),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(data.Definition.MtlManager.Resources.Limits.CPU),
									corev1.ResourceMemory: resource.MustParse(data.Definition.MtlManager.Resources.Limits.Memory),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "imtl",
									MountPath: "/var/run/imtl",
								},
								{
									Name:      "bpf",
									MountPath: "/sys/fs/bpf",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "imtl",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: data.Definition.MtlManager.VolumesMtl.ImtlHostPath,
								},
							},
						},
						{
							Name: "bpf",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: data.Definition.MtlManager.VolumesMtl.BpfPath,
								},
							},
						},
					},
					HostNetwork: true,
				},
			},
		},
	}
	affinity := &v1.Affinity{}
	AssignNodeAffinityFromConfig(affinity, data.Definition.MtlManager.ScheduleOnNode)
	depl.Spec.Template.Spec.Affinity = affinity
	noAffinity := &v1.PodAntiAffinity{}
	AssignNodeAntiAffinityFromConfig(noAffinity, data.Definition.MtlManager.DoNotScheduleOnNode)
	depl.Spec.Template.Spec.Affinity.PodAntiAffinity = noAffinity
	return depl
}

func CreateMeshAgentDeployment(cm *corev1.ConfigMap) *appsv1.Deployment {
	data, err := UnmarshalK8sConfig([]byte(cm.Data["config.yaml"]))
	if err != nil {
		fmt.Println("Error unmarshalling K8s config:", err)
		return nil
	}
	if data.Definition.MeshAgent.Resources.Requests.CPU == "" {
		data.Definition.MeshAgent.Resources.Requests.CPU = "500m"
	}
	if data.Definition.MeshAgent.Resources.Requests.Memory == "" {
		data.Definition.MeshAgent.Resources.Requests.Memory = "256Mi"
	}
	if data.Definition.MeshAgent.Resources.Limits.CPU == "" {
		data.Definition.MeshAgent.Resources.Limits.CPU = "1000m"
	}
	if data.Definition.MeshAgent.Resources.Limits.Memory == "" {
		data.Definition.MeshAgent.Resources.Limits.Memory = "512Mi"
	}
	fmt.Printf("Data: %+v\n", data)
	deploy := &appsv1.Deployment{
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
							Name:            "mesh-agent",
							Image:           data.Definition.MeshAgent.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command: []string{
								"mesh-agent", "-c", fmt.Sprintf("%d", data.Definition.MeshAgent.RestPort), "-p", fmt.Sprintf("%d", data.Definition.MeshAgent.GrpcPort),
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(data.Definition.MeshAgent.Resources.Requests.CPU),
									corev1.ResourceMemory: resource.MustParse(data.Definition.MeshAgent.Resources.Requests.Memory),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(data.Definition.MeshAgent.Resources.Limits.CPU),
									corev1.ResourceMemory: resource.MustParse(data.Definition.MeshAgent.Resources.Limits.Memory),
								},
							},
							Ports: []corev1.ContainerPort{
								{ContainerPort: int32(data.Definition.MeshAgent.RestPort)},
								{ContainerPort: int32(data.Definition.MeshAgent.GrpcPort)},
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
	affinity := &v1.Affinity{}
	AssignNodeAffinityFromConfig(affinity, data.Definition.MediaProxy.ScheduleOnNode)
	deploy.Spec.Template.Spec.Affinity = affinity
	noAffinity := &v1.PodAntiAffinity{}
	AssignNodeAntiAffinityFromConfig(noAffinity, data.Definition.MediaProxy.DoNotScheduleOnNode)
	deploy.Spec.Template.Spec.Affinity.PodAntiAffinity = noAffinity
	return deploy
}

func CreateMeshAgentService(cm *corev1.ConfigMap) *corev1.Service {
	data, err := UnmarshalK8sConfig([]byte(cm.Data["config.yaml"]))
	if err != nil {
		fmt.Println("Error unmarshalling K8s config:", err)
		return nil
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mesh-agent-service",
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
					Port:       int32(data.Definition.MeshAgent.RestPort),
					TargetPort: intstr.FromInt(data.Definition.MeshAgent.RestPort),
				},
				{
					Name:       "grpc",
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(data.Definition.MeshAgent.GrpcPort),
					TargetPort: intstr.FromInt(data.Definition.MeshAgent.GrpcPort),
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

func CreateBcsService(bcs *bcsv1.BcsConfigSpec) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bcs.Name,
			Namespace: bcs.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Selector: map[string]string{
				"app": bcs.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Name:       "nmos-node-api",
					Port:       int32(bcs.Nmos.NmosInputFile.HttpPort),
					TargetPort: intstr.FromInt(int(bcs.Nmos.NmosInputFile.HttpPort)),
					NodePort:   int32(bcs.Nmos.NmosApiNodePort),
				},
			},
		},
	}
}

func CreatePersistentVolume(cm *corev1.ConfigMap) *corev1.PersistentVolume {
	data, err := UnmarshalK8sConfig([]byte(cm.Data["config.yaml"]))
	if err != nil {
		fmt.Println("Error unmarshalling K8s config:", err)
		return nil
	}
	return &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mtl-pv",
		},
		Spec: corev1.PersistentVolumeSpec{
			Capacity: corev1.ResourceList{
				corev1.ResourceStorage: resource.MustParse(data.Definition.MediaProxy.PvStorage),
			},
			VolumeMode: func() *corev1.PersistentVolumeMode { mode := corev1.PersistentVolumeFilesystem; return &mode }(),
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimRetain,
			StorageClassName:              data.Definition.MediaProxy.PvStorageClass,
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: data.Definition.MediaProxy.PvHostPath,
				},
			},
		},
	}
}

func CreatePersistentVolumeClaim(cm *corev1.ConfigMap) *corev1.PersistentVolumeClaim {
	data, err := UnmarshalK8sConfig([]byte(cm.Data["config.yaml"]))
	if err != nil {
		fmt.Println("Error unmarshalling K8s config:", err)
		return nil
	}
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Definition.MediaProxy.PvcAssignedName,
			Namespace: "mcm",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(data.Definition.MediaProxy.PvcStorage),
				},
			},
			StorageClassName: func() *string { s := data.Definition.MediaProxy.PvStorageClass; return &s }(),
			VolumeName:       "mtl-pv",
		},
	}
}

func CreateConfigMap(bcs *bcsv1.BcsConfigSpec) *corev1.ConfigMap {
	//Override the config that is necessary for the deployment of the NMOS node
	bcs.Nmos.NmosInputFile.FfmpegGrpcServerPort = strconv.Itoa(bcs.App.GrpcPort)
	bcs.Nmos.NmosInputFile.FfmpegGrpcServerAddress = "localhost"
	data, err := json.Marshal(bcs.Nmos.NmosInputFile)
	if err != nil {
		fmt.Println("Error marshalling NMOS input file:", err)
		return nil
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bcs.Name + "-config",
			Namespace: bcs.Namespace,
		},
		Data: map[string]string{
			"config.json": string(data),
		},
	}
}

func CreateBcsDeployment(bcs *bcsv1.BcsConfigSpec) *appsv1.Deployment {

	// Assign default values if CPU or Memory requests/limits are empty for containers
	if bcs.Nmos.Resources.Requests.CPU == "" {
		bcs.Nmos.Resources.Requests.CPU = "200m"
	}
	if bcs.Nmos.Resources.Requests.Memory == "" {
		bcs.Nmos.Resources.Requests.Memory = "256Mi"
	}
	if bcs.Nmos.Resources.Limits.CPU == "" {
		bcs.Nmos.Resources.Limits.CPU = "1000m"
	}
	if bcs.Nmos.Resources.Limits.Memory == "" {
		bcs.Nmos.Resources.Limits.Memory = "512Mi"
	}
	if bcs.App.Resources.Requests.CPU == "" {
		bcs.App.Resources.Requests.CPU = "500m"
	}
	if bcs.App.Resources.Requests.Memory == "" {
		bcs.App.Resources.Requests.Memory = "256Mi"
	}
	if bcs.App.Resources.Limits.CPU == "" {
		bcs.App.Resources.Limits.CPU = "1000m"
	}
	if bcs.App.Resources.Limits.Memory == "" {
		bcs.App.Resources.Limits.Memory = "512Mi"
	}

	bcsDeploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bcs.Name,
			Namespace: bcs.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": bcs.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": bcs.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "tiber-broadcast-suite-nmos-node",
							Image:           bcs.Nmos.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Args:            bcs.Nmos.Args,
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
							Env: convertEnvVars(bcs.Nmos.EnvironmentVariables),
							Ports: []corev1.ContainerPort{
								{ContainerPort: 20000},
								{ContainerPort: 20170},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse(bcs.Nmos.Resources.Requests.Memory),
									corev1.ResourceCPU:    resource.MustParse(bcs.Nmos.Resources.Requests.CPU),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse(bcs.Nmos.Resources.Limits.Memory),
									corev1.ResourceCPU:    resource.MustParse(bcs.Nmos.Resources.Limits.CPU),
								},
							},
						},
						{
							Name:            "tiber-broadcast-suite",
							Image:           bcs.App.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Args:            []string{"localhost", fmt.Sprintf("%d", bcs.App.GrpcPort)},
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
								{Name: "hugepages-2mi", MountPath: "/tmp/hugepages"},
								{Name: "hugepages-1gi", MountPath: "/hugepages"},
								{Name: "imtl", MountPath: "/var/run/imtl"},
								{Name: "shm", MountPath: "/dev/shm"},
								{Name: "dri-dev", MountPath: "/dev/dri"},
								{Name: "vfio", MountPath: "/dev/vfio"},
							},
							Env: convertEnvVars(bcs.App.EnvironmentVariables),
							Ports: []corev1.ContainerPort{
								{ContainerPort: 20000},
								{ContainerPort: 20170},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse(bcs.App.Resources.Requests.Memory),
									corev1.ResourceCPU:    resource.MustParse(bcs.App.Resources.Requests.CPU),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse(bcs.App.Resources.Limits.Memory),
									corev1.ResourceCPU:    resource.MustParse(bcs.App.Resources.Limits.CPU),
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{Name: "videos", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: bcs.App.Volumes["videos"]}}},
						{Name: "dri", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: bcs.App.Volumes["dri"]}}},
						{Name: "kahawai-lock", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: bcs.App.Volumes["kahawaiLock"]}}},
						{Name: "dev-null", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: bcs.App.Volumes["devNull"]}}},
						{Name: "hugepages-2mi", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{Medium: "HugePages-2Mi"}}},
						{Name: "hugepages-1gi", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{Medium: "HugePages-1Gi"}}},
						{Name: "imtl", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: bcs.App.Volumes["imtl"]}}},
						{Name: "shm", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: bcs.App.Volumes["shm"]}}},
						{Name: "vfio", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: bcs.App.Volumes["vfio"]}}},
						{Name: "dri-dev", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: bcs.App.Volumes["dri-dev"]}}},
						{Name: "config", VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: bcs.Name + "-config"}}}},
					},
				},
			},
		},
	}

	// Assign huge pages if they are not empty
	if bcs.App.Resources.Requests.Hugepages1Gi != "" {
		bcsDeploy.Spec.Template.Spec.Containers[1].Resources.Requests[corev1.ResourceHugePagesPrefix+"1Gi"] = resource.MustParse(bcs.App.Resources.Requests.Hugepages1Gi)
	} else {
		bcsDeploy.Spec.Template.Spec.Containers[1].Resources.Requests[corev1.ResourceHugePagesPrefix+"1Gi"] = resource.MustParse("1Gi")
	}
	if bcs.App.Resources.Requests.Hugepages2Mi != "" {
		bcsDeploy.Spec.Template.Spec.Containers[1].Resources.Requests[corev1.ResourceHugePagesPrefix+"2Mi"] = resource.MustParse(bcs.App.Resources.Requests.Hugepages2Mi)
	} else {
		bcsDeploy.Spec.Template.Spec.Containers[1].Resources.Requests[corev1.ResourceHugePagesPrefix+"2Mi"] = resource.MustParse("2Mi")
	}
	if bcs.App.Resources.Limits.Hugepages1Gi != "" {
		bcsDeploy.Spec.Template.Spec.Containers[1].Resources.Limits[corev1.ResourceHugePagesPrefix+"1Gi"] = resource.MustParse(bcs.App.Resources.Limits.Hugepages1Gi)
	} else {
		bcsDeploy.Spec.Template.Spec.Containers[1].Resources.Limits[corev1.ResourceHugePagesPrefix+"1Gi"] = resource.MustParse("1Gi")
	}
	if bcs.App.Resources.Limits.Hugepages2Mi != "" {
		bcsDeploy.Spec.Template.Spec.Containers[1].Resources.Limits[corev1.ResourceHugePagesPrefix+"2Mi"] = resource.MustParse(bcs.App.Resources.Limits.Hugepages2Mi)
	} else {
		bcsDeploy.Spec.Template.Spec.Containers[1].Resources.Limits[corev1.ResourceHugePagesPrefix+"2Mi"] = resource.MustParse("2Mi")
	}

	affinity := &v1.Affinity{}
	AssignNodeAffinityFromConfig(affinity, bcs.ScheduleOnNode)
	bcsDeploy.Spec.Template.Spec.Affinity = affinity
	noAffinity := &v1.PodAntiAffinity{}
	AssignNodeAntiAffinityFromConfig(noAffinity, bcs.DoNotScheduleOnNode)
	bcsDeploy.Spec.Template.Spec.Affinity.PodAntiAffinity = noAffinity
	return bcsDeploy
}

func CreateDaemonSet(cm *corev1.ConfigMap) *appsv1.DaemonSet {
	data, err := UnmarshalK8sConfig([]byte(cm.Data["config.yaml"]))
	if err != nil {
		fmt.Println("Error unmarshalling K8s config:", err)
		return nil
	}
	// Assign default values if any of the resource requests or limits are empty
	if data.Definition.MediaProxy.Resources.Requests.CPU == "" {
		data.Definition.MediaProxy.Resources.Requests.CPU = "2"
	}
	if data.Definition.MediaProxy.Resources.Requests.Memory == "" {
		data.Definition.MediaProxy.Resources.Requests.Memory = "8Gi"
	}
	if data.Definition.MediaProxy.Resources.Requests.Hugepages1Gi == "" {
		data.Definition.MediaProxy.Resources.Requests.Hugepages1Gi = "1Gi"
	}
	if data.Definition.MediaProxy.Resources.Requests.Hugepages2Mi == "" {
		data.Definition.MediaProxy.Resources.Requests.Hugepages2Mi = "2Gi"
	}
	if data.Definition.MediaProxy.Resources.Limits.CPU == "" {
		data.Definition.MediaProxy.Resources.Limits.CPU = "2"
	}
	if data.Definition.MediaProxy.Resources.Limits.Memory == "" {
		data.Definition.MediaProxy.Resources.Limits.Memory = "8Gi"
	}
	if data.Definition.MediaProxy.Resources.Limits.Hugepages1Gi == "" {
		data.Definition.MediaProxy.Resources.Limits.Hugepages1Gi = "1Gi"
	}
	if data.Definition.MediaProxy.Resources.Limits.Hugepages2Mi == "" {
		data.Definition.MediaProxy.Resources.Limits.Hugepages2Mi = "2Gi"
	}
	ds := &appsv1.DaemonSet{
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
					Containers: []corev1.Container{
						{
							Name:            "media-proxy",
							Image:           data.Definition.MediaProxy.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command:         data.Definition.MediaProxy.Command,
							Args:            data.Definition.MediaProxy.Args,
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:                     resource.MustParse(data.Definition.MediaProxy.Resources.Requests.CPU),
									corev1.ResourceMemory:                  resource.MustParse(data.Definition.MediaProxy.Resources.Requests.Memory),
									corev1.ResourceHugePagesPrefix + "2Mi": resource.MustParse(data.Definition.MediaProxy.Resources.Requests.Hugepages2Mi),
									corev1.ResourceHugePagesPrefix + "1Gi": resource.MustParse(data.Definition.MediaProxy.Resources.Requests.Hugepages1Gi),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:                     resource.MustParse(data.Definition.MediaProxy.Resources.Limits.CPU),
									corev1.ResourceMemory:                  resource.MustParse(data.Definition.MediaProxy.Resources.Limits.Memory),
									corev1.ResourceHugePagesPrefix + "2Mi": resource.MustParse(data.Definition.MediaProxy.Resources.Limits.Hugepages2Mi),
									corev1.ResourceHugePagesPrefix + "1Gi": resource.MustParse(data.Definition.MediaProxy.Resources.Limits.Hugepages1Gi),
								},
							},
							SecurityContext: &corev1.SecurityContext{
								Privileged: boolPtr(true),
								RunAsUser:  int64Ptr(0),
								RunAsGroup: int64Ptr(0),
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: int32(data.Definition.MediaProxy.GrpcPort),
									HostPort:      int32(data.Definition.MediaProxy.GrpcPort),
									Protocol:      corev1.ProtocolTCP,
									Name:          "grpc-port",
								},
								{
									ContainerPort: int32(data.Definition.MediaProxy.SdkPort),
									HostPort:      int32(data.Definition.MediaProxy.SdkPort),
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
								{
									Name:      "hugepage-2mi",
									MountPath: "/hugepages-2Mi",
								},
								{
									Name:      "hugepage-1gi",
									MountPath: "/hugepages-1Gi",
								},
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
									Path: data.Definition.MediaProxy.Volumes.Memif,
								},
							},
						},
						{
							Name: "dev-vfio",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: data.Definition.MediaProxy.Volumes.Vfio,
								},
							},
						},
						{
							Name: "hugepage-2mi",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{
									Medium: "HugePages-2Mi",
								},
							},
						},
						{
							Name: "hugepage-1gi",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{
									Medium: "HugePages-1Gi",
								},
							},
						},
						{
							Name: "cache-volume",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{
									Medium:    corev1.StorageMediumMemory,
									SizeLimit: func(q resource.Quantity) *resource.Quantity { return &q }(resource.MustParse(data.Definition.MediaProxy.Volumes.CacheSize)),
								},
							},
						},
						{
							Name: "mtl-mgr",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: data.Definition.MediaProxy.PvcAssignedName,
								},
							},
						},
					},
				},
			},
		},
	}

	affinity := &v1.Affinity{}
	AssignNodeAffinityFromConfig(affinity, data.Definition.MediaProxy.ScheduleOnNode)
	ds.Spec.Template.Spec.Affinity = affinity
	noAffinity := &v1.PodAntiAffinity{}
	AssignNodeAntiAffinityFromConfig(noAffinity, data.Definition.MediaProxy.DoNotScheduleOnNode)
	ds.Spec.Template.Spec.Affinity.PodAntiAffinity = noAffinity
	return ds
}

func AssignNodeAffinityFromConfig(affinity *corev1.Affinity, scheduleOnNode []string) {
	if len(scheduleOnNode) == 0 {
		return
	}

	var nodeSelectorTerms []corev1.NodeSelectorTerm
	for _, keyValue := range scheduleOnNode {
		parts := strings.SplitN(keyValue, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]
		nodeSelectorTerms = append(nodeSelectorTerms, corev1.NodeSelectorTerm{
			MatchExpressions: []corev1.NodeSelectorRequirement{
				{
					Key:      key,
					Operator: corev1.NodeSelectorOpIn,
					Values:   []string{value},
				},
			},
		})
	}

	if len(nodeSelectorTerms) > 0 {
		if affinity.NodeAffinity == nil {
			affinity.NodeAffinity = &corev1.NodeAffinity{}
		}
		affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = &corev1.NodeSelector{
			NodeSelectorTerms: nodeSelectorTerms,
		}
	}
}

func AssignNodeAntiAffinityFromConfig(antiAffinity *corev1.PodAntiAffinity, doNotScheduleOnNode []string) {
	if len(doNotScheduleOnNode) == 0 {
		return
	}

	var podAffinityTerms []corev1.PodAffinityTerm
	for _, keyValue := range doNotScheduleOnNode {
		parts := strings.SplitN(keyValue, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]
		podAffinityTerms = append(podAffinityTerms, corev1.PodAffinityTerm{
			LabelSelector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      key,
						Operator: metav1.LabelSelectorOpNotIn,
						Values:   []string{value},
					},
				},
			},
			TopologyKey: "kubernetes.io/hostname",
		})
	}

	if len(podAffinityTerms) > 0 {
		if antiAffinity == nil {
			antiAffinity = &corev1.PodAntiAffinity{}
		}
		antiAffinity.RequiredDuringSchedulingIgnoredDuringExecution = append(antiAffinity.RequiredDuringSchedulingIgnoredDuringExecution, podAffinityTerms...)
	}
}

func convertEnvVars(envVars []bcsv1.EnvVar) []corev1.EnvVar {
	var coreEnvVars []corev1.EnvVar
	for _, envVar := range envVars {
		coreEnvVars = append(coreEnvVars, corev1.EnvVar{
			Name:  envVar.Name,
			Value: envVar.Value,
		})
	}
	return coreEnvVars
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
