package apperatorapp

import (
	"context"

	"k8s.io/apimachinery/pkg/types"

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
	log    logr.Logger
}

// Reconcile reads the state of ApperatorApp and decide further changes.
func (r *ReconcileApperatorApp) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	var err error

	reqLogger := r.log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)
	reqLogger.Info("Reconciling ApperatorApp...")

	app := &v1alpha1.ApperatorApp{}
	if err = r.client.Get(context.TODO(), req.NamespacedName, app); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	deployment := &appsv1.Deployment{}
	if err = r.client.Get(context.TODO(), req.NamespacedName, deployment); err != nil {
		if errors.IsNotFound(err) {
			deployment = r.createDeployment(app)
		} else {
			reqLogger.Error(err, "Failed to read deployment '%s'", app.ObjectMeta.Name)
		}
	}

	return reconcile.Result{}, nil
}

// createDeployment puts together a deployment based in the informed instance details.
func (r *ReconcileApperatorApp) createDeployment(app *v1alpha1.ApperatorApp) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta:   metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{Name: app.Name, Namespace: app.Namespace},
		Spec:       app.Spec.Deployment.Spec,
	}
}

// mergeEnvs execute the merge of upstream ApperatorEnv objects into current deployment.
func (r *ReconcileApperatorApp) mergeEnvs(envNames []string, deployment *appsv1.Deployment) error {
	var err error

	log.Info("Merging ApperatorEnv objects ('%#v') into deployment.", envNames)
	if len(deployment.Spec.Template.Spec.Containers) == 0 {
		log.Info("Warning! Deployment '%s' does not define any containers!", deployment.ObjectMeta.Name)
		return nil
	}

	for _, envName := range envNames {
		env := &v1alpha1.ApperatorEnv{}
		meta := types.NamespacedName{Name: envName, Namespace: deployment.ObjectMeta.Namespace}

		if err = r.client.Get(context.TODO(), meta, env); err != nil {
			return err
		}

		for _, envVar := range env.Spec.Env {
			deployment.Spec.Template.Spec.Containers[0].Env = append(
				deployment.Spec.Template.Spec.Containers[0].Env,
				envVar,
			)
		}
	}

	return nil
}

// newReconcileApperatorApp creates a new instance setting up logger first.
func newReconcileApperatorApp() *ReconcileApperatorApp {
	return &ReconcileApperatorApp{
		log: logf.Log.WithName("controller_apperatorapp_reconcile"),
	}
}
