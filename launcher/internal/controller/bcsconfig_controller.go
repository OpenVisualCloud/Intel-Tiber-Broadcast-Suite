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

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	appsv1 "k8s.io/api/apps/v1"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	bcsv1 "bcs.pod.launcher.intel/api/v1"
	"bcs.pod.launcher.intel/resources_library/utils"
)

// BcsConfigReconciler reconciles a BcsConfig object
type BcsConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	logr.Logger
}

// Information about rbac
// groups=bcs.bcs.intel,resources=bcsconfigs,verbs=get;list;watch;create;update;patch;delete
// groups=bcs.bcs.intel,resources=bcsconfigs/status,verbs=get;update;patch
// groups=bcs.bcs.intel,resources=bcsconfigs/finalizers,verbs=update
// groups=apps,resources=daemonsets;deployments,verbs=get;list;watch;create;update
// groups="",resources=services;configmaps;persistentvolumes;persistentvolumeclaims,verbs=get;list;watch;create;update

func (r *BcsConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// MCM silent start up
	createResourceIfNotExists := func(resource client.Object, namespacedName types.NamespacedName) error {
		err := r.Get(ctx, namespacedName, resource)
		if err != nil {
			if errors.IsNotFound(err) {
				// Create the resource if it doesn't exist
				err = r.Create(ctx, resource)
				if err != nil {
					log.Error(err, "Failed to create resource", "resource", resource.GetObjectKind(), "named", namespacedName)
					return err
				}
				log.Info("Resource created successfully", "resource", resource.GetObjectKind(), "name", namespacedName)
			} else {
				// Log and return any other error
				log.Error(err, "Failed to get resource", "resource", resource.GetObjectKind(), "named", namespacedName)
				return err
			}
		} else {
			// Log if the resource already exists
			log.Info("Resource already exists", "resource", resource.GetObjectKind(), "name", namespacedName)
		}
		return nil
	}

	mcmCmInfo := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: "k8s-bcs-config", Namespace: "bcs"}, mcmCmInfo)
	if err != nil {
		log.Error(err, "Failed to get resource", "resource", mcmCmInfo.GetObjectKind(), "named", "k8s-bcs-config")
		return ctrl.Result{}, err
	}

	mcmNamespace := utils.CreateNamespace("mcm")
	mcmAgentDeployment := utils.CreateMeshAgentDeployment(mcmCmInfo)
	mcmAgentService := utils.CreateMeshAgentService(mcmCmInfo)
	mcmMediaProxyPv := utils.CreatePersistentVolume(mcmCmInfo)
	mcmMediaProxyPvc := utils.CreatePersistentVolumeClaim(mcmCmInfo)
	mcmMediaProxyDs := utils.CreateDaemonSet(mcmCmInfo)
	mtlManagerDeployment := utils.CreateMtlManagerDeployment(mcmCmInfo)

	err = createResourceIfNotExists(mcmNamespace, types.NamespacedName{Name: mcmNamespace.Name})
	if err != nil {
		log.Error(err, "Failed to create resource", "resource", mcmNamespace.GetObjectKind(), "named", mcmNamespace.Name)
		return ctrl.Result{}, err
	}
	err = createResourceIfNotExists(mcmAgentDeployment, types.NamespacedName{Name: mcmAgentDeployment.Name, Namespace: "mcm"})
	if err != nil {
		log.Error(err, "Failed to create resource", "resource", mcmAgentDeployment.GetObjectKind(), "named", mcmAgentDeployment.Name)
		return ctrl.Result{}, err
	}
	err = createResourceIfNotExists(mcmAgentService, types.NamespacedName{Name: mcmAgentService.Name, Namespace: "mcm"})
	if err != nil {
		log.Error(err, "Failed to create resource", "resource", mcmAgentService.GetObjectKind(), "named", mcmAgentService.Name)
		return ctrl.Result{}, err
	}
	err = createResourceIfNotExists(mcmMediaProxyPv, types.NamespacedName{Name: mcmMediaProxyPv.Name, Namespace: "mcm"})
	if err != nil {
		log.Error(err, "Failed to create resource", "resource", mcmMediaProxyPv.GetObjectKind(), "named", mcmMediaProxyPv.Name)
		return ctrl.Result{}, err
	}
	err = createResourceIfNotExists(mcmMediaProxyPvc, types.NamespacedName{Name: mcmMediaProxyPvc.Name, Namespace: "mcm"})
	if err != nil {
		log.Error(err, "Failed to create resource", "resource", mcmMediaProxyPvc.GetObjectKind(), "named", mcmMediaProxyPvc.Name)
		return ctrl.Result{}, err
	}
	err = createResourceIfNotExists(mcmMediaProxyDs, types.NamespacedName{Name: mcmMediaProxyDs.Name, Namespace: "mcm"})
	if err != nil {
		log.Error(err, "Failed to create resource", "resource", mcmMediaProxyDs.GetObjectKind(), "named", mcmMediaProxyDs.Name)
		return ctrl.Result{}, err
	}
	err = createResourceIfNotExists(mtlManagerDeployment, types.NamespacedName{Name: mtlManagerDeployment.Name, Namespace: "mcm"})
	if err != nil {
		log.Error(err, "Failed to create resource", "resource", mtlManagerDeployment.GetObjectKind(), "named", mtlManagerDeployment.Name)
		return ctrl.Result{}, err
	}

	// Lookup the BcsConfig instance for this reconcile request
	bcsConf := &bcsv1.BcsConfig{}
	err = r.Get(ctx, req.NamespacedName, bcsConf)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("BcsConfig resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Error reading the object; Failed to get BcsConfig. \n ...Requeue...")
		return ctrl.Result{}, err
	}

	// Run all k8s resources for BCS pipeline and NMOS
	err = r.reconcileResources(ctx, bcsConf, log)
	if err != nil {
		log.Error(err, "Failed to reconcile resources for this custom resource")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *BcsConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bcsv1.BcsConfig{}).
		Complete(r)
}

