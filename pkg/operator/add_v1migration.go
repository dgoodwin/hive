package operator

import (
	"github.com/openshift/hive/pkg/operator/v1migration"
)

func init() {
	// AddToOperatorFuncs is a list of functions to create controllers and add them to an operator manager.
	AddToOperatorFuncs = append(AddToOperatorFuncs, v1migration.Add)
}
