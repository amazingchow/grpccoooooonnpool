package gpool

import (
	"context"
	"fmt"
	"os"
	"sync"

	"google.golang.org/grpc"
)

type GrpcConnPool struct {
	mu sync.RWMutex

	/* used to get logical connection */
	// TODO: use atomic val
	lConnIndex int32
	/* used to get the using logical connections
	   logical connections = physical connections * MaxConcurrentStreams */
	// TODO: use atomic val
	lConnUsedRef int32
	/* used to get current number of physical connections */
	// TODO: use atomic val
	pCurrConns int32

	/* options of the grpc connection pool */
	pOpts PoolOptions
	/* all of created physical connections */
	// TODO: clean nil conn
	// TODO: use BoundedQueue to manage conns and remove sync.RWMutex
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
		lConnIndex:   0,
		lConnUsedRef: 0,
		pOpts:        opts,
		conns:        make([]*GrpcConn, opts.MaxActives),
		addr:         addr,
	}

	for i := int32(0); i < p.pOpts.MaxIdles; i++ {
		conn, err := p.pOpts.Dial(p.addr)
		if err != nil {
			_ = p.Release()
			return nil, &PoolError{Msg: err.Error(), Err: ErrCreatePool}
		}
		p.conns[i] = &GrpcConn{
			conn: conn,
			pool: p,
			once: false,
		}
	}
	p.pCurrConns = opts.MaxIdles

	fmt.Printf("[GrpcConnPool] create a fresh grpc connection pool: %s\n", p.Status())
	return p, nil
}

// Status returns the current status of the pool.
func (p *GrpcConnPool) Status() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return fmt.Sprintf("\n<\n\tServer Address: %s\n\tLogical Conn Used Ref: %d\n\tMax Idles: %d\n\tMax Actives: %d\n>",
		p.addr, p.lConnUsedRef, p.pOpts.MaxIdles, p.pOpts.MaxActives)
}

// Get returns a available connection from the pool.
// User should use GrpcConn.Close() to put the connection back to the pool.
func (p *GrpcConnPool) Get(ctx context.Context) (*GrpcConn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.pCurrConns == 0 {
		return nil, ErrPoolAlreadyClosed
	}

	next := p.incrRef()
	if next <= p.pCurrConns*p.pOpts.MaxConcurrentStreams {
		// use round-robin to load balance
		p.lConnIndex++
		return p.conns[p.lConnIndex%p.pCurrConns], nil
	}

	// there is no available logical connection ...

	if p.pCurrConns == p.pOpts.MaxActives {
		if p.pOpts.Reuse {
			// TODO：wait for others to return grpc connection
			p.lConnIndex++
			return p.conns[p.lConnIndex%p.pCurrConns], nil
		}
		conn, err := p.pOpts.Dial(p.addr)
		return &GrpcConn{
			conn: conn,
			pool: p,
			once: true,
		}, err
	}

	increment := p.pCurrConns / 2
	if p.pCurrConns+increment > p.pOpts.MaxActives {
		increment = p.pOpts.MaxActives - p.pCurrConns
	}

	var i int32
	var err error
	var conn *grpc.ClientConn
	for i = 0; i < increment; i++ {
		conn, err = p.pOpts.Dial(p.addr)
		if err != nil {
			break
		}
		p.recycleResouce(p.pCurrConns + i)
		p.conns[p.pCurrConns+i] = &GrpcConn{
			conn: conn,
			pool: p,
			once: false,
		}
	}
	if i == 0 && err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch conn from pool, err: %v\n", err)
		return nil, err
	}

	fmt.Printf("[GrpcConnPool] grow pool: %d --> %d | increment: %d | max idles: %d | max actives: %d\n",
		p.pCurrConns, p.pCurrConns+i, i, p.pOpts.MaxIdles, p.pOpts.MaxActives)
	p.pCurrConns += i

	p.lConnIndex++
	return p.conns[p.lConnIndex%p.pCurrConns], nil
}

// Release releases all pool resource.
// After Release() the pool is no longer available.
func (p *GrpcConnPool) Release() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.lConnIndex = 0
	p.lConnUsedRef = 0
	p.pCurrConns = 0

	p.shrink(0)
	fmt.Printf("[GrpcConnPool] release grpc connection pool: %s\n", p.Status())
	return nil
}
