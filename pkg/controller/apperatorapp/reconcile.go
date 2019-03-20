package apperatorapp

import (
	"context"
	"fmt"

	logr "github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/imdario/mergo"
	v1alpha1 "github.com/otaviof/apperator/pkg/apis/apperator/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

// ReconcileApperatorApp reconciles a ApperatorApp object
type ReconcileApperatorApp struct {
	client client.Client
	scheme *runtime.Scheme
}

// getReqLogger creates a logger based the request and current file.
func (r *ReconcileApperatorApp) getReqLogger(req reconcile.Request) logr.Logger {
	return logf.Log.WithName("controller_apperatorapp_reconcile").
		WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)
}

// mergeDeployments merge two deployment objects, where the contetents of source deployment are
// overwritten in the destination. The return is the destination object.
func (r *ReconcileApperatorApp) mergeDeployments(src, dst *appsv1.Deployment) (*appsv1.Deployment, error) {
	if err := mergo.MergeWithOverwrite(dst, src); err != nil {
		return nil, err
	}
	return dst, nil
}

// deploymentEquals compare two deployments, printing out the differences in stdout, while returning
// a boolean to mark if deployments are different.
func (r *ReconcileApperatorApp) deploymentEquals(a, b *appsv1.Deployment) bool {
	if diff := cmp.Diff(a.Spec.Template.Spec, b.Spec.Template.Spec); diff != "" {
		fmt.Printf("---diff---\n%s\n---end-diff---\n", diff)
		return false
	}

	return true
}

// Reconcile reads the state of ApperatorApp and decide further changes.
func (r *ReconcileApperatorApp) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	var logger = r.getReqLogger(req)
	var err error

	logger.Info("Reconciling ApperatorApp...")

	app := &v1alpha1.ApperatorApp{}
	if err = r.client.Get(context.TODO(), req.NamespacedName, app); err != nil {
		if !errors.IsNotFound(err) {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
		},
	}

	controller := NewDeployment(r.client, app)
	if deployment.Spec, err = controller.RenderSpec(); err != nil {
		logger.Error(err, "Failed to render a new deployment")
		return reconcile.Result{}, err
	}

	logger.Info("Looking for deployment object")
	currentDeployment := &appsv1.Deployment{}
	if err = r.client.Get(context.TODO(), req.NamespacedName, currentDeployment); err != nil {
		if !errors.IsNotFound(err) {
			logger.Error(err, "Failed to read deployment")
			return reconcile.Result{}, err
		}

		logger.Info("Creating a new deployment", "name", deployment.Name)
		if err = r.client.Create(context.TODO(), deployment); err != nil {
			logger.Error(err, "Failed to create deployment")
			return reconcile.Result{}, err
		}

		return reconcile.Result{Requeue: true}, nil
	}

	logger.Info("Overwriting values in existing deployment...")
	if deployment, err = r.mergeDeployments(deployment, currentDeployment); err != nil {
		logger.Error(err, "Failed to update deployment")
		return reconcile.Result{}, err
	}

	logger.Info("Checking if Deployment Spec is up-to-date", "name", currentDeployment.Name)
	if !r.deploymentEquals(currentDeployment, deployment) {
		logger.Info("Updating deployment", "name", currentDeployment.Name)
		currentDeployment.Spec = deployment.Spec

		if err = r.client.Update(context.TODO(), currentDeployment); err != nil {
			logger.Error(err, "Failed to update deployment")
			return reconcile.Result{}, err
		}

		return reconcile.Result{Requeue: true}, nil
	}

	logger.Info("Deployment is up-to-date", "name", currentDeployment.Name)
	return reconcile.Result{}, nil
}

// newReconcileApperatorApp creates a new instance setting up logger first.
func newReconcileApperatorApp(client client.Client, scheme *runtime.Scheme) *ReconcileApperatorApp {
	return &ReconcileApperatorApp{client: client, scheme: scheme}
}
