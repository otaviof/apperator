package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ApperatorProbeSpec describes a Prometheus probe
// +k8s:openapi-gen=true
type ApperatorProbeSpec struct{}

// ApperatorProbeStatus defines the observed state of ApperatorProbe
// +k8s:openapi-gen=true
type ApperatorProbeStatus struct{}

// ApperatorProbe is the Schema for the apperatorprobes API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type ApperatorProbe struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApperatorProbeSpec   `json:"spec,omitempty"`
	Status ApperatorProbeStatus `json:"status,omitempty"`
}

// ApperatorProbeList contains a list of ApperatorProbe
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ApperatorProbeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApperatorProbe `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApperatorProbe{}, &ApperatorProbeList{})
}
