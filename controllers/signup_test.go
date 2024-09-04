package controllers

import (
	"encoding/json"
	"fmt"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ngirimana/AnnounceIT/db"
	"github.com/ngirimana/AnnounceIT/middlewares"
	"github.com/ngirimana/AnnounceIT/models"
	"github.com/stretchr/testify/assert"
)

// TestSignUp tests the signup function
func TestSignUpEmptyReq(t *testing.T) {
	router := gin.Default()
	router.POST("/users/signup", SignUp)

	req, err := http.NewRequest("POST", "/users/signup", nil)
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	expectedBody := `{"error":"could not parse the request"}`
	assert.JSONEq(t, expectedBody, resp.Body.String())
}

func TestSignUpSuccess(t *testing.T) {
	// Initialize the database connection for the test
	db.InitDB()

	// Clean up the users table before running the test
	db.TruncateUsersTable()

	// Create a new router using Gin
	router := gin.Default()
	router.POST("/users/signup", SignUp)

	// Create a request body with the correct phone number
	body := `{
		"email": "test@gmail.com",
		"password": "1234",
		"first_name": "Test",
		"last_name": "User",
		"phone_number": "+250781475108",
		"address": "KG 23 ST",
		"is_admin": false
	}`

	// Create a request to pass to our handler
	req, _ := http.NewRequest("POST", "/users/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusCreated, resp.Code)

	// Parse the actual response body
	var actualResponse map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)

	// Check the response message
	assert.Equal(t, "User created successfully", actualResponse["message"])

	// Verify the user object in the response
	user := actualResponse["user"].(map[string]interface{})
	assert.Equal(t, "test@gmail.com", user["email"])
	assert.Equal(t, "", user["password"])
	assert.Equal(t, "Test", user["first_name"])
	assert.Equal(t, "User", user["last_name"])
	assert.Equal(t, "+250781475108", user["phone_number"])
	assert.Equal(t, "KG 23 ST", user["address"])
	assert.Equal(t, false, user["is_admin"])

	// Check that the ID is a positive integer
	id, ok := user["id"].(float64) // JSON numbers are parsed as float64
	assert.True(t, ok)
	assert.Greater(t, id, float64(0))
}

func TestSignUpExistingUser(t *testing.T) {
	// Initialize the database connection for the test
	db.InitDB()

	// Clean up the users table before running the test

	// Create a new router using Gin
	router := gin.Default()
	router.POST("/users/signup", SignUp)

	// Create a request body with the correct phone number
	body := `{
		"email": "test@gmail.com",
		"password": "1234",
		"first_name": "Test",
		"last_name": "User",
		"phone_number": "+250781475108",
		"address": "KG 23 ST",
		"is_admin": false
	}`

	// Create a request to pass to our handler
	req, _ := http.NewRequest("POST", "/users/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusConflict, resp.Code)

	// Parse the actual response body
	var actualResponse map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &actualResponse)
	fmt.Println(actualResponse)
	assert.NoError(t, err)

	// Check the response message
	assert.Equal(t, "Conflict - user already exists", actualResponse["error"])
}

