package service

import (
	"cassata/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateWorkspace(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name string `json:"name" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if workspace already exists
		var existingWorkspace models.Workspace
		if err := db.Where("name = ?", input.Name).First(&existingWorkspace).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Workspace already exists"})
			return
		} else if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing workspace"})
			return
		}

		// Create new workspace
		newWorkspace := models.Workspace{
			Name: input.Name,
		}

		if err := db.Create(&newWorkspace).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workspace"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Workspace created successfully", "workspace": newWorkspace})
	}
}

func AssignUserToWorkspace(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			WorkspaceID uint `json:"workspace_id" binding:"required"`
			UserID      uint `json:"user_id" binding:"required"`
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

		// Check if user exists
		var user models.User
		if err := db.First(&user, input.UserID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Assign user to workspace
		if err := db.Model(&workspace).Association("Users").Append(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign user to workspace"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User assigned to workspace successfully"})
	}
}
