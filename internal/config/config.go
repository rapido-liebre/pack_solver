// Package config initializes Redis, get and save packs configuration
package config

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"strings"
)

var redisClient *redis.Client
var ctx = context.Background()

const PackSizesKey = "pack:sizes"

// InitRedis initializes Redis connection using REDIS_ADDR from environment.
// It expects REDIS_ADDR in URI format (e.g., redis://user:pass@host:6379).
func InitRedis() error {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		return fmt.Errorf("REDIS_ADDR is not set")
	}

	var opt *redis.Options
	var err error

	if strings.HasPrefix(addr, "redis://") {
		opt, err = redis.ParseURL(addr)
		if err != nil {
			return fmt.Errorf("invalid REDIS_ADDR format: %w", err)
		}
	} else {
		opt = &redis.Options{Addr: addr}
	}

	redisClient = redis.NewClient(opt)

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
}

// GetPackSizes retrieves the pack sizes from Redis.
func GetPackSizes() ([]int, error) {
	val, err := redisClient.Get(ctx, PackSizesKey).Result()
	if err != nil {
		return nil, err
	}
	var sizes []int
	if err := json.Unmarshal([]byte(val), &sizes); err != nil {
		return nil, err
	}
	return sizes, nil
}

// SetPackSizes stores the pack sizes in Redis as a JSON array.
func SetPackSizes(sizes []int) error {
	data, err := json.Marshal(sizes)
	if err != nil {
		return err
	}
	return redisClient.Set(ctx, PackSizesKey, data, 0).Err()
}
