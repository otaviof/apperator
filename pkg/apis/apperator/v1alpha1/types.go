package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//
// After inserting new spec fields, run:
//   $ make generate
//
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//

// ApperatorAppSpec defines the desired state of ApperatorApp
type ApperatorAppSpec struct {
	Configuration   ConfigurationSpec   `json:"configuration"`
	Deployment      DeploymentSpec      `json:"deployment"`
	Instrumentation InstrumentationSpec `json:"instrumentation"`
}

// InstrumentationSpec describes the instrumentation type employed
type InstrumentationSpec struct {
	Name   string      `json:"name"`
	Probes []ProbeSpec `json:"probes"`
}

type ProbeSpec struct {
}

// ApperatorAppStatus holds the operator status
type ApperatorAppStatus struct {
}

// ApperatorApp is the Schema for the apperatorapps API
type ApperatorApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApperatorAppSpec   `json:"spec,omitempty"`
	Status ApperatorAppStatus `json:"status,omitempty"`
}

// ApperatorAppList contains a list of ApperatorApp
type ApperatorAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApperatorApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApperatorApp{}, &ApperatorAppList{})
}
