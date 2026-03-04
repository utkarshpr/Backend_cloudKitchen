package main

import (
	"cloud-kitchen/internal/auth/controller"
	"cloud-kitchen/internal/auth/repository"
	"cloud-kitchen/internal/auth/routes"
	"cloud-kitchen/internal/auth/service"
	"cloud-kitchen/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	db, err := database.NewPostgresConnection()
	if err != nil {
		panic(err)
	}

	repo := repository.NewAuthRepository(db)

	authService := service.NewAuthService(repo)

	authController := controller.NewAuthController(authService)

	routes.RegisterAuthRoutes(router, authController)

	router.Run(":8080")
}
