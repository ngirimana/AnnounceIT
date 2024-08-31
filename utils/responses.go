package utils

import "github.com/ngirimana/AnnounceIT/models"

type SignUpSuccessResponse struct {
	Message string      `json:"message"` // The success message
	Data    models.User `json:"data"`    // The user data
}

type ErrorResponse struct {
	Error string `json:"error"` // The error message
}
