package e2e

import (
	"testing"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	apis "github.com/otaviof/apperator/pkg/apis"
	operator "github.com/otaviof/apperator/pkg/apis/apperator/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestApperatorAppAddToScheme(t *testing.T) {
	var err error

	apperatorAppList := &operator.ApperatorAppList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ApperatorApp",
			APIVersion: "apperator.otaviof.github.io/v1alpha1",
		},
	}

	if err = framework.AddToFrameworkScheme(apis.AddToScheme, apperatorAppList); err != nil {
		t.Fatalf("failed trying to add ApperatorApp to framework scheme: '%#v'", err)
	}
}
