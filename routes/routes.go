package routes

import (
	"cassata/controllers/aws"
	"cassata/controllers/gcp"
	"cassata/middleware"
	"cassata/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TODO need resource_group in route.
// Example iam.policy vs dns.policy.
// Need to update database, model and repository to support this.
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	setupLoginRoutes(router, db)

	setupAdminRoutes(router, db)

	workspaceGroup := router.Group("/workspace/:workspace_name")
	workspaceGroup.Use(middleware.AuthenticateUser())
	workspaceGroup.Use(middleware.AuthorizeUser(db))

	setupAWSRoutes(router, db, workspaceGroup)
	setupGCPRoutes(router, db, workspaceGroup)
}

func setupAdminRoutes(router *gin.Engine, db *gorm.DB) {
	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.AuthenticateUser())
	adminGroup.Use(middleware.AuthorizeUser(db))

	adminGroup.POST("/users", service.RegisterUser(db))
	adminGroup.POST("/user-assignments", service.AssignUserToWorkspace(db))
	adminGroup.POST("/workspaces", service.CreateWorkspace(db))
	adminGroup.POST("/permissions", service.CreatePermission(db))
	adminGroup.POST("/permission-assignments", service.AssignPermissionToWorkspace(db))
}

func setupLoginRoutes(router *gin.Engine, db *gorm.DB) {
	router.Group("/login")
	router.POST("/login", service.LoginUser(db))
}

func setupGCPRoutes(router *gin.Engine, db *gorm.DB, workspaceGroup *gin.RouterGroup) {
	gcpGroup := router.Group("/gcp")
	gcpGroup.Use(middleware.AuthenticateUser())
	gcpGroup.Use(middleware.AuthorizeUser(db))

	gcpEndpoints := func(r *gin.RouterGroup) {
		r.GET("/projects", gcp.ListGCPProjects)
		r.GET("/projects/:project_id/:resource_group/:resource_type", gcp.ListGCPResources)
		r.GET("/projects/:project_id/:resource_group/:resource_type/:resource_id", gcp.GetGCPResource)
		r.POST("/projects/:project_id/:resource_group/:resource_type", gcp.CreateGCPResource)
		r.PUT("/projects/:project_id/:resource_group/:resource_type/:resource_id", gcp.UpdateGCPResource)
		r.DELETE("/projects/:project_id/:resource_group/:resource_type/:resource_id", gcp.DeleteGCPResource)
	}

	gcpEndpoints(gcpGroup)
	workspaceGCP := workspaceGroup.Group("/gcp")
	gcpEndpoints(workspaceGCP)
	workspaceGCP.GET("/:resource_group/:resource_type", gcp.ListWorkspaceGCPResourcesAllAccounts)
}

func setupAWSRoutes(router *gin.Engine, db *gorm.DB, workspaceGroup *gin.RouterGroup) {
	awsGroup := router.Group("/aws")
	awsGroup.Use(middleware.AuthenticateUser())
	awsGroup.Use(middleware.AuthorizeUser(db))

	awsEndpoints := func(r *gin.RouterGroup) {
		r.GET("/accounts", aws.ListAWSAccounts)
		r.GET("/accounts/:account_id/:resource_group/:resource_type", aws.ListAWSResources)
		r.GET("/accounts/:account_id/:resource_group/:resource_type/:resource_id", aws.GetAWSResource)
		r.POST("/accounts/:account_id/:resource_group/:resource_type", aws.CreateAWSResource)
		r.PUT("/accounts/:account_id/:resource_group/:resource_type/:resource_id", aws.UpdateAWSResource)
		r.DELETE("/accounts/:account_id/:resource_group/:resource_type/:resource_id", aws.DeleteAWSResource)
	}

	awsEndpoints(awsGroup)
	workspaceAWS := workspaceGroup.Group("/aws")
	awsEndpoints(workspaceAWS)
	workspaceAWS.GET("/:resource_group/:resource_type", aws.ListWorkspaceAWSResourcesAllAccounts)
}
