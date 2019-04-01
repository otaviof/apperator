package apperatorapp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var vaultHandler *VaultHandler

func TestVaultHandlerNew(t *testing.T) {
	var err error

	vaultHandlerConfigMapObj := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "vault-handler-manifest-configmap",
			Namespace: "namespace",
		},
		Data: map[string]string{
			"authorization": `
secretName: secret-name
secretKeys:
  roleId: role-id
  secretId: secret-id`,
			"secrets": `
group:
  path: secret/path/in/vault
  data:
    - name: name
      extension: txt
      zip: true
      nameAsSubPath: true`,
		},
	}

	vaultHandler, err = NewVaultHandler(vaultHandlerConfigMapObj, "http://127.0.0.1:8200")

	assert.Nil(t, err)
	assert.NotNil(t, vaultHandler)

	assert.Equal(t, "secret-name", vaultHandler.auth.SecretName)
	assert.Equal(t, "role-id", vaultHandler.auth.SecretKeys.RoleID)
	assert.Equal(t, "secret-id", vaultHandler.auth.SecretKeys.SecretID)

	assert.True(t, len(vaultHandler.manifest) > 0)
	assert.Equal(t, "secret/path/in/vault", vaultHandler.manifest["group"].Path)
	assert.True(t, len(vaultHandler.manifest["group"].Data) > 0)
	assert.Equal(t, "name", vaultHandler.manifest["group"].Data[0].Name)
	assert.Equal(t, "txt", vaultHandler.manifest["group"].Data[0].Extension)
	assert.Equal(t, true, vaultHandler.manifest["group"].Data[0].Zip)
	assert.Equal(t, true, vaultHandler.manifest["group"].Data[0].NameAsSubPath)
}

func TestVaultHandlerVolumeEntry(t *testing.T) {
	volumes := vaultHandler.VolumeEntry()

	assert.NotNil(t, volumes)
	assert.True(t, len(volumes) > 0)

	assert.Equal(t, "vault-handler-manifest", volumes[0].Name)
	assert.Equal(t, "vault-secrets", volumes[1].Name)

	assert.Equal(t, "vault-handler-manifest-configmap", volumes[0].ConfigMap.LocalObjectReference.Name)
	assert.NotNil(t, volumes[1].EmptyDir)
}

func TestVaultHandlerGenerateEnvVarSource(t *testing.T) {
	envVarSource := vaultHandler.generateEnvVarSource("test")

	assert.NotNil(t, envVarSource)
	assert.Equal(t, "test", envVarSource.SecretKeyRef.Key)
}

func TestVaultHandlerPrepareEnv(t *testing.T) {
	envs := vaultHandler.prepareEnv()

	assert.NotNil(t, envs)
	assert.True(t, len(envs) > 0)
}

func TestVaultHandlerContainer(t *testing.T) {
	container := vaultHandler.Container()

	assert.NotNil(t, container)
	assert.Equal(t, "vault-handler", container.Name)
	assert.True(t, len(container.Args) > 0)
	assert.True(t, len(container.Env) > 0)
	assert.True(t, len(container.VolumeMounts) > 0)
}
