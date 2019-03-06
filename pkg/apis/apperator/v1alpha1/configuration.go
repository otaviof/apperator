package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

// ConfigurationSpec holds the objects that compose app's configuration
type ConfigurationSpec struct {
	Environments []EnvironmentSpec `json:"environments,omitempty"`
}

// EnvironmentSpec wrapper around core-v1 EnvVar
type EnvironmentSpec struct {
	Env []corev1.EnvVar `json:"env"`
}
