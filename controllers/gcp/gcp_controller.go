package gcp

import (
	"cassata/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListGCPProjects(c *gin.Context) {
	workspaceID := c.Param("workspace_id")

	if workspaceID != "" {
		// Workspace-specific logic
		//projects = listProjectsForWorkspace(workspaceID)
	} else {
		// General GCP logic
		//projects = listAllProjects()
	}

	c.JSON(http.StatusOK, gin.H{
		"workspace_id": workspaceID,
		"projects":     []string{}, // Replace with actual projects
	})
}

func ListGCPResources(c *gin.Context) {
	projectID := c.Param("project")
	resourceType := c.Param("resource")
	workspaceID := c.Param("workspace")

	var resources []string // This should be replaced with your actual resource type

	if workspaceID != "" {
		// Workspace-specific logic
		//resources = listResourcesForWorkspace(workspaceID, projectID, resourceType)
	} else {
		// General GCP logic
		//resources = listResourcesForProject(projectID, resourceType)
	}

	c.JSON(http.StatusOK, gin.H{
		"project_id":    projectID,
		"resource_type": resourceType,
		"workspace_id":  workspaceID,
		"resources":     resources,
	})
}

func GetGCPResource(c *gin.Context) {
	// projectID := c.Param("project_id")
	resourceType := c.Param("resource_type")
	resourceID := c.Param("resource_id")
	workspaceName := c.Param("workspace_name")

	resource, err := service.GetResource("gcp", resourceType, workspaceName, resourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"resource": resource,
	})
}

func CreateGCPResource(c *gin.Context) {
	//projectID := c.Param("project_id") will be used to get the provider
	resourceType := c.Param("resource_type")
	workspaceName := c.Param("workspace_name")

	var requestBody map[string]interface{}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resource, err := service.CreateResource("gcp", resourceType, workspaceName, requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Resource created successfully",
		"resource": resource,
	})
}

func UpdateGCPResource(c *gin.Context) {
	projectID := c.Param("project_id")
	resourceType := c.Param("resource_type")
	resourceID := c.Param("resource_id")
	workspaceID := c.Param("workspace_id")
	var resource any

	if workspaceID != "" {
		// Workspace-specific logic
		//resource = updateResourceForWorkspace(workspaceID, projectID, resourceType, resourceID)
	} else {
		// General GCP logic
		//resource = updateResourceForProject(projectID, resourceType, resourceID)
	}

	c.JSON(http.StatusOK, gin.H{
		"project_id":    projectID,
		"resource_type": resourceType,
		"resource_id":   resourceID,
		"workspace_id":  workspaceID,
		"resource":      resource,
	})
}

func DeleteGCPResource(c *gin.Context) {
	projectID := c.Param("project_id")
	resourceType := c.Param("resource_type")
	resourceID := c.Param("resource_id")
	workspaceID := c.Param("workspace_id")

	if workspaceID != "" {
		// Workspace-specific logic
		//deleteResourceForWorkspace(workspaceID, projectID, resourceType, resourceID)
	} else {
		// General GCP logic
		//deleteResourceForProject(projectID, resourceType, resourceID)
	}

	c.JSON(http.StatusOK, gin.H{
		"project_id":    projectID,
		"resource_type": resourceType,
		"resource_id":   resourceID,
		"workspace_id":  workspaceID,
		"message":       "Resource deleted successfully",
	})
}

func ListWorkspaceGCPResourcesAllAccounts(c *gin.Context) {
	workspaceID := c.Param("workspace_id")
	resourceType := c.Param("resource_type")

	// Logic to list resources across all GCP projects for the given workspace
	// var resources []string = listWorkspaceResourcesAllProjects(workspaceID, resourceType)

	c.JSON(http.StatusOK, gin.H{
		"workspace_id":  workspaceID,
		"resource_type": resourceType,
		"resources":     []string{}, // Replace with actual resources
	})
}
