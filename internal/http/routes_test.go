// internal/http/routes_test.go
package http_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	httpapi "github.com/rapido-liebre/pack_solver/internal/http"
	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
	r := httpapi.SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestOrderEndpointInvalidInput(t *testing.T) {
	r := httpapi.SetupRouter()
	w := httptest.NewRecorder()
	body := []byte(`{"quantity": -5}`)
	req, _ := http.NewRequest("POST", "/order", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestConfigPacksInvalidPayload(t *testing.T) {
	r := httpapi.SetupRouter()
	w := httptest.NewRecorder()
	body := []byte(`{"pack_sizes": [0, -1, 1000]}`)
	req, _ := http.NewRequest("POST", "/config/packs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestConfigPacksMissingField(t *testing.T) {
	r := httpapi.SetupRouter()
	w := httptest.NewRecorder()
	body := []byte(`{}`)
	req, _ := http.NewRequest("POST", "/config/packs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}
