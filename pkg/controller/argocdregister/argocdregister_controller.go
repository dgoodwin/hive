/*
Copyright (C) 2019 Red Hat, Inc.

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

// Package argocdregister provides a controller which ensures provisioned clusters are added
// to the ArgoCD cluster registry, and removed when the cluster is deprovisioned.
package argocdregister

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"k8s.io/client-go/kubernetes"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"
	kapi "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/argoproj/argo-cd/common"

	hivev1 "github.com/openshift/hive/pkg/apis/hive/v1"
	"github.com/openshift/hive/pkg/constants"
	hivemetrics "github.com/openshift/hive/pkg/controller/metrics"
	controllerutils "github.com/openshift/hive/pkg/controller/utils"
	"github.com/openshift/hive/pkg/remoteclient"
	"github.com/openshift/hive/pkg/resource"
)

const (
	controllerName = "argocdregister"
	// TODO: pass this through from HiveConfig via an EnvVar
	argoCDNamespace          = "argocd"
	argoCDServiceAccountName = "argocd-manager"

	adminKubeConfigKey = "kubeconfig"
)

// Add creates a new Argocdregister Controller and adds it to the Manager with default RBAC. The Manager will set fields on the
// Controller and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	kubeClient, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		return err
	}
	// TODO: tidy up init
	r := NewReconciler(mgr, kubeClient)
	r.remoteClientBuilder = func(cd *hivev1.ClusterDeployment) remoteclient.Builder {
		return remoteclient.NewBuilder(r.Client, cd, controllerName)
	}
	return AddToManager(mgr, r)
}

// NewReconciler returns a new reconcile.Reconciler
func NewReconciler(mgr manager.Manager, kubeClient kubernetes.Interface) *ArgoCDRegisterController {
	return &ArgoCDRegisterController{
		Client:     controllerutils.NewClientWithMetricsOrDie(mgr, controllerName),
		scheme:     mgr.GetScheme(),
		restConfig: mgr.GetConfig(),
		logger:     log.WithField("controller", controllerName),
		kubeClient: kubeClient,
	}
}

// AddToManager adds a new Controller to mgr with r as the reconcile.Reconciler
func AddToManager(mgr manager.Manager, r reconcile.Reconciler) error {

	// Create a new controller
	c, err := controller.New("argocdregister-controller", mgr, controller.Options{Reconciler: r, MaxConcurrentReconciles: controllerutils.GetConcurrentReconciles()})
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

var _ reconcile.Reconciler = &ArgoCDRegisterController{}

// ArgoCDRegisterController reconciles the MachineSets generated from a ClusterDeployment object
type ArgoCDRegisterController struct {
	client.Client
	kubeClient kubernetes.Interface
	scheme     *runtime.Scheme
	restConfig *rest.Config
	logger     log.FieldLogger

	// remoteClientBuilder is a function pointer to the function that gets a builder for building a client
	// for the remote cluster's API server
	remoteClientBuilder func(cd *hivev1.ClusterDeployment) remoteclient.Builder
}

// Reconcile checks if we can establish an API client connection to the remote cluster and maintains the unreachable condition as a result.
func (r *ArgoCDRegisterController) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	start := time.Now()
	cdLog := log.WithFields(log.Fields{
		"clusterDeployment": request.Name,
		"namespace":         request.Namespace,
		"controller":        controllerName,
	})

	// For logging, we need to see when the reconciliation loop starts and ends.
	cdLog.Info("reconciling cluster deployment")
	defer func() {
		dur := time.Since(start)
		hivemetrics.MetricControllerReconcileTime.WithLabelValues(controllerName).Observe(dur.Seconds())
		cdLog.WithField("elapsed", dur).Info("reconcile complete")
	}()

	cd := &hivev1.ClusterDeployment{}
	err := r.Get(context.TODO(), request.NamespacedName, cd)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		// Error reading the object - requeue the request
		log.WithError(err).Error("error looking up cluster deployment")
		return reconcile.Result{}, err
	}

	if !cd.Spec.Installed {
		cdLog.Info("cluster installation is not complete")
		return reconcile.Result{}, nil
	}

	if cd.Spec.ClusterMetadata == nil {
		cdLog.Error("installed cluster with no cluster metadata")
		return reconcile.Result{}, nil
	}

	if cd.DeletionTimestamp == nil {
		return r.reconcileCluster(cd, cdLog)
	} else {
		return r.reconcileDeletedCluster(cd, cdLog)
	}
}

func (r *ArgoCDRegisterController) reconcileCluster(cd *hivev1.ClusterDeployment, cdLog log.FieldLogger) (reconcile.Result, error) {
	cdLog.Info("reconciling active cluster")

	if cd.Status.APIURL == "" {
		cdLog.Info("Installed cluster does not have Status.APIURL set yet")
		return reconcile.Result{}, fmt.Errorf("installed cluster does not have Status.APIURL set yet")
	}

	h := resource.NewHelperFromRESTConfig(r.restConfig, cdLog)

	// Determine unique and predictable name for the ArgoCD secret.
	clusterSecretName, err := getPredictableSecretName(cd.Status.APIURL)
	if err != nil {
		cdLog.WithError(err).Error("error getting predictable secret name")
		return reconcile.Result{}, err
	}
	cdLog = cdLog.WithField("clusterSecret", clusterSecretName)

	kubeconfigSecretName := cd.Spec.ClusterMetadata.AdminKubeconfigSecretRef.Name
	kubeconfig, err := r.loadSecretData(kubeconfigSecretName, cd.Namespace, adminKubeConfigKey)
	if err != nil {
		cdLog.WithError(err).Error("unable to load cluster admin kubeconfig")
		return reconcile.Result{}, err
	}

	/*
	   remoteClientBuilder := r.remoteClusterAPIClientBuilder(cd)
	   // If the cluster is unreachable, do not reconcile.
	   if remoteClientBuilder.Unreachable() {
	       logger.Debug("skipping cluster with unreachable condition")
	       return reconcile.Result{}, nil
	   }
	*/

	// Parse the clusters kubeconfig so we can get the fields we need for argo's config:
	config, err := clientcmd.Load([]byte(kubeconfig))
	if err != nil {
		cdLog.WithError(err).Error("unable to load cluster kubeconfig")
		return reconcile.Result{}, err
	}
	kubeConfig := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{})
	cfg, err := kubeConfig.ClientConfig()
	if err != nil {
		cdLog.WithError(err).Error("unable to load client config")
		return reconcile.Result{}, err
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		cdLog.WithError(err).Error("unable to create clientset for remote cluster")
		return reconcile.Result{}, err
	}
	//managerBearerToken, err := r.loadArgoCDServiceAccountToken(clientset)
	managerBearerToken, err := common.InstallClusterManagerRBAC(clientset)
	if err != nil {
		cdLog.WithError(err).Error("unable to load argocd service account token")
		return reconcile.Result{}, err
	}
	// TODO: delete this. really.
	cdLog.Debug("loaded manager bearer token: %s", managerBearerToken)

	tlsClientConfig := TLSClientConfig{
		Insecure:   cfg.TLSClientConfig.Insecure,
		ServerName: cfg.TLSClientConfig.ServerName,
		CAData:     cfg.TLSClientConfig.CAData,
	}

	// Argo uses a custom format for their server config blob, not a kubeconfig:
	argoCDServerConfig := ClusterConfig{
		BearerToken:     managerBearerToken,
		TLSClientConfig: tlsClientConfig,
	}

	data := make(map[string][]byte)
	data["server"] = []byte(cd.Status.APIURL)
	data["name"] = []byte(cd.Name)
	configBytes, err := json.Marshal(argoCDServerConfig)
	if err != nil {
		return reconcile.Result{}, err
	}
	data["config"] = configBytes

	argoClusterSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterSecretName,
			Namespace: argoCDNamespace,
			Labels: map[string]string{
				"argocd.argoproj.io/secret-type": "cluster",
				constants.CreatedByHiveLabel:     "true",
			},
		},
		Data: data,
	}

	// Copy all ClusterDeployment labels onto the ArgoCD cluster secret. This will hopefully
	// allow for dynamic generation of ArgoCD Applications. (separate tooling)
	for k, v := range cd.Labels {
		argoClusterSecret.Labels[k] = v
	}

	cdLog.Info("applying ArgoCD cluster secret")
	result, err := h.ApplyRuntimeObject(argoClusterSecret, scheme.Scheme)
	if err != nil {
		cdLog.WithError(err).Error("error applying ArgoCD cluster secret")
		return reconcile.Result{}, err
	}
	// TODO: only log if changed
	cdLog.Infof("ArgoCD cluster secret applied (%s)", result)

	return reconcile.Result{}, nil
}

