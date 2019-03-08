package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//
// After inserting new spec fields, run:
//   $ make generate
//

// ApperatorAppSpec defines the desired state of ApperatorApp
// +k8s:openapi-gen=true
type ApperatorAppSpec struct {
	Deployment DeploymentSpec `json:"deployment"`
}

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
