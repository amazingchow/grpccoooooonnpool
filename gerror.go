package gpool

import (
	"errors"
)

var (
	ErrInvalidPoolSetting = errors.New("invalid pool setting")
	ErrPoolAlreadyClosed  = errors.New("pool is closed")
)
