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
	authenticated.POST("/announcements", controllers.CreateAnnouncement)

	server.GET("/announcements", controllers.GetAnnouncements)
	server.GET("/announcements/:id", controllers.GetAnnouncement)

}
