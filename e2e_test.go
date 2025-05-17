// Package tests containing e2e tests
package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type orderRequest struct {
	Quantity int `json:"quantity"`
}

type pack struct {
	Size  int `json:"size"`
	Count int `json:"count"`
}

type orderResponse struct {
	Packs      []pack `json:"packs"`
	TotalItems int    `json:"total_items"`
}

func TestOrderEndpointE2E(t *testing.T) {
	// Give backend time to start if needed
	time.Sleep(1 * time.Second)

	// Use env var or default to localhost
	baseURL := os.Getenv("PACK_SOLVER_API")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	requestData := orderRequest{Quantity: 12001}
	body, _ := json.Marshal(requestData)

	// Send real request to the endpoint /order
	resp, err := http.Post(baseURL+"/order", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	// Check response status 200
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result orderResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	// Check JSON is valid
	assert.NoError(t, err)
	// Check total packs count > 0
	assert.Greater(t, result.TotalItems, 0)
	// Check total packs count >= quantity
	assert.GreaterOrEqual(t, result.TotalItems, requestData.Quantity)
	assert.NotEmpty(t, result.Packs)
}
