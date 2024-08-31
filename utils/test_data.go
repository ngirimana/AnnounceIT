package utils

import (
	"bytes"
	"encoding/json"

	"github.com/ngirimana/AnnounceIT/models"
)

// SuccessData creates a User struct and returns an io.Reader for use in HTTP requests.
func SuccessData() (*bytes.Buffer, error) {
	// Create a new user
	user := models.User{
		ID:          1,
		Email:       "schadrack@gmail.com",
		Password:    "1234",
		FirstName:   "Safari",
		LastName:    "Schadrack",
		PhoneNumber: "+250781475108",
		Address:     "KG 22 ST",
		IsAdmin:     false,
	}

	// Convert user struct to JSON
	jsonData, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	// Return the JSON data as an io.Reader using bytes.Buffer
	return bytes.NewBuffer(jsonData), nil
}
