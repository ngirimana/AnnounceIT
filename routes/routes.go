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
	authenticated.POST("/announcements", controllers.CreateAnnouncement)

}
