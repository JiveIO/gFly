package redis

import (
	"errors"
	"fmt"
)

var (
	ErrConfigAddrEmpty       = errors.New("config Addr must not be empty")
	ErrConfigMasterNameEmpty = errors.New("config MasterName must not be empty")
)

func newErrRedisConnection(err error) error {
	return fmt.Errorf("redis connection error: %w", err)
}
