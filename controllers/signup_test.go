package controllers

import (
	"encoding/json"
	"io"

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

// TestSignUp tests the sign up functionality
func TestSignUp(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Disable Gin's default writer to prevent log output
	gin.DefaultWriter = io.Discard
	// Initialize the Gin router
	router := gin.Default()
	router.POST("/users/signup", SignUp)

	// Initialize the database connection and clean up before tests
	db.InitDB()
	db.TruncateUsersTable()

	// Define the test cases
	tests := []struct {
		name           string
		body           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Empty request",
			body:           "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "could not parse the request",
		},
		{
			name: "Successful signup",
			body: `{
				"email": "test@gmail.com",
				"password": "1234",
				"first_name": "Test",
				"last_name": "User",
				"phone_number": "+250781475108",
				"address": "KG 23 ST",
				"is_admin": false
			}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   "User created successfully",
		},
		{
			name: "User already exists",
			body: `{
				"email": "test@gmail.com",
				"password": "1234",
				"first_name": "Test",
				"last_name": "User",
				"phone_number": "+250781475108",
				"address": "KG 23 ST",
				"is_admin": false
			}`,
			expectedStatus: http.StatusConflict,
			expectedBody:   "Conflict - user already exists",
		},
		{
			name: "Unformatted request",
			body: `{
				"email": "test1@gmail.com"
				"password": "1234",
				"first_name": "Test",
				"last_name": "User",
				"phone_number": "+250781475109",
				"address": "KG 23 ST",
				"is_admin": false
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "could not parse the request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request to pass to our handler
			req, _ := http.NewRequest("POST", "/users/signup", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

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
			if tt.expectedStatus == http.StatusCreated {
				assert.Equal(t, tt.expectedBody, actualResponse["message"])

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
				id, ok := user["id"].(float64)
				assert.True(t, ok)
				assert.Greater(t, id, float64(0))
			} else {
				assert.Equal(t, tt.expectedBody, actualResponse["error"])
			}
		})
	}
}

// TestLogin tests the login functionality
func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Disable Gin's default writer to prevent log output
	gin.DefaultWriter = io.Discard
	// Initialize the Gin router
	router := gin.Default()
	router.POST("/users/login", Login)

	// Initialize the database connection and clean up before tests
	db.InitDB()
	db.TruncateUsersTable()

	// Insert a user for successful login test
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

	// Define the test cases
	tests := []struct {
		name            string
		body            string
		expectedStatus  int
		expectedError   string
		expectedMessage string
		checkToken      bool
	}{
		{
			name:           "Invalid credentials",
			body:           `{"email": "wrongemail@gmail.com", "password": "wrongpassword"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid credentials",
		},
		{
			name:           "Bad request - Invalid JSON",
			body:           `{"email": "missingpassword",}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "could not parse the request",
		},
		{
			name:           "Unauthorized - Wrong credentials",
			body:           `{"email": "missingpassword", "password": "1234"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid credentials",
		},
		{
			name:            "Successful login",
			body:            `{"email": "test@gmail.com", "password": "1234"}`,
			expectedStatus:  http.StatusOK,
			expectedMessage: "User logged in successfully with JWT token",
			checkToken:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request to pass to our handler
			req, _ := http.NewRequest("POST", "/users/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

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

			// Check the response message or error
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedMessage, actualResponse["message"])

				// Verify that JWT token is not empty if required
				if tt.checkToken {
					jwt, ok := actualResponse["jwt"].(string)
					assert.True(t, ok)
					assert.NotEmpty(t, jwt)
				}
			} else {
				assert.Equal(t, tt.expectedError, actualResponse["error"])
			}
		})
	}
}

func TestCreateAnnouncement(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Disable Gin's default writer to prevent log output
	gin.DefaultWriter = io.Discard
	// Initialize the Gin router
	router := gin.Default()
	router.POST("/announcements", middlewares.Authenticate, CreateAnnouncement)

	// Define the test cases
	tests := []struct {
		name           string
		body           string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Unauthorized - No token",
			body: `
			{
				"end_date": "2025-01-01T15:30:00.000Z",
				"start_date": "2025-01-01T13:30:00.000Z",
				"text": "This is a test announcement6"
			}`,
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Authorization token is required",
		},
		{
			name: "Unauthorized - Invalid token",
			body: `
			{
				"end_date": "2025-01-01T15:30:00.000Z",
				"start_date": "2025-01-01T13:30:00.000Z",
				"text": "This is a test announcement6"
			}`,
			authHeader:     "go",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid token",
		},
		{
			name: "Invalid request - Invalid JSON",
			body: `
			{
				"end_date": "2025-01-01T15:30:00.000Z",
				"start_date": "2025-01-01T13:30:00.000Z"
				"text": "This is a test announcement"
			}`,
			authHeader:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3R1c2VyQGdtYWlsLmNvbSIsImV4cCI6MTcyNTUxODk1OCwidXNlcklkIjoxfQ.yYNVL1Id1WKDybxmDkuaZZJdDm_6msaUtnD_1GRn_rY",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "could not parse request body",
		},
		{
			name: "Successful creation",
			body: `
			{
				"end_date": "2025-01-01T15:30:00.000Z",
				"start_date": "2025-01-01T13:30:00.000Z",
				"text": "This is a test announcement"
			}`,
			authHeader:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3R1c2VyQGdtYWlsLmNvbSIsImV4cCI6MTcyNTUxODk1OCwidXNlcklkIjoxfQ.yYNVL1Id1WKDybxmDkuaZZJdDm_6msaUtnD_1GRn_rY",
			expectedStatus: http.StatusCreated,
			expectedBody:   "Announcement created successfully",
		},
		{
			name: "Invalid request - Empty request body",
			body: `
			{
				"end_date": "2025-01-01T15:30:00.000Z",
				"start_date": "2025-01-01T13:30:00.000Z"
				"text": "This is a test announcement"
			}`,
			authHeader:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3R1c2VyQGdtYWlsLmNvbSIsImV4cCI6MTcyNTUxODk1OCwidXNlcklkIjoxfQ.yYNVL1Id1WKDybxmDkuaZZJdDm_6msaUtnD_1GRn_rY",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "could not parse request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request to pass to our handler
			req, _ := http.NewRequest(http.MethodPost, "/announcements", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", tt.authHeader)

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(rr, req)

			// Check the status code is what we expect
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Check the response body contains what we expect
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}

func TestGetUser(t *testing.T) {
	// Initialize the Gin router
	gin.SetMode(gin.TestMode)
	// Disable Gin's default writer to prevent log output
	gin.DefaultWriter = io.Discard
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
	db.TruncateUsersTable()
	db.TruncateAnnouncementsTable()
}
