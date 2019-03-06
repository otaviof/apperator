// +build !ignore_autogenerated

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/otaviof/apperator/pkg/apis/apperator/v1alpha1.ApperatorApp":       schema_pkg_apis_apperator_v1alpha1_ApperatorApp(ref),
		"github.com/otaviof/apperator/pkg/apis/apperator/v1alpha1.ApperatorAppSpec":   schema_pkg_apis_apperator_v1alpha1_ApperatorAppSpec(ref),
		"github.com/otaviof/apperator/pkg/apis/apperator/v1alpha1.ApperatorAppStatus": schema_pkg_apis_apperator_v1alpha1_ApperatorAppStatus(ref),
	}
}

func schema_pkg_apis_apperator_v1alpha1_ApperatorApp(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ApperatorApp is the Schema for the apperatorapps API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/otaviof/apperator/pkg/apis/apperator/v1alpha1.ApperatorAppSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/otaviof/apperator/pkg/apis/apperator/v1alpha1.ApperatorAppStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/otaviof/apperator/pkg/apis/apperator/v1alpha1.ApperatorAppSpec", "github.com/otaviof/apperator/pkg/apis/apperator/v1alpha1.ApperatorAppStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_apperator_v1alpha1_ApperatorAppSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ApperatorAppSpec defines the desired state of ApperatorApp",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_apperator_v1alpha1_ApperatorAppStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ApperatorAppStatus defines the observed state of ApperatorApp",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}