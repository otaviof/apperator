package apperatorapp

import (
	"testing"

	"github.com/stretchr/testify/assert"

	v1alpha1 "github.com/otaviof/apperator/pkg/apis/apperator/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var deployment *Deployment

// TestNewDeployment setup Deployment object with context. This method setup the Kubernetes related
// infrastructure so we can run the controller code.
func TestNewDeployment(t *testing.T) {
	// setting up verbose logging
	logf.SetLogger(logf.ZapLogger(true))

	var (
		name              = "app"
		namespace         = "apperator"
		matchLabels       = map[string]string{"app": "apperator"}
		replicas    int32 = 1
	)

	apperatorEnvObj := &v1alpha1.ApperatorEnv{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "env",
			Namespace: namespace,
		},
		Spec: v1alpha1.ApperatorEnvSpec{
			Env: []corev1.EnvVar{
				corev1.EnvVar{Name: "ENV_VAR_2", Value: "VALUE"},
			},
		},
	}

	apperatorInitContainerObj := &v1alpha1.ApperatorContainer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "init-container",
			Namespace: namespace,
		},
		Spec: v1alpha1.ApperatorContainerSpec{
			Spec: corev1.Container{
				Name:  "init-container",
				Image: "example/init-container:latest",
			},
		},
	}

	apperatorAppObj := &v1alpha1.ApperatorApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.ApperatorAppSpec{
			Deployment: v1alpha1.DeploymentSpec{
				Spec: appsv1.DeploymentSpec{
					Replicas: &replicas,
					Selector: &metav1.LabelSelector{
						MatchLabels: matchLabels,
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: matchLabels,
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								corev1.Container{
									Name:  "app",
									Image: "example/image:latest",
									Env: []corev1.EnvVar{
										corev1.EnvVar{Name: "ENV_VAR_1", Value: "VALUE"},
									},
								},
							},
						},
					},
				},
			},
			Envs:           []string{"env"},
			InitContainers: []string{"init-container"},
			Probes:         []string{},
			Sidecars:       []string{},
			Vault:          []string{},
		},
	}

	objects := []runtime.Object{
		apperatorEnvObj,
		apperatorInitContainerObj,
		apperatorAppObj,
	}

	s := scheme.Scheme
	s.AddKnownTypes(
		v1alpha1.SchemeGroupVersion,
		apperatorEnvObj,
		apperatorInitContainerObj,
		apperatorAppObj,
	)

	client := fake.NewFakeClient(objects...)
	deployment = NewDeployment(client, apperatorAppObj)

	assert.Equal(t, name, deployment.name)
	assert.Equal(t, namespace, deployment.namespace)
}

func TestDeploymentRender(t *testing.T) {
	_, err := deployment.Render()

	assert.Nil(t, err)
}

func TestDeploymentMergeEnvs(t *testing.T) {
	assert.Equal(t, 2, len(deployment.deployment.Template.Spec.Containers[0].Env))
	assert.Equal(t, "ENV_VAR_1", deployment.deployment.Template.Spec.Containers[0].Env[0].Name)
	assert.Equal(t, "ENV_VAR_2", deployment.deployment.Template.Spec.Containers[0].Env[1].Name)
}

func TestDeploymentMergeInitContainers(t *testing.T) {
	assert.Equal(t, 1, len(deployment.deployment.Template.Spec.InitContainers))
	assert.Equal(t, "init-container", deployment.deployment.Template.Spec.InitContainers[0].Name)
}
