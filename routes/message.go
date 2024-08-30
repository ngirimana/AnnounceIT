package routes

import "github.com/gin-gonic/gin"

// @Summary      Show a message
// @Description  Returns a "Hello World" message
// @Tags         message
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /message [get]
func Message(context *gin.Context) {
	context.JSON(200, gin.H{
		"message": "Hello World",
	})
}
