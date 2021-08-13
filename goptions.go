package gpool

import (
	"google.golang.org/grpc"
)

// PoolOptions includes pool settings for pool initialization.
type PoolOptions struct {
	/* an app-supplied-function for connection creation and configuration */
	Dial func(address string) (*grpc.ClientConn, error)
	/* maximum number of idle connections in the pool */
	MaxIdles int32
	/* maximum number of connections can be allocated by the pool at a given time */
	MaxActives int32
	/* maximum number of concurrent streams attached to a single TCP connection */
	MaxConcurrentStreams int32
	/* if set it to be true and the pool just reaches the MaxActives limitation,
	   then pool will wait for connection resycle,
	   if set it to be false and the pool just reaches the MaxActives limitation,
	   then pool will create a disposable connection */
	Reuse bool
}
