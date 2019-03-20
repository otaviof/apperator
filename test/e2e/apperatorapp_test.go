package e2e

import (
	"context"
	"testing"
	"time"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	apis "github.com/otaviof/apperator/pkg/apis"
	"github.com/otaviof/apperator/pkg/apis/apperator/v1alpha1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	retryInterval        = time.Second * 5
	timeout              = time.Second * 45
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

const (
	name = "apperator-app"
)

// TestApperatorAppAddToScheme defines the scheme and CRDs on the Kubernetes cluster.
func TestApperatorAppAddToScheme(t *testing.T) {
	var err error

	apperatorAppList := &v1alpha1.ApperatorAppList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ApperatorApp",
			APIVersion: "apperator.otaviof.github.io/v1alpha1",
		},
	}

	if err = framework.AddToFrameworkScheme(apis.AddToScheme, apperatorAppList); err != nil {
		t.Fatalf("failed trying to add ApperatorApp to framework scheme: '%#v'", err)
	}

	t.Run("end-to-end", func(t *testing.T) {
		t.Run("s1", ApperatorApp)
	})
}

// apperatorAppTest runs the actual tests against a functional workload.
func apperatorAppTest(t *testing.T, namespace string, f *framework.Framework, ctx *framework.TestCtx) {
	var name = "app"
	var replicas int32 = 1
	var matchLabels = map[string]string{"app": "apperator"}
	var err error

	t.Log("Starting Operator Tests...")

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
					corev1.Container{
						Name:  name,
						Image: "busybox:latest",
						Args:  []string{"sleep", "60"},
					},
				},
			},
		},
	}

	apperatorAppObj := &v1alpha1.ApperatorApp{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ApperatorApp",
			APIVersion: "apperator.otaviof.github.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.ApperatorAppSpec{
			Deployment: v1alpha1.ApperatorAppDeploymentSpec{
				Spec: deploymentSpecObj,
			},
			Envs:           []string{},
			InitContainers: []string{},
			Probes:         []string{},
			Sidecars:       []string{},
			Vault:          []string{},
		},
	}

	t.Logf("Creating apperator object '%s' (namespace: '%s')", name, namespace)
	// creating apperator object
	err = f.Client.Create(context.TODO(), apperatorAppObj, &framework.CleanupOptions{
		TestContext:   ctx,
		RetryInterval: retryInterval,
		Timeout:       timeout,
	})
	// should not return errors when creating the apperator object
	assert.Nil(t, err)

	t.Log("Waiting for deployment produced from ApperatorApp...")
	// waiting for deployment to reach one replica
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, name, 1, retryInterval, timeout)
	// should not have errors on waiting for deployment, the deployment must be successful
	assert.Nil(t, err)
}

// ApperatorApp creates all the underlying things to run tests against a Kuberentes cluster.
func ApperatorApp(t *testing.T) {
	var namespace string
	var err error

	t.Parallel()

	t.Log("Creating text context...")
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()

	t.Log("Initializing cluster resources...")
	if err = ctx.InitializeClusterResources(&framework.CleanupOptions{
		TestContext:   ctx,
		Timeout:       cleanupTimeout,
		RetryInterval: cleanupRetryInterval,
	}); err != nil {
		if !errors.IsAlreadyExists(err) {
			t.Fatalf("failed to setup cluster resources: '%#v'", err)
		}
	}

	t.Log("Acquiring namespace...")
	if namespace, err = ctx.GetNamespace(); err != nil {
		t.Fatalf("failed to get namespace: '%#v'", err)
	}
	t.Logf("Using namespace '%s'", namespace)

	t.Log("Copying global variables from Framework...")
	f := framework.Global

	t.Log("Waiting for ApperatorApp (operator) to be deployed...")
	if err = e2eutil.WaitForOperatorDeployment(
		t, f.KubeClient, namespace, "apperatorapp-controller", 1, retryInterval, timeout,
	); err != nil {
		t.Fatalf("failed on waiting for deployment: '%#v'", err)
	}

	apperatorAppTest(t, namespace, f, ctx)
}
