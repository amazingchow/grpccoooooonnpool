package grpccoooooonnpool

import (
	"fmt"
	"sync/atomic"

	"github.com/pkg/errors"

	boundedq "github.com/amazingchow/grpccoooooonnpool/internal/bounded-queue"
)

type GrpcClientConnPool struct {
	// Used to hold the working logical connections,
	// logical connections = physical connections * settings.MaxStreams.
	atomicLConnUsedRef int32
	// Used to hold current number of physical connections.
	atomicPConns uint32
	// Settings of the grpc connection pool.
	settings *PoolSettings
	// To store all of created logical connections.
	q *boundedq.BoundedQueue
	// To store all of created idle connections.
	conns []*GrpcConn
}

// NewGrpcClientConnPool returns a new fresh grpc client connection pool.
func NewGrpcClientConnPool(opts ...PoolSettingsOption) (p *GrpcClientConnPool, err error) {
	settings := &PoolSettings{}
	for _, opt := range opts {
		opt(settings)
	}
	if err = checkSettings(settings); err != nil {
		return
	}

	p = &GrpcClientConnPool{
		atomicLConnUsedRef: 0,
		atomicPConns:       settings.MaxIdles,
		settings:           settings,
		q:                  boundedq.NewBoundedQueue(settings.MaxIdles * settings.MaxStreams),
		conns:              make([]*GrpcConn, settings.MaxIdles),
	}

	for i := uint32(0); i < settings.MaxStreams; i++ {
		for j := uint32(0); j < settings.MaxIdles; j++ {
			p.q.Push(int(j))
		}
	}
	for i := uint32(0); i < settings.MaxIdles; i++ {
		conn, err := p.settings.Dial(p.settings.Addr)
		if err != nil {
			_ = p.Release()
			return nil, errors.Wrapf(err, "failed to create grpc client connection pool")
		}
		p.conns[i] = &GrpcConn{
			index: int(i),
			conn:  conn,
			p:     p,
		}
	}

	fmt.Println("[GrpcClientConnPool] create a new fresh grpc client connection pool")
	fmt.Printf("%s\n", p.Status())
	return p, nil
}

// Status returns the current status of the pool.
func (p *GrpcClientConnPool) Status() string {
	return fmt.Sprintf("/----------------------------------------\n Grpc Server Address: %s\n Logical Conn Used Ref: %d\n Max Idles: %d\n----------------------------------------/",
		p.settings.Addr, atomic.LoadInt32(&p.atomicLConnUsedRef), p.settings.MaxIdles)
}

// PickOne returns a available connection from the pool.
// User should use GrpcConn.Close() to put the connection back to the pool.
// If waitTime <= 0, which means block waiting for PickOne.
func (p *GrpcClientConnPool) PickOne(wait bool, waitTime int64 /* in millisecs */) (*GrpcConn, error) {
	if atomic.LoadUint32(&p.atomicPConns) == 0 {
		return nil, ErrPoolAlreadyClosed
	}

	idx := p.q.Pop(wait, waitTime)
	if idx == -1 {
		return nil, ErrPoolResourceAlreadyExhausted
	}
	p.incrRef()
	return p.conns[idx], nil
}

// Release releases all pool resources.
func (p *GrpcClientConnPool) Release() error {
	atomic.StoreInt32(&p.atomicLConnUsedRef, 0)
	atomic.StoreUint32(&p.atomicPConns, 0)

	for i := 0; i < len(p.conns); i++ {
		if p.conns[i] != nil {
			p.conns[i].Recycle()
		}
	}

	fmt.Println("[GrpcClientConnPool] release grpc client connection pool")
	fmt.Printf("%s\n", p.Status())
	return nil
}
