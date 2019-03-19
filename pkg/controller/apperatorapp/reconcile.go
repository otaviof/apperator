package apperatorapp

import (
	"context"
	"reflect"

	logr "github.com/go-logr/logr"
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

	logger.Info("Checking if Deployment Spec is up-to-date", "name", currentDeployment.Name)
	if !reflect.DeepEqual(currentDeployment.Spec, deployment.Spec) {
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
