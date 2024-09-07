package main

import (
	"log"

	"cassata/config"
	"cassata/repository"
	"cassata/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	db, _ := config.InitDB()
	// db.AutoMigrate(&models.User{}, &models.Permission{}, &models.Workspace{})
	repository.InitRepository(db)

	router := gin.Default()

	routes.SetupRoutes(router, db)

	log.Println("Starting cassata server on :8080")
	log.Fatal(router.Run(":8080"))
}
