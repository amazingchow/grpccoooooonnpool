package gpool

import (
	"sync/atomic"

	"github.com/pkg/errors"
)

func checkSettings(addr string, opts PoolOptions) error {
	if addr == "" {
		return errors.Wrapf(ErrInvalidPoolSetting, "empty address")
	}
	if opts.Dial == nil {
		return errors.Wrapf(ErrInvalidPoolSetting, "empty PoolOptions.Dial")
	}
	if opts.MaxIdles <= 0 {
		return errors.Wrapf(ErrInvalidPoolSetting, "check opts.MaxIdles")
	}
	if opts.MaxConcurrentStreams <= 0 {
		return errors.Wrapf(ErrInvalidPoolSetting, "check opts.MaxConcurrentStreams")
	}
	return nil
}

func (p *GrpcConnPool) incrRef() {
	atomic.AddInt32(&p.atomicLConnUsedRef, 1)
}

func (p *GrpcConnPool) decrRef() {
	atomic.AddInt32(&p.atomicLConnUsedRef, -1)
}
