package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ngirimana/AnnounceIT/controllers"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/message", Message)

	server.POST("/users/signup", controllers.SignUp)
}
