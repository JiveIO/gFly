package gfly

import (
	"app/core/errors"
	"app/core/log"
	"app/core/session"
	"app/core/session/providers/memory"
	"app/core/session/providers/redis"
	"app/core/utils"
	"fmt"
	"time"
)

// ===========================================================================================================
// 										Session
// ===========================================================================================================

var serverSession *session.Session

type providerType string

const (
	redisProvider  = providerType("redis")
	memoryProvider = providerType("memory")
)

func sessionFactory(provider providerType) (session.Provider, error) {
	switch provider {
	case redisProvider:
		// Build Redis connection URL.
		redisConnURL := fmt.Sprintf(
			"%s:%d",
			utils.Getenv("REDIS_HOST", "localhost"),
			utils.Getenv("REDIS_PORT", 6379),
		)

		return redis.New(redis.Config{
			KeyPrefix:       utils.Getenv("SESSION_KEY", "gfly_session"),
			Addr:            redisConnURL,
			PoolSize:        8,
			ConnMaxIdleTime: 30 * time.Second,
		})
	case memoryProvider:
		return memory.New(memory.Config{})
	}

	return nil, errors.NotYetImplemented
}

func setupSession() {
	if utils.Getenv("SESSION_TYPE", "") == "" {
		log.Trace("Disable Session")

		return
	}

	providerType := providerType(utils.Getenv("SESSION_TYPE", "memory"))

	provider, err := sessionFactory(providerType)
	if err != nil {
		log.Fatal(err)
	}

	cfg := session.NewDefaultConfig()
	cfg.EncodeFunc = session.MSGPEncode
	cfg.DecodeFunc = session.MSGPDecode
	serverSession = session.New(cfg)

	if err = serverSession.SetProvider(provider); err != nil {
		log.Fatal(err)
	}

	log.Trace("Initialize Session")
}
