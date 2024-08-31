package main

import (
	"github.com/gin-gonic/gin"

	"github.com/ngirimana/AnnounceIT/db"
	_ "github.com/ngirimana/AnnounceIT/docs" // Replace with your module name to match the generated docs import
	"github.com/ngirimana/AnnounceIT/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title AnnounceIT API
// @version 1.0
// @description AnnounceIT is a solution for broadcasting agencies to receive and manage announcements effectively.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /

// @schemes http https
func main() {
	db.InitDB()
	server := gin.Default()

	// Swagger endpoint to serve the API documentation
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.RegisterRoutes(server)

	// Start the server on port 8000
	server.Run(":8000")
}
