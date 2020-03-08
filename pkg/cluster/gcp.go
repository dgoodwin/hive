package cluster

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	installertypes "github.com/openshift/installer/pkg/types"
	installergcp "github.com/openshift/installer/pkg/types/gcp"

	hivev1 "github.com/openshift/hive/pkg/apis/hive/v1"
	hivev1gcp "github.com/openshift/hive/pkg/apis/hive/v1/gcp"
	"github.com/openshift/hive/pkg/constants"
)

const (
	gcpRegion       = "us-east1"
	gcpInstanceType = "n1-standard-4"
)

var _ CloudProvider = (*GCPCloudProvider)(nil)

type GCPCloudProvider struct {
	// ServicePrincipal is the bytes from a service account file, typically ~/.gcp/osServiceAccount.json.
	ServiceAccount []byte

	// ReuseCredsSecret is a reference to a pre-existing credentials secret for this cloud.
	ReuseCredsSecret *corev1.LocalObjectReference

	// ProjectID is the GCP project to use.
	ProjectID string
}

func (p *GCPCloudProvider) GenerateCredentialsSecret(o *Generator) *corev1.Secret {
	if p.ReuseCredsSecret != nil {
		return nil
	}
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.credsSecretName(o),
			Namespace: o.Namespace,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			constants.GCPCredentialsName: p.ServiceAccount,
		},
	}
}

func (p *GCPCloudProvider) AddClusterDeploymentPlatform(o *Generator, cd *hivev1.ClusterDeployment) {
	cd.Spec.Platform = hivev1.Platform{
		GCP: &hivev1gcp.Platform{
			CredentialsSecretRef: corev1.LocalObjectReference{
				Name: p.credsSecretName(o),
			},
			Region: gcpRegion,
		},
	}
}

func (p *GCPCloudProvider) AddMachinePoolPlatform(o *Generator, mp *hivev1.MachinePool) {
	mp.Spec.Platform.GCP = &hivev1gcp.MachinePool{
		InstanceType: gcpInstanceType,
	}

}

func (p *GCPCloudProvider) AddInstallConfigPlatform(o *Generator, ic *installertypes.InstallConfig) {
	ic.Platform = installertypes.Platform{
		GCP: &installergcp.Platform{
			ProjectID: p.ProjectID,
			Region:    gcpRegion,
		},
	}

	// Used for both control plane and workers.
	mpp := &installergcp.MachinePool{
		InstanceType: gcpInstanceType,
	}
	ic.ControlPlane.Platform.GCP = mpp
	ic.Compute[0].Platform.GCP = mpp
}

func (p *GCPCloudProvider) credsSecretName(o *Generator) string {
	return fmt.Sprintf("%s-gcp-creds", o.Name)
}
