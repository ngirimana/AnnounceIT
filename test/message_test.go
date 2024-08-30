package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ngirimana/AnnounceIT/routes"
	"github.com/stretchr/testify/assert"
)

// TestMessage tests the message function
func TestMessage(t *testing.T) {
	router := gin.Default()
	router.GET("/message", routes.Message)

	req, err := http.NewRequest("GET", "/message", nil)
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	expectedBody := `{"message":"Hello World"}`
	assert.JSONEq(t, expectedBody, resp.Body.String())
}
