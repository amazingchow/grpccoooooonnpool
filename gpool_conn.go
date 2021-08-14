package gpool

import (
	"google.golang.org/grpc"
)

// GrpcConn encapsulates the grpc.ClientConn.
type GrpcConn struct {
	conn *grpc.ClientConn
	pool *GrpcConnPool
	once bool
}

// Underlay returns the actual grpc connection.
func (c *GrpcConn) Underlay() *grpc.ClientConn {
	return c.conn
}

// Close decreases the reference of grpc connection if pool not full
// or just close the underlay TCP connection.
func (c *GrpcConn) Close() error {
	c.pool.decrRef()
	if c.once {
		return c.recycle()
	}
	return nil
}

// ForceClose decreases the reference of grpc connection and close the
// underlay TCP connection.
// If user fetch a invalid grpc connection (timeout, server restart , etc...),
// it's user's responsibility to use ForceClose() but not Close().

// Deprecated: since grpc-go client connection has implemented the reconnect retry policy.
/*
https://github.com/grpc/grpc-go/blob/master/clientconn.go

_________________________________________________________
func (ac *addrConn) connect() error {
	ac.mu.Lock()
	if ac.state == connectivity.Shutdown {
		ac.mu.Unlock()
		return errConnClosing
	}
	if ac.state != connectivity.Idle {
		ac.mu.Unlock()
		return nil
	}
	// Update connectivity state within the lock to prevent subsequent or
	// concurrent calls from resetting the transport more than once.
	ac.updateConnectivityState(connectivity.Connecting, nil)
	ac.mu.Unlock()

	ac.resetTransport()
	return nil
}

...

_________________________________________________________
*/
func (c *GrpcConn) ForceClose() error {
	c.pool.decrRef()
	return c.recycle()
}

func (c *GrpcConn) recycle() error {
	conn := c.conn
	c.conn = nil
	c.pool = nil
	c.once = false
	if conn != nil {
		return conn.Close()
	}
	return nil
}
