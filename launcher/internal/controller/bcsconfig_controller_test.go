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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	bcsv1 "bcs.pod.launcher.intel/api/v1"
)

var _ = Describe("BcsConfig Controller", func() {
    Context("When reconciling a resource", func() {


        const resourceName = "custom-resource-bcsconfig"

        ctx := context.Background()

        BeforeEach(func() {
            By("creating the namespace 'bcs'")
            ns := &corev1.Namespace{
            ObjectMeta: metav1.ObjectMeta{
                Name: "bcs",
            },
            }
            err := k8sClient.Create(ctx, ns)
            if err != nil && !errors.IsAlreadyExists(err) {
            Expect(err).NotTo(HaveOccurred())
            }
        })

        typeNamespacedName := types.NamespacedName{
            Name:      resourceName,
            Namespace: "default",
        }
        bcsconfig := &bcsv1.BcsConfig{}

        cmNamespacedName := types.NamespacedName{
            Name:      "k8s-bcs-config",
            Namespace: "bcs",
        }
        cmMcm := &corev1.ConfigMap{}

        BeforeEach(func() {
            By("creating the config map for the Cluster")
            
            err := k8sClient.Get(ctx, cmNamespacedName, cmMcm)
            if err != nil && errors.IsNotFound(err) {
                resource := &corev1.ConfigMap{
                    ObjectMeta: metav1.ObjectMeta{
                        Name:      cmNamespacedName.Name,
                        Namespace: cmNamespacedName.Namespace,
                    },
                    Data: map[string]string{
                        "config.yaml": `
                        k8s: true
                        definition:
                        meshAgent:
                            image: "mesh-agent:latest"
                            restPort: 8100
                            grpcPort: 50051
                        mediaProxy:
                            image: mcm/media-proxy:latest
                            command: ["media-proxy"]
                            args: ["-d", "kernel:eth0", "-i", "$(POD_IP)"]
                            grpcPort: 8001
                            sdkPort: 8002
                            volumes:
                            memif: /tmp/mcm/memif
                            vfio: /dev/vfio
                            pvHostPath: /var/run/imtl
                            pvStorageClass: manual
                            pvStorage: 1Gi
                            pvcStorage: 1Gi
                        `,
                    },
                }
                Expect(k8sClient.Create(ctx, resource)).To(Succeed())
            }
           
        })

        BeforeEach(func() {
            By("creating the crd for the Cluster")
            
            err := k8sClient.Get(ctx, typeNamespacedName, bcsconfig)
            if err != nil && errors.IsNotFound(err) {
                resource := &bcsv1.BcsConfig{
                    ObjectMeta: metav1.ObjectMeta{
                        Name:      typeNamespacedName.Name,
                        Namespace: typeNamespacedName.Namespace,
                    },
                    Spec: bcsv1.BcsConfigSpec{
                        Name:      typeNamespacedName.Name,
                        Namespace: typeNamespacedName.Namespace,
                        Nmos: bcsv1.Nmos{
                            NmosInputFile: bcsv1.NmosInputFile{
                                LoggingLevel:            10,
                                HttpPort:                8080,
                                Label:                   "test-label",
                                DeviceTags:              bcsv1.DeviceTags{Pipeline: []string{"rx"}},
                                Function:                "rx",
                                GpuHwAcceleration:       "none",
                                Domain:                  "local",
                                FfmpegGrpcServerAddress: "192.168.1.100",
                                FfmpegGrpcServerPort:    "50051",
                                SenderPayloadType:       96,
                                Sender: []bcsv1.Sender{
                                    {
                                        StreamPayload: bcsv1.StreamPayload{
                                            Video: bcsv1.Video{
                                                FrameWidth:  1920,
                                                FrameHeight: 1080,
                                                FrameRate: bcsv1.FrameRate{
                                                    Numerator:   30,
                                                    Denominator: 1,
                                                },
                                                PixelFormat: "yuv420p",
                                                VideoType:   "rawvideo",
                                            },
                                            Audio: bcsv1.Audio{
                                                Channels:   2,
                                                SampleRate: 48000,
                                                Format:     "pcm_s16le",
                                                PacketTime: "1ms",
                                            },
                                        },
                                        StreamType: bcsv1.StreamType{
                                            File: &bcsv1.File{
                                                Path:     "/videos",
                                                Filename: "test_video.yuv",
                                            },
                                        },
                                    },
                                },
                                Receiver: []bcsv1.Receiver{
                                    {
                                        StreamPayload: bcsv1.StreamPayload{
                                            Video: bcsv1.Video{
                                                FrameWidth:  1920,
                                                FrameHeight: 1080,
                                                FrameRate: bcsv1.FrameRate{
                                                    Numerator:   30,
                                                    Denominator: 1,
                                                },
                                                PixelFormat: "yuv420p",
                                                VideoType:   "rawvideo",
                                            },
                                            Audio: bcsv1.Audio{
                                                Channels:   2,
                                                SampleRate: 48000,
                                                Format:     "pcm_s16le",
                                                PacketTime: "1ms",
                                            },
                                        },
                                        StreamType: bcsv1.StreamType{
                                            St2110: &bcsv1.St2110{
                                                Transport:   "st2110-20",
                                                Payload_type : 112,
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
            // TODO(user): Cleanup logic after each test, like removing the resource instance.
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
            // TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
            // Example: If you expect a certain status condition after reconciliation, verify it here.
        })
        
    })

})
