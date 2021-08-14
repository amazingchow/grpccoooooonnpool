package gpool

import (
	"google.golang.org/grpc"
)

// PoolOptions includes pool settings for pool initialization.
type PoolOptions struct {
	/* an app-supplied-function for connection creation and configuration */
	Dial func(address string) (*grpc.ClientConn, error)
	/* maximum number of idle connections in the pool */
	MaxIdles uint32
	// /* maximum number of connections can be allocated by the pool at a given time */
	// MaxActives uint32
	/* maximum number of concurrent streams attached to a single TCP connection */
	MaxConcurrentStreams uint32
}
