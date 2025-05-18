package config_test

import (
	"github.com/rapido-liebre/pack_solver/internal/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSetAndGetPackSizes(t *testing.T) {
	err := os.Setenv("REDIS_ADDR", "localhost:6379")
	assert.NoError(t, err)

	err = config.InitRedis(os.Getenv("REDIS_ADDR"))
	assert.NoError(t, err)

	sizes := []int{100, 200, 300}
	err = config.SetPackSizes(sizes)
	assert.NoError(t, err)

	result, err := config.GetPackSizes()
	assert.NoError(t, err)
	assert.ElementsMatch(t, sizes, result)
}

func TestSetPackSizesInvalidRedis(t *testing.T) {
	err := os.Setenv("REDIS_ADDR", "localhost:6380")
	assert.NoError(t, err)

	err = config.InitRedis(os.Getenv("REDIS_ADDR"))
	assert.Error(t, err)

	// simulate unavailable redis
	sizes := []int{100, 200}
	err = config.SetPackSizes(sizes)
	assert.Error(t, err)
}

func TestGetPackSizesEmpty(t *testing.T) {
	err := os.Setenv("REDIS_ADDR", "localhost:6379")
	assert.NoError(t, err)

	err = config.InitRedis(os.Getenv("REDIS_ADDR"))
	assert.NoError(t, err)

	_, err = config.GetPackSizes()
	assert.NoError(t, err) // returns default config
}
