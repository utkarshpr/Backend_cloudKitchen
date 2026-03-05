package main

import (
	authController "cloud-kitchen/internal/auth/controller"
	authRepository "cloud-kitchen/internal/auth/repository"
	authRoutes "cloud-kitchen/internal/auth/routes"
	authService "cloud-kitchen/internal/auth/service"

	profileController "cloud-kitchen/internal/profile/controller"
	profileRepository "cloud-kitchen/internal/profile/repository"
	profileRoutes "cloud-kitchen/internal/profile/routes"
	profileService "cloud-kitchen/internal/profile/service"

	"cloud-kitchen/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	db, err := database.NewPostgresConnection()
	if err != nil {
		panic(err)
	}

	// AUTH MODULE
	authRepo := authRepository.NewAuthRepository(db)
	authSvc := authService.NewAuthService(authRepo)
	authCtrl := authController.NewAuthController(authSvc)
	authRoutes.RegisterAuthRoutes(router, authCtrl)

	// PROFILE MODULE
	profileRepo := profileRepository.NewProfileRepository(db)
	profileSvc := profileService.NewProfileService(profileRepo)
	profileCtrl := profileController.NewProfileController(profileSvc)
	profileRoutes.RegisterProfileRoutes(router, profileCtrl)

	router.Run(":8080")
}