package service

// import (
// 	"log"
// 	"net/http"

// 	"cassata/utils"

// 	"github.com/gin-gonic/gin"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
// 	"k8s.io/client-go/dynamic"
// 	"k8s.io/client-go/kubernetes"
// )

// type K8sService struct {
// 	clientset     *kubernetes.Clientset
// 	dynamicClient *dynamic.Interface
// }

// func NewK8sService() *K8sService {

// 	clientset, err := utils.GetKubernetesClientset()

// 	if err != nil {
// 		log.Fatalf("Error creating Kubernetes clientset: %v", err)
// 	}
// 	dynamicClient, err := utils.GetDynamicClient()
// 	if err != nil {
// 		log.Fatalf("Error creating Kubernetes clientset: %v", err)
// 	}

// 	return &K8sService{
// 		clientset:     clientset,
// 		dynamicClient: &dynamicClient,
// 	}
// }

// func (c *K8sService) ListResources(ctx *gin.Context, resourceType string) {

// 	gvr, err := utils.GetGroupVersionResource(resourceType)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource type"})
// 		return
// 	}

// 	list, err := (*c.dynamicClient).Resource(gvr).Namespace("").List(ctx, metav1.ListOptions{})
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list resources"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, list)
// }

// func (c *K8sService) CreateResource(ctx *gin.Context, resourceType string) {
// 	dynamicClient, err := utils.GetDynamicClient()
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create dynamic client"})
// 		return
// 	}

// 	gvr, err := utils.GetGroupVersionResource(resourceType)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource type"})
// 		return
// 	}

// 	var resource unstructured.Unstructured
// 	if err := ctx.ShouldBindJSON(&resource); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource data"})
// 		return
// 	}

// 	namespace := resource.GetNamespace()
// 	if namespace == "" {
// 		namespace = "default"
// 	}

// 	created, err := dynamicClient.Resource(gvr).Namespace(namespace).Create(ctx, &resource, metav1.CreateOptions{})
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create resource"})
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, created)
// }

// func (c *K8sService) GetResource(ctx *gin.Context, resourceType string) {
// 	dynamicClient, err := utils.GetDynamicClient()
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create dynamic client"})
// 		return
// 	}

// 	gvr, err := utils.GetGroupVersionResource(resourceType)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource type"})
// 		return
// 	}

// 	name := ctx.Param("name")
// 	namespace := ctx.DefaultQuery("namespace", "default")

// 	resource, err := dynamicClient.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
// 	if err != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, resource)
// }

// func (c *K8sService) UpdateResource(ctx *gin.Context, resourceType string) {
// 	dynamicClient, err := utils.GetDynamicClient()
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create dynamic client"})
// 		return
// 	}

// 	gvr, err := utils.GetGroupVersionResource(resourceType)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource type"})
// 		return
// 	}

// 	namespace := ctx.DefaultQuery("namespace", "default")

// 	var resource unstructured.Unstructured
// 	if err := ctx.ShouldBindJSON(&resource); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource data"})
// 		return
// 	}

// 	updated, err := dynamicClient.Resource(gvr).Namespace(namespace).Update(ctx, &resource, metav1.UpdateOptions{})
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update resource"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, updated)
// }

// func (c *K8sService) DeleteResource(ctx *gin.Context, resourceType string) {
// 	dynamicClient, err := utils.GetDynamicClient()
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create dynamic client"})
// 		return
// 	}

// 	gvr, err := utils.GetGroupVersionResource(resourceType)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource type"})
// 		return
// 	}

// 	name := ctx.Param("name")
// 	namespace := ctx.DefaultQuery("namespace", "default")

// 	err = dynamicClient.Resource(gvr).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{})
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete resource"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"message": "Resource deleted successfully"})
// }
