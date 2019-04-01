package apperatorapp

import (
	"context"
	"fmt"
	"testing"

	v1alpha1 "github.com/otaviof/apperator/pkg/apis/apperator/v1alpha1"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func TestReconcile(t *testing.T) {
	var name = "apperator-app"
	var image = "example/image:latest"
	var namespace = "apperator"
	var replicas int32 = 1
	var matchLabels = map[string]string{"app": "apperator"}

	logf.SetLogger(logf.ZapLogger(true))

	deploymentSpecObj := appsv1.DeploymentSpec{
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
					{
						Name:  name,
						Image: image,
					},
				},
			},
		},
	}

	apperatorAppObj := &v1alpha1.ApperatorApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.ApperatorAppSpec{
			Deployment: v1alpha1.ApperatorAppDeploymentSpec{
				Spec: deploymentSpecObj,
			},
		},
	}

	objects := []runtime.Object{apperatorAppObj}
	s := scheme.Scheme
	s.AddKnownTypes(v1alpha1.SchemeGroupVersion, apperatorAppObj)
	client := fake.NewFakeClient(objects...)

	r := &ReconcileApperatorApp{client: client, scheme: s}
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}
	res, err := r.Reconcile(req)

	assert.Nil(t, err)
	assert.NotNil(t, res)

	// making sure response is requeue, since it just created a new deployment
	assert.True(t, res.Requeue)

	// fetching deployment created by operator
	deployment := appsv1.Deployment{}
	meta := types.NamespacedName{Name: name, Namespace: namespace}
	err = client.Get(context.TODO(), meta, &deployment)
	assert.Nil(t, err)

	// produced deployment has amount of replicas set, name and image
	assert.Equal(t, &replicas, deployment.Spec.Replicas)
	assert.Equal(t, name, deployment.Spec.Template.Spec.Containers[0].Name)
	assert.Equal(t, image, deployment.Spec.Template.Spec.Containers[0].Image)

	// showing deployment as yaml
	yamlBytes, _ := yaml.Marshal(deployment)
	fmt.Printf("\n%s\n", string(yamlBytes))

	// running reconciliation again
	res, err = r.Reconcile(req)

	assert.Nil(t, err)
	assert.NotNil(t, res)

	// expecting deployment to be fully up-to-date
	assert.False(t, res.Requeue)
}
