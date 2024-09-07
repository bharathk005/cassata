package service

import (
	"cassata/config"
	"cassata/models"
	"fmt"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		type RegisterUserInput struct {
			Name       string   `json:"name" binding:"required"`
			Password   string   `json:"password" binding:"required"`
			Workspaces []string `json:"workspaces" binding:"required"`
		}

		var input RegisterUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if all the workspaces exist
		var workspaces []models.Workspace
		for _, workspaceName := range input.Workspaces {
			var workspace models.Workspace
			if err := db.Where("name = ?", workspaceName).First(&workspace).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "One or more specified workspaces do not exist"})
				return
			}
			workspaces = append(workspaces, workspace)
		}

		if len(workspaces) != len(input.Workspaces) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "One or more specified workspaces do not exist"})
			return
		}

		// Create user
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user := models.User{
			Name:         input.Name,
			PasswordHash: string(hashedPassword),
			Workspaces:   workspaces,
		}
		if err := db.Create(&user).Error; err != nil {
			if err == gorm.ErrDuplicatedKey {
				c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			}
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
	}
}

func LoginUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginData struct {
			Name     string `json:"name" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		if err := db.Where("name = ?", loginData.Name).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginData.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": fmt.Sprintf("%d", user.ID),
			"exp":    time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString([]byte(config.Env[config.JWT_SECRET_KEY]))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}
