package service

import (
	"cassata/utils"
	"context"
	"errors"
	"fmt"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func CreatePointer(gvr schema.GroupVersionResource, provider string, resourceType string, namespace string, resourceID string) error {

	// 2. Create the pointer resource
	pointerGvr := schema.GroupVersionResource{
		Group:    "pointer." + gvr.Group,
		Version:  gvr.Version,
		Resource: gvr.Resource,
	}

	pointerObj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": pointerGvr.GroupVersion().String(),
			"kind":       resourceType,
			"metadata": map[string]interface{}{
				"name":      resourceID,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"refs": []map[string]interface{}{
					{
						"name":       resourceID,
						"apiVersion": gvr.GroupVersion().String(),
						"kind":       resourceType,
					},
				},
			},
		},
	}

	// Get dynamic client
	dynamicClient, err := utils.GetDynamicClient()
	if err != nil {
		return errors.New("Failed to create dynamic client: " + err.Error())
	}

	// Check if the CRD exists
	apiextensionsClient, err := utils.GetApiextensionsClientset()
	if err != nil {
		return fmt.Errorf("failed to create apiextensions client: %w", err)
	}

	_, err = apiextensionsClient.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), pointerGvr.Resource+"."+pointerGvr.Group, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// If the resource is not found, attempt to create it
			if createErr := CreatePointerCRD(&pointerGvr, resourceType); createErr != nil {
				return fmt.Errorf("failed to create CRD: %w", createErr)
			}
		} else {
			// Return other errors directly
			return fmt.Errorf("failed to get CRD: %w", err)
		}
	}

	// Create the pointer resource
	_, err = dynamicClient.Resource(pointerGvr).Namespace(namespace).Create(context.TODO(), pointerObj, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil

}
