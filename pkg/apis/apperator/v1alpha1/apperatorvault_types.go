package v1alpha1

import (
	manifest "github.com/otaviof/vault-handler/pkg/apis/vault-handler/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ApperatorVaultSpec contains the manifest for vault-handler.
// +k8s:openapi-gen=true
type ApperatorVaultSpec struct {
	Authorization VaultAuthorizationSpec `json:"authorization"`
	Secrets       manifest.Manifest      `json:"secrets"`
}

// VaultAuthorizationSpec configuration to load vault authorization items.
// +k8s:openapi-gen=true
type VaultAuthorizationSpec struct {
	SecretName string                           `json:"secretName"` // kubernetes secret name
	SecretKeys VaultAuthorizationSecretKeysSpec `json:"secretKeys"` // kubernetes secret key mapping
}

// VaultAuthorizationSecretKeysSpec keys names in kubernetes secrets.
// +k8s:openapi-gen=true
type VaultAuthorizationSecretKeysSpec struct {
	RoleID   string `json:"roleId,omitempty"`   // key name for role-id
	SecretID string `json:"secretId,omitempty"` // key name for secret-id
	Token    string `json:"token,omitempty"`    // key name for token
}

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
