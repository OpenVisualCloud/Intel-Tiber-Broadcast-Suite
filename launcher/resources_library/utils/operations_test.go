package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	bcsv1 "bcs.pod.launcher.intel/api/v1"

	"bcs.pod.launcher.intel/resources_library/parser"
	"bcs.pod.launcher.intel/resources_library/resources/general"

	"bcs.pod.launcher.intel/resources_library/resources/nmos"
	"bcs.pod.launcher.intel/resources_library/workloads"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestFileExists(t *testing.T) {
	tempFile, err := os.CreateTemp("", "file_exists_test_*.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	assert.True(t, FileExists(tempFile.Name()))
	assert.False(t, FileExists("non-existent-file.txt"))
}

func TestUpdateNmosJsonFile_ValidFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "nmos_test_*.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	jsonData := `{
        "ffmpeg_grpc_server_address": "old-address",
        "ffmpeg_grpc_server_port": "50051"
    }`
	_, err = tempFile.Write([]byte(jsonData))
	assert.NoError(t, err)
	tempFile.Close()

	err = updateNmosJsonFile(tempFile.Name(), "new-address", "50052")
	assert.NoError(t, err)

	updatedData, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err)

	var config map[string]interface{}
	err = json.Unmarshal(updatedData, &config)
	assert.NoError(t, err)
	assert.Equal(t, "new-address", config["ffmpeg_grpc_server_address"])
	assert.Equal(t, "50052", fmt.Sprintf("%v", config["ffmpeg_grpc_server_port"]))
}

func TestUpdateNmosJsonFile_FileNotFound(t *testing.T) {
	err := updateNmosJsonFile("non-existent-file.json", "new-address", "50052")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestUpdateNmosJsonFile_InvalidJson(t *testing.T) {
	tempFile, err := os.CreateTemp("", "nmos_test_invalid_*.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte(`{invalid-json}`))
	assert.NoError(t, err)
	tempFile.Close()

	err = updateNmosJsonFile(tempFile.Name(), "new-address", "50052")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestConstructContainerConfig(t *testing.T) {
	log := logr.Discard()

	t.Run("MediaProxyAgentConfig", func(t *testing.T) {
		containerInfo := &general.Containers{Type: general.Workload(general.NotSupportedWorkload)}
		config := &parser.Configuration{
			RunOnce: parser.RunOnce{},
		}

		containerConfig, hostConfig, networkConfig := ConstructContainerConfig(containerInfo, config, log)

		assert.Nil(t, containerConfig)
		assert.Nil(t, hostConfig)
		assert.Nil(t, networkConfig)
	})

	t.Run("MediaProxyAgentConfig", func(t *testing.T) {
		containerInfo := &general.Containers{Type: general.MediaProxyAgent}
		config := &parser.Configuration{
			RunOnce: parser.RunOnce{
				MediaProxyAgent: workloads.MediaProxyAgentConfig{
					ImageAndTag: "mediaproxyagent:latest",
					RestPort:    "8080",
					GRPCPort:    "9090",
					Network: workloads.NetworkConfig{
						Enable: true,
						Name:   "test-network",
						IP:     "192.168.1.100",
					},
				},
			},
		}

		containerConfig, hostConfig, networkConfig := ConstructContainerConfig(containerInfo, config, log)

		assert.NotNil(t, containerConfig)
		assert.NotNil(t, hostConfig)
		assert.NotNil(t, networkConfig)

		assert.Equal(t, "mediaproxyagent:latest", containerConfig.Image)
		assert.Contains(t, containerConfig.Cmd, "8080")
		assert.Contains(t, containerConfig.Cmd, "9090")
		assert.Equal(t, "test-network", string(hostConfig.NetworkMode))
		assert.Equal(t, "192.168.1.100", networkConfig.EndpointsConfig["test-network"].IPAMConfig.IPv4Address)
	})

	t.Run("MediaProxyAgentConfigHostNetwork", func(t *testing.T) {
		containerInfo := &general.Containers{Type: general.MediaProxyAgent}
		config := &parser.Configuration{
			RunOnce: parser.RunOnce{
				MediaProxyAgent: workloads.MediaProxyAgentConfig{
					ImageAndTag: "mediaproxyagent:latest",
					RestPort:    "8080",
					GRPCPort:    "9090",
					Network: workloads.NetworkConfig{
						Enable: false,
						IP:     "192.168.1.100",
					},
				},
			},
		}

		containerConfig, hostConfig, networkConfig := ConstructContainerConfig(containerInfo, config, log)

		assert.NotNil(t, containerConfig)
		assert.NotNil(t, hostConfig)
		assert.NotNil(t, networkConfig)

		assert.Equal(t, "mediaproxyagent:latest", containerConfig.Image)
		assert.Contains(t, containerConfig.Cmd, "8080")
		assert.Contains(t, containerConfig.Cmd, "9090")
		assert.Equal(t, "host", string(hostConfig.NetworkMode))
		assert.Empty(t, networkConfig.EndpointsConfig)
	})

	t.Run("MediaProxyMCMConfig", func(t *testing.T) {
		containerInfo := &general.Containers{Type: general.MediaProxyMCM}
		config := &parser.Configuration{
			RunOnce: parser.RunOnce{
				MediaProxyMcm: workloads.MediaProxyMcmConfig{
					ImageAndTag:   "mediaproxymcm:latest",
					InterfaceName: "eth0",
					Volumes:       []string{"/host/path:/container/path"},
					Network: workloads.NetworkConfig{
						Enable: true,
						Name:   "test-network",
						IP:     "192.168.1.101",
					},
				},
			},
		}

		containerConfig, hostConfig, networkConfig := ConstructContainerConfig(containerInfo, config, log)

		assert.NotNil(t, containerConfig)
		assert.NotNil(t, hostConfig)
		assert.NotNil(t, networkConfig)

		assert.Equal(t, "mediaproxymcm:latest", containerConfig.Image)
		assert.Contains(t, containerConfig.Cmd, "kernel:eth0")
		assert.Equal(t, "test-network", string(hostConfig.NetworkMode))
		assert.Equal(t, "192.168.1.101", networkConfig.EndpointsConfig["test-network"].IPAMConfig.IPv4Address)
	})

	t.Run("MediaProxyMCMConfig", func(t *testing.T) {
		containerInfo := &general.Containers{Type: general.MediaProxyMCM}
		config := &parser.Configuration{
			RunOnce: parser.RunOnce{
				MediaProxyMcm: workloads.MediaProxyMcmConfig{
					ImageAndTag:   "mediaproxymcm:latest",
					InterfaceName: "eth0",
					Volumes:       []string{"/host/path:/container/path"},
					Network: workloads.NetworkConfig{
						Enable: false,
						IP:     "192.168.1.101",
					},
				},
			},
		}

		containerConfig, hostConfig, networkConfig := ConstructContainerConfig(containerInfo, config, log)

		assert.NotNil(t, containerConfig)
		assert.NotNil(t, hostConfig)
		assert.NotNil(t, networkConfig)

		assert.Equal(t, "mediaproxymcm:latest", containerConfig.Image)
		assert.Contains(t, containerConfig.Cmd, "kernel:eth0")
		assert.Equal(t, "host", string(hostConfig.NetworkMode))
		assert.Empty(t, networkConfig.EndpointsConfig)
	})

	t.Run("BcsPipelineFfmpegConfig", func(t *testing.T) {
		containerInfo := &general.Containers{Type: general.BcsPipelineFfmpeg, Id: 0}
		config := &parser.Configuration{
			WorkloadToBeRun: []workloads.WorkloadConfig{
				{
					FfmpegPipeline: workloads.FfmpegPipelineConfig{
						ImageAndTag: "ffmpegpipeline:latest",
						GRPCPort:    50051,
						Network: workloads.NetworkConfig{
							Enable: true,
							Name:   "test-network",
							IP:     "192.168.1.102",
						},
						Volumes: workloads.Volumes{
							Videos: "/host/videos",
							Dri:    "/host/dri",
						},
					},
				},
			},
		}

		containerConfig, hostConfig, networkConfig := ConstructContainerConfig(containerInfo, config, log)

		assert.NotNil(t, containerConfig)
		assert.NotNil(t, hostConfig)
		assert.NotNil(t, networkConfig)

		assert.Equal(t, "ffmpegpipeline:latest", containerConfig.Image)
		assert.Contains(t, containerConfig.Cmd, "192.168.1.102")
		assert.Contains(t, containerConfig.Cmd, "50051")
		assert.Equal(t, "test-network", string(hostConfig.NetworkMode))
		assert.Equal(t, "192.168.1.102", networkConfig.EndpointsConfig["test-network"].IPAMConfig.IPv4Address)
		assert.Equal(t, "/host/videos", hostConfig.Mounts[0].Source)
		assert.Equal(t, "/videos", hostConfig.Mounts[0].Target)
	})

	t.Run("BcsPipelineFfmpegConfigHostNetwork", func(t *testing.T) {
		containerInfo := &general.Containers{Type: general.BcsPipelineFfmpeg, Id: 0}
		config := &parser.Configuration{
			WorkloadToBeRun: []workloads.WorkloadConfig{
				{
					FfmpegPipeline: workloads.FfmpegPipelineConfig{
						ImageAndTag: "ffmpegpipeline:latest",
						GRPCPort:    50051,
						Network: workloads.NetworkConfig{
							Enable: false,
							IP:     "192.168.1.102",
						},
						Volumes: workloads.Volumes{
							Videos: "/host/videos",
							Dri:    "/host/dri",
						},
					},
				},
			},
		}

		containerConfig, hostConfig, networkConfig := ConstructContainerConfig(containerInfo, config, log)

		assert.NotNil(t, containerConfig)
		assert.NotNil(t, hostConfig)
		assert.NotNil(t, networkConfig)

		assert.Equal(t, "ffmpegpipeline:latest", containerConfig.Image)
		assert.Equal(t, "host", string(hostConfig.NetworkMode))
		assert.Empty(t, networkConfig.EndpointsConfig)
	})

}
func TestUnmarshalK8sConfig_ValidYAML(t *testing.T) {
	yamlData := `
