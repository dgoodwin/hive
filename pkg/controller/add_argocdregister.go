package controller

import (
	hiveconstants "github.com/openshift/hive/pkg/constants"
	"github.com/openshift/hive/pkg/controller/argocdregister"
	"os"
	"strings"
)

func init() {
	if strings.EqualFold(os.Getenv(hiveconstants.ArgoCDEnvVar), "true") {
		// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
		AddToManagerFuncs = append(AddToManagerFuncs, argocdregister.Add)
	}
}
