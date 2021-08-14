package gpool

import (
	"sync/atomic"
)

func checkSettings(addr string, opts PoolOptions) error {
	if addr == "" {
		return &PoolError{Msg: "empty address", Err: ErrInvalidPoolSetting}
	}
	if opts.Dial == nil {
		return &PoolError{Msg: "empty PoolOptions.Dial", Err: ErrInvalidPoolSetting}
	}
	if opts.MaxIdles <= 0 {
		return &PoolError{Msg: "check opts.MaxIdles", Err: ErrInvalidPoolSetting}
	}
	if opts.MaxConcurrentStreams <= 0 {
		return &PoolError{Msg: "check opts.MaxConcurrentStreams", Err: ErrInvalidPoolSetting}
	}
	return nil
}

func (p *GrpcConnPool) incrRef() {
	atomic.AddInt32(&p.atomicLConnUsedRef, 1)
}

func (p *GrpcConnPool) decrRef() {
	atomic.AddInt32(&p.atomicLConnUsedRef, -1)
}