k8s: true
definition:
  meshAgent:
    image: "mesh-agent:latest"
    restPort: 8080
    grpcPort: 9090
    resources:
      requests:
        cpu: "500m"
        memory: "256Mi"
      limits:
        cpu: "1000m"
        memory: "512Mi"
    scheduleOnNode:
      - "key1=value1"
    doNotScheduleOnNode:
      - "key2=value2"
  mediaProxy:
    image: "media-proxy:latest"
    grpcPort: 50051
    sdkPort: 50052
    resources:
      requests:
        cpu: "2"
        memory: "8Gi"
      limits:
        cpu: "4"
        memory: "16Gi"
    volumes:
      memif: "/run/memif"
      vfio: "/dev/vfio"
      cache-size: "2Gi"
    pvHostPath: "/mnt/hostpath"
    pvStorageClass: "standard"
    pvStorage: "10Gi"
    pvcStorage: "5Gi"
    pvcAssignedName: "media-proxy-pvc"
    scheduleOnNode:
      - "key3=value3"
    doNotScheduleOnNode:
      - "key4=value4"
  mtlManager:
    image: "mtl-manager:latest"
    resources:
      requests:
        cpu: "1"
        memory: "512Mi"
      limits:
        cpu: "2"
        memory: "1Gi"
    volumes:
      imtlHostPath: "/var/run/imtl"
      bpfPath: "/sys/fs/bpf"
    scheduleOnNode:
      - "key5=value5"
    doNotScheduleOnNode:
      - "key6=value6"
`
	config, err := UnmarshalK8sConfig([]byte(yamlData))
	assert.NoError(t, err)
	assert.NotNil(t, config)

	assert.True(t, config.K8s)
	assert.Equal(t, "mesh-agent:latest", config.Definition.MeshAgent.Image)
	assert.Equal(t, 8080, config.Definition.MeshAgent.RestPort)
	assert.Equal(t, 9090, config.Definition.MeshAgent.GrpcPort)
	assert.Equal(t, "500m", config.Definition.MeshAgent.Resources.Requests.CPU)
	assert.Equal(t, "256Mi", config.Definition.MeshAgent.Resources.Requests.Memory)
	assert.Equal(t, "1000m", config.Definition.MeshAgent.Resources.Limits.CPU)
	assert.Equal(t, "512Mi", config.Definition.MeshAgent.Resources.Limits.Memory)
	assert.Contains(t, config.Definition.MeshAgent.ScheduleOnNode, "key1=value1")
	assert.Contains(t, config.Definition.MeshAgent.DoNotScheduleOnNode, "key2=value2")

	assert.Equal(t, "media-proxy:latest", config.Definition.MediaProxy.Image)
	assert.Equal(t, 50051, config.Definition.MediaProxy.GrpcPort)
	assert.Equal(t, 50052, config.Definition.MediaProxy.SdkPort)
	assert.Equal(t, "2", config.Definition.MediaProxy.Resources.Requests.CPU)
	assert.Equal(t, "8Gi", config.Definition.MediaProxy.Resources.Requests.Memory)
	assert.Equal(t, "4", config.Definition.MediaProxy.Resources.Limits.CPU)
	assert.Equal(t, "16Gi", config.Definition.MediaProxy.Resources.Limits.Memory)
	assert.Equal(t, "/run/memif", config.Definition.MediaProxy.Volumes.Memif)
	assert.Equal(t, "/dev/vfio", config.Definition.MediaProxy.Volumes.Vfio)
	assert.Equal(t, "2Gi", config.Definition.MediaProxy.Volumes.CacheSize)
	assert.Equal(t, "/mnt/hostpath", config.Definition.MediaProxy.PvHostPath)
	assert.Equal(t, "standard", config.Definition.MediaProxy.PvStorageClass)
	assert.Equal(t, "10Gi", config.Definition.MediaProxy.PvStorage)
	assert.Equal(t, "5Gi", config.Definition.MediaProxy.PvcStorage)
	assert.Equal(t, "media-proxy-pvc", config.Definition.MediaProxy.PvcAssignedName)
	assert.Contains(t, config.Definition.MediaProxy.ScheduleOnNode, "key3=value3")
	assert.Contains(t, config.Definition.MediaProxy.DoNotScheduleOnNode, "key4=value4")

	assert.Equal(t, "mtl-manager:latest", config.Definition.MtlManager.Image)
	assert.Equal(t, "1", config.Definition.MtlManager.Resources.Requests.CPU)
	assert.Equal(t, "512Mi", config.Definition.MtlManager.Resources.Requests.Memory)
	assert.Equal(t, "2", config.Definition.MtlManager.Resources.Limits.CPU)
	assert.Equal(t, "1Gi", config.Definition.MtlManager.Resources.Limits.Memory)
	assert.Equal(t, "/var/run/imtl", config.Definition.MtlManager.VolumesMtl.ImtlHostPath)
	assert.Equal(t, "/sys/fs/bpf", config.Definition.MtlManager.VolumesMtl.BpfPath)
	assert.Contains(t, config.Definition.MtlManager.ScheduleOnNode, "key5=value5")
	assert.Contains(t, config.Definition.MtlManager.DoNotScheduleOnNode, "key6=value6")
}

func TestUnmarshalK8sConfig_InvalidYAML(t *testing.T) {
	invalidYamlData := `
