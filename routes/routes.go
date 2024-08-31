package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ngirimana/AnnounceIT/controllers"
)

func RegisterRoutes(server *gin.Engine) {

	server.POST("/users/signup", controllers.SignUp)
	server.POST("/users/login", controllers.Login)
}
