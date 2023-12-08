package gfly

import (
	"app/core/utils"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// Caching func for connecting to Redis server.
func Caching() (*redis.Client, error) {
	// Define Redis database number.
	dbNumber := utils.Getenv("REDIS_DB_NUMBER", 0)

	// Build Redis connection URL.
	redisConnURL := fmt.Sprintf(
		"%s:%d",
		utils.Getenv("REDIS_HOST", "localhost"),
		utils.Getenv("REDIS_PORT", 6379),
	)

	// Set Redis options.
	options := &redis.Options{
		Addr:     redisConnURL,
		Password: utils.Getenv("REDIS_PASSWORD", ""),
		DB:       dbNumber,
	}

	return redis.NewClient(options), nil
}
