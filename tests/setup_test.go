package tests

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/judascrow/go-apix/api/routes"
	"github.com/stretchr/testify/assert"
)

func TestHealthcheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/healthcheck", routes.Healthcheck)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthcheck", nil)
	r.ServeHTTP(w, req)

	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(w.Body.String()), &jsonMap)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, w.Body)
	assert.Equal(t, true, jsonMap["success"])
	assert.Equal(t, "API is Online", jsonMap["message"])
}
