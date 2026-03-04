package routes

import (
	"cloud-kitchen/internal/auth/controller"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.Engine, controller *controller.AuthController) {

	auth := router.Group("/auth")

	auth.POST("/signup", controller.Signup)
	auth.POST("/login", controller.Login)
}
