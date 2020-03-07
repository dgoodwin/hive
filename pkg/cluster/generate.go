package cluster

import (
	"fmt"
	hivev1 "github.com/openshift/hive/pkg/apis/hive/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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

	// Labels are labels to be added to the ClusterDeployment.
	Labels map[string]string

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

	// ServingCert is the contents of a serving certificate to be used for the cluster.
	ServingCert string

	// ServingCertKey is the contents of a key for the ServingCert.
	ServingCertKey string

	// Adopt is a flag indicating we're adopting a pre-existing cluster.
	Adopt bool

	// AdoptAdminKubeconfig is a cluster administrator admin kubeconfig typically obtained
	// from openshift-install. Should be set when adopting pre-existing clusters.
	AdoptAdminKubeconfig []byte

	// AdoptClusterID is the unique generated ID for a cluster being adopted.
	AdoptClusterID string

	// AdoptInfraID is the unique generated infrastructure ID for a cluster being adopted.
	AdoptInfraID string
}

func (o *Generator) GenerateAll() []runtime.Object {
	allObjects := []runtime.Object{}
	allObjects = append(allObjects, o.GenerateClusterDeployment())

	// TODO: maintain "include secrets" flag functionality?
	pullSecretSecret := o.GeneratePullSecretSecret()
	if pullSecretSecret != nil {
		allObjects = append(allObjects, pullSecretSecret)
	}
	sshPrivateKeySecret := o.GenerateSSHPrivateKeySecret()
	if sshPrivateKeySecret != nil {
		allObjects = append(allObjects, sshPrivateKeySecret)
	}
	servingCertSecret := o.GenerateServingCertSecret()
	if servingCertSecret != nil {
		allObjects = append(allObjects, servingCertSecret)
	}

	if o.Adopt {
		allObjects = append(allObjects, o.GenerateAdminKubeconfigSecret())
	}
	return allObjects
}

// GenerateClusterDeployment generates a new cluster deployment
func (o *Generator) GenerateClusterDeployment() *hivev1.ClusterDeployment {
	cd := &hivev1.ClusterDeployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterDeployment",
			APIVersion: hivev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        o.Name,
			Namespace:   o.Namespace,
			Annotations: map[string]string{},
			Labels:      o.Labels,
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

	if o.Adopt {
		cd.Spec.ClusterMetadata = &hivev1.ClusterMetadata{
			ClusterID:                o.AdoptClusterID,
			InfraID:                  o.AdoptInfraID,
			AdminKubeconfigSecretRef: corev1.LocalObjectReference{Name: o.getAdoptAdminKubeconfigSecretName()},
		}
	}

	return cd
}

// GeneratePullSecretSecret returns a Kubernetes Secret containing the pull secret to be
// used for pulling images.
func (o *Generator) GeneratePullSecretSecret() *corev1.Secret {
	if len(o.PullSecret) == 0 {
		return nil
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
	}
}

// GenerateSSHPrivateKeySecret returns a Kubernetes Secret containing the SSH private
// key to be used.
func (o *Generator) GenerateSSHPrivateKeySecret() *corev1.Secret {
	if o.SSHPrivateKey == "" {
		return nil
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
	}
}

func (o *Generator) GenerateServingCertSecret() *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.getServingCertSecretName(),
			Namespace: o.Namespace,
		},
		Type: corev1.SecretTypeTLS,
		StringData: map[string]string{
			"tls.crt": string(o.ServingCert),
			"tls.key": string(o.ServingCertKey),
		},
	}
}

func (o *Generator) GenerateAdminKubeconfigSecret() *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.getAdoptAdminKubeconfigSecretName(),
			Namespace: o.Namespace,
		},
		Data: map[string][]byte{
			"kubeconfig":     o.AdoptAdminKubeconfig,
			"raw-kubeconfig": o.AdoptAdminKubeconfig,
		},
	}
}

func (o *Generator) getServingCertSecretName() string {
	return fmt.Sprintf("%s-serving-cert", o.Name)
}

func (o *Generator) getAdoptAdminKubeconfigSecretName() string {
	return fmt.Sprintf("%s-adopted-admin-kubeconfig", o.Name)
}

// TODO: handle long cluster names.
func (o *Generator) getSSHPrivateKeySecretName() string {
	return fmt.Sprintf("%s-ssh-private-key", o.Name)
}

// TODO: handle long cluster names.
func (o *Generator) getPullSecretSecretName() string {
	return fmt.Sprintf("%s-pull-secret", o.Name)
}
