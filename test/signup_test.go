package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ngirimana/AnnounceIT/controllers"
	"github.com/ngirimana/AnnounceIT/db"
	"github.com/ngirimana/AnnounceIT/models"
	"github.com/stretchr/testify/assert"
)

// TestSignUp tests the signup function
func TestSignUpEmptyReq(t *testing.T) {
	router := gin.Default()
	router.POST("/users/signup", controllers.SignUp)

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
	router.POST("/users/signup", controllers.SignUp)

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

func TestSignUpUnformattedReq(t *testing.T) {
	// Initialize the database connection for the test
	db.InitDB()

	// Create a new router using Gin
	router := gin.Default()
	router.POST("/users/signup", controllers.SignUp)

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
	fmt.Print(resp.Body.String())
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
	router.POST("/users/login", controllers.Login)

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
	router.POST("/users/login", controllers.Login)

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
	router.POST("/users/login", controllers.Login)

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
	router.POST("/users/login", controllers.Login)

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