k8s: true
definition:
  meshAgent:
    image: "mesh-agent:latest"
    restPort: "invalid-port"
`
	config, err := UnmarshalK8sConfig([]byte(invalidYamlData))
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "failed to unmarshal YAML")
}
func TestUnmarshalK8sConfig_ValidData(t *testing.T) {
	yamlData := `
k8s: true
definition:
  meshAgent:
    image: "mesh-agent:latest"
    restPort: 8080
    grpcPort: 9090
    resources:
      requests:
        cpu: "500m"
        memory: "256Mi"
      limits:
        cpu: "1000m"
        memory: "512Mi"
  mediaProxy:
    image: "media-proxy:latest"
    grpcPort: 50051
    sdkPort: 50052
    resources:
      requests:
        cpu: "2"
        memory: "8Gi"
      limits:
        cpu: "4"
        memory: "16Gi"
  mtlManager:
    image: "mtl-manager:latest"
    resources:
      requests:
        cpu: "1"
        memory: "512Mi"
      limits:
        cpu: "2"
        memory: "1Gi"
`
	config, err := UnmarshalK8sConfig([]byte(yamlData))
	assert.NoError(t, err)
	assert.NotNil(t, config)

	assert.True(t, config.K8s)
	assert.Equal(t, "mesh-agent:latest", config.Definition.MeshAgent.Image)
	assert.Equal(t, 8080, config.Definition.MeshAgent.RestPort)
	assert.Equal(t, 9090, config.Definition.MeshAgent.GrpcPort)
	assert.Equal(t, "500m", config.Definition.MeshAgent.Resources.Requests.CPU)
	assert.Equal(t, "256Mi", config.Definition.MeshAgent.Resources.Requests.Memory)
	assert.Equal(t, "1000m", config.Definition.MeshAgent.Resources.Limits.CPU)
	assert.Equal(t, "512Mi", config.Definition.MeshAgent.Resources.Limits.Memory)

	assert.Equal(t, "media-proxy:latest", config.Definition.MediaProxy.Image)
	assert.Equal(t, 50051, config.Definition.MediaProxy.GrpcPort)
	assert.Equal(t, 50052, config.Definition.MediaProxy.SdkPort)
	assert.Equal(t, "2", config.Definition.MediaProxy.Resources.Requests.CPU)
	assert.Equal(t, "8Gi", config.Definition.MediaProxy.Resources.Requests.Memory)
	assert.Equal(t, "4", config.Definition.MediaProxy.Resources.Limits.CPU)
	assert.Equal(t, "16Gi", config.Definition.MediaProxy.Resources.Limits.Memory)

	assert.Equal(t, "mtl-manager:latest", config.Definition.MtlManager.Image)
	assert.Equal(t, "1", config.Definition.MtlManager.Resources.Requests.CPU)
	assert.Equal(t, "512Mi", config.Definition.MtlManager.Resources.Requests.Memory)
	assert.Equal(t, "2", config.Definition.MtlManager.Resources.Limits.CPU)
	assert.Equal(t, "1Gi", config.Definition.MtlManager.Resources.Limits.Memory)
}

func TestUnmarshalK8sConfig_InvalidData(t *testing.T) {
	invalidYamlData := `
k8s: true
definition:
  meshAgent:
    image: "mesh-agent:latest"
    restPort: "invalid-port"
`
	config, err := UnmarshalK8sConfig([]byte(invalidYamlData))
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "failed to unmarshal YAML")
}
func TestCreateMtlManagerDeployment(t *testing.T) {
	t.Run("ValidConfigMap", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			Data: map[string]string{
				"config.yaml": `
k8s: true
definition:
  mtlManager:
    image: "mtl-manager:latest"
    resources:
      requests:
        cpu: "500m"
        memory: "256Mi"
      limits:
        cpu: "1000m"
        memory: "512Mi"
    volumes:
      imtlHostPath: "/var/run/imtl"
      bpfPath: "/sys/fs/bpf"
    scheduleOnNode:
      - "key1=value1"
    doNotScheduleOnNode:
      - "key2=value2"
