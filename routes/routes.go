package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ngirimana/AnnounceIT/controllers"
	"github.com/ngirimana/AnnounceIT/middlewares"
)

func RegisterRoutes(server *gin.Engine) {

	server.POST("/users/signup", controllers.SignUp)
	server.POST("/users/login", controllers.Login)
	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.GET("/users/:email", controllers.GetUser)
	authenticated.PATCH("/users/:id/flag", controllers.FlagUser)
	authenticated.POST("/announcements", controllers.CreateAnnouncement)
	authenticated.PATCH("/announcements/:id", controllers.UpdateAnnouncement)
	authenticated.DELETE("/announcements/:id", controllers.DeleteAnnouncement)
	authenticated.PATCH("/announcements/:id/status", controllers.ChangeAnnouncementStatus)

	server.GET("/announcements", controllers.GetAnnouncements)
	server.GET("/announcements/:id", controllers.GetAnnouncement)
}