func (r *BcsConfigReconciler) reconcileResources(ctx context.Context, bcs *bcsv1.BcsConfig, log logr.Logger) error {
	for iter, specInstance := range bcs.Spec {
		log.Info("Processing BcsConfig Spec", "instance number", iter, "name", specInstance.Name, "namespace", specInstance.Namespace)
		// Check if the namespace exists, if not create it
		namespace := &corev1.Namespace{}
		err := r.Get(ctx, types.NamespacedName{Name: specInstance.Namespace}, namespace)
		if err != nil && errors.IsNotFound(err) {
			namespace = &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: specInstance.Namespace,
				},
			}
			if err := r.Create(ctx, namespace); err != nil {
				log.Error(err, "Failed to create Namespace", "name", specInstance.Namespace)
				return err
			}
			log.Info("Namespace created successfully", "name", specInstance.Namespace)
		}

		// Reconcile ConfigMap
		if err := r.reconcileConfigMap(ctx, &specInstance, log); err != nil {
			log.Error(err, "Failed to reconcile ConfigMap")
			return err
		}

		// Reconcile Deployment
		if err := r.reconcileDeployment(ctx, &specInstance, log); err != nil {
			log.Error(err, "Failed to reconcile Deployment")
			return err
		}

		// Reconcile Service
		if err := r.reconcileService(ctx, &specInstance, log); err != nil {
			log.Error(err, "Failed to reconcile Service")
			return err
		}
	}

	return nil
}

func (r *BcsConfigReconciler) reconcileConfigMap(ctx context.Context, bcs *bcsv1.BcsConfigSpec, log logr.Logger) error {
	log.Info("Processing BcsConfig Spec", "name", bcs.Name, "namespace", bcs.Namespace)
	configMapName := bcs.Name + "-config"
	bcsConfigMap := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: configMapName, Namespace: bcs.Namespace}, bcsConfigMap)
	if err != nil && errors.IsNotFound(err) {
		bcsConfigMap = utils.CreateConfigMap(bcs)
		if err := r.Create(ctx, bcsConfigMap); err != nil {
			log.Error(err, "Failed to create ConfigMap")
			return err
		}
		log.Info("ConfigMap created successfully", "name", bcsConfigMap.Name, "namespace", bcsConfigMap.Namespace)
	} else if err != nil {
		log.Error(err, "Failed to get ConfigMap")
		return err
	} else {
		updatedConfigMap := utils.CreateConfigMap(bcs)
		bcsConfigMap.Data = updatedConfigMap.Data
		if err := r.Update(ctx, bcsConfigMap); err != nil {
			log.Error(err, "Failed to update ConfigMap")
			return err
		}
		log.Info("ConfigMap updated successfully", "name", bcsConfigMap.Name, "namespace", bcsConfigMap.Namespace)
	}
	return nil
}

func (r *BcsConfigReconciler) reconcileDeployment(ctx context.Context, bcs *bcsv1.BcsConfigSpec, log logr.Logger) error {
	bcsDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: bcs.Name, Namespace: bcs.Namespace}, bcsDeployment)
	if errors.IsNotFound(err) {
		bcsDeployment = utils.CreateBcsDeployment(bcs)
		if err := r.Create(ctx, bcsDeployment); err != nil {
			log.Error(err, "Failed to create Deployment")
			return err
		}
		log.Info("Deployment is created successfully", "name", bcsDeployment.Name, "namespace", bcsDeployment.Namespace)
	} else if err != nil {
		log.Error(err, "Failed to create/update Deployment. Check your either cluster or bcs launcher configuration")
		return err
	} else {
		if err := r.Update(ctx, bcsDeployment); err != nil {
			log.Error(err, "Failed to update Deployment")
			return err
		}
		log.Info("Deployment is updated successfully", "name", bcsDeployment.Name, "namespace", bcsDeployment.Namespace)
	}
	return nil
}

func (r *BcsConfigReconciler) reconcileService(ctx context.Context, bcs *bcsv1.BcsConfigSpec, log logr.Logger) error {
	bcsSevice := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: bcs.Name, Namespace: bcs.Namespace}, bcsSevice)
	if errors.IsNotFound(err) {
		bcsSevice = utils.CreateBcsService(bcs)
		if err := r.Create(ctx, bcsSevice); err != nil {
			log.Error(err, "Failed to create Service")
			return err
		}
		log.Info("Service is created successfully", "name", bcsSevice.Name, "namespace", bcsSevice.Namespace)
	} else if err != nil {
		log.Error(err, "Failed to create/update Service. Check your either cluster or bcs launcher configuration")
		return err
	} else {
		if err := r.Update(ctx, bcsSevice); err != nil {
			log.Error(err, "Failed to update Service")
			return err
		}
		log.Info("Service is updated successfully", "name", bcsSevice.Name, "namespace", bcsSevice.Namespace)
	}

	return nil
}
