package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ApperatorVaultSpec describes a Vault ready init-container
// +k8s:openapi-gen=true
type ApperatorVaultSpec struct{}

// ApperatorVaultStatus defines the observed state of ApperatorVault
// +k8s:openapi-gen=true
type ApperatorVaultStatus struct{}

// ApperatorVault is the Schema for the apperatorvaults API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type ApperatorVault struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApperatorVaultSpec   `json:"spec,omitempty"`
	Status ApperatorVaultStatus `json:"status,omitempty"`
}

// ApperatorVaultList contains a list of ApperatorVault
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ApperatorVaultList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApperatorVault `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApperatorVault{}, &ApperatorVaultList{})
}
