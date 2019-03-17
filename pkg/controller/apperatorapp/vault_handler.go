package apperatorapp

import (
	"fmt"

	v1alpha1 "github.com/otaviof/apperator/pkg/apis/apperator/v1alpha1"
	yaml "gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VaultHandler represents the init-container for vault-handler.
type VaultHandler struct {
	container     *corev1.Container        // shaped as vault-handler init-container
	vault         *v1alpha1.ApperatorVault // CRD
	vaultAddr     string                   // vault api endpoint
	configMapName string                   // name of the config-map with vault-handler manifest
}

const (
	containerName      = "vault-handler"
	secretVolumeName   = "vault-secrets"
	secretVolumePath   = "/vault/secrets"
	manifestVolumeName = "vault-handler-manifest"
	manifestVolumePath = "/vault/manifest"
	manifestName       = "manifest.yaml"
)

func (v *VaultHandler) manifestAsString() (string, error) {
	var payload []byte
	var err error

	if payload, err = yaml.Marshal(v.vault.Spec.Secrets); err != nil {
		return "", err
	}

	return string(payload), nil
}

// ConfigMap with vault-handler extracted manifest.
func (v *VaultHandler) ConfigMap() (*corev1.ConfigMap, error) {
	var manifest string
	var err error

	if manifest, err = v.manifestAsString(); err != nil {
		return nil, err
	}

	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		Data: map[string]string{"manifest.yaml": manifest},
	}

	return configMap, nil
}

// VolumeEntry creates the volume entries to be used
func (v *VaultHandler) VolumeEntry() []corev1.Volume {
	var volumes []corev1.Volume

	// empty-dir to write secrets
	volumes = append(volumes, corev1.Volume{
		Name: secretVolumeName,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	})

	// mounting config-map with vault-handler manifest
	volumes = append(volumes, corev1.Volume{
		Name: manifestVolumeName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: v.configMapName,
				},
			},
		},
	})

	return volumes
}

// generateEnvVarSource part of prepareEnv method, creates the environment variable source reference
// so it's able to read a key from a given kubernetes secret.
func (v *VaultHandler) generateEnvVarSource(key string) *corev1.EnvVarSource {
	return &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: v.vault.Spec.Authorization.SecretName,
			},
			Key: key,
		},
	}
}

// prepareEnv prepare environment with vault-handler configuration variables.
func (v *VaultHandler) prepareEnv() []corev1.EnvVar {
	var envs []corev1.EnvVar

	envs = append(envs, corev1.EnvVar{
		Name:  "VAULT_HANDLER_VAULT_ADDR",
		Value: v.vaultAddr,
	})

	if v.vault.Spec.Authorization.SecretKeys.RoleID != "" {
		envs = append(envs, corev1.EnvVar{
			Name:      "VAULT_HANDLER_VAULT_ROLE_ID",
			ValueFrom: v.generateEnvVarSource(v.vault.Spec.Authorization.SecretKeys.RoleID),
		})
	}

	if v.vault.Spec.Authorization.SecretKeys.SecretID != "" {
		envs = append(envs, corev1.EnvVar{
			Name:      "VAULT_HANDLER_VAULT_SECRET_ID",
			ValueFrom: v.generateEnvVarSource(v.vault.Spec.Authorization.SecretKeys.SecretID),
		})
	}

	if v.vault.Spec.Authorization.SecretKeys.Token != "" {
		envs = append(envs, corev1.EnvVar{
			Name:      "VAULT_HANDLER_VAULT_TOKEN",
			ValueFrom: v.generateEnvVarSource(v.vault.Spec.Authorization.SecretKeys.Token),
		})
	}

	return envs
}

// Container generated for vault-handler.
func (v *VaultHandler) Container() *corev1.Container {
	v.container.Name = containerName
	v.container.Image = "otaviof/vault-handler:latest" // FIXME: make it configurable
	v.container.ImagePullPolicy = "Always"             // FIXME: make it configurable

	v.container.Args = []string{
		"--output-dir",
		secretVolumePath,
		fmt.Sprintf("%s/%s", manifestVolumePath, manifestName),
	}

	v.container.Env = v.prepareEnv()

	v.container.VolumeMounts = []corev1.VolumeMount{
		corev1.VolumeMount{
			Name:      secretVolumeName,
			MountPath: secretVolumePath,
		},
	}

	return v.container
}

// NewVaultHandler creates a new VaultHandler instance.
func NewVaultHandler(vault *v1alpha1.ApperatorVault, vaultAddr, configMapName string) *VaultHandler {
	return &VaultHandler{
		container:     &corev1.Container{},
		vault:         vault,
		vaultAddr:     vaultAddr,
		configMapName: configMapName,
	}
}
