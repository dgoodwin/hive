package cluster

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	hivev1 "github.com/openshift/hive/pkg/apis/hive/v1"
	hivev1azure "github.com/openshift/hive/pkg/apis/hive/v1/azure"

	installertypes "github.com/openshift/installer/pkg/types"
	azureinstallertypes "github.com/openshift/installer/pkg/types/azure"
)

const (
	azureCredFile     = "osServicePrincipal.json"
	azureRegion       = "centralus"
	azureInstanceType = "Standard_D2s_v3"
)

var _ CloudProvider = (*AzureCloudProvider)(nil)

type AzureCloudProvider struct {
	// ServicePrincipal is the bytes from a service principal file, typically ~/.azure/osServicePrincipal.json.
	ServicePrincipal []byte

	// BaseDomainResourceGroupName is the resource group where the base domain for this cluster is configured.
	BaseDomainResourceGroupName string

	// ReuseCredsSecret is a reference to a pre-existing credentials secret for this cloud.
	ReuseCredsSecret *corev1.LocalObjectReference
}

func (p *AzureCloudProvider) GenerateCredentialsSecret(o *Generator) *corev1.Secret {
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
			azureCredFile: p.ServicePrincipal,
		},
	}
}

func (p *AzureCloudProvider) AddClusterDeploymentPlatform(o *Generator, cd *hivev1.ClusterDeployment) {
	cd.Spec.Platform = hivev1.Platform{
		Azure: &hivev1azure.Platform{
			CredentialsSecretRef: corev1.LocalObjectReference{
				Name: p.credsSecretName(o),
			},
			Region:                      azureRegion,
			BaseDomainResourceGroupName: p.BaseDomainResourceGroupName,
		},
	}
}

func (p *AzureCloudProvider) AddMachinePoolPlatform(o *Generator, mp *hivev1.MachinePool) {
	mp.Spec.Platform.Azure = &hivev1azure.MachinePool{
		InstanceType: azureInstanceType,
		OSDisk: hivev1azure.OSDisk{
			DiskSizeGB: 128,
		},
	}

}

func (p *AzureCloudProvider) AddInstallConfigPlatform(o *Generator, ic *installertypes.InstallConfig) {
	// Inject platform details into InstallConfig:
	ic.Platform = installertypes.Platform{
		Azure: &azureinstallertypes.Platform{
			Region:                      azureRegion,
			BaseDomainResourceGroupName: p.BaseDomainResourceGroupName,
		},
	}

	// Used for both control plane and workers.
	mpp := &azureinstallertypes.MachinePool{}
	ic.ControlPlane.Platform.Azure = mpp
	ic.Compute[0].Platform.Azure = mpp
}

func (p *AzureCloudProvider) credsSecretName(o *Generator) string {
	return fmt.Sprintf("%s-azure-creds", o.Name)
}
