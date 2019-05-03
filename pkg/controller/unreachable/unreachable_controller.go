/*
Copyright 2018 The Kubernetes Authors.

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

// Package unreachable provides a controller which periodically checks if a remote cluster is reachable
// and maintains a condition on the cluster as a result. If the unreachable condition is true, other controllers
// can skip attempts to reach the cluster which require a 30 second timeout.
package unreachable

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"
	kapi "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	hivev1 "github.com/openshift/hive/pkg/apis/hive/v1alpha1"
	controllerutils "github.com/openshift/hive/pkg/controller/utils"
)

const (
	controllerName = "unreachable"

	adminKubeConfigKey          = "kubeconfig"
	adminCredsSecretPasswordKey = "password"
	adminSSHKeySecretKey        = "ssh-publickey"

	maxUnreachableCheck = 2 * time.Hour
)

// Add creates a new Unreachable Controller and adds it to the Manager with default RBAC. The Manager will set fields on the
// Controller and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return AddToManager(mgr, NewReconciler(mgr))
}

// NewReconciler returns a new reconcile.Reconciler
func NewReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileRemoteMachineSet{
		Client:                        mgr.GetClient(),
		scheme:                        mgr.GetScheme(),
		logger:                        log.WithField("controller", controllerName),
		remoteClusterAPIClientBuilder: controllerutils.BuildClusterAPIClientFromKubeconfig,
	}
}

// AddToManager adds a new Controller to mgr with r as the reconcile.Reconciler
func AddToManager(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("unreachable-controller", mgr, controller.Options{Reconciler: r, MaxConcurrentReconciles: controllerutils.GetConcurrentReconciles()})
	if err != nil {
		return err
	}

	// Watch for changes to ClusterDeployment
	err = c.Watch(&source.Kind{Type: &hivev1.ClusterDeployment{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileRemoteMachineSet{}

// ReconcileRemoteMachineSet reconciles the MachineSets generated from a ClusterDeployment object
type ReconcileRemoteMachineSet struct {
	client.Client
	scheme *runtime.Scheme

	logger log.FieldLogger

	// remoteClusterAPIClientBuilder is a function pointer to the function that builds a client for the
	// remote cluster's cluster-api
	remoteClusterAPIClientBuilder func(string) (client.Client, error)
}

// Reconcile checks if we can establish an API client connection to the remote cluster and maintains the unreachable condition as a result.
func (r *ReconcileRemoteMachineSet) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the ClusterDeployment instance
	cd := &hivev1.ClusterDeployment{}
	err := r.Get(context.TODO(), request.NamespacedName, cd)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request
		log.WithError(err).Error("error looking up cluster deployment")
		return reconcile.Result{}, err
	}

	cdLog := r.logger.WithFields(log.Fields{
		"clusterDeployment": cd.Name,
		"namespace":         cd.Namespace,
	})

	// If the clusterdeployment is deleted, do not reconcile.
	if cd.DeletionTimestamp != nil {
		cdLog.Debug("cluster has deletion timestamp")
		return reconcile.Result{}, nil
	}

	if !cd.Status.Installed {
		// Cluster isn't installed yet, return
		cdLog.Debug("cluster installation is not complete")
		return reconcile.Result{}, nil
	}

	// Check if we're due to recheck this clusters connectivity:
	cond := controllerutils.FindClusterDeploymentCondition(cd.Status.Conditions, hivev1.UnreachableCondition)
	if cond != nil {
		sinceLastTransition := time.Since(cond.LastTransitionTime.Time)
		sinceLastProbe := time.Since(cond.LastProbeTime.Time)
		if sinceLastProbe < (sinceLastTransition/5) && sinceLastProbe < maxUnreachableCheck {
			cdLog.WithFields(log.Fields{
				"sinceLastTransition": sinceLastTransition,
				"sinceLastProbe":      sinceLastProbe,
			}).Debug("skipping unreachable check")
			return reconcile.Result{}, nil
		}
	}

	cdLog.Info("checking if cluster is reachable")

	secretName := cd.Status.AdminKubeconfigSecret.Name
	secretData, err := r.loadSecretData(secretName, cd.Namespace, adminKubeConfigKey)
	if err != nil {
		cdLog.WithError(err).Error("unable to load admin kubeconfig")
		return reconcile.Result{}, err
	}

	_, err = r.remoteClusterAPIClientBuilder(secretData)
	if err != nil {
		cdLog.WithError(err).Info("unable to create remote API client, marking cluster unreachable")
		cd.Status.Conditions = controllerutils.SetClusterDeploymentCondition(cd.Status.Conditions, hivev1.UnreachableCondition,
			corev1.ConditionTrue, "ErrorConnectingToCluster", err.Error(), controllerutils.UpdateConditionAlways)
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileRemoteMachineSet) loadSecretData(secretName, namespace, dataKey string) (string, error) {
	s := &kapi.Secret{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: secretName, Namespace: namespace}, s)
	if err != nil {
		return "", err
	}
	retStr, ok := s.Data[dataKey]
	if !ok {
		return "", fmt.Errorf("secret %s did not contain key %s", secretName, dataKey)
	}
	return string(retStr), nil
}
