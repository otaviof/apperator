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
// +k8s:openapi-gen=true

// ApperatorAppSpec defines the desired state of ApperatorApp
type ApperatorAppSpec struct {
	Configuration ConfigurationSpec `json:"configuration"`
	Deployment    DeploymentSpec    `json:"deployment"`
}

// ConfigurationSpec holds the objects that compose app's configuration
type ConfigurationSpec struct {
	Environments []EnvironmentSpec `json:"environments,omitempty"`
}

// DeploymentSpec describes a core-v1 deployment spec object
type DeploymentSpec struct {
	Spec appsv1.DeploymentSpec `json:"spec"`
}

// EnvironmentSpec wrapper around core-v1 EnvVar
type EnvironmentSpec struct {
	Env []corev1.EnvVar `json:"env"`
}

type ApperatorAppStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApperatorApp is the Schema for the apperatorapps API
// +k8s:openapi-gen=true
type ApperatorApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApperatorAppSpec   `json:"spec,omitempty"`
	Status ApperatorAppStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApperatorAppList contains a list of ApperatorApp
type ApperatorAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApperatorApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApperatorApp{}, &ApperatorAppList{})
}
