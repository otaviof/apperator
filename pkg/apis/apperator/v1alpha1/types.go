package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//
// After inserting new spec fields, run:
//   $ make generate
//

// ApperatorAppSpec defines the desired state of ApperatorApp
// +k8s:openapi-gen=true
type ApperatorAppSpec struct {
	Deployment     DeploymentSpec  `json:"deployment"`
	Envs           []EnvSpec       `json:"envs"`
	InitContainers []ContainerSpec `json:"initContainers"`
	Sidecars       []ContainerSpec `json:"sidecars"`
	Probes         []ProbeSpec     `json:"probes"`
	Vault          []VaultSpec     `json:"vault"`
}

// DeploymentSpec describes a core-v1 deployment spec object
// +k8s:openapi-gen=true
type DeploymentSpec struct {
	Spec appsv1.DeploymentSpec `json:"spec"`
}

// EnvSpec wrapper for EnvVar
// +k8s:openapi-gen=true
type EnvSpec struct {
	Env []corev1.EnvVar `json:"env"`
}

// ContainerSpec wrapper for Container
// +k8s:openapi-gen=true
type ContainerSpec struct {
	Spec corev1.Container `json"spec"`
}

// ProbeSpec describes a Prometheus probe
// +k8s:openapi-gen=true
type ProbeSpec struct{}

// VaultSpec describes a Vault ready init-container
// +k8s:openapi-gen=true
type VaultSpec struct{}

// ApperatorAppStatus holds the operator status
// +k8s:openapi-gen=true
type ApperatorAppStatus struct {
}

// ApperatorApp is the Schema for the apperatorapps API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type ApperatorApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApperatorAppSpec   `json:"spec,omitempty"`
	Status ApperatorAppStatus `json:"status,omitempty"`
}

// ApperatorAppList contains a list of ApperatorApp
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ApperatorAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApperatorApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApperatorApp{}, &ApperatorAppList{})
}