`,
			},
		}

		deployment := CreateMtlManagerDeployment(cm)
		assert.NotNil(t, deployment)
		assert.Equal(t, "mtl-manager", deployment.ObjectMeta.Name)
		assert.Equal(t, "mcm", deployment.ObjectMeta.Namespace)
		assert.Equal(t, "mtl-manager:latest", deployment.Spec.Template.Spec.Containers[0].Image)

		resources := deployment.Spec.Template.Spec.Containers[0].Resources
		assert.Equal(t, resource.MustParse("500m"), resources.Requests[corev1.ResourceCPU])
		assert.Equal(t, resource.MustParse("256Mi"), resources.Requests[corev1.ResourceMemory])
		assert.Equal(t, resource.MustParse("1000m"), resources.Limits[corev1.ResourceCPU])
		assert.Equal(t, resource.MustParse("512Mi"), resources.Limits[corev1.ResourceMemory])

		volumeMounts := deployment.Spec.Template.Spec.Containers[0].VolumeMounts
		assert.Len(t, volumeMounts, 2)
		assert.Equal(t, "/var/run/imtl", volumeMounts[0].MountPath)
		assert.Equal(t, "/sys/fs/bpf", volumeMounts[1].MountPath)

		affinity := deployment.Spec.Template.Spec.Affinity
		assert.NotNil(t, affinity)
		assert.NotNil(t, affinity.NodeAffinity)
		assert.NotNil(t, affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution)
		assert.Len(t, affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms, 1)
		assert.Equal(t, "key1", affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Key)
		assert.Contains(t, affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Values, "value1")

		podAntiAffinity := affinity.PodAntiAffinity
		assert.NotNil(t, podAntiAffinity)
		assert.Len(t, podAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution, 1)
		assert.Equal(t, "key2", podAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution[0].LabelSelector.MatchExpressions[0].Key)
		assert.Contains(t, podAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution[0].LabelSelector.MatchExpressions[0].Values, "value2")
	})
}
func TestCreateMeshAgentService(t *testing.T) {
	t.Run("ValidConfigMap", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			Data: map[string]string{
				"config.yaml": `
k8s: true
definition:
  meshAgent:
    image: "mesh-agent:latest"
    restPort: 8080
    grpcPort: 9090
`,
			},
		}

		service := CreateMeshAgentService(cm)
		assert.NotNil(t, service)

		assert.Equal(t, "mesh-agent-service", service.ObjectMeta.Name)
		assert.Equal(t, "mcm", service.ObjectMeta.Namespace)
		assert.Equal(t, map[string]string{"app": "mesh-agent"}, service.Spec.Selector)
		assert.Len(t, service.Spec.Ports, 2)
		assert.Equal(t, "rest", service.Spec.Ports[0].Name)
		assert.Equal(t, corev1.ProtocolTCP, service.Spec.Ports[0].Protocol)
		assert.Equal(t, int32(8080), service.Spec.Ports[0].Port)
		assert.Equal(t, intstr.FromInt(8080), service.Spec.Ports[0].TargetPort)

		assert.Equal(t, "grpc", service.Spec.Ports[1].Name)
		assert.Equal(t, corev1.ProtocolTCP, service.Spec.Ports[1].Protocol)
		assert.Equal(t, int32(9090), service.Spec.Ports[1].Port)
		assert.Equal(t, intstr.FromInt(9090), service.Spec.Ports[1].TargetPort)
	})

	t.Run("InvalidConfigMap", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			Data: map[string]string{
				"config.yaml": `
k8s: true
definition:
  meshAgent:
    image: "mesh-agent:latest"
    restPort: "invalid-port"
`,
			},
		}

		service := CreateMeshAgentService(cm)
		assert.Nil(t, service)
	})
}
func TestCreateService(t *testing.T) {
	t.Run("ValidServiceName", func(t *testing.T) {
		serviceName := "test-service"
		service := CreateService(serviceName)

		assert.NotNil(t, service)
		assert.Equal(t, serviceName, service.ObjectMeta.Name)
		assert.Equal(t, "default", service.ObjectMeta.Namespace)
		assert.Equal(t, map[string]string{"app": serviceName}, service.Spec.Selector)
		assert.Len(t, service.Spec.Ports, 1)
		assert.Equal(t, corev1.ProtocolTCP, service.Spec.Ports[0].Protocol)
		assert.Equal(t, int32(80), service.Spec.Ports[0].Port)
	})

	t.Run("EmptyServiceName", func(t *testing.T) {
		serviceName := ""
		service := CreateService(serviceName)

		assert.NotNil(t, service)
		assert.Equal(t, serviceName, service.ObjectMeta.Name)
		assert.Equal(t, "default", service.ObjectMeta.Namespace)
		assert.Equal(t, map[string]string{"app": serviceName}, service.Spec.Selector)
		assert.Len(t, service.Spec.Ports, 1)
		assert.Equal(t, corev1.ProtocolTCP, service.Spec.Ports[0].Protocol)
		assert.Equal(t, int32(80), service.Spec.Ports[0].Port)
	})
}
func TestCreateBcsService(t *testing.T) {
	t.Run("ValidBcsConfigSpec", func(t *testing.T) {
		bcsConfig := &bcsv1.BcsConfigSpec{
			Name:      "test-bcs-service",
			Namespace: "test-namespace",
			Nmos: bcsv1.Nmos{
				NmosInputFile: nmos.Config{
					HttpPort: 8080,
				},
				NmosApiNodePort: 30080,
			},
		}

		service := CreateBcsService(bcsConfig)
		assert.NotNil(t, service)

		assert.Equal(t, "test-bcs-service", service.ObjectMeta.Name)
		assert.Equal(t, "test-namespace", service.ObjectMeta.Namespace)
		assert.Equal(t, corev1.ServiceTypeNodePort, service.Spec.Type)
		assert.Equal(t, map[string]string{"app": "test-bcs-service"}, service.Spec.Selector)
		assert.Len(t, service.Spec.Ports, 1)
		assert.Equal(t, "nmos-node-api", service.Spec.Ports[0].Name)
		assert.Equal(t, corev1.ProtocolTCP, service.Spec.Ports[0].Protocol)
		assert.Equal(t, int32(8080), service.Spec.Ports[0].Port)
		assert.Equal(t, intstr.FromInt(8080), service.Spec.Ports[0].TargetPort)
		assert.Equal(t, int32(30080), service.Spec.Ports[0].NodePort)
	})

	t.Run("InvalidBcsConfigSpec", func(t *testing.T) {
		bcsConfig := &bcsv1.BcsConfigSpec{
			Name:      "",
			Namespace: "",
			Nmos: bcsv1.Nmos{
				NmosInputFile: nmos.Config{
					HttpPort: 0,
				},
				NmosApiNodePort: 0,
			},
		}

		service := CreateBcsService(bcsConfig)
		assert.NotNil(t, service)

		assert.Equal(t, "", service.ObjectMeta.Name)
		assert.Equal(t, "", service.ObjectMeta.Namespace)

		assert.Equal(t, corev1.ServiceTypeNodePort, service.Spec.Type)

		assert.Equal(t, map[string]string{"app": ""}, service.Spec.Selector)

		assert.Len(t, service.Spec.Ports, 1)
		assert.Equal(t, "nmos-node-api", service.Spec.Ports[0].Name)
		assert.Equal(t, corev1.ProtocolTCP, service.Spec.Ports[0].Protocol)
		assert.Equal(t, int32(0), service.Spec.Ports[0].Port)
		assert.Equal(t, intstr.FromInt(0), service.Spec.Ports[0].TargetPort)
		assert.Equal(t, int32(0), service.Spec.Ports[0].NodePort)
	})
}
func TestCreatePersistentVolumeClaim(t *testing.T) {
	t.Run("ValidConfigMap", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			Data: map[string]string{
				"config.yaml": `