func TestSignUpUnformattedReq(t *testing.T) {
	// Initialize the database connection for the test
	db.InitDB()

	// Create a new router using Gin
	router := gin.Default()
	router.POST("/users/signup", SignUp)

	// Create a request body with the correct phone number
	body := `{
		
		"email": "test1@gmail.com"
		"password": "1234",
		"first_name": "Test",
		"last_name": "User",
		"phone_number": "+250781475109",
		"address": "KG 23 ST",
		"is_admin": false
	}`

	// Create a request to pass to our handler
	req, _ := http.NewRequest("POST", "/users/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	// Correct the expected value
	expected := `{"error":"could not parse the request"}`

	// Check the response body is what we expect
	assert.JSONEq(t, expected, resp.Body.String())
}

func TestLoginInvalidCredentials(t *testing.T) {
	// Initialize the database connection for the test
	db.InitDB()

	// Clean up the users table before running the test
	db.TruncateUsersTable()

	// Create a new router using Gin
	router := gin.Default()
	router.POST("/users/login", Login)

	// Create a request body with incorrect credentials
	body := `{
		"email": "wrongemail@gmail.com",
		"password": "wrongpassword"
	}`

	// Create a request to pass to our handler
	req, _ := http.NewRequest("POST", "/users/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusUnauthorized, resp.Code)

	// Parse the actual response body
	var actualResponse map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)

	// Check the response error message
	assert.Equal(t, "Invalid credentials", actualResponse["error"])
}

func TestLoginBadRequest(t *testing.T) {
	// Initialize the database connection for the test
	db.InitDB()

	// Create a new router using Gin
	router := gin.Default()
	router.POST("/users/login", Login)

	// Create a request body with invalid JSON
	body := `{
		"email": "missingpassword",
	}`

	// Create a request to pass to our handler
	req, _ := http.NewRequest("POST", "/users/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	// Parse the actual response body
	var actualResponse map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)

	// Check the response error message
	assert.Equal(t, "could not parse the request", actualResponse["error"])
}

func TestLoginUnauthorized(t *testing.T) {
	// Initialize the database connection for the test
	db.InitDB()

	// Create a new router using Gin
	router := gin.Default()
	router.POST("/users/login", Login)

	// Create a request body with invalid JSON
	body := `{
		"email": "missingpassword",
		"password": "1234"
	}`

	// Create a request to pass to our handler
	req, _ := http.NewRequest("POST", "/users/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusUnauthorized, resp.Code)

	// Parse the actual response body
	var actualResponse map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)

	// Check the response error message
	assert.Equal(t, "Invalid credentials", actualResponse["error"])
}

func TestLoginSuccess(t *testing.T) {
	// Initialize the database connection for the test
	db.InitDB()

	// Clean up the users table before running the test
	db.TruncateUsersTable()

	// Insert a user for the test
	user := models.User{
		Email:       "test@gmail.com",
		Password:    "1234", // Plain password for testing
		FirstName:   "Test",
		LastName:    "User",
		PhoneNumber: "+250781475108",
		Address:     "KG 23 ST",
		IsAdmin:     false,
	}
	user.Save() // Assuming that Save method also hashes the password before saving

	// Create a new router using Gin
	router := gin.Default()
	router.POST("/users/login", Login)

	// Create a request body with correct credentials
	body := `{
		"email": "test@gmail.com",
		"password": "1234"
	}`

	// Create a request to pass to our handler
	req, _ := http.NewRequest("POST", "/users/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusOK, resp.Code)

	// Parse the actual response body
	var actualResponse map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)

	// Check the response message
	assert.Equal(t, "User logged in successfully with JWT token", actualResponse["message"])

	// Verify that JWT token is not empty
	jwt, ok := actualResponse["jwt"].(string)
	assert.True(t, ok)
	assert.NotEmpty(t, jwt)

}
func TestCreateAnnouncementUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/announcements", middlewares.Authenticate, CreateAnnouncement)

	body := `
	{
  		"end_date": "2025-01-01T15:30:00.000Z",
  		"start_date": "2025-01-01T13:30:00.000Z",
  		"text": "This is a test announcement6"
	}`

	req, _ := http.NewRequest(http.MethodPost, "/announcements", strings.NewReader(body))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Authorization token is required")
}
func TestCreateAnnouncementInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/announcements", middlewares.Authenticate, CreateAnnouncement)

	body := `
	{
  		"end_date": "2025-01-01T15:30:00.000Z",
  		"start_date": "2025-01-01T13:30:00.000Z",
  		"text": "This is a test announcement6"
	}`

	req, _ := http.NewRequest(http.MethodPost, "/announcements", strings.NewReader(body))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "go")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid token")
}

func TestCreateAnnouncementSuccess(t *testing.T) {

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/announcements", middlewares.Authenticate, CreateAnnouncement)

	body := `
	{
  		"end_date": "2025-01-01T15:30:00.000Z",
  		"start_date": "2025-01-01T13:30:00.000Z",
  		"text": "This is a test announcement"
	}`

	req, _ := http.NewRequest(http.MethodPost, "/announcements", strings.NewReader(body))
	req.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3R1c2VyQGdtYWlsLmNvbSIsImV4cCI6MTcyNTUxODk1OCwidXNlcklkIjoxfQ.yYNVL1Id1WKDybxmDkuaZZJdDm_6msaUtnD_1GRn_rY")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Contains(t, rr.Body.String(), "Announcement created successfully")
}
func TestGetUser(t *testing.T) {
	// Initialize the Gin router
	router := gin.Default()

	// Define the route for GetUser
	router.GET("/users/:email", middlewares.Authenticate, GetUser)

	// Create test cases
	tests := []struct {
		name           string
		email          string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "User found",
			email:          "test@gmail.com",
			authHeader:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3R1c2VyQGdtYWlsLmNvbSIsImV4cCI6MTcyNTUxODk1OCwidXNlcklkIjoxfQ.yYNVL1Id1WKDybxmDkuaZZJdDm_6msaUtnD_1GRn_rY",
			expectedStatus: http.StatusOK,
			expectedBody:   "User retrieved successfully",
		},
		{
			name:           "User not found",
			email:          "test1@gmail.com",
			authHeader:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3R1c2VyQGdtYWlsLmNvbSIsImV4cCI6MTcyNTUxODk1OCwidXNlcklkIjoxfQ.yYNVL1Id1WKDybxmDkuaZZJdDm_6msaUtnD_1GRn_rY",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "user not found",
		},
		{
			name:           "Unauthorized request",
			email:          "test@gmail.com",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Authorization token is required",
		},
		{
			name:           "Unauthorized request",
			email:          "test@gmail.com",
			authHeader:     "dfvjdsjvdsjvbdsjbvds",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request to pass to our handler
			req, _ := http.NewRequest("GET", "/users/"+tt.email, nil)
			req.Header.Set("Content-Type", "application/json")
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Create a ResponseRecorder to record the response
			resp := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(resp, req)

			// Check the status code is what we expect
			assert.Equal(t, tt.expectedStatus, resp.Code)

			// Parse the actual response body
			var actualResponse map[string]interface{}

			err := json.Unmarshal(resp.Body.Bytes(), &actualResponse)
			assert.NoError(t, err)

			// Check the response message
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedBody, actualResponse["message"])
			} else {
				assert.Equal(t, tt.expectedBody, actualResponse["error"])
			}
		})
	}
}
