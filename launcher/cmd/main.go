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

package main

import (
	"crypto/tls"
	"flag"
	"os"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	bcsv1 "bcs.pod.launcher.intel/api/v1"
	"bcs.pod.launcher.intel/internal/controller"
	"bcs.pod.launcher.intel/resources_library/utils"
	//+kubebuilder:scaffold:imports
)

var (
	scheme            = runtime.NewScheme()
	setupLog          = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(bcsv1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	ctx := ctrl.SetupSignalHandler()
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var secureMetrics bool
	var enableHTTP2 bool
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&secureMetrics, "metrics-secure", false,
		"If set the metrics endpoint is served securely")
	flag.BoolVar(&enableHTTP2, "enable-http2", false,
		"If set, HTTP/2 will be enabled for the metrics and webhook servers")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

		 // if the enable-http2 flag is false (the default), http/2 should be disabled
		 // due to its vulnerabilities. More specifically, disabling http/2 will
		 // prevent from being vulnerable to the HTTP/2 Stream Cancellation and
		 // Rapid Reset CVEs. For more information see:
		 // - https://github.com/advisories/GHSA-qppj-fm5r-hxr3
		 // - https://github.com/advisories/GHSA-4374-p667-p6c8
		disableHTTP2 := func(c *tls.Config) {
			setupLog.Info("disabling http/2")
			c.NextProtos = []string{"http/1.1"}
		}

		tlsOpts := []func(*tls.Config){}
		if !enableHTTP2 {
			tlsOpts = append(tlsOpts, disableHTTP2)
		}

		webhookServer := webhook.NewServer(webhook.Options{
			TLSOpts: tlsOpts,
		})

		mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
			Scheme: scheme,
			Metrics: metricsserver.Options{
				BindAddress:   metricsAddr,
				SecureServing: secureMetrics,
				TLSOpts:       tlsOpts,
			},
			WebhookServer:          webhookServer,
			HealthProbeBindAddress: probeAddr,
			LeaderElection:         enableLeaderElection,
			LeaderElectionID:       "2d95eb0a.bcs.intel",
		})
		if err != nil {
			setupLog.Error(err, "unable to start manager")
			os.Exit(1)
		}

		if err = (&controller.BcsConfigReconciler{
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "BcsConfig")
			os.Exit(1)
		}
		//+kubebuilder:scaffold:builder

		if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
			setupLog.Error(err, "unable to set up health check")
			os.Exit(1)
		}
		if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
			setupLog.Error(err, "unable to set up ready check")
			os.Exit(1)
		}

        //MCM  k8s part
		setupLog.Info("Starting MCM pods launcher in Kubernetes")

		k8s_client := mgr.GetClient()

		mcmNamespace := utils.CreateNamespace("mcm")
		mcmAgentDeployment := utils.CreateDeployment("mcm-agent")
		mcmMediaProxyPv := utils.CreatePersistentVolume("media-proxy")
		mcmMediaProxyPvc := utils.CreatePersistentVolumeClaim("media-proxy")
		mcmMediaProxyDs := utils.CreateDaemonSet("media-proxy")

		err = k8s_client.Create(ctx, mcmNamespace)
		if err != nil {
			setupLog.Error(err, "Failed to create mcm namespace")
		}

		err = k8s_client.Create(ctx, mcmAgentDeployment)
		if err != nil {
			setupLog.Error(err, "Failed to create mcmAgent Deployment")
		}

		err = k8s_client.Create(ctx, mcmMediaProxyPv)
		if err != nil {
			setupLog.Error(err, "Failed to create mcm media-proxy persistent volume")
		}

		err = k8s_client.Create(ctx, mcmMediaProxyPvc)
		if err != nil {
			setupLog.Error(err, "Failed to create  mcm media-proxy persistent volume claim")
		}

		err = k8s_client.Create(ctx, mcmMediaProxyDs)
		if err != nil {
			setupLog.Error(err, "Failed to create mcm media-proxy daemonset")
		}

		setupLog.Info("Starting manager that handles pipelines")
		if err := mgr.Start(ctx); err != nil {
			setupLog.Error(err, "problem running manager")
			os.Exit(1)
	}
}