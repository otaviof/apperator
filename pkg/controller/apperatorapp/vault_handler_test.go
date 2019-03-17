package apperatorapp

import (
	"strings"
	"testing"

	v1alpha1 "github.com/otaviof/apperator/pkg/apis/apperator/v1alpha1"
	manifest "github.com/otaviof/vault-handler/pkg/apis/vault-handler/v1alpha1"
	"github.com/stretchr/testify/assert"
)

const (
	testConfigMapName = "config-map"
	testSecretName    = "secret-name"
	testVaultAddr     = "http://127.0.0.1:8200"
)

var apperatorVaultHandler *VaultHandler

func TestVaultHandlerNew(t *testing.T) {
	apperatorVaultObj := &v1alpha1.ApperatorVault{
		Spec: v1alpha1.ApperatorVaultSpec{
			Authorization: v1alpha1.VaultAuthorizationSpec{
				SecretName: testSecretName,
				SecretKeys: v1alpha1.VaultAuthorizationSecretKeysSpec{
					Token:    "token",
					RoleID:   "role-id",
					SecretID: "secret-id",
				},
			},
			Secrets: manifest.Manifest{
				map[string]manifest.Secrets{
					"group": manifest.Secrets{
						Path: "secret/path/in/vault",
						Data: []manifest.SecretData{
							manifest.SecretData{
								Name:          "name",
								Extension:     "txt",
								NameAsSubPath: true,
								Unzip:         true,
							},
						},
					},
				},
			},
		},
	}
	apperatorVaultHandler = NewVaultHandler(apperatorVaultObj, testConfigMapName, testVaultAddr)

	assert.NotNil(t, apperatorVaultHandler)
}

func TestVaultHandlerManifestAsString(t *testing.T) {
	str, err := apperatorVaultHandler.manifestAsString()

	assert.Nil(t, err)
	// new lines are kept, making sure it's using regular yaml payload
	assert.True(t, len(strings.Split(str, "\n")) > 0)
}

func TestVaultHandlerConfigMap(t *testing.T) {
	configMap, err := apperatorVaultHandler.ConfigMap()

	assert.Nil(t, err)
	assert.NotNil(t, configMap)
	assert.True(t, len(configMap.Data["manifest.yaml"]) > 0)
}

func TestVaultHandlerVolumeEntry(t *testing.T) {
	volumes := apperatorVaultHandler.VolumeEntry()

	assert.NotNil(t, volumes)
	assert.True(t, len(volumes) > 0)
}

func TestVaultHandlerGenerateEnvVarSource(t *testing.T) {
	envVarSource := apperatorVaultHandler.generateEnvVarSource("test")

	assert.NotNil(t, envVarSource)
	assert.Equal(t, "test", envVarSource.SecretKeyRef.Key)
}

func TestVaultHandlerPrepareEnv(t *testing.T) {
	envs := apperatorVaultHandler.prepareEnv()

	assert.NotNil(t, envs)
	assert.True(t, len(envs) > 0)
}

func TestVaultHandlerContainer(t *testing.T) {
	container := apperatorVaultHandler.Container()

	assert.NotNil(t, container)
	assert.Equal(t, "vault-handler", container.Name)
	assert.True(t, len(container.Args) > 0)
	assert.True(t, len(container.Env) > 0)
	assert.True(t, len(container.VolumeMounts) > 0)
}