k8s: true
definition:
  mediaProxy:
    pvcAssignedName: "test-pvc"
    pvcStorage: "5Gi"
    pvStorageClass: "standard"
`,
			},
		}

		pvc := CreatePersistentVolumeClaim(cm)
		assert.NotNil(t, pvc)

		assert.Equal(t, "test-pvc", pvc.ObjectMeta.Name)
		assert.Equal(t, "mcm", pvc.ObjectMeta.Namespace)

		assert.Equal(t, []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}, pvc.Spec.AccessModes)
		assert.Equal(t, resource.MustParse("5Gi"), pvc.Spec.Resources.Requests[corev1.ResourceStorage])
		assert.Equal(t, "standard", *pvc.Spec.StorageClassName)
		assert.Equal(t, "mtl-pv", pvc.Spec.VolumeName)
	})
}
func TestCreateConfigMap(t *testing.T) {
	t.Run("ValidBcsConfigSpec", func(t *testing.T) {
		bcsConfig := &bcsv1.BcsConfigSpec{
			Name:      "test-config-map",
			Namespace: "test-namespace",
			Nmos: bcsv1.Nmos{
				NmosInputFile: nmos.Config{
					FfmpegGrpcServerPort:    "50051",
					FfmpegGrpcServerAddress: "localhost",
				},
			},
			App: bcsv1.App{
				GrpcPort: 50051,
			},
		}

		configMap := CreateConfigMap(bcsConfig)
		assert.NotNil(t, configMap)

		assert.Equal(t, "test-config-map-config", configMap.ObjectMeta.Name)
		assert.Equal(t, "test-namespace", configMap.ObjectMeta.Namespace)

		assert.Contains(t, configMap.Data, "config.json")
		var configData nmos.Config
		err := json.Unmarshal([]byte(configMap.Data["config.json"]), &configData)
		assert.NoError(t, err)
		assert.Equal(t, "localhost", configData.FfmpegGrpcServerAddress)
		assert.Equal(t, "50051", configData.FfmpegGrpcServerPort)
	})
}
func TestAssignNodeAntiAffinityFromConfig(t *testing.T) {
	t.Run("ValidDoNotScheduleOnNode", func(t *testing.T) {
		doNotScheduleOnNode := []string{"key1=value1", "key2=value2"}
		antiAffinity := &corev1.PodAntiAffinity{}

		AssignNodeAntiAffinityFromConfig(antiAffinity, doNotScheduleOnNode)

		assert.NotNil(t, antiAffinity)
		assert.Len(t, antiAffinity.RequiredDuringSchedulingIgnoredDuringExecution, 2)

		assert.Equal(t, "key1", antiAffinity.RequiredDuringSchedulingIgnoredDuringExecution[0].LabelSelector.MatchExpressions[0].Key)
		assert.Equal(t, metav1.LabelSelectorOpNotIn, antiAffinity.RequiredDuringSchedulingIgnoredDuringExecution[0].LabelSelector.MatchExpressions[0].Operator)
		assert.Contains(t, antiAffinity.RequiredDuringSchedulingIgnoredDuringExecution[0].LabelSelector.MatchExpressions[0].Values, "value1")

		assert.Equal(t, "key2", antiAffinity.RequiredDuringSchedulingIgnoredDuringExecution[1].LabelSelector.MatchExpressions[0].Key)
		assert.Equal(t, metav1.LabelSelectorOpNotIn, antiAffinity.RequiredDuringSchedulingIgnoredDuringExecution[1].LabelSelector.MatchExpressions[0].Operator)
		assert.Contains(t, antiAffinity.RequiredDuringSchedulingIgnoredDuringExecution[1].LabelSelector.MatchExpressions[0].Values, "value2")
	})

	t.Run("EmptyDoNotScheduleOnNode", func(t *testing.T) {
		doNotScheduleOnNode := []string{}
		antiAffinity := &corev1.PodAntiAffinity{}

		AssignNodeAntiAffinityFromConfig(antiAffinity, doNotScheduleOnNode)

		assert.NotNil(t, antiAffinity)
		assert.Nil(t, antiAffinity.RequiredDuringSchedulingIgnoredDuringExecution)
	})

	t.Run("InvalidDoNotScheduleOnNode", func(t *testing.T) {
		doNotScheduleOnNode := []string{"invalid-entry"}
		antiAffinity := &corev1.PodAntiAffinity{}

		AssignNodeAntiAffinityFromConfig(antiAffinity, doNotScheduleOnNode)

		assert.NotNil(t, antiAffinity)
		assert.Nil(t, antiAffinity.RequiredDuringSchedulingIgnoredDuringExecution)
	})
}
func TestAssignNodeAffinityFromConfig(t *testing.T) {
	t.Run("ValidScheduleOnNode", func(t *testing.T) {
		scheduleOnNode := []string{"key1=value1", "key2=value2"}
		affinity := &corev1.Affinity{}

		AssignNodeAffinityFromConfig(affinity, scheduleOnNode)

		assert.NotNil(t, affinity)
		assert.NotNil(t, affinity.NodeAffinity)
		assert.NotNil(t, affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution)
		assert.Len(t, affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms, 2)

		assert.Equal(t, "key1", affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Key)
		assert.Equal(t, corev1.NodeSelectorOpIn, affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Operator)
		assert.Contains(t, affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Values, "value1")

		assert.Equal(t, "key2", affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[1].MatchExpressions[0].Key)
		assert.Equal(t, corev1.NodeSelectorOpIn, affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[1].MatchExpressions[0].Operator)
		assert.Contains(t, affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[1].MatchExpressions[0].Values, "value2")
	})

	t.Run("EmptyScheduleOnNode", func(t *testing.T) {
		scheduleOnNode := []string{}
		affinity := &corev1.Affinity{}

		AssignNodeAffinityFromConfig(affinity, scheduleOnNode)

		assert.NotNil(t, affinity)
		assert.Nil(t, affinity.NodeAffinity)
	})

	t.Run("InvalidScheduleOnNode", func(t *testing.T) {
		scheduleOnNode := []string{"invalid-entry"}
		affinity := &corev1.Affinity{}

		AssignNodeAffinityFromConfig(affinity, scheduleOnNode)

		assert.NotNil(t, affinity)
		assert.Nil(t, affinity.NodeAffinity)
	})
}
func TestConvertEnvVars(t *testing.T) {
	t.Run("ValidEnvVars", func(t *testing.T) {
		envVars := []bcsv1.EnvVar{
			{Name: "ENV_VAR_1", Value: "value1"},
			{Name: "ENV_VAR_2", Value: "value2"},
		}

		coreEnvVars := convertEnvVars(envVars)

		assert.Len(t, coreEnvVars, 2)
		assert.Equal(t, "ENV_VAR_1", coreEnvVars[0].Name)
		assert.Equal(t, "value1", coreEnvVars[0].Value)
		assert.Equal(t, "ENV_VAR_2", coreEnvVars[1].Name)
		assert.Equal(t, "value2", coreEnvVars[1].Value)
	})

	t.Run("EmptyEnvVars", func(t *testing.T) {
		envVars := []bcsv1.EnvVar{}

		coreEnvVars := convertEnvVars(envVars)

		assert.Empty(t, coreEnvVars)
	})

	t.Run("NilEnvVars", func(t *testing.T) {
		var envVars []bcsv1.EnvVar

		coreEnvVars := convertEnvVars(envVars)

		assert.Empty(t, coreEnvVars)
	})
}
func TestCreateNamespace(t *testing.T) {
	t.Run("ValidNamespaceName", func(t *testing.T) {
		namespaceName := "test-namespace"
		namespace := CreateNamespace(namespaceName)

		assert.NotNil(t, namespace)
		assert.Equal(t, namespaceName, namespace.ObjectMeta.Name)
	})

	t.Run("EmptyNamespaceName", func(t *testing.T) {
		namespaceName := ""
		namespace := CreateNamespace(namespaceName)

		assert.NotNil(t, namespace)
		assert.Equal(t, namespaceName, namespace.ObjectMeta.Name)
	})
}
func TestInt32Ptr(t *testing.T) {
	t.Run("ValidInt32", func(t *testing.T) {
		value := int32(42)
		ptr := int32Ptr(value)
		assert.NotNil(t, ptr)
		assert.Equal(t, value, *ptr)
	})

	t.Run("ZeroValue", func(t *testing.T) {
		value := int32(0)
		ptr := int32Ptr(value)
		assert.NotNil(t, ptr)
		assert.Equal(t, value, *ptr)
	})

	t.Run("NegativeValue", func(t *testing.T) {
		value := int32(-42)
		ptr := int32Ptr(value)
		assert.NotNil(t, ptr)
		assert.Equal(t, value, *ptr)
	})
}
func TestInt64Ptr(t *testing.T) {
	t.Run("ValidInt64", func(t *testing.T) {
		value := int64(42)
		ptr := int64Ptr(value)
		assert.NotNil(t, ptr)
		assert.Equal(t, value, *ptr)
	})

	t.Run("ZeroValue", func(t *testing.T) {
		value := int64(0)
		ptr := int64Ptr(value)
		assert.NotNil(t, ptr)
		assert.Equal(t, value, *ptr)
	})

	t.Run("NegativeValue", func(t *testing.T) {
		value := int64(-42)
		ptr := int64Ptr(value)
		assert.NotNil(t, ptr)
		assert.Equal(t, value, *ptr)
	})
}
func TestBcsConstructContainerConfig(t *testing.T) {
	log := logr.Discard()

	t.Run("BcsPipelineNmosClientConfig", func(t *testing.T) {
		containerInfo := &general.Containers{Type: general.BcsPipelineNmosClient, Id: 0, ContainerName: "nmosclient", Image: "nmos-client-container"}
		config := &parser.Configuration{
			WorkloadToBeRun: []workloads.WorkloadConfig{
				{
					NmosClient: workloads.NmosClientConfig{
						ImageAndTag:             "nmosclient:latest",
						NmosConfigFileName:      "config.json",
						NmosConfigPath:          "/host/config",
						FfmpegConnectionAddress: "192.168.1.103",
						FfmpegConnectionPort:    "50052",
						NmosPort:                50053,
						EnvironmentVariables:    []string{"ENV_VAR=value"},
						Network: workloads.NetworkConfig{
							Enable: true,
							Name:   "test-network",
							IP:     "192.168.1.103",
						},
					},
				},
			},
		}

		err := os.MkdirAll("/host/config", os.ModePerm)
		assert.NoError(t, err)

		tempFile, err := os.CreateTemp("/host/config", "config.json")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		_, err = tempFile.Write([]byte(`{"ffmpeg_grpc_server_address": "old-address", "ffmpeg_grpc_server_port": "50051"}`))
		assert.NoError(t, err)
		tempFile.Close()

		_, err = os.Stat(tempFile.Name())
		assert.NoError(t, err, "config.json file should exist")

		containerConfig, hostConfig, networkConfig := ConstructContainerConfig(containerInfo, config, log)

		fmt.Printf("Container Config: %+v\n", containerConfig)
		fmt.Printf("Host Config: %+v\n", hostConfig)
		fmt.Printf("Network Config: %+v\n", networkConfig)

		assert.NotNil(t, containerConfig)
		assert.NotNil(t, hostConfig)
		assert.NotNil(t, networkConfig)

		assert.Equal(t, "nmosclient:latest", containerConfig.Image)
		assert.Contains(t, containerConfig.Cmd, "config/config.json")
		assert.Equal(t, "test-network", string(hostConfig.NetworkMode))
		assert.Equal(t, "192.168.1.103", networkConfig.EndpointsConfig["test-network"].IPAMConfig.IPv4Address)
	})

	t.Run("BcsPipelineNmosClientConfigHostNetwork", func(t *testing.T) {
		containerInfo := &general.Containers{Type: general.BcsPipelineNmosClient, Id: 0, ContainerName: "nmosclient", Image: "nmos-client-container"}
		config := &parser.Configuration{
			WorkloadToBeRun: []workloads.WorkloadConfig{
				{
					NmosClient: workloads.NmosClientConfig{
						ImageAndTag:             "nmosclient:latest",
						NmosConfigFileName:      "config.json",
						NmosConfigPath:          "/host/config",
						FfmpegConnectionAddress: "192.168.1.103",
						FfmpegConnectionPort:    "50052",
						NmosPort:                50053,
						EnvironmentVariables:    []string{"ENV_VAR=value"},
						Network: workloads.NetworkConfig{
							Enable: false,
							IP:     "192.168.1.103",
						},
					},
				},
			},
		}

		err := os.MkdirAll("/host/config", os.ModePerm)
		assert.NoError(t, err)

		tempFile, err := os.CreateTemp("/host/config", "config.json")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		_, err = tempFile.Write([]byte(`{"ffmpeg_grpc_server_address": "old-address", "ffmpeg_grpc_server_port": "50051"}`))
		assert.NoError(t, err)
		tempFile.Close()

		_, err = os.Stat(tempFile.Name())
		assert.NoError(t, err, "config.json file should exist")

		containerConfig, hostConfig, networkConfig := ConstructContainerConfig(containerInfo, config, log)

		assert.NotNil(t, containerConfig)
		assert.NotNil(t, hostConfig)
		assert.NotNil(t, networkConfig)

		assert.Equal(t, "nmosclient:latest", containerConfig.Image)
		assert.Contains(t, containerConfig.Cmd, "config/config.json")
		assert.Equal(t, "host", string(hostConfig.NetworkMode)) // Host network mode should be empty
		assert.Empty(t, networkConfig.EndpointsConfig)          // No network config for host network
	})
}
func TestCreateMeshAgentDeployment(t *testing.T) {
	t.Run("ValidConfigMap", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			Data: map[string]string{
				"config.yaml": `
