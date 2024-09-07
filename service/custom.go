package service

import (
	"cassata/utils"
	"context"
	"errors"
	"strings"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func CreatePointerCRD(pointerGvr *schema.GroupVersionResource, resourceType string) error {
	crd := &apiextensionsv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: pointerGvr.Resource + "." + pointerGvr.Group,
		},
		Spec: apiextensionsv1.CustomResourceDefinitionSpec{
			Group: pointerGvr.Group,
			Names: apiextensionsv1.CustomResourceDefinitionNames{
				Plural:   pointerGvr.Resource,
				Singular: strings.TrimSuffix(pointerGvr.Resource, "s"),
				Kind:     strings.TrimSuffix(resourceType, "s"),
			},
			Scope: apiextensionsv1.NamespaceScoped,
			Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name:    pointerGvr.Version,
					Served:  true,
					Storage: true,
					Schema: &apiextensionsv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]apiextensionsv1.JSONSchemaProps{
								"spec": {
									Type: "object",
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"refs": {
											Type: "array",
											Items: &apiextensionsv1.JSONSchemaPropsOrArray{
												Schema: &apiextensionsv1.JSONSchemaProps{
													Type: "object",
													Properties: map[string]apiextensionsv1.JSONSchemaProps{
														"name":       {Type: "string"},
														"apiVersion": {Type: "string"},
														"kind":       {Type: "string"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	apiextensionsClient, err := utils.GetApiextensionsClientset()
	if err != nil {
		return errors.New("Failed to create apiextensions client: " + err.Error())
	}

	_, err = apiextensionsClient.ApiextensionsV1().CustomResourceDefinitions().Create(context.TODO(), crd, metav1.CreateOptions{})
	if err != nil {
		return errors.New("Failed to create CRD: " + err.Error())
	}
	return nil
}
