package grpccoooooonnpool

import (
	"google.golang.org/grpc"
)

// GrpcConn encapsulates the grpc.ClientConn.
type GrpcConn struct {
	index int
	conn  *grpc.ClientConn
	p     *GrpcClientConnPool
}

// Underlay returns the actual grpc connection.
func (c *GrpcConn) Underlay() *grpc.ClientConn {
	return c.conn
}

// Close decreases the reference of grpc connection if pool not full
// or just close the underlay TCP connection.
func (c *GrpcConn) Close() {
	c.p.decrRef()
	c.p.q.Push(c.index)
}

func (c *GrpcConn) Recycle() {
	conn := c.conn
	c.conn = nil
	c.p = nil
	if conn != nil {
		_ = conn.Close()
	}
}
