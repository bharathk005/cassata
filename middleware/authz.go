package middleware

import (
	"cassata/models"
	"cassata/utils"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthorizeUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, exists := c.Get("userId")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		userIdStr := userId.(string)
		//get all workspaces for user from database
		var user models.User
		userIdUint, err := strconv.ParseUint(userIdStr, 10, 64)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
			return
		}
		if err := db.Preload("Workspaces.Permissions").Where("id = ?", userIdUint).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user workspaces"})
			return
		}

		if len(user.Workspaces) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User does not have access to any workspaces"})
			return
		}

		path := c.Request.URL.Path
		methodVerb := utils.HTTPMethodToVerb[c.Request.Method]
		// Split path into workspace and resource
		pathParts := strings.Split(path, "/")
		var workspaceName string
		var resource string

		if len(pathParts) > 2 && pathParts[1] == "workspace" {
			workspaceName = pathParts[2]
			resource = "/" + strings.Join(pathParts[1:], "/")
		} else if pathParts[1] == "admin" {
			workspaceName = "admin"
			resource = "/" + strings.Join(pathParts[1:], "/")
		} else {
			workspaceName = "global"
			resource = "/" + strings.Join(pathParts[1:], "/")
		}

		// Check if the workspace is in user's workspaces
		// var userWorkspace models.Workspace
		// found := false
		// fmt.Println("user workspaces", user.Workspaces, workspaceName)
		// for _, ws := range user.Workspaces {
		// 	if ws.Name == workspaceName {
		// 		userWorkspace = ws
		// 		found = true
		// 		break
		// 	}
		// }
		// if !found {
		// 	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("User does not have access to workspace %s", workspaceName)})
		// 	return
		// }

		// Check if workspace has permission for path and method
		var allowed bool
		for _, userWorkspace := range user.Workspaces {
			for _, perm := range userWorkspace.Permissions {
				// fmt.Println("perm", perm.Resource, resource, perm.Action, methodVerb)

				// Convert the permission resource to a regex pattern
				pattern := "^" + strings.Replace(strings.Replace(regexp.QuoteMeta(perm.Resource), "\\*", ".*", -1), "/", "\\/", -1) + "$"
				matched, err := regexp.MatchString(pattern, resource)

				if err != nil {
					fmt.Printf("Error matching regex: %v\n", err)
					continue
				}

				if matched && perm.Action == methodVerb {
					allowed = true
					break
				}
			}
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("User does not have access to %s %s", methodVerb, resource)})
			return
		}

		c.Set("workspace", workspaceName)
		c.Set("resource", resource)

		c.Next()
	}
}
