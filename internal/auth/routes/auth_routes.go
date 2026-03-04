package routes

import (
	"cloud-kitchen/internal/auth/controller"
	"cloud-kitchen/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.Engine, controller *controller.AuthController) {

	router.Use(middleware.RequestIDMiddleware())
	auth := router.Group("/auth")

	auth.POST("/signup", func(c *gin.Context) {
		controller.Signup(c.Writer, c.Request)
	})
	auth.POST("/login", func(c *gin.Context) {
		controller.Login(c.Writer, c.Request)
	})
}
