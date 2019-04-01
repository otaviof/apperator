package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ApperatorContainerSpec defines the desired state of ApperatorContainer
// +k8s:openapi-gen=true
type ApperatorContainerSpec struct {
	Spec corev1.Container `json:"spec"`
}

// ApperatorContainerStatus defines the observed state of ApperatorContainer
// +k8s:openapi-gen=true
type ApperatorContainerStatus struct{}

// ApperatorContainer is the Schema for the apperatorcontainers API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type ApperatorContainer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApperatorContainerSpec   `json:"spec,omitempty"`
	Status ApperatorContainerStatus `json:"status,omitempty"`
}

// ApperatorContainerList contains a list of ApperatorContainer
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ApperatorContainerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApperatorContainer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApperatorContainer{}, &ApperatorContainerList{})
}
