package routes

import (
	"cloud-kitchen/internal/profile/controller"

	"github.com/gin-gonic/gin"
)

func RegisterProfileRoutes(router *gin.Engine, controller *controller.ProfileController) {

	profile := router.Group("/profile")
	profile.GET("/", func(c *gin.Context) {
		controller.GetProfile(c.Writer, c.Request)
	})
	profile.PUT("/", func(c *gin.Context) {
		controller.UpdateProfile(c.Writer, c.Request)
	})

	profile.DELETE("/", func(c *gin.Context) {
		controller.DeleteProfile(c.Writer, c.Request)
	})

}
