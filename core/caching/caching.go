package caching

import (
	"app/core/utils"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// New func for connecting to Redis server.
func New() (*redis.Client, error) {
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

// Key wrapper key
func Key(key string) string {
	return fmt.Sprintf("%s:%s", utils.Getenv("APP_CODE", "gfly"), key)
}