k8s: true
definition:
  meshAgent:
    image: "mesh-agent:latest"
    restPort: 8080
    grpcPort: 9090
    resources:
      requests:
        cpu: "500m"
        memory: "256Mi"
      limits:
        cpu: "1000m"
        memory: "512Mi"
    scheduleOnNode:
      - "key1=value1"
    doNotScheduleOnNode:
      - "key2=value2"
`,
			},
		}

		deployment := CreateMeshAgentDeployment(cm)
		assert.NotNil(t, deployment)

		assert.Equal(t, "mesh-agent-deployment", deployment.ObjectMeta.Name)
		assert.Equal(t, "mcm", deployment.ObjectMeta.Namespace)

		assert.Len(t, deployment.Spec.Template.Spec.Containers, 1)
		container := deployment.Spec.Template.Spec.Containers[0]
		assert.Equal(t, "mesh-agent", container.Name)
		assert.Equal(t, "mesh-agent:latest", container.Image)
		assert.Equal(t, corev1.PullIfNotPresent, container.ImagePullPolicy)
		assert.Contains(t, container.Command, "mesh-agent")
		assert.Contains(t, container.Command, "-c")
		assert.Contains(t, container.Command, "8080")
		assert.Contains(t, container.Command, "-p")
		assert.Contains(t, container.Command, "9090")

		resources := container.Resources
		assert.Equal(t, resource.MustParse("500m"), resources.Requests[corev1.ResourceCPU])
		assert.Equal(t, resource.MustParse("256Mi"), resources.Requests[corev1.ResourceMemory])
		assert.Equal(t, resource.MustParse("1000m"), resources.Limits[corev1.ResourceCPU])
		assert.Equal(t, resource.MustParse("512Mi"), resources.Limits[corev1.ResourceMemory])

		assert.Len(t, container.Ports, 2)
		assert.Equal(t, int32(8080), container.Ports[0].ContainerPort)
		assert.Equal(t, int32(9090), container.Ports[1].ContainerPort)

		assert.NotNil(t, container.SecurityContext)
		assert.True(t, *container.SecurityContext.Privileged)

	})

	t.Run("InvalidConfigMap", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			Data: map[string]string{
				"config.yaml": `
