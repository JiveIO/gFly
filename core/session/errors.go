package session

import (
	"errors"
)

var (
	ErrNotSetProvider = errors.New("not setted a session provider")
	ErrEmptySessionID = errors.New("empty session id")
)
