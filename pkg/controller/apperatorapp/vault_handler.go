package apperatorapp

import (
	"fmt"

	vaulthandler "github.com/otaviof/vault-handler/pkg/vault-handler"
	yaml "gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
)

// VaultHandler represents the init-container for vault-handler.
type VaultHandler struct {
	configMap *corev1.ConfigMap               // config-map with vault-handler construction
	vaultAddr string                          // vault api endpoint
	auth      *VaultHandlerAuth               // vault-handler authentication configuration
	manifest  map[string]vaulthandler.Secrets // vault-handler manifest
	container corev1.Container                // shaped as vault-handler init-container
}

// VaultHandlerAuth section of config-map to define how vault-handler will authenticate.
type VaultHandlerAuth struct {
	SecretName string                     `yaml:"secretName"` // kubernetes secret name
	SecretKeys VaultHandlerAuthSecretKeys `yaml:"secretKeys"` // kubernetes secret key mapping
}

// VaultHandlerAuthSecretKeys keys to map a secret to authentication config.
type VaultHandlerAuthSecretKeys struct {
	RoleID   string `yaml:"roleId"`   // key name for role-id
	SecretID string `yaml:"secretId"` // key name for secret-id
	Token    string `yaml:"token"`    // key name for token
}

const (
	containerName      = "vault-handler"
	secretVolumeName   = "vault-secrets"
	secretVolumePath   = "/vault/secrets"
	manifestVolumeName = "vault-handler-manifest"
	manifestVolumePath = "/vault/manifest"
	manifestName       = "manifest.yaml"
)

// VolumeEntry creates the volume entries to be used
func (v *VaultHandler) VolumeEntry() []corev1.Volume {
	var volumes []corev1.Volume

	// mounting config-map with vault-handler manifest
	volumes = append(volumes, corev1.Volume{
		Name: manifestVolumeName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: v.configMap.ObjectMeta.Name,
				},
			},
		},
	})

	// empty-dir to write secrets
	volumes = append(volumes, corev1.Volume{
		Name: secretVolumeName,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
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
				Name: v.auth.SecretName,
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

	if v.auth.SecretKeys.RoleID != "" {
		envs = append(envs, corev1.EnvVar{
			Name:      "VAULT_HANDLER_VAULT_ROLE_ID",
			ValueFrom: v.generateEnvVarSource(v.auth.SecretKeys.RoleID),
		})
	}

	if v.auth.SecretKeys.SecretID != "" {
		envs = append(envs, corev1.EnvVar{
			Name:      "VAULT_HANDLER_VAULT_SECRET_ID",
			ValueFrom: v.generateEnvVarSource(v.auth.SecretKeys.SecretID),
		})
	}

	if v.auth.SecretKeys.Token != "" {
		envs = append(envs, corev1.EnvVar{
			Name:      "VAULT_HANDLER_VAULT_TOKEN",
			ValueFrom: v.generateEnvVarSource(v.auth.SecretKeys.Token),
		})
	}

	return envs
}

// Container generated for vault-handler.
func (v *VaultHandler) Container() corev1.Container {
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
		{
			Name:      manifestVolumeName,
			MountPath: manifestVolumePath,
		},
		{
			Name:      secretVolumeName,
			MountPath: secretVolumePath,
		},
	}

	return v.container
}

// parseConfigMapData look for expected keys and parse them as yaml.
func (v *VaultHandler) parseConfigMapData() error {
	var found bool
	var authStr string
	var secretsStr string
	var err error

	if authStr, found = v.configMap.Data["authorization"]; !found {
		return fmt.Errorf("'authorization' section is not found in config-map")
	}
	if secretsStr, found = v.configMap.Data["secrets"]; !found {
		return fmt.Errorf("'secrets' section is not found in config-map")
	}

	if err = yaml.Unmarshal([]byte(authStr), &v.auth); err != nil {
		return err
	}
	return yaml.Unmarshal([]byte(secretsStr), &v.manifest)
}

// NewVaultHandler create a new instance by parsing config-map contents.
func NewVaultHandler(configMap *corev1.ConfigMap, vaultAddr string) (*VaultHandler, error) {
	vaultHandler := &VaultHandler{
		configMap: configMap,
		vaultAddr: vaultAddr,
		container: corev1.Container{},
	}

	if err := vaultHandler.parseConfigMapData(); err != nil {
		return nil, err
	}

	return vaultHandler, nil
}
