// Package config initializes Redis, get and save packs configuration
package config

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

const packSizesKey = "pack:sizes"

// InitRedis initializes the Redis client with the given address.
func InitRedis(addr string) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr: addr,
	})
	_, err := redisClient.Ping(ctx).Result()
	return err
}

// GetPackSizes retrieves the pack sizes from Redis.
func GetPackSizes() ([]int, error) {
	val, err := redisClient.Get(ctx, packSizesKey).Result()
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
	return redisClient.Set(ctx, packSizesKey, data, 0).Err()
}
