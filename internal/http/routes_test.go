// internal/http/routes_test.go
package http_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	routes "github.com/rapido-liebre/pack_solver/internal/http"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	routes.RegisterRoutes(r)
	return r
}

func TestHealthEndpoint(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestOrderEndpointInvalidInput(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	body := []byte(`{"quantity": -5}`)
	req, _ := http.NewRequest("POST", "/order", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestConfigPacksInvalidPayload(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	body := []byte(`{"pack_sizes": [0, -1, 1000]}`)
	req, _ := http.NewRequest("POST", "/config/packs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestConfigPacksMissingField(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	body := []byte(`{}`)
	req, _ := http.NewRequest("POST", "/config/packs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}
