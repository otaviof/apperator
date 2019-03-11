package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//
// After inserting new spec fields, run:
//   $ make generate
//

// ApperatorAppSpec defines the desired state of ApperatorApp
// +k8s:openapi-gen=true
type ApperatorAppSpec struct {
	Deployment     DeploymentSpec `json:"deployment"`
	Envs           []string       `json:"envs"`
	InitContainers []string       `json:"initContainers"`
	Sidecars       []string       `json:"sidecars"`
	Probes         []string       `json:"probes"`
	Vault          []string       `json:"vault"`
}

// DeploymentSpec describes a core-v1 deployment spec object
// +k8s:openapi-gen=true
type DeploymentSpec struct {
	Spec appsv1.DeploymentSpec `json:"spec"`
}

// ApperatorAppStatus holds the operator status
// +k8s:openapi-gen=true
type ApperatorAppStatus struct{}

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
