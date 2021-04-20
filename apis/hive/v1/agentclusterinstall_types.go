package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterInstallSpec defines the desired state of the Agent ClusterInstall.
type ClusterInstallSpec struct {

	// ImageSetRef is a reference to a ClusterImageSet. The release image specified in the ClusterImageSet will be used
	// by clusters created for this cluster pool.
	ImageSetRef ClusterImageSetReference `json:"imageSetRef"`
}

// ClusterInstallStatus defines the observed state of ClusterPool
type ClusterInstallStatus struct {
	// Conditions includes more detailed status for the cluster install.
	// +optional
	Conditions []ClusterInstallCondition `json:"conditions,omitempty"`
}

// ClusterInstallCondition contains details for the current condition of a cluster install.
type ClusterInstallCondition struct {
	// Type is the type of the condition.
	Type string `json:"type"`
	// Status is the status of the condition.
	Status corev1.ConditionStatus `json:"status"`
	// LastProbeTime is the last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Reason is a unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Message is a human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

const (
	ClusterInstallProvisioned     = "Provisioned"
	ClusterInstallFailed          = "Failed"
	ClusterInstallCompleted       = "Completed"
	ClusterInstallRequirementsMet = "RequirementsMet"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterInstall represents a request to provision an agent based cluster.
//
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type ClusterInstall struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterInstallSpec   `json:"spec"`
	Status ClusterInstallStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterInstallList contains a list of ClusterInstalls
type ClusterInstallList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterInstall `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterInstall{}, &ClusterInstallList{})
}
