#!/bin/sh

set -e

# Good candidate for doing in the operator.
kubectl apply -f config/v1alpha1-crds-tmp-apigroup


# Copy all Hive CRs to a new CRD in a temporary APIgroup so we can reclaim our main one and
# add in aggregated APIServer for v1alpha1.
oc get cd -A -o json | jq '.items | .[] |
	del(.status) |
	del(.metadata.annotations."kubectl.kubernetes.io/last-applied-configuration") |
	del(.metadata.creationTimestamp) |
	del(.metadata.generation) |
	del(.metadata.resourceVersion) |
	del(.metadata.selfLink) |
	del(.metadata.uid) |
	.apiVersion="temphive.openshift.io/v1alpha1"'

# TODO: do we need to migrate checkpoints?
# TODO: do we need to migrate deprovisionrequests?

oc get clusterimagesets -A -o json | jq '.items | .[] |
	del(.status) |
	del(.metadata.annotations."kubectl.kubernetes.io/last-applied-configuration") |
	del(.metadata.creationTimestamp) |
	del(.metadata.generation) |
	del(.metadata.resourceVersion) |
	del(.metadata.selfLink) |
	del(.metadata.uid) |
	.apiVersion="temphive.openshift.io/v1alpha1"'

# TODO: secret owner refs
