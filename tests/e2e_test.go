// Package tests containing e2e tests
package tests

import (
	"bytes"
	"encoding/json"
	"github.com/rapido-liebre/pack_solver/internal/config"
	"net/http"
	"os"
	"testing"
	"time"

	httpapi "github.com/rapido-liebre/pack_solver/internal/http"
	"github.com/stretchr/testify/assert"
)

type orderRequest struct {
	Quantity int `json:"quantity"`
}

type orderResponse struct {
	Packs []struct {
		Size  int `json:"size"`
		Count int `json:"count"`
	} `json:"packs"`
	TotalItems int `json:"total_items"`
}

func TestOrderEndpointE2E(t *testing.T) {
	// Setup environment variables
	_ = os.Setenv("REDIS_ADDR", "localhost:6379")
	_ = os.Setenv("PACK_SOLVER_API", "http://localhost:8080")

	// Initialize Redis
	err := config.InitRedis()
	assert.NoError(t, err)

	// Set sample pack sizes in Redis for test to succeed
	err = config.SetPackSizes([]int{100, 250, 500, 1000})
	assert.NoError(t, err)

	// Start backend server in goroutine
	go func() {
		r := httpapi.SetupRouter()
		_ = r.Run(":8080")
	}()
	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Prepare request
	requestData := orderRequest{Quantity: 12001}
	body, _ := json.Marshal(requestData)
	baseURL := os.Getenv("PACK_SOLVER_API")

	// Send request to /order
	resp, err := http.Post(baseURL+"/order", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
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
