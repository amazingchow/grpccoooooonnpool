package gpool

import (
	"fmt"
	"math"
)

func checkSettings(addr string, opts PoolOptions) error {
	if addr == "" {
		return &PoolError{Msg: "empty address", Err: ErrInvalidPoolSetting}
	}
	if opts.Dial == nil {
		return &PoolError{Msg: "empty PoolOptions.Dial", Err: ErrInvalidPoolSetting}
	}
	if opts.MaxIdles <= 0 || opts.MaxActives <= 0 || opts.MaxIdles > opts.MaxActives {
		return &PoolError{Msg: "check opts.MaxIdles and opts.MaxActives", Err: ErrInvalidPoolSetting}
	}
	if opts.MaxConcurrentStreams <= 0 {
		return &PoolError{Msg: "check opts.MaxConcurrentStreams", Err: ErrInvalidPoolSetting}
	}
	return nil
}

func (p *GrpcConnPool) incrRef() int32 {
	p.lConnUsedRef++
	if p.lConnUsedRef == math.MaxInt32 {
		panic(ErrRefUpOverflow.Error())
	}
	return p.lConnUsedRef
}

func (p *GrpcConnPool) decrRef() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.lConnUsedRef--
	if p.lConnUsedRef < 0 {
		panic(ErrRefDownOverflow.Error())
	}
	if p.lConnUsedRef == 0 && p.pCurrConns > p.pOpts.MaxIdles {
		/* When the pool to be zero-ref and the number of physical grpc connections in the pool
		   exceeds the MaxIdles, some connections should be closed to keep the number of physical
		   grpc connections stays at the MaxIdles */
		p.pCurrConns = p.pOpts.MaxIdles
		p.shrink(p.pOpts.MaxIdles)
		fmt.Printf("[GrpcConnPool] shrink pool: %d --> %d | decrement: %d | max idles: %d | max actives: %d\n",
			p.pCurrConns, p.pOpts.MaxIdles, p.pCurrConns-p.pOpts.MaxIdles, p.pOpts.MaxIdles, p.pOpts.MaxActives)
	}
}

func (p *GrpcConnPool) shrink(begin int32) {
	for i := begin; i < p.pOpts.MaxActives; i++ {
		p.recycleResouce(i)
	}
}

func (p *GrpcConnPool) recycleResouce(index int32) {
	conn := p.conns[index]
	if conn != nil {
		_ = conn.recycle()
		p.conns[index] = nil
	}
}
