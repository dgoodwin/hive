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
package argocdregister

import (
	"context"
	//	"errors"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/openshift/hive/pkg/apis"
	controllerutils "github.com/openshift/hive/pkg/controller/utils"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"

	//"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	configv1 "github.com/openshift/api/config/v1"
	hivev1 "github.com/openshift/hive/pkg/apis/hive/v1"
	hivev1aws "github.com/openshift/hive/pkg/apis/hive/v1/aws"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	testName                        = "foo-unreachable"
	testClusterName                 = "bar"
	testClusterID                   = "testFooClusterUUID"
	testNamespace                   = "default"
	pullSecretSecret                = "pull-secret"
	testRemoteClusterCurrentVersion = "4.0.0"
	remoteClusterVersionObjectName  = "version"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestUnreachablClusterStatusCondition(t *testing.T) {
	apis.AddToScheme(scheme.Scheme)
	configv1.Install(scheme.Scheme) //remove later
	tests := []struct {
		name        string
		existing    []runtime.Object
		expectError bool
		validate    func(*testing.T, *hivev1.ClusterDeployment)
	}{
		{
			name: "unreachable condition should be added",
			existing: []runtime.Object{
				getClusterDeployment(),
				getKubeconfigSecret(),
			},
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond := controllerutils.FindClusterDeploymentCondition(cd.Status.Conditions, hivev1.UnreachableCondition)
				if cond.Status != corev1.ConditionTrue {
					t.Errorf("Did not get expected state for unreachable condition. Expected: \n%#v\nGot: \n%#v", corev1.ConditionTrue, cond.Status)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fakeClient := fake.NewFakeClient(test.existing...)
			rcd := &ArgoCDRegisterController{
				Client: fakeClient,
				scheme: scheme.Scheme,
				logger: log.WithField("controller", "unreachable"),
			}

			namespacedName := types.NamespacedName{
				Name:      testName,
				Namespace: testNamespace,
			}

			_, err := rcd.Reconcile(reconcile.Request{NamespacedName: namespacedName})

			if test.validate != nil {
				cd := &hivev1.ClusterDeployment{}
				err := fakeClient.Get(context.TODO(), namespacedName, cd)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				test.validate(t, cd)
			}

			if err != nil && !test.expectError {
				t.Errorf("Unexpected error: %v", err)
			}
			if err == nil && test.expectError {
				t.Errorf("Expected error but got none")
			}
		})
	}
}

func getClusterDeployment() *hivev1.ClusterDeployment {
	cd := &hivev1.ClusterDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:       testName,
			Namespace:  testNamespace,
			Finalizers: []string{hivev1.FinalizerDeprovision},
			UID:        types.UID("1234"),
		},
		Spec: hivev1.ClusterDeploymentSpec{
			ClusterName: testClusterName,
			Compute:     []hivev1.MachinePool{},
			PullSecretRef: &corev1.LocalObjectReference{
				Name: pullSecretSecret,
			},
			Platform: hivev1.Platform{
				AWS: &hivev1aws.Platform{
					CredentialsSecretRef: corev1.LocalObjectReference{
						Name: "aws-credentials",
					},
					Region: "us-east-1",
				},
			},
			ClusterMetadata: &hivev1.ClusterMetadata{
				ClusterID: testClusterID,
				AdminKubeconfigSecretRef: corev1.LocalObjectReference{
					Name: "kubeconfig-secret",
				},
			},
			Installed: true,
		},
	}
	return cd
}

func getKubeconfigSecret() *corev1.Secret {
	return testSecret("kubeconfig-secret", "kubeconfig", "KUBECONFIG-DATA")
}

func testSecret(name, key, value string) *corev1.Secret {
	s := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
		Data: map[string][]byte{
			key: []byte(value),
		},
	}
	return s
}
