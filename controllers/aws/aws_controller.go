package aws

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListAWSAccounts(c *gin.Context) {
	workspaceID := c.Param("workspace_id")

	if workspaceID != "" {
		// Workspace-specific logic
		//resources = listResourcesForWorkspace(workspaceID, accountID, resourceType)
	} else {
		// General AWS logic
		//resources = listResourcesForAccount(accountID, resourceType)
	}
}

func ListAWSResources(c *gin.Context) {
	accountID := c.Param("account_id")
	resourceType := c.Param("resource_type")
	workspaceID := c.Param("workspace_id")

	var resources []string // This should be replaced with your actual resource type

	if workspaceID != "" {
		// Workspace-specific logic
		//resources = listResourcesForWorkspace(workspaceID, accountID, resourceType)
	} else {
		// General AWS logic
		//resources = listResourcesForAccount(accountID, resourceType)
	}

	c.JSON(http.StatusOK, gin.H{
		"account_id":    accountID,
		"resource_type": resourceType,
		"workspace_id":  workspaceID,
		"resources":     resources,
	})
}

func GetAWSResource(c *gin.Context) {
	accountID := c.Param("account_id")
	resourceType := c.Param("resource_type")
	resourceID := c.Param("resource_id")
	workspaceID := c.Param("workspace_id")
	var resource any
	if workspaceID != "" {
		// Workspace-specific logic
		//resources = listResourcesForWorkspace(workspaceID, accountID, resourceType)
	} else {
		// General AWS logic
		//resources = listResourcesForAccount(accountID, resourceType)
	}

	c.JSON(http.StatusOK, gin.H{
		"account_id":    accountID,
		"resource_type": resourceType,
		"resource_id":   resourceID,
		"workspace_id":  workspaceID,
		"resource":      resource,
	})
}

func CreateAWSResource(c *gin.Context) {
	accountID := c.Param("account_id")
	resourceType := c.Param("resource_type")
	workspaceID := c.Param("workspace_id")
	var resource any
	if workspaceID != "" {
		// Workspace-specific logic
		//resources = listResourcesForWorkspace(workspaceID, accountID, resourceType)
	} else {
		// General AWS logic
		//resources = listResourcesForAccount(accountID, resourceType)
	}

	c.JSON(http.StatusOK, gin.H{
		"account_id":    accountID,
		"resource_type": resourceType,
		"workspace_id":  workspaceID,
		"resource":      resource,
	})
}

func UpdateAWSResource(c *gin.Context) {
	accountID := c.Param("account_id")
	resourceType := c.Param("resource_type")
	resourceID := c.Param("resource_id")
	workspaceID := c.Param("workspace_id")
	var resource any
	if workspaceID != "" {
		// Workspace-specific logic
		//resources = listResourcesForWorkspace(workspaceID, accountID, resourceType)
	} else {
		// General AWS logic
		//resources = listResourcesForAccount(accountID, resourceType)
	}

	c.JSON(http.StatusOK, gin.H{
		"account_id":    accountID,
		"resource_type": resourceType,
		"resource_id":   resourceID,
		"workspace_id":  workspaceID,
		"resource":      resource,
	})
}

func DeleteAWSResource(c *gin.Context) {
	accountID := c.Param("account_id")
	resourceType := c.Param("resource_type")
	workspaceID := c.Param("workspace_id")
	var resource any
	if workspaceID != "" {
		// Workspace-specific logic
		//resources = listResourcesForWorkspace(workspaceID, accountID, resourceType)
	} else {
		// General AWS logic
		//resources = listResourcesForAccount(accountID, resourceType)
	}

	c.JSON(http.StatusOK, gin.H{
		"account_id":    accountID,
		"resource_type": resourceType,
		"workspace_id":  workspaceID,
		"resource":      resource,
	})
}

func ListWorkspaceAWSResourcesAllAccounts(c *gin.Context) {

}
