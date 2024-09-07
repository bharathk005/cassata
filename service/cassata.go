package service

import (
	"cassata/repository"
	"cassata/utils"
	"context"
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CreateResource(provider string, resourceType string, namespace string, obj map[string]interface{}) (interface{}, error) {
	// 1. Get GVR for provider and resource type
	gvr, err := repository.GetGvrForProviderResourceType(provider, resourceType)
	if err != nil {
		return nil, err
	}

	// 2. Create the resource based on the input payload

	if obj["spec"] == nil {
		return nil, errors.New("spec is required")
	}

	obj["spec"] = map[string]interface{}{
		"forProvider":    obj["spec"],
		"deletionPolicy": "Delete",
		"providerConfigRef": map[string]interface{}{
			"name": "gcp-provider",
		},
	}

	resource := unstructured.Unstructured{Object: obj}

	// Add necessary Kubernetes resource attributes
	resource.SetAPIVersion(gvr.GroupVersion().String())
	resource.SetKind(resourceType) // Assuming singular form

	// Ensure metadata is properly set
	metadata := resource.Object["metadata"].(map[string]interface{})
	if metadata == nil || metadata["name"] == nil {
		return nil, errors.New("metadata.name is required")
	}
	metadata["namespace"] = namespace
	resource.Object["metadata"] = metadata

	// Get dynamic client
	dynamicClient, err := utils.GetDynamicClient()
	if err != nil {
		return nil, errors.New("Failed to create dynamic client: " + err.Error())
	}

	// Create the resource
	created, err := dynamicClient.Resource(gvr).Create(context.TODO(), &resource, metav1.CreateOptions{})
	if err != nil {
		// 3. If any error has occurred, return a detailed error message
		return nil, err
	}

	err = CreatePointer(gvr, provider, resourceType, namespace, metadata["name"].(string))
	if err != nil {
		// revert the resource creation
		fmt.Println("Error creating pointer, reverting resource creation")
		dynamicClient.Resource(gvr).Delete(context.TODO(), metadata["name"].(string), metav1.DeleteOptions{})
		return nil, err
	}

	// 4. If the resource is created successfully, return a success message
	return created, nil
}

func GetResource(provider string, resourceType string, namespace string, resourceID string) (interface{}, error) {
	// 1. Get GVR for provider and resource type
	gvr, err := repository.GetGvrForProviderResourceType(provider, resourceType)
	if err != nil {
		return nil, err
	}

	// 2. Get the resource from the Kubernetes API
	dynamicClient, err := utils.GetDynamicClient()
	if err != nil {
		return nil, errors.New("Failed to create dynamic client: " + err.Error())
	}

	resource, err := dynamicClient.Resource(gvr).Get(context.TODO(), resourceID, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return resource, nil
}