k8s: true
definition:
  meshAgent:
    image: "mesh-agent:latest"
    restPort: "invalid-port"
    grpcPort: 9090
    resources:
      requests:
        cpu: "500m"
        memory: "256Mi"
      limits:
        cpu: "1000m"
        memory: "512Mi"
`,
			},
		}

		deployment := CreateMeshAgentDeployment(cm)
		assert.Nil(t, deployment)
	})
}
func TestCreatePersistentVolume(t *testing.T) {
	t.Run("ValidConfigMap", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			Data: map[string]string{
				"config.yaml": `
k8s: true
definition:
  mediaProxy:
    pvHostPath: "/mnt/hostpath"
    pvStorageClass: "standard"
    pvStorage: "10Gi"
`,
			},
		}

		pv := CreatePersistentVolume(cm)
		assert.NotNil(t, pv)

		assert.Equal(t, "mtl-pv", pv.ObjectMeta.Name)

		assert.Equal(t, resource.MustParse("10Gi"), pv.Spec.Capacity[corev1.ResourceStorage])
		assert.Equal(t, corev1.PersistentVolumeFilesystem, *pv.Spec.VolumeMode)
		assert.Equal(t, []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}, pv.Spec.AccessModes)
		assert.Equal(t, corev1.PersistentVolumeReclaimRetain, pv.Spec.PersistentVolumeReclaimPolicy)
		assert.Equal(t, "standard", pv.Spec.StorageClassName)
		assert.Equal(t, "/mnt/hostpath", pv.Spec.PersistentVolumeSource.HostPath.Path)
	})

}
func TestCreateBcsDeployment(t *testing.T) {
	t.Run("ValidBcsConfigSpec", func(t *testing.T) {
		bcsConfig := &bcsv1.BcsConfigSpec{
			Name:      "test-bcs-deployment",
			Namespace: "test-namespace",
			Nmos: bcsv1.Nmos{
				Image: "nmos-node:latest",
				Args:  []string{"--config", "/home/config/config.json"},
				EnvironmentVariables: []bcsv1.EnvVar{
					{Name: "ENV_VAR_1", Value: "value1"},
					{Name: "ENV_VAR_2", Value: "value2"},
				},
			},
			App: bcsv1.App{
				Image:    "broadcast-suite:latest",
				GrpcPort: 50051,
				Volumes: map[string]string{
					"videos":      "/host/videos",
					"dri":         "/host/dri",
					"kahawaiLock": "/host/kahawai.lock",
					"devNull":     "/dev/null",
					"imtl":        "/var/run/imtl",
					"shm":         "/dev/shm",
					"vfio":        "/dev/vfio",
					"dri-dev":     "/dev/dri",
				},
			},
			ScheduleOnNode:      []string{"key1=value1"},
			DoNotScheduleOnNode: []string{"key2=value2"},
		}

		deployment := CreateBcsDeployment(bcsConfig)
		assert.NotNil(t, deployment)

		assert.Equal(t, "test-bcs-deployment", deployment.ObjectMeta.Name)
		assert.Equal(t, "test-namespace", deployment.ObjectMeta.Namespace)

		nmosContainer := deployment.Spec.Template.Spec.Containers[0]
		assert.Equal(t, "tiber-broadcast-suite-nmos-node", nmosContainer.Name)
		assert.Equal(t, "nmos-node:latest", nmosContainer.Image)
		assert.Equal(t, []string{"--config", "/home/config/config.json"}, nmosContainer.Args)

		assert.Len(t, nmosContainer.Env, 2)
		assert.Equal(t, "ENV_VAR_1", nmosContainer.Env[0].Name)
		assert.Equal(t, "value1", nmosContainer.Env[0].Value)

		appContainer := deployment.Spec.Template.Spec.Containers[1]
		assert.Equal(t, "tiber-broadcast-suite", appContainer.Name)
		assert.Equal(t, "broadcast-suite:latest", appContainer.Image)
		assert.Equal(t, []string{"localhost", "50051"}, appContainer.Args)

		assert.Equal(t, "/host/videos", deployment.Spec.Template.Spec.Volumes[0].HostPath.Path)
		assert.Equal(t, "/host/dri", deployment.Spec.Template.Spec.Volumes[1].HostPath.Path)

		assert.NotNil(t, deployment.Spec.Template.Spec.Affinity)
		assert.NotNil(t, deployment.Spec.Template.Spec.Affinity.NodeAffinity)
		assert.NotNil(t, deployment.Spec.Template.Spec.Affinity.PodAntiAffinity)
	})

	t.Run("EmptyBcsConfigSpec", func(t *testing.T) {
		bcsConfig := &bcsv1.BcsConfigSpec{}
		deployment := CreateBcsDeployment(bcsConfig)
		assert.NotNil(t, deployment)
		assert.Equal(t, "", deployment.ObjectMeta.Name)
	})
}
func TestCreateDaemonSet(t *testing.T) {
	t.Run("ValidConfigMap", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			Data: map[string]string{
				"config.yaml": `
