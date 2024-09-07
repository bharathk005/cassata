package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"cassata/models"
	"cassata/utils"
)

func CreatePermission(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Resource string     `json:"resource" binding:"required"`
			Action   utils.Verb `json:"action" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if permission already exists
		var existingPermission models.Permission
		if err := db.Where("resource = ? AND action = ?", input.Resource, input.Action).First(&existingPermission).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Permission already exists"})
			return
		} else if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing permission"})
			return
		}

		// Create new permission
		newPermission := models.Permission{
			Resource: input.Resource,
			Action:   input.Action,
		}

		if err := db.Create(&newPermission).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create permission"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Permission created successfully", "permission": newPermission})
	}
}

func AssignPermissionToWorkspace(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			WorkspaceID  uint `json:"workspace_id" binding:"required"`
			PermissionID uint `json:"permission_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if workspace exists
		var workspace models.Workspace
		if err := db.First(&workspace, input.WorkspaceID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found"})
			return
		}

		// Check if permission exists
		var permission models.Permission
		if err := db.First(&permission, input.PermissionID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
			return
		}

		// Assign permission to workspace
		if err := db.Model(&workspace).Association("Permissions").Append(&permission); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign permission to workspace"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Permission assigned to workspace successfully"})
	}
}
