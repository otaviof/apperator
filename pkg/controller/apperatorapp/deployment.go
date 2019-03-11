package apperatorapp

import (
	"context"
	"errors"

	logr "github.com/go-logr/logr"
	v1alpha1 "github.com/otaviof/apperator/pkg/apis/apperator/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

// Deployment represents the desired deployment object produced based on ApperatorApp.
type Deployment struct {
	client     client.Client
	deployment *appsv1.DeploymentSpec
	app        *v1alpha1.ApperatorApp
	name       string
	namespace  string
	log        logr.Logger
}

// mergeEnvs retrieve ApperatorEnv objects and merge on current deployment.
func (d *Deployment) mergeEnvs() error {
	var err error

	d.log.Info("Retrieving ApperatorEnv objects...")
	for _, envName := range d.app.Spec.Envs {
		d.log.Info("Loading ApperatorEnv object", "ApperatorEnv", envName)

		env := &v1alpha1.ApperatorEnv{}
		meta := types.NamespacedName{Name: envName, Namespace: d.namespace}

		if err = d.client.Get(context.TODO(), meta, env); err != nil {
			d.log.Error(err, "Failed to read ApperatorEnv '%s'", envName)
			return err
		}

		for _, envVar := range env.Spec.Env {
			d.log.Info("Appending env variable", "name", envVar.Name)

			d.deployment.Template.Spec.Containers[0].Env = append(
				d.deployment.Template.Spec.Containers[0].Env,
				envVar,
			)
		}
	}

	return nil
}

// Render a deployment object based on ApperatorApp.
func (d *Deployment) Render() (*Deployment, error) {
	var err error

	if len(d.deployment.Template.Spec.Containers) != 1 {
		err = errors.New("invalid amount of containers in deployment, must be only one")
		d.log.Error(err, "Invalid amount of containers '%d'!", len(d.deployment.Template.Spec.Containers))
		return nil, err
	}

	if err = d.mergeEnvs(); err != nil {
		return nil, err
	}

	return nil, nil
}

// NewDeployment setup request logging and base Deployment instance.
func NewDeployment(client client.Client, app *v1alpha1.ApperatorApp) *Deployment {
	name := app.ObjectMeta.Name
	namespace := app.ObjectMeta.Namespace
	log := logf.Log.
		WithName("controller_deployment").
		WithValues("Request.Namespace", namespace, "Request.Name", name)

	return &Deployment{
		client:     client,
		deployment: &app.Spec.Deployment.Spec,
		app:        app,
		name:       name,
		namespace:  namespace,
		log:        log,
	}
}
