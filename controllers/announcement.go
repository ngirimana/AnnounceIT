package controllers

import (
	"net/http"
	"strconv"

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
// @Success 201 {object} utils.AnnouncementSuccessResponse "Announcement created successfully"
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

// @Summary Get all announcements
// @Description Retrieve all announcements
// @Tags Announcements
// @Produce json
// @Success 200 {object} []models.Announcement "Announcements retrieved successfully"
// @Failure 500 {object} utils.ErrorResponse "Could not fetch announcements"
// @Router /announcements [get]
func GetAnnouncements(context *gin.Context) {
	announcements, err := models.GetAnnouncements()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch announcements"})
		return
	}

	if len(announcements) == 0 {
		context.JSON(http.StatusNotFound, gin.H{"error": "No announcements found"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"announcements": announcements, "message": "Announcements retrieved successfully"})
}

// @Summary Get a single announcement
// @Description Retrieve an announcement by its ID
// @Tags Announcements
// @Produce json
// @Param id path int true "Announcement ID"
// @Success 200 {object} utils.AnnouncementSuccessResponse "Announcement retrieved successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid announcement ID"
// @Failure 404 {object} utils.ErrorResponse "Announcement not found"
// @Router /announcements/{id} [get]
func GetAnnouncement(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid announcement ID"})
		return
	}
	announcement, err := models.GetAnnouncementByID(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Announcement not found"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"announcement": announcement, "message": "Announcement retrieved successfully"})
}
