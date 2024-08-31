package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngirimana/AnnounceIT/models"
)

// SignUp godoc
// @Summary Sign up a new user
// @Description Create a new user in the system
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 201 {object} utils.SignUpSuccessResponse  "User created successfully"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /users/signup [post]
func SignUp(context *gin.Context) {

	var user models.User
	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "could not parse the request"})
		return
	}

	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.Password = ""
	context.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})

}
