/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	bcsv1 "bcs.pod.launcher.intel/api/v1"
	"bcs.pod.launcher.intel/resources_library/resources/nmos"
)

var _ = Describe("BcsConfig Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "bcs",
		}
		bcsconfig := &bcsv1.BcsConfig{}

		BeforeEach(func() {
			By("creating the namespace")
			namespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: typeNamespacedName.Namespace,
				},
			}
			err := k8sClient.Create(ctx, namespace)
			if err != nil && !errors.IsAlreadyExists(err) {
				Expect(err).NotTo(HaveOccurred())
			}
			By("creating the ConfigMap for the test")
			configMap := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "k8s-bcs-config",
					Namespace: "bcs",
				},
				Data: map[string]string{
					"config.yaml": `
k8s: true
definition:
  meshAgent:
    image: "mcm/mesh-agent:latest"
    restPort: 8100
    grpcPort: 50051
    requests:
      cpu: "500m"
      memory: "256Mi"
    limits:
      cpu: "1000m"
      memory: "512Mi"
    scheduleOnNode: ["node-role.kubernetes.io/worker=true"]
  mediaProxy:
    image: mcm/media-proxy:latest
    command: ["media-proxy"]
    args: ["-d", "0000:ca:11.0", "-i", $(POD_IP)]
    grpcPort: 8001
    sdkPort: 8002
    requests:
      cpu: "2"
      memory: "8Gi"
      hugepages-1Gi: "1Gi"
      hugepages-2Mi: "2Gi"
    limits:
      cpu: "2"
      memory: "8Gi"
      hugepages-1Gi: "1Gi"
      hugepages-2Mi: "2Gi"
    volumes:
      memif: /tmp/mcm/memif
      vfio: /dev/vfio
      cache-size: 4Gi
    pvHostPath: /var/run/imtl
    pvStorageClass: manual
    pvStorage: 1Gi
    pvcAssignedName: mtl-pvc
    pvcStorage: 1Gi
    scheduleOnNode: ["node-role.kubernetes.io/worker=true"]
  mtlManager:
    image: mtl-manager:latest
    requests:
      cpu: "500m"
      memory: "256Mi"
    limits:
      cpu: "1000m"
      memory: "512Mi"
    volumes:
      imtlHostPath: /var/run/imtl
      bpfPath: /sys/fs/bpf
`,
				},
			}
			err = k8sClient.Create(ctx, configMap)
			if err != nil && !errors.IsAlreadyExists(err) {
				Expect(err).NotTo(HaveOccurred())
			}
			By("creating the custom resource for the Kind BcsConfig")
			err = k8sClient.Get(ctx, typeNamespacedName, bcsconfig)
			if err != nil && errors.IsNotFound(err) {
				resource := &bcsv1.BcsConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name:      typeNamespacedName.Name,
						Namespace: typeNamespacedName.Namespace,
					},
					Spec: []bcsv1.BcsConfigSpec{
						{
							Name:      "tiber-broadcast-suite",
							Namespace: "bcs",
							App: bcsv1.App{
								Image:    "video_production_image:latest",
								GrpcPort: 50051,
								EnvironmentVariables: []bcsv1.EnvVar{
									{Name: "http_proxy", Value: ""},
									{Name: "https_proxy", Value: ""},
								},
								Volumes: map[string]string{
									"videos":      "/root/demo",
									"dri":         "/usr/lib/x86_64-linux-gnu/dri",
									"kahawaiLock": "/tmp/kahawai_lcore.lock",
									"devNull":     "/dev/null",
									"imtl":        "/var/run/imtl",
									"shm":         "/dev/shm",
									"vfio":        "/dev/vfio",
									"dri-dev":     "/dev/dri",
								},
							},
							Nmos: bcsv1.Nmos{
								Image: "tiber-broadcast-suite-nmos-node:latest",
								Args:  []string{"config/config.json"},
								EnvironmentVariables: []bcsv1.EnvVar{
									{Name: "http_proxy", Value: ""},
									{Name: "https_proxy", Value: ""},
									{Name: "VFIO_PORT_TX", Value: "0000:ca:11.0"},
								},
								NmosApiNodePort: 30084,
								NmosInputFile: nmos.Config{
									LoggingLevel:      10,
									HttpPort:          5004,
									Label:             "intel-broadcast-suite-tx",
									Function:          "tx",
									ActivateSenders:   false,
									StreamLoop:        -1,
									GpuHwAcceleration: "none",
									Domain:            "local",
									SenderPayloadType: 112,
									Sender: []nmos.Sender{
										{
											StreamPayload: nmos.StreamPayload{
												Video: nmos.Video{
													FrameWidth:  1920,
													FrameHeight: 1080,
													FrameRate: nmos.FrameRate{
														Numerator:   60,
														Denominator: 1,
													},
													PixelFormat: "yuv422p10le",
													VideoType:   "rawvideo",
												},
												Audio: nmos.Audio{
													Channels:   2,
													SampleRate: 48000,
													Format:     "pcm_s24be",
													PacketTime: "1ms",
												},
											},
											StreamType: nmos.StreamType{
												St2110: &nmos.St2110{
													Transport:    "st2110-20",
													Payload_type: 112,
													QueuesCount:  0,
												},
											},
										},
									},
									Receiver: []nmos.Receiver{
										{
											StreamPayload: nmos.StreamPayload{
												Video: nmos.Video{
													FrameWidth:  1920,
													FrameHeight: 1080,
													FrameRate: nmos.FrameRate{
														Numerator:   60,
														Denominator: 1,
													},
													PixelFormat: "yuv422p10le",
													VideoType:   "rawvideo",
												},
												Audio: nmos.Audio{
													Channels:   2,
													SampleRate: 48000,
													Format:     "pcm_s24be",
													PacketTime: "1ms",
												},
											},
											StreamType: nmos.StreamType{
												File: &nmos.File{
													Path:     "/videos",
													Filename: "1920x1080p10le_0.yuv",
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}

		})

		AfterEach(func() {
			resource := &bcsv1.BcsConfig{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance BcsConfig")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &BcsConfigReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})

			Expect(err).NotTo(HaveOccurred())

			By("Checking if mtlManager, meshAgent, mediaProxy, and bcsPipeline deployments are created")
			mtlManager := &appsv1.Deployment{}
			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      "mtl-manager",
				Namespace: "mcm",
			}, mtlManager)
			Expect(err).NotTo(HaveOccurred())

			meshAgent := &appsv1.Deployment{}
			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      "mesh-agent-deployment",
				Namespace: "mcm",
			}, meshAgent)
			Expect(err).NotTo(HaveOccurred())

			mediaProxy := &appsv1.DaemonSet{}
			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      "media-proxy",
				Namespace: "mcm",
			}, mediaProxy)
			Expect(err).NotTo(HaveOccurred())

			bcsPipeline := &appsv1.Deployment{}
			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      "tiber-broadcast-suite",
				Namespace: typeNamespacedName.Namespace,
			}, bcsPipeline)
			Expect(err).NotTo(HaveOccurred())

			By("Checking if mesh-agent-service in mcm namespace exists")
			meshAgentService := &corev1.Service{}
			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      "mesh-agent-service",
				Namespace: "mcm",
			}, meshAgentService)
			Expect(err).NotTo(HaveOccurred())

			By("Checking if tiber-broadcast-suite service in bcs namespace exists")
			tiberBroadcastSuiteService := &corev1.Service{}
			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      "tiber-broadcast-suite",
				Namespace: typeNamespacedName.Namespace,
			}, tiberBroadcastSuiteService)
			Expect(err).NotTo(HaveOccurred())

			By("Checking if PersistentVolume mtl-pv exists")
			mtlPv := &corev1.PersistentVolume{}
			err = k8sClient.Get(ctx, types.NamespacedName{
				Name: "mtl-pv",
			}, mtlPv)
			Expect(err).NotTo(HaveOccurred())

			By("Checking if PersistentVolumeClaim mtl-pvc in bcs namespace exists")
			mtlPvc := &corev1.PersistentVolumeClaim{}
			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      "mtl-pvc",
				Namespace: "mcm",
			}, mtlPvc)
			Expect(err).NotTo(HaveOccurred())

			By("Checking if ConfigMap tiber-broadcast-suite-config in bcs namespace exists")
			tiberBroadcastSuiteConfig := &corev1.ConfigMap{}
			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      "tiber-broadcast-suite-config",
				Namespace: typeNamespacedName.Namespace,
			}, tiberBroadcastSuiteConfig)
			Expect(err).NotTo(HaveOccurred())
		})

	})
})
