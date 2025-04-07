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
	"sync"

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

var fileMutex sync.Mutex

func updateNmosJsonFile(filePath string, ip string, port string) error {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

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

func ConstructContainerConfig(containerInfo *general.Containers, log logr.Logger) (*container.Config, *container.HostConfig, *network.NetworkingConfig) {
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
				{Type: mount.TypeBind, Source: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Dri, Target: "/usr/local/lib/x86_64-linux-gnu/dri"},
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
			{PathOnHost: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Devices.Dri, PathInContainer: "/dev/dri"},
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
		nmosFileNameJson := containerInfo.Configuration.WorkloadConfig.NmosClient.NmosConfigFileName
		nmosFilePathJson := containerInfo.Configuration.WorkloadConfig.NmosClient.NmosConfigPath + "/" + nmosFileNameJson
		if !FileExists(nmosFilePathJson){
			log.Error(errors.New("NMOS json file does not exist"), "NMOS json file does not exist")
			return nil, nil, nil
		}
		errUpdateJson := updateNmosJsonFile(nmosFilePathJson,
			containerInfo.Configuration.WorkloadConfig.NmosClient.FfmpegConectionAddress,
			containerInfo.Configuration.WorkloadConfig.NmosClient.FfmpegConnectionPort)
		if errUpdateJson != nil {
			log.Error(errUpdateJson, "Error updating NMOS json file")
			return nil, nil, nil
		}
		configPathContainer := "config/" + nmosFileNameJson
		containerConfig = &container.Config{
			Image: containerInfo.Configuration.WorkloadConfig.NmosClient.ImageAndTag,
			Cmd: []string{configPathContainer},
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

func boolPtr(b bool) *bool    { return &b }

type K8sConfig struct {
	K8s        bool `yaml:"k8s"`
	Definition struct {
		MeshAgent struct {
			Image    string `yaml:"image"`
			RestPort int    `yaml:"restPort"`
			GrpcPort int    `yaml:"grpcPort"`
		} `yaml:"meshAgent"`
		MediaProxy struct {
			Image       string   `yaml:"image"`
			Command     []string `yaml:"command"`
			Args        []string `yaml:"args"`
			GrpcPort    int      `yaml:"grpcPort"`
			SdkPort     int      `yaml:"sdkPort"`
			Volumes     struct {
				Memif string `yaml:"memif"`
				Vfio  string `yaml:"vfio"`
			} `yaml:"volumes"`
			PvHostPath     string `yaml:"pvHostPath"`
			PvStorageClass string `yaml:"pvStorageClass"`
			PvStorage      string `yaml:"pvStorage"`
			PvcStorage     string `yaml:"pvcStorage"`
		} `yaml:"mediaProxy"`
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

func CreateMeshAgentDeployment(cm *corev1.ConfigMap) *appsv1.Deployment {
	 data, err := UnmarshalK8sConfig([]byte(cm.Data["config.yaml"]))
	 if err != nil {
		 fmt.Println("Error unmarshalling K8s config:", err)
		 return nil
	 }
	 fmt.Printf("Data: %+v\n", data)
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
							Image: data.Definition.MeshAgent.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command: []string{
								"mesh-agent", "-c", fmt.Sprintf("%d", data.Definition.MeshAgent.RestPort), "-p", fmt.Sprintf("%d", data.Definition.MeshAgent.GrpcPort),
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
}

func CreateMeshAgentService(cm *corev1.ConfigMap) *corev1.Service {
	data, err := UnmarshalK8sConfig([]byte(cm.Data["config.yaml"]))
	 if err != nil {
		 fmt.Println("Error unmarshalling K8s config:", err)
		 return nil
	 }
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
					Port:       int32(data.Definition.MeshAgent.RestPort),
					TargetPort: intstr.FromInt(data.Definition.MeshAgent.RestPort),
				},
				{
					Name: 	    "grpc",
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

func CreateBcsService(bcs * bcsv1.BcsConfig) *corev1.Service {
	return &corev1.Service{
        ObjectMeta: metav1.ObjectMeta{
          Name: bcs.Spec.Name,
		  Namespace: bcs.Spec.Namespace,
        },
        Spec: corev1.ServiceSpec{
          Type: corev1.ServiceTypeNodePort,
          Selector: map[string]string{
            "app": bcs.Spec.Name,
          },
          Ports: []corev1.ServicePort{
            {
              Protocol:   corev1.ProtocolTCP,
			  Name: "nmos-node-api",
              Port:       int32(bcs.Spec.Nmos.NmosApiPort),
			  TargetPort: intstr.FromInt(int(bcs.Spec.Nmos.NmosApiPort)),
			  NodePort:   int32(bcs.Spec.Nmos.NmosApiNodePort),
            },
            {
              Protocol:   corev1.ProtocolTCP,
			  Name: "nmos-app-communication",
              Port:       int32(bcs.Spec.Nmos.NmosAppCommunicationPort),
			  TargetPort: intstr.FromInt(int(bcs.Spec.Nmos.NmosAppCommunicationPort)),
			  NodePort:   int32(bcs.Spec.Nmos.NmosAppCommunicationNodePort),
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
			Name:      "mtl-pvc",
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

func CreateConfigMap(bcs *bcsv1.BcsConfig) *corev1.ConfigMap {
	
	data, err := json.Marshal(bcs.Spec.Nmos.NmosInputFile)
	if err != nil {
		fmt.Println("Error marshalling NMOS input file:", err)
		return nil
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bcs.Spec.Name + "-config",
			Namespace: bcs.Spec.Namespace,
		},
		Data: map[string]string{
			"config": string(data),
		},
	}
}

func CreateBcsDeployment(bcs *bcsv1.BcsConfig) *appsv1.Deployment {
	return &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
          Name:      bcs.Spec.Name,
          Namespace: bcs.Spec.Namespace,
        },
        Spec: appsv1.DeploymentSpec{
          Replicas: int32Ptr(1),
          Selector: &metav1.LabelSelector{
            MatchLabels: map[string]string{
              "app": bcs.Spec.Name,
            },
          },
          Template: corev1.PodTemplateSpec{
            ObjectMeta: metav1.ObjectMeta{
              Labels: map[string]string{
                "app": bcs.Spec.Name,
              },
            },
            Spec: corev1.PodSpec{
              Containers: []corev1.Container{
                {
                  Name:  "tiber-broadcast-suite-nmos-node",
                  Image: bcs.Spec.Nmos.Image,
				  ImagePullPolicy: corev1.PullIfNotPresent,
                  Args:  bcs.Spec.Nmos.Args,
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
				  Env: convertEnvVars(bcs.Spec.Nmos.EnvironmentVariables),
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
                  Image: bcs.Spec.App.Image,
				  ImagePullPolicy: corev1.PullIfNotPresent,
				  Args:  []string{"localhost", fmt.Sprintf("%d", bcs.Spec.App.GrpcPort)},
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
                  Env: convertEnvVars(bcs.Spec.App.EnvironmentVariables),
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
                {Name: "videos", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: bcs.Spec.App.Volumes["videos"]}}},
                {Name: "dri", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: bcs.Spec.App.Volumes["dri"]}}},
                {Name: "kahawai-lock", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: bcs.Spec.App.Volumes["kahawaiLock"]}}},
                {Name: "dev-null", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: bcs.Spec.App.Volumes["devNull"]}}},
                {Name: "hugepages-tmp", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: bcs.Spec.App.Volumes["hugepagesTmp"]}}},
                {Name: "hugepages", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: bcs.Spec.App.Volumes["hugepages"]}}},
                {Name: "imtl", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: bcs.Spec.App.Volumes["imtl"]}}},
                {Name: "shm", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: bcs.Spec.App.Volumes["shm"]}}},
                {Name: "vfio", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: bcs.Spec.App.Volumes["vfio"]}}},
                {Name: "driDev", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: bcs.Spec.App.Volumes["driDev"]}}},
				{Name: "config", VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: bcs.Spec.Name+"-config"}}}},
              },
            },
          },
        },
      }
}

func CreateDaemonSet(cm *corev1.ConfigMap) *appsv1.DaemonSet {
    data, err := UnmarshalK8sConfig([]byte(cm.Data["config.yaml"]))
	 if err != nil {
		 fmt.Println("Error unmarshalling K8s config:", err)
		 return nil
	 }
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
					NodeSelector: map[string]string{
						"node-role.kubernetes.io/worker": "true",
					},
					Containers: []corev1.Container{
						{
							Name:    "media-proxy",
							Image:   data.Definition.MediaProxy.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command: data.Definition.MediaProxy.Command,
							Args:    data.Definition.MediaProxy.Args,
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
