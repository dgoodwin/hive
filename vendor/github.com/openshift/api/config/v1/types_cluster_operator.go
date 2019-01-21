package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterOperator is the Custom Resource object which holds the current state
// of an operator. This object is used by operators to convey their state to
// the rest of the cluster.
type ClusterOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	// spec hold the intent of how this operator should behave.
	Spec ClusterOperatorSpec `json:"spec"`

	// status holds the information about the state of an operator.  It is consistent with status information across
	// the kube ecosystem.
	Status ClusterOperatorStatus `json:"status"`
}

// ClusterOperatorSpec is empty for now, but you could imagine holding information like "pause".
type ClusterOperatorSpec struct {
}

// ClusterOperatorStatus provides information about the status of the operator.
// +k8s:deepcopy-gen=true
type ClusterOperatorStatus struct {
	// conditions describes the state of the operator's reconciliation functionality.
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []ClusterOperatorStatusCondition `json:"conditions"  patchStrategy:"merge" patchMergeKey:"type"`

	// versions is a slice of operand version tuples.  Operators which manage multiple operands will have multiple
	// entries in the array.  If an operator is Available, it must have at least one entry.  Report the version of
	// the operator itself in addition to the version of its operands is recommended.
	Versions []OperandVersion `json:"versions"`

	// relatedObjects is a list of objects that are "interesting" or related to this operator.  Common uses are:
	// 1. the detailed resource driving the operator
	// 2. operator namespaces
	// 3. operand namespaces
	RelatedObjects []ObjectReference `json:"relatedObjects"`

	// extension contains any additional status information specific to the
	// operator which owns this status object.
	Extension runtime.RawExtension `json:"extension,omitempty"`
}

type OperandVersion struct {
	// name is the name of the particular operand this version is for.  It usually matches container images, not operators.
	Name string `json:"name"`

	// version indicates which version of a particular operand is currently being manage.  It must always match the Available
	// condition.  If 1.0.0 is Available, then this must indicate 1.0.0 even if the operator is trying to rollout
	// 1.1.0
	Version string `json:"version"`
}

// ObjectReference contains enough information to let you inspect or modify the referred object.
type ObjectReference struct {
	// group of the referent.
	Group string `json:"group"`
	// resource of the referent.
	Resource string `json:"resource"`
	// namespace of the referent.
	// +optional
	Namespace string `json:"namespace,omitempty"`
	// name of the referent.
	Name string `json:"name"`
}

type ConditionStatus string

// These are valid condition statuses. "ConditionTrue" means a resource is in the condition.
// "ConditionFalse" means a resource is not in the condition. "ConditionUnknown" means kubernetes
// can't decide if a resource is in the condition or not. In the future, we could add other
// intermediate conditions, e.g. ConditionDegraded.
const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

// ClusterOperatorStatusCondition represents the state of the operator's
// reconciliation functionality.
// +k8s:deepcopy-gen=true
type ClusterOperatorStatusCondition struct {
	// type specifies the state of the operator's reconciliation functionality.
	Type ClusterStatusConditionType `json:"type"`

	// status of the condition, one of True, False, Unknown.
	Status ConditionStatus `json:"status"`

	// lastTransitionTime is the time of the last update to the current status object.
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`

	// reason is the reason for the condition's last transition.  Reasons are CamelCase
	Reason string `json:"reason,omitempty"`

	// message provides additional information about the current condition.
	// This is only to be consumed by humans.
	Message string `json:"message,omitempty"`
}

// ClusterStatusConditionType is the state of the operator's reconciliation functionality.
type ClusterStatusConditionType string

const (
	// OperatorAvailable indicates that the binary maintained by the operator (eg: openshift-apiserver for the
	// openshift-apiserver-operator), is functional and available in the cluster.
	OperatorAvailable ClusterStatusConditionType = "Available"

	// OperatorProgressing indicates that the operator is actively making changes to the binary maintained by the
	// operator (eg: openshift-apiserver for the openshift-apiserver-operator).
	OperatorProgressing ClusterStatusConditionType = "Progressing"

	// OperatorFailing indicates that the operator has encountered an error that is preventing it from working properly.
	// The binary maintained by the operator (eg: openshift-apiserver for the openshift-apiserver-operator) may still be
	// available, but the user intent cannot be fulfilled.
	OperatorFailing ClusterStatusConditionType = "Failing"
)

// ClusterOperatorList is a list of OperatorStatus resources.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ClusterOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ClusterOperator `json:"items"`
}
