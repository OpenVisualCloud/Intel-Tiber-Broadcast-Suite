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

	// Initial JSON data
	jsonData := `{
        "ffmpeg_grpc_server_address": "old-address",
        "ffmpeg_grpc_server_port": "50051"
    }`
	_, err = tempFile.Write([]byte(jsonData))
	assert.NoError(t, err)
	tempFile.Close()

	// Update the JSON file
	err = updateNmosJsonFile(tempFile.Name(), "new-address", "50052")
	assert.NoError(t, err)

	// Read and verify the updated JSON data
	updatedData, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err)

	var config map[string]interface{}
	err = json.Unmarshal(updatedData, &config)
	assert.NoError(t, err)
	assert.Equal(t, "new-address", config["ffmpeg_grpc_server_address"])
	assert.Equal(t, "50052", fmt.Sprintf("%v", config["ffmpeg_grpc_server_port"]))
}

func TestUpdateNmosJsonFile_FileNotFound(t *testing.T) {
	err := updateNmosJsonFile("non-existent-file.json", "new-address", "new-port")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestUpdateNmosJsonFile_InvalidJson(t *testing.T) {
	tempFile, err := os.CreateTemp("", "nmos_test_invalid_*.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Write invalid JSON data
	_, err = tempFile.Write([]byte(`{invalid-json}`))
	assert.NoError(t, err)
	tempFile.Close()

	// Attempt to update the JSON file
	err = updateNmosJsonFile(tempFile.Name(), "new-address", "new-port")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}
func TestConstructContainerConfig(t *testing.T) {
	log := logr.Discard()

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

		// Verify resource requests and limits
		resources := deployment.Spec.Template.Spec.Containers[0].Resources
		assert.Equal(t, resource.MustParse("500m"), resources.Requests[corev1.ResourceCPU])
		assert.Equal(t, resource.MustParse("256Mi"), resources.Requests[corev1.ResourceMemory])
		assert.Equal(t, resource.MustParse("1000m"), resources.Limits[corev1.ResourceCPU])
		assert.Equal(t, resource.MustParse("512Mi"), resources.Limits[corev1.ResourceMemory])

		// Verify volume mounts
		volumeMounts := deployment.Spec.Template.Spec.Containers[0].VolumeMounts
		assert.Len(t, volumeMounts, 2)
		assert.Equal(t, "/var/run/imtl", volumeMounts[0].MountPath)
		assert.Equal(t, "/sys/fs/bpf", volumeMounts[1].MountPath)

		// Verify node affinity
		affinity := deployment.Spec.Template.Spec.Affinity
		assert.NotNil(t, affinity)
		assert.NotNil(t, affinity.NodeAffinity)
		assert.NotNil(t, affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution)
		assert.Len(t, affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms, 1)
		assert.Equal(t, "key1", affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Key)
		assert.Contains(t, affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Values, "value1")

		// Verify pod anti-affinity
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

		// Verify metadata
		assert.Equal(t, "mesh-agent-service", service.ObjectMeta.Name)
		assert.Equal(t, "mcm", service.ObjectMeta.Namespace)

		// Verify selector
		assert.Equal(t, map[string]string{"app": "mesh-agent"}, service.Spec.Selector)

		// Verify ports
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

		// Verify selector
		assert.Equal(t, map[string]string{"app": serviceName}, service.Spec.Selector)

		// Verify ports
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

		// Verify selector
		assert.Equal(t, map[string]string{"app": serviceName}, service.Spec.Selector)

		// Verify ports
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

		// Verify metadata
		assert.Equal(t, "test-bcs-service", service.ObjectMeta.Name)
		assert.Equal(t, "test-namespace", service.ObjectMeta.Namespace)

		// Verify service type
		assert.Equal(t, corev1.ServiceTypeNodePort, service.Spec.Type)

		// Verify selector
		assert.Equal(t, map[string]string{"app": "test-bcs-service"}, service.Spec.Selector)

		// Verify ports
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

		// Verify metadata
		assert.Equal(t, "", service.ObjectMeta.Name)
		assert.Equal(t, "", service.ObjectMeta.Namespace)

		// Verify service type
		assert.Equal(t, corev1.ServiceTypeNodePort, service.Spec.Type)

		// Verify selector
		assert.Equal(t, map[string]string{"app": ""}, service.Spec.Selector)

		// Verify ports
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

		// Verify metadata
		assert.Equal(t, "test-pvc", pvc.ObjectMeta.Name)
		assert.Equal(t, "mcm", pvc.ObjectMeta.Namespace)

		// Verify spec
		assert.Equal(t, []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}, pvc.Spec.AccessModes)
		assert.Equal(t, resource.MustParse("5Gi"), pvc.Spec.Resources.Requests[corev1.ResourceStorage])
		assert.Equal(t, "standard", *pvc.Spec.StorageClassName)
		assert.Equal(t, "mtl-pv", pvc.Spec.VolumeName)
	})

	t.Run("InvalidConfigMap", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			Data: map[string]string{
				"config.yaml": `
	k8s: true
	definition:
	  mediaProxy:
	    pvcAssignedName: ""
	    pvcStorage: ""
	    pvStorageClass: ""
	`,
			},
		}

		pvc := CreatePersistentVolumeClaim(cm)
		assert.Nil(t, pvc)
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

		// Verify metadata
		assert.Equal(t, "test-config-map-config", configMap.ObjectMeta.Name)
		assert.Equal(t, "test-namespace", configMap.ObjectMeta.Namespace)

		// Verify data
		assert.Contains(t, configMap.Data, "config.json")
		var configData nmos.Config
		err := json.Unmarshal([]byte(configMap.Data["config.json"]), &configData)
		assert.NoError(t, err)
		assert.Equal(t, "localhost", configData.FfmpegGrpcServerAddress)
		assert.Equal(t, "50051", configData.FfmpegGrpcServerPort)
	})

	// t.Run("InvalidBcsConfigSpec", func(t *testing.T) {
	// 	bcsConfig := &bcsv1.BcsConfigSpec{
	// 		Name:      "",
	// 		Namespace: "",
	// 		Nmos: bcsv1.Nmos{
	// 			NmosInputFile: nmos.Config{},
	// 		},
	// 	}

	// 	configMap := CreateConfigMap(bcsConfig)
	// 	assert.NotNil(t, configMap)

	// 	// Verify metadata
	// 	assert.Equal(t, "-config", configMap.ObjectMeta.Name)
	// 	assert.Equal(t, "", configMap.ObjectMeta.Namespace)

	// 	// Verify data
	// 	assert.Contains(t, configMap.Data, "config.json")
	// 	var configData nmos.Config
	// 	err := json.Unmarshal([]byte(configMap.Data["config.json"]), &configData)
	// 	assert.NoError(t, err)
	// 	assert.Empty(t, configData.FfmpegGrpcServerAddress)
	// 	assert.Empty(t, configData.FfmpegGrpcServerPort)
	// })
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