k8s: true
definition:
  mediaProxy:
    image: "media-proxy:latest"
    command: ["media-proxy"]
    args: ["--grpc-port", "50051", "--sdk-port", "50052"]
    grpcPort: 50051
    sdkPort: 50052
    resources:
      requests:
        cpu: "2"
        memory: "8Gi"
        hugepages-1Gi: "1Gi"
        hugepages-2Mi: "2Gi"
      limits:
        cpu: "4"
        memory: "16Gi"
        hugepages-1Gi: "2Gi"
        hugepages-2Mi: "4Gi"
    volumes:
      memif: "/run/memif"
      vfio: "/dev/vfio"
      cache-size: "2Gi"
    pvHostPath: "/mnt/hostpath"
    pvStorageClass: "standard"
    pvStorage: "10Gi"
    pvcStorage: "5Gi"
    pvcAssignedName: "media-proxy-pvc"
    scheduleOnNode:
      - "key1=value1"
    doNotScheduleOnNode:
      - "key2=value2"
`,
			},
		}

		daemonSet := CreateDaemonSet(cm)
		assert.NotNil(t, daemonSet)

		assert.Equal(t, "media-proxy", daemonSet.ObjectMeta.Name)
		assert.Equal(t, "mcm", daemonSet.ObjectMeta.Namespace)

		container := daemonSet.Spec.Template.Spec.Containers[0]
		assert.Equal(t, "media-proxy", container.Name)
		assert.Equal(t, "media-proxy:latest", container.Image)
		assert.Equal(t, []string{"media-proxy"}, container.Command)
		assert.Equal(t, []string{"--grpc-port", "50051", "--sdk-port", "50052"}, container.Args)

		assert.Equal(t, resource.MustParse("2"), container.Resources.Requests[corev1.ResourceCPU])
		assert.Equal(t, resource.MustParse("8Gi"), container.Resources.Requests[corev1.ResourceMemory])
		assert.Equal(t, resource.MustParse("1Gi"), container.Resources.Requests[corev1.ResourceHugePagesPrefix+"1Gi"])
		assert.Equal(t, resource.MustParse("2Gi"), container.Resources.Requests[corev1.ResourceHugePagesPrefix+"2Mi"])
		assert.Equal(t, resource.MustParse("4"), container.Resources.Limits[corev1.ResourceCPU])
		assert.Equal(t, resource.MustParse("16Gi"), container.Resources.Limits[corev1.ResourceMemory])
		assert.Equal(t, resource.MustParse("2Gi"), container.Resources.Limits[corev1.ResourceHugePagesPrefix+"1Gi"])
		assert.Equal(t, resource.MustParse("4Gi"), container.Resources.Limits[corev1.ResourceHugePagesPrefix+"2Mi"])

		assert.Len(t, container.Ports, 2)
		assert.Equal(t, int32(50051), container.Ports[0].ContainerPort)
		assert.Equal(t, int32(50051), container.Ports[0].HostPort)
		assert.Equal(t, corev1.ProtocolTCP, container.Ports[0].Protocol)
		assert.Equal(t, "grpc-port", container.Ports[0].Name)
		assert.Equal(t, int32(50052), container.Ports[1].ContainerPort)
		assert.Equal(t, int32(50052), container.Ports[1].HostPort)
		assert.Equal(t, corev1.ProtocolTCP, container.Ports[1].Protocol)
		assert.Equal(t, "sdk-port", container.Ports[1].Name)

		assert.Len(t, container.VolumeMounts, 6)
		assert.Equal(t, "/run/mcm", container.VolumeMounts[0].MountPath)
		assert.Equal(t, "/dev/vfio", container.VolumeMounts[1].MountPath)
		assert.Equal(t, "/hugepages-2Mi", container.VolumeMounts[2].MountPath)
		assert.Equal(t, "/hugepages-1Gi", container.VolumeMounts[3].MountPath)
		assert.Equal(t, "/dev/shm", container.VolumeMounts[4].MountPath)
		assert.Equal(t, "/var/run/imtl", container.VolumeMounts[5].MountPath)

		assert.Len(t, daemonSet.Spec.Template.Spec.Volumes, 6)
		assert.Equal(t, "/run/memif", daemonSet.Spec.Template.Spec.Volumes[0].VolumeSource.HostPath.Path)
		assert.Equal(t, "/dev/vfio", daemonSet.Spec.Template.Spec.Volumes[1].VolumeSource.HostPath.Path)
		assert.Equal(t, v1.StorageMedium("HugePages-2Mi"), daemonSet.Spec.Template.Spec.Volumes[2].VolumeSource.EmptyDir.Medium)
		assert.Equal(t, v1.StorageMedium("HugePages-1Gi"), daemonSet.Spec.Template.Spec.Volumes[3].VolumeSource.EmptyDir.Medium)
		assert.Equal(t, corev1.StorageMediumMemory, daemonSet.Spec.Template.Spec.Volumes[4].VolumeSource.EmptyDir.Medium)
		assert.Equal(t, "media-proxy-pvc", daemonSet.Spec.Template.Spec.Volumes[5].VolumeSource.PersistentVolumeClaim.ClaimName)

		affinity := daemonSet.Spec.Template.Spec.Affinity
		assert.NotNil(t, affinity)
		assert.NotNil(t, affinity.NodeAffinity)
		assert.NotNil(t, affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution)
		assert.Len(t, affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms, 1)
		assert.Equal(t, "key1", affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Key)
		assert.Contains(t, affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Values, "value1")

		podAntiAffinity := affinity.PodAntiAffinity
		assert.NotNil(t, podAntiAffinity)
		assert.Len(t, podAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution, 1)
		assert.Equal(t, "key2", podAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution[0].LabelSelector.MatchExpressions[0].Key)
		assert.Contains(t, podAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution[0].LabelSelector.MatchExpressions[0].Values, "value2")
	})

	t.Run("InvalidConfigMap", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			Data: map[string]string{
				"config.yaml": `
k8s: true
definition:
  mediaProxy:
    image: "media-proxy:latest"
    grpcPort: "invalid-port"
`,
			},
		}

		daemonSet := CreateDaemonSet(cm)
		assert.Nil(t, daemonSet)
	})
}
