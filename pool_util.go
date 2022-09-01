package grpccoooooonnpool

import (
	"sync/atomic"

	"github.com/pkg/errors"
)

func checkSettings(settings *PoolSettings) error {
	if len(settings.Addr) == 0 {
		return errors.Wrapf(ErrInvalidPoolSettings, "empty PoolSettings.Addr")
	}
	if settings.Dial == nil {
		return errors.Wrapf(ErrInvalidPoolSettings, "empty PoolSettings.Dial")
	}
	if settings.MaxIdles <= 0 {
		return errors.Wrapf(ErrInvalidPoolSettings, "check PoolSettings.MaxIdles")
	}
	if settings.MaxIdles&(settings.MaxIdles-1) != 0 {
		return errors.Wrapf(ErrInvalidPoolSettings, "check PoolSettings.MaxIdles, should be the power of 2")
	}
	if settings.MaxStreams <= 0 {
		return errors.Wrapf(ErrInvalidPoolSettings, "check PoolSettings.MaxStreams")
	}
	if settings.MaxStreams&(settings.MaxStreams-1) != 0 {
		return errors.Wrapf(ErrInvalidPoolSettings, "check PoolSettings.MaxStreams, should be the power of 2")
	}
	return nil
}

func (p *GrpcClientConnPool) incrRef() {
	atomic.AddInt32(&p.atomicLConnUsedRef, 1)
}

func (p *GrpcClientConnPool) decrRef() {
	atomic.AddInt32(&p.atomicLConnUsedRef, -1)
}
