package gpool

import (
	"fmt"
	"sync/atomic"

	"github.com/pkg/errors"
)

type GrpcConnPool struct {
	/* used to get the using logical connections
	   logical connections = physical connections * MaxConcurrentStreams */
	atomicLConnUsedRef int32
	/* used to get current number of physical connections */
	atomicPCurrConns uint32

	/* options of the grpc connection pool */
	pOpts PoolOptions
	/* all of created physical connections */
	q     *BoundedQueue
	conns []*GrpcConn
	/* the server address to create connection */
	addr string
}

// NewGrpcConnPool returns a fresh grpc connection pool.
func NewGrpcConnPool(addr string, opts PoolOptions) (*GrpcConnPool, error) {
	if err := checkSettings(addr, opts); err != nil {
		return nil, err
	}

	p := &GrpcConnPool{
		atomicLConnUsedRef: 0,
		atomicPCurrConns:   opts.MaxIdles,
		pOpts:              opts,
		q:                  NewBoundedQueue(opts.MaxIdles * opts.MaxConcurrentStreams),
		conns:              make([]*GrpcConn, opts.MaxIdles),
		addr:               addr,
	}

	for i := uint32(0); i < opts.MaxConcurrentStreams; i++ {
		for j := uint32(0); j < opts.MaxIdles; j++ {
			p.q.Push(int(j))
		}
	}
	for i := uint32(0); i < opts.MaxIdles; i++ {
		conn, err := p.pOpts.Dial(p.addr)
		if err != nil {
			_ = p.Release()
			return nil, errors.Wrapf(err, "failed to create pool")
		}
		p.conns[i] = &GrpcConn{
			conn:  conn,
			pool:  p,
			index: int(i),
		}
	}

	fmt.Printf("[GrpcConnPool] create a fresh grpc connection pool: %s\n", p.Status())
	return p, nil
}

// Status returns the current status of the pool.
func (p *GrpcConnPool) Status() string {
	return fmt.Sprintf("\n<\n\tServer Address: %s\n\tLogical Conn Used Ref: %d\n\tMax Idles: %d\n>",
		p.addr, atomic.LoadInt32(&p.atomicLConnUsedRef), p.pOpts.MaxIdles)
}

// PickOne returns a available connection from the pool.
// User should use GrpcConn.Close() to put the connection back to the pool.
func (p *GrpcConnPool) PickOne(wait bool) (*GrpcConn, error) {
	if atomic.LoadUint32(&p.atomicPCurrConns) == 0 {
		return nil, ErrPoolAlreadyClosed
	}

	idx := p.q.Pop(wait)
	if idx == -1 {
		return nil, ErrPoolResourceAlreadyExhausted
	}
	p.incrRef()
	return p.conns[idx], nil
}

// Release releases all pool resource.
// After Release() the pool is no longer available.
func (p *GrpcConnPool) Release() error {
	atomic.StoreInt32(&p.atomicLConnUsedRef, 0)
	atomic.StoreUint32(&p.atomicPCurrConns, 0)

	for i := 0; i < len(p.conns); i++ {
		if p.conns[i] != nil {
			_ = p.conns[i].recycle()
		}
	}

	fmt.Printf("[GrpcConnPool] release grpc connection pool: %s\n", p.Status())
	return nil
}
