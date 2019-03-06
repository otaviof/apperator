package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
)

// DeploymentSpec describes a core-v1 deployment spec object
type DeploymentSpec struct {
	Spec appsv1.DeploymentSpec `json:"spec"`
}
