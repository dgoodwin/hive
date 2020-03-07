package cluster

import (
	"fmt"
	hivev1 "github.com/openshift/hive/pkg/apis/hive/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	deleteAfterAnnotation    = "hive.openshift.io/delete-after"
	tryInstallOnceAnnotation = "hive.openshift.io/try-install-once"
)

// Generator can be used to build all artifacts required for to create a ClusterDeployment.
type Generator struct {
	// Name is the name of your Cluster. Will be used for both the ClusterDeployment.Name and the
	// ClusterDeployment.Spec.ClusterName, which encompasses the subdomain and cloud provider resource
	// tagging.
	Name string

	// Namespace where the ClusterDeployment and all associated artifacts will be created.
	Namespace string

	// PullSecret is the secret to use when pulling images.
	PullSecret string

	// SSHPrivateKey is an optional SSH key to configure on hosts in the cluster. This would
	// typically be read from ~/.ssh/id_rsa.
	SSHPrivateKey string

	// InstallOnce indicates that the provision job should not be retried on failure.
	InstallOnce bool

	// BaseDomain is the DNS base domain to be used for the cluster.
	BaseDomain string

	// ManageDNS can be set to true to enable Hive's automatic DNS zone creation and forwarding. (assuming
	// this is properly configured in HiveConfig)
	ManageDNS bool

	// DeleteAfter is the duration after which the cluster should be automatically destroyed, relative to
	// creationTimestamp. Stored as an annotation on the ClusterDeployment.
	DeleteAfter string
}

// GenerateClusterDeployment generates a new cluster deployment
func (o *Generator) GenerateClusterDeployment(labels map[string]string) (*hivev1.ClusterDeployment, error) {
	cd := &hivev1.ClusterDeployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterDeployment",
			APIVersion: hivev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        o.Name,
			Namespace:   o.Namespace,
			Annotations: map[string]string{},
			Labels:      labels,
		},
		Spec: hivev1.ClusterDeploymentSpec{
			ClusterName:  o.Name,
			BaseDomain:   o.BaseDomain,
			ManageDNS:    o.ManageDNS,
			Provisioning: &hivev1.Provisioning{},
		},
	}

	if o.SSHPrivateKey != "" {
		cd.Spec.Provisioning.SSHPrivateKeySecretRef = &corev1.LocalObjectReference{Name: o.getSSHPrivateKeySecretName()}
	}

	if o.InstallOnce {
		cd.Annotations[tryInstallOnceAnnotation] = "true"
	}

	if o.PullSecret != "" {
		cd.Spec.PullSecretRef = &corev1.LocalObjectReference{Name: o.getPullSecretSecretName()}
	}

	if len(o.ServingCert) > 0 {
		cd.Spec.CertificateBundles = []hivev1.CertificateBundleSpec{
			{
				Name: "serving-cert",
				CertificateSecretRef: corev1.LocalObjectReference{
					Name: fmt.Sprintf("%s-serving-cert", o.Name),
				},
			},
		}
		cd.Spec.ControlPlaneConfig.ServingCertificates.Default = "serving-cert"
		cd.Spec.Ingress = []hivev1.ClusterIngress{
			{
				Name:               "default",
				Domain:             fmt.Sprintf("apps.%s.%s", o.Name, o.BaseDomain),
				ServingCertificate: "serving-cert",
			},
		}
	}

	if o.DeleteAfter != "" {
		cd.ObjectMeta.Annotations[deleteAfterAnnotation] = o.DeleteAfter
	}

	return cd, nil
}

// GeneratePullSecretSecret returns a Kubernetes Secret containing the pull secret to be
// used for pulling images.
func (o *Generator) GeneratePullSecretSecret() (*corev1.Secret, error) {
	if len(o.PullSecret) == 0 {
		return nil, nil
	}
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.getPullSecretSecretName(),
			Namespace: o.Namespace,
		},
		Type: corev1.SecretTypeDockerConfigJson,
		StringData: map[string]string{
			corev1.DockerConfigJsonKey: o.PullSecret,
		},
	}, nil
}

// GenerateSSHPrivateKeySecret returns a Kubernetes Secret containing the SSH private
// key to be used.
func (o *Generator) GenerateSSHPrivateKeySecret() (*corev1.Secret, error) {
	if o.SSHPrivateKey == "" {
		return nil, nil
	}
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.getSSHPrivateKeySecretName(),
			Namespace: o.Namespace,
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"ssh-privatekey": o.SSHPrivateKey,
		},
	}, nil
}

// TODO: handle long cluster names.
func (o *Generator) getSSHPrivateKeySecretName() string {
	return fmt.Sprintf("%s-ssh-private-key", o.Name)
}

// TODO: handle long cluster names.
func (o *Generator) getPullSecretSecretName() string {
	return fmt.Sprintf("%s-pull-secret", o.Name)
}
