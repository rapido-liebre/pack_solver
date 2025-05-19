package main

import (
	"github.com/rapido-liebre/pack_solver/internal/config"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMainHealthEndpoint(t *testing.T) {
	// Setup environment variables
	_ = os.Setenv("REDIS_ADDR", "localhost:6379")
	_ = os.Setenv("PACK_SOLVER_API", "http://localhost:8080")

	// Initialize Redis
	err := config.InitRedis()
	assert.NoError(t, err)

	// Run app as goroutine
	go func() {
		main()
	}()
	// Wait for server to start
	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://localhost:8080/health")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
