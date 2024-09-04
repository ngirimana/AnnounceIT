package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngirimana/AnnounceIT/models"
)

// CreateAnnouncement godoc
// @Summary Create an announcement
// @Description Create an announcement and save it to the database
// @Tags Announcements
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param announcement body models.Announcement true "Announcement object"
// @Success 201 {object} utils.CreateAnnouncementSuccessResponse "Announcement created successfully"
// @Failure 400 {object} utils.ErrorResponse "Could not parse request body"
// @Failure 401 {object} utils.ErrorResponse "Authorization token is required or invalid"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /announcements [post]
func CreateAnnouncement(context *gin.Context) {
	var announcement models.Announcement
	err := context.ShouldBindJSON(&announcement)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "could not parse request body"})
		return
	}

	announcement.OwnerID = context.GetInt64("userId")
	err = announcement.Create()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Announcement created successfully", "announcement": announcement})
}
