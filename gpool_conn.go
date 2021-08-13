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

// Close decreases the reference of grpc connection if pool not full or just close it.
func (c *GrpcConn) Close() error {
	c.pool.decrRef()
	if c.once {
		return c.recycle()
	}
	return nil
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
