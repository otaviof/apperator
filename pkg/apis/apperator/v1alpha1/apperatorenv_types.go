package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ApperatorEnvSpec wrapper for EnvVar
// +k8s:openapi-gen=true
type ApperatorEnvSpec struct {
	Env []corev1.EnvVar `json:"env"`
}

// ApperatorEnvStatus defines the observed state of ApperatorEnv
// +k8s:openapi-gen=true
type ApperatorEnvStatus struct{}

// ApperatorEnv is the Schema for the apperatorenvs API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type ApperatorEnv struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApperatorEnvSpec   `json:"spec,omitempty"`
	Status ApperatorEnvStatus `json:"status,omitempty"`
}

// ApperatorEnvList contains a list of ApperatorEnv
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ApperatorEnvList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApperatorEnv `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApperatorEnv{}, &ApperatorEnvList{})
}
