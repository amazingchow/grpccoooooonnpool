package grpccoooooonnpool

import (
	"errors"
)

var (
	ErrInvalidPoolSettings          = errors.New("invalid pool settings")
	ErrPoolAlreadyClosed            = errors.New("pool is already closed")
	ErrPoolResourceAlreadyExhausted = errors.New("pool has no available connection")
)
