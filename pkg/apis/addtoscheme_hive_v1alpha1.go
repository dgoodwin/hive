package apis

import (
	hivev1alpha1 "github.com/openshift/hive/pkg/apis/hive/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, hivev1alpha1.SchemeBuilder.AddToScheme)
}
