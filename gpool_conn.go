package gpool

import (
	"google.golang.org/grpc"
)

// GrpcConn encapsulates the grpc.ClientConn.
type GrpcConn struct {
	conn  *grpc.ClientConn
	pool  *GrpcConnPool
	index uint32
}

// Underlay returns the actual grpc connection.
func (c *GrpcConn) Underlay() *grpc.ClientConn {
	return c.conn
}

// Close decreases the reference of grpc connection if pool not full
// or just close the underlay TCP connection.
func (c *GrpcConn) Close() error {
	c.pool.decrRef()
	c.pool.q.Push(c.index)
	return nil
}

func (c *GrpcConn) recycle() error {
	conn := c.conn
	c.conn = nil
	c.pool = nil
	if conn != nil {
		return conn.Close()
	}
	return nil
}
