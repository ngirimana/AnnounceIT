package utils

import "github.com/ngirimana/AnnounceIT/models"

type UserSuccessResponse struct {
	Message string      `json:"message"` // The success message
	Data    models.User `json:"user"`    // The user data
}

type ErrorResponse struct {
	Error string `json:"error"` // The error message
}

type LoginSuccessResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AnnouncementSuccessResponse struct {
	Message string              `json:"message"`
	Data    models.Announcement `json:"announcement"`
}

type DeleteSuccessResponse struct {
	Message string `json:"message"`
}
