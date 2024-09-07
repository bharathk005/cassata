package repository

import (
	"cassata/models"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GetResourceMap(provider string, resourceType string) (models.ResourceMap, error) {
	db := GetDB()
	var resourceMap models.ResourceMap
	err := db.Where("provider = ? AND resource_type = ?", provider, resourceType).First(&resourceMap).Error
	return resourceMap, err
}

func GetGvrForProviderResourceType(provider string, resourceType string) (schema.GroupVersionResource, error) {
	resourceMap, err := GetResourceMap(provider, resourceType)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	return schema.GroupVersionResource{
		Group:    resourceMap.K8sApiGroup,
		Version:  resourceMap.K8sApiVersion,
		Resource: resourceMap.K8sResource,
	}, nil
}

func CreateResourceMap(resourceMap models.ResourceMap) error {
	db := GetDB()
	return db.Create(&resourceMap).Error
}

func UpdateResourceMap(resourceMap models.ResourceMap) error {
	db := GetDB()
	return db.Model(&resourceMap).Updates(&resourceMap).Error
}
