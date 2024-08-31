package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngirimana/AnnounceIT/helpers"
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

// Login godoc
// @Summary Login a user
// @Description Authenticate a user and return a JWT token
// @Tags Users
// @Accept json
// @Produce json
// @Param user body utils.LoginData true "User login credentials"
// @Success 200 {object} utils.LoginSuccessResponse "User logged in successfully with JWT token"
// @Failure 400 {object} utils.ErrorResponse "Bad Request - could not parse the request or generate token"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - invalid credentials"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error - server error"
// @Router /users/login [post]
func Login(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "could not parse the request"})
		return
	}

	err = user.Authenticate()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	jwt, err := helpers.GenerateToken(user.Email, user.ID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "could not generate token"})
		return
	}
	user.Password = ""
	context.JSON(http.StatusOK, gin.H{"message": "User logged in successfully with JWT token", "jwt": jwt, "user": user})
}
