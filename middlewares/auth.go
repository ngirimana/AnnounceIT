package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngirimana/AnnounceIT/helpers"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")
	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	userId, isAdmin, err := helpers.VerifyToken(token)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	context.Set("userId", userId)
	context.Set("isAdmin", isAdmin)
	context.Next()
}
