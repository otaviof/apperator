package apperatorapp

import (
	v1alpha1 "github.com/otaviof/apperator/pkg/apis/apperator/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_apperatorapp")

// Add creates a new ApperatorApp Controller and adds it to the Manager. The Manager will set fields
// on the Controller and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler.
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return newReconcileApperatorApp(mgr.GetClient(), mgr.GetScheme())
}

// add new controller manager, will start watching for objects this controller is interested on.
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("apperatorapp-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	if err = c.Watch(
		&source.Kind{Type: &v1alpha1.ApperatorApp{}},
		&handler.EnqueueRequestForObject{},
	); err != nil {
		return err
	}

	if err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.ApperatorApp{},
	}); err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileApperatorApp{}