func (r *ArgoCDRegisterController) loadArgoCDServiceAccountToken(clientset kubernetes.Interface) (string, error) {

	serviceAccount := &corev1.ServiceAccount{}
	err := r.Client.Get(context.Background(),
		types.NamespacedName{
			Name:      argoCDServiceAccountName,
			Namespace: "kube-system",
		}, serviceAccount)
	if err != nil {
		return "", fmt.Errorf("error looking up %s service account: %v", argoCDServiceAccountName, err)
	}
	if len(serviceAccount.Secrets) == 0 {
		return "", fmt.Errorf("%s service account has no secrets", argoCDServiceAccountName)
	}

	secretName := serviceAccount.Secrets[0].Name

	secret := &corev1.Secret{}
	err = r.Client.Get(context.Background(),
		types.NamespacedName{
			Name:      secretName,
			Namespace: "kube-system",
		}, secret)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve secret %q: %v", secretName, err)
	}
	token, ok := secret.Data["token"]
	if !ok {
		return "", fmt.Errorf("secret %q for service account %q has no token", secretName, serviceAccount)
	}
	return string(token), nil
}

func (r *ArgoCDRegisterController) reconcileDeletedCluster(cd *hivev1.ClusterDeployment, cdLog log.FieldLogger) (reconcile.Result, error) {
	cdLog.Info("reconciling deleted cluster")

	// Determine unique and predictable name for the ArgoCD secret.
	//clusterSecretName := getPredictableSecretName(cd)

	return reconcile.Result{}, nil
}

func (r *ArgoCDRegisterController) loadSecretData(secretName, namespace, dataKey string) (string, error) {
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

// getPredictableSecretName generates a unique secret name by hashing the server API URL,
// which is required as all cluster secrets land in the argocd namespace.
// This code matches what is presently done in ArgoCD (util/db/cluster.go). However the actual
// name of the secret does not matter , so we run limited risk of the implementation changing
// out from underneath us.
func getPredictableSecretName(serverAddr string) (string, error) {
	serverURL, err := url.ParseRequestURI(serverAddr)
	if err != nil {
		return "", err
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(serverAddr))
	host := strings.ToLower(strings.Split(serverURL.Host, ":")[0])
	return fmt.Sprintf("cluster-%s-%v", host, h.Sum32()), nil
}
