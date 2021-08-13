package gpool

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidPoolSetting = errors.New("invalid pool setting")
	ErrCreatePool         = errors.New("failed to create pool")
	ErrPoolAlreadyClosed  = errors.New("pool is closed")
	ErrRefUpOverflow      = errors.New("numeric overflow happens: ↑")
	ErrRefDownOverflow    = errors.New("numeric overflow happens: ↓")
)

type PoolError struct {
	Msg string
	Err error
}

func (e *PoolError) Error() string {
	return fmt.Sprintf("<err: %s | msg: %s>", e.Err.Error(), e.Msg)
}
