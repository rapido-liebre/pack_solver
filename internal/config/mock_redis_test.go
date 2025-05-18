package config_test

import (
	"encoding/json"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/rapido-liebre/pack_solver/internal/config"
	"github.com/stretchr/testify/assert"
	"os"
)

func TestGetPackSizesWithMockRedis(t *testing.T) {
	// Run miniredis on random port or run it on specific port by .RunAddr("localhost:6380")
	s, err := miniredis.Run()
	assert.NoError(t, err)
	defer s.Close()

	err = os.Setenv("REDIS_ADDR", s.Addr())
	assert.NoError(t, err)

	err = config.InitRedis()
	assert.NoError(t, err)

	value, _ := json.Marshal([]int{100, 250, 500})
	err = s.Set(config.PackSizesKey, string(value))
	assert.NoError(t, err)

	sizes, err := config.GetPackSizes()
	assert.NoError(t, err)
	assert.Equal(t, []int{100, 250, 500}, sizes)
}

func TestSetPackSizesWithMockRedis(t *testing.T) {
	s, err := miniredis.Run()
	assert.NoError(t, err)
	defer s.Close()

	err = os.Setenv("REDIS_ADDR", s.Addr())
	assert.NoError(t, err)

	err = config.InitRedis()
	assert.NoError(t, err)

	sizes := []int{300, 600, 900}
	err = config.SetPackSizes(sizes)
	assert.NoError(t, err)

	stored, _ := s.Get(config.PackSizesKey)
	var result []int
	err = json.Unmarshal([]byte(stored), &result)
	assert.NoError(t, err)
	assert.ElementsMatch(t, sizes, result)
}
