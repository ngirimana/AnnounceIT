package controllers

import (
	"fmt"
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

// UpdateAnnouncement godoc
// @Summary Update an announcement
// @Description Updates an existing announcement by its ID. The user must be the owner of the announcement to update it.
// @Tags Announcements
// @Param id path int true "Announcement ID"
// @Param Authorization header string true "Bearer token"
// @Param announcement body models.Announcement true "Update Announcement"
// @Success 200 {object} utils.AnnouncementSuccessResponse "Announcement updated successfully"
// @Failure 400 {object} utils.ErrorResponse 	"Could not parse request body"
// @Failure 403 {object} utils.ErrorResponse  "You are not allowed to update this announcement"
// @Failure 404 {object} utils.ErrorResponse  "Announcement not found"
// @Failure 500 {object} utils.ErrorResponse "Could not update announcement"
// @Router /announcements/{id} [patch]
func UpdateAnnouncement(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid announcement ID"})
		return
	}

	userId := context.GetInt64("userId")
	announcement, err := models.GetAnnouncementByID(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Announcement not found"})
		return
	}

	if announcement.OwnerID != userId {
		context.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this announcement"})
		return
	}

	var updateAnnouncement models.Announcement
	err = context.ShouldBindJSON(&updateAnnouncement)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "could not parse request body"})
		return
	}

	fmt.Println(updateAnnouncement)
	updateAnnouncement.ID = id
	err = updateAnnouncement.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update announcement"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Announcement updated successfully", "announcement": updateAnnouncement})
}

// DeleteAnnouncement godoc
// @Summary Delete an announcement
// @Description Deletes an existing announcement by its ID. The user must be the owner of the announcement to delete it.
// @Tags Announcements
// @Param id path int true "Announcement ID"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} utils.DeleteSuccessResponse "Announcement deleted successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid announcement ID"
// @Failure 403 {object} utils.ErrorResponse "You are not allowed to delete this announcement"
// @Failure 404 {object} utils.ErrorResponse  "Announcement not found"
// @Failure 500 {object} utils.ErrorResponse "Could not delete announcement"
// @Router /announcements/{id} [delete]
func DeleteAnnouncement(context *gin.Context) {
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

	isAdmin := context.GetBool("isAdmin")
	if isAdmin {
		context.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this announcement"})
		return
	}
	err = announcement.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete announcement"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Announcement deleted successfully"})

}

func parseStatus(status string) (models.Status, error) {
	val, err := strconv.Atoi(status)
	if err != nil || val < 0 || val > 4 {
		return -1, fmt.Errorf("invalid status value")
	}
	return models.Status(val), nil
}

// ChangeAnnouncementStatus godoc
// @Summary Change announcement status
// @Description Change the status of an announcement by its ID
// @Tags Announcements
// @Param id path int true "Announcement ID"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} utils.AnnouncementSuccessResponse "Announcement status updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid announcement ID or status value"
// @Failure 404 {object} utils.ErrorResponse "Announcement not found"
// @Failure 401 {object} utils.ErrorResponse "You are not allowed to update this announcement"
// @Failure 500 {object} utils.ErrorResponse "Could not update announcement status"
// @Router /announcements/{id}/status [patch]
func ChangeAnnouncementStatus(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid announcement ID"})
		return
	}

	isAdmin := context.GetBool("isAdmin")
	if isAdmin {
		context.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this announcement"})
		return
	}

	status := context.PostForm("status")
	if status == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "status is required"})
		return
	}
	statusVal, err := parseStatus(status)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}
	err = models.ChangeAnnouncementStatus(id, statusVal)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update announcement status"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Announcement status updated successfully"})
}
